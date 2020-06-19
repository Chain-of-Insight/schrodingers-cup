package handlers

import (
	"net/http"
	"time"
	"os"

	"github.com/labstack/echo/v4"
)

type PlayerList struct {
	Players string `json:"players"`
}

// @description Players
// @success 200 {object} PlayerList "List of players in the current gameplay session"
// @router /players [get]
// @produce json
func Players(c echo.Context) error {
	currentDay := time.Now().Format("2006-01-02")
	playerListKey := "PLAYERS_" + currentDay
	playerList := os.Getenv(playerListKey)

	r := &PlayerList{Players: playerList}
	return c.JSON(http.StatusOK, r)

	return c.String(http.StatusOK, "pong")
}
