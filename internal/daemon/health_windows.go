//go:build windows

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func HealthCheck() error {
	myPID := os.Getpid()

	// Use tasklist to find pink-voice.exe processes
	output, err := exec.Command("tasklist", "/FI", "IMAGENAME eq pink-voice.exe", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return fmt.Errorf("not running")
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" || strings.Contains(line, "INFO:") {
			continue
		}

		// CSV format: "pink-voice.exe","PID","Session","Mem Usage"
		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			pidStr := strings.Trim(fields[1], "\" ")
			if pid, _ := strconv.Atoi(pidStr); pid > 0 && pid != myPID {
				return nil
			}
		}
	}

	return fmt.Errorf("not running")
}

func KillExisting() {
	myPID := os.Getpid()

	// Use tasklist to find pink-voice.exe processes
	output, err := exec.Command("tasklist", "/FI", "IMAGENAME eq pink-voice.exe", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" || strings.Contains(line, "INFO:") {
			continue
		}

		// CSV format: "pink-voice.exe","PID","Session","Mem Usage"
		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			pidStr := strings.Trim(fields[1], "\" ")
			if pid, _ := strconv.Atoi(pidStr); pid > 0 && pid != myPID {
				// Kill the process using taskkill
				exec.Command("taskkill", "/PID", pidStr, "/F").Run()
			}
		}
	}
}
