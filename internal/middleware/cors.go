package middleware

import (
	"net/http"

	"lalan-be/internal/config"
)

/*
CORSMiddleware
menangani header CORS dan meneruskan ke handler berikutnya
*/
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Setup CORS berdasarkan env
		env := config.GetEnv("APP_ENV", "dev")
		origin := "*"
		// Cek environment prod
		if env == "prod" {
			origin = "https://yourdomain.com"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
