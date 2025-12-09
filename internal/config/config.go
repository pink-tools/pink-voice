package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

type Platform string

const (
	PlatformMacOS   Platform = "darwin"
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
)

type Config struct {
	Platform            Platform
	TranscriptionPrefix string
	SampleRate          int
}

func Load() *Config {
	loadEnvFile()
	return &Config{
		Platform:            Platform(runtime.GOOS),
		TranscriptionPrefix: os.Getenv("TRANSCRIPTION_PREFIX"),
		SampleRate:          16000,
	}
}

func loadEnvFile() {
	paths := []string{".env"}
	if exe, err := os.Executable(); err == nil {
		if realExe, err := filepath.EvalSymlinks(exe); err == nil {
			exe = realExe
		}
		paths = append([]string{filepath.Join(filepath.Dir(exe), ".env")}, paths...)
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			godotenv.Load(p)
			return
		}
	}
}

func (c *Config) ConvertPathForTranscribe(path string) string {
	if c.Platform == PlatformWindows {
		path = strings.ReplaceAll(path, "\\", "/")
		if len(path) >= 2 && path[1] == ':' {
			return "/mnt/" + strings.ToLower(string(path[0])) + path[2:]
		}
	}
	return path
}
