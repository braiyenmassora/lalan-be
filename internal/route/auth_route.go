// internal/route/auth_route.go
package route

import (
	"net/http"

	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
)

/*
Mendaftarkan route autentikasi.
Menyiapkan endpoint untuk registrasi, login, profil, dan test terproteksi dengan middleware JWT.
*/
func AuthRoutes(h *handler.AuthHandler) {
	v1 := "/v1"

	detailHandler := middleware.JWTMiddleware(http.HandlerFunc(h.GetProfile))

	http.HandleFunc(v1+"/auth/register", h.Register)
	http.HandleFunc(v1+"/auth/login", h.Login)

	http.Handle(v1+"/auth/detail", detailHandler)

	testHandler := middleware.JWTMiddleware(http.HandlerFunc(h.TestProtected))
	http.Handle(v1+"/test-protected", testHandler)
}
