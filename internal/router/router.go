package router

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/maythitirat/pet-log-api/internal/config"
	"github.com/maythitirat/pet-log-api/internal/handler"
	"github.com/maythitirat/pet-log-api/internal/middleware"
)

// NewRouter creates and configures the router with all routes
func NewRouter(h *handler.Handlers, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CORS)

	// Health check routes (no auth required)
	r.Get("/health", h.Health.Health)
	r.Get("/ready", h.Health.Ready)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Pet routes
		r.Route("/pets", func(r chi.Router) {
			r.Get("/", h.Pet.GetAll)     // GET /api/v1/pets
			r.Post("/", h.Pet.Create)    // POST /api/v1/pets
			r.Get("/{id}", h.Pet.GetByID)    // GET /api/v1/pets/{id}
			r.Put("/{id}", h.Pet.Update)     // PUT /api/v1/pets/{id}
			r.Delete("/{id}", h.Pet.Delete)  // DELETE /api/v1/pets/{id}
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", h.User.Create)               // POST /api/v1/users
			r.Get("/{id}", h.User.GetByID)           // GET /api/v1/users/{id}
			r.Put("/{id}", h.User.Update)            // PUT /api/v1/users/{id}
			r.Delete("/{id}", h.User.Delete)         // DELETE /api/v1/users/{id}
			r.Get("/{userId}/pets", h.Pet.GetByOwnerID) // GET /api/v1/users/{userId}/pets
		})
	})

	return r
}
