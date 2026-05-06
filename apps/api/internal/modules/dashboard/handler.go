package dashboard

import (
	"strconv"
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
	router.Get("/dashboard/summary", h.summary)
	router.Get("/dashboard/recent-movements", h.recentMovements)
}

func (h *Handler) summary(c *fiber.Ctx) error {
	item, err := h.repo.Summary(c.Context())
	if err != nil {
		return err
	}
	return shared.Success(c, fiber.StatusOK, "Dashboard summary retrieved successfully", item)
}

func (h *Handler) recentMovements(c *fiber.Ctx) error {
	limit, fieldErr := parseLimit(c.Query("limit"))
	if fieldErr != nil {
		return shared.Error(c, fiber.StatusBadRequest, "Validation error", []shared.FieldError{*fieldErr})
	}

	items, err := h.repo.RecentMovements(c.Context(), limit)
	if err != nil {
		return err
	}
	return shared.Success(c, fiber.StatusOK, "Recent movements retrieved successfully", items)
}

func parseLimit(value string) (int, *shared.FieldError) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 10, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		fieldErr := shared.FieldError{Field: "limit", Message: "Limit must be a positive number"}
		return 0, &fieldErr
	}
	return parsed, nil
}
