package middleware

import (
	"net/http"

	"lalan-be/internal/response"
)

/*
Admin memeriksa akses role admin.
Melanjutkan jika role admin, forbidden jika tidak.
*/
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek role admin
		if GetUserRole(r) != "admin" {
			response.Forbidden(w, "Admin access required")
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
			response.Forbidden(w, "Hoster access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
