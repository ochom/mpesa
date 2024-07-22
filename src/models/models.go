package models

import (
	"time"

	"gorm.io/gorm"
)

// GetSchema get schema
func GetSchema() []any {
	return []any{
		&BusinessPayment{},
		&CustomerPayment{},
		&TaxPayment{},
		&Account{},
	}
}

type Model struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
