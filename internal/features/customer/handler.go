package customer

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
	"log"
	"net/http"
	"regexp"
	"strings"
	// Tambahkan untuk image processing
)

/*
CustomerHandler menangani permintaan terkait customer.
Menyediakan layanan untuk operasi customer melalui service.
*/
type CustomerHandler struct {
	service CustomerService
}

/*
CustomerRequest berisi data untuk membuat customer baru.
Digunakan dalam permintaan pembuatan customer.
*/
type CustomerRequest struct {
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	ProfilePhoto string `json:"profile_photo"`
	Website      string `json:"website"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
UpdateCustomerRequest berisi data untuk update customer.
Hanya field yang diizinkan: full_name, phone_number, profile_photo, address.
*/
type UpdateCustomerRequest struct {
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	ProfilePhoto string `json:"profile_photo"`
	Address      string `json:"address"`
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateCustomer: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req CustomerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateCustomer: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateCustomer: full name required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateCustomer: email required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateCustomer: password required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	input := &model.CustomerModel{
		FullName: req.FullName,

		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		PasswordHash: req.Password,
		Address:      req.Address,
		ProfilePhoto: req.ProfilePhoto,
	}
	err := h.service.CreateCustomer(input)
	if err != nil {
		log.Printf("CreateCustomer: error creating customer: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, input, message.MsgSuccess)
}
func (h *CustomerHandler) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginCustomer: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginCustomer: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginCustomer	: email or password empty")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginHoster: invalid email format: %s", req.Email)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	resp, err := h.service.LoginCustomer(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginCustomer: login failed: %v", err)
		response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		return
	}
	log.Printf("LoginCustomer: login successful for email %s", req.Email)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    resp.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600,
	})
	userData := map[string]interface{}{
		"id":            resp.ID,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"token_type":    resp.TokenType,
		"expires_in":    resp.ExpiresIn,
	}
	response.OK(w, userData, message.MsgSuccess)
}
func (h *CustomerHandler) GetDetailCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailCustomer: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	ctx := r.Context()
	customer, err := h.service.GetDetailCustomer(ctx)
	if err != nil {
		log.Printf("GetDetailCustomer: error getting customer: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}
	if customer == nil {
		log.Printf("GetDetailCustomer: customer not found")
		response.Error(w, http.StatusNotFound, message.MsgCustomerNotFound)
		return
	}
	log.Printf("GetDetailCustomer: retrieved customer for ID %s", customer.ID)
	response.OK(w, customer, message.MsgSuccess)
}
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateCustomer: received request")
	if r.Method != http.MethodPut { // Gunakan PUT untuk update
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	var req UpdateCustomerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateCustomer: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi: Pastikan setidaknya satu field diisi (opsional, tergantung kebutuhan)
	if strings.TrimSpace(req.FullName) == "" && strings.TrimSpace(req.PhoneNumber) == "" &&
		strings.TrimSpace(req.ProfilePhoto) == "" && strings.TrimSpace(req.Address) == "" {
		log.Printf("UpdateCustomer: at least one field must be provided")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Buat model dari request
	updateData := &model.CustomerModel{
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		ProfilePhoto: req.ProfilePhoto,
		Address:      req.Address,
	}

	// Panggil service dengan context
	ctx := r.Context()
	err := h.service.UpdateCustomer(ctx, updateData)
	if err != nil {
		log.Printf("UpdateCustomer: error updating customer: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else if err.Error() == message.MsgCustomerNotFound {
			response.Error(w, http.StatusNotFound, message.MsgCustomerNotFound)
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("UpdateCustomer: customer updated successfully")
	response.OK(w, nil, message.MsgSuccess) // Return sukses tanpa data tambahan
}
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteCustomer: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Panggil service dengan context
	ctx := r.Context()
	err := h.service.DeleteCustomer(ctx)
	if err != nil {
		log.Printf("DeleteCustomer: error deleting customer: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else if err.Error() == message.MsgCustomerNotFound {
			response.Error(w, http.StatusNotFound, message.MsgCustomerNotFound)
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("DeleteCustomer: customer deleted successfully")
	response.OK(w, nil, message.MsgSuccess) // Return sukses tanpa data tambahan
}

func (h *CustomerHandler) UploadIdentity(w http.ResponseWriter, r *http.Request) {
	log.Printf("UploadIdentity: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		log.Printf("UploadIdentity: error parsing multipart form: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	file, header, err := r.FormFile("ktp_file")
	if err != nil {
		log.Printf("UploadIdentity: error getting file: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	defer file.Close()

	// Validasi file type (opsional)
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		log.Printf("UploadIdentity: invalid file type: %s", header.Header.Get("Content-Type"))
		response.BadRequest(w, "Invalid file type")
		return
	}

	// Upload file ke storage (misalnya S3, sini simulasikan URL)
	// Ganti dengan logic upload real
	ktpURL := "https://storage.example.com/ktp/" + header.Filename // Simulasi URL

	// Panggil service dengan URL
	ctx := r.Context()
	err = h.service.UploadIdentity(ctx, ktpURL)
	if err != nil {
		log.Printf("UploadIdentity: error uploading identity: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("UploadIdentity: identity uploaded successfully")
	response.OK(w, nil, "KTP uploaded successfully")
}

func (h *CustomerHandler) GetIdentityStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetIdentityStatus: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Panggil service
	ctx := r.Context()
	identity, err := h.service.GetIdentityStatus(ctx)
	if err != nil {
		log.Printf("GetIdentityStatus: error getting identity status: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else if err.Error() == "Identity not found" {
			response.Error(w, http.StatusNotFound, "Identity not found")
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("GetIdentityStatus: retrieved identity status for user")
	response.OK(w, identity, message.MsgSuccess)
}

func (h *CustomerHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateBooking: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	var req CreateBookingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateBooking: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input (opsional, tambahkan jika perlu)
	if req.StartDate == "" || req.EndDate == "" || len(req.Items) == 0 {
		log.Printf("CreateBooking: invalid input")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Panggil service
	ctx := r.Context()
	dto, err := h.service.CreateBooking(ctx, req)
	if err != nil {
		log.Printf("CreateBooking: error creating booking: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else if err.Error() == "Identity not found" {
			response.Error(w, http.StatusNotFound, "Identity not found")
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("CreateBooking: booking created successfully")
	response.OK(w, dto, message.MsgSuccess)
}

func (h *CustomerHandler) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetBookingsByUserID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Panggil service
	ctx := r.Context()
	bookings, err := h.service.GetBookingsByUserID(ctx)
	if err != nil {
		log.Printf("GetBookingsByUserID: error getting bookings: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("GetBookingsByUserID: retrieved %d bookings", len(bookings))
	response.OK(w, bookings, message.MsgSuccess)
}

func (h *CustomerHandler) GetListBookings(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookings: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Panggil service
	ctx := r.Context()
	bookings, err := h.service.GetListBookings(ctx)
	if err != nil {
		log.Printf("GetListBookings: error getting bookings: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}

	log.Printf("GetListBookings: retrieved %d bookings", len(bookings))
	response.OK(w, bookings, message.MsgSuccess)
}

/*
NewCustomerHandler membuat instance baru CustomerHandler.
Menginisialisasi handler dengan service yang diberikan.
*/

func NewCustomerHandler(s CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}
