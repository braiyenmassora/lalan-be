package model

import "time"

/*
Merepresentasikan data pelanggan.
Digunakan untuk mapping JSON dan database dengan field ID, nama, email, dll.
*/
type AdminModel struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
