-- name: GetWalletInfo :one
SELECT * FROM wallets
WHERE id = $1;

-- name: CreateWallet :one
INSERT INTO wallets
( balance)
VALUES ($1)
RETURNING *;

-- name: Topup :exec
UPDATE wallets
SET balance = balance + $2
WHERE id = $1;

-- name: Charge :exec
UPDATE wallets
SET balance = balance - $2
WHERE id = $1;