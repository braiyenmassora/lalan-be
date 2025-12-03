// ===================================================================
// File: public_dto.go
// Deskripsi: DTO untuk Public endpoints (Category, Item, Terms)
// Catatan: SEMUA DTO public HANYA di file ini!
// ===================================================================

package dto

import "time"

// ===================================================================
// RESPONSE DTO - PUBLIC
// ===================================================================

// CategoryPublicResponse adalah response untuk data kategori
// Endpoint: GET /public/categories
//
// Contoh JSON:
//
//	{
//	  "id": "uuid-category-123",
//	  "name": "Kamera",
//	  "description": "Kamera DSLR, Mirrorless, Action Cam, dll",
//	  "created_at": "2025-11-01T00:00:00Z",
//	  "updated_at": "2025-11-01T00:00:00Z"
//	}
type CategoryPublicResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ItemPublicResponse adalah response untuk data item publik
// Endpoint: GET /public/items
//
// Contoh JSON:
//
//	{
//	  "id": "uuid-item-123",
//	  "name": "Kamera DSLR Canon EOS 80D",
//	  "description": "Kamera DSLR 24MP dengan lensa kit",
//	  "photos": ["https://storage.com/item1.jpg", "https://storage.com/item2.jpg"],
//	  "stock": 5,
//	  "pickup_type": "pickup",
//	  "price_per_day": 100000,
//	  "deposit": 500000,
//	  "discount": 0,
//	  "created_at": "2025-11-01T00:00:00Z",
//	  "updated_at": "2025-11-01T00:00:00Z",
//	  "category_id": "uuid-category-123",
//	  "user_id": "uuid-hoster-123"
//	}
type ItemPublicResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Photos      []string  `json:"photos"`
	Stock       int       `json:"stock"`
	PickupType  string    `json:"pickup_type"`
	PricePerDay int       `json:"price_per_day"`
	Deposit     int       `json:"deposit"`
	Discount    int       `json:"discount,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CategoryID  string    `json:"category_id"`
	HosterID    string    `json:"hoster_id"`
}

// TermsAndConditionsPublicResponse adalah response untuk syarat dan ketentuan
// Endpoint: GET /public/terms
//
// Contoh JSON:
//
//	{
//	  "id": "uuid-tnc-123",
//	  "description": [
//	    "Barang harus dikembalikan dalam kondisi bersih",
//	    "Keterlambatan pengembalian dikenakan denda 10% per hari"
//	  ],
//	  "created_at": "2025-11-01T00:00:00Z",
//	  "updated_at": "2025-11-01T00:00:00Z",
//	  "user_id": "uuid-hoster-123"
//	}
type TermsAndConditionsPublicResponse struct {
	ID          string    `json:"id"`
	Description []string  `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id"`
}

// ===================================================================
// ITEM DETAIL WITH JOIN - PUBLIC
// ===================================================================

// ItemDetailResponse adalah response lengkap untuk detail item dengan JOIN
// Endpoint: GET /public/item/{id}
// Menggabungkan data item, category, hoster, dan terms & conditions dalam 1 response
//
// Contoh JSON:
//
//	{
//	  "item": {...},
//	  "category": {...},
//	  "hoster": {...},
//	  "terms_and_conditions": [...],
//	  "booked_dates": ["2025-12-05", "2025-12-06", "2025-12-07"]
//	}
type ItemDetailResponse struct {
	Item               ItemDetail     `json:"item"`
	Category           CategoryDetail `json:"category"`
	Hoster             HosterDetail   `json:"hoster"`
	TermsAndConditions []string       `json:"terms_and_conditions"`
	BookedDates        []string       `json:"booked_dates"`
}

// ItemDetail adalah detail item untuk response detail
type ItemDetail struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Photos      []string  `json:"photos"`
	Stock       int       `json:"stock"`
	PickupType  string    `json:"pickup_type"`
	PricePerDay int       `json:"price_per_day"`
	Deposit     int       `json:"deposit"`
	Discount    int       `json:"discount,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryDetail adalah detail kategori untuk response detail
type CategoryDetail struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// HosterDetail adalah detail hoster (pemilik item) untuk response detail
type HosterDetail struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	StoreName    string `json:"store_name"`
	Description  string `json:"description"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
	ProfilePhoto string `json:"profile_photo,omitempty"`
	Website      string `json:"website,omitempty"`
	Instagram    string `json:"instagram,omitempty"`
	Tiktok       string `json:"tiktok,omitempty"`
}
