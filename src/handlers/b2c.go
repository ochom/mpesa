package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ochom/mpesa/src/controllers/b2c"
	"github.com/ochom/mpesa/src/domain"
)

// HandleInitiatePayment ...
func HandleInitiatePayment(ctx fiber.Ctx) error {
	req, err := parseData[domain.B2cRequest](ctx)
	if err != nil {
		return err
	}

	go b2c.InitiatePayment(req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleB2CResult ...
func HandleB2CResult(ctx fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	req, err := parseData[domain.B2cResult](ctx)
	if err != nil {
		return err
	}

	go b2c.ResultPayment(id, &req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleB2cTimeout ...
func HandleB2cTimeout(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	go b2c.TimeoutPayment(id)
	return nil
}
