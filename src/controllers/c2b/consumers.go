package c2b

import (
	"fmt"
	"time"

	"github.com/ochom/gutils/cache"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
	"github.com/ochom/mpesa/src/utils"
	"gorm.io/gorm"
)

var messages chan map[string]any

func init() {
	messages = make(chan map[string]any, 20)

	go func() {
		logs.Info("listening for c2b payments results...")
		for msg := range messages {
			switch msg["message_type"] {
			case "callback":
				paymentID := msg["id"].(string)
				reqJson := helpers.ToBytes(msg["message"])
				req := helpers.FromBytes[domain.MpesaExpressCallback](reqJson)
				resultPayment(paymentID, &req)
			case "confirmation":
				reqJson := helpers.ToBytes(msg["message"])
				req := helpers.FromBytes[domain.ValidationRequest](reqJson)
				confirmPayment(&req)
			default:
				logs.Error("unknown message type: %v", msg["message_type"])
			}
		}
	}()
}

// AddMessage adds a message to the queue
func AddMessage(msg map[string]any) {
	messages <- msg
}

// resultPayment processes the payment result for stk push
func resultPayment(id string, req *domain.MpesaExpressCallback) {
	cacheData, err := cache.Get[domain.MpesaExpressRequest](fmt.Sprintf("stk-%s", id))
	if err != nil {
		logs.Error("failed to get stk payment cache: %v", err)
		return
	}

	if req.Body.StkCallback.ResultCode != 0 {
		logs.Error("failed to process payment: %v", req.Body.StkCallback.ResultDesc)
		return
	}

	account, err := sql.FindOne[models.Account](func(d *gorm.DB) *gorm.DB {
		return d.Where("id = ?", cacheData.AccountId)
	})
	if err != nil {
		logs.Error("failed to find account: %v", err)
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
	invoiceNumber := cacheData.InvoiceNumber

	customerPayment := models.NewCustomerPayment(account.ID, txId, txTime, txAmount, billRefNumber, invoiceNumber, billRefNumber)
	if err := customerPayment.Save(); err != nil {
		logs.Error("could not create this payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":           customerPayment.ID,
		"status":       req.Body.StkCallback.ResultCode,
		"message":      req.Body.StkCallback.ResultDesc,
		"amount":       customerPayment.Amount,
		"phone_number": customerPayment.PhoneNumber,
		"reference":    customerPayment.TransactionID,
	}

	if err := utils.NotifyClient(cacheData.CallbackUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}

// ConfirmPayment confirms payments received through REST API
func confirmPayment(req *domain.ValidationRequest) {
	account, err := sql.FindOne[models.Account](func(d *gorm.DB) *gorm.DB {
		return d.Where("short_code = ?", req.BusinessShortCode)
	})
	if err != nil {
		logs.Error("failed to find account: %v", err)
		return
	}

	customerPayment := models.NewCustomerPayment(account.ID, req.TransID, req.TransTime, req.TransAmount, req.BillRefNumber, req.InvoiceNumber, req.MSISDN)
	if err := customerPayment.Save(); err != nil {
		logs.Error("could not create this payment: %v", err)
		return
	}

	payload := map[string]any{
		"id":           customerPayment.ID,
		"status":       0,
		"message":      "Payment confirmed",
		"amount":       customerPayment.Amount,
		"phone_number": customerPayment.PhoneNumber,
		"reference":    customerPayment.TransactionID,
	}

	if err := utils.NotifyClient(account.ConfirmationUrl, payload); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}
