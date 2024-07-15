package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/controllers/tax"
	"github.com/ochom/mpesa/src/domain"
)

func HandleTaxRemittance(ctx *fiber.Ctx) error {
	req, err := parseDataValidate[domain.TaxRequest](ctx)
	if err != nil {
		logs.Error("failed to parse data: %v", err)
		return err
	}

	go tax.InitiatePayment(&req)
	return nil
}

// HandleTaxTimeout ...
func HandleTaxTimeout(ctx *fiber.Ctx) error {
	logs.Info("tax timeout => %s", string(ctx.Body()))

	id := ctx.Params("id")
	if id == "" {
		logs.Error("tax timeout => id is required")
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	go tax.TimeoutPayment(id)
	return nil
}

// HandleTaxResult ...
func HandleTaxResult(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		logs.Error("tax result => id is required")
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	req, err := parseDataValidate[domain.TaxResult](ctx)
	if err != nil {
		logs.Error("tax result parse error: => %s", err)
		return err
	}

	go tax.ResultPayment(id, &req)
	return nil
}
