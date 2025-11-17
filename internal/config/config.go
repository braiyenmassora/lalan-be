package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
Config menyimpan parameter koneksi dan instance database.
Digunakan untuk mengelola konfigurasi database PostgreSQL.
*/
type Config struct {
	DB      *sqlx.DB
	User    string
	Pass    string
	Host    string
	Port    string
	DBName  string
	SSLMode string
}

/*
DatabaseConfig menginisialisasi konfigurasi database.
Mengembalikan konfigurasi jika koneksi berhasil.
*/
func DatabaseConfig() (*Config, error) {
	user := MustGetEnv("DB_USER")
	pass := MustGetEnv("DB_PASSWORD")
	host := MustGetEnv("DB_HOST")
	port := MustGetEnv("DB_PORT")
	name := MustGetEnv("DB_NAME")
	ssl := os.Getenv("DB_SSL_MODE")
	if ssl == "" {
		ssl = "require"
	}
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, ssl,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed DB connect: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed DB ping: %w", err)
	}
	return &Config{
		DB:      db,
		User:    user,
		Pass:    pass,
		Host:    host,
		Port:    port,
		DBName:  name,
		SSLMode: ssl,
	}, nil
}
