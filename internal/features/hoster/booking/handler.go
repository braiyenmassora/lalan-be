package booking

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"

	"github.com/gorilla/mux"
)

/*
BookingHandler menangani endpoint HTTP untuk fitur booking dari perspektif hoster.
Hanya menyediakan operasi read (list & detail).
*/
type BookingHandler struct {
	service BookingService
}

/*
NewBookingHandler membuat instance handler dengan dependency injection.

Output:
- *BookingHandler siap digunakan
*/
func NewBookingHandler(s BookingService) *BookingHandler {
	return &BookingHandler{service: s}
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
func (h *BookingHandler) GetListBookings(w http.ResponseWriter, r *http.Request) {
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
func (h *BookingHandler) GetDetailBooking(w http.ResponseWriter, r *http.Request) {
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
func (h *BookingHandler) GetCustomerList(w http.ResponseWriter, r *http.Request) {
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
