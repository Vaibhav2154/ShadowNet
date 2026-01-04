package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/node"
	"github.com/Vaibhav2154/ShadowNet/internal/node/config"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/crypto"
)

func main() {
	// Parse command-line flags
	cfg := config.DefaultConfig()

	flag.StringVar(&cfg.ID, "id", getEnv("PEER_ID", ""), "Peer ID (required)")
	flag.StringVar(&cfg.ControlPlaneURL, "controlplane-url", getEnv("CONTROLPLANE_URL", "http://localhost:8080"), "Control plane URL")
	flag.StringVar(&cfg.PrivateKeyPath, "private-key-path", getEnv("PRIVATE_KEY_PATH", "./shadownet.key"), "Private key file path")
	flag.IntVar(&cfg.ListenPort, "listen-port", 51820, "WireGuard listen port")
	flag.StringVar(&cfg.STUNServer, "stun-server", getEnv("STUN_SERVER", "stun.l.google.com:19302"), "STUN server address")
	flag.DurationVar(&cfg.PunchInterval, "punch-interval", 500*time.Millisecond, "NAT hole punch interval")
	flag.StringVar(&cfg.TUNDeviceName, "tun-device", "tun0", "TUN device name")
	flag.StringVar(&cfg.VirtualIP, "virtual-ip", "", "Virtual IP address (auto-assigned if empty)")
	flag.StringVar(&cfg.VirtualNetmask, "virtual-netmask", "24", "Virtual network netmask")
	flag.DurationVar(&cfg.HeartbeatInterval, "heartbeat-interval", 30*time.Second, "Heartbeat interval")

	flag.Parse()

	// Generate peer ID if not provided
	if cfg.ID == "" {
		cfg.ID = crypto.GenerateID()
		log.Printf("Generated peer ID: %s", cfg.ID)
	}

	// Auto-assign virtual IP if not provided
	if cfg.VirtualIP == "" {
		// Simple IP assignment based on peer ID hash
		hash := 0
		for _, c := range cfg.ID {
			hash = (hash*31 + int(c)) % 254
		}
		cfg.VirtualIP = fmt.Sprintf("10.10.0.%d", hash+1)
		log.Printf("Auto-assigned virtual IP: %s", cfg.VirtualIP)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create node
	n, err := node.NewNode(cfg)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	// Start node
	if err := n.Start(); err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("Node running. Press Ctrl+C to stop.")
	<-sigChan

	// Graceful shutdown
	log.Println("Received shutdown signal")
	if err := n.Stop(); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Node stopped successfully")
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
