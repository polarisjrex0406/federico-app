package services

import (
	"github.com/polarisjrex0406/federico-app/dto"
	"github.com/polarisjrex0406/federico-app/entities"
	"github.com/polarisjrex0406/federico-app/repositories"
	"github.com/polarisjrex0406/federico-app/utils"
	"gorm.io/gorm"
)

type (
	UserService interface {
		GetBalance(userId uint) (*dto.UserGetBalanceResponse, error)
		DoTransaction(userId uint, req dto.UserDoTransactionRequest) error
	}
	userService struct {
		balanceRepo     repositories.BalanceRepository
		transactionRepo repositories.TransactionRepository
	}
)

func NewUserService(
	balanceRepo repositories.BalanceRepository,
	transactionRepo repositories.TransactionRepository,
) UserService {
	return &userService{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *userService) GetBalance(userId uint) (*dto.UserGetBalanceResponse, error) {
	existingBalance, err := s.balanceRepo.FindOneByUserID(userId)
	if err != nil {
		return nil, dto.ErrFindOneBalanceByUserID
	}

	res := dto.UserGetBalanceResponse{
		UserID:  existingBalance.UserID,
		Balance: utils.Float64ToString(existingBalance.Amount),
	}
	return &res, nil
}

func (s *userService) DoTransaction(userId uint, req dto.UserDoTransactionRequest) error {
	_, err := s.transactionRepo.FindOneByTransactionID(req.TransactionID)
	if err == nil {
		return dto.ErrTransactionAlreadyExists
	}
	if err != gorm.ErrRecordNotFound {
		return dto.ErrFindOneTransactionByTransactionID
	}

	existingBalance, err := s.balanceRepo.FindOneByUserID(userId)
	if err != nil {
		return dto.ErrFindOneBalanceByUserID
	}

	amount, err := utils.StringToFloat64(req.Amount)
	if err != nil {
		return dto.ErrConvertStringToFloat64
	}
	if req.State == "loss" && existingBalance.Amount < amount {
		return dto.ErrNotEnoughBalance
	}

	// Begin transaction
	tx := s.balanceRepo.BeginTx()
	if tx.Error != nil {
		return dto.ErrBeginDBTx
	}
	// Rollback when panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update balance
	if req.State == "loss" {
		amount *= -1
	}
	if _, err = s.balanceRepo.UpdateOneByUserID(tx, userId, amount); err != nil {
		tx.Rollback()
		return dto.ErrUpdateOneBalanceByUserID
	}

	// Create transaction
	newTransaction := entities.Transaction{
		UserID:        userId,
		State:         req.State,
		Amount:        amount,
		TransactionID: req.TransactionID,
	}
	if err = s.transactionRepo.Create(tx, &newTransaction); err != nil {
		tx.Rollback()
		return dto.ErrCreateTransaction
	}

	// Commit transaction
	if tx.Commit().Error != nil {
		return dto.ErrCommitDBTx
	}

	return nil
}
