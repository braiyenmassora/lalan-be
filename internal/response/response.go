package response

import (
	"encoding/json"
	"net/http"
)

/*
Response
struct standar untuk respons API
*/
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

/*
writeJSON
helper untuk menulis respons JSON
*/
func writeJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

/*
BadRequest
mengirim respons bad request 400
*/
func BadRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: msg,
		Success: false,
	})
}

/*
Error
mengirim respons error dengan kode tertentu
*/
func Error(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, Response{
		Code:    code,
		Message: msg,
		Success: false,
	})
}

/*
Forbidden
mengirim respons forbidden 403
*/
func Forbidden(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: msg,
		Success: false,
	})
}

/*
Unauthorized
mengirim respons unauthorized 401
*/
func Unauthorized(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Success: false,
	})
}

/*
OK
mengirim respons sukses 200
*/
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

/*
Success
mengirim respons sukses dengan kode tertentu
*/
func Success(w http.ResponseWriter, code int, data any, msg string) {
	writeJSON(w, code, Response{
		Code:    code,
		Data:    data,
		Message: msg,
		Success: true,
	})
}
