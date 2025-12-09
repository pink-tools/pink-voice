package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	hook "github.com/robotn/gohook"

	"github.com/pinkhairedboy/pink-voice/internal/config"
	"github.com/pinkhairedboy/pink-voice/internal/platform"
	"github.com/pinkhairedboy/pink-voice/internal/recorder"
	"github.com/pinkhairedboy/pink-voice/internal/transcriber"
	"github.com/pinkhairedboy/pink-voice/internal/tray"
)

type State int32

const (
	StateIdle State = iota
	StateRecording
	StateTranscribing
)

type Daemon struct {
	cfg      *config.Config
	recorder *recorder.Recorder
	tray     *tray.Tray
	state    atomic.Int32
	wg       sync.WaitGroup
}

type OTelLog struct {
	Timestamp      string            `json:"timestamp"`
	SeverityNumber int               `json:"severityNumber"`
	SeverityText   string            `json:"severityText"`
	Body           string            `json:"body"`
	Resource       map[string]string `json:"resource"`
	Attributes     map[string]string `json:"attributes,omitempty"`
}

func log(severity string, body string, attrs map[string]string) {
	sevNum := 9 // INFO
	if severity == "ERROR" {
		sevNum = 17
	}

	entry := OTelLog{
		Timestamp:      time.Now().UTC().Format(time.RFC3339Nano),
		SeverityNumber: sevNum,
		SeverityText:   severity,
		Body:           body,
		Resource:       map[string]string{"service.name": "pink-voice"},
		Attributes:     attrs,
	}

	data, _ := json.Marshal(entry)
	if severity == "ERROR" {
		fmt.Fprintln(os.Stderr, string(data))
	} else {
		fmt.Println(string(data))
	}
}

func New(cfg *config.Config) (*Daemon, error) {
	if !transcriber.HealthCheck(cfg) {
		log("ERROR", "pink-transcriber not available", nil)
		return nil, fmt.Errorf("pink-transcriber not available")
	}

	rec, err := recorder.New(cfg.SampleRate)
	if err != nil {
		log("ERROR", "recorder init failed", map[string]string{"error": err.Error()})
		return nil, err
	}

	d := &Daemon{cfg: cfg, recorder: rec}
	d.state.Store(int32(StateIdle))
	d.tray = tray.New(d.toggleRecording, d.Stop)

	return d, nil
}

func (d *Daemon) Run() {
	log("INFO", "started", nil)
	log("INFO", "hotkey registered", map[string]string{"hotkey": "Ctrl+Q"})

	hook.Register(hook.KeyDown, []string{"q", "ctrl"}, func(e hook.Event) {
		d.toggleRecording()
	})

	s := hook.Start()
	go hook.Process(s)

	d.tray.Run()
}

func (d *Daemon) toggleRecording() {
	switch State(d.state.Load()) {
	case StateIdle:
		d.startRecording()
	case StateRecording:
		d.stopRecording()
	}
}

func (d *Daemon) startRecording() {
	if err := d.recorder.Start(); err != nil {
		log("ERROR", "recording failed", map[string]string{"error": err.Error()})
		return
	}
	log("INFO", "recording started", nil)
	d.setState(StateRecording)
	platform.PlaySound(platform.SoundStart)
}

func (d *Daemon) stopRecording() {
	platform.PlaySound(platform.SoundStop)
	log("INFO", "transcribing", nil)
	d.setState(StateTranscribing)

	audioPath, err := d.recorder.Stop()
	if err != nil {
		log("ERROR", "recording stop failed", map[string]string{"error": err.Error()})
		d.setState(StateIdle)
		return
	}

	if audioPath == "" {
		d.setState(StateIdle)
		return
	}

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		defer os.Remove(audioPath)
		d.processRecording(audioPath)
	}()
}

func (d *Daemon) processRecording(audioPath string) {
	text, err := transcriber.Transcribe(d.cfg, audioPath)
	if err != nil {
		log("ERROR", "transcription failed", map[string]string{"error": err.Error()})
		d.setState(StateIdle)
		return
	}

	if text == "" {
		text = "[No speech detected]"
	}

	if d.cfg.TranscriptionPrefix != "" {
		text = d.cfg.TranscriptionPrefix + text
	}

	platform.CopyToClipboard(text)
	log("INFO", "transcribed", map[string]string{"text": text})
	platform.PlaySound(platform.SoundDone)
	d.setState(StateIdle)
}

func (d *Daemon) setState(state State) {
	d.state.Store(int32(state))
	d.tray.SetState(tray.StateType(state))
}

func (d *Daemon) Stop() {
	log("INFO", "stopping", nil)
	hook.End()
	d.wg.Wait()
	d.recorder.Close()
	d.tray.Quit()
}
