package admin

import (
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"
	"net/http"
)

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := middleware.GetUserRole(r)
		if userRole != "admin" {
			response.Forbidden(w, "Access denied: admin role required")
			return
		}
		next(w, r)
	}
}

func SetupAdminRoutes(mux *http.ServeMux, h *AdminHandler) {
	// Public admin routes - no authentication required
	mux.HandleFunc("/v1/admin/register", h.CreateAdmin)
	mux.HandleFunc("/v1/admin/login", h.LoginAdmin)

	// Category management - require JWT + admin role
	mux.HandleFunc("/v1/admin/categories/create", middleware.JWTMiddleware(http.HandlerFunc(AdminAuthMiddleware(h.CreateCategory))).ServeHTTP)
	mux.HandleFunc("/v1/admin/categories/update/", middleware.JWTMiddleware(http.HandlerFunc(AdminAuthMiddleware(h.UpdateCategory))).ServeHTTP)
	mux.HandleFunc("/v1/admin/categories/delete/", middleware.JWTMiddleware(http.HandlerFunc(AdminAuthMiddleware(h.DeleteCategory))).ServeHTTP)
}
