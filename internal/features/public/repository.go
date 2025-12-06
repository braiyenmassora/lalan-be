package public

import (
	"encoding/json"
	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"log"
	"time"

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
			price_per_day, deposit, discount, category_id, hoster_id,
			created_at, updated_at
		FROM item
		WHERE is_hidden = false
		ORDER BY created_at DESC
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
			&item.CategoryID, &item.HosterID, &item.CreatedAt, &item.UpdatedAt,
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
GetItemDetail mengambil detail lengkap item dengan JOIN ke category, hoster, dan tnc.

Parameter:
- itemID: UUID item yang ingin diambil detail lengkapnya

Alur kerja:
1. Eksekusi query JOIN 3 tabel (item, category, hoster, tnc)
2. Manual scan semua field dari hasil JOIN
3. Unmarshal JSON untuk field photos (item) dan description (tnc)
4. Mapping hasil scan ke struct ItemDetailResponse

Output sukses:
- (*dto.ItemDetailResponse, nil)
Output error:
- (nil, sql.ErrNoRows) → item tidak ditemukan
- (nil, error) → query / scan / unmarshal gagal

SQL Query:
- JOIN item dengan category via category_id
- JOIN item dengan hoster via hoster_id
- LEFT JOIN tnc via hoster.id (LEFT karena tnc bisa NULL)
*/
func (r *publicRepository) GetItemDetail(itemID string) (*dto.ItemDetailResponse, error) {
	query := `
		SELECT
			i.id, i.name, i.description, i.photos, i.stock, i.pickup_type,
			i.price_per_day, i.deposit, i.discount, i.created_at AS item_created_at,
			i.updated_at AS item_updated_at,
			
			c.id AS category_id, c.name AS category_name, c.description AS category_description,
			
			h.id AS hoster_id, h.full_name, h.store_name, h.description AS hoster_description,
			h.phone_number, h.address, h.profile_photo, h.website, h.instagram, h.tiktok,
			
			t.description AS tnc_description
		FROM item i
		INNER JOIN category c ON c.id = i.category_id
		INNER JOIN hoster h ON h.id = i.hoster_id
		LEFT JOIN tnc t ON t.hoster_id = h.id
		WHERE i.id = $1 AND i.is_hidden = false
	`

	var (
		itemDetail         dto.ItemDetailResponse
		photosJSON         []byte
		tncDescriptionJSON []byte

		// Nullable fields
		hosterProfilePhoto, hosterWebsite, hosterInstagram, hosterTiktok *string
	)

	err := r.db.QueryRow(query, itemID).Scan(
		// Item fields
		&itemDetail.Item.ID,
		&itemDetail.Item.Name,
		&itemDetail.Item.Description,
		&photosJSON,
		&itemDetail.Item.Stock,
		&itemDetail.Item.PickupType,
		&itemDetail.Item.PricePerDay,
		&itemDetail.Item.Deposit,
		&itemDetail.Item.Discount,
		&itemDetail.Item.CreatedAt,
		&itemDetail.Item.UpdatedAt,

		// Category fields
		&itemDetail.Category.ID,
		&itemDetail.Category.Name,
		&itemDetail.Category.Description,

		// Hoster fields
		&itemDetail.Hoster.ID,
		&itemDetail.Hoster.FullName,
		&itemDetail.Hoster.StoreName,
		&itemDetail.Hoster.Description,
		&itemDetail.Hoster.PhoneNumber,
		&itemDetail.Hoster.Address,
		&hosterProfilePhoto,
		&hosterWebsite,
		&hosterInstagram,
		&hosterTiktok,

		// TnC field
		&tncDescriptionJSON,
	)

	if err != nil {
		log.Printf("GetItemDetail repository error: %v", err)
		return nil, err
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(photosJSON, &itemDetail.Item.Photos); err != nil {
		log.Printf("GetItemDetail unmarshal photos error: %v", err)
		return nil, err
	}

	if tncDescriptionJSON != nil {
		if err := json.Unmarshal(tncDescriptionJSON, &itemDetail.TermsAndConditions); err != nil {
			log.Printf("GetItemDetail unmarshal tnc description error: %v", err)
			return nil, err
		}
	} else {
		itemDetail.TermsAndConditions = []string{}
	}

	// Handle nullable fields
	if hosterProfilePhoto != nil {
		itemDetail.Hoster.ProfilePhoto = *hosterProfilePhoto
	}
	if hosterWebsite != nil {
		itemDetail.Hoster.Website = *hosterWebsite
	}
	if hosterInstagram != nil {
		itemDetail.Hoster.Instagram = *hosterInstagram
	}
	if hosterTiktok != nil {
		itemDetail.Hoster.Tiktok = *hosterTiktok
	}

	// Query booked dates untuk item ini
	bookedDates, err := r.getBookedDatesForItem(itemID)
	if err != nil {
		log.Printf("GetItemDetail getBookedDatesForItem error: %v", err)
		// Tidak return error, set empty array saja
		itemDetail.BookedDates = []string{}
	} else {
		itemDetail.BookedDates = bookedDates
	}

	return &itemDetail, nil
}

/*
getBookedDatesForItem mengambil semua tanggal yang sudah di-booking untuk item tertentu.

Parameter:
- itemID: UUID item yang ingin dicek availability-nya

Alur kerja:
1. Query booking dengan status pending/confirmed yang include item ini
2. Generate semua tanggal dari start_date sampai end_date untuk setiap booking
3. Return array tanggal dalam format YYYY-MM-DD

Output:
- ([]string, nil) - Array tanggal yang sudah di-booking
- ([]string{}, error) - Empty array jika error atau tidak ada booking
*/
func (r *publicRepository) getBookedDatesForItem(itemID string) ([]string, error) {
	query := `
		SELECT DISTINCT
			b.start_date,
			b.end_date
		FROM booking b
		INNER JOIN booking_item bi ON bi.booking_id = b.id
		WHERE bi.item_id = $1
		AND (
			b.status IN ('on_progress', 'on_rent', 'completed')
			OR (b.status = 'pending' AND b.locked_until > NOW())
		)
		AND b.end_date >= CURRENT_DATE
		ORDER BY b.start_date
	`

	rows, err := r.db.Query(query, itemID)
	if err != nil {
		log.Printf("getBookedDatesForItem query error: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	bookedDatesMap := make(map[string]bool)

	for rows.Next() {
		var startDate, endDate string
		if err := rows.Scan(&startDate, &endDate); err != nil {
			log.Printf("getBookedDatesForItem scan error: %v", err)
			continue
		}

		// Generate semua tanggal dari start_date sampai end_date
		dates := generateDateRange(startDate, endDate)
		for _, date := range dates {
			bookedDatesMap[date] = true
		}
	}

	// Convert map to array
	bookedDates := make([]string, 0, len(bookedDatesMap))
	for date := range bookedDatesMap {
		bookedDates = append(bookedDates, date)
	}

	return bookedDates, nil
}

/*
generateDateRange membuat array tanggal dari start sampai end (inclusive).

Parameter:
- startStr: Tanggal mulai dalam format string (YYYY-MM-DD atau timestamp)
- endStr: Tanggal akhir dalam format string (YYYY-MM-DD atau timestamp)

Output:
- []string - Array tanggal dalam format YYYY-MM-DD
*/
func generateDateRange(startStr, endStr string) []string {
	dates := []string{}

	// Parse start dan end date
	start, err := parseDate(startStr)
	if err != nil {
		log.Printf("generateDateRange parse start error: %v", err)
		return dates
	}

	end, err := parseDate(endStr)
	if err != nil {
		log.Printf("generateDateRange parse end error: %v", err)
		return dates
	}

	// Generate semua tanggal dari start sampai end
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates
}

/*
parseDate parsing string tanggal ke time.Time.
Support format: YYYY-MM-DD dan timestamp PostgreSQL.
*/
func parseDate(dateStr string) (time.Time, error) {
	// Try parse as YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr[:10])
	if err == nil {
		return t, nil
	}

	// Try parse as RFC3339 (PostgreSQL timestamp)
	t, err = time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}

/*
PublicRepository adalah kontrak untuk operasi data publik.
Digunakan oleh service layer untuk dependency injection.
*/
type PublicRepository interface {
	GetAllCategory() ([]*domain.Category, error)
	GetAllItems() ([]*domain.Item, error)
	GetAllTermsAndConditions() ([]*domain.TermsAndConditions, error)
	GetItemDetail(itemID string) (*dto.ItemDetailResponse, error)
}

/*
NewPublicRepository membuat instance repository dengan koneksi database.

Output:
- PublicRepository siap digunakan
*/
func NewPublicRepository(db *sqlx.DB) PublicRepository {
	return &publicRepository{db: db}
}
