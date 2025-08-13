// Package routes provides route configuration for the ask-me-anything application
package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/JeanGrijp/ask-me-anything/internal/database"
	authHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/auth"
	httpHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/http"
	wsHandlers "github.com/JeanGrijp/ask-me-anything/internal/handlers/websocket"
	middlewareInternal "github.com/JeanGrijp/ask-me-anything/internal/middleware"
)

// Router holds all route configurations
type Router struct {
	db *database.DB
}

// NewRouter creates a new router instance
func NewRouter(db *database.DB) *Router {
	return &Router{
		db: db,
	}
}

// Setup configures all routes and middleware
func (rt *Router) Setup() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	rt.setupMiddleware(r)

	// Initialize handlers
	httpHandler := httpHandlers.New(rt.db)
	authHandler := authHandlers.NewAuthHandler(rt.db)
	wsHandler := wsHandlers.New(rt.db)

	// Setup route groups
	rt.setupAPIRoutes(r, httpHandler)
	rt.setupAuthRoutes(r, authHandler)
	rt.setupWebSocketRoutes(r, wsHandler)
	rt.setupHealthRoutes(r)

	return r
}

// setupMiddleware configures global middleware
func (rt *Router) setupMiddleware(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewareInternal.CORSMiddleware)
	r.Use(middlewareInternal.RateLimitMiddleware)
	r.Use(middlewareInternal.CSPMiddleware)
	r.Use(middlewareInternal.RequestIDMiddleware)

}

// setupAPIRoutes configures API v1 routes
func (rt *Router) setupAPIRoutes(r *chi.Mux, handler *httpHandlers.Handler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/questions", SetupQuestionRoutes(handler))
		r.Mount("/answers", SetupAnswerRoutes(handler))
		r.Mount("/users", SetupUserRoutes(handler))
		r.Mount("/rooms", SetupRoomRoutes(handler))
		// TODO: Add these routes when handlers are implemented
		// r.Mount("/categories", SetupCategoryRoutes(handler))
		// r.Mount("/votes", SetupVoteRoutes(handler))
	})
}

// setupAuthRoutes configures authentication routes
func (rt *Router) setupAuthRoutes(r *chi.Mux, handler *authHandlers.Handler) {
	r.Mount("/auth", SetupAuthRoutes(handler))
}

// setupWebSocketRoutes configures WebSocket routes
func (rt *Router) setupWebSocketRoutes(r *chi.Mux, handler *wsHandlers.Handler) {
	SetupWebSocketRoutes(r, handler)
}

// setupHealthRoutes configures health check routes
func (rt *Router) setupHealthRoutes(r *chi.Mux) {
	r.Get("/health", rt.healthCheck)
	r.Get("/ping", rt.pingCheck)
}
