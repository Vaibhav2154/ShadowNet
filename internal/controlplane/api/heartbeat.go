package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/service"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
)

// HeartbeatHandler handles peer heartbeats
type HeartbeatHandler struct {
	peerService *service.PeerService
}

// NewHeartbeatHandler creates a new heartbeat handler
func NewHeartbeatHandler(peerService *service.PeerService) *HeartbeatHandler {
	return &HeartbeatHandler{
		peerService: peerService,
	}
}

// ServeHTTP handles POST /heartbeat
func (h *HeartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req proto.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update heartbeat
	if err := h.peerService.UpdateHeartbeat(req.ID); err != nil {
		log.Printf("Failed to update heartbeat for %s: %v", req.ID, err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Send success response
	writeJSON(w, http.StatusOK, proto.HeartbeatResponse{
		Success: true,
		Message: "heartbeat received",
	})
}
