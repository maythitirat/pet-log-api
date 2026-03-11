package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

// Logger returns a middleware that logs the start and end of each request
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request is complete
		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("remote_addr", c.IP()).
			Int("status", c.Response().StatusCode()).
			Int("bytes", len(c.Response().Body())).
			Dur("latency", time.Since(start)).
			Msg("request completed")

		return err
	}
}
