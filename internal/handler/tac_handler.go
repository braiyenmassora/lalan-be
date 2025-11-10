package handler

import (
	"encoding/json"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
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
