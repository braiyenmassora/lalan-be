package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB *sqlx.DB
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func DatabaseConfig() (*Config, error) {
	// load .env
	_ = godotenv.Load(".env.dev")

	dsn := getEnv("DATABASE_URL", "")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping untuk verifikasi (opsional)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Config{DB: db}, nil
}
