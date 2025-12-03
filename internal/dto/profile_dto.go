package dto

import "time"

/*
HosterProfileResponse adalah response untuk menampilkan profil hoster.
Berisi informasi dasar bisnis rental outdoor.
*/
type HosterProfileResponse struct {
	ID            string    `json:"id"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	StoreName     string    `json:"store_name"`
	Description   string    `json:"description"`
	Website       string    `json:"website,omitempty"`
	Instagram     string    `json:"instagram,omitempty"`
	Tiktok        string    `json:"tiktok,omitempty"`
	ProfilePhoto  string    `json:"profile_photo,omitempty"`
	JoinedAt      time.Time `json:"joined_at"`       // Tanggal bergabung
	DaysSinceJoin int       `json:"days_since_join"` // Jumlah hari bergabung
}

/*
UpdateHosterProfileRequest adalah request untuk update profil hoster.
Field yang bisa diubah: address, phone_number, description, website, instagram, tiktok.
Semua field optional kecuali address dan phone_number.
*/
type UpdateHosterProfileRequest struct {
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Description string `json:"description,omitempty"` // Optional
	Website     string `json:"website,omitempty"`     // Optional
	Instagram   string `json:"instagram,omitempty"`   // Optional
	Tiktok      string `json:"tiktok,omitempty"`      // Optional
}
