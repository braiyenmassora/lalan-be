package hoster

import (
	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
SetupHosterRoutes
mengatur rute untuk fitur hoster dengan middleware
*/
func SetupHosterRoutes(router *mux.Router, handler *HosterHandler) {
	public := router.PathPrefix("/api/v1/hoster").Subrouter()
	public.HandleFunc("/register", handler.CreateHoster).Methods("POST")
	public.HandleFunc("/login", handler.LoginHoster).Methods("POST")

	protected := router.PathPrefix("/api/v1/hoster").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Hoster)

	protected.HandleFunc("/detail", handler.GetDetailHoster).Methods("GET")

	protected.HandleFunc("/item", handler.CreateItem).Methods("POST")
	protected.HandleFunc("/item", handler.GetAllItems).Methods("GET")
	protected.HandleFunc("/item/{id}", handler.GetItemByID).Methods("GET")
	protected.HandleFunc("/item/{id}", handler.UpdateItem).Methods("PUT")
	protected.HandleFunc("/item/{id}", handler.DeleteItem).Methods("DELETE")

	protected.HandleFunc("/tnc", handler.CreateTermsAndConditions).Methods("POST")
	protected.HandleFunc("/tnc", handler.GetAllTermsAndConditions).Methods("GET")
	protected.HandleFunc("/tnc/{id}", handler.UpdateTermsAndConditions).Methods("PUT")
	protected.HandleFunc("/tnc/{id}", handler.DeleteTermsAndConditions).Methods("DELETE")

	protected.HandleFunc("/booking", handler.GetListBookingsCustomer).Methods("GET")
	protected.HandleFunc("/booking/{bookingID}", handler.GetListBookingsCustomerByBookingID).Methods("GET")

	protected.HandleFunc("/customer", handler.GetListCustomer).Methods("GET")
}
