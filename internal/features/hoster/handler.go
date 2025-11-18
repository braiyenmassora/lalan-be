package hoster

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"

	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

/*
HosterHandler menangani permintaan terkait hoster.
Menyediakan layanan untuk operasi hoster melalui service.
*/
type HosterHandler struct {
	service HosterService
}

/*
HosterRequest berisi data untuk membuat hoster baru.
Digunakan dalam permintaan pembuatan hoster.
*/
type HosterRequest struct {
	FullName     string `json:"full_name"`
	StoreName    string `json:"store_name"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	ProfilePhoto string `json:"profile_photo"`
	Description  string `json:"description"`
	Tiktok       string `json:"tiktok"`
	Instagram    string `json:"instagram"`
	Website      string `json:"website"`
}

/*
LoginRequest berisi kredensial untuk autentikasi hoster.
Digunakan dalam permintaan login hoster.
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateIdentityStatusRequest struct {
	Status         string `json:"status"`          // "approved" or "rejected"
	RejectedReason string `json:"rejected_reason"` // Optional
}

/*
Methods untuk HosterHandler menangani operasi hoster, item, dan syarat ketentuan.
Dipanggil dari router untuk memproses permintaan HTTP.
*/
func (h *HosterHandler) CreateHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateHoster: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req HosterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateAdmin: full name required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateAdmin: email required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateAdmin: password required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	input := &model.HosterModel{
		FullName:     req.FullName,
		StoreName:    req.StoreName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		PasswordHash: req.Password,
		Address:      req.Address,
		ProfilePhoto: req.ProfilePhoto,
		Description:  req.Description,
		Tiktok:       req.Tiktok,
		Instagram:    req.Instagram,
		Website:      req.Website,
	}
	err := h.service.CreateHoster(input)
	if err != nil {
		log.Printf("CreateHoster: error creating hoster: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, input, message.MsgSuccess)
}

func (h *HosterHandler) LoginHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginHoster: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginHoster: email or password empty")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginHoster: invalid email format: %s", req.Email)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	resp, err := h.service.LoginHoster(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginHoster: login failed: %v", err)
		response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		return
	}
	log.Printf("LoginHoster: login successful for email %s", req.Email)
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

func (h *HosterHandler) GetDetailHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailHoster: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	ctx := r.Context()
	hoster, err := h.service.GetDetailHoster(ctx)
	if err != nil {
		log.Printf("GetDetailHoster: error getting hoster: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}
	if hoster == nil {
		log.Printf("GetDetailHoster: hoster not found")
		response.Error(w, http.StatusNotFound, message.MsgHosterNotFound)
		return
	}
	log.Printf("GetDetailHoster: retrieved hoster for ID %s", hoster.ID)
	response.OK(w, hoster, message.MsgSuccess)
}

func (h *HosterHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateItem: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req model.ItemModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateItem: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	item, err := h.service.CreateItem(ctx, &req)
	if err != nil {
		log.Printf("CreateItem: error creating item: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, item, message.MsgItemCreatedSuccess)
}

func (h *HosterHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetItemByID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	item, err := h.service.GetItemByID(id)
	if err != nil {
		log.Printf("GetItemByID: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}
	response.OK(w, item, message.MsgSuccess)
}

func (h *HosterHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllItems: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	items, err := h.service.GetAllItems()
	if err != nil {
		log.Printf("GetAllItems: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}
	response.OK(w, items, message.MsgSuccess)
}

func (h *HosterHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateItem: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	var req model.ItemModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateItem: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	item, err := h.service.UpdateItem(ctx, id, &req)
	if err != nil {
		log.Printf("UpdateItem: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, item, message.MsgSuccess)
}

func (h *HosterHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteItem: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	err := h.service.DeleteItem(ctx, id)
	if err != nil {
		log.Printf("DeleteItem: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, nil, message.MsgSuccess)
}

func (h *HosterHandler) CreateTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateTermsAndConditions: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	var req model.TermsAndConditionsModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateTermsAndConditions: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	tac, err := h.service.CreateTermsAndConditions(ctx, &req)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, tac, message.MsgTnCCreatedSuccess)
}

func (h *HosterHandler) FindTermsAndConditionsByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("FindTermsAndConditionsByID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	tac, err := h.service.FindTermsAndConditionsByID(id)
	if err != nil {
		log.Printf("FindTermsAndConditionsByID: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}
	response.OK(w, tac, message.MsgSuccess)
}

func (h *HosterHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllTermsAndConditions: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	tacs, err := h.service.GetAllTermsAndConditions()
	if err != nil {
		log.Printf("GetAllTermsAndConditions: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		return
	}
	response.OK(w, tacs, message.MsgSuccess)
}

func (h *HosterHandler) UpdateTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateTermsAndConditions: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	var req model.TermsAndConditionsModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateTermsAndConditions: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	tac, err := h.service.UpdateTermsAndConditions(ctx, id, &req)
	if err != nil {
		log.Printf("UpdateTermsAndConditions: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, tac, message.MsgSuccess)
}

func (h *HosterHandler) DeleteTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteTermsAndConditions: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	err := h.service.DeleteTermsAndConditions(ctx, id)
	if err != nil {
		log.Printf("DeleteTermsAndConditions: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, nil, message.MsgSuccess)
}

func (h *HosterHandler) GetIdentityCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetIdentityCustomer: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Ambil userID dari path param
	vars := mux.Vars(r)
	log.Printf("GetIdentityCustomer: vars: %+v", vars) // Tambahkan log
	userID := strings.TrimSpace(vars["userID"])
	log.Printf("GetIdentityCustomer: userID from path: %s", userID)
	if userID == "" {
		log.Printf("GetIdentityCustomer: userID required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Panggil service
	ctx := r.Context()
	identity, err := h.service.GetIdentityCustomer(ctx, userID)
	if err != nil {
		log.Printf("GetIdentityCustomer: error getting identity: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else if err.Error() == "Identity not found" {
			response.Error(w, http.StatusNotFound, "Identity not found")
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("GetIdentityCustomer: retrieved identity for user %s", userID)
	response.OK(w, identity, message.MsgSuccess)
}

func (h *HosterHandler) UpdateIdentityStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateIdentityStatus: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgMethodNotAllowed)
		return
	}

	// Ambil identityID dari path param
	vars := mux.Vars(r)
	identityID := strings.TrimSpace(vars["identityID"])
	if identityID == "" {
		log.Printf("UpdateIdentityStatus: identityID required")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	var req UpdateIdentityStatusRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateIdentityStatus: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if req.Status != "approved" && req.Status != "rejected" {
		log.Printf("UpdateIdentityStatus: invalid status")
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Panggil service
	ctx := r.Context()
	err := h.service.UpdateIdentityStatus(ctx, identityID, req.Status, req.RejectedReason)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error updating identity: %v", err)
		if err.Error() == message.MsgUnauthorized {
			response.Error(w, http.StatusUnauthorized, message.MsgUnauthorized)
		} else if err.Error() == "Identity not found" {
			response.Error(w, http.StatusNotFound, "Identity not found")
		} else {
			response.Error(w, http.StatusInternalServerError, message.MsgInternalServerError)
		}
		return
	}

	log.Printf("UpdateIdentityStatus: identity status updated successfully")
	response.OK(w, nil, message.MsgSuccess)
}

/*
NewHosterHandler membuat instance baru HosterHandler.
Menginisialisasi handler dengan service yang diberikan.
*/
func NewHosterHandler(s HosterService) *HosterHandler {
	return &HosterHandler{service: s}
}
