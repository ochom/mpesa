package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/controllers/c2b"
	"github.com/ochom/mpesa/src/domain"
)

// HandleStkPush ...
func HandleStkPush(ctx fiber.Ctx) error {
	req, err := parseDataValidate[domain.MpesaExpressRequest](ctx)
	if err != nil {
		return err
	}

	req.PhoneNumber = helpers.ParseMobile(req.PhoneNumber)

	go c2b.InitiatePayment(&req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleResult ...
func HandleC2BResult(ctx fiber.Ctx) error {
	logs.Info("c2b result => %s", string(ctx.Body()))

	id := ctx.Query("refId")
	if id == "" {
		logs.Error("c2b result => refId is required")
		return ctx.JSON(fiber.Map{"message": "failed, refId is required"})
	}

	req, err := parseDataValidate[domain.MpesaExpressCallback](ctx)
	if err != nil {
		logs.Error("c2b result parse error: => %s", err)
		return err
	}

	go c2b.ResultPayment(id, &req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleRestValidation ...
func HandleRestValidation(ctx fiber.Ctx) error {
	logs.Info("c2b rest validation => %s", string(ctx.Body()))

	req, err := parseDataValidate[domain.ValidationRequest](ctx)
	if err != nil {
		logs.Error("c2b rest validation parse error: => %s", err)
		return err
	}

	code, desc := "0", "Accepted"
	if ok := c2b.ValidatePayment(&req); !ok {
		code, desc = "C2B00016", "Rejected"
	}

	return ctx.JSON(fiber.Map{
		"ResultCode": code,
		"ResultDesc": desc,
	})
}

// HandleRestConfirmation ...
func HandleRestConfirmation(ctx fiber.Ctx) error {
	logs.Info("c2b rest confirmation => %s", string(ctx.Body()))

	req, err := parseDataValidate[domain.ValidationRequest](ctx)
	if err != nil {
		logs.Error("c2b rest confirmation parse error: => %s", err)
		return err
	}

	go c2b.ConfirmPayment(&req)
	return ctx.JSON(fiber.Map{
		"ResultCode": "0",
		"ResultDesc": "Success",
	})
}

// HandleSoapValidation ...
func HandleSoapValidation(ctx fiber.Ctx) error {
	logs.Info("c2b soap validation => %s", string(ctx.Body()))

	txId := uuid.New()
	template := strings.Replace(config.SoapValidationTemplate, "{RESULT_CODE}", "0", 1)
	template = strings.Replace(template, "{RESULT_DESCRIPTION}", "Accepted", 1)
	template = strings.Replace(template, "{THIRD_PARTY_TRANSACTION_ID}", txId, 1)

	return ctx.SendString(template)
}

// HandleSoapConfirmation ...
func HandleSoapConfirmation(ctx fiber.Ctx) error {
	logs.Info("c2b soap confirmation => %s", string(ctx.Body()))

	var req domain.SoapPaymentConfirmationRequest
	if err := ctx.Bind().XML(&req); err != nil {
		logs.Error("c2b soap confirmation parse error: => %s", err)
		return err
	}

	validationRequest := domain.ValidationRequest{
		TransactionType:   req.Body.C2BPaymentConfirmationRequest.TransType,
		TransID:           req.Body.C2BPaymentConfirmationRequest.TransID,
		TransTime:         req.Body.C2BPaymentConfirmationRequest.TransTime,
		TransAmount:       req.Body.C2BPaymentConfirmationRequest.TransAmount,
		BusinessShortCode: req.Body.C2BPaymentConfirmationRequest.BusinessShortCode,
		BillRefNumber:     req.Body.C2BPaymentConfirmationRequest.BillRefNumber,
		OrgAccountBalance: req.Body.C2BPaymentConfirmationRequest.OrgAccountBalance,
		MSISDN:            req.Body.C2BPaymentConfirmationRequest.MSISDN,
	}

	go c2b.ConfirmPayment(&validationRequest)
	template := strings.Replace(config.SoapConfirmationTemplate, "{TRANSACTION_ID}", validationRequest.TransID, 1)
	return ctx.SendString(template)
}
