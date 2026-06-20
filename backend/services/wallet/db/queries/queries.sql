-- name: InsertWallet :one
INSERT INTO wallets (id, user_id, balance, currency, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id, currency) WHERE deleted_at IS NULL DO NOTHING
RETURNING *;

-- name: GetUserActiveWalletByUserIdAndCurrency :one
SELECT * FROM wallets
WHERE user_id = $1 AND currency = $2 AND deleted_at IS NULL
LIMIT 1;
