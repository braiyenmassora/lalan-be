package auth

import (
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

	// Centralized login for admin / hoster / customer
	auth.HandleFunc("/login", h.Login).Methods("POST")
	auth.HandleFunc("/register", h.Register).Methods("POST")
	auth.HandleFunc("/verify-otp", h.VerifyEmail).Methods("POST")
	auth.HandleFunc("/resend-otp", h.ResendOTP).Methods("POST")
}
