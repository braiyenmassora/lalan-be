package model

import "time"

/*
BookingModel
struct untuk data booking utama dengan field lengkap
*/
type BookingModel struct {
	ID                   string    `json:"id" db:"id"`
	Code                 string    `json:"code" db:"code"`
	HosterID             string    `json:"hoster_id" db:"hoster_id"`
	LockedUntil          time.Time `json:"locked_until" db:"locked_until"`
	TimeRemainingMinutes int       `json:"time_remaining_minutes" db:"-"`
	StartDate            string    `json:"start_date" db:"start_date"`
	EndDate              string    `json:"end_date" db:"end_date"`
	TotalDays            int       `json:"total_days" db:"total_days"`
	DeliveryType         string    `json:"delivery_type" db:"delivery_type"`
	Rental               int       `json:"rental" db:"rental"`
	Deposit              int       `json:"deposit" db:"deposit"`
	Discount             int       `json:"discount" db:"discount"`
	Total                int       `json:"total" db:"total"`
	Outstanding          int       `json:"outstanding" db:"outstanding"`
	UserID               string    `json:"user_id" db:"user_id"`
	IdentityID           *string   `json:"identity_id" db:"identity_id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

/*
BookingItem
struct untuk item dalam booking
*/
type BookingItem struct {
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

/*
BookingCustomer
struct untuk snapshot customer dalam booking
*/
type BookingCustomer struct {
	ID        string `json:"id" db:"id"`
	BookingID string `json:"booking_id" db:"booking_id"`
	Name      string `json:"name" db:"name"`
	Phone     string `json:"phone" db:"phone"`
	Email     string `json:"email" db:"email"`
	Address   string `json:"address" db:"address"`
	Notes     string `json:"notes" db:"notes"`
}

/*
BookingIdentity
struct untuk snapshot identitas dalam booking
*/
type BookingIdentity struct {
	ID              string    `json:"id" db:"id"`
	BookingID       string    `json:"booking_id" db:"booking_id"`
	Uploaded        bool      `json:"uploaded" db:"uploaded"`
	Status          string    `json:"status" db:"status"`
	Reason          *string   `json:"reason" db:"reason"`
	ReuploadAllowed bool      `json:"reupload_allowed" db:"reupload_allowed"`
	EstimatedTime   string    `json:"estimated_time" db:"estimated_time"`
	StatusCheckURL  string    `json:"status_check_url" db:"status_check_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

/*
BookingDetailDTO
struct DTO untuk detail booking lengkap
*/
type BookingDetailDTO struct {
	Booking  BookingModel    `json:"booking"`
	Items    []BookingItem   `json:"items"`
	Customer BookingCustomer `json:"customer"`
	Identity BookingIdentity `json:"identity"`
}

/*
RentalPeriod
struct untuk periode rental booking
*/
type RentalPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	TotalDays int    `json:"total_days"`
}

/*
BookingPrice
struct untuk harga booking
*/
type BookingPrice struct {
	Rental      int `json:"rental"`
	Deposit     int `json:"deposit"`
	Discount    int `json:"discount"`
	Total       int `json:"total"`
	Outstanding int `json:"outstanding"`
}

/*
BookingListDTO
struct DTO untuk daftar booking
*/
type BookingListDTO struct {
	Code           string `json:"code" db:"code"`
	StartDate      string `json:"start_date" db:"start_date"`
	EndDate        string `json:"end_date" db:"end_date"`
	Total          int    `json:"total" db:"total"`
	IdentityStatus string `json:"identity_status" db:"identity_status"`
	ItemSummary    string `json:"item_summary" db:"item_summary"`
	Quantity       int    `json:"quantity" db:"quantity"`
}

type BookingListCustomer struct {
	BookingID      string `json:"booking_id" db:"booking_id"`
	Code           string `json:"code" db:"code"`
	CustomerID     string `json:"customer_id" db:"customer_id"`
	CustomerName   string `json:"customer_name" db:"customer_name"`
	StartDate      string `json:"start_date" db:"start_date"`
	EndDate        string `json:"end_date" db:"end_date"`
	DurationDays   int    `json:"duration_days" db:"duration_days"`
	Total          int    `json:"total" db:"total"`
	IdentityStatus string `json:"identity_status" db:"identity_status"`
	ItemSummary    string `json:"item_summary" db:"item_summary"`
	Quantity       int    `json:"quantity" db:"quantity"`
}

type BookingDetailCustomer struct {
	BookingID      string `json:"booking_id" db:"booking_id"`
	Code           string `json:"code" db:"code"`
	CustomerID     string `json:"customer_id" db:"customer_id"`
	CustomerName   string `json:"customer_name" db:"customer_name"`
	CustomerEmail  string `json:"customer_email" db:"customer_email"`
	CustomerPhone  string `json:"customer_phone" db:"customer_phone"`
	StartDate      string `json:"start_date" db:"start_date"`
	EndDate        string `json:"end_date" db:"end_date"`
	DurationDays   int    `json:"duration_days" db:"duration_days"`
	Total          int    `json:"total" db:"total"`
	IdentityStatus string `json:"identity_status" db:"identity_status"` // indikator KTP
	ItemSummary    string `json:"item_summary" db:"item_summary"`
	Quantity       int    `json:"quantity" db:"quantity"`

	// field tambahan yang bukan dari DB
	Products      string `json:"products" db:"-"`
	ItemCount     int    `json:"item_count" db:"-"`
	TotalQuantity int    `json:"total_quantity" db:"-"`
	DeliveryType  string `json:"delivery_type" db:"delivery_type"`
	CustomerAddr  string `json:"customer_address" db:"customer_address"`
	CustomerNotes string `json:"customer_notes" db:"customer_notes"`
}
