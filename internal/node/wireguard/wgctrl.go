package wireguard

// GetPeerStats returns statistics for a peer (placeholder for kernel WireGuard)
func (d *Device) GetPeerStats(publicKey *PublicKey) (map[string]interface{}, error) {
	// For kernel WireGuard, we could parse `wg show` output
	// For now, return empty stats
	return map[string]interface{}{}, nil
}

// SetFirewallMark sets the firewall mark (placeholder for kernel WireGuard)
func (d *Device) SetFirewallMark(mark uint32) error {
	// Kernel WireGuard handles this via config file
	return nil
}
