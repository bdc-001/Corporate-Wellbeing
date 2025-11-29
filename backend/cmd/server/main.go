package main

import (
	"log"
	"time"

	"github.com/convin/crae/internal/api"
	"github.com/convin/crae/internal/config"
	"github.com/convin/crae/internal/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate critical configuration
	if cfg.Environment == "production" {
		if cfg.JWTSecret == "change-me-in-production" {
			log.Fatal("JWT_SECRET must be changed in production")
		}
		if cfg.DatabaseURL == "postgres://localhost/convin_crae?sslmode=disable" {
			log.Fatal("DATABASE_URL must be configured for production")
		}
	}

	// Initialize database with connection pooling
	db, err := database.NewConnectionWithPool(
		cfg.DatabaseURL,
		cfg.DBMaxOpenConns,
		cfg.DBMaxIdleConns,
		time.Duration(cfg.DBConnMaxLifetime)*time.Second,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations (in production, use a proper migration tool)
	if err := database.RunMigrations(db); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// Initialize API server
	server, err := api.NewServer(db, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Start server with graceful shutdown
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
