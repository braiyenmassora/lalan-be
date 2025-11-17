package public

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/model"
)

/*
publicRepository menyediakan akses database untuk data publik.
Menggunakan sqlx.DB untuk query publik.
*/
type publicRepository struct {
	db *sqlx.DB
}

/*
Methods untuk publicRepository menangani query data publik kategori, item, dan terms.
Dipanggil oleh service untuk akses tanpa autentikasi.
*/
func (r *publicRepository) GetAllCategory() ([]*model.CategoryModel, error) {
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM category
		ORDER BY created_at DESC
	`
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("GetAllCategory error: %v", err)
		return nil, err
	}
	return categories, nil
}

func (r *publicRepository) GetAllItems() ([]*model.ItemModel, error) {
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

func (r *publicRepository) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
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

/*
PublicRepository mendefinisikan kontrak operasi data publik.
Diimplementasikan oleh publicRepository.
*/
type PublicRepository interface {
	GetAllCategory() ([]*model.CategoryModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
}

/*
NewPublicRepository membuat instance PublicRepository.
Menginisialisasi repository dengan koneksi database.
*/
func NewPublicRepository(db *sqlx.DB) PublicRepository {
	return &publicRepository{db: db}
}
