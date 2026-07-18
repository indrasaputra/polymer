-- Create "customers" table
CREATE TABLE public.customers (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    stripe_customer_id character varying(20) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    deleted_by uuid NULL,
    PRIMARY KEY (id)
);
-- Create index "idx_customers_user_id" to table: "customers"
CREATE UNIQUE INDEX idx_customers_user_id ON public.customers (user_id) WHERE (deleted_at IS NULL);
