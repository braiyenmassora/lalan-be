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
	adminidentity "lalan-be/internal/features/admin/identity"
	auth "lalan-be/internal/features/auth"
	booking "lalan-be/internal/features/customer/booking"
	custidentity "lalan-be/internal/features/customer/identity"
	hosterbooking "lalan-be/internal/features/hoster/booking"
	hosteritem "lalan-be/internal/features/hoster/item"
	public "lalan-be/internal/features/public"
	"lalan-be/internal/middleware"
	"lalan-be/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Lalan Backend API...")

	// 1. Load environment variables
	config.LoadEnv()
	port := config.GetEnv("APP_PORT", "8080")
	log.Printf("Running in %s mode on port %s", config.GetEnv("APP_ENV", "dev"), port)

	// 2. Inisialisasi database
	dbCfg, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbCfg.DB.Close()

	// 3. Inisialisasi Redis (opsional)
	if err := config.InitRedis(); err != nil {
		log.Printf("Redis not available: %v (continuing without cache)", err)
	} else {
		defer config.CloseRedis()
	}

	// 4. Inisialisasi storage
	cfg := config.LoadStorageConfig()
	storage := utils.NewSupabaseStorageFromEnv() // Atau NewSupabaseStorage(cfg) jika perlu custom

	// 5. Inisialisasi handler dengan dependency injection
	// Public & Auth
	pubHandler := public.NewPublicHandler(public.NewPublicService(public.NewPublicRepository(dbCfg.DB)))
	authHandler := auth.NewAuthHandler(auth.NewAuthService(auth.NewAuthRepository(dbCfg.DB)))

	// Customer
	bookingHandler := booking.NewBookingHandler(booking.NewBookingService(booking.NewBookingRepository(dbCfg.DB)))
	customerIdentityHandler := custidentity.NewIdentityHandler(
		custidentity.NewIdentityService(custidentity.NewIdentityRepository(dbCfg.DB), storage),
	)

	// Hoster
	hosterHandler := hosterbooking.NewBookingHandler(hosterbooking.NewBookingService(hosterbooking.NewHosterBookingRepository(dbCfg.DB)))
	hosterItemHandler := hosteritem.NewHosterItemHandler(hosteritem.NewItemService(hosteritem.NewHosterItemRepository(dbCfg.DB), storage, cfg))

	// Admin
	adminIdentityHandler := adminidentity.NewAdminIdentityHandler(
		adminidentity.NewAdminIdentityService(adminidentity.NewAdminIdentityRepository(dbCfg.DB)),
	)

	// 6. Setup router & routes
	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Public & Auth
	public.SetupPublicRoutes(router, pubHandler)
	auth.SetupAuthRoutes(router, authHandler)

	// Customer
	booking.SetupBookingRoutes(router, bookingHandler)
	custidentity.SetupIdentityRoutes(router, customerIdentityHandler)

	// Hoster
	hosterbooking.SetupBookingRoutes(router, hosterHandler)
	hosteritem.SetupItemRoutes(router, hosterItemHandler)
	log.Println("Hoster item routes registered (GET /item, POST /item, DELETE /item/{id})")

	// Admin
	adminidentity.SetupAdminIdentityRoutes(router, adminIdentityHandler)

	// 7. Konfigurasi HTTP server dengan timeout aman
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. Jalankan server di background
	go func() {
		log.Printf("Server listening at http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crashed: %v", err)
		}
	}()

	// 9. Tunggu sinyal shutdown (Ctrl+C / SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 10. Graceful shutdown dengan timeout 10 detik
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

/*
healthCheck adalah endpoint monitoring sederhana.

Output:
- 200 OK + JSON {"status":"healthy",...}
*/
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"healthy","service":"lalan-backend-api"}`))
}
