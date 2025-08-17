package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	UserSessionCookieName = "user_session"
	UserSessionDuration   = 24 * time.Hour
)

type UserSessionManager struct {
	store *pgstore.Queries
}

func NewUserSessionManager(store *pgstore.Queries) *UserSessionManager {
	return &UserSessionManager{store: store}
}

// generateSessionToken generates a cryptographically secure random session token
func (usm *UserSessionManager) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateSession creates a new user session and returns the session token
func (usm *UserSessionManager) CreateSession(r *http.Request) (string, error) {
	token, err := usm.generateSessionToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	expiresAt := time.Now().Add(UserSessionDuration)
	userAgent := r.Header.Get("User-Agent")

	_, err = usm.store.CreateUserSession(r.Context(), pgstore.CreateUserSessionParams{
		SessionToken: token,
		ExpiresAt:    pgtype.Timestamp{Time: expiresAt, Valid: true},
		UserAgent:    pgtype.Text{String: userAgent, Valid: true},
		IpAddress:    nil, // Will parse IP properly later if needed
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session in database: %w", err)
	}

	return token, nil
}

// GetSession retrieves a session by token and updates last activity
func (usm *UserSessionManager) GetSession(r *http.Request, token string) (*pgstore.GetUserSessionRow, error) {
	session, err := usm.store.GetUserSession(r.Context(), token)
	if err != nil {
		return nil, err
	}

	// Update last activity and extend expiration
	newExpiresAt := time.Now().Add(UserSessionDuration)
	err = usm.store.UpdateSessionActivity(r.Context(), pgstore.UpdateSessionActivityParams{
		SessionToken: token,
		ExpiresAt:    pgtype.Timestamp{Time: newExpiresAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update session activity: %w", err)
	}

	// Update the session object with new expiration
	session.ExpiresAt = pgtype.Timestamp{Time: newExpiresAt, Valid: true}

	return &session, nil
}

// DeleteSession deletes a session from the database
func (usm *UserSessionManager) DeleteSession(r *http.Request, token string) error {
	return usm.store.DeleteUserSession(r.Context(), token)
}

// SetSessionCookie sets the session cookie on the response
func (usm *UserSessionManager) SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     UserSessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(UserSessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// ClearSessionCookie clears the session cookie
func (usm *UserSessionManager) ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     UserSessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// GetSessionFromRequest gets the session token from the request cookie
func (usm *UserSessionManager) GetSessionFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(UserSessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// CleanExpiredSessions removes expired sessions from the database
func (usm *UserSessionManager) CleanExpiredSessions(r *http.Request) error {
	return usm.store.CleanExpiredSessions(r.Context())
}
