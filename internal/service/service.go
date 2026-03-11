package service

import "github.com/maythitirat/pet-log-api/internal/repository"

// Services holds all service instances
type Services struct {
	Pet   PetService
	User  UserService
	Login LoginLineService
}

// NewServices creates all services with their dependencies
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Pet:   NewPetService(repos.Pet),
		User:  NewUserService(repos.User),
		Login: NewLoginLineService(),
	}
}
