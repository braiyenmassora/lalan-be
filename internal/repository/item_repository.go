package repository

import (
	"database/sql"
	"encoding/json"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

// ItemRepository interface
type ItemRepository interface {
	FindItemNameByUserID(name string, userId string) (*model.ItemModel, error)
	CreateItem(item *model.ItemModel) error
}

// itemRepository struct
type itemRepository struct {
	db *sqlx.DB // koneksi database
}

// NewItemRepository constructur
func NewItemRepository(db *sqlx.DB) ItemRepository {
	return &itemRepository{db: db}
}

// FindItemNameByUserID method
func (r *itemRepository) FindItemNameByUserID(name string, userId string) (*model.ItemModel, error) {
	// cari item berdasarkan nama dan userid
	query := "SELECT id,name,description,photos,stock,pickup_type,price_per_day,deposit,discount,category_id,user_id,created_at,updated_at FROM items WHERE name = $1 AND user_id = $2"

	row := r.db.QueryRow(query, name, userId)
	var item model.ItemModel
	var photosJSON string
	err := row.Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock, &item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount, &item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Printf("FindItemNameByUserID, no rows for name and userId %s %s", name, userId)
		return nil, nil
	}
	if err != nil {
		log.Printf("FindItemNameByUserID: scan error %v", err)
		return nil, err
	}
	// Unmarshal photos from JSON
	if err := json.Unmarshal([]byte(photosJSON), &item.Photos); err != nil {
		log.Printf("FindItemNameByUserID: unmarshal photos error %v", err)
		return nil, err
	}
	log.Printf("FindItemNameByUserID: found item %s", item.Name, item.UserID)
	return &item, nil
}

// CreateItem method
func (r *itemRepository) CreateItem(item *model.ItemModel) error {
	// insert data item
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("CreateItem: marshal photos error %v", err)
		return err
	}
	query := "INSERT INTO items(id,name,description,photos,stock,pickup_type,price_per_day,deposit,discount,category_id,user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err = r.db.Exec(query, item.ID, item.Name, item.Description, string(photosJSON), item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount, item.CategoryID, item.UserID)
	if err != nil {
		log.Printf("CreateItem: insert error %v", err)
		return err
	}
	return nil
}
