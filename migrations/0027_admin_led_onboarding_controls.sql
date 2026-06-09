CREATE TABLE IF NOT EXISTS tenant_intake_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_name TEXT NOT NULL,
    contact_name TEXT NOT NULL DEFAULT '',
    contact_email TEXT NOT NULL DEFAULT '',
    contact_phone TEXT NOT NULL DEFAULT '',
    desired_plan_code TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'new',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    handled_by_user_id UUID REFERENCES app_users(id) ON DELETE SET NULL,
    handled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_intake_school_name_not_blank CHECK (btrim(school_name) <> ''),
    CONSTRAINT tenant_intake_contact_required CHECK (
        btrim(contact_email) <> ''
        OR btrim(contact_phone) <> ''
    ),
    CONSTRAINT tenant_intake_status_valid CHECK (
        status IN ('new', 'contacted', 'converted', 'closed')
    )
);

CREATE INDEX IF NOT EXISTS tenant_intake_requests_status_created_idx
    ON tenant_intake_requests(status, created_at DESC);

DROP TRIGGER IF EXISTS tenant_intake_requests_set_updated_at ON tenant_intake_requests;
CREATE TRIGGER tenant_intake_requests_set_updated_at
BEFORE UPDATE ON tenant_intake_requests
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS app_password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    requested_ip INET,
    user_agent TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS app_password_reset_tokens_hash_idx
    ON app_password_reset_tokens(token_hash);

CREATE INDEX IF NOT EXISTS app_password_reset_tokens_user_created_idx
    ON app_password_reset_tokens(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS tenant_email_configs (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_by_user_id UUID REFERENCES app_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DROP TRIGGER IF EXISTS tenant_email_configs_set_updated_at ON tenant_email_configs;
CREATE TRIGGER tenant_email_configs_set_updated_at
BEFORE UPDATE ON tenant_email_configs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DELETE FROM app_role_permissions arp
USING app_roles ar, app_permissions ap
WHERE arp.role_id = ar.id
  AND arp.permission_id = ap.id
  AND ar.code IN ('tenant_owner', 'tenant_admin', 'tenant_accountant')
  AND ap.code IN (
      'email_config.view',
      'email_config.update',
      'email_cron.view',
      'email_cron.update',
      'subscription.update'
  );
