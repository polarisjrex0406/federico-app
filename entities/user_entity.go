package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UserID uint `json:"user_id" gorm:"unique"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "tbl_users"
}
