package model

import "time"

type HosterModel struct {
	ID          string    `json:"id" db:"id"`
	OwnerName   string    `json:"owner_name" db:"owner_name"`
	StoreName   string    `json:"store_name" db:"store_name"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	Address     string    `json:"address" db:"address"`
	Password    string    `json:"-" db:"password"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
