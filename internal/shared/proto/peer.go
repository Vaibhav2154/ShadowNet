package proto

import "time"

// PeerInfo represents peer information exchanged via API
type PeerInfo struct {
	ID           string `json:"id"`
	WGPublicKey  string `json:"wg_public_key"`
	EndpointIP   string `json:"endpoint_ip"`
	EndpointPort int    `json:"endpoint_port"`
	LastSeen     string `json:"last_seen,omitempty"`
}

type RegisterRequest struct {
	ID           string `json:"id"`
	WGPublicKey  string `json:"wg_public_key"`
	EndpointIP   string `json:"endpoint_ip"`
	EndpointPort int    `json:"endpoint_port"`
}

// RegisterResponse is returned after successful registration
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// HeartbeatRequest is sent periodically to keep peer alive
type HeartbeatRequest struct {
	ID string `json:"id"`
}

// HeartbeatResponse confirms heartbeat receipt
type HeartbeatResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// PeersResponse contains list of active peers
type PeersResponse struct {
	Peers []PeerInfo `json:"peers"`
	Count int        `json:"count"`
}

// MetricsResponse contains control plane metrics
type MetricsResponse struct {
	TotalPeers  int       `json:"total_peers"`
	ActivePeers int       `json:"active_peers"`
	Uptime      string    `json:"uptime"`
	Timestamp   time.Time `json:"timestamp"`
}

// ErrorResponse is returned on API errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
