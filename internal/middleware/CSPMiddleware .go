// internal\middleware\CSPMiddleware .go
package middleware

import (
	"net/http"
	"os"
)

func CSPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' https://apis.google.com; style-src 'self'; img-src 'self' data:; connect-src 'self';")

		// 2. Protege contra clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// 3. Protege contra MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// 4. Obriga HTTPS (funciona apenas se tiver HTTPS rodando!)
		// Só adiciona Strict-Transport-Security se estiver em produção
		if os.Getenv("APP_ENV") == "production" {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		// 5. Protege o Referer para não vazar dados sensíveis
		w.Header().Set("Referrer-Policy", "no-referrer")

		// 6. Controla APIs do navegador (opcional mas muito bom)
		w.Header().Set("Permissions-Policy",
			"geolocation=(), microphone=(), camera=(), fullscreen=(self)")

		next.ServeHTTP(w, r)
	})
}
