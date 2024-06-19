package utils

import (
	"fmt"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
)

func NotifyClient(url string, payload any) error {
	if url == "" {
		return fmt.Errorf("no url provided")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := gttp.Post(url, headers, helpers.ToBytes(payload))
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	if res.Status > 204 {
		return fmt.Errorf("request failed status: %d body: %v", res.Status, string(res.Body))
	}

	return nil
}
