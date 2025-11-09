package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*
Struct untuk menyimpan koneksi database.
*/
type Config struct {
	DB *sqlx.DB
}

/*
Mengambil nilai environment dengan fallback.
Mengembalikan string nilai atau fallback.
*/
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

/*
Menginisialisasi koneksi database PostgreSQL.
Mengembalikan pointer ke Config atau error jika gagal.
*/
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

/*
Mengambil secret JWT dari environment.
Mengembalikan byte slice secret.
*/
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default untuk development (WARNING: harus ada di .env production!)
		return []byte("testingfordev")
	}
	return []byte(secret)
}
