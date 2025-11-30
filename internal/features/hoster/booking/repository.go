package booking

import (
	"database/sql"
	"log"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"

	"github.com/jmoiron/sqlx"
)

/*
HosterBookingRepository mendefinisikan operasi database khusus perspektif hoster.
Digunakan untuk dashboard dan detail booking yang dimiliki hoster.
*/
type HosterBookingRepository interface {
	GetListBookings(hosterID string) ([]dto.BookingListByHosterResponse, error)
	// GetCustomerList returns distinct customers (snapshot + enrichment) who made bookings
	// against the given hoster. Returned rows are unique by booking.user_id.
	GetCustomerList(hosterID string) ([]dto.CustomerListByHosterResponse, error)
	GetBookingDetail(bookingID string) (*dto.BookingDetailByHosterResponse, error)
}

/*
hosterBookingRepository adalah implementasi repository untuk hoster-facing booking.
*/
type hosterBookingRepository struct {
	db *sqlx.DB
}

/*
NewHosterBookingRepository membuat instance repository dengan koneksi database.

Output:
- HosterBookingRepository siap digunakan
*/
func NewHosterBookingRepository(db *sqlx.DB) HosterBookingRepository {
	return &hosterBookingRepository{db: db}
}

/*
GetListBookings mengambil ringkasan semua booking yang dimiliki hoster.

Alur kerja:
1. Query JOIN booking + booking_item dengan filter hoster_id
2. Agregasi item names dan total quantity
3. Konversi string_agg → slice string untuk kemudahan di frontend

Output sukses:
- ([]dto.BookingListByHosterResponse, nil)
Output error:
- (nil, error) → query gagal
*/
func (r *hosterBookingRepository) GetListBookings(hosterID string) ([]dto.BookingListByHosterResponse, error) {
	query := `
		SELECT
			b.id AS booking_id,
			b.start_date::timestamptz AS start_date,
			b.end_date::timestamptz AS end_date,
			b.total,
			b.status,
			COALESCE(items.item_name, '') AS item_name,
			COALESCE(items.total_item, 0) AS total_item,
			COALESCE(NULLIF(bc.name, ''), c.full_name, '') AS customer_name
		FROM booking b
		LEFT JOIN (
			SELECT booking_id, string_agg(name, ', ' ORDER BY name) AS item_name, SUM(quantity) AS total_item
			FROM booking_item
			GROUP BY booking_id
		) items ON items.booking_id = b.id
		LEFT JOIN booking_customer bc ON bc.booking_id = b.id
		LEFT JOIN customer c ON c.id = b.user_id
		WHERE b.hoster_id = $1
		ORDER BY b.created_at DESC
	`

	var rows []dto.BookingListByHosterResponse
	if err := r.db.Select(&rows, query, hosterID); err != nil {
		log.Printf("GetListBookings(hoster): db error hoster=%s err=%v", hosterID, err)
		return nil, err
	}

	return rows, nil
}

// GetCustomerList mengambil daftar pelanggan yang melakukan pemesanan pada hoster tertentu.
// Query akan memilih snapshot dari booking_customer jika tersedia, dan akan mencoba
// menambahkan data KTP/identity jika ada untuk user yang sama. Hasil dikembalikan
// unik per user (booking.user_id).
func (r *hosterBookingRepository) GetCustomerList(hosterID string) ([]dto.CustomerListByHosterResponse, error) {
	query := `
		SELECT DISTINCT ON (b.user_id)
			-- Prioritize booking_customer snapshot ID when available
			COALESCE(bc.id::text, c.id::text, '') AS id,
			COALESCE(bc.name, c.full_name, '') AS full_name,
			COALESCE(bc.email, c.email, '') AS email,
			COALESCE(bc.phone, c.phone_number, '') AS phone_number,
			-- Use the booking's identity_id (snapshot), NOT latest identity
			COALESCE(b.identity_id::text, '') AS ktp_id,
			COALESCE(i_by_booking.created_at, bc.created_at) AS uploaded_at,
			COALESCE(i_by_booking.ktp_url, '') AS ktp_photo,
			COALESCE(i_by_booking.status, '') AS status,
			COALESCE(i_by_booking.reason, '') AS reason
		FROM booking b
		LEFT JOIN booking_customer bc ON b.id = bc.booking_id
		LEFT JOIN customer c ON b.user_id = c.id
		-- Join with the booking's specific identity snapshot (booking.identity_id)
		LEFT JOIN identity i_by_booking ON b.identity_id = i_by_booking.id
		WHERE b.hoster_id = $1
		-- Order to get the most complete booking per user:
		-- 1. Prefer bookings that have both booking_customer AND identity_id
		-- 2. Then prefer bookings with identity_id
		-- 3. Then prefer most recent booking
		ORDER BY 
			b.user_id,
			(bc.id IS NOT NULL AND b.identity_id IS NOT NULL) DESC,
			(b.identity_id IS NOT NULL) DESC,
			b.created_at DESC
	`

	var rows []dto.CustomerListByHosterResponse
	if err := r.db.Select(&rows, query, hosterID); err != nil {
		log.Printf("GetCustomerList(hoster): db error hoster=%s err=%v", hosterID, err)
		return nil, err
	}

	return rows, nil
}

/*
GetBookingDetail mengambil detail lengkap satu booking dari perspektif hoster.

Alur kerja:
1. Ambil header booking
2. Hitung sisa waktu locked (jika ada)
3. Ambil list item
4. Ambil snapshot customer (booking_customer) → fallback ke identity/customer jika kosong
5. Mapping semua data ke DTO hoster

Output sukses:
- (*dto.BookingDetailByHosterResponse, nil)
Output error:
- (nil, error) → salah satu query gagal / data tidak ditemukan
*/
func (r *hosterBookingRepository) GetBookingDetail(bookingID string) (*dto.BookingDetailByHosterResponse, error) {
	var b domain.Booking
	queryBooking := `
		SELECT id, hoster_id, locked_until, start_date, end_date, total_days, delivery_type,
			   rental, deposit, discount, total, outstanding, user_id, identity_id, status,
			   created_at, updated_at
		FROM booking
		WHERE id = $1
	`
	if err := r.db.Get(&b, queryBooking, bookingID); err != nil {
		log.Printf("GetBookingDetail(hoster): error getting booking %s: %v", bookingID, err)
		return nil, err
	}

	// Hitung sisa waktu locked (jika masih aktif)
	now := time.Now()
	if !b.LockedUntil.IsZero() && b.LockedUntil.After(now) {
		b.TimeRemainingMinutes = int(b.LockedUntil.Sub(now).Minutes())
	} else {
		b.TimeRemainingMinutes = 0
	}

	// Ambil items
	var items []dto.BookingItemResponse
	queryItems := `
		SELECT id, booking_id, item_id, name, quantity,
			   price_per_day, deposit_per_unit, subtotal_rental, subtotal_deposit
		FROM booking_item
		WHERE booking_id = $1
	`
	if err := r.db.Select(&items, queryItems, bookingID); err != nil {
		log.Printf("GetBookingDetail(hoster): error getting items for %s: %v", bookingID, err)
		return nil, err
	}

	// Ambil snapshot customer (booking_customer)
	var cust dto.CustomerInfoResponse
	// booking_customer only contains snapshot fields: id, booking_id, name, phone, email, address, notes, created_at
	// Query only existing columns and map to DTO fields.
	queryCustPrimary := `
		SELECT id, name AS full_name, email, phone AS phone_number, created_at AS uploaded_at
		FROM booking_customer
		WHERE booking_id = $1
		LIMIT 1
	`
	err := r.db.Get(&cust, queryCustPrimary, bookingID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("GetBookingDetail(hoster): error getting booking_customer for %s: %v", bookingID, err)
		return nil, err
	}
	if err == nil {
		// try to attach identity/KTP data as enrichment
		var identity domain.Identity
		if idErr := r.db.Get(&identity, `
			SELECT id, user_id, ktp_url, verified, status, COALESCE(reason,'') AS reason, created_at, updated_at
			FROM identity WHERE user_id = $1 ORDER BY verified_at DESC NULLS LAST, created_at DESC LIMIT 1
		`, b.UserID); idErr == nil && identity.ID != "" {
			cust.KTPID = identity.ID
			cust.KTPPhoto = identity.KTPURL
			cust.Status = identity.Status
			cust.Reason = identity.Reason
			if !identity.CreatedAt.IsZero() {
				t := identity.CreatedAt
				cust.UploadedAt = &t
			}
		}
	}

	// Jika tetap tidak ada snapshot → fallback ke customer + identity
	if err == sql.ErrNoRows {
		var c domain.Customer
		fallbackErr := r.db.Get(&c, `
			SELECT id, full_name, phone_number, email
			FROM customer WHERE id = $1 LIMIT 1
		`, b.UserID)

		if fallbackErr == nil {
			cust = dto.CustomerInfoResponse{
				ID:          c.ID,
				FullName:    c.FullName,
				Email:       c.Email,
				PhoneNumber: c.PhoneNumber,
			}
		}

		// Tambahkan data KTP dari identity jika ada
		var identity domain.Identity
		if err := r.db.Get(&identity, `
			SELECT id, user_id, ktp_url, status, COALESCE(reason,'') AS reason, created_at
			FROM identity WHERE user_id = $1
			ORDER BY verified_at DESC NULLS LAST, created_at DESC LIMIT 1
		`, b.UserID); err == nil && identity.ID != "" {
			cust.KTPID = identity.ID
			cust.KTPPhoto = identity.KTPURL
			cust.Status = identity.Status
			cust.Reason = identity.Reason
			if !identity.CreatedAt.IsZero() {
				t := identity.CreatedAt
				cust.UploadedAt = &t
			}
		}
	}

	// Mapping final ke DTO
	detail := &dto.BookingDetailByHosterResponse{
		Booking: dto.BookingInfoResponse{
			ID:                   b.ID,
			HosterID:             b.HosterID,
			UserID:               b.UserID,
			IdentityID:           b.IdentityID,
			StartDate:            b.StartDate,
			EndDate:              b.EndDate,
			TotalDays:            b.TotalDays,
			DeliveryType:         b.DeliveryType,
			Rental:               b.Rental,
			Deposit:              b.Deposit,
			Discount:             b.Discount,
			Total:                b.Total,
			Outstanding:          b.Outstanding,
			Status:               b.Status,
			LockedUntil:          &b.LockedUntil,
			TimeRemainingMinutes: b.TimeRemainingMinutes,
			CreatedAt:            b.CreatedAt,
			UpdatedAt:            b.UpdatedAt,
		},
		Items: items,
		Customer: dto.CustomerInfoResponse{
			ID:          cust.ID,
			FullName:    cust.FullName,
			Email:       cust.Email,
			PhoneNumber: cust.PhoneNumber,
			KTPID:       cust.KTPID,
			KTPPhoto:    cust.KTPPhoto,
			Status:      cust.Status,
			Reason:      cust.Reason,
			UploadedAt:  cust.UploadedAt,
		},
	}

	log.Printf("GetBookingDetail(hoster): success booking=%s", bookingID)
	return detail, nil
}
