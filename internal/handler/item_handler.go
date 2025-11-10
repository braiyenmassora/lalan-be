package handler

import (
	"encoding/json"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"net/http"
	"strings"
)

/*
Menangani operasi item.
Menyediakan endpoint untuk menambah, mengambil, update, dan hapus item dengan autentikasi dan respons sukses atau error.
*/
type ItemHandler struct {
	service service.ItemService
}

/*
Merepresentasikan data request item.
Digunakan untuk decoding JSON dan validasi input.
*/
type ItemRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	Stock       int      `json:"stock"`
	PickupType  string   `json:"pickup_type"`
	PricePerDay int      `json:"price_per_day"`
	Deposit     int      `json:"deposit"`
	Discount    int      `json:"discount"`
	CategoryID  string   `json:"category_id"`
}

/*
Membuat handler item.
Mengembalikan instance ItemHandler yang siap digunakan.
*/
func NewItemHandler(s service.ItemService) *ItemHandler {
	return &ItemHandler{service: s}
}

/*
Menambahkan item baru.
Mengembalikan respons pembuatan sukses atau error validasi/autentikasi.
*/
func (h *ItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil user ID dari JWT token
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	var req ItemRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Name) == "" {
		response.BadRequest(w, message.MsgItemNameRequired)
		return
	}
	if len(req.Name) > 255 {
		response.BadRequest(w, message.MsgItemNameTooLong)
		return
	}
	if req.CategoryID == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
		return
	}
	if req.PickupType != "pickup" && req.PickupType != "delivery" {
		response.BadRequest(w, "Pickup type must be 'pickup' or 'delivery'")
		return
	}

	input := &model.ItemModel{
		Name:        req.Name,
		Description: req.Description,
		Photos:      req.Photos,
		Stock:       req.Stock,
		PickupType:  model.PickupMethod(req.PickupType),
		PricePerDay: req.PricePerDay,
		Deposit:     req.Deposit,
		Discount:    req.Discount,
		CategoryID:  req.CategoryID,
	}

	itemResp, err := h.service.AddItem(userID, input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, itemResp, message.MsgItemCreatedSuccess)
}

/*
Mengambil semua item.
Mengembalikan daftar item sukses atau error.
*/
func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	items, err := h.service.GetAllItems()
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, items, message.MsgSuccess)
}

/*
Mengambil item berdasarkan ID.
Mengembalikan data item sukses atau error jika tidak ditemukan.
*/
func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgItemIDRequired)
		return
	}

	item, err := h.service.GetItemByID(id)
	if err != nil {
		if err.Error() == message.MsgItemNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, item, message.MsgSuccess)
}

/*
Mengambil item milik user.
Mengembalikan daftar item user sukses atau error autentikasi.
*/
func (h *ItemHandler) GetMyItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil user ID dari JWT token
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	items, err := h.service.GetItemsByUserID(userID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, items, message.MsgSuccess)
}

/*
Mengupdate item.
Mengembalikan respons update sukses atau error validasi/not found/autentikasi.
*/
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil user ID dari JWT token
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgItemIDRequired)
		return
	}

	var req ItemRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Name) == "" {
		response.BadRequest(w, message.MsgItemNameRequired)
		return
	}
	if len(req.Name) > 255 {
		response.BadRequest(w, message.MsgItemNameTooLong)
		return
	}
	if req.CategoryID == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
		return
	}
	if req.PickupType != "pickup" && req.PickupType != "delivery" {
		response.BadRequest(w, "Pickup type must be 'pickup' or 'delivery'")
		return
	}

	input := &model.ItemModel{
		Name:        req.Name,
		Description: req.Description,
		Photos:      req.Photos,
		Stock:       req.Stock,
		PickupType:  model.PickupMethod(req.PickupType),
		PricePerDay: req.PricePerDay,
		Deposit:     req.Deposit,
		Discount:    req.Discount,
		CategoryID:  req.CategoryID,
	}

	itemResp, err := h.service.UpdateItem(id, userID, input)
	if err != nil {
		if err.Error() == message.MsgItemNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), message.MsgNotAllowed) {
			response.Unauthorized(w, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, itemResp, message.MsgItemUpdatedSuccess)
}

/*
Menghapus item.
Mengembalikan respons penghapusan sukses atau error jika tidak ditemukan/autentikasi.
*/
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil user ID dari JWT token
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.MsgUnauthorized)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgItemIDRequired)
		return
	}

	err := h.service.DeleteItem(id, userID)
	if err != nil {
		if err.Error() == message.MsgItemNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), message.MsgNotAllowed) {
			response.Unauthorized(w, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, nil, message.MsgItemDeletedSuccess)
}
