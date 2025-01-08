#!/bin/bash

# Create a bin directory if it doesn't exist
mkdir -p bin

# Build for macOS (Intel 64-bit)
GOOS=darwin GOARCH=amd64 go build -o bin/macintl/localhook client/client.go

# # Build for macOS (ARM 64-bit, M1/M2)
GOOS=darwin GOARCH=arm64 go build -o bin/macm1/localhook client/client.go

# # Build for Linux (Intel 64-bit)
GOOS=linux GOARCH=amd64 go build -o bin/linuxintel/localhook client/client.go

# # Build for Linux (ARM 64-bit)
GOOS=linux GOARCH=arm64 go build -o bin/linuxArm/localhook client/client.go

# # Build for Windows (Intel 64-bit)
GOOS=windows GOARCH=amd64 go build -o bin/windows/localhook-windows.exe client/client.go

echo "Build complete! Check the 'bin' directory for binaries."
