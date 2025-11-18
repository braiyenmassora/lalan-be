package customer

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

type customerRespository struct {
	db *sqlx.DB
}

func (r *customerRespository) CreateCustomer(customer *model.CustomerModel) error {
	query := `
		INSERT INTO customer (
			full_name,
			address,
			phone_number,
			email,
			password_hash,
			profile_photo,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, customer.FullName, customer.Address, customer.PhoneNumber, customer.Email, customer.PasswordHash, customer.ProfilePhoto, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
	log.Printf("CreateCustomer: inserted customer with email %s, ID %s", customer.Email, customer.ID)
	return err
}

func (r *customerRespository) FindByEmailCustomerForLogin(email string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	query := `
		SELECT
			id,
			full_name,
			phone_number,
			email,
			password_hash,
			profile_photo,
			address,
			created_at,
			updated_at
		FROM customer
		WHERE email = $1
	`
	err := r.db.Get(&customer, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindByEmailCustomerForLogin: no customer found for email %s", email)
			return nil, nil
		}
		log.Printf("FindByEmailCustomerForLogin: error querying email %s: %v", email, err)
		return nil, err
	}
	log.Printf("FindByEmailCustomerForLogin: found customer for email %s", email)
	return &customer, nil
}

func (r *customerRespository) GetDetailCustomer(id string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	query := `
        SELECT 
            id,
            full_name,
            address,
            phone_number,
            email,
            password_hash,
            profile_photo,
            created_at,
            updated_at
        FROM customer
        WHERE id = $1
    `
	err := r.db.Get(&customer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetDetailCustomer: no customer found for id %s", id)
			return nil, nil
		}
		log.Printf("GetDetailCustomer: error for id %s: %v", id, err)
		return nil, err
	}
	log.Printf("GetDetailCustomer: found customer id %s", id)
	return &customer, nil
}

func (r *customerRespository) UpdateCustomer(customer *model.CustomerModel) error {
	// Pembatasan: Hanya izinkan perubahan pada full_name, phone_number, profile_photo, dan address.
	// Pastikan di layer service bahwa hanya field ini yang diisi, dan email/password_hash tidak diubah.
	query := `
        UPDATE customer
        SET
            full_name = $1,
            phone_number = $2,
            profile_photo = $3,
            address = $4,
            updated_at = $5
        WHERE id = $6
    `
	_, err := r.db.Exec(query, customer.FullName, customer.PhoneNumber, customer.ProfilePhoto, customer.Address, customer.UpdatedAt, customer.ID)
	if err != nil {
		log.Printf("UpdateCustomer: error updating customer ID %s: %v", customer.ID, err)
	}
	return err
}

func (r *customerRespository) DeleteCustomer(id string) error {
	query := `DELETE FROM customer WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteCustomer: error deleting customer ID %s: %v", id, err)
		return err
	}
	log.Printf("DeleteCustomer: deleted customer ID %s", id)
	return nil
}

func (r *customerRespository) CreateIdentity(identity *model.IdentityModel) error {
	query := `
        INSERT INTO identity (
            user_id,
            ktp_url,
            verified,
            status,
            rejected_reason,
            verified_at,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `
	err := r.db.QueryRow(query, identity.UserID, identity.KTPURL, identity.Verified, identity.Status, identity.RejectedReason, identity.VerifiedAt, identity.CreatedAt, identity.UpdatedAt).Scan(&identity.ID, &identity.CreatedAt, &identity.UpdatedAt)
	if err != nil {
		log.Printf("CreateIdentity: error inserting identity for user %s: %v", identity.UserID, err)
		return err
	}
	log.Printf("CreateIdentity: inserted identity ID %s for user %s", identity.ID, identity.UserID)
	return nil
}

func (r *customerRespository) CheckIdentityExists(userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM identity WHERE user_id = $1 AND status != 'rejected'`
	err := r.db.Get(&count, query, userID)
	if err != nil {
		log.Printf("CheckIdentityExists: error checking identity for user %s: %v", userID, err)
		return false, err
	}
	return count > 0, nil
}

func (r *customerRespository) GetIdentityByUserID(userID string) (*model.IdentityModel, error) {
	var identity model.IdentityModel
	query := `
        SELECT
            id,
            user_id,
            ktp_url,
            verified,
            status,
            rejected_reason,
            verified_at,
            created_at,
            updated_at
        FROM identity
        WHERE user_id = $1
    `
	err := r.db.Get(&identity, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetIdentityByUserID: no identity found for user %s", userID)
			return nil, nil
		}
		log.Printf("GetIdentityByUserID: error for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("GetIdentityByUserID: found identity for user %s", userID)
	return &identity, nil
}

func (r *customerRespository) CreateBooking(booking *model.BookingModel, items []model.BookingItem, customer model.BookingCustomer, identity model.BookingIdentity) error {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("CreateBooking: error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// Insert booking
	queryBooking := `
        INSERT INTO booking (
            id, code, locked_until,
            start_date, end_date, total_days,
            delivery_type,
            rental, deposit, delivery, discount, total, outstanding,
            user_id, identity_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `
	log.Printf("CreateBooking: executing query: %s", queryBooking)
	log.Printf("CreateBooking: booking params: %s, %s, %v, %s, %s, %d, %s, %d, %d, %d, %d, %d, %d, %s, %v",
		booking.ID, booking.Code, booking.LockedUntil, booking.StartDate, booking.EndDate, booking.TotalDays, booking.DeliveryType,
		booking.Rental, booking.Deposit, booking.Delivery, booking.Discount, booking.Total, booking.Outstanding, booking.UserID, booking.IdentityID)
	_, err = tx.Exec(queryBooking, booking.ID, booking.Code, booking.LockedUntil, booking.StartDate, booking.EndDate, booking.TotalDays, booking.DeliveryType,
		booking.Rental, booking.Deposit, booking.Delivery, booking.Discount, booking.Total, booking.Outstanding, booking.UserID, booking.IdentityID)
	if err != nil {
		log.Printf("CreateBooking: error inserting booking: %v", err)
		return err
	}

	// Insert items
	queryItem := `
        INSERT INTO booking_item (
            id, booking_id, item_id, name, quantity,
            price_per_day, deposit_per_unit, subtotal_rental, subtotal_deposit
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	for i, item := range items {
		log.Printf("CreateBooking: inserting item %d: %s, %s, %s, %s, %d, %d, %d, %d, %d", i, item.ID, item.BookingID, item.ItemID, item.Name, item.Quantity, item.PricePerDay, item.DepositPerUnit, item.SubtotalRental, item.SubtotalDeposit)
		_, err = tx.Exec(queryItem, item.ID, item.BookingID, item.ItemID, item.Name, item.Quantity, item.PricePerDay, item.DepositPerUnit, item.SubtotalRental, item.SubtotalDeposit)
		if err != nil {
			log.Printf("CreateBooking: error inserting item %d: %v", i, err)
			return err
		}
	}

	// Insert customer
	queryCustomer := `
        INSERT INTO booking_customer (
            id, booking_id, name, phone, email, delivery_address, notes
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	log.Printf("CreateBooking: inserting customer: %s, %s, %s, %s, %s, %s, %s", customer.ID, customer.BookingID, customer.Name, customer.Phone, customer.Email, customer.DeliveryAddress, customer.Notes)
	_, err = tx.Exec(queryCustomer, customer.ID, customer.BookingID, customer.Name, customer.Phone, customer.Email, customer.DeliveryAddress, customer.Notes)
	if err != nil {
		log.Printf("CreateBooking: error inserting customer: %v", err)
		return err
	}

	// Insert identity
	queryIdentity := `
        INSERT INTO booking_identity (
            id, booking_id, uploaded, status, rejection_reason,
            reupload_allowed, estimated_time, status_check_url
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	log.Printf("CreateBooking: inserting identity: %s, %s, %t, %s, %v, %t, %s, %s", identity.ID, identity.BookingID, identity.Uploaded, identity.Status, identity.RejectionReason, identity.ReuploadAllowed, identity.EstimatedTime, identity.StatusCheckURL)
	_, err = tx.Exec(queryIdentity, identity.ID, identity.BookingID, identity.Uploaded, identity.Status, identity.RejectionReason, identity.ReuploadAllowed, identity.EstimatedTime, identity.StatusCheckURL)
	if err != nil {
		log.Printf("CreateBooking: error inserting identity: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("CreateBooking: error committing transaction: %v", err)
		return err
	}

	log.Printf("CreateBooking: booking %s created successfully", booking.ID)
	return nil
}

func (r *customerRespository) GetBookingsByUserID(userID string) ([]model.BookingListDTO, error) {
	query := `SELECT b.id, b.code, b.start_date, b.end_date, b.total_days, b.delivery_type, b.total, b.outstanding, b.created_at, bc.name AS customer_name, bc.phone AS customer_phone, bc.email AS customer_email, bc.delivery_address FROM booking b LEFT JOIN booking_customer bc ON b.id = bc.booking_id WHERE b.user_id = $1 ORDER BY b.created_at DESC`
	var bookings []model.BookingListDTO
	err := r.db.Select(&bookings, query, userID)
	if err != nil {
		log.Printf("GetBookingsByUserID: error for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("GetBookingsByUserID: found %d bookings for user %s", len(bookings), userID)
	return bookings, nil
}

func (r *customerRespository) GetListBookings() ([]model.BookingListDTO, error) {
	query := `
        SELECT 
            b.code, 
            b.created_at, 
            b.updated_at, 
            b.total, 
            string_agg(bi.name, ', ') AS item_name, 
            sum(bi.quantity) AS quantity, 
            bid.status AS ktp_status
        FROM booking b
        LEFT JOIN booking_item bi ON b.id = bi.booking_id
        LEFT JOIN booking_identity bid ON b.id = bid.booking_id
        GROUP BY b.id, b.code, b.created_at, b.updated_at, b.total, bid.status
        ORDER BY b.created_at DESC
    `
	var bookings []model.BookingListDTO
	err := r.db.Select(&bookings, query)
	if err != nil {
		log.Printf("GetListBookings: error: %v", err)
		return nil, err
	}
	log.Printf("GetListBookings: found %d bookings", len(bookings))
	return bookings, nil
}

func (r *customerRespository) UpdateIdentity(identity *model.IdentityModel) error {
	query := `
        UPDATE identity
        SET
            ktp_url = $1,
            verified = $2,
            status = $3,
            rejected_reason = $4,
            verified_at = $5,
            updated_at = NOW()
        WHERE id = $6
    `
	_, err := r.db.Exec(query, identity.KTPURL, identity.Verified, identity.Status, identity.RejectedReason, identity.VerifiedAt, identity.ID)
	if err != nil {
		log.Printf("UpdateIdentity: error updating identity: %v", err)
		return err
	}
	log.Printf("UpdateIdentity: updated identity %s", identity.ID)
	return nil
}

type CustomerRepository interface {
	CreateCustomer(customer *model.CustomerModel) error
	FindByEmailCustomerForLogin(email string) (*model.CustomerModel, error)
	GetDetailCustomer(id string) (*model.CustomerModel, error)
	UpdateCustomer(customer *model.CustomerModel) error
	DeleteCustomer(id string) error
	CreateIdentity(identity *model.IdentityModel) error
	CheckIdentityExists(userID string) (bool, error)
	GetIdentityByUserID(userID string) (*model.IdentityModel, error)                                                                            // Tambahkan ini
	CreateBooking(booking *model.BookingModel, items []model.BookingItem, customer model.BookingCustomer, identity model.BookingIdentity) error // Tambahkan ini
	GetBookingsByUserID(userID string) ([]model.BookingListDTO, error)
	GetListBookings() ([]model.BookingListDTO, error) // Rename dari GetAllBookings
	UpdateIdentity(identity *model.IdentityModel) error
}

/*
NewCustomerRepository membuat instance CustomerRepository.
Menginisialisasi dengan koneksi database.
*/
func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRespository{db: db}
}
