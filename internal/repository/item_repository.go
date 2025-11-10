package repository

import (
	"database/sql"
	"encoding/json"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Mendefinisikan operasi repository untuk item.
Menyediakan method untuk mencari, membuat, update, dan hapus item dengan hasil sukses atau error.
*/
type ItemRepository interface {
	FindItemNameByUserID(name string, userId string) (*model.ItemModel, error)
	CreateItem(item *model.ItemModel) error
	FindAll() ([]*model.ItemModel, error)
	FindByID(id string) (*model.ItemModel, error)
	FindByUserID(userID string) ([]*model.ItemModel, error)
	Update(item *model.ItemModel) error
	Delete(id string) error
}

/*
Implementasi repository item dengan koneksi database.
*/
type itemRepository struct {
	db *sqlx.DB
}

/*
Membuat repository item.
Mengembalikan instance ItemRepository yang siap digunakan.
*/
func NewItemRepository(db *sqlx.DB) ItemRepository {
	return &itemRepository{db: db}
}

/*
Mencari item berdasarkan nama dan user ID.
Mengembalikan data item atau nil jika tidak ditemukan.
*/
func (r *itemRepository) FindItemNameByUserID(name string, userId string) (*model.ItemModel, error) {
	query := `SELECT id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at 
	          FROM items WHERE name = $1 AND user_id = $2 LIMIT 1`

	var item model.ItemModel
	var photosJSON []byte
	err := r.db.QueryRow(query, name, userId).Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindItemNameByUserID error: %v", err)
		return nil, err
	}

	// Unmarshal JSONB ke []string
	if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
		log.Printf("Unmarshal photos error: %v", err)
		return nil, err
	}

	return &item, nil
}

/*
Membuat item baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *itemRepository) CreateItem(item *model.ItemModel) error {
	// Marshal []string ke JSON
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("Marshal photos error: %v", err)
		return err
	}

	query := `INSERT INTO items (id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`

	_, err = r.db.Exec(query, item.ID, item.Name, item.Description, photosJSON,
		item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount,
		item.CategoryID, item.UserID)

	if err != nil {
		log.Printf("CreateItem error: %v", err)
		return err
	}

	return nil
}

/*
Mengambil semua item.
Mengembalikan daftar item atau error jika gagal.
*/
func (r *itemRepository) FindAll() ([]*model.ItemModel, error) {
	query := `SELECT id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at 
	          FROM items ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("FindAll error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []*model.ItemModel
	for rows.Next() {
		var item model.ItemModel
		var photosJSON []byte
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
			&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
			&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			log.Printf("FindAll scan error: %v", err)
			return nil, err
		}

		// Unmarshal JSONB ke []string
		if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
			log.Printf("Unmarshal photos error: %v", err)
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

/*
Mencari item berdasarkan ID.
Mengembalikan data item atau nil jika tidak ditemukan.
*/
func (r *itemRepository) FindByID(id string) (*model.ItemModel, error) {
	query := `SELECT id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at 
	          FROM items WHERE id = $1 LIMIT 1`

	var item model.ItemModel
	var photosJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByID error: %v", err)
		return nil, err
	}

	// Unmarshal JSONB ke []string
	if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
		log.Printf("Unmarshal photos error: %v", err)
		return nil, err
	}

	return &item, nil
}

/*
Mengambil item berdasarkan user ID.
Mengembalikan daftar item user atau error jika gagal.
*/
func (r *itemRepository) FindByUserID(userID string) ([]*model.ItemModel, error) {
	query := `SELECT id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at 
	          FROM items WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("FindByUserID error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []*model.ItemModel
	for rows.Next() {
		var item model.ItemModel
		var photosJSON []byte
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
			&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
			&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			log.Printf("FindByUserID scan error: %v", err)
			return nil, err
		}

		// Unmarshal JSONB ke []string
		if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
			log.Printf("Unmarshal photos error: %v", err)
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

/*
Memperbarui item.
Mengembalikan error jika update gagal atau tidak ada baris terpengaruh.
*/
func (r *itemRepository) Update(item *model.ItemModel) error {
	// Marshal []string ke JSON
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("Marshal photos error: %v", err)
		return err
	}

	query := `UPDATE items 
	          SET name = $2, description = $3, photos = $4, stock = $5, pickup_type = $6, 
	              price_per_day = $7, deposit = $8, discount = $9, category_id = $10, updated_at = NOW() 
	          WHERE id = $1`

	result, err := r.db.Exec(query, item.ID, item.Name, item.Description, photosJSON,
		item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount, item.CategoryID)

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
		log.Printf("Update: no rows affected for id %s", item.ID)
	}

	return nil
}

/*
Menghapus item berdasarkan ID.
Mengembalikan error jika penghapusan gagal atau tidak ada baris terpengaruh.
*/
func (r *itemRepository) Delete(id string) error {
	query := "DELETE FROM items WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Delete error: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Delete RowsAffected error: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Delete: no rows affected for id %s", id)
	}

	return nil
}
