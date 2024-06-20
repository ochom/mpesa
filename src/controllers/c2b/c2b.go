package c2b

import (
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
	"gorm.io/gorm"
)

// hash hashes the short_code, passkey and timestamp
func hash(shortCode, passKey, timeStamp string) string {
	join := shortCode + passKey + timeStamp
	return utils.Encode([]byte(join))
}

// RegisterUrls registers c2b url
func RegisterUrls(req map[string]string) {
	account, err := sql.FindOneById[models.Account](req["account_id"])
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return
	}

	headers := map[string]string{
		"Authorization": "Bearer " + auth.Authenticate(account),
		"Content-Type":  "application/json",
	}

	payload := map[string]any{
		"ShortCode":       account.ShortCode,
		"ResponseType":    "Completed",
		"ConfirmationURL": req["confirmation_url"],
		"ValidationURL":   req["validation_url"],
	}

	url := fmt.Sprintf("%s/mpesa/c2b/v2/registerurl", config.MpesaApiUrl)
	res, err := gttp.Post(url, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return
	}

	if res.Status > 204 {
		logs.Error("failed to register url: %v", string(res.Body))
		return
	}

	logs.Info("res: %v", string(res.Body))
}

// InitiatePayment initiates an mpesa c2b stk push
func InitiatePayment(req *domain.MpesaExpressRequest) {
	account, err := sql.FindOneById[models.Account](req.AccountId)
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return
	}

	refId := uuid.New()
	cache.SetWithExpiry(fmt.Sprintf("stk-%s", refId), helpers.ToBytes(req), 5*time.Minute)

	timestamp := time.Now().Format("20060102150405")
	callbackUrl := fmt.Sprintf("%s/v1/c2b/result?refId=%s", config.BaseUrl, refId)

	payload := map[string]string{
		"BusinessShortCode": account.ShortCode,
		"Password":          hash(account.ShortCode, account.PassKey, timestamp),
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            req.Amount,
		"PartyA":            req.PhoneNumber,
		"PartyB":            account.ShortCode,
		"PhoneNumber":       req.PhoneNumber,
		"CallBackURL":       callbackUrl,
		"AccountReference":  "Customer Pay Bill Online",
		"TransactionDesc":   "Pay bill",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + auth.Authenticate(account),
		"Content-Type":  "application/json",
	}

	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", config.MpesaApiUrl)
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
	if data["ResponseCode"] != "0" {
		logs.Error("initiate stk failed: ResponseDescription=>%s", data["ResponseDescription"])
		return
	}
}

// ResultPayment processes the payment result for stk push
func ResultPayment(id string, req *domain.MpesaExpressCallback) {
	cacheBytes := cache.Get(fmt.Sprintf("stk-%s", id))
	if cacheBytes == nil {
		logs.Error("failed to get stk payment cache")
		return
	}

	if req.Body.StkCallback.ResultCode != 0 {
		logs.Error("failed to process payment: %v", req.Body.StkCallback.ResultDesc)
		return
	}

	meta := map[string]any{}
	for _, item := range req.Body.StkCallback.CallbackMetadata.Item {
		meta[item.Name] = item.Value
	}

	cacheData := helpers.FromBytes[domain.MpesaExpressRequest](cacheBytes)

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
		"id":           customerPayment.ID,
		"status":       req.Body.StkCallback.ResultCode,
		"message":      req.Body.StkCallback.ResultDesc,
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
	account, err := sql.FindOne[models.Account](func(d *gorm.DB) *gorm.DB {
		return d.Where("short_code = ?", req.BusinessShortCode)
	})

	if err != nil {
		logs.Error("failed to find account: %v", err)
		return false
	}

	if account.ValidationUrl == "" {
		return true
	}

	if err := utils.NotifyClient(account.ValidationUrl, req); err != nil {
		logs.Error("failed to notify client: %v", err)
		return false
	}

	return true
}

// ConfirmPayment confirms payments received through REST API
func ConfirmPayment(req *domain.ValidationRequest) {
	account, err := sql.FindOne[models.Account](func(d *gorm.DB) *gorm.DB {
		return d.Where("short_code = ?", req.BusinessShortCode)
	})
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return
	}

	customerPayment := models.NewCustomerPayment(req.TransID, req.TransTime, req.TransAmount, req.BillRefNumber, req.InvoiceNumber, req.MSISDN)
	if err := sql.Create(customerPayment); err != nil {
		logs.Error("could not create this payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":           customerPayment.ID,
		"status":       0,
		"message":      "Payment confirmed",
		"amount":       customerPayment.TransactionAmount,
		"phone_number": customerPayment.PhoneNumber,
		"reference":    customerPayment.TransactionID,
	}

	if err := utils.NotifyClient(account.ConfirmationUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}
