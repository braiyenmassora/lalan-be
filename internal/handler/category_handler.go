package handler

import (
	"lalan-be/internal/service"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// add new category
