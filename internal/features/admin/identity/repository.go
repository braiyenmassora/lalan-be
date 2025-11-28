package identity

import (
	"fmt"
	"time"

	"lalan-be/internal/domain"

	"github.com/jmoiron/sqlx"
)

/*
AdminIdentityRepository mengatur akses database untuk fitur verifikasi identitas oleh admin.
Berisi query-query khusus admin (pending list, update status, detail per user).
*/
type AdminIdentityRepository struct {
	db *sqlx.DB
}

/*
NewAdminIdentityRepository membuat instance repository dengan koneksi database.

Output:
- *AdminIdentityRepository siap digunakan
*/
func NewAdminIdentityRepository(db *sqlx.DB) *AdminIdentityRepository {
	return &AdminIdentityRepository{db: db}
}

/*
GetPendingIdentities mengambil semua identitas yang berstatus 'pending'.

Alur kerja:
1. Query SELECT dengan filter status = 'pending'
2. Urutkan dari yang paling lama (ASC)

Output sukses:
- ([]*model.IdentityModel, nil)
Output error:
- (nil, error) → query gagal / koneksi DB bermasalah
*/
func (r *AdminIdentityRepository) GetPendingIdentities() ([]*domain.Identity, error) {
	var identities []*domain.Identity
	query := `
		SELECT 
			id, user_id, ktp_url, verified, status, reason, 
			verified_at, created_at, updated_at 
		FROM identity 
		WHERE status = 'pending' 
		ORDER BY created_at ASC
	`

	err := r.db.Select(&identities, query)
	return identities, err
}

/*
UpdateIdentityStatus memperbarui status verifikasi identitas (digunakan oleh service).

Alur kerja:
1. Jika status = approved → verified = true, verified_at = NOW()
2. Jika rejected → reason bisa diisi
3. Update kolom status, reason, verified, verified_at, updated_at

Output sukses:
- nil
Output error:
- error → query gagal / user_id tidak ditemukan
*/
func (r *AdminIdentityRepository) UpdateIdentityStatus(userID, status, reason string) error {
	var reasonArg interface{}
	if reason == "" {
		reasonArg = nil
	} else {
		reasonArg = reason
	}

	verified := status == "approved"

	query := `
		UPDATE identity
		SET 
			status = $1,
			reason = $2,
			verified = $3,
			verified_at = CASE WHEN $1 = 'approved' THEN NOW() ELSE NULL END,
			updated_at = NOW()
		WHERE user_id = $4::uuid
	`

	_, err := r.db.Exec(query, status, reasonArg, verified, userID)
	return err
}

/*
GetIdentityByUserID mengambil satu record identitas berdasarkan user_id.

Output sukses:
- (*model.IdentityModel, nil)
Output error:
- (nil, error) → record tidak ditemukan / query error
*/
func (r *AdminIdentityRepository) GetIdentityByUserID(userID string) (*domain.Identity, error) {
	var identity domain.Identity
	query := `
		SELECT 
			id, user_id, ktp_url, verified, status, reason, 
			verified_at, created_at, updated_at 
		FROM identity 
		WHERE user_id = $1 
		LIMIT 1
	`

	err := r.db.Get(&identity, query, userID)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

/*
ValidateIdentity memperbarui status verifikasi identitas (versi lengkap di repository).

Alur kerja:
1. Validasi status hanya boleh "approved" atau "rejected"
2. Set verified = true hanya jika approved
3. Kosongkan reason jika approved
4. Update semua field terkait + timestamp

Output sukses:
- nil
Output error:
- error → status tidak valid / query gagal
*/
func (r *AdminIdentityRepository) ValidateIdentity(identityID, status, reason string) error {
	now := time.Now()
	verified := false
	var verifiedAt *time.Time

	if status == "approved" {
		verified = true
		verifiedAt = &now
	} else if status == "rejected" {
		verified = false
		verifiedAt = nil
	} else {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Explicit types for parameters reduce Postgres ambiguity when some values
	// are NULL and prevent errors like "inconsistent types deduced for parameter $N".
	query := `
		UPDATE identity
		SET
			status = $1::text,
			reason = $2::text,
			verified = $3::boolean,
			verified_at = CASE WHEN $1::text = 'approved' THEN $4::timestamptz ELSE NULL END,
			updated_at = $5::timestamptz
		WHERE id = $6::uuid
	`

	_, err := r.db.Exec(query, status, reason, verified, verifiedAt, now, identityID)
	return err
}
