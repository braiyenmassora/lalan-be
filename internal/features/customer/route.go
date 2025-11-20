package customer

import (
	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
SetupCustomerRoutes
mengatur rute public dan protected untuk customer dengan middleware
*/
func SetupCustomerRoutes(router *mux.Router, handler *CustomerHandler) {
	public := router.PathPrefix("/api/v1/customer").Subrouter()
	public.HandleFunc("/register", handler.CreateCustomer).Methods("POST")
	public.HandleFunc("/login", handler.LoginCustomer).Methods("POST")
	public.HandleFunc("/verify-otp", handler.VerifyEmail).Methods("POST")
	public.HandleFunc("/resend-otp", handler.ResendOTP).Methods("POST")

	protected := router.PathPrefix("/api/v1/customer").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.Use(middleware.Customer)

	protected.HandleFunc("/detail", handler.GetDetailCustomer).Methods("GET")
	protected.HandleFunc("/update", handler.UpdateCustomer).Methods("PUT")
	protected.HandleFunc("/delete", handler.DeleteCustomer).Methods("DELETE")
	protected.HandleFunc("/identity", handler.UploadIdentity).Methods("POST")
	protected.HandleFunc("/identity", handler.UpdateIdentity).Methods("PUT")
	protected.HandleFunc("/identity", handler.GetIdentityStatus).Methods("GET")
	protected.HandleFunc("/booking", handler.CreateBooking).Methods("POST")
	protected.HandleFunc("/booking", handler.GetListBookings).Methods("GET")
	protected.HandleFunc("/booking/{id}", handler.GetDetailBooking).Methods("GET")
}
