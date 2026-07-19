CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    balance DECIMAL(20, 5) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID,

    CONSTRAINT non_negative_balance CHECK (balance >= 0)
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wallets_id_user_id
ON wallets USING btree (id, user_id);

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_wallets_user_id_currency
ON wallets USING btree (user_id, currency)
WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    stripe_customer_id VARCHAR(20) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID
);

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_customers_user_id
ON customers USING btree (user_id)
WHERE deleted_at IS NULL;

CREATE TYPE TRANSACTION_TYPE AS ENUM ('topup');
CREATE TYPE TRANSACTION_STATUS AS ENUM ('pending', 'completed', 'failed', 'cancelled');

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type TRANSACTION_TYPE NOT NULL,
    status TRANSACTION_STATUS NOT NULL DEFAULT 'pending',

    idempotency_key UUID UNIQUE NOT NULL,

    amount DECIMAL(20, 5) NOT NULL,
    currency VARCHAR(3) NOT NULL,

    checkout_session_id VARCHAR(255) UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_transactions_user_id_status
ON transactions USING btree (user_id, status);
