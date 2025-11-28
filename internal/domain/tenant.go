// ===================================================================
// File: tenant.go
// Deskripsi: Entity Tenant dan TermsAndConditions
// Catatan: SEMUA model tenant HANYA di file ini!
// ===================================================================

package domain

import "time"

// ===================================================================
// TENANT
// ===================================================================

// Tent adalah entity untuk tenant/penyewa toko.
// Relasi one-to-one dengan Hoster (satu hoster punya satu tenant).
//
// Catatan: Model ini kemungkinan belum fully digunakan di sistem saat ini.
// TODO: Verifikasi apakah model Tenant masih diperlukan atau bisa dihapus
//
// Relasi:
// - Tenant belongs to Hoster (one-to-one via hoster_id)
type Tenant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	HosterID  string    `json:"hoster_id" db:"hoster_id"` // FK ke Hoster (one-to-one)
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ===================================================================
// TERMS AND CONDITIONS (Syarat & Ketentuan)
// ===================================================================

// TermsAndConditions adalah entity untuk syarat dan ketentuan penyewaan.
// Setiap hoster atau item bisa punya T&C sendiri.
//
// Contoh T&C:
// - "Barang harus dikembalikan dalam kondisi bersih"
// - "Keterlambatan pengembalian dikenakan denda 10% per hari"
// - "Deposit akan dikembalikan maksimal 7 hari setelah pengembalian"
//
// Relasi:
// - TermsAndConditions belongs to Hoster (user_id) → T&C umum hoster
// - TermsAndConditions belongs to Item (item_id, nullable) → T&C spesifik item
type TermsAndConditions struct {
	ID          string    `json:"id" db:"id"`
	Description []string  `json:"description" db:"description"` // Array string (poin-poin T&C)
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UserID      string    `json:"user_id" db:"user_id"`           // FK ke Hoster
	ItemID      string    `json:"item_id,omitempty" db:"item_id"` // FK ke Item (nullable, kosong jika T&C umum hoster)
}
