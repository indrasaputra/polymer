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
