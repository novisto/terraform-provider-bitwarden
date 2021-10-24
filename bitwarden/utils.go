package bitwarden

import (
	"os/exec"
)

func RunCommand(commandName string, args ...string) (string, error) {
	shellCmd := exec.Command(commandName, args...)
	out, err := shellCmd.CombinedOutput()

	return string(out), err
}
