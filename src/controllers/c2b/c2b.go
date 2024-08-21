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
func InitiatePayment(req *domain.MpesaExpressRequest) error {
	account, err := sql.FindOneById[models.Account](req.AccountId)
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return fmt.Errorf("c2b account not found")
	}

	refId := uuid.New()
	if err := cache.SetWithExpiry(fmt.Sprintf("stk-%s", refId), req, 5*time.Minute); err != nil {
		logs.Error("failed to set cache: %v", err)
		return fmt.Errorf("failed to set cache")
	}

	timestamp := time.Now().Format("20060102150405")
	callbackUrl := fmt.Sprintf("%s/v1/c2b/result?refId=%s", config.BaseUrl, refId)

	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", config.MpesaApiUrl)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", auth.Authenticate(account)),
		"Content-Type":  "application/json",
	}

	payload := map[string]string{
		"BusinessShortCode": account.ShortCode,
		"Password":          utils.Encode([]byte(account.ShortCode + account.PassKey + timestamp)),
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            req.Amount,
		"PartyA":            req.PhoneNumber,
		"PartyB":            account.ShortCode,
		"PhoneNumber":       req.PhoneNumber,
		"CallBackURL":       callbackUrl,
		"AccountReference":  req.InvoiceNumber,
		"TransactionDesc":   "Pay bill",
	}

	res, err := gttp.Post(url, headers, payload)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return fmt.Errorf("failed to make request")
	}

	if res.Status > 204 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return fmt.Errorf("request failed")
	}

	data := helpers.FromBytes[map[string]string](res.Body)
	if data["ResponseCode"] != "0" {
		logs.Error("initiate stk failed: ResponseDescription=>%s", data["ResponseDescription"])
		return fmt.Errorf("initiate stk failed")
	}

	return nil
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
