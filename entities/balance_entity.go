package entities

import "gorm.io/gorm"

type Balance struct {
	gorm.Model

	UserID uint `json:"user_id" gorm:"unique"`
	Amount float64 `json:"amount"`
}

// TableName overrides the default table name
func (Balance) TableName() string {
	return "tbl_balances"
}
