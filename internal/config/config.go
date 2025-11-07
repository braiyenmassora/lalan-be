package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB *sqlx.DB // Instance koneksi database
}

func getEnv(key, fallback string) string {
	// Ambil variabel environment atau gunakan fallback
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

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
