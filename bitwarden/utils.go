package bitwarden

import (
	"math/rand"
	"os/exec"
	"time"
)

func RunCommand(commandName string, args ...string) (string, error) {
	shellCmd := exec.Command(commandName, args...)
	out, err := shellCmd.CombinedOutput()

	return string(out), err
}

// RandSleep Util to randomly sleep before calling the BW API, as an attempt to avoid getting rate-limited
func RandSleep(maxSeconds int) {
	time.Sleep(time.Duration(rand.Intn(maxSeconds)) * time.Second)
}
