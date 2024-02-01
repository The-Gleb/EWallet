package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/The-Gleb/EWallet/internal/adapter/db/postgres/sqlc"
	"github.com/The-Gleb/EWallet/internal/domain/entity"
	postgresql "github.com/The-Gleb/EWallet/pkg/client/postgres"
	"github.com/jackc/pgx/v5/pgtype"
)

type transactionStorage struct {
	sqlc *sqlc.Queries
}

func NewTransactionStorage(client postgresql.Client) *transactionStorage {
	return &transactionStorage{
		sqlc: sqlc.New(client),
	}
}

func (s *transactionStorage) GetTransactions(ctx context.Context, walletID int64) ([]entity.Transaction, error) {

	sqlcTxs, err := s.sqlc.GetTransactions(ctx, pgtype.Int8{
		Int64: walletID,
		Valid: true,
	})
	if err != nil {
		slog.Error(err.Error())
		return make([]entity.Transaction, 0), err
	}

	txs := make([]entity.Transaction, len(sqlcTxs))
	for i, tx := range sqlcTxs {
		txs[i] = entity.Transaction{
			ID:       tx.ID,
			Sender:   tx.Sender.Int64,
			Receiver: tx.Receiver.Int64,
			Time:     tx.CreatedAt.Time,
		}
		float8, err := tx.Amount.Float64Value()
		if err != nil {
			slog.Error(err.Error())
			return make([]entity.Transaction, 0), err
		}
		txs[i].Amount = float8.Float64
	}

	return txs, nil

}
func (s *transactionStorage) AddTransaction(ctx context.Context, tx entity.Transaction) (entity.Transaction, error) {

	addTransactionParams := sqlc.AddTransactionParams{
		Sender: pgtype.Int8{
			Int64: tx.Sender,
			Valid: true,
		},
		Receiver: pgtype.Int8{
			Int64: tx.Receiver,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  tx.Time,
			Valid: true,
		},
	}
	err := addTransactionParams.Amount.Scan(fmt.Sprint(tx.Amount))
	if err != nil {
		slog.Error(err.Error())
		return entity.Transaction{}, err
	}

	sqlcTx, err := s.sqlc.AddTransaction(ctx, addTransactionParams)
	if err != nil {
		slog.Error(err.Error())
		return entity.Transaction{}, err
	}

	float8, err := sqlcTx.Amount.Float64Value()
	if err != nil {
		slog.Error(err.Error())
		return entity.Transaction{}, err
	}

	return entity.Transaction{
		ID:       sqlcTx.ID,
		Sender:   sqlcTx.Sender.Int64,
		Receiver: sqlcTx.Receiver.Int64,
		Amount:   float8.Float64,
		Time:     sqlcTx.CreatedAt.Time,
	}, err

}
