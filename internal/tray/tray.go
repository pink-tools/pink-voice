package tray

import "github.com/getlantern/systray"

type StateType int

const (
	StateIdle StateType = iota
	StateRecording
	StateTranscribing
)

var stateNames = []string{"Idle", "Recording...", "Transcribing..."}

type Tray struct {
	statusItem *systray.MenuItem
	toggleItem *systray.MenuItem
	quitItem   *systray.MenuItem
	onToggle   func()
	onQuit     func()
}

func New(onToggle, onQuit func()) *Tray {
	return &Tray{onToggle: onToggle, onQuit: onQuit}
}

func (t *Tray) Run() {
	systray.Run(t.onReady, func() {})
}

func (t *Tray) onReady() {
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTooltip("Pink Voice")

	t.statusItem = systray.AddMenuItem("Status: Idle", "")
	t.statusItem.Disable()
	systray.AddSeparator()
	t.toggleItem = systray.AddMenuItem("Start Recording", "Ctrl+Q")
	t.quitItem = systray.AddMenuItem("Quit", "")

	go t.handleClicks()
}

func (t *Tray) handleClicks() {
	for {
		select {
		case <-t.toggleItem.ClickedCh:
			if t.onToggle != nil {
				t.onToggle()
			}
		case <-t.quitItem.ClickedCh:
			if t.onQuit != nil {
				t.onQuit()
			}
			systray.Quit()
			return
		}
	}
}

func (t *Tray) SetState(state StateType) {
	t.statusItem.SetTitle("Status: " + stateNames[state])
	switch state {
	case StateIdle:
		t.toggleItem.SetTitle("Start Recording")
		t.toggleItem.Enable()
	case StateRecording:
		t.toggleItem.SetTitle("Stop Recording")
		t.toggleItem.Enable()
	case StateTranscribing:
		t.toggleItem.SetTitle("Transcribing...")
		t.toggleItem.Disable()
	}
}

func (t *Tray) Quit() {
	systray.Quit()
}
