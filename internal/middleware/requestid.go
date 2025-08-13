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
		ctx = logger.InjectUserAgent(ctx, r.UserAgent())
		ctx = logger.InjectMethod(ctx, r.Method)
		ctx = logger.InjectPath(ctx, r.URL.Path)
		ctx = logger.InjectQuery(ctx, r.URL.RawQuery)
		ctx = logger.InjectReferer(ctx, r.Referer())
		ctx = logger.InjectHost(ctx, r.Host)

		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		r = r.WithContext(ctx)
		next.ServeHTTP(recorder, r)

		latency := time.Since(start)
		statusCode := recorder.statusCode

		ctx = logger.InjectLatency(ctx, latency)
		ctx = logger.InjectStatusCode(ctx, statusCode)

		logger.Default.Info(ctx, "Request handled")
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
