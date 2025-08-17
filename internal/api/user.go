package api

import (
	"encoding/json"
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/middleware"
	"github.com/JeanGrijp/ask-me-anything/internal/responses"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/google/uuid"
)

// handleUserLogout logs out the current user by clearing their session
func (h apiHandler) handleUserLogout(w http.ResponseWriter, r *http.Request) {
	// Get session token from context
	sessionToken, ok := middleware.GetUserSessionToken(r.Context())
	if !ok {
		responses.SendError(w, http.StatusUnauthorized, "No active session")
		return
	}

	// Delete session from database
	if err := h.userSessionMgr.DeleteSession(r, sessionToken); err != nil {
		responses.SendError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Clear session cookie
	h.userSessionMgr.ClearSessionCookie(w)

	// Clean up expired sessions (housekeeping)
	h.userSessionMgr.CleanExpiredSessions(r)

	w.WriteHeader(http.StatusNoContent)
}

// UserRoomResponse represents a room created by the user
type UserRoomResponse struct {
	ID        string `json:"id"`
	Theme     string `json:"theme"`
	CreatedAt string `json:"created_at"`
}

// handleGetUserRooms returns all rooms created by the current user
func (h apiHandler) handleGetUserRooms(w http.ResponseWriter, r *http.Request) {
	// Get session token from context
	sessionToken, ok := middleware.GetUserSessionToken(r.Context())
	if !ok {
		responses.SendError(w, http.StatusUnauthorized, "No active session")
		return
	}

	// Get user's rooms from database
	rooms, err := h.q.GetUserRooms(r.Context(), sessionToken)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, "Failed to get user rooms")
		return
	}

	// Convert to response format
	var userRooms []UserRoomResponse
	for _, room := range rooms {
		userRooms = append(userRooms, UserRoomResponse{
			ID:        room.ID.String(),
			Theme:     room.Theme,
			CreatedAt: room.CreatedAt.Time.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRooms)
}

// setRoomCreator sets the current user as the creator of a room
func (h apiHandler) setRoomCreator(r *http.Request, roomID uuid.UUID) error {
	// Get session from context
	session, ok := middleware.GetUserSessionFromContext(r.Context())
	if !ok {
		return nil // No session, skip setting creator
	}

	// Set room creator in database
	return h.q.SetRoomCreator(r.Context(), pgstore.SetRoomCreatorParams{
		RoomID:           roomID,
		CreatorSessionID: session.ID,
	})
}

// isRoomCreator checks if the current user is the creator of a room
func (h apiHandler) isRoomCreator(r *http.Request, roomID uuid.UUID) (bool, error) {
	// Get session token from context
	sessionToken, ok := middleware.GetUserSessionToken(r.Context())
	if !ok {
		return false, nil
	}

	// Check if user is room creator
	isCreator, err := h.q.IsRoomCreator(r.Context(), pgstore.IsRoomCreatorParams{
		RoomID:       roomID,
		SessionToken: sessionToken,
	})
	if err != nil {
		return false, err
	}

	return isCreator, nil
}
