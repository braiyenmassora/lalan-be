// ===================================================================
// File: user.go
// Deskripsi: Entity User - Admin, Hoster, dan Customer
// Catatan: SEMUA model user HANYA di file ini. JANGAN buat di tempat lain!
// ===================================================================

package domain

import "time"

// ===================================================================
// ADMIN
// ===================================================================

// Admin adalah entity untuk user dengan role admin.
// Admin memiliki akses penuh untuk mengelola sistem, termasuk verifikasi KTP.
type Admin struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ===================================================================
// HOSTER
// ===================================================================

// Hoster adalah entity untuk user yang menyewakan item/property.
// Hoster dapat membuat item, menerima booking, dan mengelola inventori.
type Hoster struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	ProfilePhoto string    `json:"profile_photo" db:"profile_photo"`
	StoreName    string    `json:"store_name" db:"store_name"`
	Description  string    `json:"description" db:"description"`
	Website      string    `json:"website,omitempty" db:"website"`
	Instagram    string    `json:"instagram,omitempty" db:"instagram"`
	Tiktok       string    `json:"tiktok,omitempty" db:"tiktok"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	Email        string    `json:"email" db:"email"`
	Address      string    `json:"address" db:"address"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ===================================================================
// CUSTOMER
// ===================================================================

// Customer adalah entity untuk user yang melakukan booking/sewa item.
// Customer wajib verifikasi email dan upload KTP sebelum bisa booking.
type Customer struct {
	ID                    string    `json:"id" db:"id"`
	FullName              string    `json:"full_name" db:"full_name"`
	ProfilePhoto          string    `json:"profile_photo,omitempty" db:"profile_photo"`
	PhoneNumber           string    `json:"phone_number,omitempty" db:"phone_number"`
	Email                 string    `json:"email" db:"email"`
	Address               string    `json:"address,omitempty" db:"address"`
	PasswordHash          string    `json:"-" db:"password_hash"`
	EmailVerified         bool      `json:"email_verified" db:"email_verified"`
	VerificationToken     string    `json:"-" db:"verification_token"`
	VerificationExpiresAt time.Time `json:"-" db:"verification_expire"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}
