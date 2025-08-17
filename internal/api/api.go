package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/auth"
	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	custommiddleware "github.com/JeanGrijp/ask-me-anything/internal/middleware"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/JeanGrijp/ask-me-anything/internal/utils"
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
		q: q,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
			// Configurações básicas para evitar problemas de hijacking
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mu:          &sync.Mutex{},
		sessionMgr:  sessionMgr,
	}

	// Router principal com middlewares
	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)
	r.Use(custommiddleware.ContextEnrichmentMiddleware)
	r.Use(custommiddleware.RequestIDMiddleware)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Host-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Host-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rota de health check para monitoramento e Docker
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Incluir informações úteis para monitoramento
		response := map[string]interface{}{
			"status":    "ok",
			"service":   "ask-me-anything",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"uptime":    time.Since(time.Now()).String(), // Placeholder - pode ser melhorado
		}

		json.NewEncoder(w).Encode(response)
	})

	// Rota adicional para informações detalhadas do sistema (opcional)
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := map[string]interface{}{
			"service":     "ask-me-anything",
			"version":     "1.0.0",
			"environment": "development", // Pode vir de variável de ambiente
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"endpoints": map[string]interface{}{
				"health":    "/health",
				"api":       "/api/rooms",
				"websocket": "/subscribe/{room_id}",
			},
		}

		json.NewEncoder(w).Encode(status)
	})

	r.Route("/api", func(r chi.Router) {
		// Aplicar timeout apenas nas rotas da API, não no WebSocket
		r.Use(custommiddleware.TimeoutMiddleware(custommiddleware.DefaultRequestTimeout))
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

	// Exibir todas as rotas registradas no terminal
	utils.LogRoutes(r)

	a.r = r

	// Retornar um handler que separa WebSocket das outras rotas
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Default.Info(req.Context(), "Request intercepted", "method", req.Method, "path", req.URL.Path)

		// Verificar se é uma rota WebSocket usando strings.HasPrefix
		if req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/subscribe/") {
			logger.Default.Info(req.Context(), "WebSocket route detected", "path", req.URL.Path)
			a.handleSubscribeRaw(w, req)
			return
		}

		logger.Default.Info(req.Context(), "Routing to normal handler", "path", req.URL.Path)
		// Para todas as outras rotas, usar o router normal
		a.r.ServeHTTP(w, req)
	})
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

	// Criar context com timeout para notificações
	ctx, cancel := WithClientTimeout(context.Background())
	defer cancel()

	disconnectedClients := 0
	for conn, cancelFunc := range subscribers {
		// Enviar mensagem com timeout
		done := make(chan error, 1)
		go func() {
			conn.SetWriteDeadline(time.Now().Add(ClientNotificationTimeout))
			done <- conn.WriteJSON(msg)
		}()

		select {
		case err := <-done:
			if err != nil {
				logger.Default.Error(context.Background(), "failed to send message to client", "room_id", msg.RoomID, "message_kind", msg.Kind, "error", err)
				cancelFunc()
				disconnectedClients++
			}
		case <-ctx.Done():
			logger.Default.Warn(context.Background(), "timeout sending message to client", "room_id", msg.RoomID, "message_kind", msg.Kind)
			cancelFunc()
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

	// Criar um context com timeout para a conexão WebSocket
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Adicionar room_id ao context para rastreamento
	ctx = WithRoomID(ctx, rawRoomID)

	h.mu.Lock()
	if _, ok := h.subscribers[rawRoomID]; !ok {
		h.subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	logger.Default.Info(ctx, "new client connected", "room_id", rawRoomID, "client_ip", r.RemoteAddr, "total_subscribers", len(h.subscribers[rawRoomID])+1)
	h.subscribers[rawRoomID][c] = cancel
	h.mu.Unlock()

	// Configurar timeouts para WebSocket
	c.SetReadDeadline(time.Now().Add(WebSocketReadTimeout))
	c.SetWriteDeadline(time.Now().Add(WebSocketWriteTimeout))

	// Configurar pong handler para manter conexão viva
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(WebSocketReadTimeout))
		return nil
	})

	// Enviar ping periodicamente para manter conexão viva
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.SetWriteDeadline(time.Now().Add(WebSocketWriteTimeout))
				if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
					logger.Default.Debug(ctx, "failed to send ping", "error", err)
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()

	h.mu.Lock()
	delete(h.subscribers[rawRoomID], c)
	remainingSubscribers := len(h.subscribers[rawRoomID])
	h.mu.Unlock()

	logger.Default.Info(context.Background(), "client disconnected", "room_id", rawRoomID, "client_ip", r.RemoteAddr, "remaining_subscribers", remainingSubscribers)
}

// handleSubscribeRaw - Handler WebSocket completo com broadcast
func (h apiHandler) handleSubscribeRaw(w http.ResponseWriter, r *http.Request) {
	// Extrair room_id da URL
	roomID := r.URL.Path[11:] // Remove "/subscribe/"
	if roomID == "" {
		http.Error(w, "room_id is required", http.StatusBadRequest)
		return
	}

	logger.Default.Info(r.Context(), "WebSocket connection attempt", "room_id", roomID, "client_ip", r.RemoteAddr)

	// Upgrader básico
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Tentar upgrade direto
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Default.Warn(r.Context(), "upgrade failed", "error", err, "room_id", roomID)
		return
	}
	defer conn.Close()

	logger.Default.Info(r.Context(), "WebSocket connected successfully", "room_id", roomID, "client_ip", r.RemoteAddr)

	// Context para gerenciar a conexão
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Registrar cliente nos subscribers
	h.mu.Lock()
	if _, ok := h.subscribers[roomID]; !ok {
		h.subscribers[roomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	h.subscribers[roomID][conn] = cancel
	subscriberCount := len(h.subscribers[roomID])
	h.mu.Unlock()

	logger.Default.Info(ctx, "client registered", "room_id", roomID, "total_subscribers", subscriberCount)

	// Configurar timeouts e handlers para manter conexão viva
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Goroutine para enviar pings
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					logger.Default.Debug(ctx, "failed to send ping", "error", err, "room_id", roomID)
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Aguardar até a conexão ser fechada
	<-ctx.Done()

	// Limpar cliente dos subscribers
	h.mu.Lock()
	delete(h.subscribers[roomID], conn)
	remainingSubscribers := len(h.subscribers[roomID])
	h.mu.Unlock()

	logger.Default.Info(context.Background(), "client disconnected", "room_id", roomID, "client_ip", r.RemoteAddr, "remaining_subscribers", remainingSubscribers)
}
