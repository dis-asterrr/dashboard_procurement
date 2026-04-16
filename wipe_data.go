package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment from .env or backend/.env
	_ = godotenv.Load(".env", "backend/.env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		user := getEnvOrDefault("DB_USER", "postgres")
		password := os.Getenv("DB_PASSWORD")
		dbname := getEnvOrDefault("DB_NAME", "rygell_dashboard")
		sslmode := getEnvOrDefault("DB_SSLMODE", "disable")
		if password == "" {
			log.Fatal("DB_PASSWORD environment variable is required")
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Truncate tables
	err = db.Exec("TRUNCATE TABLE contract_dedicated_fixes, contract_dedicated_vars, contract_oncalls RESTART IDENTITY CASCADE").Error
	if err != nil {
		log.Fatalf("failed to truncate tables: %v", err)
	}

	fmt.Println("Successfully wiped active data from contract tables.")
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
