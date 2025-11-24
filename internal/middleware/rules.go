package middleware

import (
	"log"
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
		role := GetUserRole(r)
		userID := GetUserID(r)
		log.Printf("Admin middleware: user_role = %s, user_id = %s", role, userID)
		if role != "admin" {
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
		role := GetUserRole(r)
		userID := GetUserID(r)
		log.Printf("Hoster middleware: user_role = %s, user_id = %s", role, userID)
		if role != "hoster" {
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
		role := GetUserRole(r)
		userID := GetUserID(r)
		log.Printf("Customer middleware: user_role = %s, user_id = %s", role, userID)
		if role != "customer" {
			response.Forbidden(w, message.CustomerAccessRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}
