package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/JeanGrijp/ask-me-anything/internal/database"
	"github.com/JeanGrijp/ask-me-anything/internal/services/auth"
	"github.com/JeanGrijp/ask-me-anything/internal/services/email"
)

type Handler struct {
	db           *database.DB
	logger       *slog.Logger
	magicService *auth.MagicLinkService
	emailService *email.EmailService
}

type LoginRequest struct {
	Email string `json:"email"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

func NewAuthHandler(db *database.DB) *Handler {
	return &Handler{
		db:           db,
		logger:       slog.Default(),
		magicService: &auth.MagicLinkService{},
		emailService: email.NewEmailService(),
	}
}

func (h *Handler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Generate magic link
	token, expiresAt, err := h.magicService.CreateMagicLink(req.Email)
	if err != nil {
		h.logger.Error("Failed to create magic link", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO: Save to database using SQLC
	// For now, just send the email
	if err := h.emailService.SendMagicLink(req.Email, token); err != nil {
		h.logger.Error("Failed to send magic link email", "error", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Magic link requested", "email", req.Email, "expires_at", expiresAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Magic link sent to your email",
	})
}

func (h *Handler) VerifyMagicLink(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// TODO: Verify token in database using SQLC
	// For now, just mock success
	h.logger.Info("Magic link verified", "token", req.Token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully authenticated",
		"user": map[string]interface{}{
			"id":    1,
			"email": "user@example.com",
		},
	})
}

func (h *Handler) VerifyMagicLinkGet(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// TODO: Verify token and redirect to frontend
	h.logger.Info("Magic link clicked", "token", token)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<html>
			<body>
				<h1>✅ Login Successful!</h1>
				<p>You have been successfully authenticated.</p>
				<script>
					// In production, redirect to your frontend
					// window.location.href = 'https://yourapp.com/dashboard';
				</script>
			</body>
		</html>
	`))
}

// Logout handles POST /auth/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Successfully logged out",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser handles GET /auth/me
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user from JWT token
	response := map[string]interface{}{
		"message": "Current user info",
		"status":  "success",
		"user": map[string]interface{}{
			"id":       1,
			"email":    "user@example.com",
			"username": "testuser",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
