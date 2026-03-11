package middleware

import "github.com/gofiber/fiber/v3"

// CORS returns a middleware that handles Cross-Origin Resource Sharing
func CORS() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Set CORS headers
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
		c.Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.Next()
	}
}
