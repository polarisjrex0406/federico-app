package seeds

import (
	"github.com/polarisjrex0406/federico-app/entities"
	"gorm.io/gorm"
)

func Balance(db *gorm.DB) error {
	users := []entities.User{}
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		balance := entities.Balance{
			UserID: user.ID,
			Amount: 0,
		}
		if err := db.Create(&balance).Error; err != nil {
			return err
		}
	}

	return nil
}
