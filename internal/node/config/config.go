package config

import (
	"fmt"
	"time"
)

// Config holds node configuration
type Config struct {
	// Peer identification
	ID string
	
	// Control plane
	ControlPlaneURL string
	
	// WireGuard
	PrivateKeyPath string
	ListenPort     int
	
	// Network
	STUNServer     string
	PunchInterval  time.Duration
	
	// TUN device
	TUNDeviceName  string
	VirtualIP      string
	VirtualNetmask string
	
	// Heartbeat
	HeartbeatInterval time.Duration
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("peer ID is required")
	}
	
	if c.ControlPlaneURL == "" {
		return fmt.Errorf("control plane URL is required")
	}
	
	if c.PrivateKeyPath == "" {
		return fmt.Errorf("private key path is required")
	}
	
	if c.ListenPort < 1 || c.ListenPort > 65535 {
		return fmt.Errorf("invalid listen port: %d", c.ListenPort)
	}
	
	if c.STUNServer == "" {
		return fmt.Errorf("STUN server is required")
	}
	
	if c.TUNDeviceName == "" {
		return fmt.Errorf("TUN device name is required")
	}
	
	if c.VirtualIP == "" {
		return fmt.Errorf("virtual IP is required")
	}
	
	return nil
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		ListenPort:        51820,
		STUNServer:        "stun.l.google.com:19302",
		PunchInterval:     500 * time.Millisecond,
		TUNDeviceName:     "tun0",
		VirtualNetmask:    "24",
		HeartbeatInterval: 30 * time.Second,
	}
}
