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
	// config & DB
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer cfg.DB.Close()

	// DI: repo → service → handler
	authRepo := repository.NewAuthRepository(cfg.DB)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	// register route
	route.AuthRoutes(authHandler)

	// port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
