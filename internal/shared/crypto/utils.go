package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// GenerateID generates a unique peer ID
func GenerateID() string {
	return uuid.New().String()
}

// EncodeKey encodes a byte slice to base64 string
func EncodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// DecodeKey decodes a base64 string to byte slice
func DecodeKey(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// ValidatePublicKey validates a WireGuard public key format
func ValidatePublicKey(key string) error {
	decoded, err := DecodeKey(key)
	if err != nil {
		return fmt.Errorf("invalid base64 encoding: %w", err)
	}
	if len(decoded) != 32 {
		return fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(decoded))
	}
	return nil
}

// GenerateRandomBytes generates n random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
