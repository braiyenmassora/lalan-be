package route

import (
	"net/http"

	"lalan-be/internal/handler"
)

// Mengatur rute autentikasi
func AuthRoutes(h *handler.AuthHandler) {
	v1 := "/v1"
	// Rute untuk registrasi
	http.HandleFunc(v1+"/auth/register", h.Register)
	// Rute untuk login
	http.HandleFunc(v1+"/auth/login", h.Login)
}
