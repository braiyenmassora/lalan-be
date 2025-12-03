package dto

import "time"

// ===================================================================
// REQUEST DTO
// ===================================================================

// CreateCategoryRequest adalah payload untuk endpoint POST /admin/category
// Digunakan admin untuk membuat kategori baru
//
// Contoh JSON:
//
//	{
//	  "name": "Kamera",
//	  "description": "Kategori untuk peralatan kamera dan fotografi"
//	}
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateCategoryRequest adalah payload untuk endpoint PUT /admin/category/{id}
// Digunakan admin untuk update kategori
//
// Contoh JSON:
//
//	{
//	  "name": "Kamera & Fotografi",
//	  "description": "Kategori untuk semua peralatan kamera dan fotografi"
//	}
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ===================================================================
// RESPONSE DTO
// ===================================================================

// CategoryResponse adalah response untuk detail kategori
// Digunakan untuk response create, update, get detail
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryListResponse adalah response untuk list kategori
// Digunakan untuk endpoint GET list categories
type CategoryListResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Legacy DTO - untuk backward compatibility
type CategoryListByHosterResponse struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
