package node

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Vaibhav2154/ShadowNet/internal/node/config"
	"github.com/Vaibhav2154/ShadowNet/internal/node/control"
	"github.com/Vaibhav2154/ShadowNet/internal/node/nat"
	"github.com/Vaibhav2154/ShadowNet/internal/node/stun"
	"github.com/Vaibhav2154/ShadowNet/internal/node/transport"
	"github.com/Vaibhav2154/ShadowNet/internal/node/wireguard"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
	"github.com/Vaibhav2154/ShadowNet/internal/shared/utils"
)

// Node represents a ShadowNet node
type Node struct {
	config          *config.Config
	privateKey      *wireguard.PrivateKey
	publicKey       *wireguard.PublicKey
	wgDevice        *wireguard.Device
	udpTransport    *transport.UDPTransport
	controlClient   *control.Client
	heartbeatSender *control.HeartbeatSender
	punchManager    *nat.PunchManager
	publicIP        string
	publicPort      int
}

// NewNode creates a new node
func NewNode(cfg *config.Config) (*Node, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &Node{
		config:       cfg,
		punchManager: nat.NewPunchManager(),
	}, nil
}

// Start starts the node runtime
func (n *Node) Start() error {
	log.Println("Starting ShadowNet node...")

	// Step 1: Load or generate WireGuard keys
	if err := n.loadKeys(); err != nil {
		return fmt.Errorf("failed to load keys: %w", err)
	}
	log.Printf("Loaded WireGuard keys (public: %s...)", n.publicKey.String()[:16])

	// Step 2: Discover public endpoint via STUN (temporary UDP socket)
	if err := n.discoverEndpoint(); err != nil {
		return fmt.Errorf("failed to discover endpoint: %w", err)
	}
	log.Printf("Discovered public endpoint: %s:%d", n.publicIP, n.publicPort)

	// Step 3: Initialize WireGuard device (creates TUN and manages UDP automatically)
	if err := n.createWireGuard(); err != nil {
		return fmt.Errorf("failed to create WireGuard device: %w", err)
	}
	log.Printf("Initialized WireGuard device with IP %s", n.config.VirtualIP)

	// Step 6: Register with control plane
	if err := n.registerWithControlPlane(); err != nil {
		return fmt.Errorf("failed to register with control plane: %w", err)
	}
	log.Println("Registered with control plane")

	// Step 7: Fetch and configure peers
	if err := n.configurePeers(); err != nil {
		return fmt.Errorf("failed to configure peers: %w", err)
	}

	// Step 8: Start heartbeat
	n.startHeartbeat()
	log.Println("Started heartbeat sender")

	log.Println("ShadowNet node started successfully")
	return nil
}

// loadKeys loads or generates WireGuard keys
func (n *Node) loadKeys() error {
	privateKey, err := wireguard.LoadOrGeneratePrivateKey(n.config.PrivateKeyPath)
	if err != nil {
		return err
	}

	n.privateKey = privateKey
	n.publicKey = privateKey.PublicKey()
	return nil
}

// createTransport creates the UDP transport
func (n *Node) createTransport() error {
	transport, err := transport.NewUDPTransport(n.config.ListenPort)
	if err != nil {
		return err
	}

	n.udpTransport = transport
	return nil
}

// discoverEndpoint discovers the public endpoint using STUN with a temporary socket
func (n *Node) discoverEndpoint() error {
	// Check if we should use Docker internal IP instead of STUN
	if os.Getenv("USE_DOCKER_IP") == "true" {
		// Get Docker container's IP address
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return fmt.Errorf("failed to get interface addresses: %w", err)
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// Use first non-loopback IPv4 address (Docker internal IP)
					n.publicIP = ipnet.IP.String()
					n.publicPort = n.config.ListenPort
					log.Printf("Using Docker internal IP: %s:%d", n.publicIP, n.publicPort)
					return nil
				}
			}
		}
	}

	// Create temporary UDP socket for STUN discovery
	tempTransport, err := transport.NewUDPTransport(n.config.ListenPort)
	if err != nil {
		return err
	}
	defer tempTransport.Close() // Close immediately after STUN

	ip, port, err := stun.DiscoverEndpointWithConn(
		tempTransport.Conn(),
		n.config.STUNServer,
	)
	if err != nil {
		return err
	}

	n.publicIP = ip
	n.publicPort = port
	return nil
}

// createWireGuard initializes the WireGuard device using kernel module
func (n *Node) createWireGuard() error {
	wgDev, err := wireguard.NewDevice(
		n.config.TUNDeviceName,
		n.privateKey,
		n.config.VirtualIP,
		n.config.ListenPort,
	)
	if err != nil {
		return err
	}

	n.wgDevice = wgDev
	return nil
}

// registerWithControlPlane registers this node with the control plane
func (n *Node) registerWithControlPlane() error {
	n.controlClient = control.NewClient(n.config.ControlPlaneURL)

	peerInfo := &proto.PeerInfo{
		ID:           n.config.ID,
		WGPublicKey:  n.publicKey.String(),
		EndpointIP:   n.publicIP,
		EndpointPort: n.publicPort,
	}

	return n.controlClient.Register(peerInfo)
}

// configurePeers fetches peers and configures WireGuard
func (n *Node) configurePeers() error {
	// Fetch peers from control plane
	peers, err := n.controlClient.GetPeers(n.config.ID)
	if err != nil {
		return err
	}

	log.Printf("Found %d active peers", len(peers))

	// Configure each peer
	for _, peer := range peers {
		if err := n.addPeer(peer); err != nil {
			log.Printf("Warning: failed to add peer %s: %v", peer.ID, err)
			continue
		}
		log.Printf("Added peer: %s (%s:%d)", peer.ID, peer.EndpointIP, peer.EndpointPort)
	}

	return nil
}

// addPeer adds a peer to WireGuard
func (n *Node) addPeer(peer *proto.PeerInfo) error {
	// Parse public key
	publicKey, err := wireguard.ParsePublicKey(peer.WGPublicKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	// Use the peer's registered endpoint (works for both Docker and real deployments)
	endpoint := utils.FormatEndpoint(peer.EndpointIP, peer.EndpointPort)

	// Calculate peer's virtual IP using same hash function as main.go
	hash := 0
	for _, c := range peer.ID {
		hash = (hash*31 + int(c)) % 254
	}
	peerVirtualIP := fmt.Sprintf("10.10.0.%d", hash+1)

	// Determine allowed IPs (peer's virtual IP)
	allowedIPs := []string{fmt.Sprintf("%s/32", peerVirtualIP)}

	log.Printf("Adding peer %s with virtual IP %s, endpoint %s", peer.ID, peerVirtualIP, endpoint)

	// Add peer to WireGuard (persistent keepalive handles NAT traversal)
	if err := n.wgDevice.AddPeer(publicKey, endpoint, allowedIPs); err != nil {
		return fmt.Errorf("failed to add peer to WireGuard: %w", err)
	}

	return nil
}

// startHeartbeat starts the heartbeat sender
func (n *Node) startHeartbeat() {
	n.heartbeatSender = control.NewHeartbeatSender(
		n.controlClient,
		n.config.ID,
		n.config.HeartbeatInterval,
	)
	n.heartbeatSender.Start()
}

// Stop stops the node
func (n *Node) Stop() error {
	log.Println("Stopping ShadowNet node...")

	// Stop heartbeat
	if n.heartbeatSender != nil {
		n.heartbeatSender.Stop()
	}

	// Stop hole punching
	if n.punchManager != nil {
		n.punchManager.StopAll()
	}

	// Close WireGuard device (also closes TUN)
	if n.wgDevice != nil {
		n.wgDevice.Close()
	}

	log.Println("ShadowNet node stopped")
	return nil
}

// Wait waits for the node to stop
func (n *Node) Wait() {
	if n.wgDevice != nil {
		n.wgDevice.Wait()
	}
}

// hashPeerID creates a simple hash of peer ID for IP assignment
func hashPeerID(id string) int {
	hash := 0
	for _, c := range id {
		hash = (hash*31 + int(c)) % 254
	}
	return hash + 1 // 1-254
}
