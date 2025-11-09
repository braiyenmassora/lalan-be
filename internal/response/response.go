// internal/response/response.go
package response

import (
	"encoding/json"
	"net/http"
)

/*
Struct untuk struktur respons API.
Merepresentasikan data respons JSON.
*/
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

/*
Mengirim respons JSON sukses.
Mengatur header dan mengencode respons.
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
Mengirim respons JSON error.
Mengatur header dan mengencode respons.
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
Mengirim respons OK dengan data.
Menggunakan fungsi Success dengan status OK.
*/
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

/*
Mengirim respons Created dengan data.
Menggunakan fungsi Success dengan status Created.
*/
func Created(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusCreated, data, msg)
}

/*
Mengirim respons BadRequest dengan pesan.
Menggunakan fungsi Error dengan status BadRequest.
*/
func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}

/*
Mengirim respons Unauthorized dengan pesan.
Menggunakan fungsi Error dengan status Unauthorized.
*/
func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, msg)
}
