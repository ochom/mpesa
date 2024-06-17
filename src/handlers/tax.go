package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/controllers/tax"
	"github.com/ochom/mpesa/src/domain"
)

func HandleTaxRemittance(ctx fiber.Ctx) error {
	req, err := parseData[domain.TaxRequest](ctx)
	if err != nil {
		logs.Error("failed to parse data: %v", err)
		return err
	}

	go tax.InitiatePayment(req)
	return nil
}
