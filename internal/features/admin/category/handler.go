package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/response"
)

/*
CategoryHandler menangani HTTP requests untuk kategori.
*/
type CategoryHandler struct {
	service CategoryService
}

/*
NewCategoryHandler membuat instance handler.
*/
func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

/*
CreateCategory menangani POST /api/v1/admin/category
*/
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	resp, err := h.service.CreateCategory(req)
	if err != nil {
		log.Printf("CreateCategory: service error: %v", err)
		if err.Error() == message.CategoryAlreadyExists {
			response.BadRequest(w, message.CategoryAlreadyExists)
			return
		}
		if err.Error() == "category name required" {
			response.BadRequest(w, "category name required")
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, resp, "category created successfully")
}

/*
GetAllCategory menangani GET /api/v1/admin/category
*/
func (h *CategoryHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("GetAllCategory: service error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, categories, "categories retrieved successfully")
}

/*
GetCategoryByID menangani GET /api/v1/admin/category/{id}
*/
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		log.Printf("GetCategoryByID: service error: %v", err)
		if err.Error() == "category not found" {
			response.NotFound(w, "category not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, category, "category retrieved successfully")
}

/*
UpdateCategory menangani PUT /api/v1/admin/category/{id}
*/
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	resp, err := h.service.UpdateCategory(id, req)
	if err != nil {
		log.Printf("UpdateCategory: service error: %v", err)
		if err.Error() == message.CategoryAlreadyExists {
			response.BadRequest(w, message.CategoryAlreadyExists)
			return
		}
		if err.Error() == "category not found" {
			response.NotFound(w, "category not found")
			return
		}
		if err.Error() == "category name required" || err.Error() == "category ID required" {
			response.BadRequest(w, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, resp, "category updated successfully")
}

/*
DeleteCategory menangani DELETE /api/v1/admin/category/{id}
*/
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteCategory(id); err != nil {
		log.Printf("DeleteCategory: service error: %v", err)
		if err.Error() == "category not found" {
			response.NotFound(w, "category not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, nil, "category deleted successfully")
}
