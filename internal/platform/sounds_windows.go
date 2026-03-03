//go:build windows

package platform

import "path/filepath"

func SoundOptions() []string {
	matches, _ := filepath.Glob(`C:\Windows\Media\*.wav`)
	return matches
}

func DefaultSounds() (start, stop, done string) {
	return `C:\Windows\Media\Speech On.wav`,
		`C:\Windows\Media\Speech Sleep.wav`,
		`C:\Windows\Media\Speech Disambiguation.wav`
}
