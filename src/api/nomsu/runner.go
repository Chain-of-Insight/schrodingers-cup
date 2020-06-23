package nomsu

import (
	"io"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

func RunCode(code string) ([]byte, error) {
	cmd := exec.Command(viper.GetString("NOMSU"), "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	// cleanup code
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, "\t", strings.Repeat(" ", 4))

	// pipe code to nomsu stdin
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, code)
	}()

	// return combined output or error
	out, err := cmd.CombinedOutput()
	return out, err
}

//Hardcoded to only run master because we never need to run another file
func RunMaster() ([]byte, error) {
	command := "cd rules/; nomsu master.nom"
	cmd := exec.Command("/bin/bash", "-c", command)

	// return combined output or error
	out, err := cmd.CombinedOutput()
	return out, err
}
