package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/pink-tools/pink-core"
	"github.com/pink-tools/pink-voice/internal/config"
	"github.com/pink-tools/pink-voice/internal/daemon"
)

var version = "dev"

const serviceName = "pink-voice"

func main() {
	runtime.LockOSThread()
	core.LoadEnv(serviceName)

	core.Run(core.Config{
		Name:    serviceName,
		Version: version,
		Usage: `pink-voice - voice transcription daemon

Usage:
  pink-voice        Start daemon (Ctrl+Q hotkey)
  pink-voice stop   Stop daemon
  pink-voice status Check if daemon is running
  pink-voice help   Show this help`,
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
	}, func(ctx context.Context) error {
		cfg := config.Load()

		d, err := daemon.New(cfg)
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
