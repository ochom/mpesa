package auth

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ochom/gutils/cache"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/utils"
)

func Authenticate(tokenName, username, password string) string {
	cached := cache.Get(tokenName)
	if cached != nil {
		return string(cached)
	}

	headers := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
	}

	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", config.MpesaApiUrl)
	res, err := gttp.Get(url, headers)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return ""
	}

	tokens := helpers.FromBytes[map[string]string](res.Body)
	token := tokens["access_token"]
	if token == "" {
		logs.Error("failed to authenticate: %v", string(res.Body))
		return ""
	}

	cache.SetWithExpiry(tokenName, []byte(token), 59*time.Minute)
	return token
}

func GetSecurityCredentials(tokenName, certPath, password string) string {
	cached := cache.Get(tokenName)
	if cached != nil {
		return string(cached)
	}

	encoded := utils.HashText(certPath, password)
	cache.Set(tokenName, []byte(encoded))
	return encoded
}
