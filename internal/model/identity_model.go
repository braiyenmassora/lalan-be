package model

import "time"

type IdentityModel struct {
	ID             string     `json:"id" db:"id"`                           // Primary key
	KTPURL         string     `json:"ktp_url" db:"ktp_url"`                 // URL file KTP
	Verified       bool       `json:"verified" db:"verified"`               // Status verifikasi
	Status         string     `json:"status" db:"status"`                   // Enum: pending, approved, rejected
	RejectedReason string     `json:"rejected_reason" db:"rejected_reason"` // Alasan penolakan jika rejected
	VerifiedAt     *time.Time `json:"verified_at" db:"verified_at"`         // Timestamp verifikasi, nullable
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`

	// Foreign Key
	UserID string `json:"user_id" db:"user_id"`
}
