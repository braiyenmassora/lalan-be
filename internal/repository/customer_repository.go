package repository

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Implementasi repository pelanggan dengan koneksi database.
*/
type customerRepository struct {
	db *sqlx.DB
}

/*
Membuat pelanggan baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *customerRepository) CreateCustomer(customer *model.CustomerModel) error {
	query := `
        INSERT INTO customer (
            id, full_name, profile_photo, phone_number, email, address, password_hash
        )
        VALUES (:id, :full_name, :profile_photo, :phone_number, :email, :address, :password_hash)
    `
	_, err := r.db.NamedExec(query, customer)
	if err != nil {
		log.Printf("Error inserting customer: %v", err)
		return err // Return the error instead of nil
	}
	return nil
}

/*
Mencari pelanggan berdasarkan email untuk login.
Mengembalikan data pelanggan atau nil jika tidak ditemukan.
*/
func (r *customerRepository) FindByEmailCustomerForLogin(email string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	query := `
        SELECT
            id, email, password_hash, full_name, phone_number,
            profile_photo, address,
            created_at, updated_at
        FROM customer  
        WHERE email = $1
          AND password_hash IS NOT NULL
        LIMIT 1
    `

	err := r.db.Get(&customer, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

/*
Mengambil pelanggan berdasarkan ID.
Mengembalikan data pelanggan atau nil jika tidak ditemukan.
*/
func (r *customerRepository) GetCustomerByID(id string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	query := `
        SELECT
            id, full_name, profile_photo,
            phone_number, email, address, created_at, updated_at
        FROM customer  
        WHERE id = $1
        LIMIT 1
    `

	err := r.db.Get(&customer, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

/*
Mendefinisikan operasi repository untuk pelanggan.
Menyediakan method untuk membuat, mencari, dan mengambil pelanggan dengan hasil sukses atau error.
*/
type CustomerRepository interface {
	CreateCustomer(customer *model.CustomerModel) error
	FindByEmailCustomerForLogin(email string) (*model.CustomerModel, error)
	GetCustomerByID(id string) (*model.CustomerModel, error)
}

/*
Membuat repository pelanggan.
Mengembalikan instance CustomerRepository yang siap digunakan.
*/
func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRepository{db: db}
}
