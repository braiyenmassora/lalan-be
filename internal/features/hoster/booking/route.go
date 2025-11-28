package booking

import (
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
SetupBookingRoutes mendaftarkan semua endpoint booking untuk hoster.

Alur kerja:
1. Buat subrouter dengan prefix /api/v1/hoster
2. Terapkan middleware JWT → Hoster (protected route)
3. Daftarkan endpoint:
  - GET  /booking          → daftar semua booking milik hoster
  - GET  /booking/{id}     → detail satu booking

Output:
- Router terkonfigurasi dengan endpoint hoster yang aman dan siap digunakan
*/
func SetupBookingRoutes(router *mux.Router, h *BookingHandler) {
	protected := router.PathPrefix("/api/v1/hoster").Subrouter()

	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	protected.HandleFunc("/booking", h.GetListBookings).Methods("GET")
	protected.HandleFunc("/booking/{id}", h.GetDetailBooking).Methods("GET")
	protected.HandleFunc("/customer", h.GetCustomerList).Methods("GET")
}
