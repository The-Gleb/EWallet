package service

import (
	"context"

	"github.com/The-Gleb/EWallet/internal/domain/entity"
)

type transactionStorage interface {
	GetTransactions(ctx context.Context, walletID int64) ([]entity.Transaction, error)
	AddTransaction(ctx context.Context, tx entity.Transaction) (entity.Transaction, error)
}

type transactionService struct {
	storage transactionStorage
}

func NewTransactionService(s transactionStorage) *transactionService {
	return &transactionService{s}
}

func (s *transactionService) GetTransactions(ctx context.Context, walletID int64) ([]entity.Transaction, error) {
	return s.storage.GetTransactions(ctx, walletID)
}

func (s *transactionService) AddTransaction(ctx context.Context, tx entity.Transaction) (entity.Transaction, error) {
	return s.storage.AddTransaction(ctx, tx)
}
