package profile

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/dto"
)

/*
HosterProfileRepository mendefinisikan operasi database untuk profile hoster.
*/
type HosterProfileRepository interface {
	GetProfile(hosterID string) (*dto.HosterProfileResponse, error)
	UpdateProfile(hosterID string, req *dto.UpdateHosterProfileRequest) error
}

/*
hosterProfileRepository adalah implementasi repository untuk profile hoster.
*/
type hosterProfileRepository struct {
	db *sqlx.DB
}

/*
NewHosterProfileRepository membuat instance repository dengan koneksi database.

Output:
- HosterProfileRepository siap digunakan
*/
func NewHosterProfileRepository(db *sqlx.DB) HosterProfileRepository {
	return &hosterProfileRepository{db: db}
}

/*
GetProfile mengambil profil hoster berdasarkan hosterID.

Alur kerja:
1. Query hoster dari database
2. Hitung jumlah hari sejak bergabung
3. Map ke DTO response

Output sukses:
- (*dto.HosterProfileResponse, nil)
Output error:
- (nil, error) → hoster tidak ditemukan
*/
func (r *hosterProfileRepository) GetProfile(hosterID string) (*dto.HosterProfileResponse, error) {
	var row struct {
		ID           string       `db:"id"`
		FullName     string       `db:"full_name"`
		Email        string       `db:"email"`
		PhoneNumber  string       `db:"phone_number"`
		Address      string       `db:"address"`
		StoreName    string       `db:"store_name"`
		Description  string       `db:"description"`
		Website      string       `db:"website"`
		Instagram    string       `db:"instagram"`
		Tiktok       string       `db:"tiktok"`
		ProfilePhoto string       `db:"profile_photo"`
		CreatedAt    sql.NullTime `db:"created_at"`
	}

	query := `
		SELECT 
			id, full_name, email, phone_number, address,
			store_name, description, website, instagram, tiktok,
			profile_photo, created_at
		FROM hoster
		WHERE id = $1
	`

	err := r.db.Get(&row, query, hosterID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetProfile: hoster not found %s", hosterID)
		} else {
			log.Printf("GetProfile: database error for hoster %s: %v", hosterID, err)
		}
		return nil, err
	}

	// Hitung jumlah hari sejak bergabung
	var joinedAt time.Time
	var daysSinceJoin int
	if row.CreatedAt.Valid {
		joinedAt = row.CreatedAt.Time
		daysSinceJoin = int(time.Since(joinedAt).Hours() / 24)
	}

	// Map ke response
	response := &dto.HosterProfileResponse{
		ID:            row.ID,
		FullName:      row.FullName,
		Email:         row.Email,
		PhoneNumber:   row.PhoneNumber,
		Address:       row.Address,
		StoreName:     row.StoreName,
		Description:   row.Description,
		Website:       row.Website,
		Instagram:     row.Instagram,
		Tiktok:        row.Tiktok,
		ProfilePhoto:  row.ProfilePhoto,
		JoinedAt:      joinedAt,
		DaysSinceJoin: daysSinceJoin,
	}

	return response, nil
}

/*
UpdateProfile memperbarui profil hoster.
Update: address, phone_number, description, website, instagram, tiktok.

Alur kerja:
1. Update hoster dengan filter hosterID
2. Return error jika tidak ada row yang di-update

Output sukses:
- nil
Output error:
- error → update gagal atau hoster tidak ditemukan
*/
func (r *hosterProfileRepository) UpdateProfile(hosterID string, req *dto.UpdateHosterProfileRequest) error {
	query := `
		UPDATE hoster
		SET 
			address = $1, 
			phone_number = $2, 
			description = $3,
			website = $4,
			instagram = $5,
			tiktok = $6,
			updated_at = NOW()
		WHERE id = $7
	`

	res, err := r.db.Exec(query,
		req.Address,
		req.PhoneNumber,
		req.Description,
		req.Website,
		req.Instagram,
		req.Tiktok,
		hosterID,
	)
	if err != nil {
		log.Printf("UpdateProfile: error updating hoster %s: %v", hosterID, err)
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("UpdateProfile: no rows affected for hoster %s", hosterID)
		return sql.ErrNoRows
	}

	return nil
}
