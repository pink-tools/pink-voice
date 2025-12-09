package transcriber

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pinkhairedboy/pink-voice/internal/config"
)

func HealthCheck(cfg *config.Config) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if cfg.Platform == config.PlatformWindows {
		cmd = exec.CommandContext(ctx, "wsl", "bash", "-lc", "pink-transcriber --health")
	} else {
		cmd = exec.CommandContext(ctx, "pink-transcriber", "--health")
	}
	return cmd.Run() == nil
}

func Transcribe(cfg *config.Config, audioPath string) (string, error) {
	path := cfg.ConvertPathForTranscribe(audioPath)

	var cmd *exec.Cmd
	if cfg.Platform == config.PlatformWindows {
		cmd = exec.Command("wsl", "bash", "-lc", fmt.Sprintf("pink-transcriber %s", path))
	} else {
		cmd = exec.Command("pink-transcriber", path)
	}

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s", string(exitErr.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
