package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"golang.org/x/time/rate"
)

type UserLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	users     = make(map[string]*UserLimiter)
	usersLock sync.Mutex
	rateLimit = rate.Every(1 * time.Second)
	burst     = 5
)

func getUserLimiter(userID string) *rate.Limiter {
	usersLock.Lock()
	defer usersLock.Unlock()

	if ul, exists := users[userID]; exists {
		ul.lastSeen = time.Now().UTC()
		return ul.limiter
	}

	limiter := rate.NewLimiter(rateLimit, burst)
	users[userID] = &UserLimiter{limiter, time.Now()}
	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr // fallback
		}
		identifier := host

		userIDRaw := ctx.Value(UserIDKey)
		userIDInt, ok := userIDRaw.(int64)
		if ok && userIDInt != 0 {
			identifier = fmt.Sprintf("%d", userIDInt)
		}

		limiter := getUserLimiter(identifier)
		if !limiter.Allow() {
			logger.Default.Warn(ctx, "Rate limit exceeded", "identifier", identifier, "remote_addr", r.RemoteAddr, "method", r.Method, "path", r.URL.Path)
			http.Error(w, "Request limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CleanupUsers(ctx context.Context) {
	logger.Default.Info(ctx, "Starting user cleanup routine")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			usersLock.Lock()
			for userID, ul := range users {
				if time.Since(ul.lastSeen) > 5*time.Minute {
					delete(users, userID)
				}
			}
			usersLock.Unlock()
		case <-ctx.Done():
			logger.Default.Info(ctx, "User cleanup routine stopped")
			return
		}
	}
}
