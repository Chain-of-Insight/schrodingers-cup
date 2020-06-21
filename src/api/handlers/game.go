package handlers

import (
	"net/http"
	"time"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

// XXX TODO: Set turn duration from Tezos
const turnDuration = 300

const DELETE = "delete"
const CREATE = "create"
const UPDATE = "update"
const TRANSMUTE = "transmute"

type ProposalResult struct {
	Success bool `json:"success"`
	Round	int	`json:"round"` // Updated round value
	Message string `json:"message"` // "OK!" or error message
}

type RuleProposal struct {
	Code 			string `json:"code" form:"code"`// Nomsu code
	ProposalType	string `json:"type" form:"type"`// Update, Create, Delete, Transmute
	RuleType		string `json:"kind" form:"kind"`// Mutable / Immutable
	RuleIndex		int `json:"index" form:"index"`	// rule index of the existing rule 
													// (or -1 if creating a new rule)
}

// @description Submit a new rule proposal
// @router /game/propose [post]
// @param Authorization header string true "Bearer token"
// @success 200 {object} ProposalResult ""
// @param input body RuleProposal true "Proposed rule"
// @produce json
func SubmitProposal(c echo.Context) error {
	input := new(RuleProposal)
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	tzid := claims["tzid"].(string)

	currentDay := time.Now().Format("2006-01-02")
	roundKey := "round:" + currentDay

	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	// Current round
	round, err := redis.Int(conn.Do("GET", roundKey))
	if err != nil {
		return err
	}

	// Update to new round
	if (round == 0 || err != nil) {
		round = 1
	} else {
		round += 1
	}
	
	// Load existing logged in players
	playersListKey := "players:" + currentDay
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return err
	}

	var playerList []string;
	var times []string;

	// Return empty if no players exist in the game session
	if len(players) == 0 {
		r := &PlayerList{
			Players: playerList, 
			Turn: "", 
			NextTurn: "", 
			TurnRemaining: ""}

		return c.JSON(http.StatusOK, r)
	}

	// Build player and login times lists
	for _, s := range players {
		split := strings.Split(s, " ")
		address := split[1]
		time := split[2]
		time_trimmed := strings.TrimRight(time, "}")
		playerList = append(playerList, address)
		times = append(times, time_trimmed)
	}
	playerTurnAddress := userCan(playerList, times)

	if playerTurnAddress != tzid {
		r := &ProposalResult{
			Success: false,
			Round: round,
			Message: "Unauthorized",
		}
	
		return c.JSON(http.StatusOK, r)
	}

	// Update round storage (Redis)
	if _, err := conn.Do("SET", roundKey, round); err != nil {
		return err
	}

	// Code 			string `json:"code" form:"code"`
	// ProposalType		string `json:"type" form:"type"`
	// RuleType			string `json:"kind" form:"kind"`
	// RuleIndex		int `json:"index" form:"index"`
	
	// Create rule storage
	CreateRuleEntry(tzid, input.Code, input.ProposalType, input.RuleType, input.RuleIndex, round)

	// Update chat
	message := tzid + " proposed a rule in round " + strconv.Itoa(round)
	notification := releaseNotification(message)
	var success bool;
	var statusMsg string;
	if notification == false {
		success = false
		statusMsg = "Error releasing notification"
	} else {
		success = true
		statusMsg = "OK!"
	}

	r := &ProposalResult{
		Success: success,
		Round: round,
		Message: statusMsg,
	}

	return c.JSON(http.StatusOK, r)
}

// @description Receive and tabulate votes, stage vote outcome to be processed when game window is settled
// @router /game/vote [post]
// @param Authorization header string true "Bearer token"
// @produce json
func CastVote(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	tzid := claims["tzid"].(string)

	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	roundKey := "round:" + currentDay

	// Current round
	round, err := redis.Int(conn.Do("GET", roundKey))
	if err != nil {
		return err
	}

	if (round == 0 || err != nil) {
		round = 1
	}
	
	return c.JSON(http.StatusOK, map[string]string{
		"tzid": tzid,
	})
}

// @description Settle game window
// @router /game/settle [post]
// @param Authorization header string true "Bearer token"
// func SettleGame(c echo.Context) error {
// 	return c.String(http.StatusOK, "todo settle game")
// }

func CreateRuleEntry(author string, code string, pType string, rKind string, rIndex int, round int) {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	proposalsListKey := "proposals:" + currentDay
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	var proposal struct {
		author 		string `redis:"author"`
		code  		string `redis:"code"`
		timestamp 	string `redis:"timestamp"`
		proposal	string `redis:"proposal"`
		ruletype	string `redis:"type"`
		ruleindex	int	`redis:"index"`
		round		int `redis:"round"`
		success		bool `redis:"passed"`
	}

	proposal.author = author
	proposal.timestamp = timestamp
	proposal.proposal = pType
	proposal.round = round

	// Create New / Update Existing
	if (pType != DELETE && pType != TRANSMUTE) {
		proposal.code = code
		proposal.ruletype = "mutable" // Only mutable rules can be created / updated
		if pType == CREATE {
			proposal.ruleindex = -1
		} else {
			proposal.ruleindex = rIndex
		}
	// Delete 
	} else if (pType == DELETE) {
		proposal.ruletype = "mutable" // Only mutable rules can be created / updated
		proposal.ruleindex = rIndex
	// Transmute
	} else if (pType == TRANSMUTE) {
		proposal.author = author
		proposal.ruletype = rKind
		proposal.ruleindex = rIndex
	}

	proposal.success = false;
	
	if _, err := conn.Do("LPUSH", proposalsListKey, proposal); err != nil {
		return
	}
}

func releaseNotification(notification string) bool {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
	  return false
	}

	serviceId, ok := viper.Get("TWILIO_CHAT_SERVICE_SID").(string)
	if !ok {
		return false
	}

	channelId, ok := viper.Get("TWILIO_CHANNEL").(string)
	if !ok {
		return false
	}

	accountSid, ok := viper.Get("TWILIO_ACCOUNT_SID").(string)
	if !ok {
		return false
	}

	authToken, ok := viper.Get("TWILIO_AUTH_TOKEN").(string)
	if !ok {
		return false
	}

	body := strings.NewReader(`Body=` + notification)

	req, err := http.NewRequest("POST", "https://chat.twilio.com/v2/Services/" + serviceId + "/Channels/" + channelId + "/Messages", body)
	if err != nil {
		return false
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Et voila
	return true
}

func userCan(players []string, times []string) string {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return "Loading database failed"
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	roundKey := "round:" + currentDay
	var turn string;

	// Load current round reference
	round, err := redis.Int(conn.Do("GET", roundKey))
	if err != nil {
		return "Get round failed"
	}

	totalPlayers := len(players)
	
	if round == 0 {
		turn = players[0]
		return turn
	} else {
		// No wrap
		if (round < totalPlayers) {
			turn = players[round]
		// Wrap last
		} else if (round == totalPlayers) {
			turn = players[0]
		// Calc. wrap
		} else {
			wrap_i := round % totalPlayers;
			if (wrap_i == 0) {
				turn = players[totalPlayers - 1]
			} else {
				turn = players[wrap_i - 1]
			}
		}

		return turn
	}
}