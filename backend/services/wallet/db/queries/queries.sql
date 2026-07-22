-- name: InsertWallet :one
INSERT INTO wallets (id, user_id, balance, currency, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id, currency) WHERE deleted_at IS NULL DO NOTHING
RETURNING *;

-- name: GetActiveWalletByIDForUpdate :one
SELECT * FROM wallets WHERE id = $1 AND deleted_at IS NULL LIMIT 1 FOR NO KEY UPDATE; --noqa

-- name: GetUserActiveWalletByUserIdAndCurrency :one
SELECT * FROM wallets
WHERE user_id = $1 AND currency = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserWalletByIDAndUserID :one
SELECT * FROM wallets
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: AddWalletBalance :one
UPDATE wallets SET balance = balance + @amount WHERE id = $1 --noqa
RETURNING *;

-- name: InsertCustomer :one
INSERT INTO customers (id, user_id, stripe_customer_id, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id) WHERE deleted_at IS NULL DO NOTHING
RETURNING *;

-- name: GetCustomerByUserID :one
SELECT * FROM customers
WHERE user_id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: InsertTransaction :one
INSERT INTO transactions (id, user_id, type, status, idempotency_key, amount, currency, checkout_session_id, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetPendingTransactionByIdempotencyKey :one
SELECT * FROM transactions
WHERE idempotency_key = $1 AND status = 'pending' AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateTransactionToCompletedByCheckoutSessionID :one
UPDATE transactions
SET status = 'completed', updated_at = $1, updated_by = $2
WHERE checkout_session_id = $3 AND deleted_at IS NULL
RETURNING *;

-- name: GetActiveTransactionByCheckoutSessionIDForUpdate :one
SELECT * FROM transactions
WHERE checkout_session_id = $1 AND deleted_at IS NULL
LIMIT 1 FOR NO KEY UPDATE; --noqa
