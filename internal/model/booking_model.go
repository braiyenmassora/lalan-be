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
	Delivery             int       `json:"delivery" db:"delivery"`
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
	ID              string `json:"id" db:"id"`
	BookingID       string `json:"booking_id" db:"booking_id"`
	Name            string `json:"name" db:"name"`
	Phone           string `json:"phone" db:"phone"`
	Email           string `json:"email" db:"email"`
	DeliveryAddress string `json:"delivery_address" db:"delivery_address"`
	Notes           string `json:"notes" db:"notes"`
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
	RejectionReason *string   `json:"rejection_reason" db:"rejection_reason"`
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
	Delivery    int `json:"delivery"`
	Discount    int `json:"discount"`
	Total       int `json:"total"`
	Outstanding int `json:"outstanding"`
}

/*
BookingListDTO
struct DTO untuk daftar booking
*/
type BookingListDTO struct {
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ItemName  string    `json:"item_name" db:"item_name"`
	Quantity  int       `json:"quantity" db:"quantity"`
	KtpStatus string    `json:"ktp_status" db:"ktp_status"`
	Total     int       `json:"total" db:"total"`
}
