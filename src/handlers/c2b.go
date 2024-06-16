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

	go c2b.MpesaExpressInitiate(&req)
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

	go c2b.MpesaExpressCallback(id, &req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleC2bValidation ...
func HandleC2bValidation(ctx fiber.Ctx) error {
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	if ok := c2b.MpesaExpressValidate(&req); ok {
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

// HandleC2bConfirmation ...
func HandleC2bConfirmation(ctx fiber.Ctx) error {
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	go c2b.MpesaExpressConfirm(&req)
	return ctx.JSON(fiber.Map{
		"ResultCode": "0",
		"ResultDesc": "Success",
	})
}
