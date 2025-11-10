// internal/response/response.go
package response

import (
	"encoding/json"
	"net/http"
)

/*
Merepresentasikan struktur respons API.
Digunakan untuk encoding JSON dengan field kode, data, pesan, dan status sukses.
*/
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

/*
Mengirim respons sukses.
Mengembalikan respons JSON dengan data dan pesan sukses.
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
Mengirim respons error.
Mengembalikan respons JSON dengan pesan error.
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
Mengirim respons OK.
Mengembalikan respons JSON dengan status OK dan data.
*/
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

/*
Mengirim respons Created.
Mengembalikan respons JSON dengan status Created dan data.
*/
func Created(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusCreated, data, msg)
}

/*
Mengirim respons BadRequest.
Mengembalikan respons JSON dengan status BadRequest dan pesan.
*/
func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}

/*
Mengirim respons Unauthorized.
Mengembalikan respons JSON dengan status Unauthorized dan pesan.
*/
func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, msg)
}
