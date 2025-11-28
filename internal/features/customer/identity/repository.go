package identity

import (
	"database/sql"
	"errors"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"

	"github.com/jmoiron/sqlx"
)

/*
IdentityRepository adalah layer data access untuk fitur verifikasi identitas (KTP).
Hanya berisi operasi CRUD langsung ke tabel `identity` — tidak ada business rule.
*/
type IdentityRepository struct {
	db *sqlx.DB
}

/*
NewIdentityRepository membuat instance repository yang terhubung ke database.
Dependency injection memudahkan unit testing dengan mock DB.

Output:
- *IdentityRepository siap digunakan.
*/
func NewIdentityRepository(db *sqlx.DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

/*
UploadKTP menyimpan record KTP baru yang di-upload oleh customer (upload pertama kali).

Alur kerja:
1. Bangun model IdentityModel dengan status awal "pending"
2. Insert ke tabel identity

Output sukses:
- error = nil → insert berhasil
Output error:
- error → query gagal / constraint violation / DB error
*/
func (r *IdentityRepository) UploadKTP(req *dto.UploadIdentityByCustomerRequest) error {
	now := time.Now()
	identity := &domain.Identity{
		UserID:     req.UserID,
		KTPURL:     req.KTPURL,
		Verified:   false,
		Status:     "pending",
		Reason:     "",
		VerifiedAt: nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	query := `
		INSERT INTO identity (
			user_id, ktp_url, verified, status,
			reason, verified_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		identity.UserID,
		identity.KTPURL,
		identity.Verified,
		identity.Status,
		identity.Reason,
		identity.VerifiedAt,
		identity.CreatedAt,
		identity.UpdatedAt,
	)
	return err
}

// NOTE: We no longer update existing identity rows on re-upload. Every upload
// must create a new identity record so historical entries remain intact and
// any foreign keys referencing older identities (e.g., booking_identity_id_fkey)
// remain valid.

/*
GetStatusKTP mengambil record identity terakhir (most recent) milik user tertentu.
Jumlah record per user sekarang bisa lebih dari satu karena setiap upload
jadi entri baru. Oleh karena itu gunakan ORDER BY created_at DESC LIMIT 1
supaya selalu mendapatkan record terbaru.

Output sukses:
- (*model.IdentityModel, nil) → record ditemukan
- (nil, nil)                  → user belum pernah upload KTP
Output error:
- (nil, error)               → kesalahan database
*/
func (r *IdentityRepository) GetStatusKTP(userID string) (*domain.Identity, error) {
	var m domain.Identity
	query := `
		SELECT 
			id, user_id, ktp_url, verified, status, 
			reason, verified_at, created_at, updated_at
		FROM identity
		WHERE user_id = $1
		-- Logika baru: prioritaskan KTP yang paling baru di-upload oleh user
		-- (created_at DESC). Jika ada lebih dari satu entri dengan created_at
		-- yang sama, gunakan verified_at DESC sebagai tie-breaker sehingga jika
		-- salah satu entri sudah terverifikasi, yang terverifikasi paling baru akan
		-- diambil. Dengan urutan ini, endpoint status akan selalu menampilkan
		-- KTP hasil upload terakhir oleh user.
		ORDER BY created_at DESC, verified_at DESC NULLS LAST
		LIMIT 1
	`

	err := r.db.Get(&m, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
