package c2b

import (
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/utils"
)

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
	if config.MpesaC2BConfirmationUrl == "" {
		return
	}

	if err := utils.NotifyClient(config.MpesaC2BConfirmationUrl, req); err != nil {
		logs.Error("failed to notify client: %v", err)
	}
}
