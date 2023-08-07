package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
)

type AuthHandler struct {
	authService ports.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var request *domain.AuthRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to parse JSON"})
	}

	if err := h.validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - bad input"})
	}

	token, err := h.authService.Login(request)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - invalid email or password"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&token)
}

func (h *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	var request *domain.RefreshRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to parse JSON"})
	}

	if err := h.validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - bad input"})
	}

	token, err := h.authService.Refresh(request)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - failed to refresh token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&token)
}
