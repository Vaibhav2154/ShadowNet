package wireguard

import (
	"bufio"
	"fmt"
	"strings"
)

// Note: This file provides additional WireGuard control utilities
// The main device management is in device.go

// PeerStats represents statistics for a peer
type PeerStats struct {
	PublicKey        string
	Endpoint         string
	LastHandshake    int64
	BytesReceived    int64
	BytesTransmitted int64
	AllowedIPs       []string
}

// GetPeerStats retrieves statistics for all peers (placeholder)
// In a full implementation, this would parse IPC output
func (d *Device) GetPeerStats() ([]PeerStats, error) {
	// This would require parsing the IPC get operation
	// For now, return empty slice
	return []PeerStats{}, nil
}

// SetFwMark sets the firewall mark for the device
func (d *Device) SetFwMark(mark uint32) error {
	config := fmt.Sprintf("fwmark=%d\n", mark)

	if err := d.device.IpcSetOperation(bufio.NewReader(strings.NewReader(config))); err != nil {
		return fmt.Errorf("failed to set fwmark: %w", err)
	}

	return nil
}
