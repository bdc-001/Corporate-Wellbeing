package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func NewConnection(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sqlx.DB) error {
	// Read schema file
	schemaPath := os.Getenv("SCHEMA_PATH")
	if schemaPath == "" {
		schemaPath = "../../database/schema.sql"
	}

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		// If schema file doesn't exist, skip migrations
		// In production, use a proper migration tool
		return nil
	}

	// Execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		// Ignore errors if tables already exist
		return nil
	}

	return nil
}

// DB is a wrapper around sqlx.DB for convenience
type DB struct {
	*sqlx.DB
}

func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.DB.Get(dest, query, args...)
}

func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.DB.Select(dest, query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

