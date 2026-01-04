# Implementation Plan (CLI + Services)

This plan wires the `cmd/` entrypoints to the internal packages and delivers a working control plane and node runtime.

## Milestone A: Control Plane CLI
1. Define flags/env: `--listen :8080`, `--db ./data/controlplane.db`
2. Initialize SQLite via `internal/controlplane/store/sqlite.go`
3. Wire services (`service/*`) and API handlers (`api/*`) to an HTTP server in `server.go`
4. Start HTTP server; expose `/register`, `/peers`, `/heartbeat`, `/metrics`
5. Add graceful shutdown and logging

## Milestone B: Node CLI
1. Define flags/env: `--id`, `--controlplane-url`, `--private-key-path`, `--stun-server`, `--punch-interval`
2. Generate/load key pair via `internal/node/wireguard/keys.go`
3. Discover endpoint via `internal/node/stun/stun.go`
4. Register with control plane; fetch peers
5. Start NAT hole punching via `internal/node/nat/hole_punch.go`
6. Initialize WireGuard device/tun via `internal/node/wireguard/device.go` and `tun/tun.go`
7. Add peers/endpoints dynamically; start forwarding
8. Send heartbeats

## Milestone C: Dashboard Integration
1. Hook UI to `/metrics` and `/peers`
2. Add auto-refresh and connection state indicators
3. Provide basic settings (control plane URL)

## Milestone D: Deployment & Ops
1. Provide Dockerfile(s) for control plane and dashboard
2. Compose file (see `docs/DEPLOYMENT.md`)
3. Add TLS and reverse proxy

## Testing & Validation
- Unit tests for services and NAT routines (future)
- Manual E2E: two nodes behind different NATs

## Roadmap Enhancements
- Relay server fallback (DERP-like)
- IPv6 support
- Key rotation
- ACLs and policy engine
- Mobile clients
