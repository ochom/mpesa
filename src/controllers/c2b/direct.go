package c2b

import (
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app"
	"github.com/ochom/mpesa/src/domain"
)

// MpesaExpressValidate ...
func MpesaExpressValidate(req *domain.ValidationRequest) bool {
	if app.ClientDepositValidationUrl == "" {
		return true
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	payload := helpers.ToBytes(req)
	res, err := gttp.Post(app.ClientDepositValidationUrl, headers, payload)
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

// MpesaExpressConfirm ...
func MpesaExpressConfirm(req *domain.ValidationRequest) {
	if app.ClientDepositConfirmationUrl == "" {
		return
	}

	payload := helpers.ToBytes(req)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := gttp.Post(app.ClientDepositConfirmationUrl, headers, payload)
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
