# ShadowNet Deployment Guide

This guide shows you how to deploy ShadowNet across multiple machines for P2P mesh networking **without installing Go on every machine**.

---

## Deployment Architecture

```
┌─────────────────────┐
│  Control Plane      │  (Public Server - Cloud/VPS)
│  Port 8080          │  - Coordinates peers
│  + Dashboard 3000   │  - No data traffic
└─────────────────────┘
          │
          │ Registration & Discovery
          │
    ┌─────┴─────┬─────────────┐
    │           │             │
┌───▼───┐   ┌───▼───┐   ┌────▼────┐
│ Node 1│   │ Node 2│   │  Node 3 │
│ Home  │◄──┤ Office│◄──┤  Cloud  │
└───────┘   └───────┘   └─────────┘
     Direct P2P Encrypted Tunnels
```

---

## Option 1: Binary Distribution (Recommended - No Go Required!)

### Step 1: Build Binaries Once (On Your Dev Machine)

```bash
cd /home/vaibhi/Dev/ShadowNet

# Build for Linux (most common)
GOOS=linux GOARCH=amd64 go build -o bin/controlplane-linux-amd64 ./cmd/controlplane
GOOS=linux GOARCH=amd64 go build -o bin/node-linux-amd64 ./cmd/node

# Optional: Build for other platforms
GOOS=darwin GOARCH=amd64 go build -o bin/node-macos-amd64 ./cmd/node  # macOS Intel
GOOS=darwin GOARCH=arm64 go build -o bin/node-macos-arm64 ./cmd/node  # macOS M1/M2
GOOS=windows GOARCH=amd64 go build -o bin/node-windows-amd64.exe ./cmd/node  # Windows
```

### Step 2: Distribute Binaries

**Copy the binary to other machines:**

```bash
# Using scp
scp bin/node-linux-amd64 user@remote-machine:/usr/local/bin/shadownet-node

# Or download from a server
wget https://your-server.com/shadownet-node
chmod +x shadownet-node
```

**No Go installation needed on target machines!** ✅

---

## Option 2: Docker Deployment (Easiest!)

### Control Plane + Dashboard (Public Server)

Create `docker-compose.yml` on your VPS:

```yaml
version: '3.8'

services:
  controlplane:
    image: shadownet-controlplane:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - LISTEN_ADDR=:8080
      - DB_PATH=/data/controlplane.db
    restart: unless-stopped

  dashboard:
    build: ./web
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_CONTROLPLANE_URL=http://your-vps-ip:8080
    depends_on:
      - controlplane
    restart: unless-stopped
```

**Deploy:**
```bash
docker-compose up -d
```

### Node (On Each Machine)

**Using Docker:**
```bash
docker run -d \
  --name shadownet-node \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -e PEER_ID=$(hostname) \
  -e CONTROLPLANE_URL=http://your-vps-ip:8080 \
  shadownet-node:latest
```

---

## Option 3: Systemd Service (Production)

### Control Plane Setup (VPS/Cloud Server)

**1. Copy binary to server:**
```bash
scp bin/controlplane-linux-amd64 user@vps:/usr/local/bin/shadownet-controlplane
ssh user@vps "chmod +x /usr/local/bin/shadownet-controlplane"
```

**2. Create systemd service:**
```bash
ssh user@vps
sudo nano /etc/systemd/system/shadownet-controlplane.service
```

```ini
[Unit]
Description=ShadowNet Control Plane
After=network.target

[Service]
Type=simple
User=shadownet
ExecStart=/usr/local/bin/shadownet-controlplane --listen :8080 --db /var/lib/shadownet/controlplane.db
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**3. Start service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable shadownet-controlplane
sudo systemctl start shadownet-controlplane
```

### Node Setup (Each Machine)

**1. Copy binary:**
```bash
# On each machine
wget https://your-server.com/shadownet-node
sudo mv shadownet-node /usr/local/bin/
sudo chmod +x /usr/local/bin/shadownet-node
```

**2. Create systemd service:**
```bash
sudo nano /etc/systemd/system/shadownet-node.service
```

```ini
[Unit]
Description=ShadowNet VPN Node
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/shadownet-node \
  --id %H \
  --controlplane-url http://YOUR_VPS_IP:8080 \
  --private-key-path /etc/shadownet/node.key \
  --listen-port 51820
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**3. Start service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable shadownet-node
sudo systemctl start shadownet-node
```

---

## Quick Setup Script for New Nodes

Create this script on your dev machine and distribute it:

**`install-node.sh`:**
```bash
#!/bin/bash
set -e

CONTROLPLANE_URL="${CONTROLPLANE_URL:-http://your-vps-ip:8080}"
BINARY_URL="${BINARY_URL:-https://your-server.com/shadownet-node}"

echo "=== ShadowNet Node Installer ==="
echo "Control Plane: $CONTROLPLANE_URL"
echo ""

# Download binary
echo "Downloading ShadowNet node..."
wget -q $BINARY_URL -O /tmp/shadownet-node
sudo mv /tmp/shadownet-node /usr/local/bin/
sudo chmod +x /usr/local/bin/shadownet-node

# Create config directory
sudo mkdir -p /etc/shadownet

# Create systemd service
echo "Creating systemd service..."
sudo tee /etc/systemd/system/shadownet-node.service > /dev/null <<EOF
[Unit]
Description=ShadowNet VPN Node
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/shadownet-node \\
  --id $(hostname) \\
  --controlplane-url $CONTROLPLANE_URL \\
  --private-key-path /etc/shadownet/node.key \\
  --listen-port 51820
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# Start service
echo "Starting ShadowNet node..."
sudo systemctl daemon-reload
sudo systemctl enable shadownet-node
sudo systemctl start shadownet-node

echo ""
echo "✅ ShadowNet node installed and started!"
echo ""
echo "Check status: sudo systemctl status shadownet-node"
echo "View logs: sudo journalctl -u shadownet-node -f"
```

**Usage on any new machine:**
```bash
curl -sSL https://your-server.com/install-node.sh | sudo bash
```

---

## Real-World Deployment Example

### Scenario: 3 Machines

**Machine 1: VPS (Control Plane)**
- IP: `203.0.113.10`
- Runs: Control plane + Dashboard
- No VPN node needed

**Machine 2: Home Computer**
- Behind NAT
- Runs: VPN node
- Virtual IP: `10.10.0.1`

**Machine 3: Office Computer**
- Behind NAT
- Runs: VPN node
- Virtual IP: `10.10.0.2`

### Setup Steps

**1. On VPS (Machine 1):**
```bash
# Copy binary
scp bin/controlplane-linux-amd64 user@203.0.113.10:/usr/local/bin/shadownet-controlplane

# SSH and start
ssh user@203.0.113.10
sudo /usr/local/bin/shadownet-controlplane --listen :8080 --db /var/lib/shadownet/controlplane.db
```

**2. On Home Computer (Machine 2):**
```bash
# Download binary (no Go needed!)
wget http://203.0.113.10/shadownet-node
chmod +x shadownet-node

# Start node
sudo ./shadownet-node \
  --id home-pc \
  --controlplane-url http://203.0.113.10:8080 \
  --virtual-ip 10.10.0.1
```

**3. On Office Computer (Machine 3):**
```bash
# Download binary
wget http://203.0.113.10/shadownet-node
chmod +x shadownet-node

# Start node
sudo ./shadownet-node \
  --id office-pc \
  --controlplane-url http://203.0.113.10:8080 \
  --virtual-ip 10.10.0.2
```

**4. Test Connectivity:**

From home PC:
```bash
ping 10.10.0.2  # Ping office PC
```

From office PC:
```bash
ping 10.10.0.1  # Ping home PC
```

**Traffic flows directly between home and office - NOT through VPS!** ✅

---

## Firewall Configuration

### Control Plane (VPS)
```bash
# Allow control plane API
sudo ufw allow 8080/tcp

# Allow dashboard
sudo ufw allow 3000/tcp
```

### Nodes (All Machines)
```bash
# Allow WireGuard UDP
sudo ufw allow 51820/udp
```

---

## Hosting Binaries for Distribution

### Option 1: Simple HTTP Server

On your dev machine:
```bash
cd bin
python3 -m http.server 8000
```

Then download on other machines:
```bash
wget http://your-dev-machine-ip:8000/node-linux-amd64
```

### Option 2: GitHub Releases

1. Create a release on GitHub
2. Upload binaries as release assets
3. Download with:
```bash
wget https://github.com/Vaibhav2154/ShadowNet/releases/download/v1.0.0/node-linux-amd64
```

### Option 3: Cloud Storage

Upload to S3, Google Cloud Storage, or Dropbox and share the link.

---

## Minimal Setup Requirements Per Machine

### Control Plane Machine
- ✅ Linux server (VPS/Cloud)
- ✅ Public IP address
- ✅ Ports 8080, 3000 open
- ❌ No Go required (use binary)

### Node Machines
- ✅ Linux (Ubuntu 20.04+)
- ✅ Root/sudo access (for TUN device)
- ✅ Internet connection
- ❌ No Go required (use binary)
- ❌ No public IP needed (NAT traversal works!)

---

## Monitoring

### Check Node Status
```bash
sudo systemctl status shadownet-node
sudo journalctl -u shadownet-node -f
```

### Check WireGuard
```bash
sudo wg show
```

### View Dashboard
Open browser: `http://your-vps-ip:3000`

---

## Troubleshooting

### Node can't connect to control plane
```bash
# Test connectivity
curl http://your-vps-ip:8080/health

# Check firewall
sudo ufw status
```

### Peers can't connect to each other
```bash
# Check WireGuard status
sudo wg show

# Check if hole punching is working
sudo tcpdump -i any udp port 51820
```

### Permission denied on TUN device
```bash
# Run with sudo
sudo /usr/local/bin/shadownet-node ...

# Or set capabilities
sudo setcap cap_net_admin+ep /usr/local/bin/shadownet-node
```

---

## Summary

**You only need Go on ONE machine (your dev machine) to build the binaries.**

**Other machines just need:**
1. Download the binary
2. Run it with sudo
3. Point it to your control plane

**That's it!** No Go installation, no compilation, just copy and run! 🚀

---

## Next Steps

1. **Build binaries** on your dev machine
2. **Deploy control plane** to a VPS
3. **Copy node binary** to other machines
4. **Start nodes** pointing to control plane
5. **Test connectivity** with ping
6. **Monitor** via dashboard

**The beauty of Go: Single static binary, no dependencies!** ✨
