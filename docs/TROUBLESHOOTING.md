# Troubleshooting & FAQs

## Common Issues

- Symmetric NAT prevents direct P2P
  - Use relay fallback (future), or co-locate nodes

- STUN unreachable
  - Verify STUN server host/port; check firewall

- TUN permissions
  - Grant `CAP_NET_ADMIN` or run with `sudo`

- No WireGuard handshake
  - Check key pair, peer public keys, endpoints
  - Ensure UDP allowed by firewall/NAT

- Peer not visible in dashboard
  - Ensure `/register` and `/heartbeat` succeed
  - Control plane DB connectivity (SQLite path)

## Diagnostics
```bash
ping 8.8.8.8
nc -zu <peer-ip> <port>
wg show
journalctl -u <service>
```
