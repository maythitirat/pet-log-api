package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/maythitirat/pet-log-api/internal/config"
	"github.com/maythitirat/pet-log-api/internal/handler"
	"github.com/maythitirat/pet-log-api/internal/middleware"
)

// NewRouter creates and configures the router with all routes
func NewRouter(h *handler.Handlers, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	// Global middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recoverer())
	app.Use(middleware.CORS())

	// Health check routes (no auth required)
	app.Get("/health", h.Health.Health)
	app.Get("/ready", h.Health.Ready)

	// Authentication routes
	app.Get("/loginLine", h.LoginLine.LoginLine) // Get /loginLine
	app.Get("/callback", h.LoginLine.Callback)   // Get /callback

	// API v1 routes
	api := app.Group("/api/v1")

	// Pet routes
	pets := api.Group("/pets")
	pets.Get("/", h.Pet.GetAll)       // GET /api/v1/pets
	pets.Post("/", h.Pet.Create)      // POST /api/v1/pets
	pets.Get("/:id", h.Pet.GetByID)   // GET /api/v1/pets/:id
	pets.Put("/:id", h.Pet.Update)    // PUT /api/v1/pets/:id
	pets.Delete("/:id", h.Pet.Delete) // DELETE /api/v1/pets/:id

	// User routes
	users := api.Group("/users")
	users.Post("/", h.User.Create)                 // POST /api/v1/users
	users.Get("/:id", h.User.GetByID)              // GET /api/v1/users/:id
	users.Put("/:id", h.User.Update)               // PUT /api/v1/users/:id
	users.Delete("/:id", h.User.Delete)            // DELETE /api/v1/users/:id
	users.Get("/:userId/pets", h.Pet.GetByOwnerID) // GET /api/v1/users/:userId/pets

	return app
}
