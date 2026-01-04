package nat

import (
	"context"
	"log"
	"net"
	"time"
)

// HolePuncher manages NAT hole punching for a peer
type HolePuncher struct {
	conn     *net.UDPConn
	endpoint *net.UDPAddr
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewHolePuncher creates a new hole puncher
func NewHolePuncher(conn *net.UDPConn, remoteEndpoint string, interval time.Duration) (*HolePuncher, error) {
	addr, err := net.ResolveUDPAddr("udp4", remoteEndpoint)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &HolePuncher{
		conn:     conn,
		endpoint: addr,
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start begins sending hole punching packets
func (h *HolePuncher) Start() {
	go h.punch()
}

// Stop stops the hole punching
func (h *HolePuncher) Stop() {
	h.cancel()
}

// punch sends periodic empty UDP packets to create NAT mapping
func (h *HolePuncher) punch() {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	// Send initial packet immediately
	h.sendPacket()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.sendPacket()
		}
	}
}

// sendPacket sends a single hole punching packet
func (h *HolePuncher) sendPacket() {
	// Send empty packet (or small marker)
	_, err := h.conn.WriteToUDP([]byte{0x00}, h.endpoint)
	if err != nil {
		log.Printf("Hole punch failed to %s: %v", h.endpoint, err)
	}
}

// PunchManager manages multiple hole punchers
type PunchManager struct {
	punchers map[string]*HolePuncher
}

// NewPunchManager creates a new punch manager
func NewPunchManager() *PunchManager {
	return &PunchManager{
		punchers: make(map[string]*HolePuncher),
	}
}

// AddPeer adds a peer to punch
func (m *PunchManager) AddPeer(peerID string, conn *net.UDPConn, endpoint string, interval time.Duration) error {
	puncher, err := NewHolePuncher(conn, endpoint, interval)
	if err != nil {
		return err
	}

	m.punchers[peerID] = puncher
	puncher.Start()

	log.Printf("Started hole punching for peer %s to %s", peerID, endpoint)
	return nil
}

// RemovePeer removes a peer from punching
func (m *PunchManager) RemovePeer(peerID string) {
	if puncher, ok := m.punchers[peerID]; ok {
		puncher.Stop()
		delete(m.punchers, peerID)
		log.Printf("Stopped hole punching for peer %s", peerID)
	}
}

// StopAll stops all hole punchers
func (m *PunchManager) StopAll() {
	for peerID, puncher := range m.punchers {
		puncher.Stop()
		delete(m.punchers, peerID)
	}
}
