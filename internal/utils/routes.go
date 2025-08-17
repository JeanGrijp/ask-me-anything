// Package utils provides utility functions for the application
package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// LogRoutes percorre todas as rotas registradas e as exibe no terminal
func LogRoutes(r chi.Router) {
	fmt.Println("\n🚀 ===== API ROUTES =====")
	fmt.Println("📍 Available endpoints:")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Limpar a rota para exibição mais limpa
		cleanRoute := strings.ReplaceAll(route, "/*", "")

		// Adicionar cores para diferentes métodos HTTP
		var methodColor string
		switch method {
		case "GET":
			methodColor = "\033[32m" // Verde
		case "POST":
			methodColor = "\033[34m" // Azul
		case "PATCH":
			methodColor = "\033[33m" // Amarelo
		case "DELETE":
			methodColor = "\033[31m" // Vermelho
		case "PUT":
			methodColor = "\033[35m" // Magenta
		default:
			methodColor = "\033[37m" // Branco
		}

		fmt.Printf("  %s%-6s\033[0m %s\n", methodColor, method, cleanRoute)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("❌ Error walking routes: %v\n", err)
		return
	}

	fmt.Println("🎯 WebSocket endpoint:")
	fmt.Printf("  \033[36m%-6s\033[0m %s\n", "WS", "/subscribe/{room_id}")
	fmt.Println("========================")
}
