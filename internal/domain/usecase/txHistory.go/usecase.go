package txhistory_usecase

import (
	"context"

	"github.com/The-Gleb/EWallet/internal/domain/entity"
)

type TransactionService interface {
	GetTransactions(ctx context.Context, walletID int64) ([]entity.Transaction, error)
	AddTransaction(ctx context.Context, tx entity.Transaction) (entity.Transaction, error)
}

type WalletService interface {
	GetWallet(ctx context.Context, ID int64) (entity.Wallet, error)
	CreateWallet(ctx context.Context) (entity.Wallet, error)
	TopUp(ctx context.Context, ID int64, amount float64) error
	Charge(ctx context.Context, ID int64, amount float64) error
}

type txHistoryUsecase struct {
	transactionService TransactionService
	walletService      WalletService
}

func NewTxHistoryUsecase(ts TransactionService, ws WalletService) *txHistoryUsecase {
	return &txHistoryUsecase{ts, ws}
}

func (uc txHistoryUsecase) GetTxHistory(ctx context.Context, walletID int64) ([]entity.Transaction, error) {

	_, err := uc.walletService.GetWallet(ctx, walletID)
	if err != nil {
		return []entity.Transaction{}, err
	}

	txs, err := uc.transactionService.GetTransactions(ctx, walletID)
	if err != nil {
		return []entity.Transaction{}, err
	}

	return txs, nil
}
