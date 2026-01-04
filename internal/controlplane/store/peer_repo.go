package store

import (
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/model"
)

// PeerRepository defines the interface for peer storage operations
type PeerRepository interface {
	// CreateOrUpdate creates a new peer or updates existing one
	CreateOrUpdate(peer *model.Peer) error
	
	// GetByID retrieves a peer by ID
	GetByID(id string) (*model.Peer, error)
	
	// GetAllActive retrieves all peers active within the timeout duration
	GetAllActive(timeout time.Duration) ([]*model.Peer, error)
	
	// UpdateLastSeen updates the last seen timestamp for a peer
	UpdateLastSeen(id string) error
	
	// Delete removes a peer from storage
	Delete(id string) error
	
	// Close closes the repository connection
	Close() error
}
