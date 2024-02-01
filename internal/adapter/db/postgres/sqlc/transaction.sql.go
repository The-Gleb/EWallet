// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transaction.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addTransaction = `-- name: AddTransaction :one
INSERT INTO transactions
(sender, receiver, amount, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, sender, receiver, amount, created_at
`

type AddTransactionParams struct {
	Sender    pgtype.Int8
	Receiver  pgtype.Int8
	Amount    pgtype.Numeric
	CreatedAt pgtype.Timestamp
}

func (q *Queries) AddTransaction(ctx context.Context, arg AddTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, addTransaction,
		arg.Sender,
		arg.Receiver,
		arg.Amount,
		arg.CreatedAt,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.Sender,
		&i.Receiver,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransactions = `-- name: GetTransactions :many
SELECT id, sender, receiver, amount, created_at FROM transactions
WHERE sender = $1 OR receiver = $1
ORDER BY created_at
`

func (q *Queries) GetTransactions(ctx context.Context, sender pgtype.Int8) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, getTransactions, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.Sender,
			&i.Receiver,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}