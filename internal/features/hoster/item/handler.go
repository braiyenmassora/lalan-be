package item

import (
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lalan-be/internal/domain"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"
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
CreateItem menangani POST /api/v1/hoster/items dengan support upload photos

Alur kerja:
1. Validasi method POST
2. Parse multipart form
3. Ambil hosterID dari JWT context
4. Parse form fields dan files
5. Build domain.Item
6. Panggil service.CreateItem dengan ctx, item, dan photoFiles
7. Kembalikan hasil atau error
*/
func (h *HosterItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateItem: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // Max 10MB total
	if err != nil {
		log.Printf("CreateItem: failed to parse multipart: %v", err)
		response.BadRequest(w, "Invalid form data")
		return
	}

	// Ambil hosterID dari middleware
	hosterID := middleware.GetUserID(r)
	if hosterID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	// Parse form fields
	name := r.FormValue("name")
	description := r.FormValue("description")
	stockStr := r.FormValue("stock")
	pricePerDayStr := r.FormValue("price_per_day")
	depositStr := r.FormValue("deposit")
	pickupType := r.FormValue("pickup_type")
	categoryID := r.FormValue("category_id")
	discountStr := r.FormValue("discount")

	// Convert dan validasi
	stock, err := strconv.Atoi(stockStr)
	if err != nil || stock < 0 {
		response.BadRequest(w, "Invalid stock")
		return
	}
	pricePerDay, err := strconv.Atoi(pricePerDayStr)
	if err != nil || pricePerDay <= 0 {
		response.BadRequest(w, "Invalid price_per_day")
		return
	}
	deposit, err := strconv.Atoi(depositStr)
	if err != nil || deposit < 0 {
		response.BadRequest(w, "Invalid deposit")
		return
	}
	discount := 0
	if discountStr != "" {
		discount, err = strconv.Atoi(discountStr)
		if err != nil || discount < 0 {
			response.BadRequest(w, "Invalid discount")
			return
		}
	}

	// Parse photo files
	files := r.MultipartForm.File["photos"]
	var photoFiles []*multipart.FileHeader
	for _, fileHeader := range files {
		photoFiles = append(photoFiles, fileHeader)
	}

	// Build item
	item := &domain.Item{
		ID:          uuid.New().String(),
		HosterID:    hosterID,
		Name:        name,
		Description: description,
		Stock:       stock,
		PickupType:  domain.PickupMethod(pickupType),
		PricePerDay: pricePerDay,
		Deposit:     deposit,
		Discount:    discount,
		CategoryID:  categoryID,
	}

	// Panggil service
	detail, err := h.service.CreateItem(r.Context(), item, photoFiles)
	if err != nil {
		log.Printf("CreateItem: service error: %v", err)
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

	response.OK(w, detail, message.ItemCreated)
}

/*
DeleteItem menangani DELETE /api/v1/hoster/item/{id}
- client hanya perlu mengirim path param "id" (item id)
- hoster id diambil dari JWT/session (middleware) => tidak perlu dikirim di body
Responses:
  - 401 Unauthorized -> jika tidak authenticated
  - 400 Bad Request -> jika input tidak valid atau item tidak ada / bukan milik hoster
  - 204 No Content -> jika sukses
  - 500 Internal Server Error -> jika terjadi error internal
*/
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
	response.OK(w, nil, message.ItemDeleted)
}
