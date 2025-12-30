# Deployment & Ops

## Docker Compose (example)
```yaml
services:
  controlplane:
    image: golang:1.22
    working_dir: /app
    volumes:
      - ./:/app
    command: ["bash", "-lc", "go run ./cmd/controlplane"]
    ports:
      - "8080:8080"
    environment:
      - CP_LISTEN_ADDR=:8080
      - CP_SQLITE_PATH=/app/data/controlplane.db
    restart: unless-stopped

  web:
    image: node:20
    working_dir: /app/web
    volumes:
      - ./web:/app/web
    environment:
      - NEXT_PUBLIC_CONTROLPLANE_URL=http://localhost:8080
    command: ["bash", "-lc", "npm install && npm run dev"]
    ports:
      - "3000:3000"
    restart: unless-stopped
```

## Production Notes
- Use a dedicated SQLite or migrate to Postgres for HA
- Place control plane behind TLS (reverse proxy like Caddy/NGINX)
- Rate-limit registration endpoints; validate payloads
- Persistent keys and secure storage for node private keys

## Observability
- Minimal `/metrics` for dashboard; consider Prometheus for detailed metrics
- Logs: structure with levels; centralize if needed

## Backups
- Backup the control plane database regularly
