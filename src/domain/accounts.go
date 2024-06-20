package domain

import "github.com/ochom/mpesa/src/models"

// CreateAccountRequest the payload required to initiate an mpesa stk push
type CreateAccountRequest struct {
	ShortCode         string             `json:"short_code" validate:"required"`
	Type              models.AccountType `json:"type" validate:"required"`
	PassKey           string             `json:"pass_key"`
	ConsumerKey       string             `json:"consumer_key"`
	ConsumerSecrete   string             `json:"consumer_secrete"`
	ValidationUrl     string             `json:"validation_url"`
	ConfirmationUrl   string             `json:"confirmation_url"`
	InitiatorName     string             `json:"initiator_name"`
	InitiatorPassword string             `json:"initiator_password"`
	Certificate       string             `json:"certificate"`
}
