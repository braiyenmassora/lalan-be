// ===================================================================
// File: identity_dto.go
// Deskripsi: DTO untuk Identity/KTP (Customer & Admin)
// Catatan: SEMUA DTO identity HANYA di file ini!
// ===================================================================

package dto

import "time"

// ===================================================================
// REQUEST DTO - CUSTOMER
// ===================================================================

// UploadIdentityByCustomerRequest adalah payload saat customer upload KTP pertama kali
// Endpoint: POST /customer/identity/upload
//
// Catatan: File upload dilakukan via multipart/form-data, service akan dapat URL setelah upload
type UploadIdentityByCustomerRequest struct {
	UserID string `json:"user_id"` // ID customer yang upload KTP
	KTPURL string `json:"ktp_url"` // URL foto KTP yang sudah di-upload ke storage
}

// ReuploadIdentityByCustomerRequest adalah payload saat customer re-upload KTP
// Endpoint: PUT /customer/identity/upload
type ReuploadIdentityByCustomerRequest struct {
	KTPURL string `json:"ktp_url"` // URL KTP baru
}

// ===================================================================
// REQUEST DTO - ADMIN
// ===================================================================

// VerifyIdentityByAdminRequest adalah payload saat admin verifikasi KTP
// Endpoint: POST /admin/identity/{id}/verify
//
// Contoh JSON (Approve):
//
//	{
//	  "status": "approved",
//	  "reason": ""
//	}
//
// Contoh JSON (Reject):
//
//	{
//	  "status": "rejected",
//	  "reason": "Foto KTP buram, silakan upload ulang dengan foto yang lebih jelas"
//	}
type VerifyIdentityByAdminRequest struct {
	Status string `json:"status"` // "approved" atau "rejected"
	Reason string `json:"reason"` // Wajib diisi jika status = "rejected"
}

// ===================================================================
// RESPONSE DTO
// ===================================================================

// IdentityStatusByCustomerResponse adalah response status KTP customer
// Endpoint: GET /customer/identity/status
//
// Contoh JSON:
//
//	{
//	  "ktp_id": "uuid-ktp-123",
//	  "user_id": "uuid-customer-123",
//	  "ktp_url": "https://storage.com/ktp/customer-123.jpg",
//	  "created_at": "2025-11-28T10:00:00Z",
//	  "status": "pending",
//	  "verified": false,
//	  "reason": "",
//	  "verified_at": null
//	}
type IdentityStatusByCustomerResponse struct {
	KTPID      string     `json:"ktp_id,omitempty"`
	UserID     string     `json:"user_id"`
	KTPURL     string     `json:"ktp_url,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	Status     string     `json:"status"`                // "pending", "approved", "rejected"
	Verified   bool       `json:"verified"`              // true jika approved, false jika pending/rejected
	Reason     string     `json:"reason"`                // Alasan approve/reject dari admin
	VerifiedAt *time.Time `json:"verified_at,omitempty"` // Waktu verifikasi oleh admin
}

// IdentityListByAdminResponse adalah response untuk list semua KTP yang perlu diverifikasi
// Endpoint: GET /admin/identity/pending
type IdentityListByAdminResponse struct {
	KTPID     string    `json:"ktp_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
	KTPURL    string    `json:"ktp_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
