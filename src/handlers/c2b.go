package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/arrays"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/controllers/c2b"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
	"gorm.io/gorm"
)

// HandleGetC2BPayments ...
func HandleGetC2BPayments(ctx *fiber.Ctx) error {
	page, limit := ctx.QueryInt("page", 1), ctx.QueryInt("limit", 10)
	query := map[string]string{}
	if ctx.Query("account_id") != "" {
		query["account_id"] = ctx.Query("account_id")
	}

	if ctx.Query("phone_number") != "" {
		query["phone_number"] = helpers.ParseMobile(ctx.Query("phone_number"))
	}

	payments := sql.FindWithLimit[models.CustomerPayment](page, limit, func(d *gorm.DB) *gorm.DB {
		return d.Where(query).Order("created_at desc")
	})

	data := arrays.Map(payments, func(p *models.CustomerPayment) map[string]any {
		return map[string]any{
			"id":               p.ID,
			"created_at":       p.CreatedAt.Format(time.RFC3339),
			"account_id":       1,
			"transaction_type": p.TransactionType,
			"transaction_id":   p.TransactionID,
			"transaction_time": p.TransactionTime,
			"phone_number":     p.PhoneNumber,
			"amount":           p.Amount,
			"bill_ref_number":  p.BillRefNumber,
			"invoice_number":   p.InvoiceNumber,
		}
	})

	return ctx.JSON(data)
}

// HandleStkPush ...
func HandleStkPush(ctx *fiber.Ctx) error {
	req, err := parseDataValidate[domain.MpesaExpressRequest](ctx)
	if err != nil {
		return err
	}

	req.PhoneNumber = helpers.ParseMobile(req.PhoneNumber)
	if req.InvoiceNumber == "" {
		req.InvoiceNumber = req.PhoneNumber
	}

	go c2b.InitiatePayment(&req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleResult ...
func HandleC2BResult(ctx *fiber.Ctx) error {
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
func HandleRestValidation(ctx *fiber.Ctx) error {
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
func HandleRestConfirmation(ctx *fiber.Ctx) error {
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
func HandleSoapValidation(ctx *fiber.Ctx) error {
	logs.Info("c2b soap validation => %s", string(ctx.Body()))

	txId := uuid.New()
	template := strings.Replace(config.SoapValidationTemplate, "{RESULT_CODE}", "0", 1)
	template = strings.Replace(template, "{RESULT_DESCRIPTION}", "Accepted", 1)
	template = strings.Replace(template, "{THIRD_PARTY_TRANSACTION_ID}", txId, 1)

	return ctx.SendString(template)
}

// HandleSoapConfirmation ...
func HandleSoapConfirmation(ctx *fiber.Ctx) error {
	logs.Info("c2b soap confirmation => %s", string(ctx.Body()))

	var req domain.SoapPaymentConfirmationRequest
	if err := ctx.BodyParser(&req); err != nil {
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
