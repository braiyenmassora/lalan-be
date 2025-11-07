package handler

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ItemHandler struct {
	service service.ItemService
}

func NewItemHandlerS(s service.ItemService) *ItemHandler {
	return &ItemHandler{service: s}
}

type PickupMethod string

var (
	PickupMethodSelfPickup PickupMethod = "pickup"
	PickupMethodDelivery   PickupMethod = "delivery"
)

type ItemRequest struct {
	ID          string       `json:"id" db:"id"`                       // ID unik item
	Name        string       `json:"name" db:"name"`                   // Nama tampilan item
	Description string       `json:"description" db:"description"`     // Deskripsi detail item
	Photos      []string     `json:"photos" db:"photos"`               // Array URL atau path gambar item
	Stock       int          `json:"stock" db:"stock"`                 // Jumlah stok tersedia
	PickupType  PickupMethod `json:"pickup_type"`                      // "pickup" or "delivery"
	PricePerDay int          `json:"price_per_day" db:"price_per_day"` // Biaya sewa harian
	Deposit     int          `json:"deposit" db:"deposit"`             // Jumlah deposit keamanan
	Discount    int          `json:"discount,omitempty" db:"discount"` // Diskon persentase, dihilangkan jika nol
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`       // Waktu pembuatan awal
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`       // Waktu pembaruan terakhir

	// Foreign key
	CategoryID string `json:"category_id" db:"category_id"` // ID kategori terkait
	UserID     string `json:"user_id" db:"user_id"`         // ID pemilik item
}

// handler tambah item
func (h *ItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.ErrorResponse(w, http.StatusMethodNotAllowed, message.MsgNotAllowed)
		return
	}
	var req ItemRequest
	decoder := json.NewDecoder(r.Body)
	// membatasi field tidak di kenal
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, message.MsgBadRequest)
		return
	}
	// konversi model
	input := &model.ItemModel{
		ID:          uuid.NewString(), // generate ID unik
		Name:        req.Name,
		Description: req.Description,
		Photos:      req.Photos, // pastikan ini []string
		Stock:       req.Stock,
		PickupType:  model.PickupMethod(req.PickupType),
		PricePerDay: req.PricePerDay,
		Deposit:     req.Deposit,
		Discount:    req.Discount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CategoryID:  req.CategoryID,
		UserID:      req.UserID,
	}
	if err := h.service.AddItem(input); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Encode response JSON
	json.NewEncoder(w).Encode(response.Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: message.MsgItemCreatedSuccess,
		Data:    nil,
	})
}
