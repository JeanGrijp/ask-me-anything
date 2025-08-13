// Package http provides HTTP handlers for the ask-me-anything application
package http

import (
	"log/slog"

	"github.com/go-playground/validator/v10"

	"github.com/JeanGrijp/ask-me-anything/internal/database"
)

// Handler holds the database connection and logger for HTTP handlers
type Handler struct {
	db        *database.DB
	logger    *slog.Logger
	validator *validator.Validate
}

// New creates a new HTTP handler instance
func New(db *database.DB) *Handler {
	return &Handler{
		db:        db,
		logger:    slog.Default(),
		validator: validator.New(),
	}
}
