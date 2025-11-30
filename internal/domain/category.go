package domain

import "time"

// ===================================================================
// CATEGORY
// ===================================================================

// Category adalah entity untuk kategori item.
// Digunakan untuk mengelompokkan item agar mudah dicari oleh customer.
//
// Contoh kategori: "Kamera", "Laptop", "Tenda", "Alat Outdoor", dll
//
// Relasi:
// - Category has many Item
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
