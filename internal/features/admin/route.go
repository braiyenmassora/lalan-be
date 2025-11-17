package admin

import (
	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
SetupAdminRoutes mengatur rute untuk fitur admin.
Menggunakan mux.Router dengan middleware untuk endpoint publik dan terproteksi.
*/
func SetupAdminRoutes(router *mux.Router, h *AdminHandler) {
	admin := router.PathPrefix("/api/v1/admin").Subrouter()

	admin.HandleFunc("/register", h.CreateAdmin).Methods("POST")
	admin.HandleFunc("/login", h.LoginAdmin).Methods("POST")

	// Protected
	protected := router.PathPrefix("/api/v1/admin").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Admin)

	protected.HandleFunc("/category/create", h.CreateCategory).Methods("POST")
	protected.HandleFunc("/category/update", h.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/category/delete", h.DeleteCategory).Methods("DELETE")
	protected.HandleFunc("/category", h.GetAllCategory).Methods("GET")
}
