package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
Variabel untuk status pemuatan environment.
Menandai apakah environment sudah dimuat.
*/
var envLoaded bool

/*
GetEnv mengambil nilai environment dengan fallback.
Mengembalikan nilai atau fallback jika tidak ada.
*/
func GetEnv(key, fallback string) string {
	LoadEnv()
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

/*
GetJWTSecret mengambil rahasia JWT dari environment.
Mengembalikan sebagai byte slice dengan fallback default.
*/
func GetJWTSecret() []byte {
	secret := GetEnv("JWT_SECRET", "tesingdev")
	return []byte(secret)
}

/*
LoadEnv memuat environment dari file jika belum dimuat.
Hanya berjalan sekali per aplikasi.
*/
func LoadEnv() {
	if envLoaded {
		return
	}
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env.dev")
	}
	envLoaded = true
}

/*
MustGetEnv mengambil nilai environment yang wajib.
Menghentikan program jika tidak ada.
*/
func MustGetEnv(key string) string {
	LoadEnv()
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return v
}
