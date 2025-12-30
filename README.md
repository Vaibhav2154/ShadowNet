## ShadowNet — A Peer-to-Peer Mesh VPN using WireGuard

---

## 1. Introduction

ShadowNet is a production-grade **peer-to-peer (P2P) mesh VPN** built using **userspace WireGuard**, **UDP hole punching**, and a **centralized control plane** for peer coordination.  
The system is designed with strict **control-plane / data-plane separation**, ensuring privacy, scalability, and performance.

The control plane is responsible only for coordination and discovery, while **all encrypted traffic flows directly between peers**.

---

## 2. Architecture Overview

### 2.1 High-Level Architecture

Control Plane:
- Peer registration
- Public key exchange
- Endpoint discovery
- Health tracking

Data Plane:
- Userspace WireGuard tunnel
- UDP NAT traversal
- Encrypted peer-to-peer packet routing

Traffic never passes through the control plane.

---

## 3. Technology Stack

### Backend
- Go (Golang)
- wireguard-go
- wgctrl-go
- pion/stun
- SQLite
- Linux TUN interface

### Frontend
- Next.js
- Tailwind CSS

### Deployment
- Docker
- Docker Compose

---

## 4. Control Plane Implementation

### 4.1 Responsibilities

The control plane:
- Stores peer metadata
- Distributes peer information
- Tracks liveness via heartbeat
- Exposes metrics for UI

It does NOT:
- Handle encryption
- Relay traffic
- Inspect packets

---

### 4.2 Data Model

```go
type Peer struct {
    ID            string
    WGPublicKey   string
    EndpointIP    string
    EndpointPort  int
    LastSeen      time.Time
}
```

---

### 4.3 API Endpoints

#### POST /register
Registers a peer with its public key and discovered endpoint.

Payload:
```json
{
  "id": "peer-1",
  "public_key": "base64key",
  "endpoint_ip": "203.0.113.5",
  "endpoint_port": 51820
}
```

---

#### GET /peers
Returns all active peers except the requester.

---

#### POST /heartbeat
Keeps peer marked as online.

---

## 5. NAT Traversal

### 5.1 STUN Discovery

Each node uses STUN to discover its public-facing IP and UDP port.

Key rules:
- Use the same UDP socket for STUN and WireGuard
- Do not rebind ports
- Persist connections

---

### 5.2 UDP Hole Punching

Peers repeatedly send empty UDP packets to each other's public endpoints.
This creates NAT mappings on both sides.

Punch interval: 300–500ms

---

## 6. WireGuard Userspace Integration

### 6.1 Why Userspace WireGuard

- No kernel module dependency
- Easier debugging
- Programmatic control
- Portable deployment

---

### 6.2 WireGuard Initialization Flow

1. Create TUN device
2. Start wireguard-go with TUN FD
3. Configure interface using wgctrl-go
4. Dynamically add peers and endpoints

---

### 6.3 Key Management

- Curve25519 keys
- Generated per node
- Public keys exchanged via control plane
- Private keys never leave the node

---

## 7. TUN Device Handling

### 7.1 TUN Creation

```bash
ip tuntap add dev tun0 mode tun
ip addr add 10.10.0.1/24 dev tun0
ip link set tun0 up
```

The Go process reads raw IP packets from the TUN device.

---

## 8. Packet Flow

```
Application
   ↓
TUN Interface
   ↓
WireGuard Encryption
   ↓
UDP Socket
   ↓
Internet
   ↓
UDP Socket
   ↓
WireGuard Decryption
   ↓
TUN Interface
   ↓
Application
```

---

## 9. Node Runtime Lifecycle

1. Load configuration
2. Generate or load WireGuard keys
3. Discover public endpoint via STUN
4. Register with control plane
5. Fetch peers
6. Perform NAT hole punching
7. Establish WireGuard tunnels
8. Start packet forwarding
9. Send heartbeats

---

## 10. Dashboard Implementation

### Features
- Online peers
- Handshake timestamps
- Bytes sent/received
- Endpoint details
- Connection state

The UI polls the control plane periodically.

---

## 11. Security Model

- End-to-end encryption using WireGuard
- Control plane has zero visibility into traffic
- Replay protection via WireGuard protocol
- No plaintext metadata leakage beyond endpoints

---

## 12. Limitations

- Symmetric NATs may fail without relays
- No automatic key rotation (future work)
- No relay fallback (DERP-like)

---

## 13. Future Improvements

- Relay server fallback
- IPv6 support
- Key rotation
- ACLs and policy engine
- Mobile support

---

## 14. Conclusion

ShadowNet demonstrates a real-world implementation of a **secure, NAT-traversed, peer-to-peer mesh VPN** using modern networking primitives.  
It mirrors the architecture of industry-grade systems while remaining fully open and inspectable.

---
