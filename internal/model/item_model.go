package model

import "time"

/* PickupMethod mendefinisikan jenis metode pengambilan item. */
type PickupMethod string

const (
	PickupMethodSelfPickup PickupMethod = "pickup"
	PickupMethodDelivery   PickupMethod = "delivery"
)

/* ItemModel merepresentasikan struktur data untuk item. */
type ItemModel struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	Photos      []string     `json:"photos" db:"photos"`
	Stock       int          `json:"stock" db:"stock"`
	PickupType  PickupMethod `json:"pickup_type"`
	PricePerDay int          `json:"price_per_day" db:"price_per_day"`
	Deposit     int          `json:"deposit" db:"deposit"`
	Discount    int          `json:"discount,omitempty" db:"discount"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`

	// Foreign key
	CategoryID string `json:"category_id" db:"category_id"`
	UserID     string `json:"user_id" db:"user_id"`
}
