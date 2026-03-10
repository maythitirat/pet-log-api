package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/maythitirat/pet-log-api/internal/config"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Repositories holds all repository instances
type Repositories struct {
	Pet  PetRepository
	User UserRepository
}

// NewRepositories creates all repositories
func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Pet:  NewPetRepository(db),
		User: NewUserRepository(db),
	}
}
