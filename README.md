# Pink Voice

Voice transcription daemon with global hotkey. Press Ctrl+Q to record, release to transcribe. Result goes to clipboard.

## Requirements

- [`pink-transcriber`](https://github.com/pink-tools/pink-transcriber) running locally or accessible

## Install

```bash
git clone https://github.com/pink-tools/pink-voice.git
cd pink-voice
cp .env.example .env
# Edit .env if needed
go build .
sudo ln -sf $(pwd)/pink-voice /usr/local/bin/pink-voice
```

## Configuration

Optional `.env` file:

```
TRANSCRIPTION_PREFIX=[VOICE INPUT]
```

## Usage

```bash
pink-voice          # Start daemon (Ctrl+Q to record)
pink-voice --health # Check if running
```
