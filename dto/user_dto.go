package dto

import (
	"errors"
)

var (
	ErrFindOneBalanceByUserID            = errors.New("failed find one balance by user_id")
	ErrTransactionAlreadyExists          = errors.New("transaction already exists")
	ErrFindOneTransactionByTransactionID = errors.New("failed find one transaction by transaction_id")
	ErrConvertStringToFloat64            = errors.New("failed convert string to float64")
	ErrNotEnoughBalance                  = errors.New("not enough balance")
	ErrBeginDBTx                         = errors.New("failed begin db transaction")
	ErrUpdateOneBalanceByUserID          = errors.New("failed update one balance by user_id")
	ErrCreateTransaction                 = errors.New("failed create transaction")
	ErrCommitDBTx                        = errors.New("failed commit db transaction")
)

type (
	UserGetBalanceResponse struct {
		UserID  uint   `json:"userId" bind:"required"`
		Balance string `json:"balance" bind:"required"`
	}

	UserDoTransactionRequest struct {
		State         string `json:"state" bind:"required"`
		Amount        string `json:"amount" bind:"required"`
		TransactionID string `json:"transactionId" bind:"required"`
	}
)
