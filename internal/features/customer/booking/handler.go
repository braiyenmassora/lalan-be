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
BookingHandler adalah HTTP layer (transport layer) untuk fitur booking.
Tanggung jawab handler TERBATAS pada:
• Validasi HTTP method dan format request
• Parsing dan validasi input (JSON body, path params)
• Memanggil BookingService sebagai business logic entry point
• Mapping error dari service ke HTTP status code + response yang sesuai
Seluruh aturan bisnis, validasi domain, dan keputusan authorization
harus tetap berada di layer service — handler tidak boleh mengandung business rule.
*/
type BookingHandler struct {
	service BookingService
}

/*
NewBookingHandler membuat instance BookingHandler yang sudah terinject
dengan dependency BookingService. Digunakan saat setup router.
*/
func NewBookingHandler(s BookingService) *BookingHandler {
	return &BookingHandler{service: s}
}

/*
CreateBooking menangani request pembuatan booking baru.

Langkah-langkah:
1. Validasi method & decode JSON.
2. Validasi input (misal item tidak kosong, tanggal valid).
3. Panggil service CreateBooking.
4. Jika error spesifik "silakan upload ktp terlebih dahulu", return 400 Bad Request.
5. Jika error lain, return 500 Internal Server Error.

Output:
- 200 OK: Booking berhasil dibuat, return detail booking.
- 400 Bad Request: Validasi gagal (misal KTP belum upload).
- 500 Internal Server Error: Kesalahan sistem.
*/
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateBooking: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.CreateBookingByCustomerRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Printf("CreateBooking: decode error: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Ambil userID dari middleware context
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	// Validasi input di sini jika perlu...

	resp, err := h.service.CreateBooking(userID, req)
	if err != nil {
		log.Printf("CreateBooking: service error: %v", err)
		errMsg := err.Error()

		// Cek exact match untuk error constants
		if errMsg == message.KTPRequired {
			response.BadRequest(w, message.KTPRequired)
			return
		}
		if errMsg == message.KTPRejectedUploadNew {
			response.BadRequest(w, message.KTPRejectedUploadNew)
			return
		}

		// Semua error lain dari service = BadRequest (business logic error)
		response.BadRequest(w, errMsg)
		return
	}

	response.OK(w, resp, message.Success)
}

/*
GetBookingsByUserID menangani endpoint GET /users/me/bookings
Mengembalikan seluruh daftar booking milik user yang sedang login.

Alur kerja:
1. Validasi method GET
2. Ambil user context yang sudah divalidasi oleh middleware auth
3. Panggil service untuk mengambil daftar booking
4. Mapping error:
  - Unauthorized → 401 (token invalid/expired)
  - Lainnya     → 500

Output sukses:
- Status: 200 OK
- Body:   array of booking summary
- Message: "Success"
*/
func (h *BookingHandler) GetListBookings(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookings: received request")
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	// Ambil userID dari middleware context agar konsisten dengan handler lain
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	bookings, err := h.service.GetListBookings(userID)
	if err != nil {
		log.Printf("GetListBookings: service error: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("Handler GetListBookings: bookings data: %+v", bookings)
	response.OK(w, bookings, message.Success)
}

/*
GetDetailBooking menangani endpoint GET /bookings/{id}
Mengembalikan detail lengkap satu booking berdasarkan ID.

Alur kerja:
1. Validasi method GET
2. Ekstrak dan validasi path parameter "id" (tidak boleh kosong)
3. Panggil service yang akan melakukan:
  - Pencarian booking
  - Pengecekan kepemilikan (user harus pemilik booking)

4. Mapping error:
  - Unauthorized → 401 (bukan pemilik)
  - Not Found    → 404 (booking tidak ada)
  - Lainnya      → 500

Output sukses:
- Status: 200 OK
- Body:   object detail booking lengkap (termasuk payment, schedule, dll)
- Message: "Success"
*/
func (h *BookingHandler) GetDetailBooking(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailBooking: received request")
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	bookingID := strings.TrimSpace(vars["id"])
	if bookingID == "" {
		response.BadRequest(w, fmt.Sprintf(message.Required, "booking ID"))
		return
	}

	// Ambil userID dari middleware untuk konsistensi handler pattern
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	bookingDetail, err := h.service.GetDetailBooking(userID, bookingID)
	if err != nil {
		log.Printf("GetDetailBooking: service error: %v", err)
		switch err.Error() {
		case message.Unauthorized:
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		case fmt.Sprintf(message.NotFound, "booking"):
			response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "booking"))
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, bookingDetail, message.Success)
}
