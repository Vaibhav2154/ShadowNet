package api

import (
	"log"
	"net/http"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/service"
)

// MetricsHandler handles metrics requests
type MetricsHandler struct {
	peerService *service.PeerService
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(peerService *service.PeerService) *MetricsHandler {
	return &MetricsHandler{
		peerService: peerService,
	}
}

// ServeHTTP handles GET /metrics
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get metrics
	metrics, err := h.peerService.GetMetrics()
	if err != nil {
		log.Printf("Failed to get metrics: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to get metrics")
		return
	}

	writeJSON(w, http.StatusOK, metrics)
}
