package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"lalan-be/internal/model"
	"lalan-be/internal/service"
	"lalan-be/pkg"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// helper untuk response error konsisten
func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(pkg.Response{
		Code:    code,
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}

// RegisterRequest adalah struct untuk menerima data register dari client
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

// Register: POST /v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	// digunakan untuk mengunci agar user tidak menambhakan atau mengurangi payload
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Convert ke model
	input := &model.HosterModel{
		FullName:     req.FullName,
		ProfilePhoto: req.ProfilePhoto,
		StoreName:    req.StoreName,
		Description:  req.Description,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Address:      req.Address,
		PasswordHash: req.Password, // akan di-hash di service
	}

	if err := h.service.Register(input); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg.Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Hoster registered successfully",
		Data:    nil,
	})
}

// Login: POST /v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	log.Printf("Login request: email=%s", req.Email)

	resp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
