package service

import (
	"fmt"
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/model"
	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/store"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/crypto"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/utils"
)

// PeerService handles peer business logic
type PeerService struct {
	repo           store.PeerRepository
	activeTimeout  time.Duration
	startTime      time.Time
}

// NewPeerService creates a new peer service
func NewPeerService(repo store.PeerRepository, activeTimeout time.Duration) *PeerService {
	return &PeerService{
		repo:          repo,
		activeTimeout: activeTimeout,
		startTime:     time.Now(),
	}
}

// RegisterPeer validates and registers a new peer
func (s *PeerService) RegisterPeer(info *proto.PeerInfo) error {
	// Validate peer info
	if info.ID == "" {
		return fmt.Errorf("peer ID is required")
	}
	
	if err := crypto.ValidatePublicKey(info.WGPublicKey); err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	
	if err := utils.ValidateIP(info.EndpointIP); err != nil {
		return fmt.Errorf("invalid endpoint IP: %w", err)
	}
	
	if err := utils.ValidatePort(info.EndpointPort); err != nil {
		return fmt.Errorf("invalid endpoint port: %w", err)
	}
	
	// Create peer model
	peer := model.FromProto(info)
	peer.LastSeen = time.Now()
	
	// Store peer
	if err := s.repo.CreateOrUpdate(peer); err != nil {
		return fmt.Errorf("failed to store peer: %w", err)
	}
	
	return nil
}

// GetActivePeers returns all active peers, optionally excluding one
func (s *PeerService) GetActivePeers(excludeID string) ([]*proto.PeerInfo, error) {
	peers, err := s.repo.GetAllActive(s.activeTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to get active peers: %w", err)
	}
	
	var result []*proto.PeerInfo
	for _, peer := range peers {
		if excludeID != "" && peer.ID == excludeID {
			continue
		}
		info := peer.ToProto()
		result = append(result, &info)
	}
	
	return result, nil
}

// UpdateHeartbeat updates the last seen timestamp for a peer
func (s *PeerService) UpdateHeartbeat(id string) error {
	if id == "" {
		return fmt.Errorf("peer ID is required")
	}
	
	if err := s.repo.UpdateLastSeen(id); err != nil {
		return fmt.Errorf("failed to update heartbeat: %w", err)
	}
	
	return nil
}

// GetMetrics returns control plane metrics
func (s *PeerService) GetMetrics() (*proto.MetricsResponse, error) {
	allPeers, err := s.repo.GetAllActive(24 * time.Hour) // All peers in last 24h
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	
	activePeers, err := s.repo.GetAllActive(s.activeTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to get active peers: %w", err)
	}
	
	uptime := time.Since(s.startTime)
	
	return &proto.MetricsResponse{
		TotalPeers:  len(allPeers),
		ActivePeers: len(activePeers),
		Uptime:      uptime.String(),
		Timestamp:   time.Now(),
	}, nil
}

// GetPeerByID retrieves a specific peer by ID
func (s *PeerService) GetPeerByID(id string) (*proto.PeerInfo, error) {
	peer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer: %w", err)
	}
	
	if peer == nil {
		return nil, fmt.Errorf("peer not found: %s", id)
	}
	
	info := peer.ToProto()
	return &info, nil
}
