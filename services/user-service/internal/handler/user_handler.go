package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/service"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(userHandler service.IUserService) *UserHandler {
	return &UserHandler{
		service: userHandler,
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	user, err := h.service.GetProfile(c.Context(), id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	req, err := parseBody[domain.UpdateProfileRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.UpdateProfile(c.Context(), id, req)

	if err != nil {
		if err == domain.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
