package booking

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"

	"github.com/jmoiron/sqlx"
)

/*
BookingRepository adalah kontrak (interface) untuk semua operasi database
terkait fitur booking. Implementasinya hanya boleh berisi query SQL dan
transaksi — tidak boleh ada logika bisnis atau validasi domain.
*/
type BookingRepository interface {
	CreateBooking(booking *domain.Booking, items []domain.BookingItem, customer domain.BookingCustomer) (*dto.BookingDetailByCustomerResponse, error)
	GetListBookings(userID string) ([]dto.BookingListByCustomerResponse, error)
	GetBookingDetail(bookingID string) (*dto.BookingDetailByCustomerResponse, error)
	GetIdentityByUserID(userID string) (*domain.Identity, error)
	GetHosterIDByItemID(itemID string) (string, error)
}

/*
bookingRepository adalah implementasi konkret dari BookingRepository.
Menyimpan koneksi *sqlx.DB yang digunakan untuk semua query.
*/
type bookingRepository struct {
	db *sqlx.DB
}

/*
NewBookingRepository membuat instance repository yang siap digunakan.
Dependency injection dilakukan di sini agar mudah di-mock saat unit test.

Output:
- Implementasi BookingRepository yang terhubung ke database.
*/
func NewBookingRepository(db *sqlx.DB) BookingRepository {
	return &bookingRepository{db: db}
}

/*
CreateBooking menyimpan booking baru beserta item dan data customer ke database.
Menggunakan database transaction untuk menjamin atomicity (all-or-nothing).

Alur kerja:
1. Validasi KTP user (hanya cek keberadaan, bukan business rule)
2. Mulai transaction
3. Insert header booking → booking_item → booking_customer
4. Commit transaction
5. Query ulang detail booking untuk dikembalikan ke service

Output sukses:
- *dto.BookingDetailByCustomerResponse (data lengkap booking yang baru dibuat)
Output error:
- error validasi KTP → "silakan upload ktp terlebih dahulu"
- error DB → langsung diteruskan ke service (akan jadi 500 atau 400 sesuai konteks)
*/
func (r *bookingRepository) CreateBooking(booking *domain.Booking, items []domain.BookingItem, customer domain.BookingCustomer) (*dto.BookingDetailByCustomerResponse, error) {
	// 1. Validasi Identity (KTP)
	identity, err := r.GetIdentityByUserID(booking.UserID)
	if err != nil {
		log.Printf("CreateBooking: error checking identity for user %s: %v", booking.UserID, err)
		return nil, err
	}
	if identity == nil {
		log.Printf("CreateBooking: no identity found for user %s", booking.UserID)
		return nil, fmt.Errorf("silakan upload ktp terlebih dahulu")
	}

	// 2. Mulai Transaction
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("CreateBooking: error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	// 3. Insert Booking Header
	queryBooking := `
		INSERT INTO booking (
			id, hoster_id, locked_until, start_date, end_date, total_days,
			delivery_type, rental, deposit, discount, total, outstanding,
			user_id, identity_id, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err = tx.Exec(queryBooking,
		booking.ID, booking.HosterID, booking.LockedUntil,
		booking.StartDate, booking.EndDate, booking.TotalDays,
		booking.DeliveryType, booking.Rental, booking.Deposit,
		booking.Discount, booking.Total, booking.Outstanding,
		booking.UserID, booking.IdentityID, booking.Status,
	)
	if err != nil {
		log.Printf("CreateBooking: error inserting booking header: %v", err)
		return nil, err
	}

	// 4. Insert Booking Items
	queryItem := `
		INSERT INTO booking_item (
			id, booking_id, item_id, name, quantity,
			price_per_day, deposit_per_unit, subtotal_rental, subtotal_deposit
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	for i, item := range items {
		_, err = tx.Exec(queryItem,
			item.ID, item.BookingID, item.ItemID, item.Name, item.Quantity,
			item.PricePerDay, item.DepositPerUnit, item.SubtotalRental, item.SubtotalDeposit,
		)
		if err != nil {
			log.Printf("CreateBooking: error inserting booking_item index %d: %v", i, err)
			return nil, err
		}
	}

	// 5. Insert Booking Customer
	queryCustomer := `
		INSERT INTO booking_customer (
			id, booking_id, name, phone, email, address, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(queryCustomer,
		customer.ID, customer.BookingID, customer.Name, customer.Phone,
		customer.Email, customer.Address, customer.Notes,
	)
	if err != nil {
		log.Printf("CreateBooking: error inserting booking_customer: %v", err)
		return nil, err
	}

	// 6. Commit Transaction
	if err = tx.Commit(); err != nil {
		log.Printf("CreateBooking: error committing transaction: %v", err)
		return nil, err
	}
	log.Printf("CreateBooking: booking %s created successfully", booking.ID)

	// 7. Return Detail Booking
	detail, err := r.GetBookingDetail(booking.ID)
	if err != nil {
		log.Printf("CreateBooking: error retrieving created booking detail: %v", err)
		return nil, err
	}
	return detail, nil
}

/*
GetIdentityByUserID mengambil data KTP terverifikasi terakhir milik user.
Digunakan hanya untuk pengecekan keberadaan KTP sebelum create booking.

Output sukses:
- *domain.Identity jika ada
- nil (bukan error) jika user belum upload KTP
Output error:
- error hanya jika query database gagal
*/
func (r *bookingRepository) GetIdentityByUserID(userID string) (*domain.Identity, error) {
	var identity domain.Identity
	// Pilih hanya identity yang sudah ter-verified (verified = true) dan ambil
	// yang memiliki verified_at terbaru. Ini memastikan saat membuat booking
	// kita selalu memakai identity yang sah dan paling baru berdasarkan waktu
	// verifikasi.
	query := `
		SELECT
			id, user_id, ktp_url, verified, status,
			COALESCE(reason, '') AS reason,
			verified_at, created_at, updated_at
		FROM identity
		WHERE user_id = $1 AND verified = true
		ORDER BY verified_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`
	err := r.db.Get(&identity, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("GetIdentityByUserID: database error for user %s: %v", userID, err)
		return nil, err
	}
	return &identity, nil
}

/*
GetBookingsByUserID mengambil daftar ringkas semua booking milik user.
Hasil berupa summary (total item, nama item digabung, dll).

Output sukses:
- []dto.BookingListByCustomerResponse (slice kosong jika user belum punya booking)
Output error:
- error hanya jika query gagal (akan jadi 500 di handler)
*/
func (r *bookingRepository) GetListBookings(userID string) ([]dto.BookingListByCustomerResponse, error) {
	query := `
		SELECT
			b.id AS booking_id,
			b.start_date::timestamptz AS start_date,
			b.end_date::timestamptz AS end_date,
			b.total,
			b.status,
			COALESCE(string_agg(bi.name, ', '), '') AS item_names,
			COALESCE(SUM(bi.quantity), 0) AS total_items
		FROM booking b
		LEFT JOIN booking_item bi ON b.id = bi.booking_id
		WHERE b.user_id = $1
		GROUP BY b.id, b.start_date, b.end_date, b.total, b.status
		ORDER BY b.created_at DESC
	`
	var bookings []dto.BookingListByCustomerResponse
	err := r.db.Select(&bookings, query, userID)
	if err != nil {
		log.Printf("GetListBookings: database error for user %s: %v", userID, err)
		return nil, err
	}
	return bookings, nil
}

/*
GetHosterIDByItemID mencari user_id (hoster) yang memiliki item tertentu.
Digunakan saat pembuatan booking untuk mengisi field hoster_id.

Output sukses:
- string hoster_id
Output error:
- error jika item tidak ditemukan atau query gagal
*/
func (r *bookingRepository) GetHosterIDByItemID(itemID string) (string, error) {
	var hosterID string
	query := `SELECT hoster_id FROM item WHERE id = $1`
	err := r.db.Get(&hosterID, query, itemID)
	if err != nil {
		log.Printf("GetHosterIDByItemID: error querying item %s: %v", itemID, err)
		return "", err
	}
	return hosterID, nil
}

/*
GetBookingDetail mengambil data lengkap satu booking termasuk:
- Header booking + waktu tersisa pembayaran
- Semua booking_item
- Data customer
- Status KTP terakhir user

Alur kerja:
1. Query header booking
2. Hitung sisa menit locked_until
3. Query items
4. Query data customer
5. Query data KTP (opsional)
6. Bangun DTO lengkap

Output sukses:
- *dto.BookingDetailByCustomerResponse (semua field terisi)
Output error:
- error jika booking tidak ditemukan atau query gagal (akan jadi 404/500 di service/handler)
*/
func (r *bookingRepository) GetBookingDetail(bookingID string) (*dto.BookingDetailByCustomerResponse, error) {
	// 1. Get Booking Header
	var booking domain.Booking
	queryBooking := `
		SELECT id, hoster_id, locked_until, start_date, end_date, total_days, delivery_type,
		       rental, deposit, discount, total, outstanding, user_id, identity_id, status,
		       created_at, updated_at
		FROM booking WHERE id = $1
	`
	err := r.db.Get(&booking, queryBooking, bookingID)
	if err != nil {
		log.Printf("GetBookingDetail: error querying booking header %s: %v", bookingID, err)
		return nil, err
	}

	// Hitung waktu tersisa
	now := time.Now()
	if booking.LockedUntil.After(now) {
		booking.TimeRemainingMinutes = int(booking.LockedUntil.Sub(now).Minutes())
	} else {
		booking.TimeRemainingMinutes = 0
	}

	// 2. Get Booking Items
	var items []domain.BookingItem
	queryItems := `
		SELECT id, booking_id, item_id, name, quantity, price_per_day, deposit_per_unit,
		       subtotal_rental, subtotal_deposit
		FROM booking_item WHERE booking_id = $1
	`
	err = r.db.Select(&items, queryItems, bookingID)
	if err != nil {
		log.Printf("GetBookingDetail: error querying items for booking %s: %v", bookingID, err)
		return nil, err
	}

	// 3. Get Customer Info
	var customer domain.Customer
	queryCustomer := `
		SELECT id, full_name, phone_number, email, profile_photo, address, email_verified, created_at, updated_at
		FROM customer WHERE id = $1
	`
	err = r.db.Get(&customer, queryCustomer, booking.UserID)
	if err != nil {
		log.Printf("GetBookingDetail: error querying customer for user %s: %v", booking.UserID, err)
		return nil, err
	}

	// 4. Get Identity Info (opsional)
	var identity domain.Identity
	queryIdentity := `
		SELECT
			id, user_id, ktp_url, verified, status,
			COALESCE(reason, '') AS reason,
			verified_at, created_at, updated_at
		FROM identity
		WHERE user_id = $1
		ORDER BY verified_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`
	err = r.db.Get(&identity, queryIdentity, booking.UserID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("GetBookingDetail: error querying identity for user %s: %v", booking.UserID, err)
		return nil, err
	}

	// 5. Build Response DTO - MAPPING MANUAL

	// Mapping Customer
	customerResponse := dto.CustomerInfoResponse{
		ID:          customer.ID,
		FullName:    customer.FullName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
	}

	if err != sql.ErrNoRows && identity.ID != "" {
		customerResponse.KTPID = identity.ID
		customerResponse.KTPPhoto = identity.KTPURL
		customerResponse.Status = identity.Status
		customerResponse.Reason = identity.Reason
		customerResponse.UploadedAt = &identity.CreatedAt
	}

	// Mapping Booking Header
	bookingResponse := dto.BookingInfoResponse{
		ID:                   booking.ID,
		HosterID:             booking.HosterID,
		UserID:               booking.UserID,
		IdentityID:           booking.IdentityID,
		StartDate:            booking.StartDate,
		EndDate:              booking.EndDate,
		TotalDays:            booking.TotalDays,
		DeliveryType:         booking.DeliveryType,
		Rental:               booking.Rental,
		Deposit:              booking.Deposit,
		Discount:             booking.Discount,
		Total:                booking.Total,
		Outstanding:          booking.Outstanding,
		Status:               booking.Status,
		LockedUntil:          &booking.LockedUntil,
		TimeRemainingMinutes: booking.TimeRemainingMinutes,
		CreatedAt:            booking.CreatedAt,
		UpdatedAt:            booking.UpdatedAt,
	}

	// Mapping Items
	itemsResponse := make([]dto.BookingItemResponse, len(items))
	for i, item := range items {
		itemsResponse[i] = dto.BookingItemResponse{
			ID:              item.ID,
			BookingID:       item.BookingID,
			ItemID:          item.ItemID,
			Name:            item.Name,
			Quantity:        item.Quantity,
			PricePerDay:     item.PricePerDay,
			DepositPerUnit:  item.DepositPerUnit,
			SubtotalRental:  item.SubtotalRental,
			SubtotalDeposit: item.SubtotalDeposit,
		}
	}

	detail := &dto.BookingDetailByCustomerResponse{
		Booking:  bookingResponse,
		Items:    itemsResponse,
		Customer: customerResponse,
	}

	log.Printf("GetBookingDetail: successfully retrieved detail for booking %s", bookingID)
	return detail, nil
}
