package recorder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/gen2brain/malgo"
)

type Recorder struct {
	sampleRate uint32
	ctx        *malgo.AllocatedContext
	device     *malgo.Device
	buffer     *bytes.Buffer
	mu         sync.Mutex
	recording  bool
}

func New(sampleRate int) (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("audio context init failed: %w", err)
	}

	return &Recorder{
		sampleRate: uint32(sampleRate),
		ctx:        ctx,
		buffer:     new(bytes.Buffer),
	}, nil
}

func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.recording {
		return nil
	}

	r.buffer = new(bytes.Buffer)

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = r.sampleRate
	deviceConfig.Alsa.NoMMap = 1

	onRecvFrames := func(_, pSample []byte, _ uint32) {
		r.mu.Lock()
		r.buffer.Write(pSample)
		r.mu.Unlock()
	}

	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, malgo.DeviceCallbacks{Data: onRecvFrames})
	if err != nil {
		return fmt.Errorf("capture device init failed: %w", err)
	}

	if err := device.Start(); err != nil {
		device.Uninit()
		return fmt.Errorf("capture start failed: %w", err)
	}

	r.device = device
	r.recording = true
	return nil
}

func (r *Recorder) Stop() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.recording {
		return "", nil
	}

	r.recording = false

	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
		r.device = nil
	}

	if r.buffer.Len() == 0 {
		return "", nil
	}

	tmpFile, err := os.CreateTemp("", "pink-voice-*.wav")
	if err != nil {
		return "", fmt.Errorf("temp file failed: %w", err)
	}
	defer tmpFile.Close()

	writeWAV(tmpFile, r.buffer.Bytes(), r.sampleRate)
	return tmpFile.Name(), nil
}

func (r *Recorder) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
		r.device = nil
	}

	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
		r.ctx = nil
	}
}

func writeWAV(f *os.File, data []byte, sampleRate uint32) {
	dataSize := uint32(len(data))
	f.WriteString("RIFF")
	binary.Write(f, binary.LittleEndian, uint32(36+dataSize))
	f.WriteString("WAVE")
	f.WriteString("fmt ")
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, sampleRate)
	binary.Write(f, binary.LittleEndian, sampleRate*2)
	binary.Write(f, binary.LittleEndian, uint16(2))
	binary.Write(f, binary.LittleEndian, uint16(16))
	f.WriteString("data")
	binary.Write(f, binary.LittleEndian, dataSize)
	f.Write(data)
}
