package main

import (
	"fmt"
	"os"

	"github.com/maythitirat/pet-log-api/internal/config"
	"github.com/maythitirat/pet-log-api/internal/handler"
	"github.com/maythitirat/pet-log-api/internal/repository"
	"github.com/maythitirat/pet-log-api/internal/router"
	"github.com/maythitirat/pet-log-api/internal/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Setup logger
	setupLogger(cfg.App.Environment)

	log.Info().
		Str("app", cfg.App.Name).
		Str("env", cfg.App.Environment).
		Msg("Starting application")

	// Initialize database
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	log.Info().Msg("Database connected successfully")

	// Initialize layers (Dependency Injection)
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos)
	handlers := handler.NewHandlers(services, cfg)

	// Setup router
	r := router.NewRouter(handlers, cfg)

	// Start server
	log.Info().Int("port", cfg.App.Port).Msg("Server is starting")
	if err := r.Listen(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

// setupLogger configures zerolog based on environment
func setupLogger(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env == "development" {
		// Pretty logging for development
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// JSON logging for production
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
