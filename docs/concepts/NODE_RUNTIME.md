# Node Runtime

## Lifecycle
1. Load configuration and keys
2. Discover public endpoint via STUN
3. Register with control plane
4. Fetch peers
5. Perform NAT hole punching
6. Initialize WireGuard userspace device
7. Add peers and endpoints dynamically
8. Start packet forwarding from TUN
9. Send heartbeats periodically

## Configuration
Suggested fields:
```json
{
  "id": "peer-1",
  "private_key_path": "./keys/private.key",
  "controlplane_url": "http://localhost:8080",
  "stun_server": "stun.l.google.com:19302",
  "punch_interval_ms": 400
}
```

## STUN + NAT Traversal
- Use the same UDP socket for STUN and WireGuard
- Maintain mappings with periodic empty packets

## WireGuard
- Userspace WireGuard device initialized against TUN
- Curve25519 keys; public keys exchanged via control plane
- Replay protection and encryption handled by WireGuard protocol

## Permissions
- CAP_NET_ADMIN required for TUN operations
