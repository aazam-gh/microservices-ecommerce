package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func parseUserID(c *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(c.Get("X-User-ID"))
}

func parseBody[T any](c *fiber.Ctx) (*T, error) {
	var req T
	if err := c.BodyParser(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func parseAddressID(c *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(c.Params("id"))
}
