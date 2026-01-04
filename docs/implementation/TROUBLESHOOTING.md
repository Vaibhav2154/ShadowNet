# ShadowNet Troubleshooting Guide

## Issue: Nodes Can't Ping Each Other

### Symptoms
- Nodes show as "Active" in dashboard
- TUN interfaces exist with correct IPs
- `ping` results in 100% packet loss
- `wg show` returns no output or no peers

### Root Cause
WireGuard peers aren't being configured properly or the device state isn't persisting.

### Quick Diagnosis

```bash
# Check if WireGuard is running
docker exec shadownet-node-1 wg show

# Should show:
# - interface: tun0
# - public key
# - listening port
# - peer entries with endpoints

# If empty or missing peers, WireGuard isn't configured
```

### Solution 1: Check Node Logs

```bash
# Look for errors in node startup
docker logs shadownet-node-2 2>&1 | grep -i error

# Common errors:
# - "Unable to update bind" - Port conflict (expected, can ignore)
# - "Failed to add peer" - Peer configuration failed
# - "Failed to configure device" - WireGuard setup failed
```

### Solution 2: Verify Peer Discovery

```bash
# Check if nodes are discovering each other
docker logs shadownet-node-2 2>&1 | grep "Added peer"

# Should see:
# "Added peer: docker-node-1"
# "Added peer: docker-node-3"
```

### Solution 3: Manual WireGuard Check

```bash
# Enter container
docker exec -it shadownet-node-2 sh

# Check WireGuard interface
wg show tun0

# Check if peers are configured
wg show tun0 peers

# Check allowed IPs
wg show tun0 allowed-ips
```

### Solution 4: Restart Nodes

```bash
# Sometimes nodes need a restart after initial setup
docker-compose restart node1 node2 node3

# Wait 10 seconds for registration
sleep 10

# Try ping again
docker exec shadownet-node-2 ping -c 3 10.10.0.118
```

### Solution 5: Check Control Plane

```bash
# Verify control plane has all peers
curl http://localhost:8080/peers | jq

# Should show 3 peers with:
# - Different peer IDs
# - Different public keys
# - Recent last_seen timestamps
```

### Solution 6: Rebuild Containers

```bash
# Clean rebuild
docker-compose down
docker-compose up -d --build

# Wait for startup
sleep 15

# Check logs
docker-compose logs -f
```

## Expected Working State

When everything is working, you should see:

**1. Node Logs:**
```
✅ Loaded WireGuard keys
✅ Created TUN device: tun0
✅ Discovered public endpoint
✅ WireGuard device created
✅ Registered with control plane
✅ Found X active peers
✅ Added peer: docker-node-X
✅ Started hole punching
✅ Heartbeat sent successfully
```

**2. WireGuard Status:**
```bash
$ docker exec shadownet-node-2 wg show

interface: tun0
  public key: abc123...
  listening port: 51820

peer: def456...
  endpoint: 172.20.0.21:51820
  allowed ips: 10.10.0.118/32
  latest handshake: 5 seconds ago
  transfer: 1.2 KiB received, 892 B sent

peer: ghi789...
  endpoint: 172.20.0.23:51820
  allowed ips: 10.10.0.120/32
  latest handshake: 3 seconds ago
  transfer: 980 B received, 756 B sent
```

**3. Successful Ping:**
```bash
$ docker exec shadownet-node-2 ping -c 3 10.10.0.118

PING 10.10.0.118 (10.10.0.118): 56 data bytes
64 bytes from 10.10.0.118: seq=0 ttl=64 time=0.123 ms
64 bytes from 10.10.0.118: seq=1 ttl=64 time=0.098 ms
64 bytes from 10.10.0.118: seq=2 ttl=64 time=0.105 ms

--- 10.10.0.118 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
```

## Common Issues

### Issue: "Unable to update bind: listen udp4 :51820: bind: address already in use"

**Status:** ⚠️ Warning (can be ignored)

**Cause:** WireGuard tries to bind twice during initialization

**Fix:** None needed - this is expected behavior

### Issue: No peers showing in `wg show`

**Status:** ❌ Critical

**Cause:** Peer configuration not being applied

**Fix:** Check node logs for "Added peer" messages. If missing, control plane communication failed.

### Issue: Peers configured but ping fails

**Status:** ❌ Critical

**Cause:** NAT traversal or routing issue

**Fix:**
```bash
# Check if handshake is happening
docker exec shadownet-node-2 wg show | grep handshake

# If "latest handshake: Never", peers can't connect
# Try restarting to trigger new hole punching
docker-compose restart
```

### Issue: Handshake happens but still no ping

**Status:** ❌ Critical

**Cause:** Allowed IPs not configured correctly

**Fix:**
```bash
# Check allowed IPs
docker exec shadownet-node-2 wg show | grep "allowed ips"

# Should show other nodes' virtual IPs
# If missing, peer configuration failed
```

## Debug Commands

```bash
# Full diagnostic
echo "=== Node 1 ===" && docker exec shadownet-node-1 wg show
echo "=== Node 2 ===" && docker exec shadownet-node-2 wg show
echo "=== Node 3 ===" && docker exec shadownet-node-3 wg show

# Check all virtual IPs
for i in 1 2 3; do
  echo "Node $i:"
  docker exec shadownet-node-$i ip addr show tun0 | grep "inet "
done

# Test connectivity matrix
echo "Node 2 -> Node 1:"
docker exec shadownet-node-2 ping -c 2 -W 1 10.10.0.118

echo "Node 2 -> Node 3:"
docker exec shadownet-node-2 ping -c 2 -W 1 10.10.0.120

echo "Node 1 -> Node 2:"
docker exec shadownet-node-1 ping -c 2 -W 1 10.10.0.119
```

## Next Steps

If none of the above works, the issue is likely in the WireGuard device initialization or peer configuration code. Check:

1. `internal/node/wireguard/device.go` - Device creation
2. `internal/node/node.go` - Peer configuration logic
3. Node logs for specific error messages

The nodes are successfully:
- ✅ Creating TUN interfaces
- ✅ Registering with control plane
- ✅ Discovering peers
- ✅ Starting hole punching

But failing at:
- ❌ Configuring WireGuard peers
- ❌ Establishing encrypted tunnels
