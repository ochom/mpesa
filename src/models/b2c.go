package models

import (
	"github.com/ochom/gutils/helpers"
)

// BusinessPayment store data when a business makes payment to Customer
type BusinessPayment struct {
	Model
	AccountID                int      `json:"account_id"`
	Amount                   string   `json:"amount"`
	PhoneNumber              string   `json:"phone_number"`
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

// NewBusinessPayment create a new BusinessPayment
func NewBusinessPayment(shortCodeID int, requestId, phone, amount, cbUrl string) *BusinessPayment {
	return &BusinessPayment{
		AccountID:   shortCodeID,
		RequestId:   requestId,
		Amount:      amount,
		PhoneNumber: helpers.ParseMobile(phone),
		CallbackUrl: cbUrl,
	}
}
