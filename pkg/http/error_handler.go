package http

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := err.Error()

	var e *Exception
	if errors.As(err, &e) {
		status = e.Status
	}

	var fiberExp *fiber.Error
	if errors.As(err, &fiberExp) {
		status = fiberExp.Code
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		status = fiber.StatusBadRequest
	}

	return c.Status(status).JSON(Error(message))
}
