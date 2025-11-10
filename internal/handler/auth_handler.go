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
Menangani autentikasi hoster.
Menyediakan operasi registrasi, login, dan pengambilan profil dengan respons sukses atau error.
*/
type AuthHandler struct {
	service service.AuthService
}

/*
Merepresentasikan data registrasi hoster.
Digunakan untuk decoding request JSON dan mapping ke model.
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
Membuat handler autentikasi.
Mengembalikan instance AuthHandler yang siap digunakan.
*/
func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

/*
Memproses registrasi hoster.
Mengembalikan respons pembuatan akun sukses atau error validasi.
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
Memproses login hoster.
Mengembalikan token autentikasi sukses atau error kredensial.
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
Mengambil profil hoster.
Mengembalikan data profil jika autentikasi valid atau error jika tidak ditemukan.
*/
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	hoster, err := h.service.GetHosterProfile(userID)
	if err != nil {
		if err.Error() == message.MsgHosterNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, message.MsgInternalServerError)
		}
		return
	}

	response.OK(w, hoster, message.MsgSuccess)
}

/*
Menguji endpoint terproteksi.
Mengembalikan konfirmasi token valid dengan ID user atau error autentikasi.
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
