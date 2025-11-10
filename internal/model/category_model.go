package model

import "time"

/*
Merepresentasikan data kategori.
Digunakan untuk mapping JSON dan database dengan field ID, nama, deskripsi, dan timestamp.
*/
type CategoryModel struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
