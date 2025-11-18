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

type CustomerRepository interface {
	CreateCustomer(customer *model.CustomerModel) error
	FindByEmailCustomerForLogin(email string) (*model.CustomerModel, error)
	GetDetailCustomer(id string) (*model.CustomerModel, error)
	UpdateCustomer(customer *model.CustomerModel) error
	DeleteCustomer(id string) error
	CreateIdentity(identity *model.IdentityModel) error
	CheckIdentityExists(userID string) (bool, error)
	GetIdentityByUserID(userID string) (*model.IdentityModel, error) // Tambahkan ini
}

/*
NewCustomerRepository membuat instance CustomerRepository.
Menginisialisasi dengan koneksi database.
*/
func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRespository{db: db}
}
