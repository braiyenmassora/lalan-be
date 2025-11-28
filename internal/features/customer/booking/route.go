package booking

import (
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
SetupBookingRoutes mendaftarkan semua endpoint fitur booking ke router utama.

Alur kerja:
1. Membuat subrouter dengan prefix "/api/v1/customer"
2. Menerapkan middleware secara berurutan:
  - JWTMiddleware   → validasi dan ekstrak token JWT
  - Customer        → pastikan role user adalah "customer"

3. Register route dengan urutan spesifik-ke-umum agar tidak tertimpa:
  - GET  /booking/me          → daftar booking user login
  - GET  /booking/{id}        → detail satu booking
  - POST /booking             → buat booking baru

Output:
- Subrouter yang sudah terproteksi dan siap menerima request booking.
Semua endpoint otomatis ter-autentikasi dan ter-authorize (hanya customer yang bisa akses).
*/
func SetupBookingRoutes(router *mux.Router, h *BookingHandler) {
	// Subrouter khusus customer area
	protected := router.PathPrefix("/api/v1/customer").Subrouter()

	// Middleware stack: JWT dulu → baru role check
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Customer)

	// Urutan penting: route dengan path parameter harus didefinisikan sebelum route umum
	protected.HandleFunc("/booking", h.GetListBookings).Methods("GET")
	protected.HandleFunc("/booking/{id}", h.GetDetailBooking).Methods("GET")
	protected.HandleFunc("/booking", h.CreateBooking).Methods("POST")
}
