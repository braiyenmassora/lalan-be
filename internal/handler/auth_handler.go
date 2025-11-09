package handler

import (
	"encoding/json"
	"net/http"

	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
)

/*
Struct untuk menangani request autentikasi.
Mengandung service untuk operasi bisnis.
*/
type AuthHandler struct {
	service service.AuthService
}

/*
Membuat instance handler dengan service autentikasi.
Mengembalikan pointer ke AuthHandler.
*/
func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

/*
Struct untuk merepresentasikan request registrasi hoster.
*/
type RegisterRequest struct {
	FullName     string `json:"full_name"`
	ProfilePhoto string `json:"profile_photo"`
	StoreName    string `json:"store_name"`
	Website      string `json:"website"`
	Instagram    string `json:"instagram"`
	Tiktok       string `json:"tiktok"`
	Description  string `json:"description"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Password     string `json:"password"`
}

/*
Menangani request registrasi hoster.
Mengembalikan respons sukses atau error.
*/
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
		Website:      req.Website,
		Instagram:    req.Instagram,
		Tiktok:       req.Tiktok,
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

/*
Menangani request login hoster.
Mengembalikan respons sukses atau error.
*/
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
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

/*
Menangani request test endpoint terproteksi.
Mengembalikan informasi user jika token valid.
*/
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
