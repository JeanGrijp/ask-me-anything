package routes

import (
	httpHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/http"
	"github.com/go-chi/chi/v5"
)

// SetupAnswerRoutes configures all answer-related routes
func SetupAnswerRoutes(handler *httpHandlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.ListAnswers)
	r.Post("/", handler.CreateAnswer)
	r.Get("/{id}", handler.GetAnswer)
	r.Put("/{id}", handler.UpdateAnswer)
	r.Delete("/{id}", handler.DeleteAnswer)

	// Additional answer routes
	r.Post("/{id}/vote", handler.VoteAnswer)
	r.Post("/{id}/accept", handler.AcceptAnswer)

	return r
}
