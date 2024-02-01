package walletinfogo_usecase

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

type walletInfoUsecase struct {
	walletService WalletService
}

func NewWalletInfoUsecase(ws WalletService) *walletInfoUsecase {
	return &walletInfoUsecase{ws}
}

func (uc *walletInfoUsecase) GetWalletInfo(ctx context.Context, ID int64) (entity.Wallet, error) {
	wallet, err := uc.walletService.GetWallet(ctx, ID)
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}
