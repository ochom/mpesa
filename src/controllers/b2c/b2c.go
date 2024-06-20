package b2c

import (
	"fmt"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/controllers/auth"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
	"github.com/ochom/mpesa/src/utils"
)

// InitiatePayment initiates an mpesa b2c payment
func InitiatePayment(req domain.B2cRequest) {
	account, err := sql.FindOneById[models.Account](req.AccountId)
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return
	}

	payment := models.NewBusinessPayment(account.ID, req.RequestId, req.PhoneNumber, req.Amount, req.CallbackUrl)
	if err := sql.Create(payment); err != nil {
		logs.Error("Error creating payment: %v", err)
		return
	}

	headers := map[string]string{
		"Authorization": "Bearer " + auth.Authenticate(account),
		"Content-Type":  "application/json",
	}

	payload := map[string]string{
		"OriginatorConversationID": uuid.New(),
		"InitiatorName":            account.InitiatorName,
		"SecurityCredential":       auth.GetSecurityCredentials(account),
		"CommandID":                "BusinessPayment",
		"Amount":                   req.Amount,
		"PartyA":                   account.ShortCode,
		"PartyB":                   payment.PhoneNumber,
		"Remarks":                  "OK",
		"QueueTimeOutURL":          fmt.Sprintf("%s/v1/b2c/timeout?id=%d", config.BaseUrl, payment.ID),
		"ResultURL":                fmt.Sprintf("%s/v1/b2c/result?id=%d", config.BaseUrl, payment.ID),
		"Occassion":                "Payout",
	}

	url := fmt.Sprintf("%s/mpesa/b2c/v3/paymentrequest", config.MpesaApiUrl)
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

// TimeoutPayment processes b2c payment timeout
func TimeoutPayment(id string) {
	payment, err := sql.FindOneById[models.BusinessPayment](id)
	if err != nil {
		logs.Error("could not find payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":         payment.ID,
		"status":     2,
		"request_id": payment.RequestId,
		"amount":     payment.Amount,
	}

	if err := utils.NotifyClient(payment.CallbackUrl, payload); err != nil {
		logs.Error("could not notify client: %v", err)
	}
}

// ResultPayment processes b2c payment result
func ResultPayment(id string, req *domain.B2cResult) {
	payment, err := sql.FindOneById[models.BusinessPayment](id)
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
		"message":    req.Result.ResultDesc,
		"request_id": payment.RequestId,
		"amount":     payment.Amount,
		"reference":  payment.Meta.Get("TransactionReceipt"),
	}

	if err := utils.NotifyClient(payment.CallbackUrl, payload); err != nil {
		logs.Error("could not notify client: %v", err)
	}
}
