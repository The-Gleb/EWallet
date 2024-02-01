package newwallet_usecase

import (
	"context"

	"github.com/The-Gleb/EWallet/internal/domain/entity"
)

type WalletService interface {
	GetWallet(ctx context.Context, ID int64) (entity.Wallet, error)
	CreateWallet(ctx context.Context) (entity.Wallet, error)
	TopUp(ctx context.Context, ID int64, amount float64) error
	Charge(ctx context.Context, ID int64, amount float64) error
}

type newWalletUsecase struct {
	walletService WalletService
}

func NewNewWalletUsecase(ws WalletService) *newWalletUsecase {
	return &newWalletUsecase{ws}
}

func (nw *newWalletUsecase) NewWallet(ctx context.Context) (entity.Wallet, error) {
	newWallet, err := nw.walletService.CreateWallet(ctx)
	if err != nil {
		return entity.Wallet{}, err
	}
	return newWallet, nil
}
