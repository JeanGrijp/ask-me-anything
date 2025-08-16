package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/api"
	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/JeanGrijp/ask-me-anything/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		logger.Default.Fatal(ctx, "failed to load environment variables", "error", err)
	}

	logger.Default.Info(ctx, "starting application")

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WSRS_DATABASE_USER"),
		os.Getenv("WSRS_DATABASE_PASSWORD"),
		os.Getenv("WSRS_DATABASE_HOST"),
		os.Getenv("WSRS_DATABASE_PORT"),
		os.Getenv("WSRS_DATABASE_NAME"),
	))
	if err != nil {
		logger.Default.Fatal(ctx, "failed to create database connection pool", "error", err)
	}
	logger.Default.Info(ctx, "database connection pool created")

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Default.Fatal(ctx, "failed to ping database", "error", err)
	}

	logger.Default.Info(ctx, "database connection established")

	handler := api.NewHandler(pgstore.New(pool))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	logger.Default.Info(ctx, "starting HTTP server", "port", 8080)
	logger.Default.Info(ctx, "visit http://localhost:8080 to access the API")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Default.Fatal(ctx, "HTTP server error", "error", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Default.Info(ctx, "shutting down application")

	// Graceful shutdown com timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Default.Fatal(ctx, "server shutdown error", "error", err)
	}
}
