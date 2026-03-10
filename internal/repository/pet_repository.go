package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/maythitirat/pet-log-api/internal/model"
)

// PetRepository defines the interface for pet data operations
type PetRepository interface {
	Create(ctx context.Context, pet *model.Pet) error
	GetByID(ctx context.Context, id int64) (*model.Pet, error)
	GetAll(ctx context.Context, limit, offset int) ([]*model.Pet, error)
	GetByOwnerID(ctx context.Context, ownerID int64) ([]*model.Pet, error)
	Update(ctx context.Context, pet *model.Pet) error
	Delete(ctx context.Context, id int64) error
}

// petRepository implements PetRepository
type petRepository struct {
	db *sqlx.DB
}

// NewPetRepository creates a new pet repository
func NewPetRepository(db *sqlx.DB) PetRepository {
	return &petRepository{db: db}
}

// Create inserts a new pet into the database
func (r *petRepository) Create(ctx context.Context, pet *model.Pet) error {
	query := `
		INSERT INTO pets (name, species, breed, birth_date, weight, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		pet.Name, pet.Species, pet.Breed, pet.BirthDate, pet.Weight, pet.OwnerID,
	).Scan(&pet.ID, &pet.CreatedAt, &pet.UpdatedAt)
}

// GetByID retrieves a pet by its ID
func (r *petRepository) GetByID(ctx context.Context, id int64) (*model.Pet, error) {
	var pet model.Pet
	query := `SELECT * FROM pets WHERE id = $1`
	
	err := r.db.GetContext(ctx, &pet, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Pet not found
		}
		return nil, err
	}
	return &pet, nil
}

// GetAll retrieves all pets with pagination
func (r *petRepository) GetAll(ctx context.Context, limit, offset int) ([]*model.Pet, error) {
	var pets []*model.Pet
	query := `SELECT * FROM pets ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	err := r.db.SelectContext(ctx, &pets, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return pets, nil
}

// GetByOwnerID retrieves all pets belonging to a specific owner
func (r *petRepository) GetByOwnerID(ctx context.Context, ownerID int64) ([]*model.Pet, error) {
	var pets []*model.Pet
	query := `SELECT * FROM pets WHERE owner_id = $1 ORDER BY created_at DESC`
	
	err := r.db.SelectContext(ctx, &pets, query, ownerID)
	if err != nil {
		return nil, err
	}
	return pets, nil
}

// Update modifies an existing pet
func (r *petRepository) Update(ctx context.Context, pet *model.Pet) error {
	query := `
		UPDATE pets 
		SET name = $1, species = $2, breed = $3, birth_date = $4, weight = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		pet.Name, pet.Species, pet.Breed, pet.BirthDate, pet.Weight, pet.ID,
	).Scan(&pet.UpdatedAt)
}

// Delete removes a pet from the database
func (r *petRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pets WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}
