package b2c

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
)

func notifyClient(url string, data *models.Payment) {
	if url == "" {
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := gttp.Post(url, headers, helpers.ToBytes(data))
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return
	}

	if res.Status > 201 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return
	}

	logs.Info("request successful status: %d body: %v", res.Status, string(res.Body))
}

func InitiatePayment(req domain.B2cRequest) {
	payment := models.NewPayment(req.RequestId, req.PhoneNumber, req.Amount, req.CallbackUrl)
	if err := sql.Create(payment); err != nil {
		logs.Error("Error creating payment: %v", err)
		return
	}

	username := config.MpesaB2CConsumerKey
	password := config.MpesaB2CConsumerSecret

	headers := map[string]string{
		"Authorization": "Bearer " + auth.Authenticate("mpesa_b2c_token", username, password),
		"Content-Type":  "application/json",
	}

	certPath := config.MpesaB2CCertificatePath
	initiatorPassword := config.MpesaB2CInitiatorPassword
	securityCredential := auth.GetSecurityCredentials("mpesa_b2c_security", certPath, initiatorPassword)

	payload := map[string]string{
		"OriginatorConversationID": payment.Id,
		"InitiatorName":            config.MpesaB2CInitiatorName,
		"SecurityCredential":       securityCredential,
		"CommandID":                "BusinessPayment",
		"Amount":                   req.Amount,
		"PartyA":                   config.MpesaB2CShortCode,
		"PartyB":                   payment.PhoneNumber,
		"Remarks":                  config.MpesaB2CPaymentComment,
		"QueueTimeOutURL":          config.MpesaB2CQueueTimeoutUrl,
		"ResultURL":                config.MpesaB2CResultUrl,
		"Occassion":                "Payout",
	}

	url := fmt.Sprintf("%s/mpesa/b2c/v3/paymentrequest", config.MpesaApiUrl)
	res, err := gttp.Post(url, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return
	}

	if res.Status > 201 {
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

func ResultBusinessPayment(id string, req *domain.B2cResult) {
	payment, err := sql.FindOneById[models.Payment](id)
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

	notifyClient(payment.CallbackUrl, nil)
}
