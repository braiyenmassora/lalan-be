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

/*
Menginisialisasi dan menjalankan aplikasi server.
*/
func main() {
	// Inisialisasi konfigurasi database dan koneksi; hentikan aplikasi jika gagal.
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer cfg.DB.Close()

	// Inisialisasi komponen autentikasi dengan dependency injection.
	authRepo := repository.NewAuthRepository(cfg.DB)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	route.AuthRoutes(authHandler)

	// Inisialisasi komponen kategori
	categoryRepo := repository.NewCategoryRepository(cfg.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	route.CategoryRoutes(categoryHandler)

	// Inisialisasi komponen item
	itemRepo := repository.NewItemRepository(cfg.DB)
	itemService := service.NewItemService(itemRepo)
	itemHandler := handler.NewItemHandler(itemService)
	route.ItemRoutes(itemHandler)

	// Konfigurasi port server dari environment atau gunakan default.
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s\n", port)
	// Jalankan server HTTP pada port tertentu; hentikan aplikasi jika gagal.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
