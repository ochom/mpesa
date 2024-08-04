package config

import (
	"time"

	"github.com/ochom/gutils/env"
)

var (
	// BaseUrl used for callbacks to this application
	BaseUrl = env.Get("BASE_URL")

	// MpesaApiUrl to make requests to Mpesa API
	MpesaApiUrl = env.Get("MPESA_API_URL")

	// MpesaTokenExpiry how long the token is valid
	MpesaTokenExpiry = time.Duration(env.Int("MPESA_TOKEN_EXPIRY_MINUTES", 5)) * time.Minute
)
