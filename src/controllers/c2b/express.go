package c2b

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ochom/gutils/cache"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/app"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
)

func authenticate() string {
	tokenBytes := cache.Get("mpesa_c2b_token")
	if tokenBytes != nil {
		return string(tokenBytes)
	}

	password := []byte(app.MpesaC2bConsumerKey + ":" + app.MpesaC2bConsumerSecret)
	encoded := base64.StdEncoding.EncodeToString(password)
	headers := map[string]string{
		"Authorization": "Basic " + encoded,
	}

	res, err := gttp.Get(app.MpesaC2bAuthUrl, headers)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return ""
	}

	if res.Status > 201 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return ""
	}

	tokens := helpers.FromBytes[map[string]string](res.Body)
	if len(tokens) == 0 {
		logs.Error("failed to authenticate: %v", string(res.Body))
		return ""
	}

	token := tokens["access_token"]
	cache.SetWithExpiry("mpesa_c2b_token", []byte(token), 59*time.Minute)
	return token
}

func getPassword(shortCode, passKey, timeStamp string) string {
	join := shortCode + passKey + timeStamp
	return base64.StdEncoding.EncodeToString([]byte(join))
}

func notifyClient(url string, err error) {
	if url == "" {
		return
	}

	payload := map[string]any{
		"status":  0,
		"message": "successful",
	}

	if err != nil {
		payload["status"] = 1
		payload["message"] = err.Error()
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := gttp.Post(url, headers, payload)
	if err != nil {
		logs.Error("failed to notify client: %v", err)
		return
	}

	if res.Status > 201 {
		logs.Error("failed to notify client: %v", string(res.Body))
	}
}

func MpesaExpressInitiate(req domain.MpesaExpressRequest) {
	mpe := models.NewMpesaExpress(req.PhoneNumber, req.Amount, req.CallbackUrl, req.AccountReference)
	if err := sql.Create(mpe); err != nil {
		logs.Error("failed to create mpesa express: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	phoneNumber := helpers.ParseMobile(req.PhoneNumber)
	callbackUrl := fmt.Sprintf("%s?id=%d", app.BaseUrl, mpe.Id)

	payload := map[string]string{
		"BusinessShortCode": app.MpesaC2bShortCode,
		"Password":          getPassword(app.MpesaC2bShortCode, app.MpesaC2bPassKey, timestamp),
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            "1",
		"PartyA":            phoneNumber,
		"PartyB":            app.MpesaC2bShortCode,
		"PhoneNumber":       phoneNumber,
		"CallBackURL":       callbackUrl,
		"AccountReference":  req.AccountReference,
		"TransactionDesc":   "Pay bill",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + authenticate(),
		"Content-Type":  "application/json",
	}

	res, err := gttp.Post(app.MpesaC2bApiUrl, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		notifyClient(req.CallbackUrl, err)
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
		notifyClient(req.CallbackUrl, fmt.Errorf("%s", data["ResponseDescription"]))
		return
	}

	logs.Info("success: %v", string(res.Body))
}

func MpesaExpressCallback(id string, req domain.MpesaExpressCallback) {
	mpe, err := sql.FindOneById[models.MpesaExpress](id)
	if err != nil {
		logs.Error("failed to find mpesa express: %v", err)
		return
	}

	if req.Body.StkCallback.ResultCode != 0 {
		notifyClient(mpe.CallbackUrl, fmt.Errorf("%s", req.Body.StkCallback.ResultDescription))
		return
	}

	meta := map[string]any{}
	for _, item := range req.Body.StkCallback.CallbackMetadata.Item {
		meta[item.Name] = item.Value
	}

	mpe.Meta = meta
	if err := sql.Update(mpe); err != nil {
		logs.Error("failed to update mpesa express: %v", err)
	}

	notifyClient(mpe.CallbackUrl, nil)
}