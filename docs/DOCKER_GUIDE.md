# 🐳 Docker Deployment Guide for ShadowNet

## Quick Start (Test on Single Machine)

### Option 1: Full Test Setup (3 Nodes + Control Plane + Dashboard)

This runs everything on your machine for testing:

```bash
# Build and start all services
docker-compose up -d

# Watch logs
docker-compose logs -f

# Check status
docker-compose ps
```

**What this starts:**
- ✅ Control Plane (port 8080)
- ✅ Dashboard (port 3000)
- ✅ Node 1 (docker-node-1)
- ✅ Node 2 (docker-node-2)
- ✅ Node 3 (docker-node-3)

**Access:**
- Dashboard: http://localhost:3000
- Control Plane API: http://localhost:8080

You'll see all 3 nodes appear in the dashboard! 🎉

### Option 2: Production Setup (Control Plane + Dashboard Only)

For deploying to a VPS:

```bash
# Use production compose file
docker-compose -f docker-compose.prod.yml up -d
```

Then run nodes on separate machines (see below).

---

## Complete Testing Walkthrough

### Step 1: Start the Stack

```bash
cd /home/vaibhi/Dev/ShadowNet

# Build and start
docker-compose up -d --build

# This will:
# 1. Build control plane image
# 2. Build dashboard image
# 3. Build node image
# 4. Start all containers
```

### Step 2: Watch Services Start

```bash
# Watch all logs
docker-compose logs -f

# Or watch specific service
docker-compose logs -f node1
docker-compose logs -f controlplane
docker-compose logs -f dashboard
```

**Expected output:**
```
controlplane_1  | Starting control plane server on :8080
node1_1         | Loaded WireGuard keys (public: ...)
node1_1         | Created TUN device: tun0
node1_1         | Registered with control plane
node2_1         | Loaded WireGuard keys (public: ...)
node2_1         | Registered with control plane
node3_1         | Loaded WireGuard keys (public: ...)
node3_1         | Registered with control plane
```

### Step 3: Open Dashboard

Open browser: **http://localhost:3000**

You should see:
- **Active Peers: 3**
- **Total Peers: 3**
- **Control Plane: Online**
- All 3 nodes in the table with green status indicators

### Step 4: Test P2P Connectivity

```bash
# Enter node1 container
docker exec -it shadownet-node-1 sh

# Check WireGuard status
wg show

# Try to ping other nodes (via their virtual IPs)
# Note: Virtual IPs are auto-assigned based on peer ID hash
ping 10.10.0.X  # Replace X with actual virtual IP
```

### Step 5: Monitor in Real-Time

The dashboard auto-refreshes every 5 seconds. Watch as:
- Nodes send heartbeats
- Status indicators pulse
- Last seen timestamps update

---

## Docker Commands Cheat Sheet

### Start/Stop

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Restart a specific service
docker-compose restart node1

# Stop and remove everything (including volumes)
docker-compose down -v
```

### Logs

```bash
# All logs
docker-compose logs -f

# Specific service
docker-compose logs -f node1

# Last 100 lines
docker-compose logs --tail=100 node1
```

### Status

```bash
# Check running containers
docker-compose ps

# Check resource usage
docker stats

# Health check status
docker-compose ps
```

### Debugging

```bash
# Enter a container
docker exec -it shadownet-node-1 sh

# Check WireGuard status
docker exec shadownet-node-1 wg show

# Check network
docker network inspect shadownet_shadownet

# View control plane database
docker exec shadownet-controlplane ls -la /data
```

---

## Adding More Nodes

Want to test with more nodes? Just add to `docker-compose.yml`:

```yaml
  node4:
    build:
      context: .
      dockerfile: Dockerfile.node
    container_name: shadownet-node-4
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun
    environment:
      - PEER_ID=docker-node-4
      - CONTROLPLANE_URL=http://172.20.0.10:8080
    depends_on:
      controlplane:
        condition: service_healthy
    restart: unless-stopped
    networks:
      shadownet:
        ipv4_address: 172.20.0.24
```

Then:
```bash
docker-compose up -d node4
```

---

## Production Deployment

### Deploy to VPS

**1. Copy files to VPS:**
```bash
scp -r /home/vaibhi/Dev/ShadowNet user@vps:/opt/shadownet
```

**2. SSH to VPS:**
```bash
ssh user@vps
cd /opt/shadownet
```

**3. Update dashboard URL:**
```bash
# Edit docker-compose.prod.yml
nano docker-compose.prod.yml

# Change NEXT_PUBLIC_CONTROLPLANE_URL to your VPS IP
# - NEXT_PUBLIC_CONTROLPLANE_URL=http://YOUR_VPS_IP:8080
```

**4. Start services:**
```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

**5. Check firewall:**
```bash
sudo ufw allow 8080/tcp  # Control plane
sudo ufw allow 3000/tcp  # Dashboard
```

### Run Nodes on Other Machines

**Option A: Docker on other machines**

```bash
# On each machine
docker run -d \
  --name shadownet-node \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -e PEER_ID=$(hostname) \
  -e CONTROLPLANE_URL=http://YOUR_VPS_IP:8080 \
  --restart unless-stopped \
  shadownet-node:latest
```

**Option B: Binary on other machines**

See `DEPLOYMENT.md` for binary deployment.

---

## Troubleshooting

### Nodes won't start

**Error:** `cannot create TUN device`

**Solution:** Make sure `/dev/net/tun` exists:
```bash
# On host machine
ls -la /dev/net/tun

# If missing, create it
sudo mkdir -p /dev/net
sudo mknod /dev/net/tun c 10 200
sudo chmod 666 /dev/net/tun
```

### Dashboard shows "Error: Failed to fetch"

**Solution:** Check control plane is running:
```bash
docker-compose logs controlplane
curl http://localhost:8080/health
```

### Nodes not appearing in dashboard

**Solution:** Check node logs:
```bash
docker-compose logs node1

# Look for:
# - "Registered with control plane" ✅
# - Connection errors ❌
```

### Can't ping between nodes

**Solution:** Check WireGuard status:
```bash
docker exec shadownet-node-1 wg show

# Should show peer connections
```

---

## Architecture

```
┌─────────────────────────────────────┐
│   Docker Host (Your Machine)       │
│                                     │
│  ┌──────────────┐  ┌─────────────┐ │
│  │ Control      │  │ Dashboard   │ │
│  │ Plane :8080  │  │ :3000       │ │
│  └──────┬───────┘  └─────────────┘ │
│         │                           │
│    ┌────┴────┬─────────┬──────┐    │
│    │         │         │      │    │
│  ┌─▼──┐   ┌─▼──┐   ┌─▼──┐   │    │
│  │Node│   │Node│   │Node│   │    │
│  │ 1  │◄──┤ 2  │◄──┤ 3  │   │    │
│  └────┘   └────┘   └────┘   │    │
│    P2P Encrypted Tunnels     │    │
└─────────────────────────────────────┘
```

All running on one machine for testing!

---

## Scaling

### Single Machine Testing
- ✅ Control Plane
- ✅ Dashboard
- ✅ 3+ Nodes
- Perfect for development and testing

### Production Deployment
- ✅ Control Plane on VPS (Docker)
- ✅ Dashboard on VPS (Docker)
- ✅ Nodes on separate machines (Docker or binary)

---

## Environment Variables

### Control Plane
- `LISTEN_ADDR` - Listen address (default: `:8080`)
- `DB_PATH` - Database path (default: `/data/controlplane.db`)
- `ACTIVE_TIMEOUT` - Peer timeout (default: `5m`)

### Node
- `PEER_ID` - Unique peer identifier (required)
- `CONTROLPLANE_URL` - Control plane URL (required)

### Dashboard
- `NEXT_PUBLIC_CONTROLPLANE_URL` - Control plane URL for browser

---

## Next Steps

1. **Test locally:**
   ```bash
   docker-compose up -d
   # Open http://localhost:3000
   ```

2. **Deploy to VPS:**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Add real nodes:**
   - Run Docker containers on other machines
   - Or use binary deployment

4. **Monitor:**
   - Watch dashboard for real-time updates
   - Check logs with `docker-compose logs -f`

---

## Clean Up

```bash
# Stop all services
docker-compose down

# Remove volumes (deletes database)
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Complete cleanup
docker system prune -a
```

---

## Success! 🎉

You can now:
- ✅ Test complete mesh network on one machine
- ✅ See all nodes in the dashboard
- ✅ Monitor real-time status
- ✅ Deploy to production easily

**Start testing:**
```bash
docker-compose up -d
# Open http://localhost:3000
```

🚀 **Your P2P mesh VPN is running!**
