package service

// AuthService handles authentication (placeholder for future expansion)
type AuthService struct {
	apiKey string
}

// NewAuthService creates a new auth service
func NewAuthService(apiKey string) *AuthService {
	return &AuthService{
		apiKey: apiKey,
	}
}

// ValidateAPIKey validates an API key (optional for MVP)
func (s *AuthService) ValidateAPIKey(key string) bool {
	// If no API key is configured, allow all requests
	if s.apiKey == "" {
		return true
	}
	return key == s.apiKey
}
