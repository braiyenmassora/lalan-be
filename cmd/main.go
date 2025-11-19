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
	"lalan-be/internal/features/customer"
	"lalan-be/internal/features/hoster"
	"lalan-be/internal/features/public"
	"lalan-be/internal/middleware"
)

/*
main
menjalankan aplikasi server dengan inisialisasi dan shutdown graceful
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

	AdminRepo := admin.NewAdminRepository(db)
	AdminService := admin.NewAdminService(AdminRepo)
	AdminHandler := admin.NewAdminHandler(AdminService)

	HosterRepo := hoster.NewHosterRepository(db)
	HosterService := hoster.NewHosterService(HosterRepo)
	HosterHandler := hoster.NewHosterHandler(HosterService)

	CustomerRepo := customer.NewCustomerRepository(db)
	CustomerService := customer.NewCustomerService(CustomerRepo)
	CustomerHandler := customer.NewCustomerHandler(CustomerService)

	PublicRepo := public.NewPublicRepository(db)
	PublicService := public.NewPublicService(PublicRepo)
	PublicHandler := public.NewPublicHandler(PublicService)

	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)

	admin.SetupAdminRoutes(router, AdminHandler)
	hoster.SetupHosterRoutes(router, HosterHandler)
	customer.SetupCustomerRoutes(router, CustomerHandler)
	public.SetupPublicRoutes(router, PublicHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
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
