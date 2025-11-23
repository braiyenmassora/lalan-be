package hoster

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/model"
)

/*
type hosterRespository struct
menyediakan akses ke operasi database untuk hoster
*/
type hosterRespository struct {
	db *sqlx.DB
}

/*
CreateHoster
membuat hoster baru dengan data yang diberikan
*/
func (r *hosterRespository) CreateHoster(hoster *model.HosterModel) error {
	/*
	  CreateHoster query
	  membuat hoster baru dengan data yang diberikan
	*/
	query := `
		INSERT INTO hoster (
			full_name,
			store_name,
			address,
			phone_number,
			email,
			password_hash,
			profile_photo,
			description,
			tiktok,
			instagram,
			website,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, hoster.FullName, hoster.StoreName, hoster.Address, hoster.PhoneNumber, hoster.Email, hoster.PasswordHash, hoster.ProfilePhoto, hoster.Description, hoster.Tiktok, hoster.Instagram, hoster.Website, hoster.CreatedAt, hoster.UpdatedAt).Scan(&hoster.ID, &hoster.CreatedAt, &hoster.UpdatedAt)
	log.Printf("CreateHoster: inserted hoster with email %s, ID %s", hoster.Email, hoster.ID)
	return err
}

/*
FindByEmailHosterForLogin
mencari hoster berdasarkan email untuk login
*/
func (r *hosterRespository) FindByEmailHosterForLogin(email string) (*model.HosterModel, error) {
	var hoster model.HosterModel
	/*
	  FindByEmailHosterForLogin query
	  mencari hoster berdasarkan email untuk login
	*/
	query := `
		SELECT
			id,
			full_name,
			store_name,
			phone_number,
			email,
			password_hash,
			profile_photo,
			description,
			tiktok,
			instagram,
			website,
			created_at,
			updated_at
		FROM hoster
		WHERE email = $1
	`
	err := r.db.Get(&hoster, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindByEmailHosterForLogin: no hoster found for email %s", email)
			return nil, nil
		}
		log.Printf("FindByEmailHosterForLogin: error querying email %s: %v", email, err)
		return nil, err
	}
	log.Printf("FindByEmailHosterForLogin: found hoster for email %s", email)
	return &hoster, nil
}

/*
GetDetailHoster
mengambil detail hoster berdasarkan ID
*/
func (r *hosterRespository) GetDetailHoster(id string) (*model.HosterModel, error) {
	var hoster model.HosterModel
	/*
	  GetDetailHoster query
	  mengambil detail hoster berdasarkan ID
	*/
	query := `
        SELECT 
            id,
            full_name,
            store_name,
            address,
            phone_number,
            email,
            password_hash,
            profile_photo,
            description,
            tiktok,
            instagram,
            website,
            created_at,
            updated_at
        FROM hoster
        WHERE id = $1
    `
	err := r.db.Get(&hoster, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetDetailHoster: no hoster found for id %s", id)
			return nil, nil
		}
		log.Printf("GetDetailHoster: error for id %s: %v", id, err)
		return nil, err
	}
	log.Printf("GetDetailHoster: found hoster id %s", id)
	return &hoster, nil
}

/*
CreateItem
membuat item baru dengan data yang diberikan
*/
func (r *hosterRespository) CreateItem(item *model.ItemModel) error {
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("CreateItem: error marshaling photos: %v", err)
		return err
	}
	/*
	  CreateItem query
	  membuat item baru dengan data yang diberikan
	*/
	query := `
		INSERT INTO item (
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`
	_, err = r.db.Exec(query, item.ID, item.Name, item.Description, photosJSON,
		item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount,
		item.CategoryID, item.UserID)
	if err != nil {
		log.Printf("CreateItem: error inserting item: %v", err)
		return err
	}
	return nil
}

/*
FindItemNameByID
mencari item berdasarkan ID
*/
func (r *hosterRespository) FindItemNameByID(id string) (*model.ItemModel, error) {
	/*
	  FindItemNameByID query
	  mencari item berdasarkan ID
	*/
	query := `
		SELECT
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		FROM item
		WHERE id = $1
		LIMIT 1
	`
	var item model.ItemModel
	var photosJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByID error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
		log.Printf("Unmarshal photos error: %v", err)
		return nil, err
	}

	return &item, nil
}

/*
FindItemNameByUserID
mencari item berdasarkan nama dan user ID
*/
func (r *hosterRespository) FindItemNameByUserID(name string, userId string) (*model.ItemModel, error) {
	/*
	  FindItemNameByUserID query
	  mencari item berdasarkan nama dan user ID
	*/
	query := `
		SELECT
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		FROM item
		WHERE name = $1 AND user_id = $2
		LIMIT 1
	`
	var item model.ItemModel
	var photosJSON []byte
	err := r.db.QueryRow(query, name, userId).Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindItemNameByUserID error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
		log.Printf("Unmarshal photos error: %v", err)
		return nil, err
	}

	return &item, nil
}

/*
GetAllItems
mengambil semua item
*/
func (r *hosterRespository) GetAllItems() ([]*model.ItemModel, error) {
	/*
	  GetAllItems query
	  mengambil semua item
	*/
	query := `
		SELECT
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		FROM item
	`
	var items []*model.ItemModel
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ItemModel
		var photosJSON []byte
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock, &item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount, &item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

/*
UpdateItem
memperbarui item berdasarkan ID
*/
func (r *hosterRespository) UpdateItem(item *model.ItemModel) error {
	/*
	  UpdateItem query
	  memperbarui item berdasarkan ID
	*/
	query := `
		UPDATE item
		SET
			name = $1,
			description = $2,
			photos = $3,
			stock = $4,
			pickup_type = $5,
			price_per_day = $6,
			deposit = $7,
			discount = $8,
			category_id = $9,
			updated_at = $10
		WHERE id = $11
	`
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("UpdateItem: error marshaling photos: %v", err)
		return err
	}
	_, err = r.db.Exec(query, item.Name, item.Description, photosJSON, item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount, item.CategoryID, item.UpdatedAt, item.ID)
	if err != nil {
		log.Printf("UpdateItem: error updating item: %v", err)
		return err
	}
	return nil
}

/*
DeleteItem
menghapus item berdasarkan ID
*/
func (r *hosterRespository) DeleteItem(id string) error {
	/*
	  DeleteItem query
	  menghapus item berdasarkan ID
	*/
	query := `DELETE FROM item WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteItem: error deleting item: %v", err)
		return err
	}
	return nil
}

/*
CreateTermsAndConditions
membuat syarat dan ketentuan baru
*/
func (r *hosterRespository) CreateTermsAndConditions(tac *model.TermsAndConditionsModel) error {
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error marshaling description: %v", err)
		return err
	}
	/*
	  CreateTermsAndConditions query
	  membuat syarat dan ketentuan baru
	*/
	query := `
		INSERT INTO tnc (
			id,
			user_id,
			description,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, NOW(), NOW())
	`
	_, err = r.db.Exec(query, tac.ID, tac.UserID, descriptionJSON)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error inserting tnc: %v", err)
		return err
	}
	return nil
}

/*
FindTermsAndConditionsByID
mencari syarat dan ketentuan berdasarkan ID
*/
func (r *hosterRespository) FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
	/*
	  FindTermsAndConditionsByID query
	  mencari syarat dan ketentuan berdasarkan ID
	*/
	query := `
		SELECT
			id,
			user_id,
			description,
			created_at,
			updated_at
		FROM tnc
		WHERE id = $1
		LIMIT 1
	`
	var tac model.TermsAndConditionsModel
	var descriptionJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&tac.ID, &tac.UserID, &descriptionJSON, &tac.CreatedAt, &tac.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindTermByID error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(descriptionJSON, &tac.Description); err != nil {
		log.Printf("Unmarshal description error: %v", err)
		return nil, err
	}

	return &tac, nil
}

/*
FindTermsAndConditionsByUserIDAndDescription
mencari syarat dan ketentuan berdasarkan user ID dan deskripsi
*/
func (r *hosterRespository) FindTermsAndConditionsByUserIDAndDescription(userID string, description []string) (*model.TermsAndConditionsModel, error) {
	descriptionJSON, err := json.Marshal(description)
	if err != nil {
		log.Printf("FindTermsAndConditionsByUserIDAndDescription: error marshaling description: %v", err)
		return nil, err
	}
	/*
	  FindTermsAndConditionsByUserIDAndDescription query
	  mencari syarat dan ketentuan berdasarkan user ID dan deskripsi
	*/
	query := `
        SELECT
            id,
            user_id,
            description,
            created_at,
            updated_at
        FROM tnc
        WHERE user_id = $1 AND description = $2
        LIMIT 1
    `
	var tac model.TermsAndConditionsModel
	var descJSON []byte
	err = r.db.QueryRow(query, userID, descriptionJSON).Scan(
		&tac.ID, &tac.UserID, &descJSON, &tac.CreatedAt, &tac.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindTermsAndConditionsByUserIDAndDescription error: %v", err)
		return nil, err
	}
	if err := json.Unmarshal(descJSON, &tac.Description); err != nil {
		log.Printf("Unmarshal description error: %v", err)
		return nil, err
	}
	return &tac, nil
}

/*
GetAllTermsAndConditions
mengambil semua syarat dan ketentuan
*/
func (r *hosterRespository) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	/*
	  GetAllTermsAndConditions query
	  mengambil semua syarat dan ketentuan
	*/
	query := `
		SELECT
			id,
			user_id,
			description,
			created_at,
			updated_at
		FROM tnc
	`
	var terms []*model.TermsAndConditionsModel
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tac model.TermsAndConditionsModel
		var descriptionJSON []byte
		err := rows.Scan(&tac.ID, &tac.UserID, &descriptionJSON, &tac.CreatedAt, &tac.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(descriptionJSON, &tac.Description); err != nil {
			return nil, err
		}
		terms = append(terms, &tac)
	}
	return terms, nil
}

/*
UpdateTermsAndConditions
memperbarui syarat dan ketentuan berdasarkan ID
*/
func (r *hosterRespository) UpdateTermsAndConditions(tac *model.TermsAndConditionsModel) error {
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("UpdateTerm: error marshaling description: %v", err)
		return err
	}
	/*
	  UpdateTermsAndConditions query
	  memperbarui syarat dan ketentuan berdasarkan ID
	*/
	query := `
		UPDATE tnc
		SET
			description = $1,
			updated_at = NOW()
		WHERE id = $2
	`
	_, err = r.db.Exec(query, descriptionJSON, tac.ID)
	if err != nil {
		log.Printf("UpdateTerm: error updating tnc: %v", err)
		return err
	}
	return nil
}

/*
DeleteTermsAndConditions
menghapus syarat dan ketentuan berdasarkan ID
*/
func (r *hosterRespository) DeleteTermsAndConditions(id string) error {
	/*
	  DeleteTermsAndConditions query
	  menghapus syarat dan ketentuan berdasarkan ID
	*/
	query := `DELETE FROM tnc WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteTerm: error deleting tnc: %v", err)
		return err
	}
	return nil
}

/*
GetIdentityCustomer
mengambil identitas customer berdasarkan user ID
*/
func (r *hosterRespository) GetIdentityCustomer(userID string) (*model.IdentityModel, error) {
	var identity model.IdentityModel
	/*
	  GetIdentityCustomer query
	  mengambil identitas customer berdasarkan user ID
	*/
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
	log.Printf("GetIdentityCustomer: query executed for userID %s, err: %v, identity: %+v", userID, err, identity) // Tambahkan log
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetIdentityCustomer: no identity found for user %s", userID)
			return nil, nil
		}
		log.Printf("GetIdentityCustomer: error for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("GetIdentityCustomer: found identity for user %s", userID)
	return &identity, nil
}

/*
UpdateIdentityStatus
memperbarui status identitas berdasarkan ID
*/
func (r *hosterRespository) UpdateIdentityStatus(identityID string, status string, rejectedReason string, verified bool, verifiedAt *time.Time) error {
	/*
	  UpdateIdentityStatus query
	  memperbarui status identitas berdasarkan ID
	*/
	query := `
        UPDATE identity
        SET
            status = $1,
            rejected_reason = $2,
            verified = $3,
            verified_at = $4,
            updated_at = NOW()
        WHERE id = $5
    `
	_, err := r.db.Exec(query, status, rejectedReason, verified, verifiedAt, identityID)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error updating identity ID %s: %v", identityID, err)
		return err
	}
	log.Printf("UpdateIdentityStatus: updated identity ID %s to status %s", identityID, status)
	return nil
}

/*
UpdateBookingStatusByUserID
memperbarui status booking berdasarkan user ID
*/
func (r *hosterRespository) UpdateBookingStatusByUserID(userID, status string) error {
	/*
	  UpdateBookingStatusByUserID query
	  memperbarui status booking berdasarkan user ID
	*/
	query := `UPDATE booking SET status = $1, updated_at = NOW() WHERE user_id = $2 AND status = 'waiting_ktp_verification'`
	_, err := r.db.Exec(query, status, userID)
	if err != nil {
		log.Printf("UpdateBookingStatusByUserID: error updating bookings for user %s: %v", userID, err)
		return err
	}
	log.Printf("UpdateBookingStatusByUserID: updated bookings for user %s to status %s", userID, status)
	return nil
}

/*
UpdateBookingIdentityStatusByUserID
memperbarui status identitas booking berdasarkan user ID
*/
func (r *hosterRespository) UpdateBookingIdentityStatusByUserID(userID, status string) error {
	/*
	  UpdateBookingIdentityStatusByUserID query
	  memperbarui status identitas booking berdasarkan user ID
	*/
	query := `
        UPDATE booking_identity 
        SET status = $1, updated_at = NOW() 
        WHERE booking_id IN (
            SELECT id FROM booking WHERE user_id = $2
        )
    `
	result, err := r.db.Exec(query, status, userID)
	if err != nil {
		log.Printf("UpdateBookingIdentityStatusByUserID: error updating for user %s: %v", userID, err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("UpdateBookingIdentityStatusByUserID: rows affected: %d for user %s to status %s", rowsAffected, userID, status)
	if rowsAffected == 0 {
		log.Printf("UpdateBookingIdentityStatusByUserID: no booking identities updated for user %s", userID)
	}
	return nil
}

/*
GetListBookingsForHoster
mengambil daftar booking yang dimiliki oleh hoster berdasarkan user ID hoster dengan limit dan offset
*/
func (r *hosterRespository) GetListBookingsForHoster(userID string, limit int, offset int) ([]model.BookingListCustomer, error) {
	/*
	   GetListBookingsForHoster query
	   mengambil daftar booking dengan agregasi item untuk hoster tertentu
	*/
	query := `
SELECT
    b.id AS booking_id,
    b.code,
    b.user_id AS customer_id,
    bc.name AS customer_name,
    b.start_date,
    b.end_date,
    b.total_days AS duration_days,
    b.total,
    COALESCE(i.status, '') AS identity_status,  // Ambil dari identity global
    COALESCE(bi.item_summary, '') AS item_summary,
    COALESCE(bi.quantity, 0) AS quantity
FROM booking b
LEFT JOIN booking_customer bc ON bc.booking_id = b.id
LEFT JOIN identity i ON i.user_id = b.user_id  // Join identity berdasarkan user_id
LEFT JOIN (
    SELECT booking_id, 
           string_agg(name || ' x' || quantity, ', ') AS item_summary,
           SUM(quantity) AS quantity
    FROM booking_item
    GROUP BY booking_id
) bi ON bi.booking_id = b.id
WHERE b.hoster_id = $1
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3
`

	var bookings []model.BookingListCustomer
	err := r.db.Select(&bookings, query, userID, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsForHoster: error: %v", err)
		return nil, err
	}
	log.Printf("GetListBookingsForHoster: found %d bookings for hoster %s", len(bookings), userID)
	return bookings, nil
}

/*
GetListBookingsForHosterByCustomerID
mengambil daftar booking yang dimiliki oleh hoster berdasarkan hoster ID dan customer ID dengan limit dan offset
*/
func (r *hosterRespository) GetListBookingsForHosterByCustomerID(hosterID string, customerID string, limit int, offset int) ([]model.BookingListCustomer, error) {
	/*
	   GetListBookingsForHosterByCustomerID query
	   mengambil daftar booking dengan agregasi item untuk hoster tertentu dan customer tertentu
	*/
	query := `
SELECT
    b.id AS booking_id,
    b.code,
    b.user_id AS customer_id,
    bc.name AS customer_name,
    b.start_date,
    b.end_date,
    b.total_days AS duration_days,
    b.total,
    COALESCE(i.status, '') AS identity_status,  // Ambil dari identity global
    COALESCE(bi.item_summary, '') AS item_summary,
    COALESCE(bi.quantity, 0) AS quantity
FROM booking b
LEFT JOIN booking_customer bc ON bc.booking_id = b.id
LEFT JOIN identity i ON i.user_id = b.user_id  // Join identity berdasarkan user_id
LEFT JOIN (
    SELECT booking_id, 
           string_agg(name || ' x' || quantity, ', ') AS item_summary,
           SUM(quantity) AS quantity
    FROM booking_item
    GROUP BY booking_id
) bi ON bi.booking_id = b.id
WHERE b.hoster_id = $1 AND b.user_id = $2
ORDER BY b.created_at DESC
LIMIT $3 OFFSET $4
`

	var bookings []model.BookingListCustomer
	err := r.db.Select(&bookings, query, hosterID, customerID, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsForHosterByCustomerID: error: %v", err)
		return nil, err
	}
	log.Printf("GetListBookingsForHosterByCustomerID: found %d bookings for hoster %s and customer %s", len(bookings), hosterID, customerID)
	return bookings, nil
}

/*
GetListBookingsCustomer
mengambil daftar booking yang dimiliki oleh hoster berdasarkan user ID hoster dengan limit dan offset
*/
func (r *hosterRespository) GetListBookingsCustomer(userID string, limit int, offset int) ([]model.BookingListDTOHoster, error) {
	/*
	   GetListBookingsCustomer query
	   mengambil daftar booking dengan agregasi item untuk hoster tertentu
	*/
	query := `
SELECT
    b.id AS booking_id,
    bc.name AS customer_name,
    STRING_AGG(bi.name || ' x' || bi.quantity, ', ') AS item_summary,
    b.start_date,
    b.end_date,
    b.total_days,
    b.total,
    i.status AS identity_status
FROM booking b
JOIN booking_customer bc ON bc.booking_id = b.id
JOIN booking_item bi ON bi.booking_id = b.id
LEFT JOIN identity i ON i.id = b.identity_id
WHERE b.hoster_id = $1
GROUP BY b.id, bc.name, b.start_date, b.end_date, b.total_days, b.total, i.status
ORDER BY b.start_date DESC
LIMIT $2 OFFSET $3
`

	var bookings []model.BookingListDTOHoster
	err := r.db.Select(&bookings, query, userID, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsCustomer: error: %v", err)
		return nil, err
	}
	log.Printf("GetListBookingsCustomer: found %d bookings for hoster %s", len(bookings), userID)
	return bookings, nil
}

/*
GetListBookingsCustomerByBookingID
mengambil daftar booking yang dimiliki oleh hoster berdasarkan hoster ID dan booking ID dengan limit dan offset
*/
func (r *hosterRespository) GetListBookingsCustomerByBookingID(hosterID string, bookingID string, limit int, offset int) ([]model.BookingDetailDTOHoster, error) {
	/*
	   GetListBookingsCustomerByBookingID query
	   mengambil daftar booking dengan detail item untuk hoster tertentu dan booking tertentu
	*/
	query := `
SELECT
    b.id AS booking_id,
    b.code,
    b.hoster_id,
    b.locked_until,
    b.start_date,
    b.end_date,
    b.total_days,
    b.delivery_type,
    b.rental,
    b.deposit,
    b.discount,
    b.total,
    b.outstanding,
    b.user_id,
    b.identity_id,
    b.status,
    b.created_at AS booking_created_at,
    b.updated_at AS booking_updated_at,

    bi.id AS booking_item_id,
    bi.item_id,
    bi.name AS item_name,
    bi.quantity,
    bi.price_per_day,
    bi.deposit_per_unit,
    bi.subtotal_rental,
    bi.subtotal_deposit,

    bc.id AS booking_customer_id,
    bc.name AS customer_name,
    bc.phone,
    bc.email,
    bc.address,
    bc.notes,

    i.ktp_url,
    i.status AS identity_status,
    i.reason AS identity_reason
FROM booking b
JOIN booking_item bi ON bi.booking_id = b.id
JOIN booking_customer bc ON bc.booking_id = b.id
LEFT JOIN identity i ON i.id = b.identity_id
WHERE b.hoster_id = $1 AND b.id = $2
ORDER BY b.created_at DESC
LIMIT $3 OFFSET $4
`

	var bookings []model.BookingDetailDTOHoster
	err := r.db.Select(&bookings, query, hosterID, bookingID, limit, offset)
	if err != nil {
		log.Printf("GetListBookingsCustomerByBookingID: error: %v", err)
		return nil, err
	}
	log.Printf("GetListBookingsCustomerByBookingID: found %d bookings for hoster %s and booking %s", len(bookings), hosterID, bookingID)
	return bookings, nil
}

/*
GetListCustomer
mengambil daftar customer unik yang telah booking dengan hoster tertentu
*/
func (r *hosterRespository) GetListCustomer(hosterID string) ([]model.CustomerIdentityDTO, error) {
	/*
	   GetListCustomer query
	   mengambil daftar customer dengan identity terbaru berdasarkan hoster ID
	*/
	query := `
        SELECT
            c.id AS customer_id,
            c.full_name,
            c.email,
            c.phone_number,
            i.id AS identity_id,
            i.ktp_url,
            i.verified,
            i.status,
            i.reason,
            i.verified_at,
            i.created_at AS identity_created_at,
            i.updated_at AS identity_updated_at
        FROM customer c
        LEFT JOIN LATERAL (
            SELECT *
            FROM identity i
            WHERE i.user_id = c.id
            ORDER BY i.created_at DESC
            LIMIT 1
        ) i ON true
        WHERE c.id IN (SELECT DISTINCT user_id FROM booking WHERE hoster_id = $1)
        ORDER BY c.full_name
    `
	var customers []model.CustomerIdentityDTO
	err := r.db.Select(&customers, query, hosterID)
	if err != nil {
		log.Printf("GetListCustomer: error: %v", err)
		return nil, err
	}
	log.Printf("GetListCustomer: found %d customers for hoster %s", len(customers), hosterID)
	return customers, nil
}

/*
HosterRepository
mendefinisikan kontrak untuk akses data hoster
*/
type HosterRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmailHosterForLogin(email string) (*model.HosterModel, error)
	GetDetailHoster(id string) (*model.HosterModel, error)
	CreateItem(item *model.ItemModel) error
	FindItemNameByUserID(name string, userId string) (*model.ItemModel, error)
	FindItemNameByID(id string) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	UpdateItem(item *model.ItemModel) error
	DeleteItem(id string) error
	CreateTermsAndConditions(tac *model.TermsAndConditionsModel) error
	FindTermsAndConditionsByID(name string) (*model.TermsAndConditionsModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
	UpdateTermsAndConditions(tac *model.TermsAndConditionsModel) error
	DeleteTermsAndConditions(id string) error
	FindTermsAndConditionsByUserIDAndDescription(userID string, description []string) (*model.TermsAndConditionsModel, error)
	GetIdentityCustomer(userID string) (*model.IdentityModel, error)
	UpdateIdentityStatus(identityID string, status string, rejectedReason string, verified bool, verifiedAt *time.Time) error
	UpdateBookingStatusByUserID(userID, status string) error
	UpdateBookingIdentityStatusByUserID(userID, status string) error
	GetListBookingsCustomer(userID string, limit int, offset int) ([]model.BookingListDTOHoster, error)
	GetListBookingsCustomerByBookingID(hosterID string, bookingID string, limit int, offset int) ([]model.BookingDetailDTOHoster, error)
	GetListCustomer(hosterID string) ([]model.CustomerIdentityDTO, error)
}

/*
NewHosterRepository
membuat instance baru HosterRepository dengan database yang diberikan
*/
func NewHosterRepository(db *sqlx.DB) HosterRepository {
	return &hosterRespository{db: db}
}
