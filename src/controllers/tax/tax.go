package tax

import (
	"fmt"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/controllers/auth"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
	"github.com/ochom/mpesa/src/utils"
)

// InitiatePayment initiates an mpesa tax payment
func InitiatePayment(req *domain.TaxRequest) {
	account, err := sql.FindOneById[models.Account](req.AccountId)
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return
	}

	payment := models.NewTaxPayment(req.RequestId, req.ShortCode, req.PaymentRequestNumber, req.Amount, req.CallbackUrl)
	if err := sql.Create(payment); err != nil {
		logs.Error("Error creating payment: %v", err)
		return
	}

	headers := map[string]string{
		"Authorization": "Bearer " + auth.Authenticate(account),
		"Content-Type":  "application/json",
	}

	resultUrl := fmt.Sprintf("%s/v1/tax/result?id=%s", config.BaseUrl, payment.Uuid)
	timeoutUrl := fmt.Sprintf("%s/v1/tax/timeout?id=%s", config.BaseUrl, payment.Uuid)
	payload := map[string]string{
		"Initiator":              "TaxPayer",
		"SecurityCredential":     auth.GetSecurityCredentials(account),
		"Command ID":             "PayTaxToKRA",
		"SenderIdentifierType":   "4",
		"RecieverIdentifierType": "4",
		"Amount":                 req.Amount,
		"PartyA":                 req.ShortCode,
		"PartyB":                 "572572",
		"AccountReference":       req.PaymentRequestNumber,
		"Remarks":                "OK",
		"QueueTimeOutURL":        timeoutUrl,
		"ResultURL":              resultUrl,
	}

	url := fmt.Sprintf("%s/mpesa/b2b/v1/remittax", config.MpesaApiUrl)
	res, err := gttp.Post(url, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return
	}

	if res.Status > 204 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return
	}

	data := helpers.FromBytes[map[string]string](res.Body)
	payment.ResponseCode = data["ResponseCode"]
	payment.ResponseDescription = data["ResponseDescription"]
	payment.ConversationID = data["ConversationID"]
	payment.OriginatorConversationID = data["OriginatorConversationID"]
	if err := sql.Update(payment); err != nil {
		logs.Error("Error updating payment: %v", err)
		return
	}

	if data["ResponseCode"] != "0" {
		logs.Error(
			"failed to initiate payment =>ResponseCode: %v, ResponseDescription: %v, errorCode: %v, errorMessage: %v",
			data["ResponseCode"], data["ResponseDescription"], data["errorCode"], data["ResponseCode"],
		)
	}
}

// TimeoutPayment processes tax payment timeout
func TimeoutPayment(id string) {
	payment, err := sql.FindOneById[models.TaxPayment](id)
	if err != nil {
		logs.Error("could not find payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":         payment.ID,
		"status":     2,
		"request_id": payment.RequestId,
		"amount":     payment.Amount,
		"reference":  payment.PaymentRequestNumber,
	}

	if err := utils.NotifyClient(payment.CallbackUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}

// ResultPayment processes tax payment result
func ResultPayment(id string, req *domain.TaxResult) {
	payment, err := sql.FindOneById[models.TaxPayment](id)
	if err != nil {
		logs.Error("could not find payment: %v", err)
		return
	}

	meta := map[string]any{}
	for _, item := range req.Result.ResultParameters.ResultParameter {
		meta[item.Key] = item.Value
	}

	payment.Meta = meta
	payment.ResultCode = req.Result.ResultCode
	payment.ResultDescription = req.Result.ResultDesc
	payment.TransactionID = req.Result.TransactionID

	if err := sql.Update(payment); err != nil {
		logs.Error("could not update payment: %v", err)
		return
	}

	if req.Result.ResultCode != 0 {
		logs.Error(
			"failed to initiate payment =>ResultCode: %v, ResultDescription: %v",
			req.Result.ResultCode, req.Result.ResultDesc,
		)
		return
	}

	payload := map[string]any{
		"id":         payment.ID,
		"status":     req.Result.ResultCode,
		"request_id": payment.RequestId,
		"amount":     payment.Amount,
		"reference":  payment.PaymentRequestNumber,
		"message":    req.Result.ResultDesc,
	}

	if err := utils.NotifyClient(payment.CallbackUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}
