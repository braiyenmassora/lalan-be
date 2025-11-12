package handler

import (
	"context"
	"encoding/json"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"net/http"
)

/*
Menangani operasi pelanggan.
Menyediakan endpoint untuk registrasi, login, dan pengambilan profil pelanggan dengan respons sukses atau error.
*/
type CustomerHandler struct {
	customerService service.CustomerService
}

/*
Merepresentasikan data request registrasi pelanggan.
Digunakan untuk decoding JSON dan validasi input.
*/
type CustomerRegisterRequest struct {
	FullName     string `json:"full_name"`
	ProfilePhoto string `json:"profile_photo"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Password     string `json:"password"`
}

/*
Mendaftarkan pelanggan baru.
Mengembalikan respons pembuatan sukses atau error validasi.
*/
func (h *CustomerHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	var req CustomerRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	customer := &model.CustomerModel{
		FullName:     req.FullName,
		ProfilePhoto: req.ProfilePhoto,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Address:      req.Address,
		PasswordHash: req.Password,
	}

	customerResp, err := h.customerService.RegisterCustomer(customer)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, customerResp, message.MsgCustomerCreatedSuccess)
}

/*
Melakukan login pelanggan.
Mengembalikan token atau respons error jika gagal.
*/
func (h *CustomerHandler) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	customerResp, err := h.customerService.LoginCustomer(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}
	response.OK(w, customerResp, message.MsgSuccess)
}

/*
Mengambil profil pelanggan.
Mengembalikan data profil sukses atau error jika tidak ditemukan.
*/
func (h *CustomerHandler) GetCustomerProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	customerID := middleware.GetUserID(r)
	if customerID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	// Buat context dengan userID
	ctx := context.WithValue(context.Background(), "user_id", customerID)

	customer, err := h.customerService.GetCustomerProfile(ctx)
	if err != nil {
		if err.Error() == message.MsgCustomerNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}
	response.OK(w, customer, message.MsgSuccess)
}

/*
Membuat handler pelanggan.
Mengembalikan instance CustomerHandler yang siap digunakan.
*/
func NewCustomerHandler(s service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: s}
}
