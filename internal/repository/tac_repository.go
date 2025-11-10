package repository

import (
	"database/sql"
	"encoding/json"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Mendefinisikan operasi repository untuk terms and conditions.
Menyediakan method untuk membuat dan mengambil TAC dengan hasil sukses atau error.
*/
type TermsAndConditionsRepository interface {
	CreateTermAndConditions(tac *model.TermsAndConditionsModel) error
	FindByUserID(userID string) (*model.TermsAndConditionsModel, error)
}

/*
Implementasi repository TAC dengan koneksi database.
*/
type termsAndConditionsRepository struct {
	db *sqlx.DB
}

/*
Membuat repository TAC.
Mengembalikan instance TermsAndConditionsRepository yang siap digunakan.
*/
func NewTermsAndConditionsRepository(db *sqlx.DB) TermsAndConditionsRepository {
	return &termsAndConditionsRepository{db: db}
}

/*
Mencari TAC berdasarkan user ID.
Mengembalikan data TAC atau nil jika tidak ditemukan.
*/
func (r *termsAndConditionsRepository) FindByUserID(userID string) (*model.TermsAndConditionsModel, error) {
	query := `SELECT id, user_id, description, created_at, updated_at 
			  FROM terms_and_conditions WHERE user_id = $1 LIMIT 1`
	var tnc model.TermsAndConditionsModel
	var descriptionJSON []byte
	err := r.db.QueryRow(query, userID).Scan(
		&tnc.ID, &tnc.UserID, &descriptionJSON, &tnc.CreatedAt, &tnc.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByUserID error: %v", err)
		return nil, err
	}

	// Unmarshal JSONB ke []string
	if err := json.Unmarshal(descriptionJSON, &tnc.Description); err != nil {
		log.Printf("Unmarshal description error: %v", err)
		return nil, err
	}

	return &tnc, nil
}

/*
Membuat TAC baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *termsAndConditionsRepository) CreateTermAndConditions(tac *model.TermsAndConditionsModel) error {
	// Marshal []string ke JSON
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("Marshal description error: %v", err)
		return err
	}

	query := `INSERT INTO terms_and_conditions (id, user_id, description, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW())`
	_, err = r.db.Exec(query, tac.ID, tac.UserID, descriptionJSON)

	if err != nil {
		log.Printf("CreateTermAndConditions error: %v", err)
		return err
	}
	return nil
}
