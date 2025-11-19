package model

import "time"

/*
TenantModel
struct untuk data tenant dengan field JSON dan database
*/
type TenantModel struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	HosterID  string    `json:"hoster_id" db:"hoster_id"` // relasi one-to-one
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
