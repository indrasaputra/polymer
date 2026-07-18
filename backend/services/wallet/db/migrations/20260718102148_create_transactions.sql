-- Create enum type "transaction_type"
CREATE TYPE public.transaction_type AS ENUM ('topup');
-- Create enum type "transaction_status"
CREATE TYPE public.transaction_status AS ENUM ('pending', 'completed', 'failed', 'cancelled');
-- Create "transactions" table
CREATE TABLE public.transactions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    type public.transaction_type NOT NULL,
    status public.transaction_status NOT NULL DEFAULT 'pending',
    idempotency_key uuid NOT NULL,
    amount numeric(20, 5) NOT NULL,
    currency character varying(3) NOT NULL,
    payment_session_id character varying(255) NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    deleted_by uuid NULL,
    PRIMARY KEY (id),
    CONSTRAINT transactions_idempotency_key_key UNIQUE (idempotency_key),
    CONSTRAINT transactions_payment_session_id_key UNIQUE (payment_session_id)
);
-- Create index "idx_transactions_user_id_status" to table: "transactions"
CREATE INDEX idx_transactions_user_id_status ON public.transactions (user_id, status);
