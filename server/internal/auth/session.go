package auth

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// HostSession representa uma sessão de host para uma sala
type HostSession struct {
	RoomID    uuid.UUID `json:"room_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionManager gerencia as sessões de host
type SessionManager struct {
	sessions  map[string]*HostSession // token -> session
	roomHosts map[string]string       // room_id -> token
	mu        sync.RWMutex
}

// NewSessionManager cria um novo gerenciador de sessões
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions:  make(map[string]*HostSession),
		roomHosts: make(map[string]string),
	}

	// Limpeza automática de sessões expiradas (a cada 30 minutos)
	go sm.cleanupExpiredSessions()

	return sm
}

// CreateHostSession cria uma nova sessão de host para uma sala
func (sm *SessionManager) CreateHostSession(roomID uuid.UUID) *HostSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	token := uuid.New().String()
	now := time.Now()

	session := &HostSession{
		RoomID:    roomID,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour), // Token válido por 24 horas
	}

	sm.sessions[token] = session
	sm.roomHosts[roomID.String()] = token

	return session
}

// ValidateHostToken verifica se o token é válido para a sala
func (sm *SessionManager) ValidateHostToken(roomID uuid.UUID, token string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[token]
	if !exists {
		return false
	}

	// Verifica se não expirou
	if time.Now().After(session.ExpiresAt) {
		return false
	}

	// Verifica se o token corresponde à sala
	return session.RoomID == roomID
}

// IsRoomHost verifica se um token representa o host de uma sala específica
func (sm *SessionManager) IsRoomHost(roomID uuid.UUID, token string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	expectedToken, exists := sm.roomHosts[roomID.String()]
	if !exists {
		return false
	}

	return expectedToken == token && sm.ValidateHostToken(roomID, token)
}

// GetHostToken retorna o token do host de uma sala (se existir)
func (sm *SessionManager) GetHostToken(roomID uuid.UUID) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	token, exists := sm.roomHosts[roomID.String()]
	if !exists {
		return "", false
	}

	// Verifica se ainda é válido
	if sm.ValidateHostToken(roomID, token) {
		return token, true
	}

	return "", false
}

// RevokeHostSession revoga a sessão de host de uma sala
func (sm *SessionManager) RevokeHostSession(roomID uuid.UUID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	token, exists := sm.roomHosts[roomID.String()]
	if exists {
		delete(sm.sessions, token)
		delete(sm.roomHosts, roomID.String())
	}
}

// cleanupExpiredSessions remove sessões expiradas periodicamente
func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()

		var expiredTokens []string
		for token, session := range sm.sessions {
			if now.After(session.ExpiresAt) {
				expiredTokens = append(expiredTokens, token)
			}
		}

		for _, token := range expiredTokens {
			session := sm.sessions[token]
			delete(sm.sessions, token)
			delete(sm.roomHosts, session.RoomID.String())
		}

		sm.mu.Unlock()
	}
}

// GetSessionInfo retorna informações sobre uma sessão
func (sm *SessionManager) GetSessionInfo(token string) (*HostSession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[token]
	if !exists {
		return nil, false
	}

	// Verifica se não expirou
	if time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}
