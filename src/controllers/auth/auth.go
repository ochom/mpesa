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

	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", config.MpesaAuthUrl)
	res, err := gttp.Get(url, headers)
	if err != nil {
		logs.Error("failed to make request: %v", err)
		return ""
	}

	if res.Status > 201 {
		logs.Error("request failed status: %d body: %v", res.Status, string(res.Body))
		return ""
	}

	tokens := helpers.FromBytes[map[string]string](res.Body)
	if len(tokens) == 0 {
		logs.Error("failed to authenticate: %v", string(res.Body))
		return ""
	}

	token := tokens["access_token"]
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
