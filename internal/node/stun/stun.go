package stun

import (
	"fmt"
	"net"
	"time"

	"github.com/pion/stun"
)

// DiscoverEndpoint discovers the public endpoint using STUN
func DiscoverEndpoint(stunServer string, localPort int) (string, int, error) {
	// Create UDP connection on specific port
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: localPort,
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Resolve STUN server address
	serverAddr, err := net.ResolveUDPAddr("udp4", stunServer)
	if err != nil {
		return "", 0, fmt.Errorf("failed to resolve STUN server: %w", err)
	}

	// Create STUN binding request
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	// Send request
	_, err = conn.WriteToUDP(message.Raw, serverAddr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send STUN request: %w", err)
	}

	// Receive response
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return "", 0, fmt.Errorf("failed to receive STUN response: %w", err)
	}

	// Parse response
	var response stun.Message
	response.Raw = buf[:n]
	if err := response.Decode(); err != nil {
		return "", 0, fmt.Errorf("failed to decode STUN response: %w", err)
	}

	// Extract mapped address
	var xorAddr stun.XORMappedAddress
	if err := xorAddr.GetFrom(&response); err != nil {
		return "", 0, fmt.Errorf("failed to get XOR-MAPPED-ADDRESS: %w", err)
	}

	return xorAddr.IP.String(), xorAddr.Port, nil
}

// DiscoverEndpointWithConn discovers the public endpoint using an existing connection
func DiscoverEndpointWithConn(conn *net.UDPConn, stunServer string) (string, int, error) {
	// Resolve STUN server address
	serverAddr, err := net.ResolveUDPAddr("udp4", stunServer)
	if err != nil {
		return "", 0, fmt.Errorf("failed to resolve STUN server: %w", err)
	}

	// Create STUN binding request
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	// Send request
	_, err = conn.WriteToUDP(message.Raw, serverAddr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send STUN request: %w", err)
	}

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetReadDeadline(time.Time{})

	// Receive response
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return "", 0, fmt.Errorf("failed to receive STUN response: %w", err)
	}

	// Parse response
	var response stun.Message
	response.Raw = buf[:n]
	if err := response.Decode(); err != nil {
		return "", 0, fmt.Errorf("failed to decode STUN response: %w", err)
	}

	// Extract mapped address
	var xorAddr stun.XORMappedAddress
	if err := xorAddr.GetFrom(&response); err != nil {
		return "", 0, fmt.Errorf("failed to get XOR-MAPPED-ADDRESS: %w", err)
	}

	return xorAddr.IP.String(), xorAddr.Port, nil
}
