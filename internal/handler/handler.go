package handler

import "github.com/maythitirat/pet-log-api/internal/service"

// Handlers holds all HTTP handlers
type Handlers struct {
	Health *HealthHandler
	Pet    *PetHandler
	User   *UserHandler
}

// NewHandlers creates all handlers with their dependencies
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(),
		Pet:    NewPetHandler(services.Pet),
		User:   NewUserHandler(services.User),
	}
}
