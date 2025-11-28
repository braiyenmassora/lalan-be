// ===================================================================
// File: item.go
// Deskripsi: Entity Item dan Category
// Catatan: SEMUA model item HANYA di file ini!
// ===================================================================

package domain

import "time"

// ===================================================================
// PICKUP METHOD (Enum/Constant)
// ===================================================================

// PickupMethod adalah enum untuk metode pengambilan item
type PickupMethod string

const (
	// PickupMethodSelfPickup: Customer ambil sendiri ke lokasi hoster
	PickupMethodSelfPickup PickupMethod = "pickup"

	// PickupMethodDelivery: Hoster kirim ke alamat customer
	PickupMethodDelivery PickupMethod = "delivery"
)

// ===================================================================
// CATEGORY
// ===================================================================

// Category adalah entity untuk kategori item.
// Digunakan untuk mengelompokkan item agar mudah dicari oleh customer.
//
// Contoh kategori: "Kamera", "Laptop", "Tenda", "Alat Outdoor", dll
//
// Relasi:
// - Category has many Item
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ===================================================================
// ITEM
// ===================================================================

// Item adalah entity untuk barang/jasa yang bisa disewakan.
// Hoster membuat item untuk ditampilkan di platform dan bisa di-booking oleh customer.
//
// Field penting:
// - Stock: Jumlah unit yang tersedia. Jika 0, item tidak bisa di-booking
// - PickupType: Metode pengambilan (pickup/delivery)
// - PricePerDay: Harga sewa per hari per unit
// - Deposit: Uang jaminan per unit (dikembalikan jika item tidak rusak)
// - Discount: Diskon dalam nominal (opsional)
//
// Relasi:
// - Item belongs to Hoster (user_id)
// - Item belongs to Category (category_id)
// - Item has many TermsAndConditions
type Item struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	Photos      []string     `json:"photos" db:"photos"`               // Array URL foto item
	Stock       int          `json:"stock" db:"stock"`                 // Jumlah unit tersedia
	PickupType  PickupMethod `json:"pickup_type" db:"pickup_type"`     // "pickup" atau "delivery"
	PricePerDay int          `json:"price_per_day" db:"price_per_day"` // Harga sewa per hari (dalam satuan terkecil, misal: rupiah)
	Deposit     int          `json:"deposit" db:"deposit"`             // Deposit per unit
	Discount    int          `json:"discount,omitempty" db:"discount"` // Diskon (opsional)
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	CategoryID  string       `json:"category_id" db:"category_id"` // FK ke Category
	UserID      string       `json:"user_id" db:"user_id"`         // FK ke Hoster
}
