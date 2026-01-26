# pink-voice

Voice input daemon with global hotkey. Ctrl+Q to record, Ctrl+Q to stop. Result copied to clipboard.

## Install

Download binary from [Releases](https://github.com/pink-tools/pink-voice/releases).

## Usage

```bash
pink-voice        # start daemon
pink-voice stop   # stop daemon
pink-voice status # check status
pink-voice help   # show help
```

## Requirements

- [pink-transcriber](https://github.com/pink-tools/pink-transcriber) running

## Configuration

Optional `.env` file:

```
TRANSCRIPTION_PREFIX="[VOICE INPUT] "
```
