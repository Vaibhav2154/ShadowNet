package api

import (
	"log"
	"net/http"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/service"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
)

// PeersHandler handles peer listing
type PeersHandler struct {
	peerService *service.PeerService
}

// NewPeersHandler creates a new peers handler
func NewPeersHandler(peerService *service.PeerService) *PeersHandler {
	return &PeersHandler{
		peerService: peerService,
	}
}

// ServeHTTP handles GET /peers
func (h *PeersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get optional exclude parameter
	excludeID := r.URL.Query().Get("exclude")

	// Get active peers
	peers, err := h.peerService.GetActivePeers(excludeID)
	if err != nil {
		log.Printf("Failed to get peers: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to get peers")
		return
	}

	// Handle nil slice
	if peers == nil {
		peers = []*proto.PeerInfo{}
	}

	// Convert to response format
	var peerList []proto.PeerInfo
	for _, p := range peers {
		peerList = append(peerList, *p)
	}

	response := proto.PeersResponse{
		Peers: peerList,
		Count: len(peerList),
	}

	writeJSON(w, http.StatusOK, response)
}
