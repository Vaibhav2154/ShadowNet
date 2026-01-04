# Control Plane API

Base URL: `http://<host>:<port>`

## POST /register
Registers a peer with its public key and discovered endpoint.

Request
```json
{
  "id": "peer-1",
  "public_key": "base64key",
  "endpoint_ip": "203.0.113.5",
  "endpoint_port": 51820
}
```

Response
```json
{ "status": "ok" }
```

## GET /peers
Returns all active peers except the requester.

Response
```json
[
  {
    "id": "peer-2",
    "public_key": "base64key",
    "endpoint_ip": "198.51.100.23",
    "endpoint_port": 51820,
    "last_seen": "2025-12-30T12:34:56Z"
  }
]
```

## POST /heartbeat
Marks the peer as online.

Request
```json
{ "id": "peer-1" }
```

Response
```json
{ "status": "ok" }
```

## GET /metrics
Exposes minimal counters for the dashboard (implementation-dependent).

Notes
- All payloads should be validated; control plane does not inspect encrypted traffic.
- Consider auth tokens for write endpoints in production.
