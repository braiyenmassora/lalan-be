package middleware

import (
	"net/http"

	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

/*
Admin memeriksa akses role admin.
Melanjutkan jika role admin, forbidden jika tidak.
*/
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek role admin
		if GetUserRole(r) != "admin" {
			response.Forbidden(w, message.MsgAdminAccessRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
Hoster memeriksa akses role hoster.
Melanjutkan jika role hoster, forbidden jika tidak.
*/
func Hoster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek role hoster
		if GetUserRole(r) != "hoster" {
			response.Forbidden(w, message.MsgHosterAccessRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Customer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cek role customer
		if GetUserRole(r) != "customer" {
			response.Forbidden(w, message.MsgCustomerAccessRequired)
		}
		next.ServeHTTP(w, r)
	})
}
