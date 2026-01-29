package main

import (
	"context"
	"fmt"

	"github.com/pink-tools/pink-core"
	"github.com/pink-tools/pink-voice/internal/config"
	"github.com/pink-tools/pink-voice/internal/daemon"
)

var version = "dev"

const serviceName = "pink-voice"

func main() {
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
		// Stop any existing instance via IPC
		if core.IsRunning(serviceName) {
			core.SendStop(serviceName)
		}

		cfg := config.Load()

		d, err := daemon.New(cfg)
		if err != nil {
			return err
		}

		// Run daemon in goroutine (tray.Run blocks)
		done := make(chan struct{})
		go func() {
			d.Run()
			close(done)
		}()

		// Wait for shutdown signal
		<-ctx.Done()
		d.Stop()

		// Wait for daemon to finish
		<-done
		return nil
	})
}
