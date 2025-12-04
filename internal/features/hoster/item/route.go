package item

import (
	"lalan-be/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

/*
SetupItemRoutes mendaftarkan semua endpoint item untuk hoster.

Alur kerja:
1. Buat subrouter dengan prefix /api/v1/hoster
2. Terapkan middleware JWT → Hoster (protected route)
3. Daftarkan endpoint:
  - GET  /item         → daftar semua item milik hoster
  - POST /item         → buat item baru oleh hoster
  - PUT  /item/{id}    → update item milik hoster berdasarkan ID
  - DELETE /item/{id}  → hapus item milik hoster berdasarkan ID

Output:
- Router terkonfigurasi dengan endpoint hoster yang aman dan siap digunakan
*/
func SetupItemRoutes(router *mux.Router, h *HosterItemHandler) {
	protected := router.PathPrefix("/api/v1/hoster").Subrouter()

	// JWT + Role check
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	// Route normal (use singular "item")
	protected.HandleFunc("/item", h.GetListItem).Methods("GET", "OPTIONS")
	protected.HandleFunc("/item", h.CreateItem).Methods("POST", "OPTIONS")
	protected.HandleFunc("/item/category", h.GetCategory).Methods("GET", "OPTIONS") // Dropdown categories
	protected.HandleFunc("/item/{id}", h.GetItemDetail).Methods("GET", "OPTIONS")
	protected.HandleFunc("/item/{id}", h.UpdateItem).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/item/{id}", h.DeleteItem).Methods("DELETE", "OPTIONS")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
