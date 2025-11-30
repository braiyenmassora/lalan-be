package identity

import (
	"lalan-be/internal/middleware"
	"net/http"

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

	protected.HandleFunc("/pending", handler.GetPendingIdentities).Methods("GET")
	protected.HandleFunc("/user/{userID}", handler.GetIdentity).Methods("GET")
	protected.HandleFunc("/{id}/validate", handler.ValidateIdentity).Methods("POST")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
