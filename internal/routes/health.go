package routes

import (
	"encoding/json"
	"net/http"
	"time"
)

// healthCheck returns the health status of the application
func (rt *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "ask-me-anything",
		"version":   "1.0.0",
	}

	// Test database connection
	if err := rt.db.Ping(); err != nil {
		health["status"] = "error"
		health["database"] = "disconnected"
		health["error"] = err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		health["database"] = "connected"
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// pingCheck returns a simple pong response
func (rt *Router) pingCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message":   "pong",
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
