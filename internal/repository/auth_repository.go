package repository

import (
	"database/sql"
	"log"

	"lalan-be/internal/model"

	"github.com/jmoiron/sqlx"
)

/*
Interface untuk operasi autentikasi hoster.
Mendefinisikan method untuk membuat dan mengambil data hoster.
*/
type AuthRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmail(email string) (*model.HosterModel, error)
	FindByEmailForLogin(email string) (*model.HosterModel, error)
}

/*
Struct untuk menyimpan koneksi database.
Digunakan sebagai implementasi interface AuthRepository.
*/
type authRepository struct {
	db *sqlx.DB
}

/*
Membuat instance repository dengan koneksi database.
Mengembalikan interface AuthRepository.
*/
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

/*
Menyisipkan data hoster ke tabel hosters.
Mengembalikan error jika gagal.
*/
func (r *authRepository) CreateHoster(h *model.HosterModel) error {
	query := `
        INSERT INTO hosters (
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
Mengambil data hoster lengkap dari tabel hosters berdasarkan email.
Mengembalikan pointer ke model dan error; (nil, nil) jika tidak ada baris.
*/
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

/*
Mengambil data hoster untuk autentikasi dari tabel hosters berdasarkan email.
Mengembalikan pointer ke model dan error; (nil, nil) jika tidak ada baris.
*/
func (r *authRepository) FindByEmailForLogin(email string) (*model.HosterModel, error) {
	var h model.HosterModel

	query := `
        SELECT
            id, email, password_hash, full_name, phone_number,
            store_name, description, address, profile_photo, website, instagram, tiktok,
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
