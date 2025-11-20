package admin

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
type AdminHandler struct
struct yang menangani permintaan terkait admin dan kategori menggunakan service
*/
type AdminHandler struct {
	service AdminService
}

/*
CreateAdmin
membuat admin baru dengan validasi input dan response sukses atau error
*/
func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateAdmin: received request")
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}
	var req AdminRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateAdmin: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi full name
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateAdmin: full name required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "full name"))
		return
	}

	// Validasi email
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateAdmin: email required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "email"))
		return
	}

	// Validasi password
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateAdmin: password required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "password"))
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
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, input, fmt.Sprintf(message.Created, "admin"))
}

/*
LoginAdmin
melakukan login admin dengan validasi kredensial dan response token atau error
*/
func (h *AdminHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginAdmin: received request")
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginAdmin: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi email
	if req.Email == "" {
		log.Printf("LoginAdmin: email required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "email"))
		return
	}

	// Validasi password
	if req.Password == "" {
		log.Printf("LoginAdmin: password required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "password"))
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// Validasi format email
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginAdmin: invalid email format: %s", req.Email)
		response.BadRequest(w, fmt.Sprintf(message.InvalidFormat, "email"))
		return
	}

	resp, err := h.service.LoginAdmin(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginAdmin: login failed: %v", err)
		response.Unauthorized(w, message.LoginFailed)
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

	response.OK(w, userData, message.Success)
}

/*
CreateCategory
membuat kategori baru dengan validasi input dan response sukses atau error
*/
func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateCategory: received request")
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi name
	if strings.TrimSpace(req.Name) == "" {
		log.Printf("CreateCategory: name required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "name"))
		return
	}

	// Validasi panjang name
	if len(req.Name) > 255 {
		log.Printf("CreateCategory: name too long")
		response.BadRequest(w, fmt.Sprintf(message.TooLong, "name"))
		return
	}

	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.CreateCategory(input)
	if err != nil {
		log.Printf("CreateCategory: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, input, fmt.Sprintf(message.Created, "category"))
}

/*
UpdateCategory
memperbarui kategori dengan validasi input dan response sukses atau error
*/
func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateCategory: received request")
	if r.Method != http.MethodPut {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	// Validasi ID
	if strings.TrimSpace(id) == "" {
		log.Printf("UpdateCategory: id required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "id"))
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi name
	if strings.TrimSpace(req.Name) == "" {
		log.Printf("UpdateCategory: name required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "name"))
		return
	}

	// Validasi panjang name
	if len(req.Name) > 255 {
		log.Printf("UpdateCategory: name too long")
		response.BadRequest(w, fmt.Sprintf(message.TooLong, "name"))
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
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, input, fmt.Sprintf(message.Updated, "category"))
}

/*
DeleteCategory
menghapus kategori berdasarkan ID dan response sukses atau error
*/
func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteCategory: received request")
	if r.Method != http.MethodDelete {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	// Validasi ID
	if strings.TrimSpace(id) == "" {
		log.Printf("DeleteCategory: id required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "id"))
		return
	}

	err := h.service.DeleteCategory(id)
	if err != nil {
		log.Printf("DeleteCategory: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, nil, fmt.Sprintf(message.Deleted, "category"))
}

/*
GetAllCategory
mengambil semua kategori dan response data atau error
*/
func (h *AdminHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetCategories: received request")
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("GetCategories: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, categories, message.Success)
}

/*
UpdateIdentityStatus
memperbarui status identitas berdasarkan user ID untuk approval admin
*/
func (h *AdminHandler) UpdateIdentityStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateIdentityStatus: received request")
	if r.Method != http.MethodPut {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	// Validasi user ID
	if strings.TrimSpace(userID) == "" {
		log.Printf("UpdateIdentityStatus: user_id required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "user_id"))
		return
	}

	var req UpdateIdentityRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateIdentityStatus: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi status
	if strings.TrimSpace(req.Status) == "" {
		log.Printf("UpdateIdentityStatus: status required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "status"))
		return
	}

	if req.Status != "approved" && req.Status != "rejected" {
		log.Printf("UpdateIdentityStatus: invalid status %s", req.Status)
		response.BadRequest(w, message.InvalidStatus)
		return
	}

	err := h.service.UpdateIdentityStatus(r.Context(), userID, req.Status, req.Reason)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error: %v", err)
		if err.Error() == message.Unauthorized {
			response.Unauthorized(w, message.Unauthorized)
			return
		}
		if err.Error() == message.InvalidStatus {
			response.BadRequest(w, message.InvalidStatus)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(w, fmt.Sprintf(message.NotFound, "identity"))
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	if req.Status == "approved" {
		response.OK(w, nil, message.IdentityApproved)
	} else {
		response.OK(w, nil, fmt.Sprintf(message.IdentityRejected, req.Reason))
	}
}

/*
GetIdentityByCustomerID
mengambil data identitas berdasarkan user ID
*/
func (h *AdminHandler) GetIdentityByCustomerID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetIdentityByCustomerID: received request")
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, message.MethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	// Validasi user ID
	if strings.TrimSpace(userID) == "" {
		log.Printf("GetIdentityByCustomerID: user_id required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "user_id"))
		return
	}

	identity, err := h.service.GetIdentityByCustomerID(userID)
	if err != nil {
		log.Printf("GetIdentityByCustomerID: error: %v", err)
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(w, fmt.Sprintf(message.NotFound, "identity"))
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, identity, message.Success)
}

/*
type AdminRequest struct
struct yang berisi data untuk membuat admin baru
*/
type AdminRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
type LoginRequest struct
struct yang berisi kredensial untuk autentikasi admin
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
type CategoryRequest struct
struct yang berisi data untuk operasi kategori
*/
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
type UpdateIdentityRequest struct
struct yang berisi data untuk update status identitas
*/
type UpdateIdentityRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

/*
NewAdminHandler
membuat instance baru AdminHandler dengan service yang diberikan
*/
func NewAdminHandler(s AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}
