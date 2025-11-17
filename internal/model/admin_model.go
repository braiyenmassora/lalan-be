package model

import "time"

/*
Merepresentasikan data admin dengan field yang diperlukan.
Digunakan untuk serialisasi JSON dan interaksi database.
*/
type AdminModel struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
