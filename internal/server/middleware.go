package server

import (
	"net/http"
)

// withMiddleware wraps a handler with common middleware
func withMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers for development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set JSON content type for API responses
		w.Header().Set("Content-Type", "application/json")

		handler(w, r)
	}
}
