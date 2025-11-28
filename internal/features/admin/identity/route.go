package identity

import (
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
SetupAdminIdentityRoutes mendaftarkan semua endpoint admin untuk verifikasi identitas (KTP).

Alur kerja:
1. Buat subrouter dengan prefix /api/v1/admin/identities
2. Terapkan middleware JWT + role Admin (protected route)
3. Daftarkan endpoint:
  - GET    /pending           → daftar identitas yang menunggu verifikasi
  - GET    /{userID}          → detail identitas user tertentu
  - POST   /{id}/validate     → approve / reject KTP (targerkan record identity by id)

Output:
- Router terkonfigurasi dengan 3 endpoint admin yang aman dan siap digunakan
*/
func SetupAdminIdentityRoutes(router *mux.Router, handler *AdminIdentityHandler) {
	// Gunakan singular noun "identity" sesuai best practice
	protected := router.PathPrefix("/api/v1/admin/identity").Subrouter()

	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Admin)

	// Resource khusus untuk pending identities
	protected.HandleFunc("/pending", handler.GetPendingIdentities).Methods("GET")

	// Explicit path untuk mengambil identity berdasarkan User ID
	protected.HandleFunc("/user/{userID}", handler.GetIdentity).Methods("GET")

	// Validate by specific identity record ID
	protected.HandleFunc("/{id}/validate", handler.ValidateIdentity).Methods("POST")
}
