// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: wallet.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const charge = `-- name: Charge :exec
UPDATE wallets
SET balance = balance - $2
WHERE id = $1
`

type ChargeParams struct {
	ID      int64
	Balance pgtype.Numeric
}

func (q *Queries) Charge(ctx context.Context, arg ChargeParams) error {
	_, err := q.db.Exec(ctx, charge, arg.ID, arg.Balance)
	return err
}

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets
( balance)
VALUES ($1)
RETURNING id, balance
`

func (q *Queries) CreateWallet(ctx context.Context, balance pgtype.Numeric) (Wallet, error) {
	row := q.db.QueryRow(ctx, createWallet, balance)
	var i Wallet
	err := row.Scan(&i.ID, &i.Balance)
	return i, err
}

const getWalletInfo = `-- name: GetWalletInfo :one
SELECT id, balance FROM wallets
WHERE id = $1
`

func (q *Queries) GetWalletInfo(ctx context.Context, id int64) (Wallet, error) {
	row := q.db.QueryRow(ctx, getWalletInfo, id)
	var i Wallet
	err := row.Scan(&i.ID, &i.Balance)
	return i, err
}

const topup = `-- name: Topup :exec
UPDATE wallets
SET balance = balance + $2
WHERE id = $1
`

type TopupParams struct {
	ID      int64
	Balance pgtype.Numeric
}

func (q *Queries) Topup(ctx context.Context, arg TopupParams) error {
	_, err := q.db.Exec(ctx, topup, arg.ID, arg.Balance)
	return err
}
