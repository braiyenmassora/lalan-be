// ===================================================================
// File: tnc_dto.go
// Deskripsi: DTO untuk Terms and Conditions (Syarat & Ketentuan)
// Catatan: SEMUA DTO TnC HANYA di file ini!
// ===================================================================

package dto

import "time"

// ===================================================================
// REQUEST DTO - HOSTER
// ===================================================================

// CreateTnCRequest adalah payload untuk create T&C baru oleh hoster
// Endpoint: POST /api/v1/hoster/tnc
//
// Contoh JSON:
//
//	{
//	  "description": [
//	    "Penyewa wajib mengembalikan barang dalam kondisi baik",
//	    "Keterlambatan pengembalian dikenakan denda"
//	  ]
//	}
type CreateTnCRequest struct {
	Description []string `json:"description"` // Array poin-poin T&C
}

// UpdateTnCRequest adalah payload untuk update T&C oleh hoster
// Endpoint: PUT /api/v1/hoster/tnc/{id}
//
// Contoh JSON:
//
//	{
//	  "description": [
//	    "Penyewa wajib mengembalikan barang dalam kondisi baik dan bersih",
//	    "Keterlambatan pengembalian dikenakan denda 10% per hari"
//	  ]
//	}
type UpdateTnCRequest struct {
	Description []string `json:"description"` // Array poin-poin T&C yang baru
}

// ===================================================================
// RESPONSE DTO - HOSTER
// ===================================================================

// TnCResponse adalah response untuk T&C (create, update, detail)
// Endpoint: POST/PUT/GET /api/v1/hoster/tnc
//
// Contoh JSON:
//
//	{
//	  "id": "uuid-tnc-123",
//	  "hoster_id": "uuid-hoster-123",
//	  "description": [
//	    "Penyewa wajib mengembalikan barang dalam kondisi baik",
//	    "Keterlambatan pengembalian dikenakan denda"
//	  ],
//	  "created_at": "2025-11-17T14:37:45Z",
//	  "updated_at": "2025-11-17T15:06:08Z"
//	}
type TnCResponse struct {
	ID          string    `json:"id" db:"id"`
	HosterID    string    `json:"hoster_id" db:"hoster_id"`
	Description []string  `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
