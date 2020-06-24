package handlers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
)

// @description Leaderboard
// @success 200 {object} PlayerList "List of players sorted by points"
// @router /leaderboard [get]
// @produce json
func Leaderboard(c echo.Context) error {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer conn.Close()

	playerStartPts := 10 // TODO: placeholder, for now

	currentDay := time.Now().Format("2006-01-02")
	playersListKey := "players:" + currentDay

	// Load existing logged in players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return err
	}

	pointsItem := struct {
		Player string
		Points int
	}{}
	pointsList := []struct {
		Player string
		Points int
	}{}

	// Return empty if no players exist in the game session
	if len(players) == 0 {
		return c.JSON(http.StatusOK, pointsList)
	}

	for _, playerItem := range players {
		parts := strings.Fields(playerItem)
		player := parts[1]
		pKey := player + ":points:" + currentDay
		points, err := redis.Int(conn.Do("GET", pKey))
		if err != nil {
			points = playerStartPts
		}
		pointsItem.Player = player
		pointsItem.Points = points
		pointsList = append(pointsList, pointsItem)
	}

	sort.Slice(pointsList, func(i, j int) bool {
		return pointsList[i].Points > pointsList[j].Points
	})

	return c.JSON(http.StatusOK, pointsList)
}
