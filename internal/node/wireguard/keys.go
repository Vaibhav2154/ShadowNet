package wireguard

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/curve25519"
)

const (
	// KeyLength is the length of WireGuard keys in bytes
	KeyLength = 32
)

// PrivateKey represents a WireGuard private key
type PrivateKey [KeyLength]byte

// PublicKey represents a WireGuard public key
type PublicKey [KeyLength]byte

// GeneratePrivateKey generates a new random private key
func GeneratePrivateKey() (*PrivateKey, error) {
	var key PrivateKey
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	// Clamp the key as per Curve25519 spec
	key[0] &= 248
	key[31] &= 127
	key[31] |= 64

	return &key, nil
}

// PublicKey derives the public key from a private key
func (k *PrivateKey) PublicKey() *PublicKey {
	var pub PublicKey
	curve25519.ScalarBaseMult((*[32]byte)(&pub), (*[32]byte)(k))
	return &pub
}

// String returns the base64-encoded private key
func (k *PrivateKey) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}

// HexString returns the hex-encoded private key (for WireGuard IPC)
func (k *PrivateKey) HexString() string {
	return fmt.Sprintf("%x", k[:])
}

// String returns the base64-encoded public key
func (k *PublicKey) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}

// HexString returns the hex-encoded public key (for WireGuard IPC)
func (k *PublicKey) HexString() string {
	return fmt.Sprintf("%x", k[:])
}

// ParsePrivateKey parses a base64-encoded private key
func ParsePrivateKey(encoded string) (*PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	if len(decoded) != KeyLength {
		return nil, fmt.Errorf("invalid key length: expected %d, got %d", KeyLength, len(decoded))
	}

	var key PrivateKey
	copy(key[:], decoded)
	return &key, nil
}

// ParsePublicKey parses a base64-encoded public key
func ParsePublicKey(encoded string) (*PublicKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	if len(decoded) != KeyLength {
		return nil, fmt.Errorf("invalid key length: expected %d, got %d", KeyLength, len(decoded))
	}

	var key PublicKey
	copy(key[:], decoded)
	return &key, nil
}

// SaveToFile saves the private key to a file
func (k *PrivateKey) SaveToFile(path string) error {
	// Create file with restricted permissions (0600)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(k.String())
	if err != nil {
		return fmt.Errorf("failed to write key: %w", err)
	}

	return nil
}

// LoadPrivateKeyFromFile loads a private key from a file
func LoadPrivateKeyFromFile(path string) (*PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	return ParsePrivateKey(string(data))
}

// LoadOrGeneratePrivateKey loads a private key from file or generates a new one
func LoadOrGeneratePrivateKey(path string) (*PrivateKey, error) {
	// Try to load existing key
	key, err := LoadPrivateKeyFromFile(path)
	if err == nil {
		return key, nil
	}

	// Check if file doesn't exist (unwrap the error)
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		key, err := GeneratePrivateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate key: %w", err)
		}

		if err := key.SaveToFile(path); err != nil {
			return nil, fmt.Errorf("failed to save key: %w", err)
		}

		return key, nil
	}

	return nil, err
}
