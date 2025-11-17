package model

import "time"

/*
Merepresentasikan data terms and conditions dengan field yang diperlukan.
Digunakan untuk serialisasi JSON dan interaksi database.
*/
type TermsAndConditionsModel struct {
	ID          string    `json:"id" db:"id"`
	Description []string  `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UserID      string    `json:"user_id" db:"user_id"`
	ItemID      string    `json:"item_id,omitempty" db:"item_id"`
}
