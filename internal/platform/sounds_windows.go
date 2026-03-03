//go:build windows

package platform

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

var winmm = syscall.NewLazyDLL("winmm.dll")
var playSoundW = winmm.NewProc("PlaySoundW")

const (
	sndFilename = 0x00020000
	sndAsync    = 0x00000001
)

func PlaySound(path string, _ float64) {
	ptr, _ := syscall.UTF16PtrFromString(path)
	playSoundW.Call(uintptr(unsafe.Pointer(ptr)), 0, sndFilename|sndAsync)
}

func SoundOptions() []string {
	matches, _ := filepath.Glob(`C:\Windows\Media\*.wav`)
	return matches
}

func DefaultSounds() (start, stop, done string) {
	return `C:\Windows\Media\Speech On.wav`,
		`C:\Windows\Media\Speech Sleep.wav`,
		`C:\Windows\Media\Speech Disambiguation.wav`
}
