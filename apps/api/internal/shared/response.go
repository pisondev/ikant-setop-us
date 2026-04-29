package shared

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string, errs interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func ErrorHandler(log *logrus.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return Error(c, fiberErr.Code, fiberErr.Message, nil)
		}

		log.WithError(err).Error("unhandled request error")
		return Error(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
}
