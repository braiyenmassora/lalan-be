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

type ItemHandler struct {
	service service.ItemService
}

func NewItemHandler(s service.ItemService) *ItemHandler {
	return &ItemHandler{service: s}
}

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

// AddItem menambahkan item baru
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
		response.BadRequest(w, "Item name is required")
		return
	}
	if len(req.Name) > 255 {
		response.BadRequest(w, "Item name must not exceed 255 characters")
		return
	}
	if req.CategoryID == "" {
		response.BadRequest(w, "Category ID is required")
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
		UserID:      userID,
	}

	itemResp, err := h.service.AddItem(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, itemResp, message.MsgItemCreatedSuccess)
}

// GetAllItems mendapatkan semua item
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

// GetItemByID mendapatkan item berdasarkan ID
func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "Item ID is required")
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

// GetMyItems mendapatkan semua item milik user yang login
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

// UpdateItem mengupdate item
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
		response.BadRequest(w, "Item ID is required")
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
		response.BadRequest(w, "Item name is required")
		return
	}
	if len(req.Name) > 255 {
		response.BadRequest(w, "Item name must not exceed 255 characters")
		return
	}
	if req.CategoryID == "" {
		response.BadRequest(w, "Category ID is required")
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
		} else if strings.Contains(err.Error(), "unauthorized") {
			response.Unauthorized(w, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, itemResp, "Item updated successfully")
}

// DeleteItem menghapus item
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
		response.BadRequest(w, "Item ID is required")
		return
	}

	err := h.service.DeleteItem(id, userID)
	if err != nil {
		if err.Error() == message.MsgItemNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "unauthorized") {
			response.Unauthorized(w, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, nil, "Item deleted successfully")
}
