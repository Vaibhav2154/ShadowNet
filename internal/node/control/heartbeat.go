package control

import (
	"context"
	"log"
	"time"
)

// HeartbeatSender sends periodic heartbeats to the control plane
type HeartbeatSender struct {
	client   *Client
	peerID   string
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewHeartbeatSender creates a new heartbeat sender
func NewHeartbeatSender(client *Client, peerID string, interval time.Duration) *HeartbeatSender {
	ctx, cancel := context.WithCancel(context.Background())

	return &HeartbeatSender{
		client:   client,
		peerID:   peerID,
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start begins sending heartbeats
func (h *HeartbeatSender) Start() {
	go h.run()
}

// Stop stops sending heartbeats
func (h *HeartbeatSender) Stop() {
	h.cancel()
}

// run is the main heartbeat loop
func (h *HeartbeatSender) run() {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	// Send initial heartbeat immediately
	h.sendHeartbeat()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.sendHeartbeat()
		}
	}
}

// sendHeartbeat sends a single heartbeat
func (h *HeartbeatSender) sendHeartbeat() {
	if err := h.client.SendHeartbeat(h.peerID); err != nil {
		log.Printf("Failed to send heartbeat: %v", err)
	} else {
		log.Printf("Heartbeat sent successfully")
	}
}
