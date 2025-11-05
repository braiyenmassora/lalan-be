package helper

import (
	hoster "lalan-be/pkg"
	"net/http"

	"github.com/goccy/go-json"
)

// EncoderJsonData mengirimkan response dalam format JSON.
func EncoderJsonData(w http.ResponseWriter, res hoster.Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}

// SuccessResponse mengirimkan response sukses dalam format JSON.
func SuccessResponse(w http.ResponseWriter, data any, message string, statusCode int) {
	res := hoster.Response{
		Code:    statusCode,
		Data:    data,
		Message: message,
		Status:  "success",
	}
	EncoderJsonData(w, res, statusCode)
}

// ErrorResponse mengirimkan response error dalam format JSON.
func ErrorResponse(w http.ResponseWriter, data any, message string, statusCode int) {
	res := hoster.Response{
		Code:    statusCode,
		Data:    data,
		Message: message,
		Status:  "error",
	}
	EncoderJsonData(w, res, statusCode)
}
