package email

import (
	"fmt"
	"log/slog"
)

// EmailService handles sending emails
type EmailService struct {
	logger *slog.Logger
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		logger: slog.Default(),
	}
}

// SendMagicLink sends a magic link email to the user
func (e *EmailService) SendMagicLink(email, token string) error {
	// For now, just log the magic link
	// In production, integrate with SendGrid, AWS SES, etc.

	magicLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)

	e.logger.Info("Magic link generated",
		"email", email,
		"token", token,
		"link", magicLink,
	)

	// TODO: Implement actual email sending
	fmt.Printf("\n🔗 Magic Link for %s:\n%s\n\n", email, magicLink)

	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (e *EmailService) SendWelcomeEmail(email, username string) error {
	e.logger.Info("Welcome email sent", "email", email, "username", username)
	fmt.Printf("👋 Welcome %s! Check your email at %s\n", username, email)
	return nil
}
