# Makefile para API Go + React Server
# ====================================

# Configurações
APP_NAME = wsrs-server
BINARY_DIR = bin
BINARY_PATH = $(BINARY_DIR)/$(APP_NAME)
MAIN_PATH = ./cmd/wsrs/main.go
MIGRATIONS_PATH = ./internal/store/pgstore/migrations

# Variáveis de ambiente padrão para desenvolvimento
export LOG_LEVEL ?= info
export WSRS_DATABASE_HOST ?= localhost
export WSRS_DATABASE_PORT ?= 5432
export WSRS_DATABASE_USER ?= postgres
export WSRS_DATABASE_PASSWORD ?= password
export WSRS_DATABASE_NAME ?= wsrs_db

# Cores para output
GREEN = \033[32m
YELLOW = \033[33m
RED = \033[31m
BLUE = \033[34m
NC = \033[0m # No Color

.PHONY: help
help: ## Mostra este menu de ajuda
	@echo "$(BLUE)API Go + React Server - Comandos Disponíveis:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# ===========================
# 🏗️  BUILD & DEVELOPMENT
# ===========================

.PHONY: build
build: ## Compila a aplicação
	@echo "$(YELLOW)📦 Compilando aplicação...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)✅ Aplicação compilada em $(BINARY_PATH)$(NC)"

.PHONY: build-linux
build-linux: ## Compila para Linux (útil para Docker/deploy)
	@echo "$(YELLOW)📦 Compilando para Linux...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH)-linux $(MAIN_PATH)
	@echo "$(GREEN)✅ Aplicação compilada para Linux em $(BINARY_PATH)-linux$(NC)"

.PHONY: clean
clean: ## Remove binários e arquivos temporários
	@echo "$(YELLOW)🧹 Limpando arquivos...$(NC)"
	@rm -rf $(BINARY_DIR)
	@rm -rf internal/logger/logs/*.log*
	@rm -rf tmp/
	@go clean
	@echo "$(GREEN)✅ Limpeza concluída$(NC)"

.PHONY: deps
deps: ## Instala/atualiza dependências
	@echo "$(YELLOW)📥 Instalando dependências...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependências atualizadas$(NC)"

.PHONY: deps-upgrade
deps-upgrade: ## Atualiza todas as dependências para versões mais recentes
	@echo "$(YELLOW)⬆️  Atualizando dependências...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✅ Dependências atualizadas$(NC)"

# ===========================
# 🚀 RUN & DEVELOPMENT
# ===========================

.PHONY: run
run: ## Executa a aplicação
	@echo "$(YELLOW)🚀 Iniciando servidor...$(NC)"
	@go run $(MAIN_PATH)

.PHONY: run-debug
run-debug: ## Executa em modo debug
	@echo "$(YELLOW)🐛 Iniciando servidor em modo debug...$(NC)"
	@LOG_LEVEL=debug go run $(MAIN_PATH)

.PHONY: run-bin
run-bin: build ## Compila e executa o binário
	@echo "$(YELLOW)🚀 Executando binário...$(NC)"
	@$(BINARY_PATH)

.PHONY: dev
dev: ## Modo desenvolvimento com Docker Compose
	@echo "$(YELLOW)🔥 Iniciando desenvolvimento com Docker Compose...$(NC)"
	@docker-compose up --build

.PHONY: dev-stop
dev-stop: ## Para o ambiente de desenvolvimento
	@echo "$(YELLOW)🛑 Parando ambiente de desenvolvimento...$(NC)"
	@docker-compose down

.PHONY: dev-logs
dev-logs: ## Mostra logs do ambiente de desenvolvimento
	@echo "$(YELLOW)📋 Logs do ambiente de desenvolvimento...$(NC)"
	@docker-compose logs -f

# ===========================
# 🐘 DATABASE
# ===========================

.PHONY: db-up
db-up: ## Inicia containers do banco de dados
	@echo "$(YELLOW)🐘 Iniciando banco de dados...$(NC)"
	@docker compose up -d db
	@echo "$(GREEN)✅ PostgreSQL iniciado na porta $(WSRS_DATABASE_PORT)$(NC)"

.PHONY: db-up-all
db-up-all: ## Inicia banco + pgAdmin
	@echo "$(YELLOW)🐘 Iniciando banco de dados e pgAdmin...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)✅ PostgreSQL: localhost:$(WSRS_DATABASE_PORT)$(NC)"
	@echo "$(GREEN)✅ pgAdmin: http://localhost:8081 (admin@admin.com / password)$(NC)"

.PHONY: db-down
db-down: ## Para containers do banco
	@echo "$(YELLOW)⏹️  Parando banco de dados...$(NC)"
	@docker compose down
	@echo "$(GREEN)✅ Banco de dados parado$(NC)"

.PHONY: db-restart
db-restart: db-down db-up ## Reinicia o banco de dados

.PHONY: db-logs
db-logs: ## Mostra logs do banco
	@docker compose logs -f db

.PHONY: db-shell
db-shell: ## Acessa shell do PostgreSQL
	@echo "$(YELLOW)🐘 Conectando ao PostgreSQL...$(NC)"
	@docker compose exec db psql -U $(WSRS_DATABASE_USER) -d $(WSRS_DATABASE_NAME)

.PHONY: db-reset
db-reset: ## Remove volumes e reinicia banco (⚠️  APAGA TODOS OS DADOS)
	@echo "$(RED)⚠️  ATENÇÃO: Isso apagará todos os dados do banco!$(NC)"
	@read -p "Tem certeza? [y/N]: " confirm && [ "$$confirm" = "y" ]
	@docker compose down -v
	@docker compose up -d db
	@echo "$(GREEN)✅ Banco de dados resetado$(NC)"

# ===========================
# 🔄 DATABASE MIGRATIONS
# ===========================

.PHONY: migrate-up
migrate-up: ## Executa migrações para cima
	@echo "$(YELLOW)⬆️  Executando migrações...$(NC)"
	@cd $(MIGRATIONS_PATH) && tern migrate
	@echo "$(GREEN)✅ Migrações executadas$(NC)"

.PHONY: migrate-down
migrate-down: ## Reverte última migração
	@echo "$(YELLOW)⬇️  Revertendo migração...$(NC)"
	@cd $(MIGRATIONS_PATH) && tern migrate --destination -1
	@echo "$(GREEN)✅ Migração revertida$(NC)"

.PHONY: migrate-status
migrate-status: ## Mostra status das migrações
	@echo "$(YELLOW)📊 Status das migrações:$(NC)"
	@cd $(MIGRATIONS_PATH) && tern status

.PHONY: migrate-new
migrate-new: ## Cria nova migração (uso: make migrate-new NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)❌ Uso: make migrate-new NAME=nome_da_migracao$(NC)"; \
		exit 1; \
	fi
	@cd $(MIGRATIONS_PATH) && tern new $(NAME)
	@echo "$(GREEN)✅ Nova migração criada: $(NAME)$(NC)"

# ===========================
# 🏗️  CODE GENERATION
# ===========================

.PHONY: generate
generate: ## Executa go generate (migrations + sqlc)
	@echo "$(YELLOW)🔧 Executando geradores...$(NC)"
	@go generate ./...
	@echo "$(GREEN)✅ Código gerado$(NC)"

.PHONY: sqlc-generate
sqlc-generate: ## Gera código SQLC apenas
	@echo "$(YELLOW)🔧 Gerando código SQLC...$(NC)"
	@sqlc generate -f ./internal/store/pgstore/sqlc.yaml
	@echo "$(GREEN)✅ Código SQLC gerado$(NC)"

# ===========================
# 🧪 TESTING
# ===========================

.PHONY: test
test: ## Executa todos os testes
	@echo "$(YELLOW)🧪 Executando testes...$(NC)"
	@go test ./...

.PHONY: test-verbose
test-verbose: ## Executa testes com output detalhado
	@echo "$(YELLOW)🧪 Executando testes (verbose)...$(NC)"
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## Executa testes com coverage
	@echo "$(YELLOW)🧪 Executando testes com coverage...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage gerado em coverage.html$(NC)"

.PHONY: test-race
test-race: ## Executa testes com detecção de race conditions
	@echo "$(YELLOW)🧪 Executando testes (race detection)...$(NC)"
	@go test -race ./...

.PHONY: benchmark
benchmark: ## Executa benchmarks
	@echo "$(YELLOW)📊 Executando benchmarks...$(NC)"
	@go test -bench=. ./...

# ===========================
# 📝 LINTING & FORMATTING
# ===========================

.PHONY: fmt
fmt: ## Formata código Go
	@echo "$(YELLOW)✨ Formatando código...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Código formatado$(NC)"

.PHONY: lint
lint: ## Executa linting com golangci-lint
	@echo "$(YELLOW)🔍 Executando linting...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(RED)❌ golangci-lint não encontrado$(NC)"; \
		echo "$(YELLOW)💡 Instale com: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2$(NC)"; \
	fi

.PHONY: lint-fix
lint-fix: ## Executa linting e corrige problemas automaticamente
	@echo "$(YELLOW)🔧 Executando linting com correções...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix; \
	else \
		echo "$(RED)❌ golangci-lint não encontrado$(NC)"; \
	fi

.PHONY: vet
vet: ## Executa go vet
	@echo "$(YELLOW)🔍 Executando go vet...$(NC)"
	@go vet ./...

.PHONY: check
check: fmt vet lint test ## Executa todas as verificações (fmt + vet + lint + test)

# ===========================
# 📊 MONITORING & LOGS
# ===========================

.PHONY: logs
logs: ## Mostra logs da aplicação
	@echo "$(YELLOW)📋 Logs da aplicação:$(NC)"
	@if [ -f "internal/logger/logs/api.log" ]; then \
		tail -f internal/logger/logs/api.log; \
	else \
		echo "$(RED)❌ Arquivo de log não encontrado$(NC)"; \
	fi

.PHONY: logs-errors
logs-errors: ## Mostra apenas logs de erro
	@echo "$(YELLOW)🚨 Logs de erro:$(NC)"
	@if [ -f "internal/logger/logs/api.log" ]; then \
		grep -i "error\|fatal" internal/logger/logs/api.log | tail -20; \
	else \
		echo "$(RED)❌ Arquivo de log não encontrado$(NC)"; \
	fi

.PHONY: logs-clean
logs-clean: ## Limpa arquivos de log
	@echo "$(YELLOW)🧹 Limpando logs...$(NC)"
	@rm -f internal/logger/logs/*.log*
	@echo "$(GREEN)✅ Logs limpos$(NC)"

# ===========================
# 🚀 DEPLOYMENT
# ===========================

.PHONY: docker-build
docker-build: ## Constrói imagem Docker (requer Dockerfile)
	@echo "$(YELLOW)🐳 Construindo imagem Docker...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)✅ Imagem Docker criada: $(APP_NAME):latest$(NC)"

.PHONY: docker-run
docker-run: ## Executa aplicação no Docker
	@echo "$(YELLOW)🐳 Executando no Docker...$(NC)"
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

.PHONY: docker-reload
docker-reload: ## Rebuilda apenas a aplicação preservando banco de dados
	@echo "$(BLUE)🔄 ===== DOCKER RELOAD (PRESERVANDO DADOS) =====$(NC)"
	@echo "$(YELLOW)⏹️  Parando apenas o container da aplicação...$(NC)"
	@docker compose stop app 2>/dev/null || true
	@docker compose rm -f app 2>/dev/null || true
	@echo "$(GREEN)✅ Container da app removido$(NC)"
	@echo ""
	@echo "$(YELLOW)🧹 Removendo imagem antiga da aplicação...$(NC)"
	@docker rmi ask-me-anything-app:latest 2>/dev/null || true
	@echo "$(GREEN)✅ Imagem antiga removida$(NC)"
	@echo ""
	@echo "$(YELLOW)🏗️  Reconstruindo apenas a aplicação...$(NC)"
	@docker compose build --no-cache app
	@echo "$(GREEN)✅ Build da aplicação concluído$(NC)"
	@echo ""
	@echo "$(YELLOW)🚀 Iniciando aplicação...$(NC)"
	@docker compose up -d app
	@echo ""
	@echo "$(BLUE)⏳ Aguardando aplicação ficar pronta...$(NC)"
	@sleep 3
	@echo ""
	@echo "$(GREEN)🎉 ===== APLICAÇÃO ATUALIZADA =====$(NC)"
	@echo "$(GREEN)✅ API: http://localhost:8080$(NC)"
	@echo "$(BLUE)💾 Banco de dados preservado$(NC)"
	@echo ""
	@echo "$(BLUE)� Status dos containers:$(NC)"
	@docker compose ps
	@echo ""
	@echo "$(BLUE)🔗 Testando API...$(NC)"
	@sleep 2
	@curl -s -o /dev/null -w "Status: %{http_code} | Tempo: %{time_total}s\n" http://localhost:8080/api/rooms || echo "$(YELLOW)⚠️  API ainda não respondeu (aguarde alguns segundos)$(NC)"
	@echo ""
	@echo "$(BLUE)📋 Logs da aplicação (Ctrl+C para sair):$(NC)"
	@echo "$(YELLOW)💡 Use 'make docker-logs' para ver logs novamente$(NC)"
	@echo ""
	@docker compose logs -f app

.PHONY: docker-logs
docker-logs: ## Mostra logs dos containers em tempo real
	@echo "$(BLUE)📋 Logs dos serviços Docker (Ctrl+C para sair):$(NC)"
	@docker compose logs -f

.PHONY: docker-quick
docker-quick: ## Restart rápido da aplicação (sem rebuild)
	@echo "$(BLUE)⚡ ===== RESTART RÁPIDO =====$(NC)"
	@echo "$(YELLOW)🔄 Reiniciando apenas a aplicação...$(NC)"
	@docker compose restart app
	@echo "$(GREEN)✅ Aplicação reiniciada$(NC)"
	@echo "$(GREEN)✅ API: http://localhost:8080$(NC)"
	@sleep 2
	@docker compose logs --tail=10 app

.PHONY: docker-full-restart
docker-full-restart: ## Reinicia todos os serviços (preservando dados)
	@echo "$(BLUE)🔄 ===== RESTART COMPLETO (PRESERVANDO DADOS) =====$(NC)"
	@echo "$(YELLOW)⏹️  Parando todos os containers...$(NC)"
	@docker compose down
	@echo "$(YELLOW)🚀 Iniciando todos os serviços...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)✅ Todos os serviços reiniciados$(NC)"
	@echo "$(GREEN)✅ API: http://localhost:8080$(NC)"
	@echo "$(GREEN)✅ PostgreSQL: localhost:5432$(NC)"
	@echo "$(GREEN)✅ pgAdmin: http://localhost:8081$(NC)"
	@docker compose ps

.PHONY: docker-clean
docker-clean: ## Limpeza completa (remove volumes, imagens órfãs)
	@echo "$(BLUE)🧹 ===== LIMPEZA COMPLETA =====$(NC)"
	@echo "$(YELLOW)⚠️  Isso vai remover containers, volumes e imagens não utilizadas$(NC)"
	@read -p "Continuar? [y/N]: " confirm && [ "$$confirm" = "y" ] || exit 1
	@docker compose down -v --remove-orphans
	@docker system prune -f
	@docker volume prune -f
	@echo "$(GREEN)✅ Limpeza concluída$(NC)"

.PHONY: docker-fresh-start
docker-fresh-start: ## ⚠️  RESET COMPLETO - Apaga TODOS os dados e recria tudo
	@echo "$(RED)⚠️  ===== RESET COMPLETO - APAGA TODOS OS DADOS =====$(NC)"
	@echo "$(RED)⚠️  Isso vai APAGAR todos os dados do banco!$(NC)"
	@read -p "Tem certeza que quer APAGAR TODOS OS DADOS? [y/N]: " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "$(YELLOW)🗑️  Removendo tudo...$(NC)"
	@docker compose down -v --remove-orphans
	@docker rmi ask-me-anything-app:latest 2>/dev/null || true
	@echo "$(YELLOW)🏗️  Recriando tudo do zero...$(NC)"
	@docker compose up --build -d
	@echo "$(GREEN)✅ Sistema recriado do zero$(NC)"
	@echo "$(GREEN)✅ API: http://localhost:8080$(NC)"
	@echo "$(GREEN)✅ PostgreSQL: localhost:5432$(NC)"
	@echo "$(GREEN)✅ pgAdmin: http://localhost:8081$(NC)"
	@echo "$(YELLOW)💡 Execute 'make migrate-up' para criar as tabelas$(NC)"

.PHONY: docker-status
docker-status: ## Mostra status detalhado dos containers
	@echo "$(BLUE)📊 ===== STATUS DOS CONTAINERS =====$(NC)"
	@docker compose ps -a
	@echo ""
	@echo "$(BLUE)💾 Uso de recursos:$(NC)"
	@docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
	@echo ""
	@echo "$(BLUE)🔗 Endpoints disponíveis:$(NC)"
	@echo "$(GREEN)✅ API: http://localhost:8080$(NC)"
	@echo "$(GREEN)✅ API Health: http://localhost:8080/api/rooms$(NC)"
	@echo "$(GREEN)✅ PostgreSQL: localhost:5432$(NC)"
	@echo "$(GREEN)✅ pgAdmin: http://localhost:8081$(NC)"

# ===========================
# 🧪 API TESTING & MONITORING
# ===========================

.PHONY: test-api
test-api: ## Testa endpoints principais da API
	@echo "$(BLUE)🧪 ===== TESTANDO API =====$(NC)"
	@echo "$(YELLOW)📡 Testando conexão...$(NC)"
	@curl -s -o /dev/null -w "GET /api/rooms - Status: %{http_code} | Tempo: %{time_total}s\n" http://localhost:8080/api/rooms || echo "❌ API não está respondendo"
	@echo ""
	@echo "$(YELLOW)📝 Criando sala de teste...$(NC)"
	@curl -s -X POST http://localhost:8080/api/rooms \
		-H "Content-Type: application/json" \
		-d '{"theme": "Sala de Teste - $(shell date +%H:%M:%S)"}' \
		-w "POST /api/rooms - Status: %{http_code} | Tempo: %{time_total}s\n" || echo "❌ Erro ao criar sala"
	@echo ""
	@echo "$(YELLOW)📋 Listando salas...$(NC)"
	@curl -s http://localhost:8080/api/rooms | jq '.' 2>/dev/null || curl -s http://localhost:8080/api/rooms

.PHONY: monitor
monitor: ## Monitor em tempo real dos logs da API
	@echo "$(BLUE)📡 ===== MONITOR DA API =====$(NC)"
	@echo "$(YELLOW)💡 Pressione Ctrl+C para sair$(NC)"
	@echo ""
	@docker compose logs -f app | grep --line-buffered -E "(INFO|ERROR|WARN|🚀|📍|🎯)" --color=always

# ===========================
# 🔧 UTILITIES
# ===========================

.PHONY: env-copy
env-copy: ## Cria .env (se existir um template, usa; senão, cria vazio)
	@if [ -f .env ]; then \
		echo "$(YELLOW)⚠️  Arquivo .env já existe$(NC)"; \
	elif [ -f .env.example ]; then \
		cp .env.example .env; \
		echo "$(GREEN)✅ Arquivo .env criado a partir do .env.example$(NC)"; \
	else \
		touch .env; \
		echo "$(YELLOW)ℹ️  .env.example não encontrado; criado .env vazio$(NC)"; \
	fi

.PHONY: env-show
env-show: ## Mostra variáveis de ambiente
	@echo "$(BLUE)📋 Variáveis de ambiente atuais:$(NC)"
	@echo "LOG_LEVEL=$(LOG_LEVEL)"
	@echo "WSRS_DATABASE_HOST=$(WSRS_DATABASE_HOST)"
	@echo "WSRS_DATABASE_PORT=$(WSRS_DATABASE_PORT)"
	@echo "WSRS_DATABASE_USER=$(WSRS_DATABASE_USER)"
	@echo "WSRS_DATABASE_NAME=$(WSRS_DATABASE_NAME)"

.PHONY: install-tools
install-tools: ## Instala ferramentas de desenvolvimento
	@echo "$(YELLOW)🔧 Instalando ferramentas de desenvolvimento...$(NC)"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/jackc/tern/v2@latest
	@echo "$(GREEN)✅ Ferramentas instaladas:$(NC)"
	@echo "  - sqlc (geração de código SQL)"
	@echo "  - tern (migrações)"

.PHONY: size
size: build ## Mostra tamanho do binário
	@echo "$(BLUE)📏 Tamanho do binário:$(NC)"
	@ls -lh $(BINARY_PATH) | awk '{print $$5 " " $$9}'

.PHONY: info
info: ## Mostra informações do projeto
	@echo "$(BLUE)ℹ️  Informações do Projeto:$(NC)"
	@echo "Nome: $(APP_NAME)"
	@echo "Go Version: $$(go version)"
	@echo "Binário: $(BINARY_PATH)"
	@echo "Main: $(MAIN_PATH)"
	@echo "Migrações: $(MIGRATIONS_PATH)"
	@echo ""
	@make env-show

# ===========================
# 🎯 SHORTCUTS ÚTEIS
# ===========================

.PHONY: setup
setup: deps env-copy db-up migrate-up install-tools ## Setup completo do projeto
	@echo "$(GREEN)🎉 Setup completo! Execute 'make run' para iniciar$(NC)"

.PHONY: start
start: docker-reload test-api ## Rebuilda a aplicação e testa
	@echo "$(GREEN)🎉 ===== APLICAÇÃO PRONTA! =====$(NC)"
	@echo "$(BLUE)📝 Comandos úteis para desenvolvimento:$(NC)"
	@echo "  $(GREEN)make docker-reload$(NC)      - Rebuilda apenas a app (preserva dados)"
	@echo "  $(GREEN)make docker-quick$(NC)       - Restart rápido (sem rebuild)"
	@echo "  $(GREEN)make docker-full-restart$(NC) - Restart completo (preserva dados)"
	@echo "  $(GREEN)make monitor$(NC)            - Monitor dos logs"
	@echo "  $(GREEN)make test-api$(NC)           - Testar API"
	@echo "  $(GREEN)make docker-status$(NC)      - Status dos containers"

.PHONY: restart
restart: db-restart run ## Reinicia banco e aplicação

.PHONY: reset
reset: clean db-reset setup ## Reset completo do projeto (⚠️  APAGA DADOS)

# Comandos padrão
.DEFAULT_GOAL := help
