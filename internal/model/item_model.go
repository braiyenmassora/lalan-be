package model

import "time"

/*
Mendefinisikan metode pengambilan item yang tersedia.
Digunakan dalam validasi dan pemilihan cara pengambilan.
*/
const (
	PickupMethodSelfPickup PickupMethod = "pickup"
	PickupMethodDelivery   PickupMethod = "delivery"
)

/*
Mewakili tipe metode pengambilan item.
Digunakan untuk menentukan cara pengambilan dalam model item.
*/
type PickupMethod string

/*
Merepresentasikan data item dengan field yang diperlukan.
Digunakan untuk serialisasi JSON dan interaksi database.
*/
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
