package tnc

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"
)

/*
HosterTnCHandler menangani endpoint HTTP untuk T&C dari perspektif hoster.
*/
type HosterTnCHandler struct {
	service TnCService
}

/*
NewHosterTnCHandler membuat instance handler dengan dependency injection.

Output:
- *HosterTnCHandler siap digunakan
*/
func NewHosterTnCHandler(s TnCService) *HosterTnCHandler {
	return &HosterTnCHandler{service: s}
}

/*
CreateTnC menangani POST /api/v1/hoster/tnc

Alur kerja:
1. Validasi method POST
2. Ambil userID dari JWT context
3. Parse JSON body
4. Panggil service
5. Return response

Output sukses:
- 201 Created + T&C yang baru dibuat
Output error:
- 400 Bad Request / 401 Unauthorized / 500 Internal Server Error
*/
func (h *HosterTnCHandler) CreateTnC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	var req dto.CreateTnCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateTnC: failed to decode JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	result, err := h.service.CreateTnC(hosterID, &req)
	if err != nil {
		log.Printf("CreateTnC handler: service error hoster=%s err=%v", hosterID, err)
		switch err.Error() {
		case message.BadRequest:
			response.BadRequest(w, message.BadRequest)
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, result, message.Success)
}

/*
UpdateTnC menangani PUT /api/v1/hoster/tnc/{id}

Alur kerja:
1. Validasi method PUT
2. Ambil userID dari JWT context
3. Ambil tncID dari path parameter
4. Parse JSON body
5. Panggil service
6. Return response

Output sukses:
- 200 OK + T&C yang diupdate
Output error:
- 400 Bad Request / 401 Unauthorized / 404 Not Found / 500 Internal Server Error
*/
func (h *HosterTnCHandler) UpdateTnC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	vars := mux.Vars(r)
	tncID := vars["id"]
	if tncID == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	var req dto.UpdateTnCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateTnC: failed to decode JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	result, err := h.service.UpdateTnC(hosterID, tncID, &req)
	if err != nil {
		log.Printf("UpdateTnC handler: service error hoster=%s tnc=%s err=%v", hosterID, tncID, err)
		switch err.Error() {
		case message.BadRequest:
			response.BadRequest(w, message.BadRequest)
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case message.NotFound:
			response.NotFound(w, message.TnCNotFound)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, result, message.TnCUpdated)
}

/*
GetTnC menangani GET /api/v1/hoster/tnc

Alur kerja:
1. Validasi method GET
2. Ambil userID dari JWT context
3. Panggil service
4. Return response

Output sukses:
- 200 OK + T&C data
Output error:
- 401 Unauthorized / 404 Not Found / 500 Internal Server Error
*/
func (h *HosterTnCHandler) GetTnC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	result, err := h.service.GetTnC(hosterID)
	if err != nil {
		log.Printf("GetTnC handler: service error hoster=%s err=%v", hosterID, err)
		switch err.Error() {
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case message.NotFound:
			response.NotFound(w, message.TnCNotFound)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, result, message.TnCRetrieved)
}
