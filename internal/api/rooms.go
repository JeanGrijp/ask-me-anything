package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/auth"
	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/JeanGrijp/ask-me-anything/internal/middleware"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	logger.Default.Info(r.Context(), "creating new room")

	type _body struct {
		Theme string `json:"theme"`
	}
	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Default.Warn(r.Context(), "invalid JSON in create room request", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.Default.Debug(r.Context(), "creating room with theme", "theme", body.Theme)

	// Adicionar timeout para operação de banco de dados
	dbCtx, cancel := WithDatabaseTimeout(r.Context())
	defer cancel()

	roomID, err := h.q.InsertRoom(dbCtx, body.Theme)
	if err != nil {
		logger.Default.Error(r.Context(), "failed to insert room", "error", err)

		// Verificar se foi erro de timeout
		if dbCtx.Err() == context.DeadlineExceeded {
			logger.Default.Warn(r.Context(), "database operation timed out", "error", dbCtx.Err())
			// Retornar erro de timeout específico
			http.Error(w, "request timeout", http.StatusRequestTimeout)
			return
		}

		logger.Default.Error(r.Context(), "failed to insert room", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	logger.Default.Info(r.Context(), "room created successfully", "room_id", roomID.String())

	// Set the current user as the room creator
	if err := h.setRoomCreator(r, roomID); err != nil {
		logger.Default.Warn(r.Context(), "failed to set room creator", "room_id", roomID.String(), "error", err)
		// Continue execution - this is not a critical failure
	}

	// Criar sessão de host para o criador da sala
	hostSession := h.sessionMgr.CreateHostSession(roomID)

	logger.Default.Info(r.Context(), "host session created", "room_id", roomID.String(), "token", hostSession.Token[:8]+"...")

	type response struct {
		ID        string `json:"id"`
		HostToken string `json:"host_token"`
	}

	sendJSON(w, response{
		ID:        roomID.String(),
		HostToken: hostSession.Token,
	})
}

func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request) {
	logger.Default.Debug(r.Context(), "fetching all rooms")

	rooms, err := h.q.GetRooms(r.Context())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		logger.Default.Error(r.Context(), "failed to get rooms", "error", err)
		return
	}

	if rooms == nil {
		rooms = []pgstore.Room{}
	}

	logger.Default.Debug(r.Context(), "rooms fetched successfully", "count", len(rooms))
	sendJSON(w, rooms)
}

func (h apiHandler) handleGetRoom(w http.ResponseWriter, r *http.Request) {
	room, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	logger.Default.Debug(r.Context(), "fetching room details", "room_id", rawRoomID)
	sendJSON(w, room)
}

func (h apiHandler) handleGetHostStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isHost := auth.IsHost(ctx)

	rawRoomID := chi.URLParam(r, "room_id")
	logger.Default.Debug(ctx, "checking host status", "room_id", rawRoomID, "is_host", isHost)

	type response struct {
		IsHost bool   `json:"is_host"`
		RoomID string `json:"room_id"`
	}

	sendJSON(w, response{
		IsHost: isHost,
		RoomID: rawRoomID,
	})
}

func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, roomID, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	logger.Default.Info(r.Context(), "creating new message", "room_id", rawRoomID)

	type _body struct {
		Message string `json:"message"`
	}
	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Default.Warn(r.Context(), "invalid JSON in create message request", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.Default.Debug(r.Context(), "creating message", "room_id", rawRoomID, "message_length", len(body.Message))

	messageID, err := h.q.InsertMessage(r.Context(), pgstore.InsertMessageParams{RoomID: roomID, Message: body.Message})
	if err != nil {
		logger.Default.Error(r.Context(), "failed to insert message", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	logger.Default.Info(r.Context(), "message created successfully", "room_id", rawRoomID, "message_id", messageID.String())

	type response struct {
		ID string `json:"id"`
	}

	sendJSON(w, response{ID: messageID.String()})

	go h.notifyClients(Message{
		Kind:   MessageKindMessageCreated,
		RoomID: rawRoomID,
		Value: MessageMessageCreated{
			ID:      messageID.String(),
			Message: body.Message,
		},
	})
}

func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, roomID, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	logger.Default.Debug(r.Context(), "fetching messages for room", "room_id", rawRoomID)

	// Get session token to check user reactions
	sessionToken, hasSession := middleware.GetUserSessionToken(r.Context())

	if hasSession {
		// Use enhanced query with user reaction info
		messages, err := h.q.GetRoomMessagesWithUserReactions(r.Context(), pgstore.GetRoomMessagesWithUserReactionsParams{
			RoomID:       roomID,
			SessionToken: sessionToken,
		})
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			logger.Default.Error(r.Context(), "failed to get room messages with user reactions", "room_id", rawRoomID, "error", err)
			return
		}

		if messages == nil {
			messages = []pgstore.GetRoomMessagesWithUserReactionsRow{}
		}

		logger.Default.Debug(r.Context(), "messages with user reactions fetched successfully", "room_id", rawRoomID, "count", len(messages))
		sendJSON(w, messages)
	} else {
		// Fallback to basic query without user reactions
		messages, err := h.q.GetRoomMessages(r.Context(), roomID)
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			logger.Default.Error(r.Context(), "failed to get room messages", "room_id", rawRoomID, "error", err)
			return
		}

		if messages == nil {
			messages = []pgstore.Message{}
		}

		logger.Default.Debug(r.Context(), "basic messages fetched successfully", "room_id", rawRoomID, "count", len(messages))
		sendJSON(w, messages)
	}
}

func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawMessageID := chi.URLParam(r, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		logger.Default.Warn(r.Context(), "invalid message ID in get message request", "room_id", rawRoomID, "message_id", rawMessageID, "error", err)
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	logger.Default.Debug(r.Context(), "fetching specific message", "room_id", rawRoomID, "message_id", rawMessageID)

	messages, err := h.q.GetMessage(r.Context(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Default.Warn(r.Context(), "message not found", "room_id", rawRoomID, "message_id", rawMessageID)
			http.Error(w, "message not found", http.StatusBadRequest)
			return
		}

		logger.Default.Error(r.Context(), "failed to get message", "room_id", rawRoomID, "message_id", rawMessageID, "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	logger.Default.Debug(r.Context(), "message fetched successfully", "room_id", rawRoomID, "message_id", rawMessageID)
	sendJSON(w, messages)
}
