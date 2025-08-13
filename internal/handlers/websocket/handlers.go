package websocket

import (
	"log/slog"
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/database"
	"github.com/gorilla/websocket"
)

type Handler struct {
	db       *database.DB
	logger   *slog.Logger
	upgrader websocket.Upgrader
}

func New(db *database.DB) *Handler {
	return &Handler{
		db:     db,
		logger: slog.Default(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", "error", err)
		return
	}
	defer conn.Close()

	h.logger.Info("WebSocket connection established")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("Failed to read message", "error", err)
			break
		}

		h.logger.Info("Received message", "type", messageType, "message", string(message))

		// Echo the message back
		if err := conn.WriteMessage(messageType, message); err != nil {
			h.logger.Error("Failed to write message", "error", err)
			break
		}
	}
}
