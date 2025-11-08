package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Paket config untuk handle konfigurasi aplikasi dan koneksi database.

// Struktur Config untuk simpan instance koneksi database.
type Config struct {
	DB *sqlx.DB // Instance koneksi database
}

// Ambil nilai environment dengan fallback untuk konfigurasi aman.
func getEnv(key, fallback string) string {
	// Ambil variabel environment atau gunakan fallback
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Inisialisasi koneksi database PostgreSQL dan verifikasi konektivitas.
func DatabaseConfig() (*Config, error) {
	// Muat variabel environment dari .env.dev
	_ = godotenv.Load(".env.dev")

	// Ambil DSN database dari environment
	dsn := getEnv("DATABASE_URL", "")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Buat koneksi database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verifikasi konektivitas database
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Kembalikan instance konfigurasi
	return &Config{DB: db}, nil
}

// Ambil secret JWT dari environment untuk autentikasi aman.
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default untuk development (WARNING: harus ada di .env production!)
		return []byte("super-rahasia-123-ubah-nanti")
	}
	return []byte(secret)
}
