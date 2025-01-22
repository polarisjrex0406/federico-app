package repositories

import (
	"github.com/polarisjrex0406/federico-app/entities"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	BeginTx() *gorm.DB

	Create(tx *gorm.DB, transaction *entities.Transaction) error

	FindOneByTransactionID(transactionId string) (*entities.Transaction, error)
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{DB: db}
}

func (r *transactionRepository) BeginTx() *gorm.DB {
	return r.DB.Begin()
}

func (r *transactionRepository) Create(tx *gorm.DB, transaction *entities.Transaction) error {
	dbInst := r.DB
	if tx != nil {
		dbInst = tx
	}
	result := dbInst.Create(transaction)
	return result.Error
}

func (r *transactionRepository) FindOneByTransactionID(transactionId string) (*entities.Transaction, error) {
	transaction := &entities.Transaction{}
	result := r.DB.Where("transaction_id = ?", transactionId).First(transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}
