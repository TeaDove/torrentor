package api_presentation

import (
	"torrentor/backend/utils/validators_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func parseJSON[T any](c *fiber.Ctx) (T, error) {
	var v = new(T)

	err := c.BodyParser(v)
	if err != nil {
		return *v, sentUnprocessable(errors.Wrap(err, "failed to body parse"))
	}

	err = validators_utils.Validate.Struct(v)
	if err != nil {
		return *v, sentUnprocessable(errors.Wrap(err, "failed to validate"))
	}

	return *v, nil
}

func sentUnprocessable(err error) error {
	return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
}
