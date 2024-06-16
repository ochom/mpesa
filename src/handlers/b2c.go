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

// HandlePaymentCallback ...
func HandlePaymentCallback(ctx fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	req, err := parseData[domain.B2cResult](ctx)
	if err != nil {
		return err
	}

	go b2c.ResultBusinessPayment(id, &req)
	return ctx.JSON(fiber.Map{"message": "success"})
}
