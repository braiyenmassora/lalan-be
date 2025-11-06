package route

import (
	"net/http"

	"lalan-be/internal/handler"
)

// auth routes
func AuthRoutes(h *handler.AuthHandler) {
	v1 := "/v1"
	http.HandleFunc(v1+"/auth/register", h.Register)
	http.HandleFunc(v1+"/auth/login", h.Login)
}
