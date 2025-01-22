package entities

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model

	UserID        uint    `json:"user_id"`
	TransactionID string  `json:"transaction_id"`
	State         string  `json:"state"`
	Amount        float64 `json:"amount"`
}

// TableName overrides the default table name
func (Transaction) TableName() string {
	return "tbl_transactions"
}
