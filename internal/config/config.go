package config

import (
	"fmt"
	"lalan-be/internal/message"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
type Config
menyimpan parameter koneksi dan instance database untuk PostgreSQL
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
DatabaseConfig
menginisialisasi dan mengembalikan konfigurasi database jika koneksi berhasil
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
		return nil, fmt.Errorf("%s: %w", message.InternalError, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", message.InternalError, err)
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
