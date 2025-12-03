package config

import (
	"fmt"
	"log"
	"os"

	"lalan-be/internal/message"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
DBConfig menyimpan koneksi database dan detail kredensial.
Digunakan oleh repository layer untuk berinteraksi dengan PostgreSQL.
*/
type DBConfig struct {
	DB      *sqlx.DB
	User    string
	Pass    string
	Host    string
	Port    string
	Name    string
	SSLMode string
}

/*
StorageConfig berisi semua konfigurasi untuk Supabase Storage (S3-compatible).
Digunakan oleh utils.Storage untuk upload, delete, dan presigned URL.
*/
type StorageConfig struct {
	AccessKey      string
	SecretKey      string
	Endpoint       string
	Region         string
	CustomerBucket string // Existing untuk KTP
	HosterBucket   string // Tambah untuk item
	ProjectID      string
	Domain         string
	ItemBucket     string // Map ke STORAGE_HOSTER_BUCKET
	Bucket         string // Tambah ini
}

/*
InitDatabase menginisialisasi koneksi ke PostgreSQL menggunakan sqlx.

Alur kerja:
1. Ambil kredensial wajib dari env via MustGetEnv
2. Bangun DSN (Data Source Name)
3. Buka koneksi dengan sqlx.Connect
4. Test koneksi dengan Ping
5. Atur connection pool (max open/idle, lifetime)
6. Return DBConfig siap pakai

Output sukses:
- (*DBConfig, nil) → koneksi berhasil dan siap digunakan
Output error:
- (nil, error) → gagal koneksi / ping / env tidak lengkap
*/
func InitDatabase() (*DBConfig, error) {
	user := MustGetEnv("DB_USER")
	pass := MustGetEnv("DB_PASSWORD")
	host := MustGetEnv("DB_HOST")
	port := MustGetEnv("DB_PORT")
	name := MustGetEnv("DB_NAME")

	ssl := os.Getenv("DB_SSL_MODE")
	if ssl == "" {
		ssl = "require" // Supabase & production wajib SSL
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, ssl,
	)

	log.Println("Connecting to PostgreSQL...")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", message.InternalError, err)
	}

	// Verifikasi koneksi aktif
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Optimasi connection pool untuk production
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60) // 5 menit

	log.Println("Database connected successfully")
	return &DBConfig{
		DB:      db,
		User:    user,
		Pass:    pass,
		Host:    host,
		Port:    port,
		Name:    name,
		SSLMode: ssl,
	}, nil
}

/*
LoadStorageConfig mengembalikan konfigurasi Supabase Storage dari environment.

Alur kerja:
1. Baca semua variabel wajib via MustGetEnv
2. Kembalikan struct StorageConfig lengkap

Output sukses:
- StorageConfig → semua field terisi
Output error:
- log.Fatal (panic) → jika ada env yang hilang (via MustGetEnv)
*/
func LoadStorageConfig() StorageConfig {
	cfg := StorageConfig{
		AccessKey:      MustGetEnv("STORAGE_ACCESS_KEY"),
		SecretKey:      MustGetEnv("STORAGE_SECRET_KEY"),
		Endpoint:       MustGetEnv("STORAGE_ENDPOINT"),
		Region:         MustGetEnv("STORAGE_REGION"),
		CustomerBucket: MustGetEnv("STORAGE_CUSTOMER_BUCKET"),
		HosterBucket:   MustGetEnv("STORAGE_HOSTER_BUCKET"),
		ProjectID:      MustGetEnv("STORAGE_PROJECT_ID"),
		Domain:         MustGetEnv("STORAGE_DOMAIN"),
		ItemBucket:     getEnv("STORAGE_HOSTER_BUCKET", "hoster"), // Map ke STORAGE_HOSTER_BUCKET
		Bucket:         getEnv("STORAGE_HOSTER_BUCKET", "hoster"), // Tambah ini
	}
	cfg.HosterBucket = getEnv("STORAGE_HOSTER_BUCKET", "hoster")
	cfg.CustomerBucket = getEnv("STORAGE_CUSTOMER_BUCKET", "customer")
	return cfg
}

// getEnv mengembalikan nilai env dengan default jika tidak ada
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
