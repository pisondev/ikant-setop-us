package storage

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/shared"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/cold-storages", h.list)
	router.Post("/cold-storages", h.create)
}

func (h *Handler) list(c *fiber.Ctx) error {
	items, err := h.repo.List(c.Context())
	if err != nil {
		return err
	}
	return shared.Success(c, fiber.StatusOK, "Cold storages retrieved successfully", items)
}

type createBody struct {
	Name          string  `json:"name"`
	LocationLabel *string `json:"location_label"`
	Description   *string `json:"description"`
}

func (h *Handler) create(c *fiber.Ctx) error {
	var body createBody
	if err := c.BodyParser(&body); err != nil {
		return shared.Error(c, fiber.StatusBadRequest, "Invalid JSON body", nil)
	}

	input, fieldErrors := validateCreate(body)
	if len(fieldErrors) > 0 {
		return shared.Error(c, fiber.StatusBadRequest, "Validation error", fieldErrors)
	}

	item, err := h.repo.Create(c.Context(), input)
	if err != nil {
		return err
	}
	return shared.Success(c, fiber.StatusCreated, "Cold storage created successfully", item)
}

func validateCreate(body createBody) (CreateInput, []shared.FieldError) {
	name := strings.TrimSpace(body.Name)
	if name == "" {
		return CreateInput{}, []shared.FieldError{{Field: "name", Message: "Name is required"}}
	}

	return CreateInput{
		Name:          name,
		LocationLabel: cleanOptionalString(body.LocationLabel),
		Description:   cleanOptionalString(body.Description),
	}, nil
}

func cleanOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
