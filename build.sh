#!/bin/bash

VERSION=$1

# Build script for water-repo
set -e

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "Building water-repo for all platforms..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build for Linux amd64
echo "Building for Linux amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/wt-$VERSION-linux-amd64 ./cmd/wt

# Build for Termux Android aarch64 (linux/arm64)
echo "Building for Termux Android aarch64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/wt-$VERSION-termux-aarch64 ./cmd/wt

# Build for macOS arm64 (Apple Silicon)
echo "Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/wt-$VERSION-darwin-arm64 ./cmd/wt

# Build for Windows amd64
echo "Building for Windows amd64..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/wt-$VERSION-windows-amd64.exe ./cmd/wt

# Build generic binary for current platform
echo "Building generic binary..."
CGO_ENABLED=0 go build -o bin/wt ./cmd/wt

echo "Build completed successfully!"
echo "Binaries are located in the bin/ directory:"
ls -la bin/
