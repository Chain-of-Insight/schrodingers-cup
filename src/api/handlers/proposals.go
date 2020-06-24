package handlers

import (
	"net/http"
	"time"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/gomodule/redigo/redis"
)

type ProposalObject struct {
	Author    	string `json:"author"`
	Code      	string `json:"code"`
	Timestamp 	string `json:"timestamp"`
	Proposal  	string `json:"proposal"`
	Ruletype  	string `json:"type"`
	Ruleindex 	int    `json:"index"`
	Round     	int    `json:"round"`
	RulePassed	bool   `json:"passed"`
}

// @description Fetches a specific rule proposal
// @router /game/proposal/:round [get]
// @param Authorization header string true "Bearer token"
// @success 200 {object} ProposalObject "A rule proposal object"
// @produce json
func Proposal(c echo.Context) error {
	round, err := strconv.Atoi(c.Param("round"))
	if err != nil {
		r := new(ProposalObject)
		return c.JSON(http.StatusOK, r)
	}
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	proposalItemKey := "proposals:" + currentDay + ":" + strconv.Itoa(round)

	p, err := redis.Strings(conn.Do("HMGET", proposalItemKey, "author", "code", "timestamp", "proposal", "ruletype", "ruleindex", "round", "success"))
	if err != nil {
		r := new(ProposalObject)
		return c.JSON(http.StatusOK, r)
	}

	// Create empty model
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

	// Instantiate struct
	proposal.author = p[0]
	proposal.code = p[1]
	proposal.timestamp = p[2]
	proposal.proposal = p[3]
	proposal.ruletype = p[4]
	proposal.ruleindex, err = strconv.Atoi(p[5])
	if err != nil {
		r := new(ProposalObject)
		return c.JSON(http.StatusOK, r)
	}
	proposal.round, err = strconv.Atoi(p[6])
	if err != nil {
		r := new(ProposalObject)
		return c.JSON(http.StatusOK, r)
	}
	proposal.success, err = strconv.ParseBool(p[7])
	if err != nil {
		r := new(ProposalObject)
		return c.JSON(http.StatusOK, r)
	}

	// Return model / resp.
	r := &ProposalObject{
		Author: proposal.author,
		Code: proposal.code,
		Timestamp: proposal.timestamp,
		Proposal: proposal.proposal,
		Ruletype: proposal.ruletype,
		Ruleindex: proposal.ruleindex,
		Round: proposal.round,
		RulePassed: proposal.success,
	}

	return c.JSON(http.StatusOK, r)
}