package service

import (
	"context"
	"errors"

	"github.com/maythitirat/pet-log-api/internal/model"
	"github.com/maythitirat/pet-log-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Common errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserService defines the interface for user business logic
type UserService interface {
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	GetByID(ctx context.Context, id int64) (*model.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, id int64) error
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id int64) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user.ToResponse(), nil
}

// GetByEmail retrieves a user by email (returns full user for auth purposes)
func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check if new email already exists (if changing email)
	if req.Email != nil && *req.Email != user.Email {
		existing, err := s.repo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}

	// Update only provided fields
	if req.Name != nil {
		user.Name = *req.Name
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// Delete removes a user
func (s *userService) Delete(ctx context.Context, id int64) error {
	// Check if user exists
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
}

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
