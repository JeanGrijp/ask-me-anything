package routes

import (
	wsHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/websocket"
	"github.com/go-chi/chi/v5"
)

// SetupWebSocketRoutes configures WebSocket routes
func SetupWebSocketRoutes(r *chi.Mux, handler *wsHandlers.Handler) {
	r.HandleFunc("/ws", handler.HandleConnection)
	// TODO: Add specific WebSocket endpoints when handlers are implemented
	// r.HandleFunc("/ws/questions", handler.HandleQuestionUpdates)
	// r.HandleFunc("/ws/answers", handler.HandleAnswerUpdates)
}
