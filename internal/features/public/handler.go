package public

import (
	"log"
	"net/http"

	"lalan-be/internal/message"
	"lalan-be/internal/response"
)

/*
PublicHandler menangani semua endpoint publik yang tidak memerlukan autentikasi.
Hanya berfungsi sebagai adapter antara HTTP request dan layer service.
*/
type PublicHandler struct {
	service PublicService
}

/*
GetAllCategories menangani endpoint GET /public/categories.

Alur kerja:
1. Validasi method HTTP
2. Panggil service untuk ambil semua kategori
3. Return data atau error

Output sukses:
- 200 OK + list kategori (DTO)
Output error:
- 405 Method Not Allowed / 400 Bad Request jika service error
*/
func (h *PublicHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllCategories: request from %s", r.RemoteAddr)

	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("GetAllCategories: service error: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	response.OK(w, categories, message.Success)
}

/*
GetAllItems menangani endpoint GET /public/items.

Alur kerja:
1. Validasi method
2. Panggil service untuk ambil semua item publik
3. Return data atau 500 jika service gagal

Output sukses:
- 200 OK + list item (DTO)
Output error:
- 405 / 500 Internal Server Error
*/
func (h *PublicHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllItems: request from %s", r.RemoteAddr)

	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	items, err := h.service.GetAllItems()
	if err != nil {
		log.Printf("GetAllItems: service error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, items, message.Success)
}

/*
GetAllTermsAndConditions menangani endpoint GET /public/terms.

Alur kerja:
1. Validasi method
2. Ambil data syarat & ketentuan dari service
3. Return data atau 500 jika gagal

Output sukses:
- 200 OK + data terms (DTO)
Output error:
- 405 / 500 Internal Server Error
*/
func (h *PublicHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllTermsAndConditions: request from %s", r.RemoteAddr)

	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	tacs, err := h.service.GetAllTermsAndConditions()
	if err != nil {
		log.Printf("GetAllTermsAndConditions: service error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, tacs, message.Success)
}

/*
NewPublicHandler membuat instance PublicHandler dengan dependency injection.

Output:
- *PublicHandler siap digunakan dengan service yang sudah disuntikkan
*/
func NewPublicHandler(s PublicService) *PublicHandler {
	return &PublicHandler{service: s}
}
