package transcriber

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pink-tools/pink-core"
)

func Transcribe(audioPath string) (string, error) {
	cmd := exec.Command(core.BinaryPath("pink-transcriber"), "transcribe", audioPath)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("transcribe: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("transcribe: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
