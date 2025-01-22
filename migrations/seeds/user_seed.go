package seeds

import (
	"github.com/polarisjrex0406/federico-app/entities"
	"gorm.io/gorm"
)

func User(db *gorm.DB) error {
	userIdSlice := []uint{1, 2, 3}

	for _, userId := range userIdSlice {
		user := entities.User{
			UserID: userId,
		}
		if err := db.Where("user_id = ?", userId).Attrs(user).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
