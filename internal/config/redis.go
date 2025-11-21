// internal/config/redis.go
package config

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var RedisCtx = context.Background()

/*
InitRedis
menginisialisasi koneksi Redis dengan timeout tinggi & TLS fleksibel untuk Render
*/
func InitRedis() {
	url := MustGetEnv("REDIS_URL")

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	// FIX 100% UNTUK RENDER REDIS
	opt.DialTimeout = 30 * time.Second // dari 5s → 30s
	opt.ReadTimeout = 10 * time.Second
	opt.WriteTimeout = 10 * time.Second
	opt.TLSConfig = &tls.Config{
		InsecureSkipVerify: true, // WAJIB TRUE untuk Render Redis internal (cert issue)
		MinVersion:         tls.VersionTLS12,
	}

	Redis = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(RedisCtx, 30*time.Second)
	defer cancel()

	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
