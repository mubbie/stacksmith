# Makefile for Stacksmith

BINARY_NAME=stacksmith
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X github.com/mubbie/stacksmith/cmd.Version=${VERSION} -X github.com/mubbie/stacksmith/cmd.BuildTime=${BUILD_TIME}"

# Default build is for current OS
.PHONY: build
build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} main.go

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-macos

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 main.go

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe main.go

.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 main.go

.PHONY: generate-completions
generate-completions:
	mkdir -p build/completions
	go run main.go completion bash > build/completions/stacksmith.bash
	go run main.go completion zsh > build/completions/stacksmith.zsh
	go run main.go completion fish > build/completions/stacksmith.fish

.PHONY: clean
clean:
	rm -rf bin/