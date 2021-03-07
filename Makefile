.PHONY: build
build: GOOS ?= linux
build: GOARCH ?= amd64
build: BIN ?= httpapi
build: DEST ?= build/$(BIN)
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DEST) cmd/$(BIN)/*.go
