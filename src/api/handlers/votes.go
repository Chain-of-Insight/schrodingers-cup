package handlers

import (
	"net/http"
	"time"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/gomodule/redigo/redis"
)

type VoteObject struct {
	Address 	string `json:"player"`
	Timestamp	string `json:"timestamp"`
	Vote		string `json:"vote"`
}

type VotesList struct {
	Votes	[]VoteObject
}

// @description Fetches votes for a specific round
// @router /game/votes/:round [get]
// @param Authorization header string true "Bearer token"
// @success 200 {object} VotesList "List of votes for a target round"
// @produce json
func GetVotes(c echo.Context) error {
	round, err := strconv.Atoi(c.Param("round"))
	if err != nil {
		return err
	}
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	votesKey := "votes:" + currentDay + ":" + strconv.Itoa(round)

	votes, err := redis.Strings(conn.Do("LRANGE", votesKey, 0, -1))
	if err != nil {
		return err
	}

	v := new(VotesList)

	// Build player and login times lists
	for _, s := range votes {
		split := strings.Split(s, " ")
		address := split[0]
		time := split[1]
		vote := split[2]
		address_trimmed := strings.TrimLeft(address, "{")
		vote_trimmed := strings.TrimRight(vote, "}")
		
		VoteObj := VoteObject{}
		
		VoteObj.Address = address_trimmed
		VoteObj.Timestamp = time

		if vote_trimmed == "true" {
			VoteObj.Vote = "YES"
		} else {
			VoteObj.Vote = "NO"
		}

		v.Votes = append(v.Votes, VoteObj)
	}

	// Return model / resp.
	r := &VotesList{Votes: v.Votes,}

	return c.JSON(http.StatusOK, r)
}