package domain

// TaxRequest the payload required to initiate an mpesa stk push
type TaxRequest struct {
	RequestId            string `json:"request_id" validate:"required"`
	Amount               string `json:"amount" validate:"required"`
	PaymentRequestNumber string `json:"prn" validate:"required"`
	ShortCode            string `json:"short_code" validate:"required"`
	CallbackUrl          string `json:"callback_url" validate:"required"`
}
