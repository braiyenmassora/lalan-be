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
Config
menyimpan koneksi database dan parameter konfigurasi PostgreSQL
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
membuat koneksi ke PostgreSQL menggunakan variabel environment dan mengembalikan konfigurasi database
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

	log.Println("Database connected successfully")

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
