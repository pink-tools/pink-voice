# pink-voice

Voice input daemon with global hotkey. Ctrl+Q to record, Ctrl+Q to stop. Result copied to clipboard.

**Repository:** https://github.com/pink-tools/pink-voice

**Language:** Go 1.24

## File Structure

```
pink-voice/
├── main.go                           # Entry point
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
├── README.md
├── RELEASE_NOTES.md
├── rsrc_windows_386.syso             # Windows resource (32-bit)
├── rsrc_windows_amd64.syso           # Windows resource (64-bit)
├── winres/                           # Windows resources source
│   ├── icon.ico
│   ├── icon.png
│   ├── icon16.png
│   └── winres.json
├── internal/
│   ├── config/
│   │   └── config.go                 # Configuration loading
│   ├── daemon/
│   │   ├── daemon.go                 # Main daemon logic
│   │   ├── health_unix.go            # Unix health check/kill
│   │   └── health_windows.go         # Windows health check/kill
│   ├── platform/
│   │   ├── clipboard.go              # Clipboard operations
│   │   ├── sounds.go                 # macOS/Linux sounds
│   │   └── sounds_windows.go         # Windows sounds
│   ├── recorder/
│   │   └── recorder.go               # Audio recording (malgo)
│   ├── transcriber/
│   │   └── transcriber.go            # TCP client for pink-whisper
│   └── tray/
│       ├── tray.go                   # System tray menu
│       ├── icon.ico
│       ├── icon.png
│       ├── icon_other.go
│       └── icon_windows.go
├── ai-docs/
│   └── CLAUDE.md                     # This file
└── .github/
    └── workflows/
        └── build.yml
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                           main.go                               │
│  - Parse CLI (stop, status, help)                               │
│  - Kill existing instances                                      │
│  - Initialize daemon, handle signals                            │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         daemon.Daemon                           │
│  State machine: Idle → Recording → Transcribing → Idle         │
│  Hotkey: Ctrl+Q (gohook)                                        │
└─────────────────────────────────────────────────────────────────┘
            │                   │                   │
            ▼                   ▼                   ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│ recorder.Recorder│ │   tray.Tray      │ │ transcriber      │
│  - malgo context │ │  - systray menu  │ │  - TCP client    │
│  - Capture device│ │  - Status display│ │  - WAV→PCM       │
│  - WAV output    │ │  - Toggle button │ │  - pink-whisper  │
└──────────────────┘ └──────────────────┘ └──────────────────┘
            │                   │                   │
            ▼                   ▼                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                        platform.*                               │
│  - CopyToClipboard() (atotto/clipboard)                         │
│  - PlaySound() (afplay on macOS, winmm.dll on Windows)          │
└─────────────────────────────────────────────────────────────────┘
```

### State Machine

```
                 ┌──────────────┐
                 │    Idle      │◄────────────────────────┐
                 └──────┬───────┘                         │
                        │ Ctrl+Q                          │
                        ▼                                 │
                 ┌──────────────┐                         │
                 │  Recording   │ → Play start sound      │
                 └──────┬───────┘                         │
                        │ Ctrl+Q                          │
                        ▼                                 │
                 ┌──────────────┐                         │
                 │ Transcribing │ → Play stop sound       │
                 └──────┬───────┘                         │
                        │ Complete                        │
                        ▼                                 │
                 ┌──────────────┐                         │
                 │  Clipboard   │ → Play done sound ──────┘
                 │   Updated    │
                 └──────────────┘
```

## CLI

```bash
pink-voice        # Start daemon
pink-voice stop   # Stop daemon
pink-voice status # Check status (JSON output)
pink-voice help   # Show usage
```

**Health check output:**
```json
{"service":"pink-voice","status":"running"}
{"service":"pink-voice","status":"stopped"}
```

## Configuration

Optional `.env` file (executable directory or cwd):

```
TRANSCRIPTION_PREFIX="[VOICE INPUT] "
```

## Data Flow

1. User presses **Ctrl+Q** → StateRecording, play start sound
2. Audio captured via malgo (16kHz mono PCM)
3. User presses **Ctrl+Q** → StateTranscribing, play stop sound
4. Recording saved as temporary WAV
5. PCM extracted, sent via TCP to localhost:7465
6. pink-whisper returns transcribed text
7. Prefix added, text copied to clipboard
8. Play done sound, return to StateIdle
9. Temporary WAV deleted

## Components

### config/config.go

```go
type Config struct {
    TranscriptionPrefix string  // Prepended to transcribed text
    SampleRate          int     // Always 16000 Hz
}
```

Environment loaded via `core.LoadEnv()` from `/Users/pink-tools/pink-voice/.env`.

### daemon/daemon.go

```go
type Daemon struct {
    cfg      *config.Config
    recorder *recorder.Recorder
    tray     *tray.Tray
    state    atomic.Int32  // Thread-safe state
    wg       sync.WaitGroup
}
```

**Hotkey registration:**
```go
hook.Register(hook.KeyDown, []string{"q", "ctrl"}, func(e hook.Event) {
    d.toggleRecording()
})
```

### daemon/health_unix.go

- `HealthCheck()` - `pgrep -x pink-voice`
- `KillExisting()` - `pgrep` + `syscall.Kill(SIGTERM)`

### daemon/health_windows.go

- `HealthCheck()` - `tasklist /FI "IMAGENAME eq pink-voice.exe"`
- `KillExisting()` - `taskkill /PID /F`

### recorder/recorder.go

Audio capture using miniaudio (malgo):

```go
deviceConfig.Capture.Format = malgo.FormatS16  // 16-bit signed
deviceConfig.Capture.Channels = 1               // Mono
deviceConfig.SampleRate = 16000                 // 16 kHz
```

Outputs WAV file with 44-byte header to temp directory.

### transcriber/transcriber.go

TCP client for pink-whisper (localhost:7465):

```
Request:  [4B uint32 LE: size][N bytes: PCM data]
Response: [4B uint32 LE: size][M bytes: UTF-8 text]
```

Reads PCM from WAV (skips 44-byte header).

### platform/sounds.go (macOS/Linux)

```go
var macOSSounds = []string{
    "/System/Library/Sounds/Ping.aiff",   // Start
    "/System/Library/Sounds/Ping.aiff",   // Stop
    "/System/Library/Sounds/Glass.aiff",  // Done
}
```

Uses `afplay` command, runs asynchronously.

### platform/sounds_windows.go

```go
var windowsSounds = []string{
    `C:\Windows\Media\Speech On.wav`,              // Start
    `C:\Windows\Media\Speech Sleep.wav`,           // Stop
    `C:\Windows\Media\Speech Disambiguation.wav`,  // Done
}
```

Uses `winmm.dll` PlaySoundW with SND_ASYNC.

### tray/tray.go

System tray using `getlantern/systray`:

```
[Icon] Pink Voice
├── Status: Idle
├── ─────────────
├── Start Recording (Ctrl+Q)
└── Quit
```

## Dependencies

### Direct (go.mod)

| Package | Purpose |
|---------|---------|
| github.com/atotto/clipboard | Cross-platform clipboard |
| github.com/gen2brain/malgo | miniaudio bindings (audio) |
| github.com/getlantern/systray | System tray |
| github.com/pink-tools/pink-core | CLI, IPC, .env loading |
| github.com/pink-tools/pink-otel | JSON logging |
| github.com/robotn/gohook | Global hotkey |

### External

| Dependency | Purpose |
|------------|---------|
| pink-transcriber | STT server manager |
| pink-whisper | whisper.cpp TCP server |
| ggml-large-v3.bin | Whisper model (~3GB) |

### Linux Build Dependencies

```bash
sudo apt-get install -y \
    libx11-dev libx11-xcb-dev libxtst-dev \
    libxkbcommon-dev libxkbcommon-x11-dev \
    libayatana-appindicator3-dev
```

## Build Process

**Trigger:** Manual (workflow_dispatch)

| Platform | Runner | Output |
|----------|--------|--------|
| macOS ARM64 | macos-latest | pink-voice-darwin-arm64 |
| Linux AMD64 | ubuntu-latest | pink-voice-linux-amd64 |
| Windows AMD64 | windows-latest | pink-voice-windows-amd64.exe |

```bash
go build -o pink-voice .
```

## Concurrency Model

- **Main goroutine:** systray event loop
- **Hook goroutine:** gohook keyboard events
- **Transcription goroutine:** spawned per recording
- **Sound goroutines:** non-blocking playback

**Thread safety:**
- `atomic.Int32` for state
- `sync.Mutex` for recorder buffer
- `sync.WaitGroup` for shutdown

## Platform Differences

| Feature | macOS | Linux | Windows |
|---------|-------|-------|---------|
| Sounds | afplay | afplay | winmm.dll |
| Tray icon | PNG | PNG | ICO |
| Process check | pgrep | pgrep | tasklist |
| Process kill | SIGTERM | SIGTERM | taskkill |
| Clipboard | pbcopy | xclip/xsel | WinAPI |

## Error Handling

- Recorder init failure → exit 1
- Recording start/stop failure → log, return to Idle
- Transcription failure → log, return to Idle
- Empty audio → skip transcription
- No speech detected → copy "[No speech detected]"

## Temp Files

Pattern: `/tmp/pink-voice-*.wav` (Unix) or `%TEMP%\pink-voice-*.wav` (Windows)

Deleted after transcription (defer os.Remove).

## Related Projects

- **pink-transcriber** - STT server manager (github.com/pink-tools/pink-transcriber)
- **pink-whisper** - Whisper.cpp TCP server (github.com/pink-tools/pink-whisper)
- **pink-otel** - OTEL JSON logging (github.com/pink-tools/pink-otel)
