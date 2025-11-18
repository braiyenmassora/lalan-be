package customer

import (
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
SetupCustomerRoutes mengatur rute untuk customer.
Mendaftarkan endpoint public dan protected dengan middleware.
*/
func SetupCustomerRoutes(router *mux.Router, handler *CustomerHandler) {
	public := router.PathPrefix("/api/v1/customer").Subrouter()
	public.HandleFunc("/register", handler.CreateCustomer).Methods("POST")
	public.HandleFunc("/login", handler.LoginCustomer).Methods("POST")

	protected := router.PathPrefix("/api/v1/customer").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Customer)

	protected.HandleFunc("/detail", handler.GetDetailCustomer).Methods("GET")
	protected.HandleFunc("/update", handler.UpdateCustomer).Methods("PUT")
	protected.HandleFunc("/delete", handler.DeleteCustomer).Methods("DELETE")
	protected.HandleFunc("/upload-identity", handler.UploadIdentity).Methods("POST")
	protected.HandleFunc("/identity-status", handler.GetIdentityStatus).Methods("GET")

}
