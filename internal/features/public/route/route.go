package route

import (
	"lalan-be/internal/features/public/handler"
	"net/http"
)

func SetupPublicRoutes(mux *http.ServeMux, h *handler.PublicHandler) {
	mux.HandleFunc("/v1/public/categories", h.GetCategories)
}
