package models

import (
	"time"
)

// MpesaExpress store mpesa stk-push requests and respective responses
type MpesaExpress struct {
	Id                  int       `json:"id"`
	Amount              string    `json:"amount"`
	PhoneNumber         string    `json:"phone_number"`
	AccountReference    string    `json:"account_reference"`
	CallbackUrl         string    `json:"callback_url"`
	MerchantRequestId   string    `json:"merchant_request_id"`
	CheckoutRequestId   string    `json:"checkout_request_id"`
	ResponseCode        string    `json:"response_code"`
	ResponseDescription string    `json:"response_description"`
	ResultCode          int       `json:"result_code"`
	ResultDescription   string    `json:"result_description"`
	Meta                MetaData  `json:"meta" gorm:"type:json"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           time.Time `json:"deleted_at"`
}

// NewMpesaExpress create a new MpesaExpress
func NewMpesaExpress(phone, amount, cbUrl, AccountReference string) *MpesaExpress {
	return &MpesaExpress{
		Amount:           amount,
		PhoneNumber:      phone,
		AccountReference: AccountReference,
		CallbackUrl:      cbUrl,
	}
}
