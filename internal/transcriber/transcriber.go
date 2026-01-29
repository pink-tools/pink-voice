package transcriber

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const tcpAddr = "localhost:7465"

func Transcribe(audioPath string) (string, error) {
	pcmData, err := readPCMFromWAV(audioPath)
	if err != nil {
		return "", fmt.Errorf("read wav: %w", err)
	}

	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		return "", fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	sizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBytes, uint32(len(pcmData)))
	if _, err := conn.Write(sizeBytes); err != nil {
		return "", fmt.Errorf("write size: %w", err)
	}

	if _, err := conn.Write(pcmData); err != nil {
		return "", fmt.Errorf("write pcm: %w", err)
	}

	respSizeBytes := make([]byte, 4)
	if _, err := io.ReadFull(conn, respSizeBytes); err != nil {
		return "", fmt.Errorf("read response size: %w", err)
	}
	respSize := binary.LittleEndian.Uint32(respSizeBytes)

	textBytes := make([]byte, respSize)
	if _, err := io.ReadFull(conn, textBytes); err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	return strings.TrimSpace(string(textBytes)), nil
}

func readPCMFromWAV(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header := make([]byte, 44)
	if _, err := io.ReadFull(f, header); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, fmt.Errorf("invalid wav format")
	}

	pcmData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return pcmData, nil
}
