package route

import (
	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
	"net/http"
)

// CategoryRoutes mengatur semua rute untuk kategori dengan JWT protection
func CategoryRoutes(h *handler.CategoryHandler) {
	v1 := "/v1/category"

	// Endpoint yang membutuhkan JWT (hanya admin/user yang login)
	addHandler := middleware.JWTMiddleware(http.HandlerFunc(h.AddCategory))
	updateHandler := middleware.JWTMiddleware(http.HandlerFunc(h.UpdateCategory))
	deleteHandler := middleware.JWTMiddleware(http.HandlerFunc(h.DeleteCategory))

	// Register protected routes
	http.Handle(v1+"/add", addHandler)
	http.Handle(v1+"/update", updateHandler)
	http.Handle(v1+"/delete", deleteHandler)

	// Public endpoints (tidak perlu login untuk melihat kategori)
	http.HandleFunc(v1+"/list", h.GetAllCategories)
	http.HandleFunc(v1+"/detail", h.GetCategoryByID)
}
