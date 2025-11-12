package repository

import (
	"database/sql"
	"encoding/json"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Implementasi repository TAC dengan koneksi database.
*/
type termsAndConditionsRepository struct {
	db *sqlx.DB
}

/*
Mencari TAC berdasarkan user ID.
Mengembalikan data TAC atau nil jika tidak ditemukan.
*/
func (r *termsAndConditionsRepository) FindByUserID(userID string) (*model.TermsAndConditionsModel, error) {
	query := `SELECT id, user_id, description, created_at, updated_at 
			  FROM tnc WHERE user_id = $1 LIMIT 1`
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

	query := `INSERT INTO tnc (id, user_id, description, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW())`
	_, err = r.db.Exec(query, tac.ID, tac.UserID, descriptionJSON)

	if err != nil {
		log.Printf("CreateTermAndConditions error: %v", err)
		return err
	}
	return nil
}

/*
Mencari TAC berdasarkan ID.
Mengembalikan data TAC atau nil jika tidak ditemukan.
*/
func (r *termsAndConditionsRepository) FindTermAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
	query := `SELECT id, user_id, description, created_at, updated_at 
			  FROM tnc WHERE id = $1 LIMIT 1`
	var tnc model.TermsAndConditionsModel
	var descriptionJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&tnc.ID, &tnc.UserID, &descriptionJSON, &tnc.CreatedAt, &tnc.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByID error: %v", err)
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
Update TAC di database.
Mengembalikan error jika update gagal.
*/
func (r *termsAndConditionsRepository) UpdateTermAndConditions(tnc *model.TermsAndConditionsModel) error {
	// Marshal []string ke JSON
	descriptionJSON, err := json.Marshal(tnc.Description)
	if err != nil {
		log.Printf("Marshal description error: %v", err)
		return err
	}

	query := `UPDATE tnc 
			  SET description = $2, updated_at = NOW() 
			  WHERE id = $1`
	result, err := r.db.Exec(query, tnc.ID, descriptionJSON)
	if err != nil {
		log.Printf("Update error: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Update RowsAffected error: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Update: no rows affected for id %s", tnc.ID)
	}

	return nil
}

/*
Mengambil semua TAC.
Mengembalikan list TAC atau error jika gagal.
*/
func (r *termsAndConditionsRepository) FindAllTermAndConditions() ([]*model.TermsAndConditionsModel, error) {
	query := `SELECT id, user_id, description, created_at, updated_at 
			  FROM tnc ORDER BY created_at DESC`
	var tncs []*model.TermsAndConditionsModel
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("FindAll error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tnc model.TermsAndConditionsModel
		var descriptionJSON []byte
		err := rows.Scan(&tnc.ID, &tnc.UserID, &descriptionJSON, &tnc.CreatedAt, &tnc.UpdatedAt)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}

		// Unmarshal JSONB ke []string
		if err := json.Unmarshal(descriptionJSON, &tnc.Description); err != nil {
			log.Printf("Unmarshal description error: %v", err)
			return nil, err
		}

		tncs = append(tncs, &tnc)
	}

	return tncs, nil
}

/*
Mendefinisikan operasi repository untuk terms and conditions.
Menyediakan method untuk membuat dan mengambil TAC dengan hasil sukses atau error.
*/
type TermsAndConditionsRepository interface {
	CreateTermAndConditions(tnc *model.TermsAndConditionsModel) error
	FindByUserID(userID string) (*model.TermsAndConditionsModel, error)
	UpdateTermAndConditions(tnc *model.TermsAndConditionsModel) error
	FindTermAndConditionsByID(id string) (*model.TermsAndConditionsModel, error)
	FindAllTermAndConditions() ([]*model.TermsAndConditionsModel, error)
}

/*
Membuat repository TAC.
Mengembalikan instance TermsAndConditionsRepository yang siap digunakan.
*/
func NewTermsAndConditionsRepository(db *sqlx.DB) TermsAndConditionsRepository {
	return &termsAndConditionsRepository{db: db}
}
