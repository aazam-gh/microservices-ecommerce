package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/service"
)

type AddressHandler struct {
	service service.IUserService
}

func NewAddressHandler(addressHandler service.IUserService) *AddressHandler {
	return &AddressHandler{
		service: addressHandler,
	}
}

func (h *AddressHandler) List(c *fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	addresses, err := h.service.GetAddresses(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(addresses)
}

func (h *AddressHandler) Create(c *fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	req, err := parseBody[domain.CreateAddressRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	address, err := h.service.CreateAddress(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(address)
}

func (h *AddressHandler) Update(c *fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	req, err := parseBody[domain.UpdateAddressRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	addressID, err := parseAddressID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	address, err := h.service.UpdateAddress(c.Context(), id, addressID, req)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case domain.ErrAddressNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case domain.ErrAddressNotOwned:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(address)
}

func (h *AddressHandler) Delete(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	addressID, err := parseAddressID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err = h.service.DeleteAddress(c.Context(), userID, addressID)
	if err != nil {
		switch err {
		case domain.ErrAddressNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case domain.ErrAddressNotOwned:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AddressHandler) SetDefault(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	addressID, err := parseAddressID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.service.SetDefaultAddress(c.Context(), userID, addressID)
	if err != nil {
		switch err {
		case domain.ErrAddressNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case domain.ErrAddressNotOwned:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
