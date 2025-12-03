package public

import (
	"github.com/gorilla/mux"
)

/*
SetupPublicRoutes mendaftarkan endpoint publik (tanpa autentikasi).

Route:
- GET /api/v1/public/item        -> GetAllItems (list item untuk halaman home)
- GET /api/v1/public/item/{id}   -> GetItemDetail (detail item dengan JOIN: category + hoster + tnc)
*/
func SetupPublicRoutes(router *mux.Router, h *PublicHandler) {
	public := router.PathPrefix("/api/v1/public").Subrouter()

	public.HandleFunc("/item", h.GetAllItems).Methods("GET")
	public.HandleFunc("/item/{id}", h.GetItemDetail).Methods("GET")
}
