VERSION := dev-$(shell date +%Y-%m-%d_%H:%M:%S)
INSTALL_DIR := ~/pink-tools/pink-voice

build:
	go build -ldflags="-X main.version=$(VERSION)" -o pink-voice .

install: build
	cp pink-voice $(INSTALL_DIR)/pink-voice
