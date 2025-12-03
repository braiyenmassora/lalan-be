package booking

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"

	"github.com/gorilla/mux"
)

/*
HosterBookingHandler menangani endpoint HTTP untuk fitur booking dari perspektif hoster.
Hanya menyediakan operasi read (list & detail).
*/
type HosterBookingHandler struct {
	service BookingService
}

/*
NewHosterBookingHandler membuat instance handler dengan dependency injection.

Output:
- *HosterBookingHandler siap digunakan
*/
func NewHosterBookingHandler(s BookingService) *HosterBookingHandler {
	return &HosterBookingHandler{service: s}
}

/*
GetListBookings menangani GET /api/v1/hoster/bookings

Alur kerja:
1. Validasi method GET
2. Ambil hosterID dari JWT context
3. Panggil service untuk ambil daftar booking milik hoster

Output sukses:
- 200 OK + list booking ringkas
Output error:
- 401 Unauthorized / 500 Internal Server Error
*/
func (h *HosterBookingHandler) GetListBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	bookings, err := h.service.GetListBookings(hosterID)
	if err != nil {
		log.Printf("GetListBookings handler: service error hoster=%s err=%v", hosterID, err)
		if err.Error() == message.Unauthorized {
			response.Unauthorized(w, message.Unauthorized)
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, bookings, message.Success)
}

/*
GetDetailBooking menangani GET /api/v1/hoster/bookings/{bookingID}

Alur kerja:
1. Validasi method GET
2. Ambil bookingID dari path parameter
3. Ambil hosterID dari JWT context
4. Panggil service (sudah termasuk authorization check)

Output sukses:
- 200 OK + detail booking lengkap
Output error:
- 400 Bad Request (ID kosong)
- 401 Unauthorized (bukan pemilik)
- 404 Not Found
- 500 Internal Server Error
*/
func (h *HosterBookingHandler) GetDetailBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	bookingID := strings.TrimSpace(vars["bookingID"])
	if bookingID == "" {
		// route may use {id} — support both variants
		bookingID = strings.TrimSpace(vars["id"])
	}
	if bookingID == "" {
		response.BadRequest(w, fmt.Sprintf(message.Required, "booking ID"))
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	detail, err := h.service.GetDetailBooking(hosterID, bookingID)
	if err != nil {
		log.Printf("GetDetailBooking handler: service error booking=%s hoster=%s err=%v", bookingID, hosterID, err)
		switch err.Error() {
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case fmt.Sprintf(message.NotFound, "booking"):
			response.NotFound(w, fmt.Sprintf(message.NotFound, "booking"))
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, detail, message.Success)
}

// GetCustomerList menangani GET /api/v1/hoster/booking/customers
// Mengembalikan daftar pelanggan yang melakukan booking pada hoster yang sedang login
func (h *HosterBookingHandler) GetCustomerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	list, err := h.service.GetCustomerList(hosterID)
	if err != nil {
		log.Printf("GetCustomerList handler: service error hoster=%s err=%v", hosterID, err)
		if err.Error() == message.Unauthorized {
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, list, message.Success)
}

/*
UpdateBookingStatus menangani PUT /api/v1/hoster/booking/{id}/status

Alur kerja:
1. Validasi method PUT
2. Ambil bookingID dari path parameter
3. Ambil hosterID dari JWT context
4. Parse request body untuk status baru
5. Panggil service (include authorization & validation)

Output sukses:
- 200 OK + message sukses
Output error:
- 400 Bad Request (ID kosong / invalid JSON / invalid status)
- 401 Unauthorized (bukan pemilik)
- 404 Not Found (booking tidak ada)
- 500 Internal Server Error
*/
func (h *HosterBookingHandler) UpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	// Get bookingID from path
	vars := mux.Vars(r)
	bookingID := strings.TrimSpace(vars["id"])
	if bookingID == "" {
		response.BadRequest(w, fmt.Sprintf(message.Required, "booking ID"))
		return
	}

	// Get hosterID from JWT
	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	// Parse request body
	var req dto.UpdateBookingStatusByHosterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateBookingStatus: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi status tidak kosong
	if strings.TrimSpace(req.Status) == "" {
		response.BadRequest(w, fmt.Sprintf(message.Required, "status"))
		return
	}

	// Call service
	err := h.service.UpdateBookingStatus(hosterID, bookingID, req.Status)
	if err != nil {
		log.Printf("UpdateBookingStatus: service error booking=%s status=%s err=%v", bookingID, req.Status, err)

		switch err.Error() {
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case fmt.Sprintf(message.NotFound, "booking"):
			response.NotFound(w, fmt.Sprintf(message.NotFound, "booking"))
		case message.InvalidStatus:
			response.BadRequest(w, message.InvalidStatus)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, nil, message.BookingStatusUpdated)
}
