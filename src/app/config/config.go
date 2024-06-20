package config

import "github.com/ochom/gutils/env"

var (
	// BaseUrl used for callbacks to this application
	DbDriver    = env.Int("DATABASE_DRIVER", 0)
	DbUrl       = env.Get("DATABASE_URL")
	BaseUrl     = env.Get("BASE_URL")
	MpesaApiUrl = env.Get("MPESA_API_URL")

	BasicAuthUsername = env.Get("BASIC_AUTH_USERNAME")
	BasicAuthPassword = env.Get("BASIC_AUTH_PASSWORD")

	//MpesaB2C credentials
	B2CAllowedOrigins = env.Get("B2C_ALLOWED_ORIGINS")

	// Tax Remittance credentials
	TaxAllowedOrigins = env.Get("TAX_REMITTANCE_ALLOWED_ORIGINS")
)
