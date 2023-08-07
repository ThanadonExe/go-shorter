package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
)

type UrlHandler struct {
	urlService ports.UrlService
}

func NewUrlHandler(urlService ports.UrlService) *UrlHandler {
	return &UrlHandler{
		urlService: urlService,
	}
}

func (h *UrlHandler) Create(ctx *fiber.Ctx) error {
	var request *domain.UrlCreateRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to parse JSON"})
	}

	response, err := h.urlService.Create(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - failed to create url"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *UrlHandler) Get(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid url id format"})
	}

	u, err := h.urlService.GetById(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "url not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&u)
}

func (h *UrlHandler) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error - invalid url id format"})
	}

	err = h.urlService.Delete(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - failed to delete url"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Url deleted."})
}

func (h *UrlHandler) GetAll(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	urls, err := h.urlService.GetAll(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - urls not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&urls)
}

func (h *UrlHandler) Resolve(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	u, err := h.urlService.GetByCode(code)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Error - url not found"})
	}

	return ctx.Redirect(u.FullUrl, fiber.StatusFound)
}
