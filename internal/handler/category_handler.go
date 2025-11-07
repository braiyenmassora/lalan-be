package handler

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"net/http"
)

// Struct untuk menangani request HTTP kategori
type CategoryHandler struct {
	service service.CategoryService // Service untuk operasi kategori
}

// Membuat instance CategoryHandler baru dengan dependency injection
func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// Struct untuk payload request JSON kategori
type CategoryRequest struct {
	Name        string `json:"name"`        // Nama kategori
	Description string `json:"description"` // Deskripsi kategori
}

// Handler untuk menambah kategori
func (h *CategoryHandler) AddCategory(w http.ResponseWriter, r *http.Request) {
	// Validasi metode HTTP POST
	if r.Method != http.MethodPost {
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Message not allowed")
		return
	}
	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)

	// Mencegah field tidak dikenal
	decoder.DisallowUnknownFields()
	// Decode payload JSON
	if err := decoder.Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Konversi request ke model domain
	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	// Panggil service untuk menambah kategori
	if err := h.service.AddCategory(input); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Encode response JSON
	json.NewEncoder(w).Encode(response.Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: message.MsgCategoryCreatedSuccess,
		Data:    nil,
	})
}
