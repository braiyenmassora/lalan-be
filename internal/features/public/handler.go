package public

import (
	"log"
	"net/http"

	"lalan-be/internal/message"
	"lalan-be/internal/response"
)

/*
PublicHandler
menangani permintaan publik tanpa autentikasi
*/
type PublicHandler struct {
	service PublicService
}

/*
GetAllCategories
mengambil semua kategori publik
*/
func (h *PublicHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetCategories: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		response.BadRequest(w, message.BadRequest)
		return
	}

	response.OK(w, categories, message.Success)
}

/*
GetAllItems
mengambil semua item publik
*/
func (h *PublicHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
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
GetAllTermsAndConditions
mengambil semua syarat dan ketentuan publik
*/
func (h *PublicHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
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
NewPublicHandler
membuat instance PublicHandler dengan service
*/
func NewPublicHandler(s PublicService) *PublicHandler {
	return &PublicHandler{service: s}
}
