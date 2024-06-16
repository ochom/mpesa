package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func parseData[T any](ctx fiber.Ctx) (T, error) {
	var data T
	err := ctx.Bind().Body(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
