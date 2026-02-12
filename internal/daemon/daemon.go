package daemon

import (
	"context"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	hook "github.com/robotn/gohook"

	"github.com/pink-tools/pink-core/log"
	"github.com/pink-tools/pink-voice/internal/config"
	"github.com/pink-tools/pink-voice/internal/platform"
	"github.com/pink-tools/pink-voice/internal/recorder"
	"github.com/pink-tools/pink-voice/internal/transcriber"
	"github.com/pink-tools/pink-voice/internal/tray"
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

func New(cfg *config.Config) (*Daemon, error) {
	rec, err := recorder.New(cfg.SampleRate)
	if err != nil {
		log.Error(context.Background(), "recorder init failed", log.Attr{K: "error", V: err.Error()})
		return nil, err
	}

	d := &Daemon{cfg: cfg, recorder: rec}
	d.state.Store(int32(StateIdle))
	d.tray = tray.New(d.toggleRecording, d.Stop)

	return d, nil
}

func (d *Daemon) Run() {
	log.Info(context.Background(), "hotkey registered", log.Attr{K: "hotkey", V: "Ctrl+Q"})

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
		log.Error(context.Background(), "recording failed", log.Attr{K: "error", V: err.Error()})
		return
	}
	d.setState(StateRecording)
	platform.PlaySound(platform.SoundStart)
}

func (d *Daemon) stopRecording() {
	platform.PlaySound(platform.SoundStop)
	d.setState(StateTranscribing)

	audioPath, err := d.recorder.Stop()
	if err != nil {
		log.Error(context.Background(), "recording stop failed", log.Attr{K: "error", V: err.Error()})
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
	start := time.Now()
	text, err := transcriber.Transcribe(audioPath)
	duration := time.Since(start).Seconds()

	if err != nil {
		log.Error(context.Background(), "transcription failed", log.Attr{K: "error", V: err.Error()})
		d.setState(StateIdle)
		return
	}

	if text == "" {
		text = "[No speech detected]"
	}

	var logAttrs []log.Attr
	clipboardText := text
	if d.cfg.TranscriptionPrefix != "" {
		prefix := d.cfg.TranscriptionPrefix
		if !strings.HasSuffix(prefix, " ") {
			prefix += " "
		}
		logAttrs = append(logAttrs, log.Attr{K: "prefix", V: strings.TrimSpace(d.cfg.TranscriptionPrefix)})
		clipboardText = prefix + text
	}
	logAttrs = append(logAttrs, log.Attr{K: "text", V: text}, log.Attr{K: "chars", V: len(text)}, log.Attr{K: "seconds", V: duration})

	platform.CopyToClipboard(clipboardText)
	log.Info(context.Background(), "transcribed", logAttrs...)
	platform.PlaySound(platform.SoundDone)
	d.setState(StateIdle)
}

func (d *Daemon) setState(state State) {
	d.state.Store(int32(state))
	d.tray.SetState(tray.StateType(state))
}

func (d *Daemon) Stop() {
	log.Info(context.Background(), "stopping")
	hook.End()
	d.wg.Wait()
	d.recorder.Close()
	d.tray.Quit()
}
