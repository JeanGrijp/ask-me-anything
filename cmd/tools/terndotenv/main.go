package main

import (
	"context"
	"os/exec"

	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		logger.Default.Fatal(ctx, "error loading .env file", err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)
	if err := cmd.Run(); err != nil {
		logger.Default.Fatal(ctx, "error running tern migrate", err)
	}
}
