CREATE TABLE IF NOT EXISTS tenant_payment_settings (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    default_provider_code TEXT NOT NULL DEFAULT 'manual_vietqr',
    updated_by_user_id UUID REFERENCES app_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_payment_settings_provider_check CHECK (
        default_provider_code IN ('manual_vietqr', 'sepay', 'payos')
    )
);

DROP TRIGGER IF EXISTS tenant_payment_settings_set_updated_at ON tenant_payment_settings;
CREATE TRIGGER tenant_payment_settings_set_updated_at
BEFORE UPDATE ON tenant_payment_settings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

INSERT INTO tenant_payment_settings (tenant_id, default_provider_code)
SELECT id, 'manual_vietqr'
FROM tenants
ON CONFLICT (tenant_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS tenant_email_cron_states (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    state JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_by_user_id UUID REFERENCES app_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_email_cron_states_state_object CHECK (jsonb_typeof(state) = 'object')
);

DROP TRIGGER IF EXISTS tenant_email_cron_states_set_updated_at ON tenant_email_cron_states;
CREATE TRIGGER tenant_email_cron_states_set_updated_at
BEFORE UPDATE ON tenant_email_cron_states
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS tenant_subscription_change_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    request_type TEXT NOT NULL,
    desired_plan_code TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    effective_at DATE,
    refund_amount INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'new',
    requested_by_user_id UUID REFERENCES app_users(id) ON DELETE SET NULL,
    processed_by_user_id UUID REFERENCES app_users(id) ON DELETE SET NULL,
    processed_at TIMESTAMPTZ,
    admin_note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_subscription_request_type_check CHECK (
        request_type IN ('upgrade', 'downgrade', 'cancel', 'refund')
    ),
    CONSTRAINT tenant_subscription_request_status_check CHECK (
        status IN ('new', 'approved', 'rejected', 'processed')
    ),
    CONSTRAINT tenant_subscription_request_refund_amount_check CHECK (refund_amount >= 0)
);

CREATE INDEX IF NOT EXISTS tenant_subscription_change_requests_tenant_status_idx
    ON tenant_subscription_change_requests(tenant_id, status, created_at DESC);

DROP TRIGGER IF EXISTS tenant_subscription_change_requests_set_updated_at ON tenant_subscription_change_requests;
CREATE TRIGGER tenant_subscription_change_requests_set_updated_at
BEFORE UPDATE ON tenant_subscription_change_requests
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
