package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
envLoaded adalah flag internal untuk memastikan file .env hanya dimuat sekali.
Mencegah pemanggilan godotenv.Load() berulang kali pada setiap GetEnv().
*/
var envLoaded bool

/*
GetEnv mengambil nilai environment variable dengan fallback.

Alur kerja:
1. Pastikan .env sudah dimuat (via LoadEnv)
2. Cek os.Getenv(key)
3. Jika kosong → kembalikan fallback

Output sukses:
- string nilai dari environment
- string fallback jika key tidak ada
*/
func GetEnv(key, fallback string) string {
	LoadEnv()
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

/*
GetJWTSecret mengembalikan secret key JWT dalam bentuk []byte.

Alur kerja:
1. Baca JWT_SECRET dari env
2. Validasi ketat di production:
  - Wajib ada
  - Minimal 32 karakter

3. Di development:
  - Jika kosong → pakai default + warning
  - Jika pendek → warning tapi tetap lanjut

Output sukses:
- []byte secret key yang aman digunakan untuk signing JWT
Output error:
- log.Fatal → aplikasi berhenti (hanya di production jika tidak valid)
*/
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	env := os.Getenv("APP_ENV")

	if env == "production" {
		if secret == "" {
			log.Fatal("FATAL: JWT_SECRET environment variable is required in production")
		}
		if len(secret) < 32 {
			log.Fatal("FATAL: JWT_SECRET must be at least 32 characters in production")
		}
	} else {
		// Development mode
		if secret == "" {
			log.Println("WARNING: JWT_SECRET not set, using development default (NOT FOR PRODUCTION!)")
			secret = "dev-secret-12345678901234567890" // 32 chars
		} else if len(secret) < 32 {
			log.Printf("WARNING: JWT_SECRET is weak (%d chars), recommended minimum 32 characters\n", len(secret))
		}
	}

	return []byte(secret)
}

/*
LoadEnv memuat file .env.dev secara otomatis jika bukan mode production.

Alur kerja:
1. Cek flag envLoaded → skip jika sudah pernah dimuat
2. Jika APP_ENV != "production" → load .env.dev
3. Set envLoaded = true

Output:
- Tidak ada return value
- Environment variables tersedia via os.Getenv()
- Tidak ada error (godotenv.Load bersifat idempotent & ignore missing file)
*/
func LoadEnv() {
	if envLoaded {
		return
	}

	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env.dev")
		log.Println("Loaded environment variables from .env.dev")
	}

	envLoaded = true
}

/*
MustGetEnv mengambil environment variable yang WAJIB ada.

Alur kerja:
1. LoadEnv() dipanggil otomatis
2. Cek nilai → jika kosong → log.Fatal (aplikasi langsung mati)

Digunakan untuk konfigurasi kritis: DB_URL, REDIS_URL, STORAGE credentials, dll.

Output sukses:
- string nilai environment
Output error:
- log.Fatal → aplikasi berhenti dengan pesan jelas
*/
func MustGetEnv(key string) string {
	LoadEnv()
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("FATAL: missing required environment variable: %s", key)
	}
	return v
}
