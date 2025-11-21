package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lalan-be/internal/config"
	"lalan-be/internal/features/admin"
	"lalan-be/internal/features/customer"
	"lalan-be/internal/features/hoster"
	"lalan-be/internal/features/public"
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
main
entry point aplikasi yang menginisialisasi semua komponen dan menjalankan server dengan graceful shutdown
*/
func main() {

	/*
		LoadEnv & InitRedis
		memuat environment dan menginisialisasi koneksi Redis
	*/
	config.LoadEnv()
	config.InitRedis()

	/*
		DatabaseConfig
		membuat koneksi ke PostgreSQL
	*/
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer cfg.DB.Close()

	/*
		Handler
		membuat semua handler dengan service dan repository
	*/
	AdminHandler := admin.NewAdminHandler(admin.NewAdminService(admin.NewAdminRepository(cfg.DB)))
	HosterHandler := hoster.NewHosterHandler(hoster.NewHosterService(hoster.NewHosterRepository(cfg.DB)))
	CustomerHandler := customer.NewCustomerHandler(customer.NewCustomerService(customer.NewCustomerRepository(cfg.DB)))
	PublicHandler := public.NewPublicHandler(public.NewPublicService(public.NewPublicRepository(cfg.DB)))

	/*
		Router
		membuat router dan menambahkan middleware CORS
	*/
	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)

	/*
		SetupRoutes
		mendaftarkan semua route dari setiap fitur
	*/
	admin.SetupAdminRoutes(r, AdminHandler)
	hoster.SetupHosterRoutes(r, HosterHandler)
	customer.SetupCustomerRoutes(r, CustomerHandler)
	public.SetupPublicRoutes(r, PublicHandler)

	/*
		http.Server
		menjalankan server di port 8080
	*/
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server running at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crashed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	/*
		Graceful Shutdown
		menghentikan server dengan timeout 30 detik
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
