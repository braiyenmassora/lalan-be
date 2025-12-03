package profile

import (
	"net/http"

	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
SetupProfileRoutes mendaftarkan endpoint profile untuk hoster.

Alur kerja:
1. Buat subrouter dengan prefix /api/v1/hoster
2. Terapkan middleware JWT → Hoster (protected route)
3. Daftarkan endpoint:
  - GET /profile → tampilkan profil hoster
  - PUT /profile → update profil hoster (address, phone_number, description, website, instagram, tiktok)

Output:
- Router terkonfigurasi dengan endpoint hoster yang aman dan siap digunakan
*/
func SetupProfileRoutes(router *mux.Router, h *HosterProfileHandler) {
	protected := router.PathPrefix("/api/v1/hoster").Subrouter()

	// JWT + Role check
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	// Route normal
	protected.HandleFunc("/profile", h.GetProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/profile", h.UpdateProfile).Methods("PUT", "OPTIONS")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
