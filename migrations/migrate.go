package migrations

import (
	"github.com/polarisjrex0406/federico-app/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	tables := []interface{}{
		&entities.Balance{},
		&entities.Transaction{},
		&entities.User{},
	}

	tableList, err := db.Migrator().GetTables()
	if err != nil {
		return err
	}

	for i := range tableList {
		db.Migrator().DropTable(tableList[i])
	}

	for _, table := range tables {
		if err := db.Migrator().CreateTable(table); err != nil {
			return err
		}
	}

	return nil
}
