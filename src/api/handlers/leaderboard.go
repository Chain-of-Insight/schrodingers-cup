package handlers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
)

type PointsItem struct {
	Player string `json:"player"`
	Points int    `json:"points"`
}

type PointsList []PointsItem

// @description Leaderboard
// @success 200 {object} PointsList "List of players sorted by points"
// @router /leaderboard [get]
// @produce json
func Leaderboard(c echo.Context) error {
	pointsList, err := getPlayerPoints()
	if err != nil {
		return err
	}

	sort.Slice(*pointsList, func(i, j int) bool {
		return (*pointsList)[i].Points > (*pointsList)[j].Points
	})

	return c.JSON(http.StatusOK, *pointsList)
}

func getPlayerPoints() (*PointsList, error) {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	playersListKey := "players:" + currentDay

	// Load existing logged in players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return nil, err
	}

	pointsList := &PointsList{}

	// Return empty if no players exist in the game session
	if len(players) == 0 {
		// TODO: More code here?
		return pointsList, nil
	}

	for _, playerItem := range players {
		parts := strings.Fields(playerItem)
		player := parts[1]
		pKey := player + ":points:" + currentDay
		points, err := redis.Int(conn.Do("GET", pKey))
		if err != nil {
			points = playerStartPts
		}

		pointsItem := PointsItem{Player: player, Points: points}
		*pointsList = append(*pointsList, pointsItem)
	}

	return pointsList, nil
}
