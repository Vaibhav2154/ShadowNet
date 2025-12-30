# Setup & Development

## Prerequisites
- Linux with TUN/TAP support
- Go 1.22+
- Node.js 20+ and a package manager (pnpm recommended)
- Docker (optional for deployment)

## Repository Layout
- `internal/controlplane`: HTTP API, storage (SQLite), services
- `internal/node`: STUN, NAT hole punching, WireGuard (userspace), TUN
- `web`: Next.js dashboard
- `cmd/controlplane`, `cmd/node`: CLI entrypoints (to be wired)

## Local Dev Environment
```bash
# Go tools
go version
go mod tidy

# Web UI
cd web
pnpm install
pnpm dev
```

## Environment Variables
- `NEXT_PUBLIC_CONTROLPLANE_URL`: Dashboard base URL for the control plane
- Control plane service variables (suggested):
  - `CP_LISTEN_ADDR` (e.g., `:8080`)
  - `CP_SQLITE_PATH` (e.g., `./data/controlplane.db`)
- Node runtime variables (suggested):
  - `NODE_ID`, `NODE_PRIVATE_KEY_PATH`
  - `CONTROLPLANE_URL`
  - `STUN_SERVER` (e.g., `stun.l.google.com:19302`)

## Permissions
WireGuard userspace + TUN requires CAP_NET_ADMIN. Run with elevated privileges or grant capabilities to the binary.

```bash
sudo setcap cap_net_admin+ep ./node
```

## Code Quality
```bash
gofmt -s -w .
go vet ./...

cd web
pnpm lint
```
