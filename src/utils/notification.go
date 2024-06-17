package utils

import (
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

func NotifyClient(url string, payload any) {
	if url == "" {
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := gttp.Post(url, headers, helpers.ToBytes(payload))
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
