package handler

import (
	"net/http"

	"lalan-be/internal/helper"
	"lalan-be/internal/model"
	"lalan-be/internal/service"

	"github.com/goccy/go-json"
)

type HosterHandler struct {
	service service.HosterService
}

func NewHosterHandler(s service.HosterService) *HosterHandler {
	return &HosterHandler{service: s}
}

func (h *HosterHandler) RegisterHosterHandler(w http.ResponseWriter, r *http.Request) {
	var input model.HosterModel

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.ErrorResponse(w, nil, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterHoster(&input); err != nil {
		helper.ErrorResponse(w, nil, err.Error(), http.StatusBadRequest)
		return
	}

	helper.SuccessResponse(w, nil, "hoster registered successfully", http.StatusCreated)
}
