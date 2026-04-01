package config

import "os"

// Config holds all configuration for the application.
type Config struct {
	DatabaseURL     string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	ServerPort      string
	UploadDir       string
	FrontendOrigins string
	JWTSecret       string
	JWTExpiryHours  string
	AdminName       string
	AdminUsername   string
	AdminPassword   string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "rygell"),
		DBPassword:      getEnv("DB_PASSWORD", "rygell_secret"),
		DBName:          getEnv("DB_NAME", "rygell_dashboard"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		ServerPort:      getEnv("SERVER_PORT", getEnv("PORT", "8080")),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		FrontendOrigins: getEnv("FRONTEND_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		JWTSecret:       getEnv("JWT_SECRET", "rygell_super_secret_change_me"),
		JWTExpiryHours:  getEnv("JWT_EXPIRY_HOURS", "24"),
		AdminName:       getEnv("ADMIN_NAME", "Administrator"),
		AdminUsername:   getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:   getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
