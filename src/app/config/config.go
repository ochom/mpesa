package config

import "github.com/ochom/gutils/env"

var (
	// BaseUrl used for callbacks to this application
	DbDriver    = env.Int("DATABASE_DRIVER", 0)
	DbUrl       = env.Get("DATABASE_URL")
	BaseUrl     = env.Get("BASE_URL")
	MpesaApiUrl = env.Get("MPESA_API_URL")

	// MpesaC2B credentials
	MpesaC2BShortCode            = env.Get("MPESA_C2B_SHORT_CODE")
	MpesaC2BPassKey              = env.Get("MPESA_C2B_PASSKEY")
	MpesaC2BConsumerKey          = env.Get("MPESA_C2B_CONSUMER_KEY")
	MpesaC2BConsumerSecret       = env.Get("MPESA_C2B_CONSUMER_SECRET")
	ClientDepositValidationUrl   = env.Get("CLIENT_DEPOSIT_VALIDATION_URL")
	ClientDepositConfirmationUrl = env.Get("CLIENT_DEPOSIT_CONFIRMATION_URL")

	//MpesaB2C credentials
	B2CAllowedOrigins         = env.Get("B2C_ALLOWED_ORIGINS")
	MpesaB2CShortCode         = env.Get("MPESA_B2C_SHORT_CODE")
	MpesaB2CPassKey           = env.Get("MPESA_B2C_PASSKEY")
	MpesaB2CConsumerKey       = env.Get("MPESA_B2C_CONSUMER_KEY")
	MpesaB2CConsumerSecret    = env.Get("MPESA_B2C_CONSUMER_SECRET")
	MpesaB2CInitiatorName     = env.Get("MPESA_B2C_INITIATOR_NAME")
	MpesaB2CInitiatorPassword = env.Get("MPESA_B2C_INITIATOR_PASSWORD")
	MpesaB2CCertificatePath   = env.Get("MPESA_B2C_CERTIFICATE_PATH")
	MpesaB2CPaymentComment    = env.Get("MPESA_B2C_PAYMENT_COMMENT")

	// Tax Remittance credentials
	TaxAllowedOrigins         = env.Get("TAX__REMITTANCE_ALLOWED_ORIGINS")
	MpesaTaxConsumerSecrete   = env.Get("MPESA_TAX_CONSUMER_SECRET")
	MpesaTaxConsumerKey       = env.Get("MPESA_TAX_CONSUMER_KEY")
	MpesaTaxInitiatorPassword = env.Get("MPESA_TAX_INITIATOR_PASSWORD")
	MpesaTaxCertificatePath   = env.Get("MPESA_TAX_CERTIFICATE_PATH")
)
