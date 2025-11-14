package main

import (
	"log"
	"net/http"

	"lalan-be/internal/config"
	"lalan-be/internal/features/admin"

	publicHandler "lalan-be/internal/features/public/handler"
	publicRepo "lalan-be/internal/features/public/repository"
	publicRoute "lalan-be/internal/features/public/route"
	publicService "lalan-be/internal/features/public/service"

	_ "github.com/lib/pq"
)

/* Memulai aplikasi web dengan koneksi database, pengaturan dependency injection, rute, dan peluncuran server. Server berjalan di port 8080 tanpa kesalahan. */
func main() {
	db, err := config.DatabaseConfig()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	log.Println("Database connected")

	// Admin DI chain
	aRepo := admin.NewAdminRepository(db.DB)
	aService := admin.NewAdminService(aRepo)
	aHandler := admin.NewAdminHandler(aService)

	// Public DI chain
	pRepo := publicRepo.NewPublicRepository(db.DB)
	pService := publicService.NewPublicService(pRepo)
	pHandler := publicHandler.NewPublicHandler(pService)

	// Setup routes
	mux := http.NewServeMux()
	admin.SetupAdminRoutes(mux, aHandler)
	publicRoute.SetupPublicRoutes(mux, pHandler)

	// Start server
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
