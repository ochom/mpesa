package models

import (
	"database/sql/driver"
	"time"

	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/uuid"
)

// Payment store mpesa b2c requests and respective responses
type Payment struct {
	Id                       string      `json:"id"`
	Amount                   string      `json:"amount"`
	PhoneNumber              string      `json:"phone_number"`
	RequestId                string      `json:"request_id"`
	CallbackUrl              string      `json:"callback_url"`
	ConversationID           string      `json:"conversation_id"`
	OriginatorConversationID string      `json:"originator_conversation_id"`
	TransactionID            string      `json:"transaction_id"`
	ResponseCode             string      `json:"response_code"`
	ResponseDescription      string      `json:"response_description"`
	ResultCode               int         `json:"result_code"`
	ResultDescription        string      `json:"result_description"`
	Meta                     B2cMetaData `json:"meta" gorm:"type:json"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
	DeletedAt                time.Time   `json:"deleted_at"`
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

// B2cMetaData the metadata for an mpesa express request
type B2cMetaData map[string]any

// Scan implements the sql.Scanner interface.
func (m *B2cMetaData) Scan(value any) error {
	v, ok := value.([]byte)
	if !ok {
		return nil
	}

	*m = helpers.FromBytes[map[string]any](v)
	return nil
}

// Value implements the driver.Valuer interface.
func (m B2cMetaData) Value() (driver.Value, error) {
	return helpers.ToBytes(m), nil
}
