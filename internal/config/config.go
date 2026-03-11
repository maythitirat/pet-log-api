package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Line     LineConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Environment string
	Port        int
	Debug       bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LineConfig struct {
	ChannelID     string
	ChannelSecret string
	RedirectURI   string
	State         string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (for local development)
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("❌ .env not found")
	} else {
		fmt.Println("✅ .env loaded")
	}

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "pet-log-api"),
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnvAsInt("APP_PORT", 8080),
			Debug:       getEnvAsBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "pet_log"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Line: LineConfig{
			ChannelID:     getEnv("LINE_CHANNEL_ID", ""),
			ChannelSecret: getEnv("LINE_CHANNEL_SECRET", ""),
			RedirectURI:   getEnv("LINE_REDIRECT_URI", ""),
			State:         getEnv("LINE_STATE", "12345"),
		},
	}

	return cfg, nil
}

// DSN returns the database connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// Helper functions to read environment variables

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
