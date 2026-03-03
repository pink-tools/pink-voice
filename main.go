package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/pink-tools/pink-core"
	"github.com/pink-tools/pink-voice/internal/config"
	"github.com/pink-tools/pink-voice/internal/daemon"
	"github.com/pink-tools/pink-voice/internal/platform"
)

var version = "dev"

const serviceName = "pink-voice"

func main() {
	runtime.LockOSThread()
	core.LoadEnv(serviceName)

	cfg := core.Config{
		Name:    serviceName,
		Version: version,
		Usage: `pink-voice - voice transcription daemon

Usage:
  pink-voice          Start daemon
  pink-voice stop     Stop daemon
  pink-voice status   Check if daemon is running
  pink-voice settings Configure service parameters
  pink-voice help     Show this help`,
		Commands: map[string]core.Command{
			"stop": {
				Desc: "Stop daemon",
				Run: func(args []string) error {
					if !core.IsRunning(serviceName) {
						fmt.Println("not running")
						return nil
					}
					return core.SendStop(serviceName)
				},
			},
			"status": {
				Desc: "Check if daemon is running",
				Run: func(args []string) error {
					if core.IsRunning(serviceName) {
						fmt.Println(`{"service":"pink-voice","status":"running"}`)
					} else {
						fmt.Println(`{"service":"pink-voice","status":"stopped"}`)
					}
					return nil
				},
			},
		},
	}

	actions := []core.Action{
		{Name: "settings", Label: "Settings", Desc: "Configure service parameters"},
	}
	handlers := map[string]core.ActionHandler{
		"settings": {Describe: describeSettings, Execute: executeSettings},
	}
	core.HandleActions(&cfg, actions, handlers)

	core.Run(cfg, func(ctx context.Context) error {
		voiceCfg := config.Load()

		d, err := daemon.New(voiceCfg)
		if err != nil {
			return err
		}

		go func() {
			<-ctx.Done()
			d.Stop()
		}()

		d.Run()
		return nil
	})
}

func describeSettings() core.FormSpec {
	voiceCfg := config.Load()
	sounds := platform.SoundOptions()

	return core.FormSpec{
		Title: "Pink Voice Settings",
		Fields: []core.Field{
			{
				Name:    "HOTKEY",
				Type:    "hotkey",
				Label:   "Hotkey",
				Current: voiceCfg.Hotkey,
				Default: "alt+q",
			},
			{
				Name:    "SOUND_START",
				Type:    "sound",
				Label:   "Start recording sound",
				Current: voiceCfg.SoundStart,
				Options: sounds,
			},
			{
				Name:    "SOUND_STOP",
				Type:    "sound",
				Label:   "Stop recording sound",
				Current: voiceCfg.SoundStop,
				Options: sounds,
			},
			{
				Name:    "SOUND_DONE",
				Type:    "sound",
				Label:   "Transcription done sound",
				Current: voiceCfg.SoundDone,
				Options: sounds,
			},
			{
				Name:    "SOUND_VOLUME",
				Type:    "range",
				Label:   "Sound volume",
				Current: voiceCfg.SoundVolume,
				Default: 1.0,
				Min:     0,
				Max:     1,
				Step:    0.01,
			},
			{
				Name:    "TRANSCRIPTION_PREFIX",
				Type:    "text",
				Label:   "Transcription prefix",
				Hint:    "Text prepended to transcribed output",
				Current: voiceCfg.TranscriptionPrefix,
			},
		},
	}
}

func executeSettings(values map[string]any) error {
	env := make(map[string]string)
	for k, v := range values {
		env[k] = fmt.Sprintf("%v", v)
	}
	return core.SaveEnv(serviceName, env)
}
