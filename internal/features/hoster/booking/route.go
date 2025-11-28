package booking

import (
	"lalan-be/internal/middleware"
	"net/http"

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

	// JWT + Role check
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	// Route normal
	protected.HandleFunc("/booking", h.GetListBookings).Methods("GET")
	protected.HandleFunc("/booking/{id}", h.GetDetailBooking).Methods("GET")
	protected.HandleFunc("/customer", h.GetCustomerList).Methods("GET")

	// Tambahan ini yang paling penting: tangkap SEMUA OPTIONS di dalam subrouter ini
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS middleware sudah jalan duluan (di main router), jadi header sudah ada
		// Kita cukup balas 204 supaya tidak 404
		w.WriteHeader(http.StatusNoContent)
	})
}
