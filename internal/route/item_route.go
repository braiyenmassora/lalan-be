package route

import (
	"lalan-be/internal/handler"
	"net/http"
)

// Mengatur rute kategori
func ItemRoute(h *handler.ItemHandler) {
	v1 := "/v1"
	// Rute untuk menambah kategori
	http.HandleFunc(v1+"/item/add", h.AddItem)
}
