package app

import "github.com/ochom/gutils/env"

var (
	// BaseUrl used for callbacks to this application
	BaseUrl = env.Get("BASE_URL")

	// MpesaC2b credentials
	MpesaC2bAuthUrl              = env.Get("MPESA_C2B_AUTH_URL")
	MpesaC2bApiUrl               = env.Get("MPESA_C2B_API_URL")
	MpesaC2bShortCode            = env.Get("MPESA_C2B_SHORT_CODE")
	MpesaC2bPassKey              = env.Get("MPESA_C2B_PASSKEY")
	MpesaC2bConsumerKey          = env.Get("MPESA_C2B_CONSUMER_KEY")
	MpesaC2bConsumerSecret       = env.Get("MPESA_C2B_CONSUMER_SECRET")
	ClientDepositValidationUrl   = env.Get("CLIENT_DEPOSIT_VALIDATION_URL")
	ClientDepositConfirmationUrl = env.Get("CLIENT_DEPOSIT_CONFIRMATION_URL")
)
