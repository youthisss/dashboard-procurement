package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/database"
	"rygell-dashboard/internal/repositories"
	"rygell-dashboard/internal/router"
	"rygell-dashboard/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Prefer real local env files; fallback to .env.example only for dev convenience.
	loadedEnv := false
	for _, envPath := range []string{".env", "../../.env"} {
		if err := godotenv.Load(envPath); err == nil {
			loadedEnv = true
			break
		}
	}
	if !loadedEnv {
		log.Println("WARN: No .env file found, falling back to .env.example (dev only)")
		_ = godotenv.Load("../../.env.example")
	}

	// Load configuration
	cfg := config.Load()

	// Validate JWT_SECRET: must not be empty, must not be a known default, must be >= 32 chars.
	jwtTrimmed := strings.TrimSpace(cfg.JWTSecret)
	knownDefaults := []string{"", "change_me", "rygell_super_secret_change_me", "abcdefghijk",
		"change_me_to_a_strong_random_secret_at_least_32_chars"}
	for _, def := range knownDefaults {
		if jwtTrimmed == def {
			log.Fatalf("Invalid JWT_SECRET: set a strong non-default value (at least 32 characters)")
		}
	}
	if len(jwtTrimmed) < 32 {
		log.Fatalf("JWT_SECRET is too short (%d chars): must be at least 32 characters", len(jwtTrimmed))
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Ensure default admin account exists (and optionally sync password)
	userRepo := repositories.NewUserRepository(db)
	totalUsers, err := userRepo.Count()
	if err != nil {
		log.Fatalf("Failed to count existing users: %v", err)
	}
	if totalUsers == 0 || cfg.AdminSyncPassword {
		if err := services.ValidatePassword(cfg.AdminPassword); err != nil {
			log.Fatalf("Invalid ADMIN_PASSWORD: %v", err)
		}
	}

	userService := services.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	if err := userService.EnsureDefaultAdminWithSync(cfg.AdminName, cfg.AdminUsername, cfg.AdminPassword, cfg.AdminSyncPassword); err != nil {
		log.Fatalf("Failed to ensure default admin user: %v", err)
	}

	// Setup router
	r := router.Setup(db, cfg)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Rygell Dashboard API starting on %s", addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
