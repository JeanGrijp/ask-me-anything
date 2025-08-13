package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// MagicLinkService handles magic link authentication
type MagicLinkService struct {
	// Will be injected later
}

// GenerateToken creates a secure random token for magic links
func (s *MagicLinkService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GetExpirationTime returns the expiration time for magic links (15 minutes from now)
func (s *MagicLinkService) GetExpirationTime() time.Time {
	return time.Now().Add(15 * time.Minute)
}

// IsValidEmail basic email validation
func (s *MagicLinkService) IsValidEmail(email string) bool {
	// Basic email validation - in production use a proper library
	return len(email) > 3 && len(email) < 255
}

// CreateMagicLink creates a new magic link for the given email
func (s *MagicLinkService) CreateMagicLink(email string) (string, time.Time, error) {
	if !s.IsValidEmail(email) {
		return "", time.Time{}, fmt.Errorf("invalid email format")
	}

	token, err := s.GenerateToken()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := s.GetExpirationTime()

	return token, expiresAt, nil
}
