# pink-voice

Voice input daemon with configurable global hotkey. Push-to-talk: hold hotkey to record, release to transcribe. Result copied to clipboard.

## Install

Download binary from [Releases](https://github.com/pink-tools/pink-voice/releases), or via pink-orchestrator:

```bash
pink-orchestrator --service-download pink-voice
```

## Requirements

- [pink-transcriber](https://github.com/pink-tools/pink-transcriber)

## Usage

```bash
pink-voice            # Start daemon (system tray)
pink-voice stop       # Stop daemon
pink-voice status     # Check status
pink-voice settings   # Configure hotkey, sounds, prefix
```

## Configuration

Optional `.env` file in `~/pink-tools/pink-voice/`:

| Variable | Default | Description |
|----------|---------|-------------|
| `HOTKEY` | `alt+q` | Global hotkey (e.g. `ctrl+shift+r`) |
| `SOUND_VOLUME` | `1.0` | Sound volume (0-3) |
| `SOUND_START` | system | Sound on recording start |
| `SOUND_STOP` | system | Sound on recording stop |
| `SOUND_DONE` | system | Sound on transcription complete |
| `TRANSCRIPTION_PREFIX` | | Text prepended to output |

## Build from Source

```bash
git clone https://github.com/pink-tools/pink-voice.git
cd pink-voice
go build .
```
