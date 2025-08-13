.PHONY: build run test clean docker-up docker-down migrate-up migrate-down sqlc-generate

# Variables
APP_NAME=ask-me-anything
DOCKER_COMPOSE=docker-compose
DB_URL=postgres://ama_user:ama_password@localhost:5432/ask_me_anything?sslmode=disable

# Build the application
build:
	go build -o bin/$(APP_NAME) cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-reload:
	@echo "Stopping Docker containers..."
	docker-compose down
	# @echo "Restoring swagger docs..."
	# rm -rf docs/
	# swag init --parseDependency --parseInternal -g cmd/api/main.go
	@echo "Starting Docker containers..."
	COMPOSE_BAKE=true docker-compose up -d --build
	@echo "Watching logs..."
	docker compose logs -f redis n8n postgres app
	@echo "Watching for file changes..."

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# SQLC
sqlc-generate:
	sqlc generate

# Development
dev: docker-up migrate-up sqlc-generate run

# Install dependencies
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Go mod
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download
