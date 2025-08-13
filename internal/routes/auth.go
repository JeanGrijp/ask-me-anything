package routes

import (
	authHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/auth"
	"github.com/go-chi/chi/v5"
)

// SetupAuthRoutes configures all authentication routes
func SetupAuthRoutes(handler *authHandlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/login", handler.RequestMagicLink)
	r.Post("/verify", handler.VerifyMagicLink)
	r.Get("/verify", handler.VerifyMagicLinkGet) // For email clicks
	r.Post("/logout", handler.Logout)
	r.Get("/me", handler.GetCurrentUser)

	return r
}
