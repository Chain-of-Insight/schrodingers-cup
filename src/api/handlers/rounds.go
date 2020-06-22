package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/gomodule/redigo/redis"
)

type CurrentRound struct {
	Round int `json:"round"`
}

// @description Players
// @success 200 {object} CurrentRound "The current round, 0 if game has not started, or -1 if no players are online"
// @router /round [get]
// @produce json
func Round(c echo.Context) error {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	roundKey := "round:" + currentDay

	// Load current round reference
	round, err := redis.Int(conn.Do("GET", roundKey))
	if err != nil {
		round = -1
	}

	r := &CurrentRound{Round: round,}

	return c.JSON(http.StatusOK, r)
}