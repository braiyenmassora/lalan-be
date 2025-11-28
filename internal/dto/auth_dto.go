// ===================================================================
// File: auth_dto.go
// Deskripsi: DTO untuk Auth (Login, Register, Verify Email, Resend OTP)
// Catatan: Semua DTO auth HANYA di file ini! Jangan buat di handler!
// ===================================================================

package dto

// ===================================================================
// REQUEST DTO
// ===================================================================

// LoginRequest adalah payload untuk endpoint POST /auth/login
// Digunakan oleh semua role (admin, hoster, customer)
//
// Contoh JSON:
//
//	{
//	  "email": "customer@example.com",
//	  "password": "rahasia123"
//	}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest adalah payload untuk endpoint POST /auth/register
// Mendukung multiple roles: admin, hoster, customer
//
// Field wajib untuk semua role:
// - role, full_name, email, password
//
// Field opsional tergantung role:
// - Hoster: store_name, profile_photo
// - Customer: phone_number, address, profile_photo
//
// Contoh JSON (Customer):
//
//	{
//	  "role": "customer",
//	  "full_name": "Budi Santoso",
//	  "email": "budi@example.com",
//	  "password": "rahasia123",
//	  "phone_number": "081234567890",
//	  "address": "Jakarta Selatan"
//	}
//
// Contoh JSON (Hoster):
//
//	{
//	  "role": "hoster",
//	  "full_name": "Toko Rental ABC",
//	  "email": "rental@example.com",
//	  "password": "rahasia123",
//	  "store_name": "Rental ABC",
//	  "phone_number": "081234567890"
//	}
type RegisterRequest struct {
	Role         string `json:"role"` // "admin", "hoster", atau "customer"
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Address      string `json:"address,omitempty"`
	StoreName    string `json:"store_name,omitempty"` // Khusus hoster
	ProfilePhoto string `json:"profile_photo,omitempty"`
}

// VerifyEmailRequest adalah payload untuk endpoint POST /auth/verify-email
// Digunakan customer untuk verifikasi email dengan kode OTP
//
// Contoh JSON:
//
//	{
//	  "email": "customer@example.com",
//	  "otp": "123456"
//	}
type VerifyEmailRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// ResendOTPRequest adalah payload untuk endpoint POST /auth/resend-otp
// Digunakan customer jika tidak menerima OTP atau OTP kadaluarsa
//
// Contoh JSON:
//
//	{
//	  "email": "customer@example.com"
//	}
type ResendOTPRequest struct {
	Email string `json:"email"`
}

// ===================================================================
// RESPONSE DTO
// ===================================================================

// AuthResponse adalah response sukses untuk endpoint login
// Berisi token dan informasi dasar user
//
// Contoh JSON:
//
//	{
//	  "id": "uuid-customer-123",
//	  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	  "token_type": "Bearer",
//	  "expires_in": 86400,
//	  "role": "customer"
//	}
type AuthResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Role         string `json:"role"`
}

// CreateCustomerResponse adalah response sukses setelah register customer
// Berisi ID customer dan kode OTP (untuk development/testing)
//
// Contoh JSON:
//
//	{
//	  "customer_id": "uuid-customer-123",
//	  "otp": "123456"
//	}
type CreateCustomerResponse struct {
	CustomerID string `json:"customer_id"`
	OTP        string `json:"otp"` // Untuk development, production jangan return OTP!
}
