package identity

import (
	"lalan-be/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

/*
SetupAdminIdentityRoutes mendaftarkan semua endpoint admin untuk verifikasi identitas (KTP).

Alur kerja:
1. Buat subrouter dengan prefix /api/v1/admin/identity
2. Terapkan middleware JWT + role Admin (protected route)
3. Daftarkan endpoint:
  - GET    /pending           → daftar identitas yang menunggu verifikasi
  - GET    /{id}              → detail identitas berdasarkan KTP ID
  - POST   /validate/{id}     → approve / reject KTP berdasarkan ID

Output:
- Router terkonfigurasi dengan endpoint admin yang aman dan konsisten
*/
func SetupAdminIdentityRoutes(router *mux.Router, handler *AdminIdentityHandler) {
	// Gunakan singular "identity" untuk resource endpoint
	protected := router.PathPrefix("/api/v1/admin/identity").Subrouter()

	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Admin)

	protected.HandleFunc("/pending", handler.GetPendingIdentities).Methods("GET")
	protected.HandleFunc("/{id}", handler.GetIdentity).Methods("GET")
	protected.HandleFunc("/validate/{id}", handler.ValidateIdentity).Methods("POST")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
