package dto

import "time"

// ===================================================================
// RESPONSE DTO - HOSTER
// ===

type PickupMethod string

const (
	// PickupMethodSelfPickup: Customer ambil sendiri ke lokasi hoster
	PickupMethodSelfPickup PickupMethod = "self_pickup"

	// PickupMethodDelivery: Hoster kirim ke alamat customer
	PickupMethodDelivery PickupMethod = "delivery"
)

type ItemListByHosterResponse struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Stock       int          `json:"stock" db:"stock"`
	PricePerDay int          `json:"price_per_day" db:"price_per_day"`
	PickupType  PickupMethod `json:"pickup_type" db:"pickup_type"`
}

type ItemDetailByHosterResponse struct {
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
	HosterID    string       `json:"hoster_id" db:"hoster_id"`     // FK ke Hoster
}

type CreateItemByCustomerRequest struct {
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	Photos      []string     `json:"photos" db:"photos"`               // Array URL foto item
	Stock       int          `json:"stock" db:"stock"`                 // Jumlah unit tersedia
	PickupType  PickupMethod `json:"pickup_type" db:"pickup_type"`     // "pickup" atau "delivery"
	PricePerDay int          `json:"price_per_day" db:"price_per_day"` // Harga sewa per hari (dalam satuan terkecil, misal: rupiah)
	Deposit     int          `json:"deposit" db:"deposit"`             // Deposit per unit
	Discount    int          `json:"discount,omitempty" db:"discount"` // Diskon (opsional)
	CategoryID  string       `json:"category_id" db:"category_id"`     // FK ke Category
}

type UpdateItemRequestRequest struct {
	Stock      *int          `json:"stock,omitempty"`       // Pointer untuk opsional
	PickupType *PickupMethod `json:"pickup_type,omitempty"` // Gunakan PickupMethod untuk type-safe
	Deposit    *int          `json:"deposit,omitempty"`
	Discount   *int          `json:"discount,omitempty"`
}
