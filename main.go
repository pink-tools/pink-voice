package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pinkhairedboy/pink-voice/internal/config"
	"github.com/pinkhairedboy/pink-voice/internal/daemon"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--health":
			if err := daemon.HealthCheck(); err != nil {
				fmt.Println(`{"service":"pink-voice","status":"stopped"}`)
				os.Exit(1)
			}
			fmt.Println(`{"service":"pink-voice","status":"running"}`)
			return
		case "--help", "-h":
			fmt.Println("pink-voice - Voice transcription daemon")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  pink-voice          Start daemon")
			fmt.Println("  pink-voice --health Check status")
			return
		}
	}

	daemon.KillExisting()

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
