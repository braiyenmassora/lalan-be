// internal/route/auth_route.go
package route

import (
	"net/http"

	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
)

/*
Mendaftarkan route untuk autentikasi.
*/
func AuthRoutes(h *handler.AuthHandler) {
	v1 := "/v1"
	http.HandleFunc(v1+"/auth/register", h.Register)
	http.HandleFunc(v1+"/auth/login", h.Login)

	testHandler := middleware.JWTMiddleware(http.HandlerFunc(h.TestProtected))
	http.Handle(v1+"/test-protected", testHandler)
}
