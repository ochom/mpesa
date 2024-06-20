package auth

import (
	"fmt"
	"time"

	"github.com/ochom/gutils/cache"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/models"
	"github.com/ochom/mpesa/src/utils"
)

func Authenticate(account *models.Account) string {
	tokenName := fmt.Sprintf("mpesa_%s_token_%d", account.Type, account.ID)
	cached := cache.Get(tokenName)
	if cached != nil {
		return string(cached)
	}

	headers := map[string]string{
		"Authorization": "Basic " + utils.Encode([]byte(account.ConsumerKey+":"+account.ConsumerSecrete)),
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

func GetSecurityCredentials(account *models.Account) string {
	tokenName := fmt.Sprintf("mpesa_%s_token_%d", account.Type, account.ID)
	cached := cache.Get(tokenName)
	if cached != nil {
		return string(cached)
	}

	encoded := utils.HashText(account.Certificate, account.InitiatorPassword)
	cache.Set(tokenName, []byte(encoded))
	return encoded
}
