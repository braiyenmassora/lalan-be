package item

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"lalan-be/internal/domain"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"

	"github.com/google/uuid"
)

/*
HosterItemHandler menangani endpoint HTTP untuk fitur item dari perspektif hoster.
Menyediakan operasi read (list) dan create.
*/
type HosterItemHandler struct {
	service ItemService
}

/*
NewHosterItemHandler membuat instance handler dengan dependency injection.

Output:
- *HosterItemHandler siap digunakan
*/
func NewHosterItemHandler(s ItemService) *HosterItemHandler {
	return &HosterItemHandler{service: s}
}

/*
GetListItem menangani GET /api/v1/hoster/items

Alur kerja:
1. Validasi method GET
2. Ambil hosterID dari JWT context
3. Panggil service untuk ambil daftar item milik hoster

Output sukses:
- 200 OK + list item ringkas
Output error:
- 401 Unauthorized / 500 Internal Server Error
*/
func (h *HosterItemHandler) GetListItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	items, err := h.service.GetListItem(hosterID)
	if err != nil {
		log.Printf("GetListItem handler: service error hoster=%s err=%v", hosterID, err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, items, message.Success)
}

/*
CreateItem menangani POST /api/v1/hoster/items

Alur kerja:
1. Validasi method POST
2. Ambil hosterID dari JWT context
3. Decode body -> domain.Item, set id/user/timestamps
4. Panggil service.CreateItem
5. Kembalikan hasil atau error yang sesuai
*/
func (h *HosterItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MethodNotAllowed)
		return
	}

	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	var req domain.Item
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateItem handler: invalid request body hoster=%s err=%v", hosterID, err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// set required server-side fields
	req.HosterID = hosterID
	req.ID = uuid.New().String()
	now := time.Now().UTC()
	req.CreatedAt = now
	req.UpdatedAt = now

	created, err := h.service.CreateItem(&req)
	if err != nil {
		log.Printf("CreateItem handler: service error hoster=%s err=%v", hosterID, err)
		switch err.Error() {
		case message.BadRequest:
			response.BadRequest(w, message.BadRequest)
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	response.OK(w, created, message.ItemCreated)
}

// DeleteItem menangani DELETE /api/v1/hoster/item/{id}
// - client hanya perlu mengirim path param "id" (item id)
// - hoster id diambil dari JWT/session (middleware) => tidak perlu dikirim di body
// Responses:
//   - 401 Unauthorized -> jika tidak authenticated
//   - 400 Bad Request -> jika input tidak valid atau item tidak ditemukan / bukan milik hoster
//   - 204 No Content -> jika sukses
//   - 500 Internal Server Error -> jika terjadi error internal
func (h *HosterItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	hosterID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	itemID := vars["id"]

	log.Printf("DeleteItem handler: hosterID=%s, itemID=%s", hosterID, itemID)

	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}
	if itemID == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	if err := h.service.DeleteItem(hosterID, itemID); err != nil {
		switch err.Error() {
		case message.BadRequest:
			response.BadRequest(w, message.BadRequest)
		case message.Unauthorized:
			response.Unauthorized(w, message.Unauthorized)
		case message.ItemNotFound:
			response.NotFound(w, message.ItemNotFound)
		default:
			response.Error(w, http.StatusInternalServerError, message.InternalError)
		}
		return
	}

	log.Printf("DeleteItem handler: deleted item %s for hoster %s", itemID, hosterID)
	// return JSON body with code/message/success for clients (Postman)
	response.OK(w, nil, message.ItemDeleted)
}
