package auth

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/domain"
)

/*
AuthUser adalah proyeksi sederhana dari user untuk keperluan autentikasi.
Struct ini digunakan untuk menampung data user dari berbagai tabel (admin, hoster, customer)
dalam format yang seragam.
*/
type AuthUser struct {
	ID            string
	Email         string
	PasswordHash  string
	Role          string
	EmailVerified bool
}

/*
authRepository menangani interaksi database untuk fitur autentikasi.
Struct ini menyimpan koneksi database (*sqlx.DB) yang digunakan untuk query.
*/
type authRepository struct {
	db *sqlx.DB
}

/*
NewAuthRepository membuat instance repository baru.

Output:
- Pointer ke authRepository yang siap digunakan.
*/
func NewAuthRepository(db *sqlx.DB) *authRepository {
	return &authRepository{db: db}
}

/*
FindByEmail mencari user berdasarkan email di semua tabel role.

Fungsi ini akan mencari secara berurutan di tabel:
1. Admin
2. Hoster
3. Customer

Jika ditemukan, akan mengembalikan data user beserta role-nya.
Jika tidak ditemukan di semua tabel, mengembalikan nil.

Output:
- Pointer ke AuthUser jika ditemukan.
- nil jika tidak ditemukan.
- error jika terjadi kesalahan database.
*/
func (r *authRepository) FindByEmail(email string) (*AuthUser, error) {
	// 1. Cek tabel Admin
	var aid struct {
		ID           string `db:"id"`
		Email        string `db:"email"`
		PasswordHash string `db:"password_hash"`
	}
	queryAdmin := `SELECT id, email, password_hash FROM admin WHERE email = $1 LIMIT 1`
	err := r.db.Get(&aid, queryAdmin, email)
	if err == nil {
		log.Printf("authRepository: found admin %s", aid.ID)
		return &AuthUser{ID: aid.ID, Email: aid.Email, PasswordHash: aid.PasswordHash, Role: "admin"}, nil
	}
	// Pada titik ini, jika err==nil kita sudah return di blok sebelumnya,
	// jadi cukup cek apakah error bukan ErrNoRows.
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// 2. Cek tabel Hoster
	var hid struct {
		ID           string `db:"id"`
		Email        string `db:"email"`
		PasswordHash string `db:"password_hash"`
	}
	queryHoster := `SELECT id, email, password_hash FROM hoster WHERE email = $1 LIMIT 1`
	err = r.db.Get(&hid, queryHoster, email)
	if err == nil {
		log.Printf("authRepository: found hoster %s", hid.ID)
		return &AuthUser{ID: hid.ID, Email: hid.Email, PasswordHash: hid.PasswordHash, Role: "hoster"}, nil
	}
	// Sama seperti di atas: tidak perlu periksa err != nil lagi karena
	// kontrol aliran sudah memastikan err != nil jika kita sampai di sini.
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// 3. Cek tabel Customer
	var cid struct {
		ID            string `db:"id"`
		Email         string `db:"email"`
		PasswordHash  string `db:"password_hash"`
		EmailVerified bool   `db:"email_verified"`
	}
	queryCustomer := `SELECT id, email, password_hash, email_verified FROM customer WHERE email = $1 LIMIT 1`
	err = r.db.Get(&cid, queryCustomer, email)
	if err == nil {
		log.Printf("authRepository: found customer %s", cid.ID)
		return &AuthUser{ID: cid.ID, Email: cid.Email, PasswordHash: cid.PasswordHash, Role: "customer", EmailVerified: cid.EmailVerified}, nil
	}
	// Konsisten — jika bukan ErrNoRows, kembalikan error.
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Tidak ditemukan di manapun
	return nil, nil
}

/*
CreateCustomer menyimpan data customer baru ke database.

Fungsi ini melakukan insert data customer lengkap termasuk token verifikasi email.
Menggunakan RETURNING untuk mendapatkan ID dan timestamp yang digenerate database.

Output:
- error jika insert gagal (misal email duplikat).
- nil jika berhasil.
*/
func (r *authRepository) CreateCustomer(c *domain.Customer) error {
	query := `
		INSERT INTO customer (
			full_name,
			address,
			phone_number,
			email,
			password_hash,
			profile_photo,
			email_verified,
			verification_token,
			verification_expire,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		c.FullName, c.Address, c.PhoneNumber, c.Email, c.PasswordHash,
		c.ProfilePhoto, c.EmailVerified, c.VerificationToken, c.VerificationExpiresAt,
		c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		// Pesan error DB untuk constraint unique kadang berbeda antar DB/drivers,
		// contohnya "duplicate key value" atau hanya "duplicate". Cukup cek
		// "duplicate" sebagai substring tunggal untuk menutup kedua kasus ini
		// tanpa menulis kondisi redundan (tautological).
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("email already exists")
		}
		log.Printf("CreateCustomer (auth): %v", err)
		return err
	}
	return nil
}

/*
CreateHoster menyimpan data hoster baru ke database.

Fungsi ini melakukan insert data hoster lengkap. Hoster biasanya tidak
memerlukan verifikasi email di tahap awal (tergantung business logic).

Output:
- error jika insert gagal.
- nil jika berhasil.
*/
func (r *authRepository) CreateHoster(h *domain.Hoster) error {
	query := `
		INSERT INTO hoster (
			full_name,
			store_name,
			address,
			phone_number,
			email,
			password_hash,
			profile_photo,
			description,
			tiktok,
			instagram,
			website,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		h.FullName, h.StoreName, h.Address, h.PhoneNumber, h.Email,
		h.PasswordHash, h.ProfilePhoto, h.Description, h.Tiktok,
		h.Instagram, h.Website, h.CreatedAt, h.UpdatedAt,
	).Scan(&h.ID, &h.CreatedAt, &h.UpdatedAt)

	if err != nil {
		log.Printf("CreateHoster (auth): %v", err)
		return err
	}
	return nil
}

/*
CreateAdmin menyimpan data admin baru ke database.

Output:
- error jika insert gagal (misal duplikat email).
- nil jika berhasil.
*/
func (r *authRepository) CreateAdmin(a *domain.Admin) error {
	query := `
		INSERT INTO admin (
			email,
			password_hash,
			full_name,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		a.Email, a.PasswordHash, a.FullName, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("duplicate")
		}
		log.Printf("CreateAdmin (auth): %v", err)
		return err
	}
	return nil
}

/*
SendOTP memverifikasi OTP dan mengaktifkan akun customer.

Fungsi ini mencari customer dengan email dan token yang sesuai,
serta memastikan token belum kadaluarsa. Jika valid, status
email_verified diubah menjadi true.

Output:
- error jika OTP salah, kadaluarsa, atau terjadi kesalahan database.
- nil jika verifikasi berhasil.
*/
func (r *authRepository) SendOTP(email string, otp string) error {
	query := `
		UPDATE customer
		SET email_verified = true,
			verification_token = NULL,
			verification_expire = NULL,
			updated_at = NOW()
		WHERE email = $1 AND verification_token = $2 AND verification_expire > NOW()
	`
	res, err := r.db.Exec(query, email, otp)
	if err != nil {
		log.Printf("SendOTP (auth): error verifying email %s: %v", email, err)
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		log.Printf("SendOTP (auth): invalid or expired OTP for email %s", email)
		return errors.New("invalid or expired OTP")
	}
	return nil
}

/*
ResendOTP memperbarui token OTP untuk customer.

Fungsi ini mengupdate verification_token dan verification_expire
untuk customer dengan email tertentu.

Output:
- error jika customer tidak ditemukan atau update gagal.
- nil jika berhasil.
*/
func (r *authRepository) ResendOTP(email string, newOTP string, expireTime time.Time) error {
	query := `
		UPDATE customer 
		SET verification_token = $1, 
			verification_expire = $2, 
			updated_at = NOW() 
		WHERE email = $3
	`
	res, err := r.db.Exec(query, newOTP, expireTime, email)
	if err != nil {
		log.Printf("ResendOTP (auth): error updating OTP for email %s: %v", email, err)
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		log.Printf("ResendOTP (auth): no customer found for email %s", email)
		return errors.New("customer not found")
	}
	return nil
}
