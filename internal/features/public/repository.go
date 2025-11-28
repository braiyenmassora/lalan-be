package public

import (
	"encoding/json"
	"lalan-be/internal/domain"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
publicRepository adalah implementasi repository untuk data publik.
Bertanggung jawab atas query langsung ke database tanpa business logic.
*/
type publicRepository struct {
	db *sqlx.DB
}

/*
GetAllCategory mengambil semua kategori dari tabel category.

Alur kerja:
1. Eksekusi query SELECT sederhana
2. Mapping hasil ke slice []*domain.CategoryModel menggunakan sqlx.Select

Output sukses:
- ([]*domain.CategoryModel, nil)
Output error:
- (nil, error) → query gagal / koneksi DB bermasalah
*/
func (r *publicRepository) GetAllCategory() ([]*domain.Category, error) {
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

	var categories []*domain.Category
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("GetAllCategory repository error: %v", err)
		return nil, err
	}

	if categories == nil {
		categories = make([]*domain.Category, 0)
	}

	return categories, nil
}

/*
GetAllItems mengambil semua item publik beserta foto dalam format JSON.

Alur kerja:
1. Query semua kolom dari tabel item
2. Manual scan + json.Unmarshal untuk field photos (karena tipe []string di DB disimpan sebagai JSON)
3. Append ke slice hasil

Output sukses:
- ([]*domain.ItemModel, nil)
Output error:
- (nil, error) → query / scan / unmarshal gagal
*/
func (r *publicRepository) GetAllItems() ([]*domain.Item, error) {
	query := `
		SELECT
			id, name, description, photos, stock, pickup_type,
			price_per_day, deposit, discount, category_id, user_id,
			created_at, updated_at
		FROM item
	`

	var items []*domain.Item
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("GetAllItems query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Item
		var photosJSON []byte

		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
			&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
			&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			log.Printf("GetAllItems scan error: %v", err)
			return nil, err
		}

		if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
			log.Printf("GetAllItems unmarshal photos error: %v", err)
			return nil, err
		}

		items = append(items, &item)
	}

	if items == nil {
		items = make([]*domain.Item, 0)
	}

	return items, nil
}

/*
GetAllTermsAndConditions mengambil semua syarat & ketentuan dengan unmarshal JSON description.

Alur kerja:
1. Query semua record dari tabel tnc
2. Manual scan + json.Unmarshal untuk field description (tipe map/string disimpan sebagai JSON)

Output sukses:
- ([]*domain.TermsAndConditionsModel, nil)
Output error:
- (nil, error) → query / scan / unmarshal gagal
*/
func (r *publicRepository) GetAllTermsAndConditions() ([]*domain.TermsAndConditions, error) {
	query := `
		SELECT
			id, user_id, description, created_at, updated_at
		FROM tnc
	`

	var terms []*domain.TermsAndConditions
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("GetAllTermsAndConditions query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tac domain.TermsAndConditions
		var descriptionJSON []byte

		err := rows.Scan(&tac.ID, &tac.UserID, &descriptionJSON, &tac.CreatedAt, &tac.UpdatedAt)
		if err != nil {
			log.Printf("GetAllTermsAndConditions scan error: %v", err)
			return nil, err
		}

		if err := json.Unmarshal(descriptionJSON, &tac.Description); err != nil {
			log.Printf("GetAllTermsAndConditions unmarshal description error: %v", err)
			return nil, err
		}

		terms = append(terms, &tac)
	}

	if terms == nil {
		terms = make([]*domain.TermsAndConditions, 0)
	}

	return terms, nil
}

/*
PublicRepository adalah kontrak untuk operasi data publik.
Digunakan oleh service layer untuk dependency injection.
*/
type PublicRepository interface {
	GetAllCategory() ([]*domain.Category, error)
	GetAllItems() ([]*domain.Item, error)
	GetAllTermsAndConditions() ([]*domain.TermsAndConditions, error)
}

/*
NewPublicRepository membuat instance repository dengan koneksi database.

Output:
- PublicRepository siap digunakan
*/
func NewPublicRepository(db *sqlx.DB) PublicRepository {
	return &publicRepository{db: db}
}
