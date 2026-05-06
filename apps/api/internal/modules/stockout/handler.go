package stockout

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/shared"
)

type Handler struct {
	store   Store
	service *Service
}

func NewHandler(store Store, service *Service) *Handler {
	return &Handler{store: store, service: service}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/stock-outs", h.list)
	router.Post("/stock-outs", h.create)
}

type createBody struct {
	FishTypeID    string  `json:"fish_type_id"`
	TotalWeightKG float64 `json:"total_weight_kg"`
	Destination   string  `json:"destination"`
	OutAt         string  `json:"out_at"`
	Notes         *string `json:"notes"`
}

func (h *Handler) create(c *fiber.Ctx) error {
	var body createBody
	if err := c.BodyParser(&body); err != nil {
		return shared.Error(c, fiber.StatusBadRequest, "Invalid JSON body", nil)
	}

	input, fieldErrors := validateCreate(body)
	if len(fieldErrors) > 0 {
		return validationError(c, fieldErrors)
	}

	item, err := h.service.Create(c.Context(), input)
	if err != nil {
		return mapError(c, err)
	}
	return shared.Success(c, fiber.StatusCreated, "Stock out created successfully", item)
}

func (h *Handler) list(c *fiber.Ctx) error {
	filter, fieldErrors := parseListFilter(c)
	if len(fieldErrors) > 0 {
		return validationError(c, fieldErrors)
	}

	items, err := h.store.List(c.Context(), filter)
	if err != nil {
		return err
	}
	return shared.Success(c, fiber.StatusOK, "Stock outs retrieved successfully", items)
}

func validateCreate(body createBody) (CreateInput, []shared.FieldError) {
	fieldErrors := []shared.FieldError{}
	if err := validateUUID("fish_type_id", body.FishTypeID); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if body.TotalWeightKG <= 0 {
		fieldErrors = append(fieldErrors, fieldError("total_weight_kg", "Total weight must be greater than 0"))
	}
	destination := strings.TrimSpace(body.Destination)
	if destination == "" {
		fieldErrors = append(fieldErrors, fieldError("destination", "Destination is required"))
	}
	outAt, timeErr := parseRequiredTime("out_at", body.OutAt)
	if timeErr != nil {
		fieldErrors = append(fieldErrors, *timeErr)
	}

	return CreateInput{
		FishTypeID:    body.FishTypeID,
		TotalWeightKG: body.TotalWeightKG,
		Destination:   destination,
		OutAt:         outAt,
		Notes:         cleanOptionalString(body.Notes),
	}, fieldErrors
}

func parseListFilter(c *fiber.Ctx) (ListFilter, []shared.FieldError) {
	filter := ListFilter{}
	fieldErrors := []shared.FieldError{}

	if value := strings.TrimSpace(c.Query("fish_type_id")); value != "" {
		if err := validateUUID("fish_type_id", value); err != nil {
			fieldErrors = append(fieldErrors, *err)
		} else {
			filter.FishTypeID = &value
		}
	}
	if value := strings.TrimSpace(c.Query("destination")); value != "" {
		filter.Destination = &value
	}
	if value := strings.TrimSpace(c.Query("date_from")); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			fieldErrors = append(fieldErrors, fieldError("date_from", "Date must use YYYY-MM-DD format"))
		} else {
			filter.DateFrom = &parsed
		}
	}
	if value := strings.TrimSpace(c.Query("date_to")); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			fieldErrors = append(fieldErrors, fieldError("date_to", "Date must use YYYY-MM-DD format"))
		} else {
			filter.DateTo = &parsed
		}
	}

	return filter, fieldErrors
}

func parseRequiredTime(field string, value string) (time.Time, *shared.FieldError) {
	value = strings.TrimSpace(value)
	if value == "" {
		err := fieldError(field, "Time is required")
		return time.Time{}, &err
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fieldErr := fieldError(field, "Time must use ISO 8601/RFC3339 format")
		return time.Time{}, &fieldErr
	}
	return parsed, nil
}

func validateUUID(field string, value string) *shared.FieldError {
	if strings.TrimSpace(value) == "" {
		err := fieldError(field, "UUID is required")
		return &err
	}
	if _, err := uuid.Parse(value); err != nil {
		fieldErr := fieldError(field, "Value must be a valid UUID")
		return &fieldErr
	}
	return nil
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

func fieldError(field string, message string) shared.FieldError {
	return shared.FieldError{Field: field, Message: message}
}

func validationError(c *fiber.Ctx, errs []shared.FieldError) error {
	return shared.Error(c, fiber.StatusBadRequest, "Validation error", errs)
}

func mapError(c *fiber.Ctx, err error) error {
	var insufficient InsufficientStockError
	if errors.As(err, &insufficient) {
		return shared.Error(c, fiber.StatusBadRequest, "Insufficient stock", []shared.FieldError{
			fieldError(
				"total_weight_kg",
				fmt.Sprintf("Requested %.2f kg, but only %.2f kg is available", insufficient.Requested, insufficient.Available),
			),
		})
	}
	return err
}
