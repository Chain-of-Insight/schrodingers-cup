package handlers

import (
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

type Result struct {
	Result string `json:"result"`
}

// @description Run some arbitrary nomsu code
// @success 200 {string} string "nomsu execution result"
// @router /test [post]
// @param code formData string true "Nomsu code"
// @produce json
func TestNomsu(c echo.Context) error {
	cmd := exec.Command("nomsu", "-") // @todo add config to specify nomsu location
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	code := c.FormValue("code")

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, code)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	r := &Result{
		Result: strings.TrimSpace(string(out)),
	}

	return c.JSON(http.StatusOK, r)
}
