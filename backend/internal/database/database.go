package database

import (
	"fmt"
	"strings"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// parseLogLevel converts DB_LOG_LEVEL string to a GORM logger.LogLevel.
func parseLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

// Connect establishes a connection to the PostgreSQL database.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)
	}

	logLevel := parseLogLevel(cfg.DBLogLevel)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs auto-migration for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Mill{},
		&models.Vendor{},
		&models.Product{},
		&models.Zone{},
		&models.Mot{},
		&models.Uom{},
		&models.ContractDedicatedFix{},
		&models.ContractDedicatedVar{},
		&models.ContractOncall{},
		&models.AuditLog{},
		&models.User{},
	)
}
