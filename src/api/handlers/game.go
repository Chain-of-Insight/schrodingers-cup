package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// @description Receive and tabulate votes, stage vote outcome to be processed when game window is settled
// @router /game/vote [post]
// @param Authorization header string true "Bearer token"
// @produce json
func CastVote(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	tzid := claims["tzid"].(string)

	// Current round
	currentRound, err := strconv.Atoi(os.Getenv("CURRENT_ROUND"))
	if (currentRound == 0 || err != nil) {
		currentRound = 1
	}
	
	return c.JSON(http.StatusOK, map[string]string{
		"tzid": tzid,
	})
}

// @description Settle game window
// @router /game/settle [post]
// @param Authorization header string true "Bearer token"
func SettleGame(c echo.Context) error {
	return c.String(http.StatusOK, "todo settle game")
}
