package model

import "time"

/*
CustomerModel
struct untuk data customer dengan field JSON dan database
*/
type CustomerModel struct {
	ID                    string     `json:"id" db:"id"`
	FullName              string     `json:"full_name" db:"full_name"`
	ProfilePhoto          string     `json:"profile_photo,omitempty" db:"profile_photo"`
	PhoneNumber           string     `json:"phone_number,omitempty" db:"phone_number"`
	Email                 string     `json:"email" db:"email"`
	Address               string     `json:"address,omitempty" db:"address"`
	PasswordHash          string     `json:"-" db:"password_hash"`
	EmailVerified         bool       `json:"email_verified" db:"email_verified"`
	VerificationToken     string     `json:"-" db:"verification_token"`  // OTP 6 digit
	VerificationExpiresAt *time.Time `json:"-" db:"verification_expire"` // OTP expiry time
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}
