package handler

import (
	"github.com/maythitirat/pet-log-api/internal/config"
	"github.com/maythitirat/pet-log-api/internal/service"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Health    *HealthHandler
	Pet       *PetHandler
	User      *UserHandler
	LoginLine *LoginLineHandler
}

// NewHandlers creates all handlers with their dependencies
func NewHandlers(services *service.Services, cfg *config.Config) *Handlers {
	return &Handlers{
		Health:    NewHealthHandler(),
		Pet:       NewPetHandler(services.Pet),
		User:      NewUserHandler(services.User),
		LoginLine: NewLoginLineHandler(services.Login, &cfg.Line),
	}
}
