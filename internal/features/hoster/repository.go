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
hosterRespository menyediakan akses database untuk hoster.
Menggunakan sqlx.DB untuk operasi CRUD.
*/
type hosterRespository struct {
	db *sqlx.DB
}

/*
Methods untuk hosterRespository menangani operasi database hoster, item, dan terms.
Dipanggil oleh service untuk akses data.
*/
func (r *hosterRespository) CreateHoster(hoster *model.HosterModel) error {
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

func (r *hosterRespository) FindByEmailHosterForLogin(email string) (*model.HosterModel, error) {
	var hoster model.HosterModel
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

func (r *hosterRespository) GetDetailHoster(id string) (*model.HosterModel, error) {
	var hoster model.HosterModel
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

func (r *hosterRespository) CreateItem(item *model.ItemModel) error {
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("CreateItem: error marshaling photos: %v", err)
		return err
	}
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

func (r *hosterRespository) FindItemNameByID(id string) (*model.ItemModel, error) {
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

func (r *hosterRespository) FindItemNameByUserID(name string, userId string) (*model.ItemModel, error) {
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

func (r *hosterRespository) GetAllItems() ([]*model.ItemModel, error) {
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

func (r *hosterRespository) UpdateItem(item *model.ItemModel) error {
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

func (r *hosterRespository) DeleteItem(id string) error {
	query := `DELETE FROM item WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteItem: error deleting item: %v", err)
		return err
	}
	return nil
}

func (r *hosterRespository) CreateTermsAndConditions(tac *model.TermsAndConditionsModel) error {
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error marshaling description: %v", err)
		return err
	}
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

func (r *hosterRespository) FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
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

func (r *hosterRespository) FindTermsAndConditionsByUserIDAndDescription(userID string, description []string) (*model.TermsAndConditionsModel, error) {
	descriptionJSON, err := json.Marshal(description)
	if err != nil {
		log.Printf("FindTermsAndConditionsByUserIDAndDescription: error marshaling description: %v", err)
		return nil, err
	}
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

func (r *hosterRespository) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
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

func (r *hosterRespository) UpdateTermsAndConditions(tac *model.TermsAndConditionsModel) error {
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("UpdateTerm: error marshaling description: %v", err)
		return err
	}
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

func (r *hosterRespository) DeleteTermsAndConditions(id string) error {
	query := `DELETE FROM tnc WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteTerm: error deleting tnc: %v", err)
		return err
	}
	return nil
}

func (r *hosterRespository) GetIdentityCustomer(userID string) (*model.IdentityModel, error) {
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

func (r *hosterRespository) UpdateIdentityStatus(identityID string, status string, rejectedReason string, verified bool, verifiedAt *time.Time) error {
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

func (r *hosterRespository) GetIdentityCustomerByID(identityID string) (*model.IdentityModel, error) {
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
        WHERE id = $1
    `
	err := r.db.Get(&identity, query, identityID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetIdentityCustomerByID: no identity found for ID %s", identityID)
			return nil, nil
		}
		log.Printf("GetIdentityCustomerByID: error for ID %s: %v", identityID, err)
		return nil, err
	}
	log.Printf("GetIdentityCustomerByID: found identity for ID %s", identityID)
	return &identity, nil
}

/*
HosterRepository mendefinisikan kontrak operasi database hoster.
Diimplementasikan oleh hosterRespository.
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
	GetIdentityCustomerByID(identityID string) (*model.IdentityModel, error) // Tambahkan ini
}

/*
NewHosterRepository membuat instance HosterRepository.
Menginisialisasi dengan koneksi database.
*/
func NewHosterRepository(db *sqlx.DB) HosterRepository {
	return &hosterRespository{db: db}
}
