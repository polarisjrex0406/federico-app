package repositories

import (
	"github.com/polarisjrex0406/federico-app/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceRepository interface {
	BeginTx() *gorm.DB

	FindOneByUserID(userId uint) (*entities.Balance, error)

	UpdateOneByUserID(tx *gorm.DB, userId uint, amount float64) (*entities.Balance, error)
}

type balanceRepository struct {
	DB *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{DB: db}
}

func (r *balanceRepository) BeginTx() *gorm.DB {
	return r.DB.Begin()
}

func (r *balanceRepository) FindOneByUserID(userId uint) (*entities.Balance, error) {
	balance := entities.Balance{}
	result := r.DB.Where("user_id = ?", userId).First(&balance)
	if result.Error != nil {
		return nil, result.Error
	}
	return &balance, nil
}

func (r *balanceRepository) UpdateOneByUserID(tx *gorm.DB, userId uint, amount float64) (*entities.Balance, error) {
	dbInst := r.DB
	if tx != nil {
		dbInst = tx
	}

	balance := entities.Balance{}
	result := dbInst.Model(&balance).
		Clauses(clause.Returning{}).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"amount": gorm.Expr("amount + ?", amount),
		})
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &balance, nil
}
