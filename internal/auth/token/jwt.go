// internal/auth/token/jwt.go
package tokenAuth

import (
	"errors"
	"os"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

/* ------------------------------------------------------------------
   Claims comuns
-------------------------------------------------------------------*/

type fullClaims struct {
	UserID int64 `json:"user_id"`

	Role     string `json:"role,omitempty"`
	LoginTok bool   `json:"lt,omitempty"` // true nos login-tokens
	jwt.RegisteredClaims
}

/* ------------------------------------------------------------------
   Token da aplicação (contém tenant_id)
-------------------------------------------------------------------*/

func (m *JWTManager) Generate(userID int64, role *models.UserRole) (string, error) {
	roleStr := ""
	if role != nil {
		roleStr = string(*role)
	}
	claims := &fullClaims{
		UserID: userID,

		Role: roleStr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(m.secret)
}

/* ------------------------------------------------------------------
   Login-token (etapa 1)  –  TTL curto, sem tenant_id
-------------------------------------------------------------------*/

func (m *JWTManager) GenerateLoginToken(userID int64) (string, error) {
	claims := &fullClaims{
		UserID:   userID,
		LoginTok: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(m.secret)
}

// devolve userID contido no login-token
func (m *JWTManager) ParseLoginToken(tok string) (int64, error) {
	var claims fullClaims
	_, err := jwt.ParseWithClaims(tok, &claims,
		func(t *jwt.Token) (any, error) { return m.secret, nil })
	if err != nil {
		return 0, err
	}
	if !claims.LoginTok {
		return 0, errors.New("not a login token")
	}
	return claims.UserID, nil
}

/* ------------------------------------------------------------------*/

var Default = New(os.Getenv("JWT_SECRET"), 24*time.Hour)
