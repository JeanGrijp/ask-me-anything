// Package middleware provides HTTP middleware functions for the API.
package middleware

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const (
	RequestStartTimeKey contextKey = "request_start_time"
	ClientIPKey         contextKey = "client_ip"
	UserAgentKey        contextKey = "user_agent"
)

const DefaultRequestTimeout = 30 * time.Second

// TimeoutMiddleware adiciona um timeout global para todas as requisições HTTP
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Canal para capturar quando a requisição termina
			done := make(chan struct{})

			// Executar handler em goroutine
			go func() {
				defer close(done)
				next.ServeHTTP(w, r.WithContext(ctx))
			}()

			select {
			case <-done:
				// Requisição terminou normalmente
				return
			case <-ctx.Done():
				// Timeout ou cancelamento
				if ctx.Err() == context.DeadlineExceeded {
					http.Error(w, "Request timeout", http.StatusRequestTimeout)
				} else {
					http.Error(w, "Request cancelled", http.StatusRequestTimeout)
				}
				return
			}
		})
	}
}

// ContextEnrichmentMiddleware adiciona informações úteis ao context
func ContextEnrichmentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Adicionar timestamp da requisição
		ctx = context.WithValue(ctx, RequestStartTimeKey, time.Now())

		// Adicionar IP do cliente
		ctx = context.WithValue(ctx, ClientIPKey, r.RemoteAddr)

		// Adicionar User-Agent
		ctx = context.WithValue(ctx, UserAgentKey, r.UserAgent())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
