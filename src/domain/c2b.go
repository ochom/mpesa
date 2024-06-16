package domain

// MpesaExpressRequest the payload required to initiate an mpesa stk push
type MpesaExpressRequest struct {
	Amount           string `json:"amount"`
	PhoneNumber      string `json:"phone_number"`
	AccountReference string `json:"account_reference"`
	CallbackUrl      string `json:"callback_url"`
}

// MpesaExpressCallback the response from an mpesa stk push
type MpesaExpressCallback struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDescription string `json:"ResultDescription"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string `json:"Name"`
					Value any    `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// ValidationRequest store mpesa stk-push requests and respective responses
type ValidationRequest struct {
	TransactionType   string
	TransID           string
	TransTime         string
	TransAmount       string
	BusinessShortCode string
	BillRefNumber     string
	InvoiceNumber     string
	OrgAccountBalance string
	ThirdPartyTransID string
	MSISDN            string
	FirstName         string
	MiddleName        string
	LastName          string
}
