package public

import (
	"github.com/gorilla/mux"
)

/*
SetupPublicRoutes
mengatur rute publik untuk kategori, item, dan terms tanpa autentikasi
*/
func SetupPublicRoutes(router *mux.Router, h *PublicHandler) {
	public := router.PathPrefix("/api/v1/public").Subrouter()
	public.HandleFunc("/category", h.GetAllCategories).Methods("GET")
	public.HandleFunc("/item", h.GetAllItems).Methods("GET")
	public.HandleFunc("/tnc", h.GetAllTermsAndConditions).Methods("GET")
}
