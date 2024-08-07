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
	logs.Info("getting token: %s", tokenName)

	tokens, err := cache.Get[map[string]string](tokenName)
	if err != nil {
		logs.Warn("token not found, generating a new token")
		return setToken(account, tokenName)
	}

	if tokens["access_token"] == "" {
		logs.Warn("empty token, generating a new token")
		return setToken(account, tokenName)
	}

	return tokens["access_token"]
}

func setToken(account *models.Account, tokenName string) string {
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

	if err := cache.SetWithExpiry(tokenName, tokens, config.MpesaTokenExpiry); err != nil {
		logs.Error("failed to set token: %v", err)
		return ""
	}

	return token
}

func GetSecurityCredentials(account *models.Account) string {
	tokenName := fmt.Sprintf("mpesa_%s_token_%d", account.Type, account.ID)
	cached, err := cache.Get[map[string]string](tokenName)
	if err != nil {
		return setSecurityToken(account, tokenName)
	}

	if cached["security_token"] == "" {
		return setSecurityToken(account, tokenName)
	}

	return cached["security_token"]
}

func setSecurityToken(account *models.Account, tokenName string) string {
	encoded := utils.HashText(account.Certificate, account.InitiatorPassword)
	data := map[string]string{
		"security_token": encoded,
	}

	if err := cache.SetWithExpiry(tokenName, data, 50*time.Minute); err != nil {
		logs.Error("failed to set token: %v", err)
	}

	return encoded
}
