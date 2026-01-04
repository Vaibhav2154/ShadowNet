package model

import (
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
)

// Peer represents a peer in the database
type Peer struct {
	ID           string
	WGPublicKey  string
	EndpointIP   string
	EndpointPort int
	LastSeen     time.Time
}

// ToProto converts database model to API proto
func (p *Peer) ToProto() proto.PeerInfo {
	return proto.PeerInfo{
		ID:           p.ID,
		WGPublicKey:  p.WGPublicKey,
		EndpointIP:   p.EndpointIP,
		EndpointPort: p.EndpointPort,
		LastSeen:     p.LastSeen.Format(time.RFC3339),
	}
}

// FromProto creates a Peer from API proto
func FromProto(info *proto.PeerInfo) *Peer {
	lastSeen := time.Now()
	if info.LastSeen != "" {
		if t, err := time.Parse(time.RFC3339, info.LastSeen); err == nil {
			lastSeen = t
		}
	}
	
	return &Peer{
		ID:           info.ID,
		WGPublicKey:  info.WGPublicKey,
		EndpointIP:   info.EndpointIP,
		EndpointPort: info.EndpointPort,
		LastSeen:     lastSeen,
	}
}
