package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @description Ping/Health Check
// @success 200 {string} string "pong"
// @router /ping [get]
// @produce plain
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
