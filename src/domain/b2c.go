package domain

type B2cRequest struct {
	RequestId   string `json:"request_id"`
	PhoneNumber string `json:"phone_number"`
	Amount      string `json:"amount"`
	CallbackUrl string `json:"callback_url"`
}

type B2cResult struct {
	Result struct {
		ResultType               int    `json:"ResultType"`
		ResultCode               int    `json:"ResultCode"`
		ResultDesc               string `json:"ResultDesc"`
		OriginatorConversationID string `json:"OriginatorConversationID"`
		ConversationID           string `json:"ConversationID"`
		TransactionID            string `json:"TransactionID"`
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			}
		} `json:"ResultParameters"`
		ReferenceData struct {
			ReferenceItem []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			}
		} `json:"ReferenceData"`
	} `json:"Result"`
}
