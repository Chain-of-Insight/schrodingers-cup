package handlers

import (
	"net/http"
	"time"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/gomodule/redigo/redis"
)

type PlayerList struct {
	Players 		[]string `json:"players"`
	Turn 			string `json:"currentTurn"`
	NextTurn 		string `json:"nextTurn"`
	TurnRemaining	string `json:"nextTurnAt"`
}

// @description Players
// @success 200 {object} PlayerList "List of players in the current gameplay session"
// @router /players [get]
// @produce json
func Players(c echo.Context) error {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	playersListKey := "players:" + currentDay
	
	// Load existing logged in players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
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

	// Build turns data
	currentTurn, nextTurn, timeRemaining := getTurnsData(playerList, times)

	r := &PlayerList{
		Players: playerList, 
		Turn: currentTurn, 
		NextTurn: nextTurn, 
		TurnRemaining: timeRemaining}

	return c.JSON(http.StatusOK, r)
}

/*
TODO: FIX NOMIC TIME
*/
func getTurnsData(players []string, times []string) (string, string, string) {
	// XXX TODO: Set Turn duration from Tezos value
	turnDuration := 300

	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return "Loading database failed", "", ""
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	gameStartKey := "game:" + currentDay
	roundKey := "round:" + currentDay
	var turn string;
	var nextTurn string;
	
	// Load game start reference
	gameStart, err := redis.Int(conn.Do("GET", gameStartKey))
	if err != nil {
		return "Get gameStart failed", "", ""
	}

	// Load current round reference
	round, err := redis.Int(conn.Do("GET", roundKey))
	if err != nil {
		return "Get round failed", "", ""
	}

	totalPlayers := len(players)
	
	if round == 0 {
		turn = players[0]
		if totalPlayers > 1 {
			nextTurn = players[1]
		} else {
			nextTurn = ""
		}
		timeRemaining := ""
		return turn, nextTurn, timeRemaining
	} else {
		nextRuleProposal := (turnDuration * (round + 1)) + gameStart
		m := int64(nextRuleProposal)
		t := time.Until(time.Unix(m, 0))
		timeRemaining := t.String()
		// No wrap
		if (round < totalPlayers) {
			turn = players[round]
			nextTurn = players[round + 1]
		// Wrap last
		} else if (round == totalPlayers) {
			turn = players[0]
			if totalPlayers > 1 {
				nextTurn = players[1]
			} else {
				nextTurn = ""
			}
		// Calc. wrap
		} else {
			wrap_i := round % totalPlayers;
			if (wrap_i == 0) {
				turn = players[totalPlayers - 1]
				nextTurn = players[0]
			} else {
				turn = players[wrap_i - 1]
				nextTurn = players[wrap_i]
			}
		}

		return turn, nextTurn, timeRemaining
	}
}