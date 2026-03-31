package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func KeycloakMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		if userID == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Locals("userID", userID)

		role := c.Get("X-User-Role")
		if role == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Locals("role", role)

		return c.Next()
	}
}
