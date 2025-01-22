package migrations

import (
	"fmt"

	"github.com/polarisjrex0406/federico-app/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.User(db); err != nil {
		fmt.Println("Failed to seed User with error: ", err.Error())
	}
	if err := seeds.Balance(db); err != nil {
		fmt.Println("Failed to seed Balance with error: ", err.Error())
	}
	return nil
}
