package route

import (
	"lalan-be/internal/handler"
	"net/http"
)

// Mengatur rute kategori
func CategoryRoute(h *handler.CategoryHandler) {
	v1 := "/v1"
	// Rute untuk menambah kategori
	http.HandleFunc(v1+"/category/add", h.AddCategory)
}
