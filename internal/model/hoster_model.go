package model

import "time"

// Model untuk entitas hoster
type HosterModel struct {
	ID           string    `json:"id" db:"id"`                       // ID unik hoster
	FullName     string    `json:"full_name" db:"full_name"`         // Nama lengkap hoster
	ProfilePhoto string    `json:"profile_photo" db:"profile_photo"` // URL foto profil
	StoreName    string    `json:"store_name" db:"store_name"`       // Nama toko
	Description  string    `json:"description" db:"description"`     // Deskripsi toko
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`   // Nomor telepon
	Email        string    `json:"email" db:"email"`                 // Alamat email
	Address      string    `json:"address" db:"address"`             // Alamat
	PasswordHash string    `json:"-" db:"password_hash"`             // Hash kata sandi (tidak diekspor JSON)
	CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Waktu pembuatan
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`       // Waktu pembaruan
}
