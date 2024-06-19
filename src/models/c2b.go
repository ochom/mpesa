package models

import (
	"time"

	"github.com/ochom/gutils/helpers"
	"gorm.io/gorm"
)

// CustomerPayment store data when customer makes a payment to business
type CustomerPayment struct {
	Id                      int            `json:"id,omitempty" gorm:"primaryKey"`
	TransactionType         string         `json:"transaction_type,omitempty"`
	TransactionID           string         `json:"transaction_id,omitempty" gorm:"unique"`
	TransactionTime         string         `json:"transaction_time,omitempty"`
	PhoneNumber             string         `json:"phone_number,omitempty" gorm:"index"`
	TransactionAmount       string         `json:"transaction_amount,omitempty"`
	BillRefNumber           string         `json:"bill_ref_number,omitempty"`
	InvoiceNumber           string         `json:"invoice_number,omitempty"`
	ThirdPartyTransactionID string         `json:"third_party_transaction_id,omitempty"`
	CreatedAt               time.Time      `json:"created_at,omitempty"`
	UpdatedAt               time.Time      `json:"updated_at,omitempty"`
	DeletedAt               gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// NewCustomerPayment create a new CustomerPayment
func NewCustomerPayment(txId, txTime, txAmount, billRefNumber, invoiceNumber, msisdn string) *CustomerPayment {
	return &CustomerPayment{
		TransactionType:   "CustomerPayBillOnline",
		TransactionID:     txId,
		TransactionTime:   txTime,
		TransactionAmount: txAmount,
		BillRefNumber:     billRefNumber,
		InvoiceNumber:     invoiceNumber,
		PhoneNumber:       getPhoneNumber(billRefNumber, msisdn),
	}
}

func getPhoneNumber(billRef, msisdn string) string {
	if phone := helpers.ParseMobile(billRef); phone != "" {
		return phone
	}

	if phone := helpers.ParseMobile(msisdn); phone != "" {
		return phone
	}

	// TODO implement query phone number using hash
	return ""
}
