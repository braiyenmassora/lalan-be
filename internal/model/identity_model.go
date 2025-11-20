package model

import "time"

/*
IdentityModel
struct untuk data identitas dengan field JSON dan database
*/
type IdentityModel struct {
	ID         string     `json:"id" db:"id"`
	KTPURL     string     `json:"ktp_url" db:"ktp_url"`
	Verified   bool       `json:"verified" db:"verified"`
	Status     string     `json:"status" db:"status"`
	Reason     string     `json:"reason" db:"reason"`
	VerifiedAt *time.Time `json:"verified_at" db:"verified_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	UserID     string     `json:"user_id" db:"user_id"`
}
