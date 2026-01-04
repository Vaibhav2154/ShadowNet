package tun

import (
	"fmt"
	"os/exec"

	"golang.zx2c4.com/wireguard/tun"
)

// Device represents a TUN device
type Device struct {
	tunDevice tun.Device
	name      string
}

// CreateTUN creates and configures a TUN device
func CreateTUN(name, ip, netmask string) (*Device, error) {
	// Create TUN device
	tunDevice, err := tun.CreateTUN(name, 1420) // MTU 1420 for WireGuard
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN device: %w", err)
	}

	// Get actual device name (may differ from requested)
	actualName, err := tunDevice.Name()
	if err != nil {
		tunDevice.Close()
		return nil, fmt.Errorf("failed to get TUN device name: %w", err)
	}

	// Configure IP address
	if err := configureIP(actualName, ip, netmask); err != nil {
		tunDevice.Close()
		return nil, fmt.Errorf("failed to configure TUN device: %w", err)
	}

	return &Device{
		tunDevice: tunDevice,
		name:      actualName,
	}, nil
}

// configureIP configures the IP address and brings the interface up
func configureIP(name, ip, netmask string) error {
	// Add IP address
	cmd := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", ip, netmask), "dev", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add IP address: %w", err)
	}

	// Bring interface up
	cmd = exec.Command("ip", "link", "set", "dev", name, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring interface up: %w", err)
	}

	return nil
}

// Name returns the device name
func (d *Device) Name() string {
	return d.name
}

// Device returns the underlying TUN device
func (d *Device) Device() tun.Device {
	return d.tunDevice
}

// Close closes the TUN device
func (d *Device) Close() error {
	return d.tunDevice.Close()
}

// File returns the file descriptor (for compatibility)
func (d *Device) File() (int, error) {
	// Get file descriptor from TUN device
	// This is platform-specific and may need adjustment
	return -1, fmt.Errorf("file descriptor access not implemented")
}
