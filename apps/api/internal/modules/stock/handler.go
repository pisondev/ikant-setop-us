package stock

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
	router.Get("/stocks/fifo", h.listFIFO)
	router.Get("/stocks/:id", h.detail)
	router.Get("/stocks", h.list)
	router.Post("/stocks", h.create)
	router.Patch("/stocks/:id/quality", h.updateQuality)
	router.Patch("/stocks/:id/location", h.updateLocation)
}

func (h *Handler) list(c *fiber.Ctx) error {
	filter, fieldErrors := parseListFilter(c)
	if len(fieldErrors) > 0 {
		return validationError(c, fieldErrors)
	}

	items, err := h.store.List(c.Context(), filter)
	if err != nil {
		return mapError(c, err)
	}
	return shared.Success(c, fiber.StatusOK, "Stocks retrieved successfully", items)
}

type createBody struct {
	FishTypeID      string  `json:"fish_type_id"`
	ColdStorageID   string  `json:"cold_storage_id"`
	Quality         string  `json:"quality"`
	InitialWeightKG float64 `json:"initial_weight_kg"`
	EnteredAt       string  `json:"entered_at"`
	Notes           *string `json:"notes"`
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
	return shared.Success(c, fiber.StatusCreated, "Stock batch created successfully", item)
}

func (h *Handler) detail(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := validateUUID("id", id); err != nil {
		return validationError(c, []shared.FieldError{*err})
	}

	item, err := h.store.GetDetail(c.Context(), id)
	if err != nil {
		return mapError(c, err)
	}
	return shared.Success(c, fiber.StatusOK, "Stock detail retrieved successfully", item)
}

type updateQualityBody struct {
	Quality string  `json:"quality"`
	Notes   *string `json:"notes"`
}

func (h *Handler) updateQuality(c *fiber.Ctx) error {
	id := c.Params("id")
	fieldErrors := []shared.FieldError{}
	if err := validateUUID("id", id); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}

	var body updateQualityBody
	if err := c.BodyParser(&body); err != nil {
		return shared.Error(c, fiber.StatusBadRequest, "Invalid JSON body", nil)
	}
	if !isValidQuality(body.Quality) {
		fieldErrors = append(fieldErrors, fieldError("quality", "Quality must be baik, sedang, or buruk"))
	}
	if len(fieldErrors) > 0 {
		return validationError(c, fieldErrors)
	}

	item, err := h.service.UpdateQuality(c.Context(), id, UpdateQualityInput{
		Quality: body.Quality,
		Notes:   cleanOptionalString(body.Notes),
	})
	if err != nil {
		return mapError(c, err)
	}
	return shared.Success(c, fiber.StatusOK, "Stock quality updated successfully", item)
}

type updateLocationBody struct {
	ColdStorageID string  `json:"cold_storage_id"`
	Notes         *string `json:"notes"`
}

func (h *Handler) updateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	fieldErrors := []shared.FieldError{}
	if err := validateUUID("id", id); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}

	var body updateLocationBody
	if err := c.BodyParser(&body); err != nil {
		return shared.Error(c, fiber.StatusBadRequest, "Invalid JSON body", nil)
	}
	if err := validateUUID("cold_storage_id", body.ColdStorageID); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if len(fieldErrors) > 0 {
		return validationError(c, fieldErrors)
	}

	item, err := h.service.UpdateLocation(c.Context(), id, UpdateLocationInput{
		ColdStorageID: body.ColdStorageID,
		Notes:         cleanOptionalString(body.Notes),
	})
	if err != nil {
		return mapError(c, err)
	}
	return shared.Success(c, fiber.StatusOK, "Stock location updated successfully", item)
}

func (h *Handler) listFIFO(c *fiber.Ctx) error {
	filter, fieldErrors := parseFIFOFilter(c)
	if len(fieldErrors) > 0 {
		return validationError(c, fieldErrors)
	}

	items, err := h.store.ListFIFO(c.Context(), filter)
	if err != nil {
		return mapError(c, err)
	}
	return shared.Success(c, fiber.StatusOK, "FIFO stocks retrieved successfully", items)
}

func parseListFilter(c *fiber.Ctx) (ListFilter, []shared.FieldError) {
	filter := ListFilter{Sort: strings.TrimSpace(c.Query("sort"))}
	fieldErrors := []shared.FieldError{}

	if value := strings.TrimSpace(c.Query("fish_type_id")); value != "" {
		if err := validateUUID("fish_type_id", value); err != nil {
			fieldErrors = append(fieldErrors, *err)
		} else {
			filter.FishTypeID = &value
		}
	}
	if value := strings.TrimSpace(c.Query("cold_storage_id")); value != "" {
		if err := validateUUID("cold_storage_id", value); err != nil {
			fieldErrors = append(fieldErrors, *err)
		} else {
			filter.ColdStorageID = &value
		}
	}
	if value := strings.TrimSpace(c.Query("quality")); value != "" {
		if !isValidQuality(value) {
			fieldErrors = append(fieldErrors, fieldError("quality", "Quality must be baik, sedang, or buruk"))
		} else {
			filter.Quality = &value
		}
	}
	if value := strings.TrimSpace(c.Query("status")); value != "" {
		if !isValidStatus(value) {
			fieldErrors = append(fieldErrors, fieldError("status", "Status must be available or depleted"))
		} else {
			filter.Status = &value
		}
	}
	if filter.Sort != "" && filter.Sort != "fifo" && filter.Sort != "latest" {
		fieldErrors = append(fieldErrors, fieldError("sort", "Sort must be fifo or latest"))
	}
	return filter, fieldErrors
}

func parseFIFOFilter(c *fiber.Ctx) (FIFOFilter, []shared.FieldError) {
	filter := FIFOFilter{Limit: 100}
	fieldErrors := []shared.FieldError{}

	if value := strings.TrimSpace(c.Query("fish_type_id")); value != "" {
		if err := validateUUID("fish_type_id", value); err != nil {
			fieldErrors = append(fieldErrors, *err)
		} else {
			filter.FishTypeID = &value
		}
	}
	if value := strings.TrimSpace(c.Query("limit")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			fieldErrors = append(fieldErrors, fieldError("limit", "Limit must be a positive number"))
		} else {
			filter.Limit = parsed
		}
	}
	if value := strings.TrimSpace(c.Query("offset")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			fieldErrors = append(fieldErrors, fieldError("offset", "Offset must be zero or greater"))
		} else {
			filter.Offset = parsed
		}
	}
	return filter, fieldErrors
}

func validateCreate(body createBody) (CreateInput, []shared.FieldError) {
	fieldErrors := []shared.FieldError{}
	if err := validateUUID("fish_type_id", body.FishTypeID); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validateUUID("cold_storage_id", body.ColdStorageID); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if !isValidQuality(body.Quality) {
		fieldErrors = append(fieldErrors, fieldError("quality", "Quality must be baik, sedang, or buruk"))
	}
	if body.InitialWeightKG <= 0 {
		fieldErrors = append(fieldErrors, fieldError("initial_weight_kg", "Initial weight must be greater than 0"))
	}
	enteredAt, timeErr := parseRequiredTime("entered_at", body.EnteredAt)
	if timeErr != nil {
		fieldErrors = append(fieldErrors, *timeErr)
	}

	return CreateInput{
		FishTypeID:      body.FishTypeID,
		ColdStorageID:   body.ColdStorageID,
		Quality:         body.Quality,
		InitialWeightKG: body.InitialWeightKG,
		EnteredAt:       enteredAt,
		Notes:           cleanOptionalString(body.Notes),
	}, fieldErrors
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

func isValidQuality(value string) bool {
	switch value {
	case QualityGood, QualityMedium, QualityBad:
		return true
	default:
		return false
	}
}

func isValidStatus(value string) bool {
	switch value {
	case StatusAvailable, StatusDepleted:
		return true
	default:
		return false
	}
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
	if errors.Is(err, ErrNotFound) {
		return shared.Error(c, fiber.StatusNotFound, "Resource not found", nil)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return shared.Error(c, fiber.StatusBadRequest, "Referenced resource does not exist", nil)
		case "22P02":
			return shared.Error(c, fiber.StatusBadRequest, "Invalid request value", nil)
		}
	}

	return fmt.Errorf("stock module: %w", err)
}
