package hoster

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"lalan-be/internal/config"
	"lalan-be/internal/message"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
)

/*
type HosterHandler struct
menangani permintaan terkait hoster melalui service
*/
type HosterHandler struct {
	service HosterService
}

/*
type HosterRequest struct
berisi data untuk membuat hoster baru
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
type LoginRequest struct
berisi kredensial untuk autentikasi hoster
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
CreateHoster
membuat hoster baru dengan validasi input dan response sukses atau error
*/
func (h *HosterHandler) CreateHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateHoster: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req HosterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateHoster: full name required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "full name"))
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateHoster: email required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "email"))
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateHoster: password required")
		response.BadRequest(w, fmt.Sprintf(message.Required, "password"))
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
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, input, fmt.Sprintf(message.Created, "hoster"))
}

/*
LoginHoster
melakukan login hoster dengan validasi kredensial dan response token atau error
*/
// LoginHoster — versi DEV (refresh_token masih keliatan di body)
func (h *HosterHandler) LoginHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginHoster: received request")

	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	// === BAGIAN INI WAJIB ADA ===
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Validasi email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		response.BadRequest(w, fmt.Sprintf(message.InvalidFormat, "email"))
		return
	}
	// === SAMPE SINI ===

	// Panggil service
	resp, err := h.service.LoginHoster(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginHoster: login failed: %v", err)
		response.Unauthorized(w, message.LoginFailed)
		return
	}

	accessToken := resp.AccessToken
	refreshToken := resp.RefreshToken

	// Simpan ke Redis
	config.Redis.Set(config.RedisCtx, "refresh:"+refreshToken, resp.ID+":"+resp.Role, 30*24*time.Hour)

	// httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // dev
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth/refresh",
	})

	// Hapus cookie lama
	http.SetCookie(w, &http.Cookie{
		Name: "auth_token", Value: "", MaxAge: -1, Path: "/",
	})

	log.Printf("LoginHoster: success → %s (id: %s)", req.Email, resp.ID)

	// DEV MODE: tampilkan refresh_token di body
	response.OK(w, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken, // ← buat testing dev
		"expires_in":    resp.ExpiresIn,
		"id":            resp.ID,
		"role":          resp.Role,
	}, "Login successful (DEV MODE)")
}

/*
GetDetailHoster
mengambil detail hoster berdasarkan konteks dan response data atau error
*/
func (h *HosterHandler) GetDetailHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailHoster: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	ctx := r.Context()
	hoster, err := h.service.GetDetailHoster(ctx)
	if err != nil {
		log.Printf("GetDetailHoster: error getting hoster: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	if hoster == nil {
		log.Printf("GetDetailHoster: hoster not found")
		response.Error(w, http.StatusNotFound, fmt.Sprintf(message.NotFound, "hoster"))
		return
	}
	log.Printf("GetDetailHoster: retrieved hoster for ID %s", hoster.ID)
	response.OK(w, hoster, message.Success)
}

/*
CreateItem
membuat item baru dengan validasi input dan response sukses atau error
*/
func (h *HosterHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateItem: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req model.ItemModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateItem: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	item, err := h.service.CreateItem(ctx, &req)
	if err != nil {
		log.Printf("CreateItem: error creating item: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, item, fmt.Sprintf(message.Created, "item"))
}

/*
GetItemByID
mengambil item berdasarkan ID dan response data atau error
*/
func (h *HosterHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetItemByID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}
	item, err := h.service.GetItemByID(id)
	if err != nil {
		log.Printf("GetItemByID: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, item, message.Success)
}

/*
GetAllItems
mengambil semua item dan response data atau error
*/
func (h *HosterHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllItems: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	items, err := h.service.GetAllItems()
	if err != nil {
		log.Printf("GetAllItems: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, items, message.Success)
}

/*
UpdateItem
memperbarui item berdasarkan ID dan response sukses atau error
*/
func (h *HosterHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateItem: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}
	var req model.ItemModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateItem: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	item, err := h.service.UpdateItem(ctx, id, &req)
	if err != nil {
		log.Printf("UpdateItem: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, item, fmt.Sprintf(message.Updated, "item"))
}

/*
DeleteItem
menghapus item berdasarkan ID dan response sukses atau error
*/
func (h *HosterHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteItem: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	err := h.service.DeleteItem(ctx, id)
	if err != nil {
		log.Printf("DeleteItem: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, nil, fmt.Sprintf(message.Deleted, "item"))
}

/*
CreateTermsAndConditions
membuat syarat dan ketentuan baru dan response sukses atau error
*/
func (h *HosterHandler) CreateTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateTermsAndConditions: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	var req model.TermsAndConditionsModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateTermsAndConditions: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	tac, err := h.service.CreateTermsAndConditions(ctx, &req)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, tac, fmt.Sprintf(message.Created, "terms and conditions"))
}

/*
FindTermsAndConditionsByID
mencari syarat dan ketentuan berdasarkan ID dan response data atau error
*/
func (h *HosterHandler) FindTermsAndConditionsByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("FindTermsAndConditionsByID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}
	tac, err := h.service.FindTermsAndConditionsByID(id)
	if err != nil {
		log.Printf("FindTermsAndConditionsByID: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, tac, message.Success)
}

/*
GetAllTermsAndConditions
mengambil semua syarat dan ketentuan dan response data atau error
*/
func (h *HosterHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllTermsAndConditions: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	tacs, err := h.service.GetAllTermsAndConditions()
	if err != nil {
		log.Printf("GetAllTermsAndConditions: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, tacs, message.Success)
}

/*
UpdateTermsAndConditions
memperbarui syarat dan ketentuan berdasarkan ID dan response sukses atau error
*/
func (h *HosterHandler) UpdateTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateTermsAndConditions: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}
	var req model.TermsAndConditionsModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateTermsAndConditions: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	tac, err := h.service.UpdateTermsAndConditions(ctx, id, &req)
	if err != nil {
		log.Printf("UpdateTermsAndConditions: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, tac, fmt.Sprintf(message.Updated, "terms and conditions"))
}

/*
DeleteTermsAndConditions
menghapus syarat dan ketentuan berdasarkan ID dan response sukses atau error
*/
func (h *HosterHandler) DeleteTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteTermsAndConditions: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}
	ctx := r.Context()
	err := h.service.DeleteTermsAndConditions(ctx, id)
	if err != nil {
		log.Printf("DeleteTermsAndConditions: error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, nil, fmt.Sprintf(message.Deleted, "terms and conditions"))
}

/*
GetListBookingsForHoster
mengambil daftar booking yang dimiliki oleh hoster berdasarkan konteks dengan pagination dan response data atau error
*/
func (h *HosterHandler) GetListBookingsForHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookingsForHoster: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	ctx := r.Context()
	bookings, err := h.service.GetListBookingsCustomer(ctx, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsForHoster: error getting bookings: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	log.Printf("GetListBookingsForHoster: retrieved %d bookings for page %d", len(bookings), page)
	response.OK(w, bookings, message.Success)
}

/*
GetListBookingsForHosterByCustomerID
mengambil daftar booking yang dimiliki oleh hoster berdasarkan konteks dan customer ID dengan pagination dan response data atau error
*/
func (h *HosterHandler) GetListBookingsForHosterByCustomerID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookingsForHosterByCustomerID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	customerID := vars["customerID"]
	if customerID == "" {
		response.BadRequest(w, "customer ID is required")
		return
	}

	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	ctx := r.Context()
	bookings, err := h.service.GetListBookingsCustomerByBookingID(ctx, customerID, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsForHosterByCustomerID: error getting bookings: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	log.Printf("GetListBookingsForHosterByCustomerID: retrieved %d bookings for customer %s", len(bookings), customerID)
	response.OK(w, bookings, message.Success)
}

/*
GetListBookingsCustomer
mengambil daftar booking yang dimiliki oleh hoster berdasarkan konteks dengan pagination dan response data atau error
*/
func (h *HosterHandler) GetListBookingsCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookingsCustomer: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	ctx := r.Context()
	bookings, err := h.service.GetListBookingsCustomer(ctx, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsCustomer: error getting bookings: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	log.Printf("GetListBookingsCustomer: retrieved %d bookings for page %d", len(bookings), page)
	response.OK(w, bookings, message.Success)
}

/*
GetListBookingsCustomerByBookingID
mengambil daftar booking yang dimiliki oleh hoster berdasarkan konteks dan booking ID dengan pagination dan response data atau error
*/
func (h *HosterHandler) GetListBookingsCustomerByBookingID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetListBookingsCustomerByBookingID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	bookingID := vars["bookingID"]
	if bookingID == "" {
		response.BadRequest(w, "booking ID is required")
		return
	}

	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	ctx := r.Context()
	bookings, err := h.service.GetListBookingsCustomerByBookingID(ctx, bookingID, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsCustomerByBookingID: error getting bookings: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	log.Printf("GetListBookingsCustomerByBookingID: retrieved %d bookings for booking %s", len(bookings), bookingID)
	response.OK(w, bookings, message.Success)
}

/*
NewHosterHandler
membuat instance baru HosterHandler dengan service yang diberikan
*/
func NewHosterHandler(s HosterService) *HosterHandler {
	return &HosterHandler{service: s}
}
