package tnc

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
)

/*
TnCRepository mendefinisikan operasi database untuk Terms and Conditions.
*/
type TnCRepository interface {
	CreateTnC(tnc *domain.TermsAndConditions) error
	UpdateTnC(tncID, hosterID string, description []string) error
	GetTnCByHosterID(hosterID string) (*dto.TnCResponse, error)
}

/*
tncRepository adalah implementasi repository untuk TnC.
*/
type tncRepository struct {
	db *sqlx.DB
}

/*
NewTnCRepository membuat instance repository dengan koneksi database.

Output:
- TnCRepository siap digunakan
*/
func NewTnCRepository(db *sqlx.DB) TnCRepository {
	return &tncRepository{db: db}
}

/*
CreateTnC menyimpan T&C baru ke database.

Alur kerja:
1. Marshal description array ke JSONB
2. Insert ke tabel tnc
3. Return error jika gagal

Output sukses:
- nil
Output error:
- error → insert gagal
*/
func (r *tncRepository) CreateTnC(tnc *domain.TermsAndConditions) error {
	descriptionJSON, err := json.Marshal(tnc.Description)
	if err != nil {
		log.Printf("CreateTnC: failed to marshal description: %v", err)
		return err
	}

	query := `
		INSERT INTO tnc (id, hoster_id, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`

	_, err = r.db.Exec(query, tnc.ID, tnc.UserID, descriptionJSON)
	if err != nil {
		log.Printf("CreateTnC: error inserting tnc %s: %v", tnc.ID, err)
		return err
	}

	return nil
}

/*
UpdateTnC memperbarui description T&C.

Alur kerja:
1. Marshal description array ke JSONB
2. Update dengan filter tncID dan userID (ownership check)
3. Return error jika tidak ada row yang di-update

Output sukses:
- nil
Output error:
- error → update gagal atau TnC tidak ditemukan
*/
func (r *tncRepository) UpdateTnC(tncID, hosterID string, description []string) error {
	descriptionJSON, err := json.Marshal(description)
	if err != nil {
		log.Printf("UpdateTnC: failed to marshal description: %v", err)
		return err
	}

	query := `
		UPDATE tnc
		SET description = $1, updated_at = NOW()
		WHERE id = $2 AND hoster_id = $3
	`

	res, err := r.db.Exec(query, descriptionJSON, tncID, hosterID)
	if err != nil {
		log.Printf("UpdateTnC: error updating tnc %s: %v", tncID, err)
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("UpdateTnC: no rows affected for tnc %s hoster %s", tncID, hosterID)
		return sql.ErrNoRows
	}

	return nil
}

/*
GetTnCByHosterID mengambil T&C berdasarkan hoster_id.

Alur kerja:
1. Query tnc dengan filter hoster_id
2. Unmarshal description dari JSONB
3. Map ke DTO response

Output sukses:
- (*dto.TnCResponse, nil)
Output error:
- (nil, error) → TnC tidak ditemukan
*/
func (r *tncRepository) GetTnCByHosterID(hosterID string) (*dto.TnCResponse, error) {
	var (
		response dto.TnCResponse
		row      struct {
			ID          string          `db:"id"`
			HosterID    string          `db:"hoster_id"`
			Description json.RawMessage `db:"description"`
			CreatedAt   sql.NullTime    `db:"created_at"`
			UpdatedAt   sql.NullTime    `db:"updated_at"`
		}
	)

	query := `
		SELECT id, hoster_id, description, created_at, updated_at
		FROM tnc
		WHERE hoster_id = $1
		LIMIT 1
	`

	err := r.db.Get(&row, query, hosterID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetTnCByHosterID: tnc not found for hoster %s", hosterID)
		} else {
			log.Printf("GetTnCByHosterID: database error for hoster %s: %v", hosterID, err)
		}
		return nil, err
	}

	// Unmarshal description
	if len(row.Description) > 0 {
		var description []string
		if err = json.Unmarshal(row.Description, &description); err != nil {
			log.Printf("GetTnCByHosterID: failed to unmarshal description: %v", err)
			return nil, err
		}
		response.Description = description
	}

	// Map fields
	response.ID = row.ID
	response.HosterID = row.HosterID
	if row.CreatedAt.Valid {
		response.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		response.UpdatedAt = row.UpdatedAt.Time
	}

	return &response, nil
}
