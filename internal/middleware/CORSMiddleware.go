package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
)

// CORSMiddleware valida a origem e permite apenas o FRONTEND_URL.
// Também habilita cookies (credentials) via Same-Origin.
func CORSMiddleware(next http.Handler) http.Handler {
	allowedOrigins := []string{
		"http://localhost:30004",
		"http://localhost:33333/api/v1/auth/google/login",
		"http://localhost:33333/api/v1/auth/google/callback",
		"http://localhost:33333",
	}

	ctx := context.Background()
	logger.Default.Info(ctx, "CORS origins are "+strings.Join(allowedOrigins, ","))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		origin := r.Header.Get("Origin")

		logger.Default.Info(ctx, "CORS origin is "+origin)

		// Se não for CORS, simplesmente segue
		if origin == "" {
			logger.Default.Info(ctx, "CORS origin is empty")
			next.ServeHTTP(w, r)
			return
		}

		// Valida whitelist
		allowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}
		logger.Default.Info(ctx, "CORS origin allowed: "+origin, "allowed", allowed)
		if !allowed {
			logger.Default.Error(ctx, "CORS origin not allowed", "origin", origin)
			http.Error(w, "Origin não permitida", http.StatusForbidden)
			return
		}

		// Seta headers CORS
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Só trata preflight se for CORS
		if r.Method == http.MethodOptions {
			logger.Default.Info(ctx, "CORS preflight request")
			w.WriteHeader(http.StatusOK)
			return
		}

		logger.Default.Info(ctx, "CORS preflight request", "origin", origin)
		next.ServeHTTP(w, r)
	})
}
