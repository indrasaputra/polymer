-- Create "wallets" table
CREATE TABLE public.wallets (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    balance numeric(20, 5) NOT NULL DEFAULT 0,
    currency character varying(3) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    deleted_by uuid NULL,
    PRIMARY KEY (id),
    CONSTRAINT non_negative_balance CHECK (balance >= (0)::numeric)
);
-- Create index "idx_wallets_id_user_id" to table: "wallets"
CREATE INDEX idx_wallets_id_user_id ON public.wallets (id, user_id);
-- Create index "idx_wallets_user_id_currency" to table: "wallets"
CREATE UNIQUE INDEX idx_wallets_user_id_currency ON public.wallets (user_id, currency) WHERE (deleted_at IS NULL);
