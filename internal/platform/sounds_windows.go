//go:build windows

package platform

import (
	"syscall"
	"unsafe"
)

type SoundType int

const (
	SoundStart SoundType = iota
	SoundStop
	SoundDone
)

var winmm = syscall.NewLazyDLL("winmm.dll")
var playSound = winmm.NewProc("PlaySoundW")

const (
	SND_FILENAME = 0x00020000
	SND_ASYNC    = 0x00000001
)

var windowsSounds = []string{
	`C:\Windows\Media\Speech On.wav`,
	`C:\Windows\Media\Speech Sleep.wav`,
	`C:\Windows\Media\Speech Disambiguation.wav`,
}

func PlaySound(soundType SoundType) {
	path := windowsSounds[soundType]
	ptr, _ := syscall.UTF16PtrFromString(path)
	playSound.Call(uintptr(unsafe.Pointer(ptr)), 0, SND_FILENAME|SND_ASYNC)
}
