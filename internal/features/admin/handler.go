package admin

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

/*
Menangani operasi admin.
Menyediakan endpoint untuk menambah, login, dan operasi kategori dengan respons sukses atau error.
*/
type AdminHandler struct {
	service AdminService
}

/*
Merepresentasikan data request admin.
Digunakan untuk decoding JSON dan validasi input.
*/
type AdminRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
Merepresentasikan data request login admin.
Digunakan untuk decoding JSON dan validasi input.
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
Merepresentasikan data request kategori.
Digunakan untuk decoding JSON dan validasi input.
*/
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
Menambahkan admin baru.
Mengembalikan respons pembuatan sukses atau error validasi.
*/
func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Removed role check, handled in middleware/route

	var req AdminRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.FullName) == "" {
		response.BadRequest(w, "Full name is required")
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		response.BadRequest(w, "Email is required")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		response.BadRequest(w, "Password is required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.BadRequest(w, "Failed to hash password")
		return
	}

	input := &model.AdminModel{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = h.service.CreateAdmin(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, input, "Admin created successfully")
}

/*
Login admin.
Mengembalikan respons login sukses atau error.
*/
func (h *AdminHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Email) == "" {
		response.BadRequest(w, "Email is required")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		response.BadRequest(w, "Password is required")
		return
	}

	admin, err := h.service.FindByEmailAdminForLogin(req.Email)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	if admin == nil {
		response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password))
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response.OK(w, admin, "Login successful")
}

/*
Menambahkan kategori baru.
Mengembalikan respons pembuatan sukses atau error validasi.
*/
func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Removed role check, handled in middleware/route

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Name) == "" {
		response.BadRequest(w, message.MsgCategoryNameRequired)
		return
	}

	if len(req.Name) > 255 {
		response.BadRequest(w, message.MsgCategoryNameTooLong)
		return
	}

	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.CreateCategory(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, input, message.MsgCategoryCreatedSuccess)
}

/*
Mengupdate kategori.
Mengembalikan respons update sukses atau error validasi/not found.
*/
func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Removed role check, handled in middleware/route

	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Name) == "" {
		response.BadRequest(w, message.MsgCategoryNameRequired)
		return
	}

	if len(req.Name) > 255 {
		response.BadRequest(w, message.MsgCategoryNameTooLong)
		return
	}

	input := &model.CategoryModel{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.UpdateCategory(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, input, message.MsgCategoryUpdatedSuccess)
}

/*
Menghapus kategori.
Mengembalikan respons penghapusan sukses atau error jika tidak ditemukan.
*/
func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Removed role check, handled in middleware/route

	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
		return
	}

	err := h.service.DeleteCategory(id)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, nil, message.MsgCategoryDeletedSuccess)
}

/*
Membuat handler admin.
Mengembalikan instance AdminHandler yang siap digunakan.
*/
func NewAdminHandler(s AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}
