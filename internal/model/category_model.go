package model

import "time"

// Model untuk entitas kategori
type CategoryModel struct {
	ID          string    `json:"id" db:"id"`                   // ID unik kategori
	Name        string    `json:"name" db:"name"`               // Nama kategori
	Description string    `json:"description" db:"description"` // Deskripsi kategori
	CreatedAt   time.Time `json:"created_at" db:"created_at"`   // Waktu pembuatan
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // Waktu pembaruan
}
