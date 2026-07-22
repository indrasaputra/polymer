-- Drop index "idx_transactions_user_id_status" from table: "transactions"
DROP INDEX public.idx_transactions_user_id_status;
-- Create index "idx_transactions_checkout_session_id" to table: "transactions"
CREATE INDEX idx_transactions_checkout_session_id ON public.transactions (checkout_session_id) WHERE ((checkout_session_id IS NOT NULL) AND (deleted_at IS NULL));
-- Drop index "idx_wallets_id_user_id" from table: "wallets"
DROP INDEX public.idx_wallets_id_user_id;
