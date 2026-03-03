//go:build !windows

package platform

import "path/filepath"

func SoundOptions() []string {
	matches, _ := filepath.Glob("/System/Library/Sounds/*.aiff")
	return matches
}

func DefaultSounds() (start, stop, done string) {
	return "/System/Library/Sounds/Ping.aiff",
		"/System/Library/Sounds/Ping.aiff",
		"/System/Library/Sounds/Glass.aiff"
}
