package config

import "os"

type Config struct {
	TranscriptionPrefix string
	SampleRate          int
}

func Load() *Config {
	return &Config{
		TranscriptionPrefix: os.Getenv("TRANSCRIPTION_PREFIX"),
		SampleRate:          16000,
	}
}
