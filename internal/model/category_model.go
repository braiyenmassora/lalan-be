package model

import "time"

/*
Merepresentasikan data kategori dengan field yang diperlukan.
Digunakan untuk serialisasi JSON dan interaksi database.
*/
type CategoryModel struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
