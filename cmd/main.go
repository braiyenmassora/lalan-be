package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"lalan-be/internal/config"
	"lalan-be/internal/features/admin"
	"lalan-be/internal/features/hoster"
	"lalan-be/internal/features/public"
	"lalan-be/internal/middleware"
)

/*
Main menjalankan aplikasi server.
Menginisialisasi dan menjalankan dengan shutdown graceful.
*/
func main() {
	config.LoadEnv()
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	db := cfg.DB
	defer db.Close()
	log.Printf(
		"Database connected → host=%s port=%d db=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	// admin setup
	aRepo := admin.NewAdminRepository(db)
	aService := admin.NewAdminService(aRepo)
	aHandler := admin.NewAdminHandler(aService)
	// public setup
	pRepo := public.NewPublicRepository(db)
	pService := public.NewPublicService(pRepo)
	pHandler := public.NewPublicHandler(pService)
	// hoster setup
	hRepo := hoster.NewHosterRepository(db)
	hService := hoster.NewHosterService(hRepo)
	hHandler := hoster.NewHosterHandler(hService)

	router := mux.NewRouter()
	// Setup CORS Middleware
	router.Use(middleware.CORSMiddleware)

	admin.SetupAdminRoutes(router, aHandler)
	hoster.SetupHosterRoutes(router, hHandler)
	public.SetupPublicRoutes(router, pHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Println("Server running at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	<-c
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
