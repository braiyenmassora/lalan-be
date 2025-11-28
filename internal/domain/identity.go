// ===================================================================
// File: identity.go
// Deskripsi: Entity Identity (KTP/Verifikasi Identitas)
// Catatan: INI SATU-SATUNYA tempat untuk model Identity. JANGAN duplikasi!
// ===================================================================

package domain

import "time"

// ===================================================================
// IDENTITY (KTP)
// ===================================================================

// Identity adalah entity untuk data KTP user (customer/hoster).
// Digunakan untuk proses verifikasi identitas sebelum user bisa melakukan booking.
//
// Flow:
// 1. User upload foto KTP → status = "pending", verified = false
// 2. Admin verifikasi:
//   - Approve → status = "approved", verified = true, verified_at diisi
//   - Reject → status = "rejected", verified = false, reason diisi alasan penolakan
//
// 3. User bisa re-upload jika ditolak (data lama di-override atau buat baru)
//
// Relasi:
// - Satu user bisa punya banyak record Identity (history upload)
// - Hanya yang verified=true yang dipakai untuk booking
type Identity struct {
	ID         string     `json:"id" db:"id"`
	UserID     string     `json:"user_id" db:"user_id"`         // ID customer/hoster
	KTPURL     string     `json:"ktp_url" db:"ktp_url"`         // URL foto KTP (dari cloud storage)
	Verified   bool       `json:"verified" db:"verified"`       // true jika approved, false jika pending/rejected
	Status     string     `json:"status" db:"status"`           // "pending", "approved", "rejected"
	Reason     string     `json:"reason" db:"reason"`           // Alasan reject (kosong jika pending/approved)
	VerifiedAt *time.Time `json:"verified_at" db:"verified_at"` // Waktu admin approve/reject (nullable)
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
