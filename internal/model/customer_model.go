package model

import "time"

type CustomerModel struct {
	ID                    string     `json:"id" db:"id"`
	FullName              string     `json:"full_name" db:"full_name"`
	ProfilePhoto          string     `json:"profile_photo,omitempty" db:"profile_photo"`
	PhoneNumber           string     `json:"phone_number,omitempty" db:"phone_number"`
	Email                 string     `json:"email" db:"email"`
	Address               string     `json:"address,omitempty" db:"address"`
	PasswordHash          string     `json:"-" db:"password_hash"`
	EmailVerified         bool       `json:"email_verified" db:"email_verified"`
	VerificationToken     string     `json:"-" db:"verification_token"`
	VerificationExpiresAt *time.Time `json:"-" db:"verification_expire"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

/*
CustomerIdentityDTO
struct untuk data customer dengan identity
*/
type CustomerIdentityDTO struct {
	CustomerID  string `json:"customer_id" db:"customer_id"` // customer.id
	FullName    string `json:"full_name" db:"full_name"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`

	IdentityID string    `json:"identity_id" db:"identity_id"` // identity.id
	KTPURL     string    `json:"ktp_url" db:"ktp_url"`
	Verified   bool      `json:"verified" db:"verified"`
	Status     string    `json:"status" db:"status"`
	Reason     string    `json:"reason" db:"reason"`
	VerifiedAt time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt  time.Time `json:"created_at" db:"identity_created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"identity_updated_at"`
}
