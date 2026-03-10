package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	handlers := handler.NewHandlers(services)

	// Setup router
	r := router.NewRouter(handlers, cfg)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Int("port", cfg.App.Port).Msg("Server is starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
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
