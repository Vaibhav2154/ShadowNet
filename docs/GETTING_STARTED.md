# 🚀 Getting Started with ShadowNet

This guide will help you quickly set up and run ShadowNet on your system.

## Prerequisites

- **Linux** (Ubuntu 20.04+ recommended)
- **Go 1.21+** for building from source
- **Root/sudo access** for running nodes (TUN device creation)
- **Node.js 20+** for dashboard (optional)

## Quick Start (3 Steps)

### 1. Build the Binaries

```bash
# Clone the repository (if not already done)
cd /path/to/ShadowNet

# Build both binaries
mkdir -p bin
go build -o bin/controlplane ./cmd/controlplane
go build -o bin/node ./cmd/node
```

This creates:
- `bin/controlplane` (12MB) - Control plane server
- `bin/node` (10MB) - VPN node client

### 2. Start the Control Plane

In terminal 1:

```bash
./bin/controlplane --listen :8080 --db ./data/controlplane.db
```

You should see:
```
Starting control plane server on :8080
Database: ./data/controlplane.db
Active timeout: 5m0s
```

Test it:
```bash
curl http://localhost:8080/health
# Should return: OK
```

### 3. Start a Node

In terminal 2 (requires sudo):

```bash
sudo ./bin/node \
  --id node1 \
  --controlplane-url http://localhost:8080 \
  --virtual-ip 10.10.0.1
```

You should see:
```
Loaded WireGuard keys (public: ...)
Created TUN device: tun0 (10.10.0.1/24)
Discovered public endpoint: <your-ip>:51820
Registered with control plane
ShadowNet node started successfully
```

**That's it!** Your first node is running. 🎉

---

## Testing with Two Nodes

To test P2P connectivity, run a second node on a different machine (or same machine with different port):

**On Machine 2:**

```bash
sudo ./bin/node \
  --id node2 \
  --controlplane-url http://<machine1-ip>:8080 \
  --virtual-ip 10.10.0.2 \
  --listen-port 51821
```

**Test connectivity:**

From Machine 1:
```bash
ping 10.10.0.2
```

From Machine 2:
```bash
ping 10.10.0.1
```

If pings work, you have a working P2P mesh VPN! 🚀

---

## Using the Quick Start Script

We've included a helper script for easier deployment:

```bash
# Make it executable
chmod +x quickstart.sh

# Start control plane
./quickstart.sh controlplane

# Start node (in another terminal)
sudo ./quickstart.sh node

# Start dashboard (optional)
./quickstart.sh dashboard
```

---

## Dashboard (Optional)

The web dashboard lets you monitor your mesh network.

```bash
cd web
npm install
export NEXT_PUBLIC_CONTROLPLANE_URL="http://localhost:8080"
npm run dev
```

Open http://localhost:3000 in your browser.

---

## Docker Deployment

For production deployment:

```bash
docker-compose up -d
```

This starts:
- Control plane on port 8080
- Dashboard on port 3000

---

## Command-Line Options

### Control Plane

```bash
./bin/controlplane [options]

Options:
  --listen string          Listen address (default ":8080")
  --db string             SQLite database path (default "./data/controlplane.db")
  --active-timeout duration  Peer active timeout (default 5m)
  --api-key string        Optional API key for authentication
```

### Node

```bash
sudo ./bin/node [options]

Options:
  --id string                  Peer ID (auto-generated if empty)
  --controlplane-url string    Control plane URL (default "http://localhost:8080")
  --private-key-path string    Private key file (default "./shadownet.key")
  --listen-port int           WireGuard listen port (default 51820)
  --stun-server string        STUN server address (default "stun.l.google.com:19302")
  --virtual-ip string         Virtual IP address (auto-assigned if empty)
  --tun-device string         TUN device name (default "tun0")
  --heartbeat-interval duration  Heartbeat interval (default 30s)
```

---

## Environment Variables

You can also use environment variables:

**Control Plane:**
```bash
export LISTEN_ADDR=":8080"
export DB_PATH="./data/controlplane.db"
./bin/controlplane
```

**Node:**
```bash
export PEER_ID="my-node"
export CONTROLPLANE_URL="http://controlplane.example.com:8080"
sudo ./bin/node
```

---

## Troubleshooting

### "Permission denied" when starting node

**Solution:** Node requires root for TUN device creation:
```bash
sudo ./bin/node ...
```

### "Failed to create TUN device"

**Solution:** Ensure you're on Linux and have CAP_NET_ADMIN capability:
```bash
sudo setcap cap_net_admin+ep ./bin/node
./bin/node ...
```

### "Connection refused" to control plane

**Solution:** Ensure control plane is running and accessible:
```bash
curl http://<controlplane-ip>:8080/health
```

### Peers can't connect

**Possible causes:**
1. **Firewall blocking UDP** - Open port 51820 (or your listen port)
2. **Symmetric NAT** - Current implementation may not work (future: relay server)
3. **STUN server unreachable** - Try different STUN server

**Debug steps:**
```bash
# Check if STUN works
dig stun.l.google.com

# Check WireGuard status
sudo wg show

# Check control plane logs
curl http://localhost:8080/metrics
```

---

## Network Architecture

```
┌─────────────┐
│   Node 1    │
│ 10.10.0.1   │
└──────┬──────┘
       │
       │ Encrypted WireGuard Tunnel
       │ (Direct P2P Connection)
       │
┌──────┴──────┐         ┌─────────────────┐
│   Node 2    │◄────────┤ Control Plane   │
│ 10.10.0.2   │  Coord  │  (Discovery)    │
└─────────────┘         └─────────────────┘
```

**Key Points:**
- Control plane only handles **coordination** (peer discovery, endpoints)
- All **data traffic** flows directly between peers (P2P)
- **End-to-end encryption** via WireGuard
- **NAT traversal** via STUN + UDP hole punching

---

## What's Next?

1. **Add more nodes** - Scale your mesh network
2. **Deploy control plane** - Put it on a public server
3. **Configure firewall** - Open UDP ports for WireGuard
4. **Monitor with dashboard** - Track peer status and metrics
5. **Automate deployment** - Use systemd services or Docker

---

## Production Deployment Tips

### 1. Run Control Plane as Service

Create `/etc/systemd/system/shadownet-controlplane.service`:

```ini
[Unit]
Description=ShadowNet Control Plane
After=network.target

[Service]
Type=simple
User=shadownet
ExecStart=/usr/local/bin/controlplane --listen :8080 --db /var/lib/shadownet/controlplane.db
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable shadownet-controlplane
sudo systemctl start shadownet-controlplane
```

### 2. Run Node as Service

Create `/etc/systemd/system/shadownet-node.service`:

```ini
[Unit]
Description=ShadowNet Node
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/node --id %H --controlplane-url http://controlplane.example.com:8080
Restart=always

[Install]
WantedBy=multi-user.target
```

### 3. Use TLS for Control Plane

Put control plane behind nginx with Let's Encrypt:

```nginx
server {
    listen 443 ssl;
    server_name controlplane.example.com;
    
    ssl_certificate /etc/letsencrypt/live/controlplane.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/controlplane.example.com/privkey.pem;
    
    location / {
        proxy_pass http://localhost:8080;
    }
}
```

---

## Support

- **Documentation:** See `docs/` directory
- **Issues:** Report bugs on GitHub
- **Architecture:** Read `docs/ARCHITECTURE.md`
- **API Reference:** See `docs/CONTROL_PLANE_API.md`

---

## License

MIT License - See LICENSE file for details.

---

**Happy meshing! 🌐**
