package api

import (
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/JeanGrijp/ask-me-anything/internal/middleware"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, roomID, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawID := chi.URLParam(r, "message_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		logger.Default.Warn(r.Context(), "invalid message ID in react request", "message_id", rawID, "error", err)
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	logger.Default.Debug(r.Context(), "adding reaction to message", "room_id", rawRoomID, "message_id", rawID)

	// Get user session for tracking reactions
	session, hasSession := middleware.GetUserSessionFromContext(r.Context())
	if !hasSession {
		logger.Default.Warn(r.Context(), "no user session found for reaction", "room_id", rawRoomID, "message_id", rawID)
		http.Error(w, "session required", http.StatusUnauthorized)
		return
	}

	// Add user reaction to tracking table
	err = h.q.AddUserReaction(r.Context(), pgstore.AddUserReactionParams{
		SessionID:    session.ID,
		RoomID:       roomID,
		MessageID:    id,
		ReactionType: "like", // For now, we only support "like" reactions
	})
	if err != nil {
		logger.Default.Error(r.Context(), "failed to add user reaction", "error", err)
		// Continue anyway - this is not critical for the reaction count
	}

	// Increment reaction count
	count, err := h.q.ReactToMessage(r.Context(), id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		logger.Default.Error(r.Context(), "failed to react to message", "error", err)
		return
	}

	logger.Default.Info(r.Context(), "reaction added successfully", "room_id", rawRoomID, "message_id", rawID, "new_count", count)

	type response struct {
		Count int64 `json:"count"`
	}

	sendJSON(w, response{Count: count})

	go h.notifyClients(Message{
		Kind:   MessageKindMessageRactionIncreased,
		RoomID: rawRoomID,
		Value: MessageMessageReactionIncreased{
			ID:    rawID,
			Count: count,
		},
	})
}

func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawID := chi.URLParam(r, "message_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		logger.Default.Warn(r.Context(), "invalid message ID in remove reaction request", "message_id", rawID, "error", err)
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	logger.Default.Debug(r.Context(), "removing reaction from message", "room_id", rawRoomID, "message_id", rawID)

	count, err := h.q.RemoveReactionFromMessage(r.Context(), id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		logger.Default.Error(r.Context(), "failed to remove reaction from message", "room_id", rawRoomID, "message_id", rawID, "error", err)
		return
	}

	logger.Default.Info(r.Context(), "reaction removed successfully", "room_id", rawRoomID, "message_id", rawID, "new_count", count)

	type response struct {
		Count int64 `json:"count"`
	}

	sendJSON(w, response{Count: count})

	go h.notifyClients(Message{
		Kind:   MessageKindMessageRactionDecreased,
		RoomID: rawRoomID,
		Value: MessageMessageReactionDecreased{
			ID:    rawID,
			Count: count,
		},
	})
}

func (h apiHandler) handleMarkMessageAsAnswered(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawID := chi.URLParam(r, "message_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		logger.Default.Warn(r.Context(), "invalid message ID in mark answered request", "message_id", rawID, "error", err)
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	logger.Default.Info(r.Context(), "marking message as answered", "room_id", rawRoomID, "message_id", rawID)

	err = h.q.MarkMessageAsAnswered(r.Context(), id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		logger.Default.Error(r.Context(), "failed to mark message as answered", "room_id", rawRoomID, "message_id", rawID, "error", err)
		return
	}

	logger.Default.Info(r.Context(), "message marked as answered successfully", "room_id", rawRoomID, "message_id", rawID)

	w.WriteHeader(http.StatusOK)

	go h.notifyClients(Message{
		Kind:   MessageKindMessageAnswered,
		RoomID: rawRoomID,
		Value: MessageMessageAnswered{
			ID: rawID,
		},
	})
}
