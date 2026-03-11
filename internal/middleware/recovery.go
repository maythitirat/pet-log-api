package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

// Recoverer returns a middleware that recovers from panics
func Recoverer() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if rvr := recover(); rvr != nil {
				// Log the panic with stack trace
				log.Error().
					Interface("panic", rvr).
					Str("stack", string(debug.Stack())).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Msg("panic recovered")

				// Return 500 Internal Server Error
				_ = c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %v", rvr))
			}
		}()

		return c.Next()
	}
}
