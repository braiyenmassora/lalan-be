package main

import (
	"lalan-be/internal/config"
	"lalan-be/internal/handler"
	"lalan-be/internal/repository"
	"lalan-be/internal/route"
	"lalan-be/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	// Konfigurasi dan koneksi database
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer cfg.DB.Close()

	// Inisialisasi komponen autentikasi
	authRepo := repository.NewAuthRepository(cfg.DB)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	route.AuthRoutes(authHandler)

	// Inisialisasi komponen kategori
	categoryRepo := repository.NewCategoryRepository(cfg.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	route.CategoryRoute(categoryHandler)

	// Konfigurasi port server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s\n", port)
	// Menjalankan server HTTP
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
