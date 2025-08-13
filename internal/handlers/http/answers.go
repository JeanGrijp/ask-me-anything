package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListAnswers handles GET /api/v1/answers
func (h *Handler) ListAnswers(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "List answers",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateAnswer handles POST /api/v1/answers
func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Answer created",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAnswer handles GET /api/v1/answers/{id}
func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	answerID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message": "Get answer",
		"status":  "success",
		"id":      answerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateAnswer handles PUT /api/v1/answers/{id}
func (h *Handler) UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	answerID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message": "Answer updated",
		"status":  "success",
		"id":      answerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteAnswer handles DELETE /api/v1/answers/{id}
func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	answerID := chi.URLParam(r, "id")

	h.logger.Info("Answer deleted", "id", answerID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// VoteAnswer handles POST /api/v1/answers/{id}/vote
func (h *Handler) VoteAnswer(w http.ResponseWriter, r *http.Request) {
	answerID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message":   "Vote registered",
		"status":    "success",
		"answer_id": answerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// AcceptAnswer handles POST /api/v1/answers/{id}/accept
func (h *Handler) AcceptAnswer(w http.ResponseWriter, r *http.Request) {
	answerID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message":   "Answer accepted",
		"status":    "success",
		"answer_id": answerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
