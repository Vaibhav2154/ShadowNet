#!/bin/bash
# Build ShadowNet binaries for distribution

set -e

echo "=== Building ShadowNet Binaries ==="
echo ""

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="bin/release"

mkdir -p $BUILD_DIR

# Build for Linux (most common)
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/shadownet-controlplane-linux-amd64 ./cmd/controlplane
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/shadownet-node-linux-amd64 ./cmd/node

# Build for Linux ARM64 (Raspberry Pi, etc.)
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $BUILD_DIR/shadownet-node-linux-arm64 ./cmd/node

# Build for macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/shadownet-node-macos-amd64 ./cmd/node
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $BUILD_DIR/shadownet-node-macos-arm64 ./cmd/node

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/shadownet-node-windows-amd64.exe ./cmd/node

echo ""
echo "✅ Build complete!"
echo ""
echo "Binaries created in $BUILD_DIR:"
ls -lh $BUILD_DIR/
echo ""
echo "To distribute:"
echo "  1. Upload to GitHub releases"
echo "  2. Host on a web server"
echo "  3. Copy directly to target machines"
