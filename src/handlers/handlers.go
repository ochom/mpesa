package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func parseData[T any](ctx *fiber.Ctx) (T, error) {
	var req T
	err := ctx.BodyParser(&req)
	return req, err
}

func parseDataValidate[T any](ctx *fiber.Ctx) (T, error) {
	var req T
	err := ctx.BodyParser(&req)
	if err != nil {
		return req, err
	}

	if err := validate.Struct(req); err != nil {
		return req, err
	}

	return req, nil
}
