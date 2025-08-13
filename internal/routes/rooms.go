package routes

import (
	"github.com/go-chi/chi/v5"

	httpHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/http"
)

// SetupRoomRoutes configures room-related routes
func SetupRoomRoutes(handler *httpHandlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Room routes
	r.Get("/", handler.ListRooms)                        // GET /api/v1/rooms
	r.Post("/", handler.CreateRoom)                      // POST /api/v1/rooms
	r.Get("/{id}", handler.GetRoom)                      // GET /api/v1/rooms/{id}
	r.Put("/{id}", handler.UpdateRoom)                   // PUT /api/v1/rooms/{id}
	r.Delete("/{id}", handler.DeleteRoom)                // DELETE /api/v1/rooms/{id}
	r.Get("/owner/{owner_id}", handler.ListRoomsByOwner) // GET /api/v1/rooms/owner/{owner_id}

	return r
}
