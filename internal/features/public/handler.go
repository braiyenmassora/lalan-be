package public

import (
	"log"
	"net/http"

	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

/*
PublicHandler menangani permintaan publik.
Menyediakan endpoint tanpa autentikasi.
*/
type PublicHandler struct {
	service PublicService
}

/*
Methods untuk PublicHandler menangani operasi publik kategori, item, dan terms.
Dipanggil dari router untuk akses umum.
*/
func (h *PublicHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
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

func (h *PublicHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
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

func (h *PublicHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
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

/*
NewPublicHandler membuat instance PublicHandler.
Menginisialisasi handler dengan service.
*/
func NewPublicHandler(s PublicService) *PublicHandler {
	return &PublicHandler{service: s}
}
