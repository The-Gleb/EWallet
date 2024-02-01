package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/The-Gleb/EWallet/internal/adapter/db/postgres/sqlc"
	"github.com/The-Gleb/EWallet/internal/domain/entity"
	postgresql "github.com/The-Gleb/EWallet/pkg/client/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNoDataFound = errors.New("no data found")
)

type walletStorage struct {
	client postgresql.Client
	sqlc   *sqlc.Queries
}

func NewWalletStorage(client postgresql.Client) *walletStorage {
	return &walletStorage{
		client: client,
		sqlc:   sqlc.New(client),
	}
}

func (s *walletStorage) GetWallet(ctx context.Context, ID int64) (entity.Wallet, error) {

	sqlcWallet, err := s.sqlc.GetWalletInfo(ctx, ID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Wallet{}, ErrNoDataFound
		}
		return entity.Wallet{}, err
	}
	w := entity.Wallet{
		ID: sqlcWallet.ID,
	}
	float8, err := sqlcWallet.Balance.Float64Value()
	if err != nil {
		slog.Error(err.Error())
		return entity.Wallet{}, err
	}
	w.Balance = float8.Float64

	return w, nil

}
func (s *walletStorage) CreateWallet(ctx context.Context, wallet entity.Wallet) (entity.Wallet, error) {

	var balance pgtype.Numeric
	err := balance.Scan(fmt.Sprint(wallet.Balance))
	if err != nil {
		slog.Error(err.Error())
		return entity.Wallet{}, err
	}

	sqlcWallet, err := s.sqlc.CreateWallet(ctx, balance)
	if err != nil {
		slog.Error(err.Error())
		return entity.Wallet{}, err
	}

	return entity.Wallet{ID: sqlcWallet.ID, Balance: wallet.Balance}, nil

}
func (s *walletStorage) TopUp(ctx context.Context, ID int64, amount float64) error {

	topupParams := sqlc.TopupParams{
		ID: ID,
	}
	err := topupParams.Balance.Scan(fmt.Sprint(amount))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = s.sqlc.Topup(ctx, topupParams)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil

}
func (s *walletStorage) Charge(ctx context.Context, ID int64, amount float64) error {

	chargeParams := sqlc.ChargeParams{
		ID: ID,
	}
	err := chargeParams.Balance.Scan(fmt.Sprint(amount))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = s.sqlc.Charge(ctx, chargeParams)
	if err != nil {
		slog.Error(err.Error())
		return err // TODO
	}

	return nil

}

func (s *walletStorage) MakeTransaction(ctx context.Context, tx entity.Transaction) error {
	dbTx, err := s.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	sqlcTx := s.sqlc.WithTx(dbTx)

	chargeParams := sqlc.ChargeParams{
		ID: tx.Sender,
	}

	err = chargeParams.Balance.Scan(fmt.Sprint(tx.Amount))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = sqlcTx.Charge(ctx, chargeParams)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	topupParams := sqlc.TopupParams{
		ID: tx.Receiver,
	}
	err = topupParams.Balance.Scan(fmt.Sprint(tx.Amount))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = sqlcTx.Topup(ctx, topupParams)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return dbTx.Commit(ctx)

}
