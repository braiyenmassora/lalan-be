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

/*
Menangani operasi kategori.
Menyediakan endpoint untuk menambah, mengambil, update, dan hapus kategori dengan respons sukses atau error.
*/
type CategoryHandler struct {
	service service.CategoryService
}

/*
Merepresentasikan data request kategori.
Digunakan untuk decoding JSON dan validasi input.
*/
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
Membuat handler kategori.
Mengembalikan instance CategoryHandler yang siap digunakan.
*/
func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

/*
Menambahkan kategori baru.
Mengembalikan respons pembuatan sukses atau error validasi.
*/
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
		response.BadRequest(w, message.MsgCategoryNameRequired)
		return
	}

	if len(req.Name) > 255 {
		response.BadRequest(w, message.MsgCategoryNameTooLong)
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

/*
Mengambil semua kategori.
Mengembalikan daftar kategori sukses atau error.
*/
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

/*
Mengambil kategori berdasarkan ID.
Mengembalikan data kategori sukses atau error jika tidak ditemukan.
*/
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
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

/*
Mengupdate kategori.
Mengembalikan respons update sukses atau error validasi/not found.
*/
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
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
		response.BadRequest(w, message.MsgCategoryNameRequired)
		return
	}

	if len(req.Name) > 255 {
		response.BadRequest(w, message.MsgCategoryNameTooLong)
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

	response.OK(w, categoryResp, message.MsgCategoryUpdatedSuccess)
}

/*
Menghapus kategori.
Mengembalikan respons penghapusan sukses atau error jika tidak ditemukan.
*/
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, message.MsgCategoryIDRequired)
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

	response.OK(w, nil, message.MsgCategoryDeletedSuccess)
}
