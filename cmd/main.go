package main

import (
	"log"
	"net/http"
	"os"

	"lalan-be/internal/config"
	"lalan-be/internal/handler"
	"lalan-be/internal/repository"
	"lalan-be/internal/route"
	"lalan-be/internal/service"
)

func main() {
	// config & DB
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer cfg.DB.Close()

	// DI: repo → service → handler
	hosterRepo := repository.NewHosterRepository(cfg.DB)
	hosterService := service.NewHosterService(hosterRepo)
	hosterHandler := handler.NewHosterHandler(hosterService)

	// register route
	route.RegisterHosterRoutes(hosterHandler)

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
