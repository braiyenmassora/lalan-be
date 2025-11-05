package repository

import (
	"database/sql"

	"lalan-be/internal/model"

	"github.com/jmoiron/sqlx"
)

// kontrak untuk operasi database hoster (CreateHoster, FindByEmail).
type HosterRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmail(email string) (*model.HosterModel, error)
}

// implementasi konkret pakai PostgreSQL (sqlx.DB)
type hosterRepository struct {
	db *sqlx.DB
}

// inisialisasi repository dengan koneksi DB
func NewHosterRepository(db *sqlx.DB) HosterRepository {
	return &hosterRepository{db: db}
}

// INSERT data hoster baru ke tabel hosters
func (r *hosterRepository) CreateHoster(h *model.HosterModel) error {
	query := `
		INSERT INTO hosters (
			id, owner_name, store_name, phone_number, email, address, password
		)
		VALUES (:id, :owner_name, :store_name, :phone_number, :email, :address, :password)
	`
	_, err := r.db.NamedExec(query, h)
	return err
}

// cari hoster berdasarkan email; return nil kalau tidak ditemukan
func (r *hosterRepository) FindByEmail(email string) (*model.HosterModel, error) {
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
