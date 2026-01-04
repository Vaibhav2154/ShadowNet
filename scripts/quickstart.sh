#!/bin/bash
# Quick start script for ShadowNet

set -e

echo "=== ShadowNet Quick Start ==="
echo ""

# Check if running as root for node
if [ "$1" == "node" ] && [ "$EUID" -ne 0 ]; then 
   echo "Error: Node runtime requires root privileges for TUN device creation"
   echo "Please run: sudo $0 node"
   exit 1
fi

# Build binaries if they don't exist
if [ ! -f "bin/controlplane" ] || [ ! -f "bin/node" ]; then
    echo "Building binaries..."
    mkdir -p bin
    go build -o bin/controlplane ./cmd/controlplane
    go build -o bin/node ./cmd/node
    echo "✓ Binaries built successfully"
    echo ""
fi

# Start control plane
if [ "$1" == "controlplane" ] || [ "$1" == "all" ]; then
    echo "Starting Control Plane..."
    mkdir -p data
    ./bin/controlplane \
        --listen :8080 \
        --db ./data/controlplane.db \
        --active-timeout 5m
fi

# Start node
if [ "$1" == "node" ]; then
    echo "Starting Node..."
    
    # Generate peer ID if not provided
    PEER_ID=${PEER_ID:-"node-$(hostname)"}
    CONTROLPLANE_URL=${CONTROLPLANE_URL:-"http://localhost:8080"}
    
    echo "Peer ID: $PEER_ID"
    echo "Control Plane: $CONTROLPLANE_URL"
    echo ""
    
    ./bin/node \
        --id "$PEER_ID" \
        --controlplane-url "$CONTROLPLANE_URL" \
        --private-key-path ./shadownet.key \
        --listen-port 51820 \
        --stun-server stun.l.google.com:19302
fi

# Start dashboard
if [ "$1" == "dashboard" ]; then
    echo "Starting Dashboard..."
    cd web
    export NEXT_PUBLIC_CONTROLPLANE_URL="http://localhost:8080"
    npm run dev
fi

# Show usage if no arguments
if [ -z "$1" ]; then
    echo "Usage: $0 [controlplane|node|dashboard|all]"
    echo ""
    echo "Examples:"
    echo "  $0 controlplane    # Start control plane server"
    echo "  sudo $0 node       # Start VPN node (requires root)"
    echo "  $0 dashboard       # Start web dashboard"
    echo ""
    echo "Environment variables for node:"
    echo "  PEER_ID            # Unique peer identifier"
    echo "  CONTROLPLANE_URL   # Control plane URL (default: http://localhost:8080)"
fi
