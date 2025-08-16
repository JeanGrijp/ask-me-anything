package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/JeanGrijp/ask-me-anything/internal/auth"
	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	_middleware "github.com/JeanGrijp/ask-me-anything/internal/middleware"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

type apiHandler struct {
	q           *pgstore.Queries
	r           *chi.Mux
	upgrader    websocket.Upgrader
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	mu          *sync.Mutex
	sessionMgr  *auth.SessionManager
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	sessionMgr := auth.NewSessionManager()

	a := apiHandler{
		q:           q,
		upgrader:    websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mu:          &sync.Mutex{},
		sessionMgr:  sessionMgr,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	r.Use(_middleware.RequestIDMiddleware)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Host-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/subscribe/{room_id}", a.handleSubscribe)

	r.Route("/api", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", a.handleCreateRoom)
			r.Get("/", a.handleGetRooms)

			r.Route("/{room_id}", func(r chi.Router) {
				r.Get("/", a.handleGetRoom)

				// Rota para verificar se é host (com middleware opcional)
				r.With(auth.OptionalHostMiddleware(sessionMgr)).Get("/host-status", a.handleGetHostStatus)

				r.Route("/messages", func(r chi.Router) {
					r.Post("/", a.handleCreateRoomMessage)
					r.Get("/", a.handleGetRoomMessages)

					r.Route("/{message_id}", func(r chi.Router) {
						r.Get("/", a.handleGetRoomMessage)
						r.Patch("/react", a.handleReactToMessage)
						r.Delete("/react", a.handleRemoveReactFromMessage)

						// Apenas o host pode marcar mensagens como respondidas
						r.With(auth.HostOnlyMiddleware(sessionMgr)).Patch("/answer", a.handleMarkMessageAsAnswered)
					})
				})
			})
		})
	})

	a.r = r
	return a
}

const (
	MessageKindMessageCreated          = "message_created"
	MessageKindMessageRactionIncreased = "message_reaction_increased"
	MessageKindMessageRactionDecreased = "message_reaction_decreased"
	MessageKindMessageAnswered         = "message_answered"
)

type MessageMessageReactionIncreased struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type MessageMessageReactionDecreased struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type MessageMessageAnswered struct {
	ID string `json:"id"`
}

type MessageMessageCreated struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type Message struct {
	Kind   string `json:"kind"`
	Value  any    `json:"value"`
	RoomID string `json:"-"`
}

func (h apiHandler) notifyClients(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	subscribers, ok := h.subscribers[msg.RoomID]
	if !ok || len(subscribers) == 0 {
		logger.Default.Debug(context.Background(), "no subscribers for room", "room_id", msg.RoomID, "message_kind", msg.Kind)
		return
	}

	logger.Default.Debug(context.Background(), "notifying clients", "room_id", msg.RoomID, "message_kind", msg.Kind, "subscriber_count", len(subscribers))

	disconnectedClients := 0
	for conn, cancel := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			logger.Default.Error(context.Background(), "failed to send message to client", "room_id", msg.RoomID, "message_kind", msg.Kind, "error", err)
			cancel()
			disconnectedClients++
		}
	}

	if disconnectedClients > 0 {
		logger.Default.Warn(context.Background(), "some clients disconnected", "room_id", msg.RoomID, "disconnected_count", disconnectedClients)
	}
}

func (h apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	logger.Default.Info(r.Context(), "WebSocket connection attempt", "room_id", rawRoomID, "client_ip", r.RemoteAddr)

	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Default.Warn(r.Context(), "failed to upgrade connection", "room_id", rawRoomID, "client_ip", r.RemoteAddr, "error", err)
		http.Error(w, "failed to upgrade to ws connection", http.StatusBadRequest)
		return
	}

	defer c.Close()

	ctx, cancel := context.WithCancel(r.Context())

	h.mu.Lock()
	if _, ok := h.subscribers[rawRoomID]; !ok {
		h.subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	logger.Default.Info(r.Context(), "new client connected", "room_id", rawRoomID, "client_ip", r.RemoteAddr, "total_subscribers", len(h.subscribers[rawRoomID])+1)
	h.subscribers[rawRoomID][c] = cancel
	h.mu.Unlock()

	<-ctx.Done()

	h.mu.Lock()
	delete(h.subscribers[rawRoomID], c)
	remainingSubscribers := len(h.subscribers[rawRoomID])
	h.mu.Unlock()

	logger.Default.Info(context.Background(), "client disconnected", "room_id", rawRoomID, "client_ip", r.RemoteAddr, "remaining_subscribers", remainingSubscribers)
}
