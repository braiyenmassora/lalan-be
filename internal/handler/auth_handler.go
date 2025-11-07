package handler

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"log"
	"net/http"
)

// Struct untuk menangani request HTTP autentikasi
type AuthHandler struct {
	service service.AuthService // Service untuk operasi autentikasi
}

// Membuat instance AuthHandler baru dengan dependency injection
func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Struct untuk payload request JSON register
type RegisterRequest struct {
	FullName     string `json:"full_name"`     // Nama lengkap hoster
	ProfilePhoto string `json:"profile_photo"` // URL foto profil
	StoreName    string `json:"store_name"`    // Nama toko
	Description  string `json:"description"`   // Deskripsi toko
	PhoneNumber  string `json:"phone_number"`  // Nomor telepon
	Email        string `json:"email"`         // Alamat email
	Address      string `json:"address"`       // Alamat
	Password     string `json:"password"`      // Kata sandi
}

// Handler untuk registrasi hoster
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Validasi metode HTTP POST
	if r.Method != http.MethodPost {
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	// Mencegah field tidak dikenal
	decoder.DisallowUnknownFields()
	// Decode payload JSON
	if err := decoder.Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Konversi request ke model domain
	input := &model.HosterModel{
		FullName:     req.FullName,
		ProfilePhoto: req.ProfilePhoto,
		StoreName:    req.StoreName,
		Description:  req.Description,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Address:      req.Address,
		PasswordHash: req.Password,
	}

	// Panggil service untuk registrasi
	if err := h.service.Register(input); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Kembalikan response JSON
	json.NewEncoder(w).Encode(response.Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: message.MsgHosterCreatedSuccess,
		Data:    nil,
	})
}

// Handler untuk login hoster
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Validasi metode HTTP POST
	if r.Method != http.MethodPost {
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Decode payload JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Log request untuk debugging
	log.Printf("Login request: email=%s", req.Email)

	// Panggil service untuk login
	resp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Kembalikan response JSON
	json.NewEncoder(w).Encode(resp)
}
