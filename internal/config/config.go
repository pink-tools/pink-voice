package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	TranscriptionPrefix string
	SampleRate          int
}

func Load() *Config {
	loadEnvFile()
	return &Config{
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
