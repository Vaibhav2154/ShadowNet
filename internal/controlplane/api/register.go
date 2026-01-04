package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/service"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
)

// RegisterHandler handles peer registration
type RegisterHandler struct {
	peerService *service.PeerService
}

// NewRegisterHandler creates a new register handler
func NewRegisterHandler(peerService *service.PeerService) *RegisterHandler {
	return &RegisterHandler{
		peerService: peerService,
	}
}

// ServeHTTP handles POST /register
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req proto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Convert to PeerInfo
	peerInfo := &proto.PeerInfo{
		ID:           req.ID,
		WGPublicKey:  req.WGPublicKey,
		EndpointIP:   req.EndpointIP,
		EndpointPort: req.EndpointPort,
	}

	// Register peer
	if err := h.peerService.RegisterPeer(peerInfo); err != nil {
		log.Printf("Failed to register peer %s: %v", req.ID, err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Registered peer: %s (%s:%d)", req.ID, req.EndpointIP, req.EndpointPort)

	// Send success response
	writeJSON(w, http.StatusOK, proto.RegisterResponse{
		Success: true,
		Message: "peer registered successfully",
	})
}

// Helper functions
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, proto.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
