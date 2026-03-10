package model

import "time"

// Pet represents a pet entity in the system
type Pet struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Species   string    `json:"species" db:"species"`     // e.g., "dog", "cat", "bird"
	Breed     string    `json:"breed" db:"breed"`         // e.g., "Golden Retriever", "Persian"
	BirthDate *time.Time `json:"birth_date" db:"birth_date"`
	Weight    *float64  `json:"weight" db:"weight"`       // in kilograms
	OwnerID   int64     `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePetRequest represents the request body for creating a new pet
type CreatePetRequest struct {
	Name      string     `json:"name" validate:"required,min=1,max=100"`
	Species   string     `json:"species" validate:"required,min=1,max=50"`
	Breed     string     `json:"breed" validate:"max=100"`
	BirthDate *time.Time `json:"birth_date"`
	Weight    *float64   `json:"weight" validate:"omitempty,gt=0"`
	OwnerID   int64      `json:"owner_id" validate:"required"`
}

// UpdatePetRequest represents the request body for updating a pet
type UpdatePetRequest struct {
	Name      *string    `json:"name" validate:"omitempty,min=1,max=100"`
	Species   *string    `json:"species" validate:"omitempty,min=1,max=50"`
	Breed     *string    `json:"breed" validate:"omitempty,max=100"`
	BirthDate *time.Time `json:"birth_date"`
	Weight    *float64   `json:"weight" validate:"omitempty,gt=0"`
}

// PetResponse represents the response body for a pet
type PetResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Species   string     `json:"species"`
	Breed     string     `json:"breed,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Weight    *float64   `json:"weight,omitempty"`
	OwnerID   int64      `json:"owner_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse converts Pet to PetResponse
func (p *Pet) ToResponse() *PetResponse {
	return &PetResponse{
		ID:        p.ID,
		Name:      p.Name,
		Species:   p.Species,
		Breed:     p.Breed,
		BirthDate: p.BirthDate,
		Weight:    p.Weight,
		OwnerID:   p.OwnerID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
