package platform

import (
	"os/exec"
	"runtime"
)

type SoundType int

const (
	SoundStart SoundType = iota
	SoundStop
	SoundDone
)

var macOSSounds = []string{
	"/System/Library/Sounds/Ping.aiff",
	"/System/Library/Sounds/Ping.aiff",
	"/System/Library/Sounds/Glass.aiff",
}

var windowsSounds = []string{
	`C:\Windows\Media\Speech On.wav`,
	`C:\Windows\Media\Speech Sleep.wav`,
	`C:\Windows\Media\Speech Disambiguation.wav`,
}

func PlaySound(soundType SoundType) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("afplay", macOSSounds[soundType]).Start()
	case "windows":
		exec.Command("powershell", "-c", `(New-Object Media.SoundPlayer '`+windowsSounds[soundType]+`').Play()`).Start()
	}
}
