package config

import (
	"os"
	"strconv"

	"github.com/pink-tools/pink-voice/internal/platform"
)

type Config struct {
	Hotkey              string
	SoundStart          string
	SoundStop           string
	SoundDone           string
	SoundVolume         float64
	TranscriptionPrefix string
	SampleRate          int
}

func Load() *Config {
	defStart, defStop, defDone := platform.DefaultSounds()

	return &Config{
		Hotkey:              envOr("HOTKEY", "alt+q"),
		SoundStart:          envOr("SOUND_START", defStart),
		SoundStop:           envOr("SOUND_STOP", defStop),
		SoundDone:           envOr("SOUND_DONE", defDone),
		SoundVolume:         envFloat("SOUND_VOLUME", 1.0),
		TranscriptionPrefix: os.Getenv("TRANSCRIPTION_PREFIX"),
		SampleRate:          16000,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envFloat(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}
