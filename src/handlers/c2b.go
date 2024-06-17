package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ochom/mpesa/src/controllers/c2b"
	"github.com/ochom/mpesa/src/domain"
)

// HandleStkPush ...
func HandleStkPush(ctx fiber.Ctx) error {
	req, err := parseData[domain.MpesaExpressRequest](ctx)
	if err != nil {
		return err
	}

	go c2b.InitiatePayment(&req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleStkCallback ...
func HandleStkCallback(ctx fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.JSON(fiber.Map{"message": "failed"})
	}

	req, err := parseData[domain.MpesaExpressCallback](ctx)
	if err != nil {
		return err
	}

	go c2b.ResultPayment(id, &req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleRestValidation ...
func HandleRestValidation(ctx fiber.Ctx) error {
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	if ok := c2b.ValidateStkRest(&req); ok {
		return ctx.JSON(fiber.Map{
			"ResultCode": "0",
			"ResultDesc": "Accepted",
		})
	}

	return ctx.JSON(fiber.Map{
		"ResultCode": "C2B00016",
		"ResultDesc": "Rejected",
	})
}

// HandleRestConfirmation ...
func HandleRestConfirmation(ctx fiber.Ctx) error {
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	go c2b.ConfirmStkRest(&req)
	return ctx.JSON(fiber.Map{
		"ResultCode": "0",
		"ResultDesc": "Success",
	})
}

// HandleSoapValidation ...
func HandleSoapValidation(ctx fiber.Ctx) error {
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	if ok := c2b.ValidateStkSoap(&req); ok {
		return ctx.JSON(fiber.Map{
			"ResultCode": "0",
			"ResultDesc": "Accepted",
		})
	}

	return ctx.JSON(fiber.Map{
		"ResultCode": "C2B00016",
		"ResultDesc": "Rejected",
	})
}

// HandleSoapConfirmation ...
func HandleSoapConfirmation(ctx fiber.Ctx) error {
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	go c2b.ConfirmStkSoap(&req)
	return ctx.JSON(fiber.Map{
		"ResultCode": "0",
		"ResultDesc": "Success",
	})
}
