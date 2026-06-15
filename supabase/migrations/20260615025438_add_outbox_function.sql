-- 1. Ensure the necessary outbound HTTP network queue extension is active
CREATE EXTENSION IF NOT EXISTS pg_net;
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- 2. Create the Dedicated Transactional Outbox Table
CREATE TABLE IF NOT EXISTS public.outbox_events (
    id BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PROCESSING', 'PROCESSED', 'FAILED')),
    net_request_id BIGINT NULL, -- Captures pg_net's internal tracking identifier
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2a. Indexes to support the cron worker's polling patterns and avoid full-table scans
--     (status, retry_count) covers the PENDING dispatch query in Step A
--     (status, net_request_id) covers the PROCESSING reconcile query in Step B
--     (status, updated_at) covers the stale-PROCESSING timeout check in Step B
CREATE INDEX IF NOT EXISTS idx_outbox_status_retry
    ON public.outbox_events (status, retry_count);

CREATE INDEX IF NOT EXISTS idx_outbox_status_net_request
    ON public.outbox_events (status, net_request_id);

CREATE INDEX IF NOT EXISTS idx_outbox_status_updated_at
    ON public.outbox_events (status, updated_at);

-- 3. The Signup Trigger (Fast, Local, and 100% Failure Proof)
CREATE OR REPLACE FUNCTION public.handle_auth_signup_outbox()
RETURNS TRIGGER AS $$
BEGIN
    -- Simply append to your local table. This will NEVER fail due to network drops.
    INSERT INTO public.outbox_events (event_type, payload)
    VALUES (
        'USER.SIGNUP',
        jsonb_build_object(
            'id', NEW.id,
            'email', NEW.email,
            'created_at', NEW.created_at
        )
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- 4. Bind trigger to auth.users
DROP TRIGGER IF EXISTS trigger_on_auth_user_signup ON auth.users;
CREATE TRIGGER trigger_on_auth_user_signup
    AFTER INSERT ON auth.users
    FOR EACH ROW EXECUTE FUNCTION public.handle_auth_signup_outbox();

-- 5. Create the Outbox Dispatcher & Reconciler
CREATE OR REPLACE FUNCTION public.dispatch_and_reconcile_outbox()
RETURNS void AS $$
DECLARE
    -- Single source of truth for the retry budget.
    -- Total attempts = k_max_retries (retries 0 through k_max_retries - 1).
    -- To change the budget, update only this constant.
    k_max_retries CONSTANT INT := 5;

    user_service_signup_webhook_url TEXT;
    user_service_webhook_secret TEXT;
    event_record RECORD;
    response_record RECORD;
    new_net_id BIGINT;
    table_exists BOOLEAN;
BEGIN
    -- Load active cluster configuration keys
    SELECT decrypted_secret INTO user_service_signup_webhook_url
    FROM vault.decrypted_secrets WHERE name = 'user_service_signup_webhook_url';

    SELECT decrypted_secret INTO user_service_webhook_secret
    FROM vault.decrypted_secrets WHERE name = 'user_service_webhook_secret';

    -- Short-circuit if either config value is missing.
    -- Without this guard, every PENDING row would call net.http_post with a NULL URL
    -- on each cron tick, burning through the retry budget with no chance of success.
    IF user_service_signup_webhook_url IS NULL OR user_service_webhook_secret IS NULL THEN
        RAISE WARNING 'dispatch_and_reconcile_outbox: missing webhook config (url=%, secret set=%). Skipping run.',
            user_service_signup_webhook_url,
            (user_service_webhook_secret IS NOT NULL);
        RETURN;
    END IF;

    -- Step A: Dispatch new pending events
    -- FOR UPDATE SKIP LOCKED prevents double-dispatch if pg_cron fires while a
    -- previous invocation is still running (concurrency safety).
    FOR event_record IN
        SELECT id, payload FROM public.outbox_events
        WHERE status = 'PENDING' AND retry_count < k_max_retries
        FOR UPDATE SKIP LOCKED
    LOOP
        BEGIN
            -- Pass payload to pg_net queue and capture its query token ID
            SELECT net.http_post(
                url := user_service_signup_webhook_url,
                body := event_record.payload,
                headers := jsonb_build_object(
                    'Content-Type', 'application/json',
                    'X-Webhook-Secret', user_service_webhook_secret
                ),
                timeout_milliseconds := 2000
            ) INTO new_net_id;

            -- Move status to 'PROCESSING' so we wait for the network daemon feedback loop
            UPDATE public.outbox_events
            SET status = 'PROCESSING', net_request_id = new_net_id, updated_at = NOW()
            WHERE id = event_record.id;
        EXCEPTION WHEN OTHERS THEN
            -- If pg_net crashes internally, log a retry instantly
            UPDATE public.outbox_events
            SET retry_count = retry_count + 1, updated_at = NOW()
            WHERE id = event_record.id;
        END;
    END LOOP;

    -- Step B: Reconcile in-flight events against real network results
    SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'net' AND table_name = '_http_response')
    INTO table_exists;
    IF NOT table_exists THEN RETURN; END IF;

    -- Added FOR UPDATE SKIP LOCKED here too for the same concurrency reason
    FOR event_record IN
        SELECT id, net_request_id, retry_count, updated_at FROM public.outbox_events
        WHERE status = 'PROCESSING'
        FOR UPDATE SKIP LOCKED
    LOOP
        -- Look inside pg_net's logging tables for our matching request ID token
        SELECT status_code, error_msg INTO response_record
        FROM net._http_response WHERE id = event_record.net_request_id LIMIT 1;

        IF FOUND THEN
            -- Explicitly guard against NULL status_code (e.g. pg_net timeout),
            -- treating it as a failure rather than relying on NULL comparison behaviour.
            IF response_record.status_code IS NOT NULL
               AND response_record.status_code >= 200
               AND response_record.status_code < 300 THEN
                -- Success: user service responded with a 2xx status
                UPDATE public.outbox_events
                SET status = 'PROCESSED', updated_at = NOW()
                WHERE id = event_record.id;
            ELSE
                -- Failure: 4xx/5xx, NULL (timeout), or internal pg_net error
                IF event_record.retry_count >= k_max_retries - 1 THEN
                    -- Dead Letter Queue: hard-fail after k_max_retries total attempts
                    UPDATE public.outbox_events
                    SET status = 'FAILED', updated_at = NOW()
                    WHERE id = event_record.id;
                ELSE
                    -- Micro-retry: return to PENDING so the next cron cycle retries
                    UPDATE public.outbox_events
                    SET status = 'PENDING', retry_count = retry_count + 1, net_request_id = NULL, updated_at = NOW()
                    WHERE id = event_record.id;
                END IF;
            END IF;

            -- Clean up the processed row from pg_net's system log to save disk space
            DELETE FROM net._http_response WHERE id = event_record.net_request_id;

        ELSIF event_record.updated_at < NOW() - INTERVAL '90 seconds' THEN
            UPDATE public.outbox_events
            SET status = 'PENDING', retry_count = retry_count + 1, net_request_id = NULL, updated_at = NOW()
            WHERE id = event_record.id;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;


-- 6. Configure pg_cron to fire the outbox dispatcher once per minute
DO $$
BEGIN
    PERFORM cron.unschedule('dispatch-outbox-webhooks-job');
EXCEPTION WHEN OTHERS THEN
    -- Absorb error gracefully if the job wasn't scheduled yet
END $$;

DO $$
BEGIN
    PERFORM cron.schedule(
        'dispatch-outbox-webhooks-job',
        '* * * * *', -- Runs exactly once per minute
        'SELECT public.dispatch_and_reconcile_outbox();'
    );
END $$;
