package model

import "time"

type PickupMethod string

const (
	PickupMethodSelfPickup PickupMethod = "pickup"
	PickupMethodDelivery   PickupMethod = "delivery"
)

// ItemModel represents the data structure for an item entity, used for JSON serialization and database interactions.
type ItemModel struct {
	ID          string       `json:"id" db:"id"`                       // ID unik item
	Name        string       `json:"name" db:"name"`                   // Nama tampilan item
	Description string       `json:"description" db:"description"`     // Deskripsi detail item
	Photos      []string     `json:"photos" db:"photos"`               // Array URL atau path gambar item
	Stock       int          `json:"stock" db:"stock"`                 // Jumlah stok tersedia
	PickupType  PickupMethod `json:"pickup_type"`                      // "pickup" or "delivery"
	PricePerDay int          `json:"price_per_day" db:"price_per_day"` // Biaya sewa harian
	Deposit     int          `json:"deposit" db:"deposit"`             // Jumlah deposit keamanan
	Discount    int          `json:"discount,omitempty" db:"discount"` // Diskon persentase, dihilangkan jika nol
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`       // Waktu pembuatan awal
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`       // Waktu pembaruan terakhir

	// Foreign key
	CategoryID string `json:"category_id" db:"category_id"` // ID kategori terkait
	UserID     string `json:"user_id" db:"user_id"`         // ID pemilik item
}
