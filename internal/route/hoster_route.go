package route

import (
	"net/http"

	"lalan-be/internal/handler"
)

// RegisterHosterRoutes mendaftarkan semua route terkait hoster
func RegisterHosterRoutes(h *handler.HosterHandler) {
	http.HandleFunc("/hoster/register", h.RegisterHosterHandler)
}
