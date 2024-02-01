package maketransaction_usecase

import (
	"context"
	"time"

	"github.com/The-Gleb/EWallet/internal/domain/entity"
)

type WalletService interface {
	GetWallet(ctx context.Context, ID int64) (entity.Wallet, error)
	CreateWallet(ctx context.Context) (entity.Wallet, error)
	TopUp(ctx context.Context, ID int64, amount float64) error
	Charge(ctx context.Context, ID int64, amount float64) error
	MakeTransaction(ctx context.Context, tx entity.Transaction) error
}

type TransactionService interface {
	GetTransactions(ctx context.Context, walletID int64) ([]entity.Transaction, error)
	AddTransaction(ctx context.Context, tx entity.Transaction) (entity.Transaction, error)
}

type makeTransactionUsecase struct {
	walletService      WalletService
	transactionService TransactionService
}

func NewMakeTransactionUsecase(ws WalletService, ts TransactionService) *makeTransactionUsecase {
	return &makeTransactionUsecase{ws, ts}
}

func (uc *makeTransactionUsecase) MakeTransaction(ctx context.Context, dto MakeTransactionDTO) error {

	newTransaction := entity.Transaction{
		Sender:   dto.From,
		Receiver: dto.To,
		Amount:   dto.Amount,
		Time:     time.Now(),
	}

	// err := uc.walletService.Charge(ctx, newTransaction.Sender, newTransaction.Amount)
	// if err != nil {
	// 	return err
	// }

	// err = uc.walletService.TopUp(ctx, newTransaction.Receiver, newTransaction.Amount)
	// if err != nil {

	// 	return err
	// }

	err := uc.walletService.MakeTransaction(ctx, newTransaction)
	if err != nil {
		return err
	}

	_, err = uc.transactionService.AddTransaction(ctx, newTransaction)
	if err != nil {
		return err
	}

	return nil

}
