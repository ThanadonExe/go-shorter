package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

func (h *UserHandler) Create(ctx *fiber.Ctx) error {
	var request *domain.UserCreateRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to parse JSON"})
	}

	if err := h.validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - bad input"})
	}

	response, err := h.userService.Create(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to create user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *UserHandler) Get(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - invalid user id"})
	}

	user, err := h.userService.GetById(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - user not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&user)
}

func (h *UserHandler) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - invalid user id"})
	}

	var request *domain.UserUpdateRequest
	err = ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to parse JSON"})
	}

	if err := h.validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - bad input"})
	}

	user, err := h.userService.Update(id, request)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - failed to update user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&user)
}

func (h *UserHandler) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - invalid user id format"})
	}

	err = h.userService.Delete(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - failed to delete user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted."})
}

func (h *UserHandler) GetAll(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	urls, err := h.userService.GetAll(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - users not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&urls)
}
