package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pink-tools/pink-otel"
	"github.com/pink-tools/pink-voice/internal/config"
	"github.com/pink-tools/pink-voice/internal/daemon"
)

var version = "dev"

func main() {
	otel.Init("pink-voice", version)

	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-V" {
			fmt.Printf("pink-voice v%s\n", version)
			return
		}

		if os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help" {
			printUsage()
			return
		}

		switch os.Args[1] {
		case "stop":
			daemon.KillExisting()
			otel.Info(context.Background(), "stopped")
			return
		case "status":
			if err := daemon.HealthCheck(); err != nil {
				fmt.Println(`{"service":"pink-voice","status":"stopped"}`)
				os.Exit(1)
			}
			fmt.Println(`{"service":"pink-voice","status":"running"}`)
			return
		default:
			fmt.Printf("unknown command: %s\n", os.Args[1])
			printUsage()
			os.Exit(1)
		}
	}

	daemon.KillExisting()
	otel.Info(context.Background(), "starting", map[string]any{"version": version})

	cfg := config.Load()

	d, err := daemon.New(cfg)
	if err != nil {
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		d.Stop()
		os.Exit(0)
	}()

	d.Run()
}

func printUsage() {
	fmt.Println(`pink-voice - voice transcription daemon

Usage:
  pink-voice        Start daemon (Ctrl+Q hotkey)
  pink-voice stop   Stop daemon
  pink-voice status Check if daemon is running
  pink-voice help   Show this help`)
}
