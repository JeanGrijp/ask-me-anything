package routes

import (
	httpHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/http"
	"github.com/go-chi/chi/v5"
)

// SetupQuestionRoutes configures all question-related routes
func SetupQuestionRoutes(handler *httpHandlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.ListQuestions)
	r.Post("/", handler.CreateQuestion)
	r.Get("/{id}", handler.GetQuestion)
	r.Put("/{id}", handler.UpdateQuestion)
	r.Delete("/{id}", handler.DeleteQuestion)

	// Additional question routes
	r.Get("/{id}/answers", handler.GetQuestionAnswers)
	r.Post("/{id}/vote", handler.VoteQuestion)
	r.Get("/search", handler.SearchQuestions)
	r.Get("/by-category/{categoryId}", handler.GetQuestionsByCategory)
	r.Get("/by-user/{userId}", handler.GetQuestionsByUser)

	return r
}
