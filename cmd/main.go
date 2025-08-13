package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JeanGrijp/ask-me-anything/internal/config"
	"github.com/JeanGrijp/ask-me-anything/internal/database"
	"github.com/JeanGrijp/ask-me-anything/internal/logger"
	"github.com/JeanGrijp/ask-me-anything/internal/routes"
	"github.com/JeanGrijp/ask-me-anything/internal/validators"
)

func main() {
	// Contexto base para logging
	ctx := context.Background()

	// Usar o logger padrão do projeto
	logger.Default.Info(ctx, "🚀 Iniciando Ask Me Anything API...")

	// Carregar configurações
	logger.Default.Info(ctx, "📋 Carregando configurações...")
	cfg := config.Load()
	logger.Default.Info(ctx, "✅ Configurações carregadas com sucesso")

	// Inicializar validators
	logger.Default.Info(ctx, "🔍 Inicializando validators...")
	if err := validators.InitValidator(); err != nil {
		logger.Default.Error(ctx, "❌ Erro ao inicializar validators", "error", err)
		os.Exit(1)
	}
	logger.Default.Info(ctx, "✅ Validators inicializados")

	// Conectar ao banco de dados
	logger.Default.Info(ctx, "🗄️ Conectando ao banco de dados...", "host", cfg.Database.Host, "port", cfg.Database.Port, "database", cfg.Database.DBName)
	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Default.Error(ctx, "❌ Erro ao conectar com o banco de dados", "error", err)
		os.Exit(1)
	}
	defer func() {
		logger.Default.Info(ctx, "🔒 Fechando conexão com o banco de dados...")
		if err := db.Close(); err != nil {
			logger.Default.Error(ctx, "❌ Erro ao fechar conexão com o banco", "error", err)
		}
	}()

	// Testar conexão com o banco
	logger.Default.Info(ctx, "🏥 Testando conexão com o banco de dados...")
	if err := db.Ping(); err != nil {
		logger.Default.Error(ctx, "❌ Erro ao fazer ping no banco de dados", "error", err)
		os.Exit(1)
	}
	logger.Default.Info(ctx, "✅ Conexão com banco de dados estabelecida com sucesso")

	// Configurar roteador
	logger.Default.Info(ctx, "🛣️ Configurando rotas e middlewares...")
	router := routes.NewRouter(db)
	handler := router.Setup()
	logger.Default.Info(ctx, "✅ Rotas configuradas com sucesso")

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: handler,
	}

	// Canal para capturar sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		logger.Default.Info(ctx, "🌐 Servidor HTTP iniciado",
			"port", cfg.Server.Port,
			"address", fmt.Sprintf("http://localhost:%s", cfg.Server.Port),
			"health_check", fmt.Sprintf("http://localhost:%s/health", cfg.Server.Port))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Default.Error(ctx, "❌ Erro ao iniciar servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	logger.Default.Info(ctx, "🎯 Aplicação iniciada com sucesso! Pressione Ctrl+C para parar")

	// Aguardar sinal de parada
	<-quit
	logger.Default.Info(ctx, "🛑 Sinal de parada recebido, iniciando graceful shutdown...")

	// Contexto com timeout para graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parar servidor gracefully
	logger.Default.Info(ctx, "⏳ Parando servidor HTTP...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Default.Error(ctx, "❌ Erro durante graceful shutdown", "error", err)
		os.Exit(1)
	}

	logger.Default.Info(ctx, "👋 Aplicação finalizada com sucesso!")
}
