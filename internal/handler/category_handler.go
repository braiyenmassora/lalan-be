package handler

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/internal/service"
	"lalan-be/pkg/message"
	"net/http"
	"strings"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddCategory menambahkan kategori baru
func (h *CategoryHandler) AddCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Name) == "" {
		response.BadRequest(w, "Category name is required")
		return
	}

	if len(req.Name) > 255 {
		response.BadRequest(w, "Category name must not exceed 255 characters")
		return
	}

	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	categoryResp, err := h.service.AddCategory(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, categoryResp, message.MsgCategoryCreatedSuccess)
}

// GetAllCategories mendapatkan semua kategori
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	categories, err := h.service.GetAllCategories()
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, categories, message.MsgSuccess)
}

// GetCategoryByID mendapatkan kategori berdasarkan ID
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "Category ID is required")
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		if err.Error() == message.MsgCategoryNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, category, message.MsgSuccess)
}

// UpdateCategory mengupdate kategori
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "Category ID is required")
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(req.Name) == "" {
		response.BadRequest(w, "Category name is required")
		return
	}

	if len(req.Name) > 255 {
		response.BadRequest(w, "Category name must not exceed 255 characters")
		return
	}

	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	categoryResp, err := h.service.UpdateCategory(id, input)
	if err != nil {
		if err.Error() == message.MsgCategoryNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, categoryResp, "Category updated successfully")
}

// DeleteCategory menghapus kategori
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "Category ID is required")
		return
	}

	err := h.service.DeleteCategory(id)
	if err != nil {
		if err.Error() == message.MsgCategoryNotFound {
			response.Error(w, http.StatusNotFound, err.Error())
		} else {
			response.BadRequest(w, err.Error())
		}
		return
	}

	response.OK(w, nil, "Category deleted successfully")
}
