package handlers

import (
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Input struct {
	Code string `json:"code" form:"code"`
}

type Result struct {
	Result string `json:"result"`
}

// @description Run some arbitrary nomsu code
// @success 200 {string} string "nomsu execution result"
// @router /test [post]
// @param code formData string true "Nomsu code"
// @produce json
func TestNomsu(c echo.Context) error {
	cmd := exec.Command(viper.GetString("NOMSU"), "-") // @todo add config to specify nomsu location
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	input := new(Input)
	if err = c.Bind(input); err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input.Code)
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
