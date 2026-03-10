package service

import (
	"context"
	"errors"

	"github.com/maythitirat/pet-log-api/internal/model"
	"github.com/maythitirat/pet-log-api/internal/repository"
)

// Common errors
var (
	ErrPetNotFound = errors.New("pet not found")
)

// PetService defines the interface for pet business logic
type PetService interface {
	Create(ctx context.Context, req *model.CreatePetRequest) (*model.PetResponse, error)
	GetByID(ctx context.Context, id int64) (*model.PetResponse, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*model.PetResponse, error)
	GetByOwnerID(ctx context.Context, ownerID int64) ([]*model.PetResponse, error)
	Update(ctx context.Context, id int64, req *model.UpdatePetRequest) (*model.PetResponse, error)
	Delete(ctx context.Context, id int64) error
}

// petService implements PetService
type petService struct {
	repo repository.PetRepository
}

// NewPetService creates a new pet service
func NewPetService(repo repository.PetRepository) PetService {
	return &petService{repo: repo}
}

// Create creates a new pet
func (s *petService) Create(ctx context.Context, req *model.CreatePetRequest) (*model.PetResponse, error) {
	pet := &model.Pet{
		Name:      req.Name,
		Species:   req.Species,
		Breed:     req.Breed,
		BirthDate: req.BirthDate,
		Weight:    req.Weight,
		OwnerID:   req.OwnerID,
	}

	if err := s.repo.Create(ctx, pet); err != nil {
		return nil, err
	}

	return pet.ToResponse(), nil
}

// GetByID retrieves a pet by ID
func (s *petService) GetByID(ctx context.Context, id int64) (*model.PetResponse, error) {
	pet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pet == nil {
		return nil, ErrPetNotFound
	}

	return pet.ToResponse(), nil
}

// GetAll retrieves all pets with pagination
func (s *petService) GetAll(ctx context.Context, page, pageSize int) ([]*model.PetResponse, error) {
	// Default values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	pets, err := s.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.PetResponse, len(pets))
	for i, pet := range pets {
		responses[i] = pet.ToResponse()
	}

	return responses, nil
}

// GetByOwnerID retrieves all pets belonging to an owner
func (s *petService) GetByOwnerID(ctx context.Context, ownerID int64) ([]*model.PetResponse, error) {
	pets, err := s.repo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.PetResponse, len(pets))
	for i, pet := range pets {
		responses[i] = pet.ToResponse()
	}

	return responses, nil
}

// Update updates an existing pet
func (s *petService) Update(ctx context.Context, id int64, req *model.UpdatePetRequest) (*model.PetResponse, error) {
	// Get existing pet
	pet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pet == nil {
		return nil, ErrPetNotFound
	}

	// Update only provided fields
	if req.Name != nil {
		pet.Name = *req.Name
	}
	if req.Species != nil {
		pet.Species = *req.Species
	}
	if req.Breed != nil {
		pet.Breed = *req.Breed
	}
	if req.BirthDate != nil {
		pet.BirthDate = req.BirthDate
	}
	if req.Weight != nil {
		pet.Weight = req.Weight
	}

	if err := s.repo.Update(ctx, pet); err != nil {
		return nil, err
	}

	return pet.ToResponse(), nil
}

// Delete removes a pet
func (s *petService) Delete(ctx context.Context, id int64) error {
	// Check if pet exists
	pet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pet == nil {
		return ErrPetNotFound
	}

	return s.repo.Delete(ctx, id)
}
