// internal/response/response.go
package response

import (
	"encoding/json"
	"net/http"
)

// Paket response untuk kirim respons HTTP berformat JSON standar.

// Struktur Response untuk format respons API konsisten.
type Response struct {
	Code    int    `json:"code"`           // Kode status HTTP numerik
	Data    any    `json:"data,omitempty"` // Payload respons; dihilangkan jika kosong
	Message string `json:"message"`        // Pesan ringkas menjelaskan status atau error
	Success bool   `json:"success"`        // True jika operasi berhasil
}

// Kirim respons JSON sukses dan set header Content-Type.
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

// Kirim respons JSON error dan set header Content-Type.
func Error(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
		Success: false,
	})
}

// Pintasan untuk kirim respons OK (200) dengan data dan pesan.
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

// Pintasan untuk kirim respons Created (201) dengan data dan pesan.
func Created(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusCreated, data, msg)
}

// Pintasan untuk kirim respons BadRequest (400) dengan pesan error.
func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}

// Pintasan untuk kirim respons Unauthorized (401) dengan pesan error.
func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, msg)
}
