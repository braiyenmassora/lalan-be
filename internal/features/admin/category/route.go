package category

import (
	"net/http"

	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
SetupCategoryRoutes mendaftarkan endpoint category admin ke router.

Daftar Endpoint:
- POST   /api/v1/admin/category      : Create category (admin only)
- GET    /api/v1/admin/category      : Get all categories (admin only)
- GET    /api/v1/admin/category/{id} : Get category by ID (admin only)
- PUT    /api/v1/admin/category/{id} : Update category (admin only)
- DELETE /api/v1/admin/category/{id} : Delete category (admin only)
*/
func SetupCategoryRoutes(router *mux.Router, h *CategoryHandler) {
	protected := router.PathPrefix("/api/v1/admin/category").Subrouter()

	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Admin)

	protected.HandleFunc("", h.GetAllCategory).Methods("GET")
	protected.HandleFunc("", h.CreateCategory).Methods("POST")
	protected.HandleFunc("/{id}", h.GetCategoryByID).Methods("GET")
	protected.HandleFunc("/{id}", h.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/{id}", h.DeleteCategory).Methods("DELETE")

	// Opsional: handler khusus OPTIONS biar return 204 (lebih bersih)
	protected.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
