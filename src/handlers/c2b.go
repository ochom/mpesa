package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/mpesa/src/controllers/c2b"
	"github.com/ochom/mpesa/src/domain"
)

// HandleStkPush ...
func HandleStkPush(ctx fiber.Ctx) error {
	req, err := parseData[domain.MpesaExpressRequest](ctx)
	if err != nil {
		return err
	}

	req.PhoneNumber = helpers.ParseMobile(req.PhoneNumber)

	go c2b.InitiatePayment(&req)
	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleResult ...
func HandleC2BResult(ctx fiber.Ctx) error {
	id := ctx.Query("refId")
	if id == "" {
		return ctx.JSON(fiber.Map{"message": "failed, refId is required"})
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

	if ok := c2b.ValidatePayment(&req); ok {
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

	go c2b.ConfirmPayment(&req)
	return ctx.JSON(fiber.Map{
		"ResultCode": "0",
		"ResultDesc": "Success",
	})
}

// HandleSoapValidation ...
func HandleSoapValidation(ctx fiber.Ctx) error {
	// TODO parse the data into SOAP data
	req, err := parseData[domain.ValidationRequest](ctx)
	if err != nil {
		return err
	}

	template := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
    xmlns:c2b="http://cps.huawei.com/cpsinterface/c2bpayment">
    <soapenv:Header/>
    <soapenv:Body>
        <c2b:C2BPaymentValidationResult>
            <ResultCode>{RESULT_CODE}</ResultCode>
            <ResultDesc>{RESULT_DESCRIPTION}</ResultDesc>
            <ThirdPartyTransactionID>{THIRD_PARTY_TRANSACTION_ID}</ThirdPartyTransactionID>
        </c2b:C2BPaymentValidationResult>
    </soapenv:Body>
</soapenv:Envelope>`

	txId := uuid.New()
	if ok := c2b.ValidatePayment(&req); ok {
		template = strings.Replace(template, "{RESULT_CODE}", "0", 1)
		template = strings.Replace(template, "{RESULT_DESCRIPTION}", "Accepted", 1)
		template = strings.Replace(template, "{THIRD_PARTY_TRANSACTION_ID}", txId, 1)
	} else {
		template = strings.Replace(template, "{RESULT_CODE}", "C2B00016", 1)
		template = strings.Replace(template, "{RESULT_DESCRIPTION}", "Rejected", 1)
		template = strings.Replace(template, "{THIRD_PARTY_TRANSACTION_ID}", txId, 1)
	}

	return ctx.SendString(template)
}

// HandleSoapConfirmation ...
func HandleSoapConfirmation(ctx fiber.Ctx) error {
	var req domain.SoapPaymentConfirmationRequest
	if err := ctx.Bind().XML(&req); err != nil {
		logs.Error("Error decoding XML: %v", err)
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

	template := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
    xmlns:c2b="http://cps.huawei.com/cpsinterface/c2bpayment">
    <soapenv:Header/>
    <soapenv:Body>
        <c2b:C2BPaymentConfirmationResult>C2B Payment result received | Transaction ID: {TRANSACTION_ID}</c2b:C2BPaymentConfirmationResult>
    </soapenv:Body>
</soapenv:Envelope>`

	template = strings.Replace(template, "{TRANSACTION_ID}", validationRequest.TransID, 1)

	return ctx.SendString(template)
}
