package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		ctx = logger.InjectRequestID(ctx)
		ctx = logger.InjectClientIP(ctx, getClientIP(r))
		ctx = logger.InjectMethod(ctx, r.Method)
		ctx = logger.InjectPath(ctx, r.URL.Path)

		// Só adicionar query se não estiver vazio
		if r.URL.RawQuery != "" {
			ctx = logger.InjectQuery(ctx, r.URL.RawQuery)
		}

		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		r = r.WithContext(ctx)
		next.ServeHTTP(recorder, r)

		latency := time.Since(start)
		statusCode := recorder.statusCode

		ctx = logger.InjectLatency(ctx, latency)
		ctx = logger.InjectStatusCode(ctx, statusCode)

		// Log simplificado - apenas informações essenciais
		logFields := []interface{}{
			"method", r.Method,
			"path", r.URL.Path,
			"status", statusCode,
		}

		// Adicionar query apenas se existir
		if r.URL.RawQuery != "" {
			logFields = append(logFields, "query", r.URL.RawQuery)
		}

		// Adicionar latência apenas se for significativa (> 1ms)
		if latency > time.Millisecond {
			logFields = append(logFields, "latency", latency.String())
		}

		logger.Default.Info(ctx, "Request handled", logFields...)
	})
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
