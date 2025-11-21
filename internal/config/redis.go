// internal/config/redis.go
package config

import (
	"context"
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

	// Debug: pastikan ENV kebaca di Render
	log.Println("REDIS_URL =", url)

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	opt.DialTimeout = 30 * time.Second
	opt.ReadTimeout = 10 * time.Second
	opt.WriteTimeout = 10 * time.Second

	Redis = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(RedisCtx, 30*time.Second)
	defer cancel()

	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
