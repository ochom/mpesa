package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func parseData[T any](ctx fiber.Ctx) (T, error) {
	var data T
	err := ctx.Bind().Body(&data)
	if err != nil {
		return data, err
	}

	if err := validate.Struct(data); err != nil {
		return data, err
	}

	return data, nil
}
