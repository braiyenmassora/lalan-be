package public

import (
	"github.com/gorilla/mux"
)

/*
SetupPublicRoutes mendaftarkan endpoint publik (tanpa autentikasi).

Route:
- GET /api/v1/public/category -> GetAllCategories
- GET /api/v1/public/item     -> GetAllItems
- GET /api/v1/public/tnc      -> GetAllTermsAndConditions
*/
func SetupPublicRoutes(router *mux.Router, h *PublicHandler) {
	public := router.PathPrefix("/api/v1/public").Subrouter()

	public.HandleFunc("/category", h.GetAllCategories).Methods("GET")
	public.HandleFunc("/item", h.GetAllItems).Methods("GET")
	public.HandleFunc("/tnc", h.GetAllTermsAndConditions).Methods("GET")
}
