package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/JeanGrijp/ask-me-anything/internal/models"
	"github.com/JeanGrijp/ask-me-anything/internal/responses"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
	TenantIDKey contextKey = "tenant_id"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthClaims struct {
	UserID   int64           `json:"user_id"`
	TenantID int64           `json:"tenant_id"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1) tenta Bearer <token> no header
		auth := r.Header.Get("Authorization")
		parts := strings.Fields(auth)
		logger.Default.Info(ctx, "Parsing authorization header", "parts", parts, "auth", auth)
		var tokenString string
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			logger.Default.Info(ctx, "Parsing bearer token")
			tokenString = parts[1]

		}

		// 2) se não veio no header, tenta cookie "jwt_token"
		if tokenString == "" {
			logger.Default.Info(ctx, "Trying to get token from cookie")
			if c, err := r.Cookie("jwt_token"); err == nil {
				tokenString = c.Value
				logger.Default.Info(ctx, "Found cookie jwt_token", "token", tokenString)
			}
		}

		if tokenString == "" {
			logger.Default.Warn(ctx, "Missing or invalid token")
			unauthorized(ctx, w, "Missing or invalid token")
			return
		}

		logger.Default.Info(ctx, "Parsing token")
		// 3) parse + validação do JWT como você já fez
		var claims AuthClaims
		parser := jwt.Parser{}
		token, err := parser.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Default.Error(ctx, "Invalid signing method in token")
				return nil, jwt.ErrSignatureInvalid
			}
			logger.Default.Info(ctx, "Parsing token claims", "claims", claims)
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			logger.Default.Warn(ctx, "Invalid token", "error", err)
			unauthorized(ctx, w, "Invalid token")
			return
		}

		// 4) checa expiração / role / etc...
		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			logger.Default.Warn(ctx, "Token expired")
			unauthorized(ctx, w, "Token expired")
			return
		}
		if !claims.Role.IsValid() {
			logger.Default.Warn(ctx, "Invalid role in token")
			unauthorized(ctx, w, "Invalid role")
			return
		}

		// 5) injeta no contexto e chama next
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, TenantIDKey, claims.TenantID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		logger.Default.Info(ctx, "Valid token", "claims", claims, "user_id", claims.UserID, "tenant_id", claims.TenantID, "user_id_key", UserIDKey, "tenant_id_key", TenantIDKey, "role_key", UserRoleKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/* ---- Helper para respostas 401 em JSON ---- */

func unauthorized(ctx context.Context, w http.ResponseWriter, msg string) {
	logger.Default.Warn(ctx, msg)
	responses.JSON(w, http.StatusUnauthorized, responses.NewErrorResponse(msg))
}
