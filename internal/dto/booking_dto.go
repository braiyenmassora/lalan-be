// ===================================================================
// File: booking_dto.go
// Deskripsi: DTO untuk Booking - Customer dan Hoster
// Catatan: SEMUA DTO booking (request & response) HANYA di file ini!
//          Penamaan SUPER JELAS: [Action][Entity]By[Actor][Type]
// ===================================================================

package dto

import "time"

// ===================================================================
// REQUEST DTO - CUSTOMER
// ===================================================================

// CreateBookingByCustomerRequest adalah payload saat customer membuat booking baru
// Endpoint: POST /customer/booking
//
// Contoh JSON:
//
//	{
//	  "start_date": "2025-12-20",
//	  "end_date": "2025-12-25",
//	  "delivery_type": "pickup",
//	  "items": [
//	    {
//	      "item_id": "uuid-item-123",
//	      "name": "Kamera DSLR Canon",
//	      "quantity": 2,
//	      "price_per_day": 100000,
//	      "deposit_per_unit": 500000,
//	      "subtotal_rental": 1000000,
//	      "subtotal_deposit": 1000000
//	    }
//	  ],
//	  "customer": {
//	    "name": "Budi Santoso",
//	    "phone": "081234567890",
//	    "email": "budi@example.com",
//	    "delivery_address": "Jakarta Selatan",
//	    "notes": "Tolong kirim pagi hari"
//	  },
//	  "delivery": 50000,
//	  "discount": 0
//	}
type CreateBookingByCustomerRequest struct {
	StartDate    string                                 `json:"start_date"`    // Format: YYYY-MM-DD
	EndDate      string                                 `json:"end_date"`      // Format: YYYY-MM-DD
	DeliveryType string                                 `json:"delivery_type"` // "pickup" atau "delivery"
	Items        []CreateBookingItemByCustomerRequest   `json:"items"`
	Customer     CreateBookingCustomerByCustomerRequest `json:"customer"`
	Delivery     int                                    `json:"delivery"`
	Discount     int                                    `json:"discount"`
}

// CreateBookingItemByCustomerRequest adalah detail item dalam booking request
type CreateBookingItemByCustomerRequest struct {
	ItemID          string `json:"item_id"`
	Name            string `json:"name"`
	Quantity        int    `json:"quantity"`
	PricePerDay     int    `json:"price_per_day"`
	DepositPerUnit  int    `json:"deposit_per_unit"`
	SubtotalRental  int    `json:"subtotal_rental"`
	SubtotalDeposit int    `json:"subtotal_deposit"`
}

// CreateBookingCustomerByCustomerRequest adalah data kontak penerima booking
type CreateBookingCustomerByCustomerRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"delivery_address"`
	Notes   string `json:"notes"`
}

// ===================================================================
// RESPONSE DTO - CUSTOMER
// ===================================================================

// BookingDetailByCustomerResponse adalah response detail lengkap booking dari sisi customer
// Endpoint: GET /customer/booking/{id}
//
// Contoh JSON:
//
//	{
//	  "booking": { ... },
//	  "items": [ ... ],
//	  "customer": { ... }
//	}
type BookingDetailByCustomerResponse struct {
	Booking  BookingInfoResponse   `json:"booking"`
	Items    []BookingItemResponse `json:"items"`
	Customer CustomerInfoResponse  `json:"customer"`
}

// BookingListByCustomerResponse adalah response untuk list booking customer
// Endpoint: GET /customer/booking/me
type BookingListByCustomerResponse struct {
	BookingID  string    `json:"booking_id" db:"booking_id"`
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	Total      int64     `json:"total" db:"total"`
	Status     string    `json:"status" db:"status"`
	ItemNames  string    `json:"item_name" db:"item_names"`
	TotalItems int       `json:"total_item" db:"total_items"`
}

// ===================================================================
// RESPONSE DTO - HOSTER
// ===================================================================

// BookingDetailByHosterResponse adalah response detail lengkap booking dari sisi hoster
// Endpoint: GET /hoster/booking/{id}
//
// Sama seperti customer, tapi hoster juga bisa lihat data KTP customer
type BookingDetailByHosterResponse struct {
	Booking  BookingInfoResponse   `json:"booking"`
	Items    []BookingItemResponse `json:"items"`
	Customer CustomerInfoResponse  `json:"customer"` // Termasuk KTP jika sudah verified
}

// BookingListByHosterResponse adalah response untuk list booking hoster
// Endpoint: GET /hoster/booking
type BookingListByHosterResponse struct {
	BookingID    string    `json:"booking_id" db:"booking_id"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	Total        float64   `json:"total" db:"total"`
	Status       string    `json:"status" db:"status"`
	ItemName     string    `json:"item_name" db:"item_name"`
	TotalItem    int       `json:"total_item" db:"total_item"`
	CustomerName string    `json:"customer_name" db:"customer_name"`
}

// CustomerListByHosterResponse adalah response untuk list customer yang pernah booking
// Endpoint: GET /hoster/booking/customers
type CustomerListByHosterResponse = CustomerInfoResponse

// ===================================================================
// SHARED RESPONSE DTO (dipakai customer & hoster)
// ===================================================================

// BookingInfoResponse berisi informasi lengkap tentang booking header
// Digunakan untuk customer dan hoster detail view
type BookingInfoResponse struct {
	ID                   string     `json:"id"`
	HosterID             string     `json:"hoster_id,omitempty"`
	UserID               string     `json:"user_id,omitempty"`
	IdentityID           *string    `json:"identity_id,omitempty"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              time.Time  `json:"end_date"`
	TotalDays            int        `json:"total_days"`
	DeliveryType         string     `json:"delivery_type"`
	Rental               int        `json:"rental"`
	Deposit              int        `json:"deposit"`
	Discount             int        `json:"discount"`
	Total                int        `json:"total"`
	Outstanding          int        `json:"outstanding"`
	Status               string     `json:"status"`
	LockedUntil          *time.Time `json:"locked_until,omitempty"`
	TimeRemainingMinutes int        `json:"time_remaining_minutes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// BookingItemResponse berisi informasi item dalam booking
// Digunakan untuk customer dan hoster view
type BookingItemResponse struct {
	ID              string `json:"id" db:"id"`
	BookingID       string `json:"booking_id" db:"booking_id"`
	ItemID          string `json:"item_id" db:"item_id"`
	Name            string `json:"name" db:"name"`
	Quantity        int    `json:"quantity" db:"quantity"`
	PricePerDay     int    `json:"price_per_day" db:"price_per_day"`
	DepositPerUnit  int    `json:"deposit_per_unit" db:"deposit_per_unit"`
	SubtotalRental  int    `json:"subtotal_rental" db:"subtotal_rental"`
	SubtotalDeposit int    `json:"subtotal_deposit" db:"subtotal_deposit"`
}

// CustomerInfoResponse berisi informasi dasar customer yang aman untuk ditampilkan
// Digunakan di berbagai tempat: booking detail, customer list, dll
type CustomerInfoResponse struct {
	ID          string     `json:"id" db:"id"`
	FullName    string     `json:"full_name" db:"full_name"`
	Email       string     `json:"email" db:"email"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	KTPID       string     `json:"ktp_id,omitempty" db:"ktp_id"`
	KTPPhoto    string     `json:"ktp_photo,omitempty" db:"ktp_photo"`
	Status      string     `json:"status,omitempty" db:"status"` // Status KTP: pending, approved, rejected
	Reason      string     `json:"reason,omitempty" db:"reason"` // Alasan approve/reject KTP
	UploadedAt  *time.Time `json:"uploaded_at,omitempty" db:"uploaded_at"`
}
