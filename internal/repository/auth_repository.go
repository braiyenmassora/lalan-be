package repository

import (
	"database/sql"
	"log"

	"lalan-be/internal/model"

	"github.com/jmoiron/sqlx"
)

// kontrak untuk operasi database auth (CreateHoster, FindByEmail).
type AuthRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmail(email string) (*model.HosterModel, error)
	FindByEmailForLogin(email string) (*model.HosterModel, error)
}

// implementasi konkret pakai PostgreSQL (sqlx.DB)
type authRepository struct {
	db *sqlx.DB
}

// inisialisasi repository dengan koneksi DB
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

// INSERT data hoster baru ke tabel hosters
func (r *authRepository) CreateHoster(h *model.HosterModel) error {
	query := `
		INSERT INTO hosters (
			id, full_name, profile_photo, store_name, description, phone_number, email, address, password_hash
		)
		VALUES (:id, :full_name, :profile_photo, :store_name, :description, :phone_number, :email, :address, :password_hash)
	`
	_, err := r.db.NamedExec(query, h)
	return err
}

// cari hoster berdasarkan email; return nil kalau tidak ditemukan
func (r *authRepository) FindByEmail(email string) (*model.HosterModel, error) {
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

// cari hoster berdasarkan email
func (r *authRepository) FindByEmailForLogin(email string) (*model.HosterModel, error) {
	var h model.HosterModel

	// Gunakan $1 → PostgreSQL
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
		// Log untuk debug (hapus di prod)
		log.Printf("Login query failed: %v | email: %s", err, email)
		return nil, err
	}

	return &h, nil
}
