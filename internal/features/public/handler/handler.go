package handler

import (
	"lalan-be/internal/features/public/service"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
	"net/http"
)

type PublicHandler struct {
	service service.PublicService
}

func (h *PublicHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	categories, err := h.service.GetListCategory()
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, categories, message.MsgSuccess)
}

func NewPublicHandler(s service.PublicService) *PublicHandler {
	return &PublicHandler{service: s}
}
