package models

import (
	"github.com/ochom/gutils/uuid"
	"gorm.io/gorm"
)

// TaxPayment store data when business makes a tax payment
type TaxPayment struct {
	gorm.Model
	Uuid                     string   `json:"uuid"`
	Amount                   string   `json:"amount"`
	ShortCode                string   `json:"short_code"`
	PaymentRequestNumber     string   `json:"payment_request_number"` // prn
	RequestId                string   `json:"request_id"`
	CallbackUrl              string   `json:"callback_url"`
	ConversationID           string   `json:"conversation_id"`
	OriginatorConversationID string   `json:"originator_conversation_id"`
	TransactionID            string   `json:"transaction_id"`
	ResponseCode             string   `json:"response_code"`
	ResponseDescription      string   `json:"response_description"`
	ResultCode               int      `json:"result_code"`
	ResultDescription        string   `json:"result_description"`
	Meta                     MetaData `json:"meta" gorm:"type:json"`
}

// NewTaxPayment create a new TaxPayment
func NewTaxPayment(requestId, shortCode, prn, amount, cbUrl string) *TaxPayment {
	return &TaxPayment{
		Uuid:                 uuid.New(),
		RequestId:            requestId,
		ShortCode:            shortCode,
		Amount:               amount,
		PaymentRequestNumber: prn,
		CallbackUrl:          cbUrl,
	}
}
