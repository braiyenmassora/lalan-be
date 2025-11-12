package handler

import (
	"encoding/json"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"log"
	"net/http"
)

/*
Menangani operasi terms and conditions.
Menyediakan endpoint untuk menambah terms and conditions dengan autentikasi dan respons sukses atau error.
*/
type TermsAndConditionsHandler struct {
	service service.TermsAndConditionsService
}

/*
Merepresentasikan data request terms and conditions.
Digunakan untuk decoding JSON dan validasi input.
*/
type TermsAndConditionsRequest struct {
	Description []string `json:"description"`
}

/*
Membuat handler terms and conditions.
Mengembalikan instance TermsAndConditionsHandler yang siap digunakan.
*/
func NewTermsAndConditionsHandler(s service.TermsAndConditionsService) *TermsAndConditionsHandler {
	return &TermsAndConditionsHandler{service: s}
}

/*
Menambahkan terms and conditions.
Mengembalikan respons pembuatan sukses atau error validasi/autentikasi.
*/
func (h *TermsAndConditionsHandler) AddTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Validasi role: hanya hoster yang bisa
	userRole := middleware.GetUserRole(r)
	if userRole != "hoster" {
		response.Error(w, http.StatusForbidden, "Access denied: only hosters can add terms and conditions")
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	var req TermsAndConditionsRequest
	decoder := json.NewDecoder(r.Body)
	// Hapus DisallowUnknownFields agar tidak error jika ada field tambahan seperti user_id
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if len(req.Description) == 0 {
		response.BadRequest(w, message.MsgTermAndConditionsDescriptionRequired)
		return
	}

	input := &model.TermsAndConditionsModel{
		Description: req.Description,
		UserID:      userID,
	}

	tncResp, err := h.service.AddTermsAndConditions(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, tncResp, message.MsgTermAndConditionsCreatedSuccess)
}

// Tambahkan method lain jika diperlukan untuk CRUD lengkap, dengan validasi role serupa
// Contoh: GetTermsAndConditions, UpdateTermsAndConditions, DeleteTermsAndConditions
func (h *TermsAndConditionsHandler) GetAllTermAndConditions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	categories, err := h.service.GetAllTermAndConditions()
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, categories, message.MsgSuccess)
}

func (h *TermsAndConditionsHandler) UpdateTermAndConditions(w http.ResponseWriter, r *http.Request) {
	// Tambahkan log untuk debug
	log.Printf("UpdateTermAndConditions called for userID: %s", middleware.GetUserID(r))

	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Validasi role: hanya hoster yang bisa
	userRole := middleware.GetUserRole(r)
	if userRole != "hoster" {
		response.Forbidden(w, "Access denied: only hosters can update terms and conditions")
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	var req TermsAndConditionsRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	if len(req.Description) == 0 {
		response.BadRequest(w, message.MsgTermAndConditionsDescriptionRequired)
		return
	}

	input := &model.TermsAndConditionsModel{
		Description: req.Description,
	}

	tncResp, err := h.service.UpdateTermsAndConditionsByUserID(userID, input)
	if err != nil {
		log.Printf("Error updating TAC: %v", err) // Tambahkan log error
		if err.Error() == message.MsgTermAndConditionsNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, tncResp, message.MsgTermAndConditionsUpdatedSuccess)
}
