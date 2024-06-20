package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/controllers/b2c"
	"github.com/ochom/mpesa/src/domain"
)

// HandleInitiatePayment ...
func HandleInitiatePayment(ctx fiber.Ctx) error {
	req, err := parseDataValidate[domain.B2cRequest](ctx)
	if err != nil {
		return err
	}

	go b2c.InitiatePayment(req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleB2CResult ...
func HandleB2CResult(ctx fiber.Ctx) error {
	logs.Info("b2c result => %s", string(ctx.Body()))

	id := ctx.Query("id")
	if id == "" {
		logs.Error("b2c result => id is required")
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	req, err := parseDataValidate[domain.B2cResult](ctx)
	if err != nil {
		logs.Error("b2c result parse error: => %s", err)
		return err
	}

	go b2c.ResultPayment(id, &req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleB2cTimeout ...
func HandleB2cTimeout(ctx fiber.Ctx) error {
	logs.Info("b2c timeout => %s", string(ctx.Body()))

	id := ctx.Params("id")
	if id == "" {
		logs.Error("b2c timeout => id is required")
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	go b2c.TimeoutPayment(id)
	return nil
}
