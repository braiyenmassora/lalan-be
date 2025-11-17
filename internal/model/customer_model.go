package model

import "time"

/*
Merepresentasikan data pelanggan dengan field yang diperlukan.
Digunakan untuk serialisasi JSON dan interaksi database.
*/
type CustomerModel struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	ProfilePhoto string    `json:"profile_photo,omitempty" db:"profile_photo"`
	PhoneNumber  string    `json:"phone_number,omitempty" db:"phone_number"`
	Email        string    `json:"email" db:"email"`
	Address      string    `json:"address,omitempty" db:"address"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
