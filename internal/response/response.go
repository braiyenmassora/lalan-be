package response

import (
	"encoding/json"
	"net/http"
)

/*
Mewakili struktur respons API standar.
Digunakan untuk format JSON respons sukses atau error.
*/
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

/*
Mengirim respons bad request dengan pesan error.
Mengembalikan status 400 dan JSON error.
*/
func BadRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Code:    http.StatusBadRequest,
		Message: msg,
		Success: false,
	})
}

/*
Mengirim respons error dengan kode dan pesan.
Mengembalikan JSON dengan success false.
*/
func Error(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
		Success: false,
	})
}

/*
Mengirim respons forbidden dengan pesan.
Mengembalikan status 403 dan JSON error.
*/
func Forbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(Response{
		Code:    http.StatusForbidden,
		Message: message,
		Success: false,
	})
}

/*
Mengirim respons OK dengan data dan pesan.
Mengembalikan status 200 dan JSON sukses.
*/
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

/*
Mengirim respons sukses dengan kode, data, dan pesan.
Mengembalikan JSON dengan success true.
*/
func Success(w http.ResponseWriter, code int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Data:    data,
		Message: message,
		Success: true,
	})
}

/*
Mengirim respons unauthorized dengan pesan.
Mengembalikan status 401 dan JSON error.
*/
func Unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Success: false,
	})
}
