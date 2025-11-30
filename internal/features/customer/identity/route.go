package identity

import (
	"lalan-be/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

/*
SetupIdentityRoutes mendaftarkan semua endpoint fitur verifikasi identitas (KTP) khusus customer.

Alur kerja:
1. Membuat subrouter dengan prefix "/api/v1/customer/identity"
2. Menerapkan middleware berurutan:
  - JWTMiddleware → validasi & ekstrak token JWT
  - Customer     → pastikan role user adalah "customer"

3. Register endpoint:
  - POST  /upload → upload KTP pertama kali
  - PUT   /update → re-upload KTP (reset status jadi pending)
  - GET   /status → ambil status verifikasi KTP user login

Output:
  - Subrouter terproteksi penuh — semua route otomatis ter-autentikasi dan ter-authorize.
    Hanya customer yang sedang login yang bisa mengakses endpoint ini.
*/
func SetupIdentityRoutes(router *mux.Router, h *IdentityHandler) {
	// Subrouter khusus area customer identity
	identity := router.PathPrefix("/api/v1/customer").Subrouter()

	// Middleware stack: JWT dulu → baru role check
	identity.Use(middleware.JWTMiddleware)
	identity.Use(middleware.Customer)

	// Route registration
	identity.HandleFunc("/identity", h.UploadKTP).Methods("POST")
	identity.HandleFunc("/identity", h.UpdateKTP).Methods("PUT")
	identity.HandleFunc("/identity", h.GetStatusKTP).Methods("GET")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	identity.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
