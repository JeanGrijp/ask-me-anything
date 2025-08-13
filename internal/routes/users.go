package routes

import (
	httpHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/http"
	"github.com/go-chi/chi/v5"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(handler *httpHandlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.ListUsers)
	r.Post("/", handler.CreateUser)
	r.Get("/{id}", handler.GetUser)
	r.Put("/{id}", handler.UpdateUser)
	r.Delete("/{id}", handler.DeleteUser)

	// Additional user routes
	r.Get("/{id}/profile", handler.GetUserProfile)
	r.Get("/{id}/questions", handler.GetUserQuestions)
	r.Get("/{id}/answers", handler.GetUserAnswers)

	return r
}
