package customer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"

	"lalan-be/internal/message"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
)

/*
CustomerHandler
menangani permintaan terkait customer melalui service
*/
type CustomerHandler struct {
	service CustomerService
}

/*
CustomerRequest
berisi data untuk membuat customer baru
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

/*
LoginRequest
berisi data untuk login customer
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
UpdateCustomerRequest
berisi data untuk update customer dengan field terbatas
*/
type UpdateCustomerRequest struct {
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	ProfilePhoto string `json:"profile_photo"`
	Address      string `json:"address"`
}

/*
CreateCustomer
membuat customer baru dan mengembalikan data customer
*/
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateCustomer: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req CustomerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateCustomer: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateCustomer: full name required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "full name"))
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateCustomer: email required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "email"))
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateCustomer: password required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "password"))
		return
	}
	input := &model.CustomerModel{
		FullName:     req.FullName,
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
	response.OK(w, input, fmt.Sprintf(message.Created, "customer"))
}

/*
LoginCustomer
melakukan login customer dan mengembalikan token serta data user
*/
func (h *CustomerHandler) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginCustomer: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginCustomer: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginCustomer: email or password empty")
		response.BadRequest(w, message.BadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginCustomer: invalid email format: %s", req.Email)
		response.BadRequest(w, fmt.Sprintf(message.InvalidFormat, "email"))
		return
	}
	resp, err := h.service.LoginCustomer(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginCustomer: login failed: %v", err)
		response.Error(w, http.StatusUnauthorized, message.Unauthorized)
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
	response.OK(w, userData, message.Success)
}

/*
GetDetailCustomer
mengambil detail customer berdasarkan context
*/
func (h *CustomerHandler) GetDetailCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailCustomer: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	ctx := r.Context()
	customer, err := h.service.GetDetailCustomer(ctx)
	if err != nil {
		log.Printf("GetDetailCustomer: error getting customer: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	if customer == nil {
		log.Printf("GetDetailCustomer: customer not found")
		response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "customer"))
		return
	}
	log.Printf("GetDetailCustomer: retrieved customer for ID %s", customer.ID)
	response.OK(w, customer, message.Success)
}

/*
UpdateCustomer
memperbarui data customer dengan field terbatas
*/
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateCustomer: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req UpdateCustomerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateCustomer: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if strings.TrimSpace(req.FullName) == "" && strings.TrimSpace(req.PhoneNumber) == "" &&
		strings.TrimSpace(req.ProfilePhoto) == "" && strings.TrimSpace(req.Address) == "" {
		log.Printf("UpdateCustomer: at least one field must be provided")
		response.BadRequest(w, message.BadRequest)
		return
	}
	updateData := &model.CustomerModel{
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		ProfilePhoto: req.ProfilePhoto,
		Address:      req.Address,
	}
	ctx := r.Context()
	err := h.service.UpdateCustomer(ctx, updateData)
	if err != nil {
		log.Printf("UpdateCustomer: error updating customer: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else if err.Error() == fmt.Sprintf(message.NotFound, "customer") {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "customer"))
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("UpdateCustomer: customer updated successfully")
	response.OK(w, nil, fmt.Sprintf(message.Updated, "customer"))
}

/*
DeleteCustomer
menghapus customer berdasarkan context
*/
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteCustomer: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	ctx := r.Context()
	err := h.service.DeleteCustomer(ctx)
	if err != nil {
		log.Printf("DeleteCustomer: error deleting customer: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else if err.Error() == fmt.Sprintf(message.NotFound, "customer") {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "customer"))
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("DeleteCustomer: customer deleted successfully")
	response.OK(w, nil, fmt.Sprintf(message.Deleted, "customer"))
}

/*
UploadIdentity
mengunggah file identitas customer
*/
func (h *CustomerHandler) UploadIdentity(w http.ResponseWriter, r *http.Request) {
	log.Printf("UploadIdentity: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("UploadIdentity: error parsing multipart form: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	file, header, err := r.FormFile("ktp_file")
	if err != nil {
		log.Printf("UploadIdentity: error getting file: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	defer file.Close()
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		log.Printf("UploadIdentity: invalid file type: %s", header.Header.Get("Content-Type"))
		response.BadRequest(w, fmt.Sprintf(message.InvalidFormat, "file type"))
		return
	}
	ktpURL := "https://storage.example.com/ktp/" + header.Filename
	ctx := r.Context()
	err = h.service.UploadIdentity(ctx, ktpURL)
	if err != nil {
		log.Printf("UploadIdentity: error uploading identity: %v", err)
		if err.Error() == message.IdentityAlreadyUploaded {
			response.BadRequest(w, message.IdentityAlreadyUploaded)
			return
		}
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("UploadIdentity: identity uploaded successfully")
	response.OK(w, nil, fmt.Sprintf(message.Created, "identity"))
}

/*
UpdateIdentity
memperbarui identitas customer dengan upload baru
*/
func (h *CustomerHandler) UpdateIdentity(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateIdentity: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("UpdateIdentity: error parsing multipart form: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	file, header, err := r.FormFile("ktp_file")
	if err != nil {
		log.Printf("UpdateIdentity: error getting file: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	defer file.Close()
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		log.Printf("UpdateIdentity: invalid file type: %s", header.Header.Get("Content-Type"))
		response.BadRequest(w, fmt.Sprintf(message.InvalidFormat, "file type"))
		return
	}
	ktpURL := "https://storage.example.com/ktp/" + header.Filename
	ctx := r.Context()
	err = h.service.UpdateIdentity(ctx, ktpURL)
	if err != nil {
		log.Printf("UpdateIdentity: error updating identity: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else if err.Error() == fmt.Sprintf(message.NotFound, "identity") {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "identity"))
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("UpdateIdentity: identity updated successfully")
	response.OK(w, nil, fmt.Sprintf(message.Updated, "identity"))
}

/*
GetIdentityStatus
mengambil status identitas customer
*/
func (h *CustomerHandler) GetIdentityStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetIdentityStatus: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	ctx := r.Context()
	identity, err := h.service.GetIdentityStatus(ctx)
	if err != nil {
		log.Printf("GetIdentityStatus: error getting identity status: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else if err.Error() == fmt.Sprintf(message.NotFound, "identity") {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "identity"))
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("GetIdentityStatus: retrieved identity status for user")
	response.OK(w, identity, message.Success)
}

/*
CreateBooking
membuat booking baru untuk customer
*/
func (h *CustomerHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateBooking: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req CreateBookingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateBooking: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if req.StartDate == "" || req.EndDate == "" || len(req.Items) == 0 {
		log.Printf("CreateBooking: invalid input")
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	dto, err := h.service.CreateBooking(ctx, req)
	if err != nil {
		log.Printf("CreateBooking: error creating booking: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else if err.Error() == fmt.Sprintf(message.NotFound, "identity") {
			response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "identity"))
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("CreateBooking: booking created successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"data":    dto,
		"message": message.Success,
		"success": true,
	})
}

/*
GetBookingsByUserID
mengambil daftar booking berdasarkan user ID
*/
func (h *CustomerHandler) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetBookingsByUserID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	ctx := r.Context()
	bookings, err := h.service.GetBookingsByUserID(ctx)
	if err != nil {
		log.Printf("GetBookingsByUserID: error getting bookings: %v", err)
		if err.Error() == message.Unauthorized {
			response.Error(w, http.StatusUnauthorized, message.Unauthorized)
		} else {
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}
	log.Printf("GetBookingsByUserID: retrieved %d bookings", len(bookings))
	response.OK(w, bookings, message.Success)
}

/*
GetListBookings
mengambil daftar semua booking
*/
func (h *CustomerHandler) GetListBookings(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookings: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	ctx := r.Context()
	bookings, err := h.service.GetListBookings(ctx)
	if err != nil {
		log.Printf("GetListBookings: error getting bookings: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	log.Printf("GetListBookings: retrieved %d bookings", len(bookings))
	response.OK(w, bookings, message.Success)
}

/*
GetDetailBooking
menangani request untuk mendapatkan detail booking
*/
func (h *CustomerHandler) GetDetailBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID := vars["id"]

	ctx := r.Context()
	bookingDetail, err := h.service.GetDetailBooking(ctx, bookingID)
	if err != nil {
		// Handle error, e.g., send error response
		http.Error(w, message.Unauthorized, http.StatusUnauthorized) // Adjust status code as needed
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"data":    bookingDetail,
		"message": message.Success,
		"success": true,
	})
}

/*
NewCustomerHandler
membuat instance baru CustomerHandler dengan service
*/
func NewCustomerHandler(s CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}
