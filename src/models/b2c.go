package models

import (
	"time"

	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/uuid"
)

// Payment store mpesa b2c requests and respective responses
type Payment struct {
	Id                       string    `json:"id"`
	Amount                   string    `json:"amount"`
	PhoneNumber              string    `json:"phone_number"`
	RequestId                string    `json:"request_id"`
	CallbackUrl              string    `json:"callback_url"`
	ConversationID           string    `json:"conversation_id"`
	OriginatorConversationID string    `json:"originator_conversation_id"`
	TransactionID            string    `json:"transaction_id"`
	ResponseCode             string    `json:"response_code"`
	ResponseDescription      string    `json:"response_description"`
	ResultCode               int       `json:"result_code"`
	ResultDescription        string    `json:"result_description"`
	Meta                     MetaData  `json:"meta" gorm:"type:json"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
	DeletedAt                time.Time `json:"deleted_at"`
}

// NewPayment create a new Payment
func NewPayment(requestId, phone, amount, cbUrl string) *Payment {
	return &Payment{
		Id:          uuid.New(),
		RequestId:   requestId,
		Amount:      amount,
		PhoneNumber: helpers.ParseMobile(phone),
		CallbackUrl: cbUrl,
	}
}
