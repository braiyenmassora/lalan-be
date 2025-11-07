package response

import (
	"encoding/json"
	"net/http"
)

// Struct untuk format response JSON standar
type Response struct {
	Code    int    `json:"code"`           // Kode status HTTP
	Data    any    `json:"data,omitempty"` // Data respons (opsional)
	Message string `json:"message"`        // Pesan respons
	Status  string `json:"status"`         // Status respons
}

// Mengirim response JSON dengan format standar
func JSON(w http.ResponseWriter, code int, status string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Mengirim response sukses standar (HTTP 200)
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, "success", "OK", data)
}

// Mengirim response error dengan kode dan pesan khusus
func ErrorResponse(w http.ResponseWriter, code int, message string) {
	JSON(w, code, "Error", message, nil)
}
