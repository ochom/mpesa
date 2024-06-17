package c2b

import (
	"encoding/base64"
	"fmt"
	"time"

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

func hash(shortCode, passKey, timeStamp string) string {
	join := shortCode + passKey + timeStamp
	return base64.StdEncoding.EncodeToString([]byte(join))
}

func InitiatePayment(req *domain.MpesaExpressRequest) {
	mpe := models.NewCustomerPayment(req.PhoneNumber, req.Amount, req.CallbackUrl, req.AccountReference)
	if err := sql.Create(mpe); err != nil {
		logs.Error("failed to create mpesa express: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	phoneNumber := helpers.ParseMobile(req.PhoneNumber)
	callbackUrl := fmt.Sprintf("%s/c2b/result?id=%d", config.BaseUrl, mpe.Id)

	payload := map[string]string{
		"BusinessShortCode": config.MpesaC2BShortCode,
		"Password":          hash(config.MpesaC2BShortCode, config.MpesaC2BPassKey, timestamp),
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            "1",
		"PartyA":            phoneNumber,
		"PartyB":            config.MpesaC2BShortCode,
		"PhoneNumber":       phoneNumber,
		"CallBackURL":       callbackUrl,
		"AccountReference":  req.AccountReference,
		"TransactionDesc":   "Pay bill",
	}

	username := config.MpesaC2BConsumerKey
	password := config.MpesaC2BConsumerSecret

	headers := map[string]string{
		"Authorization": "Bearer " + auth.Authenticate("mpesa_c2b_token", username, password),
		"Content-Type":  "application/json",
	}

	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", config.MpesaApiUrl)
	res, err := gttp.Post(url, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return
	}

	defer func() {
		if err := sql.Update(mpe); err != nil {
			logs.Error("failed to create mpesa express: %v", err)
		}
	}()

	data := helpers.FromBytes[map[string]string](res.Body)
	mpe.MerchantRequestId = data["MerchantRequestID"]
	mpe.CheckoutRequestId = data["CheckoutRequestID"]
	mpe.ResponseCode = data["ResponseCode"]
	mpe.ResponseDescription = data["ResponseDescription"]

	if data["ResponseCode"] != "0" {
		return
	}

	logs.Info("success: %v", string(res.Body))
}

func ResultPayment(id string, req *domain.MpesaExpressCallback) {
	mpe, err := sql.FindOneById[models.CustomerPayment](id)
	if err != nil {
		logs.Error("failed to find mpesa express: %v", err)
		return
	}

	if req.Body.StkCallback.ResultCode != 0 {
		logs.Error("failed to process payment: %v", req.Body.StkCallback.ResultDescription)
	}

	meta := map[string]any{}
	for _, item := range req.Body.StkCallback.CallbackMetadata.Item {
		meta[item.Name] = item.Value
	}

	mpe.Meta = meta
	mpe.ResultCode = req.Body.StkCallback.ResultCode
	mpe.ResultDescription = req.Body.StkCallback.ResultDescription
	if err := sql.Update(mpe); err != nil {
		logs.Error("failed to update mpesa express: %v", err)
	}

	payload := map[string]any{
		"id":           mpe.Id,
		"status":       req.Body.StkCallback.ResultCode,
		"message":      req.Body.StkCallback.ResultDescription,
		"amount":       mpe.Amount,
		"phone_number": mpe.PhoneNumber,
		"reference":    mpe.Meta.Get("MpesaReceiptNumber"),
	}

	if err := utils.NotifyClient(mpe.CallbackUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}
