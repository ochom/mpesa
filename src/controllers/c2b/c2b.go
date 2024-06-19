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
	"github.com/ochom/gutils/uuid"
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
	refId := uuid.New()
	cache.SetWithExpiry(fmt.Sprintf("stk-%s", refId), helpers.ToBytes(req), 5*time.Minute)

	timestamp := time.Now().Format("20060102150405")
	phoneNumber := helpers.ParseMobile(req.PhoneNumber)
	callbackUrl := fmt.Sprintf("%s/v1/c2b/result?refId=%s", config.BaseUrl, refId)

	payload := map[string]string{
		"BusinessShortCode": config.MpesaC2BShortCode,
		"Password":          hash(config.MpesaC2BShortCode, config.MpesaC2BPassKey, timestamp),
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            req.Amount,
		"PartyA":            phoneNumber,
		"PartyB":            config.MpesaC2BShortCode,
		"PhoneNumber":       phoneNumber,
		"CallBackURL":       callbackUrl,
		"AccountReference":  "Customer Pay Bill Online",
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

	if res.Status > 201 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return
	}

	data := helpers.FromBytes[map[string]string](res.Body)
	if data["ResponseCode"] != "0" {
		logs.Error("initiate stk failed: ResponseDescription=>%s", data["ResponseDescription"])
		return
	}
}

func ResultPayment(id string, req *domain.MpesaExpressCallback) {
	cacheBytes := cache.Get(fmt.Sprintf("stk-%s", id))
	if cacheBytes == nil {
		logs.Error("failed to get stk payment cache")
		return
	}

	cacheData := helpers.FromBytes[domain.MpesaExpressRequest](cacheBytes)

	if req.Body.StkCallback.ResultCode != 0 {
		logs.Error("failed to process payment: %v", req.Body.StkCallback.ResultDescription)
		return
	}

	meta := map[string]any{}
	for _, item := range req.Body.StkCallback.CallbackMetadata.Item {
		meta[item.Name] = item.Value
	}

	txId := meta["MpesaReceiptNumber"].(string)
	txTime := time.Now().Format("20060102150405")
	txAmount := cacheData.Amount
	billRefNumber := cacheData.PhoneNumber
	invoiceNumber := req.Body.StkCallback.MerchantRequestID

	customerPayment := models.NewCustomerPayment(txId, txTime, txAmount, billRefNumber, invoiceNumber, billRefNumber)
	if err := sql.Create(customerPayment); err != nil {
		logs.Error("could not create this payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":           customerPayment.Id,
		"status":       req.Body.StkCallback.ResultCode,
		"message":      req.Body.StkCallback.ResultDescription,
		"amount":       customerPayment.TransactionAmount,
		"phone_number": customerPayment.PhoneNumber,
		"reference":    customerPayment.TransactionID,
	}

	if err := utils.NotifyClient(cacheData.CallbackUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}

// ValidatePayment  validates payments received through REST API
func ValidatePayment(req *domain.ValidationRequest) bool {
	if config.MpesaC2BValidationUrl == "" {
		return true
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	payload := helpers.ToBytes(req)
	res, err := gttp.Post(config.MpesaC2BValidationUrl, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return false
	}

	if res.Status > 201 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return false
	}

	return true
}

// ConfirmPayment confirms payments received through REST API
func ConfirmPayment(req *domain.ValidationRequest) {
	customerPayment := models.NewCustomerPayment(req.TransID, req.TransTime, req.TransAmount, req.BillRefNumber, req.InvoiceNumber, req.MSISDN)
	if err := sql.Create(customerPayment); err != nil {
		logs.Error("could not create this payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":           customerPayment.Id,
		"status":       0,
		"message":      "Payment confirmed",
		"amount":       customerPayment.TransactionAmount,
		"phone_number": customerPayment.PhoneNumber,
		"reference":    customerPayment.TransactionID,
	}

	if err := utils.NotifyClient(config.MpesaC2BConfirmationUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}
