// internal/handler/auth_handler.go
package handler

import (
	"encoding/json"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"net/http"
)

// Paket handler untuk handle request HTTP autentikasi.

// Struktur handler autentikasi dengan dependency service.
type AuthHandler struct {
	service service.AuthService
}

// Buat instance handler autentikasi dengan service.
func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Struktur request untuk registrasi hoster.
type RegisterRequest struct {
	FullName     string `json:"full_name"`
	ProfilePhoto string `json:"profile_photo"`
	StoreName    string `json:"store_name"`
	Description  string `json:"description"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Password     string `json:"password"`
}

// Handle registrasi hoster dan kembalikan token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

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

	authResp, err := h.service.Register(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, authResp, message.MsgHosterCreatedSuccess)
}

// Handle login hoster dan kembalikan token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON")
		return
	}

	authResp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}

	response.OK(w, authResp, "Hoster logged in successfully.")
}

// Handle test endpoint terproteksi dengan JWT.
func (h *AuthHandler) TestProtected(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, "Token gagal diverifikasi")
		return
	}

	response.OK(w, map[string]string{
		"message": "Token VALID! Kamu login sebagai:",
		"user_id": userID,
	}, "Middleware JWT berhasil!")
}
