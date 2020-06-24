package handlers

import (
	"io/ioutil"
	"net/http"
	"nomsu-api/nomsu"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// XXX TODO: Set all these vars from Tezos
var turnDuration int = 300
var quorumRatio float64 = 100
var pointsToWin int = 100
var playerStartPts int = 0
var rulePassPts int = 10
var voteAgainstPts int = 10
var ruleFailedPenalty int = 10

const DELETE = "delete"
const CREATE = "create"
const UPDATE = "update"
const TRANSMUTE = "transmute"

type ProposalResult struct {
	Success bool   `json:"success"`
	Round   int    `json:"round"`   // Updated round value
	Message string `json:"message"` // "OK!" or error message
}

type RuleProposal struct {
	Code         string `json:"code" form:"code"`   // Nomsu code
	ProposalType string `json:"type" form:"type"`   // Update, Create, Delete, Transmute
	RuleType     string `json:"kind" form:"kind"`   // Mutable / Immutable
	RuleIndex    int    `json:"index" form:"index"` // rule index of the existing rule
	// (or -1 if creating a new rule)
}

type VoteResult struct {
	Success bool   `json:"success"`
	Round   int    `json:"round"`   //
	Message string `json:"message"` // "OK!" or error message
}

type Vote struct {
	Vote  bool `json:"vote"`
	Round int  `json:"round"`
}

// @description Submit a new rule proposal
// @router /game/propose [post]
// @param Authorization header string true "Bearer token"
// @success 200 {object} ProposalResult ""
// @param input body RuleProposal true "Proposed rule"
// @produce json
func SubmitProposal(c echo.Context) error {
	input := new(RuleProposal)
	if err := c.Bind(input); err != nil {
		return err
	}

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
	if round == 0 || err != nil {
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

	var playerList []string
	var times []string

	// Return empty if no players exist in the game session
	if len(players) == 0 {
		r := &ProposalResult{
			Success: false,
			Round:   round,
			Message: "Unauthorized",
		}

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
			Round:   round,
			Message: "Unauthorized",
		}

		return c.JSON(http.StatusOK, r)
	}

	// Code         string `json:"code" form:"code"`   // Nomsu code
	// ProposalType string `json:"type" form:"type"`   // Update, Create, Delete, Transmute
	// RuleType     string `json:"kind" form:"kind"`   // Mutable / Immutable
	// RuleIndex    int    `json:"index" form:"index"` // rule index of the existing rule

	// Create rule storage
	ruleWasCreated := CreateRuleEntry(tzid, input.Code, input.ProposalType, input.RuleType, input.RuleIndex, round)
	if !ruleWasCreated {
		r := &ProposalResult{
			Success: false,
			Round:   round,
			Message: "Failed to create rule entry of type " + input.RuleType,
		}
		return c.JSON(http.StatusOK, r)
	}

	// Update chat
	message := tzid + " proposed a rule in round " + strconv.Itoa(round)
	notification := releaseNotification(message)
	var success bool
	var statusMsg string
	if notification == false {
		success = false
		statusMsg = "Error releasing notification"
	} else {
		success = true
		statusMsg = "OK!"
	}

	// Update round storage (Redis)
	if _, err := conn.Do("SET", roundKey, round); err != nil {
		return err
	}

	r := &ProposalResult{
		Success: success,
		Round:   round,
		Message: statusMsg,
	}

	return c.JSON(http.StatusOK, r)
}

// @description Receive and tabulate votes, stage vote outcome to be processed when game window is settled
// @router /game/vote [post]
// @param Authorization header string true "Bearer token"
// @produce json
func CastVote(c echo.Context) error {
	input := new(Vote)
	if err := c.Bind(input); err != nil {
		return err
	}

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
	playersListKey := "players:" + currentDay

	// Current round
	round, err := redis.Int(conn.Do("GET", roundKey))
	if err != nil {
		return err
	}

	// Players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return err
	}

	// Game not started
	if round == 0 {
		r := &VoteResult{
			Success: false,
			Round:   round,
			Message: "Unauthorized",
		}
		return c.JSON(http.StatusOK, r)
		// Voting on incorrect round
	} else if input.Round != round {
		r := &VoteResult{
			Success: false,
			Round:   round,
			Message: "Unauthorized",
		}
		return c.JSON(http.StatusOK, r)
	}

	// If user can vote
	canVote := userCanVote(tzid, round)
	var success bool
	var message string
	var votetype string
	if canVote {
		if input.Vote == true {
			votetype = "YES"
		} else {
			votetype = "NO"
		}

		success = true
		message = tzid + " successfully voted " + votetype + " in round " + strconv.Itoa(round)
	} else {
		r := &VoteResult{
			Success: false,
			Round:   round,
			Message: "Unauthorized",
		}
		return c.JSON(http.StatusOK, r)
	}

	// Cast vote
	voteKey := "votes:" + currentDay + ":" + strconv.Itoa(round)
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	var vote struct {
		author    string `redis:"author"`
		timestamp string `redis:"timestamp"`
		vote      bool   `redis:"vote"`
	}

	vote.author = tzid
	vote.timestamp = timestamp
	vote.vote = input.Vote

	// Update vote
	if _, err := conn.Do("RPUSH", voteKey, vote); err != nil {
		return err
	}

	// Handle quorum vote
	votes, err := redis.Strings(conn.Do("LRANGE", voteKey, 0, -1))
	if err != nil {
		return err
	}
	totalPlayers := len(players)
	totalVotes := len(votes)
	isQuorum := isValidQuorum(totalPlayers, totalVotes)

	if isQuorum {
		// Handle round proccessing
		// XXX TODO: THIS
		processRound, respMsg := processRound(round)
		// @pogo @drew

		// State of paradox :o
		if processRound == false {
			statusMsg := "Error proccessing round. The game has reached a state of paradox!"
			paradoxMessage := message + " which has brought the game to a state of paradox"
			notification := releaseNotification(paradoxMessage)
			if notification == false {
				success = false
				statusMsg = "Error releasing notification"
			} else {
				success = true
				statusMsg = respMsg
			}
			EndNomic(tzid)
			r := &VoteResult{
				Success: false,
				Round:   round,
				Message: statusMsg,
			}
			return c.JSON(http.StatusOK, r)
		}
		// Check game over
		isGameOver := checkGameOver()
		if isGameOver {
			EndNomic(tzid)
		}

		// Release notification
		var statusMsg string
		notification := releaseNotification(message)
		if notification == false {
			statusMsg = "Error releasing notification"
		} else {
			statusMsg = respMsg
		}

		// HTTP repsonse
		r := &VoteResult{
			Success: success,
			Round:   round,
			Message: statusMsg,
		}

		return c.JSON(http.StatusOK, r)
	} else {
		respMsg := "OK!"
		// Release notification
		var statusMsg string
		notification := releaseNotification(message)
		if notification == false {
			statusMsg = "Error releasing notification"
		} else {
			statusMsg = respMsg
		}

		// HTTP repsonse
		r := &VoteResult{
			Success: success,
			Round:   round,
			Message: statusMsg,
		}

		return c.JSON(http.StatusOK, r)
	}
}

// @description Settle game window
// @router /game/settle [post]
// @param Authorization header string true "Bearer token"
// func SettleGame(c echo.Context) error {
// 	return c.String(http.StatusOK, "todo settle game")
// }

func CreateRuleEntry(author string, code string, pType string, rKind string, rIndex int, round int) bool {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return false
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	proposalItemKey := "proposals:" + currentDay + ":" + strconv.Itoa(round)
	voteKey := "votes:" + currentDay + ":" + strconv.Itoa(round)
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Check if rule proposal exists for round
	// Checking by votes is easier since their key
	// is instanced by round. This works because
	// submitting a proposal automatically casts the
	// first vote as e.g. "yes" for proposing player
	votes, err := redis.Strings(conn.Do("LRANGE", voteKey, 0, -1)) //  lrange votes:2020-06-23:1 0 -1
	if err != nil {
		return false
	}
	if len(votes) != 0 {
		return false
	}

	var proposal struct {
		author    string `redis:"author"`
		code      string `redis:"code"`
		timestamp string `redis:"timestamp"`
		proposal  string `redis:"proposal"`
		ruletype  string `redis:"type"`
		ruleindex int    `redis:"index"`
		round     int    `redis:"round"`
		success   bool   `redis:"passed"`
	}

	proposal.author = author
	proposal.timestamp = timestamp
	proposal.proposal = pType
	proposal.round = round

	var vote struct {
		author    string `redis:"author"`
		timestamp string `redis:"timestamp"`
		vote      bool   `redis:"vote"`
	}

	vote.author = author
	vote.timestamp = timestamp
	vote.vote = true

	// Create New / Update Existing
	if pType != DELETE && pType != TRANSMUTE {
		// Proposal type must be valid
		if pType != CREATE && pType != UPDATE {
			return false
		}
		proposal.code = code
		proposal.ruletype = "mutable" // Only mutable rules can be created / updated
		if pType == CREATE {
			proposal.ruleindex = -1
		} else {
			proposal.ruleindex = rIndex
		}
		// Delete
	} else if pType == DELETE {
		proposal.ruletype = "mutable" // Only mutable rules can be created / updated
		proposal.ruleindex = rIndex
		// Transmute
	} else if pType == TRANSMUTE {
		proposal.author = author
		proposal.ruletype = rKind
		proposal.ruleindex = rIndex
	}

	proposal.success = false

	// Update proposal
	if _, err := conn.Do("HSET", proposalItemKey, "author", proposal.author, "code", proposal.code, "timestamp", proposal.timestamp, "proposal", proposal.proposal, "ruletype", proposal.ruletype, "ruleindex", proposal.ruleindex, "round", proposal.round, "success", proposal.success); err != nil {
		return false
	}

	// Update vote
	if _, err := conn.Do("RPUSH", voteKey, vote); err != nil {
		return false
	}

	return true
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

	req, err := http.NewRequest("POST", "https://chat.twilio.com/v2/Services/"+serviceId+"/Channels/"+channelId+"/Messages", body)
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

// XXX TODO - Process round
func processRound(round int) (bool, string) {
	if round < 1 {
		return false, "Game has not started. Current round is " + strconv.Itoa(round)
	}

	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return false, "Error connecting to db"
	}

	var proposal struct {
		author    string `redis:"author"`
		code      string `redis:"code"`
		timestamp string `redis:"timestamp"`
		proposal  string `redis:"proposal"`
		ruletype  string `redis:"type"`
		ruleindex int    `redis:"index"`
		round     int    `redis:"round"`
		success   bool   `redis:"passed"`
	}

	// Set some defaults
	currentDay := time.Now().Format("2006-01-02")
	proposalItemKey := "proposals:" + currentDay + ":" + strconv.Itoa(round)
	roundCounterKey := "round:" + currentDay
	roundKey := "round" + currentDay + ":" + strconv.Itoa(round)
	round = round + 1

	p, err := redis.Strings(conn.Do("HMGET", proposalItemKey, "author", "code", "timestamp", "proposal", "ruletype", "ruleindex", "round", "success"))
	if err != nil {
		return false, "Error getting proposal at key " + proposalItemKey
	}
	proposal.author = p[0]
	proposal.code = p[1]
	proposal.timestamp = p[2]
	proposal.proposal = p[3]
	proposal.ruletype = p[4]
	proposal.ruleindex, err = strconv.Atoi(p[5])
	if err != nil {
		return false, "Error converting proposal.ruleindex to int"
	}
	proposal.round, err = strconv.Atoi(p[6])
	if err != nil {
		return false, "Error converting proposal.round to int"
	}
	proposal.success, err = strconv.ParseBool(p[7])
	if err != nil {
		return false, "Error converting proposal.success to bool"
	}

	proposingPlayer := proposal.author

	// 0) Determine if the vote was voted in or voted down by the quorum
	
	// Determine if passing / failing proposal
	voteKey := "votes:" + currentDay + ":" + strconv.Itoa(round)
	votes, err := redis.Strings(conn.Do("LRANGE", voteKey, 0, -1))
	if err != nil {
		return false, "Error retrieving votes from db with LRANGE for key " + voteKey
	}

	var votedYes int = 0;
	var votedNo int = 0;
	var votedAgainstPlayers []string;
	for _, s := range votes {
		split := strings.Split(s, " ")
		address := split[0]
		address_trimmed := strings.TrimLeft(address, "{")
		vote := split[2]
		vote_trimmed := strings.TrimRight(vote, "}")
		if vote_trimmed == "true" {
			votedYes++
		} else {
			votedNo++
			// Store voted against ref.
			votedAgainstPlayers = append(votedAgainstPlayers, address_trimmed)
		}
	}

	var rulePassed bool;
	if votedYes > rulePassPts {
		// Rule successfully passed
		rulePassed = true
	} else {
		rulePassed = false
	}

	// 1) Call your file system functions to change the target rule (see: proposal)
	
	if rulePassed {
		switch proposal.proposal {
		case UPDATE:
			err, _ := Update(proposal.code, proposal.ruleindex, proposal.ruletype)
			if err != nil {
				return false, "Error performing nomsu UPDATE with proposal.code " + proposal.code + "at proposal.ruleindex " + strconv.Itoa(proposal.ruleindex) + " for proposal.ruletype " + proposal.ruletype
			}
		case CREATE:
			err, _, _ := Create(proposal.code, proposal.ruletype)
			if err != nil {
				return false, "Error performing nomsu CREATE with proposal.code " + proposal.code + " for proposal.ruletype " + proposal.ruletype
			}
		case DELETE:
			err := Delete(proposal.ruleindex, proposal.ruletype)
			if err != nil {
				return false, "Error performing nomsu DELETE with proposal.ruleindex " + strconv.Itoa(proposal.ruleindex) + " for proposal.ruletype " + proposal.ruletype
			}
		case TRANSMUTE:
			err, _ := Transmute(proposal.ruleindex, proposal.ruletype)
			if err != nil {
				return false, "Error performing nomsu TRANSMUTE with proposal.ruleindex " + strconv.Itoa(proposal.ruleindex) + " for proposal.ruletype " + proposal.ruletype
			}
		}
		
		// 2) (rules loop) Run updated ruleset using master.nom
		output, err := nomsu.RunMaster()
		if err != nil {
			return false, "Error running nomsu.RunMaster()"
		}

		// 3) Use the output from master to get replacement values for vars.nom
		var ruleSet []int
		words := strings.Fields(string(output))
		for idx, word := range words {
			if word == "=" {
				val, err := strconv.Atoi(words[idx+1])
				if err != nil {
					return false, "Error converting master.nom output to integer at word index " + strconv.Itoa(idx + 1) + " with value: " + words[idx+1]
				}
				ruleSet = append(ruleSet, val)
			}
		}

		// 4) Replace vars.nom with updated values
		b, err := ioutil.ReadFile("nomsu/rules/vars.nom") // read the original file contents
		if err != nil {
			return false, "Error reading vars.nom file with ioutil.ReadFile()"
		}

		// replace the integer values
		newVars := ""
		count := 0
		words = strings.Fields(string(b))
		for _, word := range words {
			_, err := strconv.Atoi(word)
			if err != nil {
				newVars += word + " "
				continue
			}
			newVars += strconv.Itoa(ruleSet[count]) + "\n"
			count++
		}

		f, err := os.Create("nomsu/rules/vars.nom") // write the file
		if err != nil {
			return false, "Error creating file vars.nom using os.Create()"
		}
		defer f.Close()

		// write the code string
		_, err = f.WriteString(newVars)
		if err != nil {
			return false, "Error writing to file vars.nom with f.WriteString(newVars)"
		}

		// Updated vars.nom values if rule was changed
		playerStartPts = ruleSet[0]
		pointsToWin = ruleSet[1] //  <= @see checkGameover()
		rulePassPts = ruleSet[2]
		voteAgainstPts = ruleSet[3]
		ruleFailedPenalty = ruleSet[4]
		// gameWindowDuration := ruleSet[5]
		// gameWindowStartUTC := ruleSet[6]
		turnDuration = ruleSet[7]
		quorumRatio = float64(ruleSet[8])
	}

	// 5 a) (players loop) Apply point changes to each user (as necessary)
	
	// say("Player starting points = " + $bl_startPoints)
	// say("Points required to win the game = " + $bl_winPoints)
	// say("Points gained for passing a rule = " + $bl_rulePassPts)
	// say("Points gained for voting against a passed rule = " + $bl_voteAgainstPts)
	// say("Points lost for a rule rejection = " + $bl_ruleFailedPenalty)
	// say("Daily game window duration = " + $bl_gameWindowDuration)
	// say("Game start hour (UTC) = " + $bl_gameWindowHourUTC)
	// say("Player turn duration = " + $bl_turnWindowDuration)
	// say("Percentage of votes required to pass a rule = " + $bl_voteRatioRequired)
	// say("Percentage of votes required to transmutate a rule = " + $bl_transmutationVoteRatioRequired)

	// Get player refs and award points
	type PlayerRoundResult struct {
		player string `redis:"player"`
		points int `redis:"points"`
	}

	// Parse points
	// Proposing player
	// Handle points if rule passed
	pKey := proposingPlayer + ":points:" + currentDay
	points, err := redis.Int(conn.Do("GET", pKey))
	if err != nil {
		points = playerStartPts
	}
	playerUpdate := new(PlayerRoundResult)
	playerUpdate.player = proposingPlayer
	if rulePassed {
		// Points increase
		playerUpdate.points = rulePassPts
		points = points + rulePassPts
	} else {
		// Points decrease
		playerUpdate.points = 0 - ruleFailedPenalty
		points = points + (0 - ruleFailedPenalty)
	}
	// Proposing player points stored
	if _, err := conn.Do("SET", pKey, points); err != nil {
		return false, "Error setting proposing player's points in the db to " + strconv.Itoa(points)
	}
	// Update round deltas ref.
	if _, err := conn.Do("RPUSH", roundKey, playerUpdate); err != nil {
		return false, "Error pushing playerUpdate into " + roundKey + " using RPUSH for proposing player"
	}

	// Determine if player gains points
	for _, s := range votedAgainstPlayers {				
		if rulePassed {
			if voteAgainstPts > 0 {
				againstKey := s + ":points:" + currentDay
				points, err := redis.Int(conn.Do("GET", againstKey))
				if err != nil {
					points = playerStartPts
				}
				playerUpdate := new(PlayerRoundResult)
				playerUpdate.player = s
				
				// Points increase
				playerUpdate.points = voteAgainstPts
				points = points + voteAgainstPts

				// Vote against player points stored
				if _, err := conn.Do("SET", againstKey, points); err != nil {
					return false, "Error setting vote against points on " + againstKey + " for " + strconv.Itoa(points) + " points"
				}
				// Update round deltas ref.
				if _, err := conn.Do("RPUSH", roundKey, playerUpdate); err != nil {
					return false, "Error pushing playerUpdate into " + roundKey + " using RPUSH for vote against player: " + s + " for " + strconv.Itoa(points) + " points"
				}
			}
		}
	}


	// 5 b) (rounds loop) TODO: Apply changes to rules (e.g. necro rule strategy)


	// 6) Release chat notifications
	general := releaseNotification("Round " + strconv.Itoa(round) + " has concluded")
	if !general {
		return false, "Error releasing general Twilio message for round conclusion"
	}
	thePeanutGallery := strings.Join(votedAgainstPlayers, ", ")
	if rulePassed {
		gg := proposingPlayer + "'s rule has been successfully passed in round " + strconv.Itoa(round)
		m := releaseNotification(gg)
		if !m {
			return true, "Round concluded successfully but encountered an error updating chat with message : " + gg
		}
		aChallengerApproaches := proposingPlayer + " gained " + strconv.Itoa(rulePassPts) + " points"
		m2 := releaseNotification(aChallengerApproaches)
		if !m2 {
			return true, "Round concluded successfully but encountered an error updating chat with message : " + aChallengerApproaches
		}
		luckyOnesMsg := thePeanutGallery + " each gain " + strconv.Itoa(voteAgainstPts) + " points for challenging the mentality of the herd"
		m3 := releaseNotification(luckyOnesMsg)
		if !m3 {
			return true, "Round concluded successfully but encountered an error updating chat with message : " + luckyOnesMsg
		}
	} else {
		bm := proposingPlayer + "'s rule was deemed useless by a jury of their peers in round " + strconv.Itoa(round)
		m := releaseNotification(bm)
		if !m {
			return true, "Round concluded successfully but encountered an error updating chat with message : " + bm
		}
		bm2 := proposingPlayer + "loses " + strconv.Itoa(ruleFailedPenalty) + " points carelessly"
		m1 := releaseNotification(bm2)
		if !m1 {
			return true, "Round concluded successfully but encountered an error updating chat with message : " + bm2
		}
		bm3 := thePeanutGallery + " each snicker and stomp their feet in delight"
		m2 := releaseNotification(bm3)
		if !m2 {
			return true, "Round concluded successfully but encountered an error updating chat with message : " + bm3
		}
	}

	// Update round storage (Redis)
	if _, err := conn.Do("SET", roundCounterKey, round); err != nil {
		return false, "Error incrementing round storage in db"
	}

	return true, "OK!"
}

func isValidQuorum(totalPlayers int, totalVotes int) bool {
	// Shift percentage to ratio multer
	q := quorumRatio * 0.01
	var players float64 = float64(totalPlayers)
	var votes float64 = float64(totalVotes)
	var isQuorum bool
	// Minimum votes required (float64 with decimals)
	votesRequired := players * q
	if votes >= votesRequired {
		isQuorum = true
	} else {
		isQuorum = false
	}
	return isQuorum
}

func userCanVote(playerAddress string, round int) bool {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return false
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	voteKey := "votes:" + currentDay + ":" + strconv.Itoa(round)
	playersListKey := "players:" + currentDay

	// Load existing logged in players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return false
	}

	// Check if player has already logged in
	playerExists := ItemExists(players, playerAddress)
	if !playerExists {
		return false
	}

	// Load existing votes of target round
	votes, err := redis.Strings(conn.Do("LRANGE", voteKey, 0, -1))
	if err != nil {
		return false
	}

	// Check if player has already voted in this round
	voteExists := ItemExists(votes, playerAddress)
	if voteExists {
		return false
	}

	// Else, can vote
	return true
}

func checkGameOver() bool {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return false
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	playersListKey := "players:" + currentDay

	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return false
	}

	var playerList []string

	// Build player list
	for _, s := range players {
		split := strings.Split(s, " ")
		address := split[1]
		playerList = append(playerList, address)
	}

	// Check points
	var gameover bool = false
	for _, p := range playerList {
		pKey := p + ":points:" + currentDay
		points, err := redis.Int(conn.Do("GET", pKey))
		if err != nil {
			points = 0
		}
		if points >= pointsToWin {
			gameover = true
			break
		}
	}

	return gameover
}

/*
TODO: finish this...
- Update chat with game result
- Add environment variables middleware to disable the server when Nomic is won
*/
func EndNomic(player string) {
	os.Setenv("GAME_OVER", "1")
	gameoverMsg := "This game of nomic has just been concluded by " + player
	gameoverNotification := releaseNotification(gameoverMsg)
	if !gameoverNotification {
		return
	}
	ggMsg := "Congrats and GG to everyone who participated!"
	ggNotification := releaseNotification(ggMsg)
	if !ggNotification {
		return
	}
}

func ItemExists(slice []string, entry string) bool {
	for _, s := range slice {
		if strings.Contains(s, entry) {
			return true
		}
	}
	return false
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
	var turn string

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
		if round < totalPlayers {
			turn = players[round - 1]
			// Wrap last
		} else if round == totalPlayers {
			turn = players[0]
			// Calc. wrap
		} else {
			wrap_i := round % totalPlayers
			if wrap_i == 0 {
				turn = players[totalPlayers-1]
			} else {
				turn = players[wrap_i-1]
			}
		}

		return turn
	}
}

// helper function to loop through the next files to move them up
func MoveUpFiles(index int, ruleType string) error {
	for true {
		// check if file exists
		nextRuleFile := "nomsu/rules/" + ruleType + "/rule" + strconv.Itoa((index + 1)) + ".nom"
		if _, err := os.Stat(nextRuleFile); os.IsNotExist(err) {
			break
		}
		// rename it
		currRuleFile := "nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom"
		err := os.Rename(nextRuleFile, currRuleFile)
		if err != nil {
			return err
		}
		index += 1
	}
	return nil
}

func Update(code string, index int, ruleType string) (error, string) {
	// cleanup code
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, "\t", strings.Repeat(" ", 4))

	// check if code is valid Nomsu code before overwriting anything! (might be redundant)
	_, err := nomsu.RunCode(code)
	if err != nil {
		return err, ""
	}

	// check if file exists (otherwise it's not really Updating, is it?)
	if _, err := os.Stat("nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom"); os.IsNotExist(err) {
		return err, ""
	}

	// read the file
	byte_contents, err := ioutil.ReadFile("schrodinger/schrodingers-cup/src/api/nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom") // read the original file contents
	if err != nil {
		return err, ""
	}
	contents := string(byte_contents)

	idx1 := strings.Index(contents, "external: $")
	idx2 := strings.Index(contents[idx1+11:], "=")
	varname := contents[idx1+11 : idx1+11+idx2]
	varname = strings.ReplaceAll(varname, " ", "")
	idx3 := strings.Index(code, varname)
	if idx3 < 0 {
		return err, ""
	}
	code = strings.ReplaceAll(code, " ", "")
	idx4 := strings.Index(code[idx3:], "=")
	newVal, err := strconv.Atoi(code[idx4+2:])
	newContents := contents[:idx1+11+idx2+1] + " " + strconv.Itoa(newVal)

	// overwrite the file
	f, err := os.Create("nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom")
	if err != nil {
		return err, ""
	}
	defer f.Close()

	// write the code string
	_, err = f.WriteString(newContents)
	if err != nil {
		return err, ""
	}

	// return success
	return nil, code
}

func Create(code string, ruleType string) (error, int, string) {
	// cleanup code
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, "\t", strings.Repeat(" ", 4))

	//check if the user added the "basic stuff a rule needs" aka 'use "vars"' and 'external: '
	//this check is very basic and will not work with complex rules (but probably won't need to either)
	idx := strings.Index(code, "external: ")
	if idx < 0 {
		code = "external: " + code
	}

	// check if code is valid Nomsu code before overwriting anything! (might be redundant)
	_, err = nomsu.RunCode(code)
	if err != nil {
		return err, -1, ""
	}

	//we can't run the code in limbo with this, so we're adding it after
	idx = strings.Index(code, "use \"vars\"")
	if idx < 0 {
		code = "use \"vars\"\n" + code
	}

	// find the first free available index
	index := 0
	for true {
		if _, err := os.Stat("nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom"); os.IsNotExist(err) {
			break
		}
		index += 1
	}

	// create the new rule file with the free index
	f, err := os.Create("nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom")
	if err != nil {
		return err, -1, ""
	}
	defer f.Close()

	// write the code string
	_, err = f.WriteString(code)
	if err != nil {
		return err, -1, ""
	}

	// return created rule index
	return nil, index, code
}

func Delete(index int, ruleType string) error {
	// check if the file exists
	if _, err := os.Stat("nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom"); os.IsNotExist(err) {
		return err
	}

	// remove the file
	err := os.Remove("nomsu/rules/" + ruleType + "/rule" + strconv.Itoa(index) + ".nom")
	if err != nil {
		return err
	}

	// loop through the next files to move them up
	MoveUpFiles(index, ruleType)

	// return nil if everything worked
	return nil
}

func Transmute(indexFrom int, ruleTypeFrom string) (error, int) {
	// check if the file exists
	fileFrom := "nomsu/rules/" + ruleTypeFrom + "/rule" + strconv.Itoa(indexFrom) + ".nom"
	if _, err := os.Stat(fileFrom); os.IsNotExist(err) {
		return err, -1
	}

	// find out what we're transmuting to
	ruleTypeTo := "mutable"
	if ruleTypeFrom == ruleTypeTo {
		ruleTypeTo = "immutable"
	}

	// find the first free available index in the new folder
	indexTo := 0
	for true {
		if _, err := os.Stat("nomsu/rules/" + ruleTypeTo + "/rule" + strconv.Itoa(indexTo) + ".nom"); os.IsNotExist(err) {
			break
		}
		indexTo += 1
	}

	// move the file from the old folder to the new one
	err := os.Rename(fileFrom, "nomsu/rules/"+ruleTypeTo+"/rule"+strconv.Itoa(indexTo)+".nom")
	if err != nil {
		return err, -1
	}

	// clean up the original folder (which lost a rule)
	MoveUpFiles(indexFrom, ruleTypeFrom)

	// return nil if everything worked
	return nil, indexTo
}
