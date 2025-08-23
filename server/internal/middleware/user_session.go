package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/auth"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
)

type userSessionContextKey string

const (
	UserSessionContextKey userSessionContextKey = "user_session"
)

// UserSessionMiddleware automatically manages user sessions via cookies
func UserSessionMiddleware(sessionManager *auth.UserSessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var session *pgstore.GetUserSessionRow

			// Try to get existing session from cookie
			if token, err := sessionManager.GetSessionFromRequest(r); err == nil {
				// Validate and refresh session
				if existingSession, err := sessionManager.GetSession(r, token); err == nil {
					session = existingSession
					// Refresh cookie with updated expiration
					sessionManager.SetSessionCookie(w, token)
				}
			}

			// If no valid session exists, create a new one
			if session == nil {
				token, err := sessionManager.CreateSession(r)
				if err != nil {
					slog.Error("Failed to create user session", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Set cookie for new session
				sessionManager.SetSessionCookie(w, token)

				// Get the newly created session
				if newSession, err := sessionManager.GetSession(r, token); err == nil {
					session = newSession
				}
			}

			// Add session to context if we have one
			if session != nil {
				ctx := context.WithValue(r.Context(), UserSessionContextKey, session)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserSessionFromContext retrieves the user session from the request context
func GetUserSessionFromContext(ctx context.Context) (*pgstore.GetUserSessionRow, bool) {
	session, ok := ctx.Value(UserSessionContextKey).(*pgstore.GetUserSessionRow)
	return session, ok
}

// GetUserSessionID retrieves just the session ID from the request context
func GetUserSessionID(ctx context.Context) (string, bool) {
	if session, ok := GetUserSessionFromContext(ctx); ok {
		return session.ID.String(), true
	}
	return "", false
}

// GetUserSessionToken retrieves just the session token from the request context
func GetUserSessionToken(ctx context.Context) (string, bool) {
	if session, ok := GetUserSessionFromContext(ctx); ok {
		return session.SessionToken, true
	}
	return "", false
}
