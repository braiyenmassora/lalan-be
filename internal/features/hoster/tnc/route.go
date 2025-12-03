package tnc

import (
	"net/http"

	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
SetupTnCRoutes mendaftarkan endpoint T&C untuk hoster.

Alur kerja:
1. Buat subrouter dengan prefix /api/v1/hoster
2. Terapkan middleware JWT → Hoster (protected route)
3. Daftarkan endpoint:
  - GET  /tnc      → tampilkan T&C milik hoster
  - POST /tnc      → buat T&C baru oleh hoster
  - PUT  /tnc/{id} → update T&C milik hoster berdasarkan ID

Output:
- Router terkonfigurasi dengan endpoint hoster yang aman dan siap digunakan
*/
func SetupTnCRoutes(router *mux.Router, h *HosterTnCHandler) {
	protected := router.PathPrefix("/api/v1/hoster").Subrouter()

	// JWT + Role check
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	// Route normal (use singular "tnc")
	protected.HandleFunc("/tnc", h.GetTnC).Methods("GET", "OPTIONS")
	protected.HandleFunc("/tnc", h.CreateTnC).Methods("POST", "OPTIONS")
	protected.HandleFunc("/tnc/{id}", h.UpdateTnC).Methods("PUT", "OPTIONS")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
