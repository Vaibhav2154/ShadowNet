# Security Model & Considerations

## Principles
- End-to-end encryption with WireGuard
- Control plane does not relay traffic; only coordinates
- Minimal metadata: endpoints and liveness only

## Key Management
- Curve25519 keys generated per node
- Public keys exchanged via control plane
- Private keys never leave the node

## Hardening
- TLS for control plane
- Auth tokens for write endpoints (`/register`, `/heartbeat`)
- Input validation and rate limiting
- Secure storage for keys and configs

## Threats
- NAT traversal exposure of UDP endpoints
- Denial of service on control plane
- Compromised node exposing metadata (not traffic)
