//go:build !windows

package platform

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func PlaySound(path string, volume float64) {
	cmd := exec.Command("afplay", "--volume", fmt.Sprintf("%.2f", volume), path)
	cmd.Start()
	go cmd.Wait()
}

func SoundOptions() []string {
	matches, _ := filepath.Glob("/System/Library/Sounds/*.aiff")
	return matches
}

func DefaultSounds() (start, stop, done string) {
	return "/System/Library/Sounds/Ping.aiff",
		"/System/Library/Sounds/Ping.aiff",
		"/System/Library/Sounds/Glass.aiff"
}
