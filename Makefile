# REST API Summarizer Makefile

.PHONY: build clean install test release help

BINARY_NAME=restapisummarizer
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Default target
all: build

# Build for current platform
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for all platforms
build-all:
	mkdir -p dist
	# Linux
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	# macOS
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	# Windows
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .

# Install locally
install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/

# Install to user bin
install-user: build
	mkdir -p ~/bin
	cp ${BINARY_NAME} ~/bin/

# Clean build artifacts
clean:
	rm -f ${BINARY_NAME}
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Create release archives
release: build-all
	cd dist && \
	tar -czf ${BINARY_NAME}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64 && \
	tar -czf ${BINARY_NAME}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64 && \
	tar -czf ${BINARY_NAME}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64 && \
	tar -czf ${BINARY_NAME}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64 && \
	zip ${BINARY_NAME}-windows-amd64.zip ${BINARY_NAME}-windows-amd64.exe

# Create checksums
checksums: release
	cd dist && sha256sum *.tar.gz *.zip > checksums.txt

# Development build with race detection
dev:
	go build -race -o ${BINARY_NAME} .

# Quick test on current directory
demo: build
	./${BINARY_NAME} sum .

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build for current platform"
	@echo "  build-all  - Build for all platforms"
	@echo "  install    - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user - Install to ~/bin"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  release    - Create release archives"
	@echo "  checksums  - Create checksums for release"
	@echo "  dev        - Development build with race detection"
	@echo "  demo       - Quick demo on current directory"
	@echo "  help       - Show this help" 