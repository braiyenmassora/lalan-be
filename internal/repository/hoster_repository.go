package repository

import (
	"database/sql"
	"log"

	"lalan-be/internal/model"

	"github.com/jmoiron/sqlx"
)

/*
Implementasi repository autentikasi dengan koneksi database.
*/
type authRepository struct {
	db *sqlx.DB
}

/*
Membuat hoster baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *authRepository) CreateHoster(h *model.HosterModel) error {
	query := `
        INSERT INTO hoster (
            id, full_name, profile_photo, store_name, description, phone_number, email, address, password_hash, website, instagram, tiktok
        )
        VALUES (:id, :full_name, :profile_photo, :store_name, :description, :phone_number, :email, :address, :password_hash, :website, :instagram, :tiktok)
    `
	_, err := r.db.NamedExec(query, h)
	if err != nil {
		log.Printf("Error inserting hoster: %v", err)
		return err
	}

	return nil
}

/*
Mengambil hoster berdasarkan email.
Mengembalikan data hoster atau nil jika tidak ditemukan.
*/
func (r *authRepository) FindByEmail(email string) (*model.HosterModel, error) {
	var hoster model.HosterModel
	query := `SELECT * FROM hoster WHERE email = $1 LIMIT 1`

	err := r.db.Get(&hoster, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &hoster, nil
}

/*
Mengambil hoster untuk login berdasarkan email.
Mengembalikan data hoster lengkap atau nil jika tidak ditemukan.
*/
func (r *authRepository) FindByEmailForLogin(email string) (*model.HosterModel, error) {
	var h model.HosterModel
	query := `
        SELECT
            id, email, password_hash, full_name, phone_number,
            store_name, description, address, profile_photo, website, instagram, tiktok,
            created_at, updated_at
        FROM hoster
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
		&h.Website,
		&h.Instagram,
		&h.Tiktok,
		&h.CreatedAt,
		&h.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &h, nil
}

/*
Mengambil hoster berdasarkan ID.
Mengembalikan data hoster tanpa password atau nil jika tidak ditemukan.
*/
func (r *authRepository) GetHosterByID(id string) (*model.HosterModel, error) {
	var hoster model.HosterModel
	query := `
        SELECT
            id, full_name, profile_photo, store_name, description, website, instagram, tiktok,
            phone_number, email, address, created_at, updated_at
        FROM hoster
        WHERE id = $1
        LIMIT 1
    `

	err := r.db.Get(&hoster, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &hoster, nil
}

/*
Mendefinisikan operasi repository untuk autentikasi hoster.
Menyediakan method untuk membuat dan mengambil data hoster dengan hasil sukses atau error.
*/
type AuthRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmail(email string) (*model.HosterModel, error)
	FindByEmailForLogin(email string) (*model.HosterModel, error)
	GetHosterByID(id string) (*model.HosterModel, error)
}

/*
Membuat repository autentikasi.
Mengembalikan instance AuthRepository yang siap digunakan.
*/
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}
