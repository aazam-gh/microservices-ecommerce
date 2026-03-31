package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/service"
)

type AuthHandler struct {
	service service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{
		service: authService,
	}
}

func (h *AuthHandler) Validate(c *fiber.Ctx) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	validateTokenResp, err := h.service.ValidateToken(c.UserContext(), token)
	if err != nil {
		if err == domain.ErrKeycloakUnreachable {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set("X-User-ID", validateTokenResp.UserId.String())
	c.Set("X-User-Role", string(validateTokenResp.Role))
	c.Set("X-User-Email", validateTokenResp.Email)
	return nil
}
