package middleware

import (
	"net/http"

	"lalan-be/internal/message"
	"lalan-be/internal/response"
)

/*
Admin
memeriksa akses role admin dan melanjutkan jika valid
*/
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserRole(r) != "admin" {
			response.Forbidden(w, message.AdminAccessRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
Hoster
memeriksa akses role hoster dan melanjutkan jika valid
*/
func Hoster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserRole(r) != "hoster" {
			response.Forbidden(w, message.HosterAccessRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
Customer
memeriksa akses role customer dan melanjutkan jika valid
*/
func Customer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserRole(r) != "customer" {
			response.Forbidden(w, message.CustomerAccessRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}
