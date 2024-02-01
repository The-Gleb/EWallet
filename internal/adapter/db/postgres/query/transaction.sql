-- name: AddTransaction :one
INSERT INTO transactions
(sender, receiver, amount, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTransactions :many
SELECT * FROM transactions
WHERE sender = $1 OR receiver = $1
ORDER BY created_at;