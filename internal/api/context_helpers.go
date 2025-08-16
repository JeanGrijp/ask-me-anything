// Package api provides context helpers and timeout utilities for API operations.
package api

import (
	"context"
	"time"
)

const (
	// Timeouts para diferentes operações
	DatabaseOperationTimeout  = 5 * time.Second
	WebSocketReadTimeout      = 30 * time.Second
	WebSocketWriteTimeout     = 10 * time.Second
	ClientNotificationTimeout = 2 * time.Second
)

// WithDatabaseTimeout cria um context com timeout para operações de banco de dados
func WithDatabaseTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DatabaseOperationTimeout)
}

// WithClientTimeout cria um context com timeout para operações do cliente
func WithClientTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, ClientNotificationTimeout)
}

// contextValue define chaves para valores no context
type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	RoomIDKey    contextKey = "room_id"
)

// GetRequestID extrai o request ID do context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithRequestID adiciona um request ID ao context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithRoomID adiciona um room ID ao context
func WithRoomID(ctx context.Context, roomID string) context.Context {
	return context.WithValue(ctx, RoomIDKey, roomID)
}

// GetRoomID extrai o room ID do context
func GetRoomID(ctx context.Context) string {
	if id, ok := ctx.Value(RoomIDKey).(string); ok {
		return id
	}
	return ""
}
