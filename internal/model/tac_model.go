package model

import "time"

/*
Merepresentasikan data terms and conditions.
Digunakan untuk mapping JSON dan database dengan field deskripsi, timestamp, dan relasi user.
*/
type TermsAndConditionsModel struct {
	ID          string    `json:"id" db:"id"`
	Description []string  `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Foreign key
	UserID string `json:"user_id" db:"user_id"`
}
