package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
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
	_ = godotenv.Load(".env.dev")

	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "")
	port := getEnv("DB_PORT", "")
	name := getEnv("DB_NAME", "")

	if user == "" || password == "" || host == "" || port == "" || name == "" {
		return nil, fmt.Errorf("DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME are required")
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, name)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Config{DB: db}, nil
}

/*
Test koneksi database menggunakan pgx untuk verifikasi tambahan.
*/
func TestDatabaseConnection(dsn string) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var result int
	err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		log.Fatalf("Failed to execute test query: %v", err)
	}

	log.Println("Database connection test successful")
}

/*
Mengambil secret JWT dari environment.
Mengembalikan byte slice secret.
*/
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("testingfordev")
	}
	return []byte(secret)
}
