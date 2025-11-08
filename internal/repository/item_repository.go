package repository

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ItemRepository interface
type ItemRepository interface {
	FindItemNameByUserID(name string, userId string) (*model.ItemModel, error)
	CreateItem(item *model.ItemModel) error
	FindAll() ([]*model.ItemModel, error)
	FindByID(id string) (*model.ItemModel, error)
	FindByUserID(userID string) ([]*model.ItemModel, error)
	Update(item *model.ItemModel) error
	Delete(id string) error
}

// itemRepository struct
type itemRepository struct {
	db *sqlx.DB // koneksi database
}

// NewItemRepository constructor
func NewItemRepository(db *sqlx.DB) ItemRepository {
	return &itemRepository{db: db}
}

// FindItemNameByUserID method
func (r *itemRepository) FindItemNameByUserID(name string, userId string) (*model.ItemModel, error) {
	query := `SELECT id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at 
	          FROM items WHERE name = $1 AND user_id = $2 LIMIT 1`

	var item model.ItemModel
	err := r.db.QueryRow(query, name, userId).Scan(
		&item.ID, &item.Name, &item.Description, pq.Array(&item.Photos), &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindItemNameByUserID error: %v", err)
		return nil, err
	}

	return &item, nil
}

// CreateItem method
func (r *itemRepository) CreateItem(item *model.ItemModel) error {
	query := `INSERT INTO items (id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`

	_, err := r.db.Exec(query, item.ID, item.Name, item.Description, pq.Array(item.Photos),
		item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount,
		item.CategoryID, item.UserID)

	if err != nil {
		log.Printf("CreateItem error: %v", err)
		return err
	}

	return nil
}

// FindAll method
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
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, pq.Array(&item.Photos), &item.Stock,
			&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
			&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			log.Printf("FindAll scan error: %v", err)
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

// FindByID method
func (r *itemRepository) FindByID(id string) (*model.ItemModel, error) {
	query := `SELECT id, name, description, photos, stock, pickup_type, price_per_day, deposit, discount, category_id, user_id, created_at, updated_at 
	          FROM items WHERE id = $1 LIMIT 1`

	var item model.ItemModel
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &item.Description, pq.Array(&item.Photos), &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByID error: %v", err)
		return nil, err
	}

	return &item, nil
}

// FindByUserID method
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
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, pq.Array(&item.Photos), &item.Stock,
			&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
			&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			log.Printf("FindByUserID scan error: %v", err)
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

// Update method
func (r *itemRepository) Update(item *model.ItemModel) error {
	query := `UPDATE items 
	          SET name = $2, description = $3, photos = $4, stock = $5, pickup_type = $6, 
	              price_per_day = $7, deposit = $8, discount = $9, category_id = $10, updated_at = NOW() 
	          WHERE id = $1`

	result, err := r.db.Exec(query, item.ID, item.Name, item.Description, pq.Array(item.Photos),
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

// Delete method
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
