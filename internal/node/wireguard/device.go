package wireguard

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

// Device represents a WireGuard device
type Device struct {
	device     *device.Device
	tunDevice  tun.Device
	listenPort int
}

// NewDevice creates a new WireGuard device
func NewDevice(tunDevice tun.Device, privateKey *PrivateKey, listenPort int) (*Device, error) {
	// Create logger
	logger := device.NewLogger(
		device.LogLevelError,
		fmt.Sprintf("(%s) ", "wireguard"),
	)

	// Create WireGuard device
	wgDevice := device.NewDevice(tunDevice, conn.NewDefaultBind(), logger)

	// Configure device
	config := fmt.Sprintf("private_key=%s\nlisten_port=%d\n",
		privateKey.String(),
		listenPort,
	)

	if err := wgDevice.IpcSetOperation(bufio.NewReader(strings.NewReader(config))); err != nil {
		return nil, fmt.Errorf("failed to configure device: %w", err)
	}

	// Bring device up
	wgDevice.Up()

	log.Printf("WireGuard device created on port %d", listenPort)

	return &Device{
		device:     wgDevice,
		tunDevice:  tunDevice,
		listenPort: listenPort,
	}, nil
}

// AddPeer adds a peer to the WireGuard device
func (d *Device) AddPeer(publicKey *PublicKey, endpoint string, allowedIPs []string) error {
	// Build peer configuration
	config := fmt.Sprintf("public_key=%s\n", publicKey.String())

	if endpoint != "" {
		config += fmt.Sprintf("endpoint=%s\n", endpoint)
	}

	// Add allowed IPs
	for _, ip := range allowedIPs {
		config += fmt.Sprintf("allowed_ip=%s\n", ip)
	}

	// Enable persistent keepalive (25 seconds)
	config += "persistent_keepalive_interval=25\n"

	// Apply configuration
	if err := d.device.IpcSetOperation(bufio.NewReader(strings.NewReader(config))); err != nil {
		return fmt.Errorf("failed to add peer: %w", err)
	}

	log.Printf("Added peer with public key %s...", publicKey.String()[:16])
	return nil
}

// UpdatePeerEndpoint updates a peer's endpoint
func (d *Device) UpdatePeerEndpoint(publicKey *PublicKey, endpoint string) error {
	config := fmt.Sprintf("public_key=%s\nendpoint=%s\n",
		publicKey.String(),
		endpoint,
	)

	if err := d.device.IpcSetOperation(bufio.NewReader(strings.NewReader(config))); err != nil {
		return fmt.Errorf("failed to update peer endpoint: %w", err)
	}

	return nil
}

// RemovePeer removes a peer from the WireGuard device
func (d *Device) RemovePeer(publicKey *PublicKey) error {
	config := fmt.Sprintf("public_key=%s\nremove=true\n", publicKey.String())

	if err := d.device.IpcSetOperation(bufio.NewReader(strings.NewReader(config))); err != nil {
		return fmt.Errorf("failed to remove peer: %w", err)
	}

	return nil
}

// Close closes the WireGuard device
func (d *Device) Close() error {
	d.device.Close()
	return nil
}

// Wait waits for the device to close
func (d *Device) Wait() {
	d.device.Wait()
}
