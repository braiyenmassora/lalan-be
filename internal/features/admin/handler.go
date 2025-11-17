package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

/*
AdminHandler menangani permintaan terkait admin.
Menggunakan service untuk operasi bisnis admin dan kategori.
*/
type AdminHandler struct {
	service AdminService
}

/*
Methods AdminHandler mengelola admin dan kategori.
Menyediakan endpoint untuk CRUD dan autentikasi.
*/
func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateAdmin: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req AdminRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateAdmin: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi full name
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateAdmin: full name required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi email
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateAdmin: email required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi password
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateAdmin: password required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	input := &model.AdminModel{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	err := h.service.CreateAdmin(input)
	if err != nil {
		log.Printf("CreateAdmin: error creating admin: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, input, message.MsgSuccess)
}

func (h *AdminHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginAdmin: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginAdmin: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi email dan password
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginAdmin: email or password empty")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// Validasi format email
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginAdmin: invalid email format: %s", req.Email)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	resp, err := h.service.LoginAdmin(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginAdmin: login failed: %v", err)
		response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		return
	}

	log.Printf("LoginAdmin: login successful for email %s", req.Email)
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

func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateCategory: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi name
	if strings.TrimSpace(req.Name) == "" {
		log.Printf("CreateCategory: name required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi panjang name
	if len(req.Name) > 255 {
		log.Printf("CreateCategory: name too long")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.CreateCategory(input)
	if err != nil {
		log.Printf("CreateCategory: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, input, message.MsgSuccess)
}

func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateCategory: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	// Validasi ID
	if strings.TrimSpace(id) == "" {
		log.Printf("UpdateCategory: id required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi name
	if strings.TrimSpace(req.Name) == "" {
		log.Printf("UpdateCategory: name required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi panjang name
	if len(req.Name) > 255 {
		log.Printf("UpdateCategory: name too long")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	input := &model.CategoryModel{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.UpdateCategory(input)
	if err != nil {
		log.Printf("UpdateCategory: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, input, message.MsgSuccess)
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteCategory: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	// Validasi ID
	if strings.TrimSpace(id) == "" {
		log.Printf("DeleteCategory: id required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	err := h.service.DeleteCategory(id)
	if err != nil {
		log.Printf("DeleteCategory: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, nil, message.MsgSuccess)
}

func (h *AdminHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetCategories: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	response.OK(w, categories, message.MsgSuccess)
}

/*
AdminRequest berisi data untuk membuat admin baru.
Digunakan dalam endpoint pembuatan admin.
*/
type AdminRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
LoginRequest berisi kredensial untuk autentikasi admin.
Digunakan dalam endpoint login admin.
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
CategoryRequest berisi data untuk operasi kategori.
Digunakan dalam endpoint CRUD kategori.
*/
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
NewAdminHandler membuat instance baru AdminHandler.
Mengembalikan pointer ke AdminHandler dengan service yang diberikan.
*/
func NewAdminHandler(s AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}
