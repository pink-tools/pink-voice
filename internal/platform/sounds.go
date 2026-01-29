//go:build !windows

package platform

import "os/exec"

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

func PlaySound(soundType SoundType) {
	cmd := exec.Command("afplay", macOSSounds[soundType])
	cmd.Start()
	go cmd.Wait()
}
