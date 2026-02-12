//go:build !windows

package main

import (
	"os"
	"syscall"
)

func silenceStderr() {
	if f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0); err == nil {
		syscall.Dup2(int(f.Fd()), 2)
		f.Close()
	}
}
