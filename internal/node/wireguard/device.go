package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Device represents a WireGuard device using kernel module
type Device struct {
	interfaceName string
	configPath    string
}

// NewDevice creates a new WireGuard device using kernel module
func NewDevice(interfaceName string, privateKey *PrivateKey, virtualIP string, listenPort int) (*Device, error) {
	// Create WireGuard config file
	configPath := fmt.Sprintf("/etc/wireguard/%s.conf", interfaceName)

	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/24
ListenPort = %d
`, privateKey.String(), virtualIP, listenPort)

	// Write config file
	if err := os.MkdirAll("/etc/wireguard", 0700); err != nil {
		return nil, fmt.Errorf("failed to create wireguard directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		return nil, fmt.Errorf("failed to write config: %w", err)
	}

	// Bring up interface using wg-quick
	cmd := exec.Command("wg-quick", "up", interfaceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to bring up interface: %w (output: %s)", err, string(output))
	}

	return &Device{
		interfaceName: interfaceName,
		configPath:    configPath,
	}, nil
}

// AddPeer adds a peer to the WireGuard device
func (d *Device) AddPeer(publicKey *PublicKey, endpoint string, allowedIPs []string) error {
	// Build wg command to add peer
	args := []string{
		"set", d.interfaceName,
		"peer", publicKey.String(),
	}

	if endpoint != "" {
		args = append(args, "endpoint", endpoint)
	}

	if len(allowedIPs) > 0 {
		args = append(args, "allowed-ips", strings.Join(allowedIPs, ","))
	}

	// Add persistent keepalive
	args = append(args, "persistent-keepalive", "25")

	cmd := exec.Command("wg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add peer: %w (output: %s)", err, string(output))
	}

	return nil
}

// UpdatePeerEndpoint updates a peer's endpoint
func (d *Device) UpdatePeerEndpoint(publicKey *PublicKey, endpoint string) error {
	cmd := exec.Command("wg", "set", d.interfaceName,
		"peer", publicKey.String(),
		"endpoint", endpoint,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update endpoint: %w (output: %s)", err, string(output))
	}

	return nil
}

// RemovePeer removes a peer from the WireGuard device
func (d *Device) RemovePeer(publicKey *PublicKey) error {
	cmd := exec.Command("wg", "set", d.interfaceName,
		"peer", publicKey.String(),
		"remove",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove peer: %w (output: %s)", err, string(output))
	}

	return nil
}

// Close closes the WireGuard device
func (d *Device) Close() error {
	// Bring down interface
	cmd := exec.Command("wg-quick", "down", d.interfaceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to bring down interface: %w (output: %s)", err, string(output))
	}

	// Remove config file
	if err := os.Remove(d.configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config: %w", err)
	}

	return nil
}

// Wait is a no-op for kernel WireGuard
func (d *Device) Wait() {
	// Kernel WireGuard doesn't need to wait
}
