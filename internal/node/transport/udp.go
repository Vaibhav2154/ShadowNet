package transport

import (
	"fmt"
	"net"
)

// UDPTransport manages UDP socket for WireGuard and STUN
type UDPTransport struct {
	conn *net.UDPConn
	port int
}

// NewUDPTransport creates a new UDP transport
func NewUDPTransport(port int) (*UDPTransport, error) {
	addr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP socket: %w", err)
	}

	// Set socket options for better performance
	// Enable address reuse
	file, err := conn.File()
	if err == nil {
		defer file.Close()
		// Additional socket options could be set here via syscall
	}

	return &UDPTransport{
		conn: conn,
		port: port,
	}, nil
}

// Conn returns the UDP connection
func (t *UDPTransport) Conn() *net.UDPConn {
	return t.conn
}

// Port returns the listen port
func (t *UDPTransport) Port() int {
	return t.port
}

// Close closes the UDP connection
func (t *UDPTransport) Close() error {
	return t.conn.Close()
}

// LocalAddr returns the local address
func (t *UDPTransport) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}
