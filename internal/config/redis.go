// config/redis.go
package config

import (
	"context"
	"crypto/tls"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
Redis
client Redis global untuk seluruh aplikasi
*/
var Redis *redis.Client

/*
RedisCtx
context default untuk operasi Redis
*/
var RedisCtx = context.Background()

/*
InitRedis
menginisialisasi koneksi Redis dari environment dan memverifikasi dengan ping
*/

func InitRedis() {
	// Baca dari .env
	host := MustGetEnv("REDIS_HOST")
	portStr := GetEnv("REDIS_PORT", "6379")
	username := GetEnv("REDIS_USERNAME", "")
	password := GetEnv("REDIS_PASSWORD", "")

	port, _ := strconv.Atoi(portStr)

	Redis = redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Username: username,
		Password: password,
		DB:       0,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // Render free tier butuh ini
		},
		DialTimeout:  20 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	_, err := Redis.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
