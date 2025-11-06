package model

import "time"

type HosterModel struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	ProfilePhoto string    `json:"profile_photo" db:"profile_photo"`
	StoreName    string    `json:"store_name" db:"store_name"`
	Description  string    `json:"description" db:"description"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	Email        string    `json:"email" db:"email"`
	Address      string    `json:"address" db:"address"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
