package main

import (
	"log"
	"strings"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/database"
	"rygell-dashboard/internal/repositories"
	"rygell-dashboard/internal/router"
	"rygell-dashboard/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (try local over relative path)
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load("../../.env")
	}

	// Load configuration
	cfg := config.Load()
	if strings.TrimSpace(cfg.JWTSecret) == "" || cfg.JWTSecret == "change_me" || cfg.JWTSecret == "rygell_super_secret_change_me" {
		log.Fatalf("Invalid JWT_SECRET: set a strong non-default value")
	}
	if strings.TrimSpace(cfg.AdminPassword) == "" || cfg.AdminPassword == "admin123" || cfg.AdminPassword == "change_me" {
		log.Fatalf("Invalid ADMIN_PASSWORD: set a strong non-default value")
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

	// Ensure default admin account exists
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	if err := userService.EnsureDefaultAdmin(cfg.AdminName, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Fatalf("Failed to ensure default admin user: %v", err)
	}

	// Setup router
	r := router.Setup(db, cfg)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Rygell Dashboard API starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
