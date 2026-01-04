package utils

import (
	"fmt"
	"net"
	"strconv"
)

// ValidateIP validates an IP address string
func ValidateIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}
	return nil
}

// ValidatePort validates a port number
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", port)
	}
	return nil
}

// FormatEndpoint formats IP and port as "ip:port"
func FormatEndpoint(ip string, port int) string {
	return net.JoinHostPort(ip, strconv.Itoa(port))
}

// ParseEndpoint parses "ip:port" string into IP and port
func ParseEndpoint(endpoint string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return "", 0, fmt.Errorf("invalid endpoint format: %w", err)
	}
	
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port number: %w", err)
	}
	
	if err := ValidateIP(host); err != nil {
		return "", 0, err
	}
	
	if err := ValidatePort(port); err != nil {
		return "", 0, err
	}
	
	return host, port, nil
}

// GetLocalIP attempts to get the local IP address
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
