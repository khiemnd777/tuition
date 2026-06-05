-- Tenant billing, plan lifecycle, and subscription enforcement foundation.

CREATE TABLE IF NOT EXISTS subscription_plans (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	name text NOT NULL,
	status text NOT NULL DEFAULT 'active',
	description text NOT NULL DEFAULT '',
	limits jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT subscription_plans_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT subscription_plans_code_format CHECK (code ~ '^[a-z][a-z0-9_:-]*$'),
	CONSTRAINT subscription_plans_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT subscription_plans_status_check CHECK (status IN ('active', 'archived')),
	CONSTRAINT subscription_plans_limits_object CHECK (jsonb_typeof(limits) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS subscription_plans_code_key ON subscription_plans (code);

CREATE TABLE IF NOT EXISTS tenant_subscriptions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
	plan_id uuid NOT NULL REFERENCES subscription_plans(id) ON DELETE RESTRICT,
	status text NOT NULL,
	trial_ends_at timestamptz,
	current_period_starts_at timestamptz,
	current_period_ends_at timestamptz,
	billing_metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	updated_by_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	CONSTRAINT tenant_subscriptions_status_check CHECK (status IN ('trial', 'active', 'past_due', 'suspended', 'cancelled')),
	CONSTRAINT tenant_subscriptions_billing_metadata_object CHECK (jsonb_typeof(billing_metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS tenant_subscriptions_tenant_id_key ON tenant_subscriptions (tenant_id);
CREATE INDEX IF NOT EXISTS tenant_subscriptions_status_idx ON tenant_subscriptions (status, tenant_id);

DROP TRIGGER IF EXISTS subscription_plans_set_updated_at ON subscription_plans;
CREATE TRIGGER subscription_plans_set_updated_at
	BEFORE UPDATE ON subscription_plans
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS tenant_subscriptions_set_updated_at ON tenant_subscriptions;
CREATE TRIGGER tenant_subscriptions_set_updated_at
	BEFORE UPDATE ON tenant_subscriptions
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();

INSERT INTO subscription_plans (code, name, status, description, limits)
VALUES
	('free_trial', 'Free Trial', 'active', 'Initial tenant onboarding and evaluation period', '{"schools":1,"operators":5}'::jsonb),
	('standard', 'Standard', 'active', 'Standard paid tenant plan', '{"schools":10,"operators":100}'::jsonb)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	status = EXCLUDED.status,
	description = EXCLUDED.description,
	limits = EXCLUDED.limits,
	updated_at = now();

INSERT INTO tenant_subscriptions (
	tenant_id,
	plan_id,
	status,
	trial_ends_at,
	current_period_starts_at,
	current_period_ends_at
)
SELECT tenant.id,
	plan.id,
	CASE
		WHEN tenant.status = 'trial' THEN 'trial'
		WHEN tenant.status = 'suspended' THEN 'suspended'
		WHEN tenant.status = 'archived' THEN 'cancelled'
		ELSE 'active'
	END,
	CASE WHEN tenant.status = 'trial' THEN tenant.created_at + interval '30 day' ELSE NULL END,
	CASE WHEN tenant.status IN ('active', 'trial') THEN tenant.created_at ELSE NULL END,
	CASE WHEN tenant.status = 'trial' THEN tenant.created_at + interval '30 day' WHEN tenant.status = 'active' THEN tenant.created_at + interval '1 year' ELSE NULL END
FROM tenants tenant
JOIN subscription_plans plan ON plan.code = CASE WHEN tenant.status = 'trial' THEN 'free_trial' ELSE 'standard' END
ON CONFLICT (tenant_id) DO NOTHING;

INSERT INTO app_permissions (code, description)
VALUES
	('subscription.view', 'View tenant subscription plans and lifecycle'),
	('subscription.update', 'Update the active tenant subscription state')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('subscription.view', 'subscription.update')
WHERE role.code IN ('admin', 'super_admin')
ON CONFLICT (role_id, permission_id) DO NOTHING;
