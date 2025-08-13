package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetRoom handles GET /api/v1/rooms/{id}
func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message": "Get room",
		"status":  "success",
		"id":      roomID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ListRooms handles GET /api/v1/rooms
func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "List rooms",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ListRoomsByOwner handles GET /api/v1/rooms/owner/{owner_id}
func (h *Handler) ListRoomsByOwner(w http.ResponseWriter, r *http.Request) {
	ownerID := chi.URLParam(r, "owner_id")

	response := map[string]interface{}{
		"message":  "List rooms by owner",
		"status":   "success",
		"owner_id": ownerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateRoom handles POST /api/v1/rooms
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Room created",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateRoom handles PUT /api/v1/rooms/{id}
func (h *Handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message": "Room updated",
		"status":  "success",
		"id":      roomID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteRoom handles DELETE /api/v1/rooms/{id}
func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	response := map[string]interface{}{
		"message": "Room deleted",
		"status":  "success",
		"id":      roomID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
