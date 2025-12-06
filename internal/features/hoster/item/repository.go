package item

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

/*
HosterItemRepository mendefinisikan operasi database khusus perspektif hoster.
Digunakan untuk dashboard item yang dimiliki hoster.
*/
type HosterItemRepository interface {
	GetListItem(hosterID string) ([]dto.ItemListByHosterResponse, error)
	GetItemDetail(hosterID, itemID string) (*dto.ItemDetailByHosterResponse, error)
	CreateItem(item *domain.Item) (*dto.ItemDetailByHosterResponse, error)
	DeleteItem(hosterID, itemID string) error
	GetItemPhotos(itemID string) ([]string, error) // Return array URL photos
	UpdateItem(hosterID, itemID string, req *dto.UpdateItemRequestRequest) error
	GetCategory() ([]dto.CategoryResponse, error)                  // Get all categories for dropdown
	HasActiveBookings(itemID string) (bool, error)                 // Cek apakah item punya booking aktif
	UpdateVisibility(hosterID, itemID string, isHidden bool) error // Toggle visibility item
}

/*
hosterItemRepository adalah implementasi repository untuk hoster-facing item.
*/
type hosterItemRepository struct {
	db *sqlx.DB
}

/*
NewHosterItemRepository membuat instance repository dengan koneksi database.

Output:
- HosterItemRepository siap digunakan
*/
func NewHosterItemRepository(db *sqlx.DB) HosterItemRepository {
	return &hosterItemRepository{db: db}
}

/*
GetListItem mengambil ringkasan semua item yang dimiliki hoster.

Alur kerja:
1. Query item dengan filter hoster_id
2. Order by name untuk kemudahan di frontend

Output sukses:
- ([]dto.ItemListByHosterResponse, nil)
Output error:
- (nil, error) → query gagal
*/
func (r *hosterItemRepository) GetListItem(hosterID string) ([]dto.ItemListByHosterResponse, error) {
	query := `
        SELECT
            id,
            name,
            stock,
            price_per_day,
            pickup_type,
            is_hidden
        FROM item
        WHERE hoster_id = $1
        ORDER BY created_at DESC
    `

	var items []dto.ItemListByHosterResponse
	err := r.db.Select(&items, query, hosterID)
	if err != nil {
		log.Printf("GetListItem: database error for hoster %s: %v", hosterID, err)
		return nil, err
	}

	return items, nil
}

/*
GetItemDetail mengambil detail lengkap item berdasarkan ID.

Alur kerja:
1. Query item dengan filter id dan hoster_id (untuk memastikan ownership)
2. Unmarshal photos dari JSONB
3. Map ke DTO response

Output sukses:
- (*dto.ItemDetailByHosterResponse, nil)
Output error:
- (nil, error) → item tidak ditemukan atau bukan milik hoster
*/
func (r *hosterItemRepository) GetItemDetail(hosterID, itemID string) (*dto.ItemDetailByHosterResponse, error) {
	var (
		detail dto.ItemDetailByHosterResponse
		row    struct {
			ID          string          `db:"id"`
			Name        string          `db:"name"`
			Description sql.NullString  `db:"description"`
			Photos      json.RawMessage `db:"photos"`
			Stock       int             `db:"stock"`
			PickupType  string          `db:"pickup_type"`
			PricePerDay int             `db:"price_per_day"`
			Deposit     int             `db:"deposit"`
			Discount    sql.NullInt64   `db:"discount"`
			CreatedAt   sql.NullTime    `db:"created_at"`
			UpdatedAt   sql.NullTime    `db:"updated_at"`
			CategoryID  string          `db:"category_id"`
			HosterID    string          `db:"hoster_id"`
			IsHidden    bool            `db:"is_hidden"`
		}
	)

	query := `
		SELECT id, name, description, photos, stock, pickup_type,
		       price_per_day, deposit, discount, created_at, updated_at, category_id, hoster_id, is_hidden
		FROM item 
		WHERE id = $1 AND hoster_id = $2
	`

	err := r.db.Get(&row, query, itemID, hosterID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetItemDetail: item %s not found for hoster %s", itemID, hosterID)
		} else {
			log.Printf("GetItemDetail: database error for item %s hoster %s: %v", itemID, hosterID, err)
		}
		return nil, err
	}

	// Convert nullable fields
	if row.Description.Valid {
		detail.Description = row.Description.String
	}
	if row.Discount.Valid {
		detail.Discount = int(row.Discount.Int64)
	}
	if row.CreatedAt.Valid {
		detail.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		detail.UpdatedAt = row.UpdatedAt.Time
	}

	// Unmarshal photos
	if len(row.Photos) > 0 {
		var photos []string
		if err = json.Unmarshal(row.Photos, &photos); err != nil {
			log.Printf("GetItemDetail: failed to unmarshal photos for item %s: %v", itemID, err)
			return nil, err
		}
		detail.Photos = photos
	}

	// Map remaining fields
	detail.ID = row.ID
	detail.Name = row.Name
	detail.Stock = row.Stock
	detail.PickupType = dto.PickupMethod(row.PickupType)
	detail.PricePerDay = row.PricePerDay
	detail.Deposit = row.Deposit
	detail.CategoryID = row.CategoryID
	detail.HosterID = row.HosterID
	detail.IsHidden = row.IsHidden

	return &detail, nil
}

/*
CreateItem menyimpan item baru ke dalam database.

Alur kerja:
1. Menyusun query INSERT untuk menambah item
2. Menjalankan query dengan parameter yang diberikan
3. Mengembalikan ID item yang baru dibuat

Output sukses:
- (id_item_baru, nil)
Output error:
- ("", error) → query gagal
*/
func (r *hosterItemRepository) CreateItem(item *domain.Item) (*dto.ItemDetailByHosterResponse, error) {
	// mulai transaction
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("CreateItem: error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO item (
            id, hoster_id, name, description, photos, stock, pickup_type,
            price_per_day, deposit, discount, category_id, created_at, updated_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW())
    `
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("CreateItem: failed to marshal photos for item %s: %v", item.ID, err)
		return nil, err
	}

	_, err = tx.Exec(query,
		item.ID, item.HosterID, item.Name, item.Description, photosJSON,
		item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount,
		item.CategoryID,
	)
	if err != nil {
		log.Printf("CreateItem: error inserting item %s: %v", item.ID, err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("CreateItem: error committing transaction: %v", err)
		return nil, err
	}

	// ambil detail item yang baru dibuat dan return
	// jsonb `photos` returns []byte from driver — scan into json.RawMessage then unmarshal into []string
	var (
		detail dto.ItemDetailByHosterResponse
		row    struct {
			ID          string          `db:"id"`
			Name        string          `db:"name"`
			Description sql.NullString  `db:"description"`
			Photos      json.RawMessage `db:"photos"`
			Stock       int             `db:"stock"`
			PickupType  string          `db:"pickup_type"`
			PricePerDay int             `db:"price_per_day"`
			Deposit     int             `db:"deposit"`
			Discount    sql.NullInt64   `db:"discount"`
			CreatedAt   sql.NullTime    `db:"created_at"`
			UpdatedAt   sql.NullTime    `db:"updated_at"`
			CategoryID  string          `db:"category_id"`
			HosterID    string          `db:"hoster_id"`
		}
	)

	getQuery := `
        SELECT id, name, description, photos, stock, pickup_type,
               price_per_day, deposit, discount, created_at, updated_at, category_id, hoster_id
        FROM item WHERE id = $1
    `
	err = r.db.Get(&row, getQuery, item.ID)
	if err != nil {
		log.Printf("CreateItem: failed to load created item detail %s: %v", item.ID, err)
		return nil, err
	}

	// convert nullable fields + unmarshal photos
	if row.Description.Valid {
		detail.Description = row.Description.String
	}
	if row.Discount.Valid {
		detail.Discount = int(row.Discount.Int64)
	}
	if row.CreatedAt.Valid {
		detail.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		detail.UpdatedAt = row.UpdatedAt.Time
	}
	if len(row.Photos) > 0 {
		var photos []string
		if err = json.Unmarshal(row.Photos, &photos); err != nil {
			log.Printf("CreateItem: failed to unmarshal photos for item %s: %v", item.ID, err)
			return nil, err
		}
		detail.Photos = photos
	}

	detail.ID = row.ID

	detail.Name = row.Name
	detail.Stock = row.Stock
	detail.PickupType = dto.PickupMethod(row.PickupType)
	detail.PricePerDay = row.PricePerDay
	detail.Deposit = row.Deposit
	detail.CategoryID = row.CategoryID
	detail.HosterID = row.HosterID

	return &detail, nil
}

/*
DeleteItem menghapus item milik hoster.
Alur:
1. Mulai transaction
2. DELETE dari tabel item dengan filter id (diasumsikan milik hoster yang login)
3. Jika tidak ada baris yang terhapus -> kembalikan sql.ErrNoRows
4. Commit bila sukses
*/
func (r *hosterItemRepository) DeleteItem(hosterID, itemID string) error {
	// mulai transaction
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("DeleteItem: error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// cek pemilik item dulu (debug & keamanan)
	var ownerID string
	if err := tx.Get(&ownerID, `SELECT hoster_id FROM item WHERE id = $1`, itemID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DeleteItem: item %s not found", itemID)
			return sql.ErrNoRows
		}
		log.Printf("DeleteItem: failed to query owner for item %s: %v", itemID, err)
		return err
	}
	if ownerID != hosterID {
		log.Printf("DeleteItem: ownership mismatch for item %s owner=%s requester=%s", itemID, ownerID, hosterID)
		return sql.ErrNoRows
	}

	query := `DELETE FROM item WHERE id = $1 AND hoster_id = $2`
	res, err := tx.Exec(query, itemID, hosterID)
	if err != nil {
		log.Printf("DeleteItem: error deleting item %s for hoster %s: %v", itemID, hosterID, err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("DeleteItem: failed to get rows affected for item %s: %v", itemID, err)
		return err
	}
	log.Printf("DeleteItem: rows affected = %d for item %s hoster %s", rows, itemID, hosterID)
	if rows == 0 {
		// tidak ada item yang terhapus -> kemungkinan tidak ada atau bukan milik hoster
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		log.Printf("DeleteItem: error committing transaction for item %s: %v", itemID, err)
		return err
	}

	return nil
}

/*
GetItemPhotos mengambil daftar URL foto untuk item tertentu.

Alur kerja:
1. Query item dengan filter id
2. Mengembalikan array URL foto

Output sukses:
- ([]string, nil)
Output error:
- (nil, error) → query gagal
*/
func (r *hosterItemRepository) GetItemPhotos(itemID string) ([]string, error) {
	var photosJSON string
	query := `SELECT photos FROM item WHERE id = $1`
	err := r.db.Get(&photosJSON, query, itemID)
	if err != nil {
		return nil, err
	}

	var photos []string
	if err := json.Unmarshal([]byte(photosJSON), &photos); err != nil {
		return nil, err
	}
	return photos, nil
}

/*
UpdateItem mengubah data item yang dimiliki hoster.

Alur kerja:
1. Mulai transaction
2. Cek kepemilikan item
3. Siapkan query UPDATE dinamis sesuai field yang diubah
4. Jalankan query
5. Commit transaction

Output sukses:
- nil
Output error:
- error → gagal di tengah jalan
*/
func (r *hosterItemRepository) UpdateItem(hosterID, itemID string, req *dto.UpdateItemRequestRequest) error {
	// Mulai transaction
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("UpdateItem: error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// Cek ownership
	var ownerID string
	if err := tx.Get(&ownerID, `SELECT hoster_id FROM item WHERE id = $1`, itemID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("UpdateItem: item %s not found", itemID)
			return sql.ErrNoRows
		}
		log.Printf("UpdateItem: failed to query owner for item %s: %v", itemID, err)
		return err
	}
	if ownerID != hosterID {
		log.Printf("UpdateItem: ownership mismatch for item %s owner=%s requester=%s", itemID, ownerID, hosterID)
		return sql.ErrNoRows
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Stock != nil {
		setParts = append(setParts, fmt.Sprintf("stock = $%d", argIndex))
		args = append(args, *req.Stock)
		argIndex++
	}
	if req.PickupType != nil {
		setParts = append(setParts, fmt.Sprintf("pickup_type = $%d", argIndex))
		args = append(args, *req.PickupType)
		argIndex++
	}
	if req.Deposit != nil {
		setParts = append(setParts, fmt.Sprintf("deposit = $%d", argIndex))
		args = append(args, *req.Deposit)
		argIndex++
	}
	if req.Discount != nil {
		setParts = append(setParts, fmt.Sprintf("discount = $%d", argIndex))
		args = append(args, *req.Discount)
		argIndex++
	}
	if req.CategoryID != nil {
		setParts = append(setParts, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.PricePerDay != nil {
		setParts = append(setParts, fmt.Sprintf("price_per_day = $%d", argIndex))
		args = append(args, *req.PricePerDay)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE item SET %s, updated_at = NOW() WHERE id = $%d AND hoster_id = $%d",
		strings.Join(setParts, ", "), argIndex, argIndex+1)
	args = append(args, itemID, hosterID)

	_, err = tx.Exec(query, args...)
	if err != nil {
		log.Printf("UpdateItem: error updating item %s: %v", itemID, err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("UpdateItem: error committing transaction for item %s: %v", itemID, err)
		return err
	}

	log.Printf("UpdateItem: successfully updated item %s for hoster %s", itemID, hosterID)
	return nil
}

/*
GetCategory mengambil semua kategori yang aktif untuk dropdown saat create/update item.

Output:
- ([]dto.CategoryResponse, nil) - List semua kategori
- (nil, error) - Query gagal
*/
func (r *hosterItemRepository) GetCategory() ([]dto.CategoryResponse, error) {
	query := `
		SELECT id, name, description
		FROM category
		ORDER BY name ASC
	`

	var categories []dto.CategoryResponse
	if err := r.db.Select(&categories, query); err != nil {
		log.Printf("GetCategory: db error err=%v", err)
		return nil, err
	}

	log.Printf("GetCategory: found %d categories", len(categories))
	return categories, nil
}

/*
HasActiveBookings memeriksa apakah item memiliki booking aktif.
Booking aktif adalah booking dengan status: pending, on_progress, atau on_rent.

Output:
- (true, nil) - Item memiliki booking aktif
- (false, nil) - Item tidak memiliki booking aktif
- (false, error) - Query gagal
*/
func (r *hosterItemRepository) HasActiveBookings(itemID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM booking_item bi
			JOIN booking b ON bi.booking_id = b.id
			WHERE bi.item_id = $1 
			AND b.status IN ('pending', 'on_progress', 'on_rent')
		) AS has_bookings
	`

	var hasBookings bool
	if err := r.db.Get(&hasBookings, query, itemID); err != nil {
		log.Printf("HasActiveBookings: db error item=%s err=%v", itemID, err)
		return false, err
	}

	log.Printf("HasActiveBookings: item=%s has_bookings=%v", itemID, hasBookings)
	return hasBookings, nil
}

/*
UpdateVisibility mengubah status visibility item (is_hidden).

Alur kerja:
1. Mulai transaction
2. Cek ownership item
3. Update field is_hidden
4. Commit transaction

Output:
- nil - Sukses
- error - Gagal (item not found atau bukan milik hoster)
*/
func (r *hosterItemRepository) UpdateVisibility(hosterID, itemID string, isHidden bool) error {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("UpdateVisibility: error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// Cek ownership
	var ownerID string
	if err := tx.Get(&ownerID, `SELECT hoster_id FROM item WHERE id = $1`, itemID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("UpdateVisibility: item %s not found", itemID)
			return sql.ErrNoRows
		}
		log.Printf("UpdateVisibility: failed to query owner for item %s: %v", itemID, err)
		return err
	}
	if ownerID != hosterID {
		log.Printf("UpdateVisibility: ownership mismatch for item %s owner=%s requester=%s", itemID, ownerID, hosterID)
		return sql.ErrNoRows
	}

	// Update visibility
	query := `UPDATE item SET is_hidden = $1, updated_at = NOW() WHERE id = $2 AND hoster_id = $3`
	_, err = tx.Exec(query, isHidden, itemID, hosterID)
	if err != nil {
		log.Printf("UpdateVisibility: error updating item %s: %v", itemID, err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("UpdateVisibility: error committing transaction for item %s: %v", itemID, err)
		return err
	}

	log.Printf("UpdateVisibility: successfully updated item %s is_hidden=%v", itemID, isHidden)
	return nil
}
