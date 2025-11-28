package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
RedisClient adalah singleton instance dari redis.Client.
Digunakan bersama oleh seluruh aplikasi (cache, rate limit, OTP, session, dll).
Hanya ada satu koneksi Redis selama lifetime aplikasi.
*/
var RedisClient *redis.Client

/*
InitRedis menginisialisasi koneksi Redis dari environment REDIS_URL.

Alur kerja:
1. Baca REDIS_URL dari env (jika kosong → skip, Redis opsional)
2. Parse URL (support redis:// dan rediss://)
3. Set timeout dial/read/write untuk mencegah hanging
4. Buat client dan lakukan PING untuk verifikasi koneksi

Output sukses:
- error = nil → koneksi berhasil / Redis di-skip
Output error:
- error → URL invalid / timeout / server tidak respons
*/
func InitRedis() error {
	log.Println("Connecting to Redis...")

	url := GetEnv("REDIS_URL", "")
	if url == "" {
		log.Println("REDIS_URL not set, skipping Redis initialization")
		return nil // Redis bersifat opsional
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return fmt.Errorf("invalid Redis URL: %w", err)
	}

	// Timeout proteksi agar aplikasi tidak hang jika Redis down
	opt.DialTimeout = 10 * time.Second
	opt.ReadTimeout = 5 * time.Second
	opt.WriteTimeout = 5 * time.Second

	RedisClient = redis.NewClient(opt)

	// Test koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

/*
CloseRedis menutup koneksi Redis dengan graceful shutdown.

Alur kerja:
1. Cek apakah RedisClient ada
2. Panggil Close() dan log hasilnya

Output sukses:
- error = nil → koneksi ditutup bersih
Output error:
- error → gagal menutup koneksi (jarang terjadi)
*/
func CloseRedis() error {
	if RedisClient != nil {
		log.Println("Closing Redis connection...")
		if err := RedisClient.Close(); err != nil {
			return fmt.Errorf("failed to close Redis: %w", err)
		}
		log.Println("Redis connection closed")
	}
	return nil
}

/*
GetRedis mengembalikan instance Redis client yang sedang aktif.

Output:
- *redis.Client → client siap pakai
- nil → Redis belum di-init / gagal koneksi
*/
func GetRedis() *redis.Client {
	return RedisClient
}

/*
IsRedisAvailable mengecek apakah Redis sedang hidup dan responsif.

Alur kerja:
1. Cek apakah client sudah ada
2. Lakukan PING dengan timeout 1 detik

Output:
- true  → Redis aktif dan bisa dihubungi
- false → Redis mati / belum di-init / network error
*/
func IsRedisAvailable() bool {
	if RedisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := RedisClient.Ping(ctx).Err()
	return err == nil
}
