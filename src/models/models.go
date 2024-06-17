package models

// GetSchema get schema
func GetSchema() []any {
	return []any{
		&BusinessPayment{},
		&CustomerPayment{},
		&TaxPayment{},
	}
}
