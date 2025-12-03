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
  - PUT  /booking/{id}/status → update status booking

Output:
- Router terkonfigurasi dengan endpoint hoster yang aman dan siap digunakan
*/
func SetupBookingRoutes(router *mux.Router, h *HosterBookingHandler) {
	protected := router.PathPrefix("/api/v1/hoster").Subrouter()

	// JWT + Role check
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	// Route normal
	protected.HandleFunc("/booking", h.GetListBooking).Methods("GET", "OPTIONS")
	protected.HandleFunc("/booking/{id}", h.GetDetailBooking).Methods("GET", "OPTIONS")
	protected.HandleFunc("/booking/status/{id}", h.UpdateBookingStatus).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/customer", h.GetCustomerList).Methods("GET", "OPTIONS")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
