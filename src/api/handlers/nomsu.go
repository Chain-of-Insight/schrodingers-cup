package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @description Run some arbitrary nomsu code
// @success 200 {string} string "nomsu execution result"
// @router /test [post]
// @param code formData string true "Nomsu code"
func TestNomsu(c echo.Context) error {
	return c.String(http.StatusOK, "todo nomsu code")
}
