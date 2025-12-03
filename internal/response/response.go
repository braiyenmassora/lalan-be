package response

import (
	"encoding/json"
	"net/http"
)

/*
Response adalah format standar JSON yang digunakan oleh seluruh API.
Semua response (sukses/error) harus mengikuti struktur ini agar frontend konsisten.
*/
type Response struct {
	Code         int    `json:"code"`                    // HTTP status code
	Data         any    `json:"data,omitempty"`          // Payload (hanya ada jika success)
	Message      string `json:"message"`                 // Pesan untuk user / developer
	Success      bool   `json:"success"`                 // true = sukses, false = error
	ErrorDetails any    `json:"error_details,omitempty"` // Detail error tambahan (optional)
}

/*
writeJSON adalah helper internal untuk menulis header JSON + encode struct Response.
Dipanggil oleh semua fungsi publik di package ini.
*/
func writeJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

/*
BadRequest mengembalikan response 400 Bad Request.

Output:
- Status: 400
- Body: { "code": 400, "message": msg, "success": false }
*/
func BadRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: msg,
		Success: false,
	})
}

/*
BadRequestWithDetails mengembalikan response 400 Bad Request dengan detail tambahan.

Output:
- Status: 400
- Body: { "code": 400, "message": msg, "error_details": details, "success": false }
*/
func BadRequestWithDetails(w http.ResponseWriter, msg string, details any) {
	writeJSON(w, http.StatusBadRequest, Response{
		Code:         http.StatusBadRequest,
		Message:      msg,
		ErrorDetails: details,
		Success:      false,
	})
}

/*
Error mengembalikan response error dengan status code custom (contoh: 500, 422, dll).

Output:
- Status: code yang diberikan
- Body: { "code": code, "message": msg, "success": false }
*/
func Error(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, Response{
		Code:    code,
		Message: msg,
		Success: false,
	})
}

/*
Forbidden mengembalikan response 403 Forbidden (biasanya role tidak cukup).

Output:
- Status: 403
- Body: { "code": 403, "message": msg, "success": false }
*/
func Forbidden(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: msg,
		Success: false,
	})
}

/*
Unauthorized mengembalikan response 401 Unauthorized (token salah / expired / missing).

Output:
- Status: 401
- Body: { "code": 401, "message": msg, "success": false }
*/
func Unauthorized(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Success: false,
	})
}

/*
OK mengembalikan response sukses 200 OK (paling sering dipakai).

Output:
- Status: 200
- Body: { "code": 200, "data": data, "message": msg, "success": true }
*/
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

/*
Success mengembalikan response sukses dengan status code custom (contoh: 201 Created).

Output:
- Status: code yang diberikan
- Body: { "code": code, "data": data, "message": msg, "success": true }
*/
func Success(w http.ResponseWriter, code int, data any, msg string) {
	writeJSON(w, code, Response{
		Code:    code,
		Data:    data,
		Message: msg,
		Success: true,
	})
}

/*
NotFound mengembalikan response 404 Not Found.

Output:
- Status: 404
- Body: { "code": 404, "message": msg, "success": false }
*/
func NotFound(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: msg,
		Success: false,
	})
}

/*
MethodNotAllowed mengembalikan response 405 Method Not Allowed.

Output:
- Status: 405
- Body: { "code": 405, "message": msg, "success": false }
*/
func MethodNotAllowed(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusMethodNotAllowed, Response{
		Code:    http.StatusMethodNotAllowed,
		Message: msg,
		Success: false,
	})
}
