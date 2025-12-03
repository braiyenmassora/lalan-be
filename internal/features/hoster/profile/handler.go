package profile

import (
	"encoding/json"
	"log"
	"net/http"

	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"
)

/*
HosterProfileHandler menangani endpoint HTTP untuk profile hoster.
*/
type HosterProfileHandler struct {
	service HosterProfileService
}

/*
NewHosterProfileHandler membuat instance handler dengan dependency injection.

Output:
- *HosterProfileHandler siap digunakan
*/
func NewHosterProfileHandler(s HosterProfileService) *HosterProfileHandler {
	return &HosterProfileHandler{service: s}
}

/*
GetProfile menangani GET /api/v1/hoster/profile

Alur kerja:
1. Validasi method GET
2. Ambil userID dari JWT context
3. Panggil service
4. Return response

Output sukses:
- 200 OK + profile data
Output error:
- 401 Unauthorized / 404 Not Found / 500 Internal Server Error
*/
func (h *HosterProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	result, err := h.service.GetProfile(hosterID)
	if err != nil {
		log.Printf("GetProfile handler: service error hoster=%s err=%v", hosterID, err)
		switch err.Error() {
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case message.HosterNotFound:
			response.NotFound(w, message.ProfileNotFound)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, result, message.ProfileRetrieved)
}

/*
UpdateProfile menangani PUT /api/v1/hoster/profile

Alur kerja:
1. Validasi method PUT
2. Ambil userID dari JWT context
3. Parse JSON body
4. Panggil service
5. Return response

Output sukses:
- 200 OK + updated profile
Output error:
- 400 Bad Request / 401 Unauthorized / 404 Not Found / 500 Internal Server Error
*/
func (h *HosterProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	var req dto.UpdateHosterProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateProfile: failed to decode JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	result, err := h.service.UpdateProfile(hosterID, &req)
	if err != nil {
		log.Printf("UpdateProfile handler: service error hoster=%s err=%v", hosterID, err)
		switch err.Error() {
		case message.BadRequest:
			response.BadRequest(w, message.BadRequest)
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case message.HosterNotFound:
			response.NotFound(w, message.ProfileNotFound)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, result, message.ProfileUpdated)
}
