package repository

import (
	"database/sql"
	"log"

	"lalan-be/internal/model"

	"github.com/jmoiron/sqlx"
)

// Paket repository untuk akses data database.

// Interface untuk operasi autentikasi hoster.
type AuthRepository interface {
	CreateHoster(hoster *model.HosterModel) error                 // Simpan hoster baru
	FindByEmail(email string) (*model.HosterModel, error)         // Cari hoster berdasarkan email
	FindByEmailForLogin(email string) (*model.HosterModel, error) // Cari hoster untuk login
}

// Struktur implementasi repository autentikasi.
type authRepository struct {
	db *sqlx.DB // Koneksi database
}

// Buat instance repository autentikasi dengan koneksi database.
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

// Simpan hoster baru ke database.
func (r *authRepository) CreateHoster(h *model.HosterModel) error {
	// Insert data hoster
	query := `
		INSERT INTO hosters (
			id, full_name, profile_photo, store_name, description, phone_number, email, address, password_hash
		)
		VALUES (:id, :full_name, :profile_photo, :store_name, :description, :phone_number, :email, :address, :password_hash)
	`
	_, err := r.db.NamedExec(query, h)
	return err
}

// Cari hoster berdasarkan email.
func (r *authRepository) FindByEmail(email string) (*model.HosterModel, error) {
	// Cari hoster berdasarkan email
	var hoster model.HosterModel
	query := `SELECT * FROM hosters WHERE email = $1 LIMIT 1`

	err := r.db.Get(&hoster, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &hoster, nil
}

// Cari hoster untuk proses login dengan field tertentu.
func (r *authRepository) FindByEmailForLogin(email string) (*model.HosterModel, error) {
	// Cari hoster untuk login
	var h model.HosterModel

	query := `
		SELECT
			id, email, password_hash, full_name, phone_number,
			store_name, description, address, profile_photo,
			created_at, updated_at
		FROM hosters
		WHERE email = $1
		  AND password_hash IS NOT NULL
		LIMIT 1
	`

	// Scan hasil query
	err := r.db.QueryRow(query, email).Scan(
		&h.ID,
		&h.Email,
		&h.PasswordHash,
		&h.FullName,
		&h.PhoneNumber,
		&h.StoreName,
		&h.Description,
		&h.Address,
		&h.ProfilePhoto,
		&h.CreatedAt,
		&h.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Login query failed: %v | email: %s", err, email)
		return nil, err
	}

	return &h, nil
}
