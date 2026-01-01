//go:build !windows

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func HealthCheck() error {
	myPID := os.Getpid()

	output, err := exec.Command("pgrep", "-x", "pink-voice").Output()
	if err != nil {
		return fmt.Errorf("not running")
	}

	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if pid, _ := strconv.Atoi(strings.TrimSpace(line)); pid > 0 && pid != myPID {
			return nil
		}
	}

	return fmt.Errorf("not running")
}

func KillExisting() {
	myPID := os.Getpid()

	output, err := exec.Command("pgrep", "-x", "pink-voice").Output()
	if err != nil {
		return
	}

	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if pid, _ := strconv.Atoi(strings.TrimSpace(line)); pid > 0 && pid != myPID {
			syscall.Kill(pid, syscall.SIGTERM)
		}
	}
}
