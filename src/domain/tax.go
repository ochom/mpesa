package domain

// TaxRequest the payload required to initiate an mpesa stk push
type TaxRequest struct {
	AccountId            string `json:"account_id" validate:"required"`
	ShortCode            string `json:"short_code" validate:"required"`
	RequestId            string `json:"request_id" validate:"required"`
	Amount               string `json:"amount" validate:"required"`
	PaymentRequestNumber string `json:"prn" validate:"required"`
	CallbackUrl          string `json:"callback_url" validate:"required"`
}

// TaxResult the payload required to initiate an mpesa stk push
type TaxResult struct {
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
