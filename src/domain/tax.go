package domain

// TaxRequest the payload required to initiate an mpesa stk push
type TaxRequest struct {
	RequestId            string `json:"request_id"`
	Amount               string `json:"amount"`
	PaymentRequestNumber string `json:"prn"`
	ShortCode            string `json:"short_code"`
	CallbackUrl          string `json:"callback_url"`
}
