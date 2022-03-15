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

func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}
