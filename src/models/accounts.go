package models

import (
	"gorm.io/gorm"
)

type AccountType string

const (
	C2B AccountType = "c2b"
	B2C AccountType = "b2c"
	B2B AccountType = "b2b"
	TAX AccountType = "tax"
)

// Account ...
type Account struct {
	gorm.Model
	ShortCode         string      `json:"short_code" gorm:"index"`
	Type              AccountType `json:"type"`
	PassKey           string      `json:"pass_key"`
	ConsumerKey       string      `json:"consumer_key"`
	ConsumerSecrete   string      `json:"-"`
	ValidationUrl     string      `json:"validation_url"`
	ConfirmationUrl   string      `json:"confirmation_url"`
	InitiatorName     string      `json:"initiator_name"`
	InitiatorPassword string      `json:"-"`
	Certificate       string      `json:"-"`
}

// NewAccount ...
func NewAccount(accType AccountType, shortCode, passKey, consumerKey, consumerSecrete string) *Account {
	return &Account{
		ShortCode:         shortCode,
		PassKey:           passKey,
		ConsumerKey:       consumerKey,
		ConsumerSecrete:   consumerSecrete,
		Type:              accType,
		ValidationUrl:     "",
		ConfirmationUrl:   "",
		InitiatorName:     "",
		InitiatorPassword: "",
		Certificate:       "",
	}
}
