package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
)

// HosterHandler menangani autentikasi hoster.
// Menyediakan operasi registrasi, login, dan pengambilan profil dengan respons sukses atau error.
type HosterHandler struct {
	service service.HosterService
}

// RegisterRequest merepresentasikan data registrasi hoster.
// Digunakan untuk decoding request JSON dan mapping ke model.
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

// NewHosterHandler membuat handler autentikasi.
// Mengembalikan instance HosterHandler yang siap digunakan.
func NewHosterHandler(s service.HosterService) *HosterHandler {
	return &HosterHandler{service: s}
}

// Register memproses registrasi hoster.
// Mengembalikan respons pembuatan akun sukses atau error validasi.
func (h *HosterHandler) Register(w http.ResponseWriter, r *http.Request) {
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

	hosterResp, err := h.service.Register(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, hosterResp, message.MsgHosterCreatedSuccess)
}

// LoginHoster memproses login hoster.
// Mengembalikan token autentikasi sukses atau error kredensial.
func (h *HosterHandler) LoginHoster(w http.ResponseWriter, r *http.Request) {
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

	hosterResp, err := h.service.LoginHoster(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}

	response.OK(w, hosterResp, message.MsgSuccess)
}

// GetHosterProfile mengambil profil hoster.
// Mengembalikan data profil jika autentikasi valid atau error jika tidak ditemukan.
func (h *HosterHandler) GetHosterProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	// Buat context dengan userID
	ctx := context.WithValue(context.Background(), "user_id", userID)

	hoster, err := h.service.GetHosterProfile(ctx)
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

// TestProtected menguji endpoint terproteksi.
// Mengembalikan konfirmasi token valid dengan ID user atau error autentikasi.
func (h *HosterHandler) TestProtected(w http.ResponseWriter, r *http.Request) {
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
