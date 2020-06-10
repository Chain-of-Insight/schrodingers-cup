package handlers

import (
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/buildkite/terminal-to-html"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type TestInput struct {
	Code string `json:"code" form:"code"`
}

type TestResult struct {
	Result     string `json:"result"`
	ResultHtml string `json:"resultHtml"`
}

// @description Run some arbitrary nomsu code
// @router /test [post]
// @success 200 {object} TestResult "Nomsu execution result"
// @param input body TestInput true "Nomsu code to run"
// @produce json
func TestNomsu(c echo.Context) error {
	cmd := exec.Command(viper.GetString("NOMSU"), "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	input := new(TestInput)
	if err = c.Bind(input); err != nil {
		return err
	}

	// cleanup input
	input.Code = strings.TrimSpace(input.Code)
	input.Code = strings.ReplaceAll(input.Code, "\t", "    ")

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input.Code)
	}()

	out, err := cmd.CombinedOutput()

	// render ANSI output from nomsu as HTML
	outHtml := string(terminal.Render(out))

	// Strip ANSI from plain text output
	outText := strings.TrimSpace(stripansi.Strip(string(out)))

	r := &TestResult{
		Result:     outText,
		ResultHtml: outHtml,
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, r)
	}

	return c.JSON(http.StatusOK, r)
}
