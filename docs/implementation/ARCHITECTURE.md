# Architecture Deep Dive

## Separation of Planes
- Control plane: peer registry, discovery, health/metrics
- Data plane: end-to-end encrypted traffic via WireGuard, direct P2P

## Components
- Control Plane API (`internal/controlplane/api/*`)
- Services (`internal/controlplane/service/*`)
- Storage (`internal/controlplane/store/*`) using SQLite
- Node runtime (`internal/node/*`): STUN, NAT traversal, WireGuard userspace, TUN
- Dashboard (`web/`): Next.js app consuming control plane metrics and peer data

## Data Model
```go
type Peer struct {
    ID           string
    WGPublicKey  string
    EndpointIP   string
    EndpointPort int
    LastSeen     time.Time
}
```

## Flows
1. Node discovers public endpoint via STUN
2. Node registers with control plane (`/register`)
3. Node fetches peers (`/peers`)
4. NAT hole punching to direct endpoints
5. WireGuard peers configured; traffic flows P2P
6. Node sends heartbeat (`/heartbeat`)

## NAT Traversal
- STUN on the same UDP socket as data
- Frequent empty UDP packets 300–500ms to maintain mappings

## Userspace WireGuard
- TUN device reading/writing raw IP
- WireGuard device initialized, peers added dynamically
- Keys managed per node; private key never leaves the node
