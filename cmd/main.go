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
Menginisialisasi aplikasi dan menjalankan server HTTP.
Hasil: Server berjalan di port tertentu.
*/
func main() {
	/*
		Mengatur konfigurasi database.
		Hasil: Koneksi database tersedia atau aplikasi berhenti.
	*/
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer cfg.DB.Close()

	/*
		Menginisialisasi komponen autentikasi.
		Hasil: Handler autentikasi siap digunakan.
	*/
	authRepo := repository.NewAuthRepository(cfg.DB)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	route.AuthRoutes(authHandler)

	/*
		Menginisialisasi komponen kategori.
		Hasil: Handler kategori siap digunakan.
	*/
	categoryRepo := repository.NewCategoryRepository(cfg.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	route.CategoryRoutes(categoryHandler)

	/*
		Menginisialisasi komponen item.
		Hasil: Handler item siap digunakan.
	*/
	itemRepo := repository.NewItemRepository(cfg.DB)
	itemService := service.NewItemService(itemRepo)
	itemHandler := handler.NewItemHandler(itemService)
	route.ItemRoutes(itemHandler)

	/*
		Menginisialisasi komponen syarat dan ketentuan.
		Hasil: Handler syarat dan ketentuan siap digunakan.
	*/
	tncRepo := repository.NewTermsAndConditionsRepository(cfg.DB)
	tncService := service.NewTermsAndConditionsService(tncRepo)
	tncHandler := handler.NewTermsAndConditionsHandler(tncService)
	route.TermsAndConditionsRoutes(tncHandler)

	/*
		Mengatur port server.
		Hasil: Port server ditentukan.
	*/
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s\n", port)
	/*
		Menjalankan server HTTP.
		Hasil: Server mendengarkan permintaan atau aplikasi berhenti.
	*/
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
