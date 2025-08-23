package auth

import (
	"context"
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	HostTokenKey contextKey = "host_token"
	IsHostKey    contextKey = "is_host"
)

// HostOnlyMiddleware middleware que permite apenas hosts executarem a ação
func HostOnlyMiddleware(sm *SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extrair room_id da URL
			rawRoomID := chi.URLParam(r, "room_id")
			if rawRoomID == "" {
				logger.Default.Warn(ctx, "room_id not found in URL")
				http.Error(w, "room_id required", http.StatusBadRequest)
				return
			}

			roomID, err := uuid.Parse(rawRoomID)
			if err != nil {
				logger.Default.Warn(ctx, "invalid room_id format", "room_id", rawRoomID)
				http.Error(w, "invalid room_id", http.StatusBadRequest)
				return
			}

			// Extrair token do header
			token := r.Header.Get("X-Host-Token")
			if token == "" {
				logger.Default.Warn(ctx, "host token not provided", "room_id", rawRoomID)
				http.Error(w, "host token required", http.StatusUnauthorized)
				return
			}

			// Verificar se é o host da sala
			if !sm.IsRoomHost(roomID, token) {
				logger.Default.Warn(ctx, "unauthorized host action", "room_id", rawRoomID, "token", token[:8]+"...")
				http.Error(w, "only room host can perform this action", http.StatusForbidden)
				return
			}

			logger.Default.Info(ctx, "host action authorized", "room_id", rawRoomID)

			// Adicionar informações no contexto
			ctx = context.WithValue(ctx, HostTokenKey, token)
			ctx = context.WithValue(ctx, IsHostKey, true)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalHostMiddleware middleware que identifica se o usuário é host (mas não requer)
func OptionalHostMiddleware(sm *SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extrair room_id da URL
			rawRoomID := chi.URLParam(r, "room_id")
			token := r.Header.Get("X-Host-Token")

			isHost := false
			if rawRoomID != "" && token != "" {
				if roomID, err := uuid.Parse(rawRoomID); err == nil {
					isHost = sm.IsRoomHost(roomID, token)
				}
			}

			// Adicionar informações no contexto
			ctx = context.WithValue(ctx, HostTokenKey, token)
			ctx = context.WithValue(ctx, IsHostKey, isHost)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IsHost verifica se o contexto atual representa um host
func IsHost(ctx context.Context) bool {
	isHost, ok := ctx.Value(IsHostKey).(bool)
	return ok && isHost
}

// GetHostToken extrai o token do host do contexto
func GetHostToken(ctx context.Context) string {
	token, _ := ctx.Value(HostTokenKey).(string)
	return token
}
