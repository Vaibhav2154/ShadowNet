#!/bin/bash
# ShadowNet Node Installer
# Usage: curl -sSL https://your-server.com/install-node.sh | sudo bash

set -e

# Configuration
CONTROLPLANE_URL="${CONTROLPLANE_URL:-http://localhost:8080}"
BINARY_URL="${BINARY_URL:-https://github.com/Vaibhav2154/ShadowNet/releases/latest/download/shadownet-node-linux-amd64}"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/shadownet"
SERVICE_NAME="shadownet-node"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== ShadowNet Node Installer ===${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
  echo -e "${RED}Error: This script must be run as root${NC}"
  echo "Please run: sudo $0"
  exit 1
fi

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
  x86_64)
    BINARY_URL="${BINARY_URL/linux-amd64/linux-amd64}"
    ;;
  aarch64|arm64)
    BINARY_URL="${BINARY_URL/linux-amd64/linux-arm64}"
    ;;
  *)
    echo -e "${RED}Unsupported architecture: $ARCH${NC}"
    exit 1
    ;;
esac

echo "Architecture: $ARCH"
echo "Control Plane: $CONTROLPLANE_URL"
echo "Binary URL: $BINARY_URL"
echo ""

# Download binary
echo -e "${YELLOW}Downloading ShadowNet node...${NC}"
if command -v wget &> /dev/null; then
  wget -q --show-progress "$BINARY_URL" -O /tmp/shadownet-node
elif command -v curl &> /dev/null; then
  curl -L "$BINARY_URL" -o /tmp/shadownet-node
else
  echo -e "${RED}Error: Neither wget nor curl found${NC}"
  exit 1
fi

# Install binary
echo -e "${YELLOW}Installing binary...${NC}"
mv /tmp/shadownet-node "$INSTALL_DIR/shadownet-node"
chmod +x "$INSTALL_DIR/shadownet-node"

# Create config directory
mkdir -p "$CONFIG_DIR"

# Get hostname for peer ID
PEER_ID=$(hostname)

# Create systemd service
echo -e "${YELLOW}Creating systemd service...${NC}"
cat > /etc/systemd/system/$SERVICE_NAME.service <<EOF
[Unit]
Description=ShadowNet VPN Node
After=network.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/shadownet-node \\
  --id $PEER_ID \\
  --controlplane-url $CONTROLPLANE_URL \\
  --private-key-path $CONFIG_DIR/node.key \\
  --listen-port 51820
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
systemctl daemon-reload

# Enable and start service
echo -e "${YELLOW}Starting ShadowNet node...${NC}"
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

# Wait a moment for service to start
sleep 2

# Check status
if systemctl is-active --quiet $SERVICE_NAME; then
  echo ""
  echo -e "${GREEN}✅ ShadowNet node installed and started successfully!${NC}"
  echo ""
  echo "Peer ID: $PEER_ID"
  echo "Control Plane: $CONTROLPLANE_URL"
  echo ""
  echo "Useful commands:"
  echo "  Check status:  sudo systemctl status $SERVICE_NAME"
  echo "  View logs:     sudo journalctl -u $SERVICE_NAME -f"
  echo "  Restart:       sudo systemctl restart $SERVICE_NAME"
  echo "  Stop:          sudo systemctl stop $SERVICE_NAME"
  echo "  Check WG:      sudo wg show"
  echo ""
else
  echo -e "${RED}❌ Service failed to start${NC}"
  echo "Check logs: sudo journalctl -u $SERVICE_NAME -n 50"
  exit 1
fi
