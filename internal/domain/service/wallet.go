package service

import (
	"context"
	"errors"
	"log/slog"

	storage "github.com/The-Gleb/EWallet/internal/adapter/db/postgres"
	"github.com/The-Gleb/EWallet/internal/domain/entity"
)

var (
	ErrInsufficientFunds      = errors.New("not enough monet in the wallet")
	ErrWalletNotFound         = errors.New("wallet not found")
	ErrSenderWalletNotFound   = errors.New("sender wallet not found")
	ErrReceiverWalletNotFound = errors.New("receiver wallet not found")
)

type WalletStorage interface {
	GetWallet(ctx context.Context, ID int64) (entity.Wallet, error)
	CreateWallet(ctx context.Context, wallet entity.Wallet) (entity.Wallet, error)
	TopUp(ctx context.Context, ID int64, amount float64) error
	Charge(ctx context.Context, ID int64, amount float64) error
	MakeTransaction(ctx context.Context, tx entity.Transaction) error
}

type walletService struct {
	storage WalletStorage
}

func NewWalletService(s WalletStorage) *walletService {
	return &walletService{s}
}

func (ws *walletService) GetWallet(ctx context.Context, ID int64) (entity.Wallet, error) {
	wallet, err := ws.storage.GetWallet(ctx, ID)
	if err != nil {
		if errors.Is(err, storage.ErrNoDataFound) {
			return entity.Wallet{}, ErrWalletNotFound
		}
		return entity.Wallet{}, err
	}
	return wallet, nil
}

func (ws *walletService) CreateWallet(ctx context.Context) (entity.Wallet, error) {
	newWallet := entity.Wallet{
		Balance: 100,
	}
	return ws.storage.CreateWallet(ctx, newWallet)
}

func (ws *walletService) TopUp(ctx context.Context, walletID int64, amount float64) error {
	return ws.storage.TopUp(ctx, walletID, amount)
}

func (ws *walletService) Charge(ctx context.Context, walletID int64, amount float64) error {
	w, err := ws.GetWallet(ctx, walletID)
	if err != nil {
		return err
	}
	if w.Balance < amount {
		return ErrInsufficientFunds
	}
	return ws.storage.Charge(ctx, walletID, amount)
}

func (ws *walletService) MakeTransaction(ctx context.Context, tx entity.Transaction) error {

	wallet, err := ws.storage.GetWallet(ctx, tx.Sender)
	if err != nil {
		if errors.Is(err, storage.ErrNoDataFound) {
			return ErrSenderWalletNotFound
		}
		return err
	}

	if wallet.Balance < tx.Amount {
		slog.Error("insufficient funds", "error", ErrInsufficientFunds)
		return ErrInsufficientFunds
	}

	wallet, err = ws.storage.GetWallet(ctx, tx.Receiver)
	if err != nil {
		if errors.Is(err, storage.ErrNoDataFound) {
			return ErrReceiverWalletNotFound
		}
		slog.Error(err.Error())
		return err
	}

	return ws.storage.MakeTransaction(ctx, tx)
}
