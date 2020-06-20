package handlers

import (
	"net/http"
	"time"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/gomodule/redigo/redis"
)

type PlayerList struct {
	Players []string `json:"players"`
}

// @description Players
// @success 200 {object} PlayerList "List of players in the current gameplay session"
// @router /players [get]
// @produce json
func Players(c echo.Context) error {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	playersListKey := "players:" + currentDay
	
	// Load existing logged in players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return err
	}

	// Build list
	var playerList []string;
	for _, s := range players {
		split := strings.Split(s, " ")
		address := split[1]
		playerList = append(playerList, address)
	}

	r := &PlayerList{Players: playerList}
	return c.JSON(http.StatusOK, r)
}