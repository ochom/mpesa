package config

import "github.com/ochom/gutils/env"

var (
	// BaseUrl used for callbacks to this application
	BaseUrl = env.Get("BASE_URL")

	// MpesaApiUrl to make requests to Mpesa API
	MpesaApiUrl = env.Get("MPESA_API_URL")
)
