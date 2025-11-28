package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*
SetupAuthRoutes mendaftarkan endpoint autentikasi ke router utama.

Fungsi ini membuat subrouter dengan prefix "/api/v1/auth" dan
menghubungkan path URL dengan handler function yang sesuai.

Daftar Endpoint:
- POST /api/v1/auth/login      : Login user (semua role)
- POST /api/v1/auth/register   : Registrasi user baru
- POST /api/v1/auth/verify-otp : Verifikasi email customer
- POST /api/v1/auth/resend-otp : Kirim ulang OTP

Output:
- Router yang sudah dikonfigurasi dengan route auth.
*/
func SetupAuthRoutes(router *mux.Router, h *AuthHandler) {
	auth := router.PathPrefix("/api/v1/auth").Subrouter()

	// PUBLIC ROUTES — tambahkan OPTIONS di semua endpoint yang dipanggil dari browser
	auth.HandleFunc("/login", h.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/register", h.Register).Methods("POST", "OPTIONS")
	auth.HandleFunc("/verify-otp", h.VerifyEmail).Methods("POST", "OPTIONS")
	auth.HandleFunc("/resend-otp", h.ResendOTP).Methods("POST", "OPTIONS")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	auth.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
