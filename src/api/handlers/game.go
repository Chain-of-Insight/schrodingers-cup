package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @description Receive and tabulate votes, stage vote outcome to be processed when game window is settled
// @router /game/vote [post]
func CastVote(c echo.Context) error {
	return c.String(http.StatusOK, "todo vote")
}

// @description Settle game window
// @router /game/settle [post]
func SettleGame(c echo.Context) error {
	return c.String(http.StatusOK, "todo settle game")
}
