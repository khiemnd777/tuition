-- Subscription invoicing, manual collections, and dunning tracking.

CREATE TABLE IF NOT EXISTS subscription_invoices (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
	subscription_id uuid NOT NULL REFERENCES tenant_subscriptions(id) ON DELETE CASCADE,
	invoice_code text NOT NULL,
	plan_code text NOT NULL,
	plan_name text NOT NULL DEFAULT '',
	amount integer NOT NULL,
	currency text NOT NULL DEFAULT 'VND',
	period_starts_at timestamptz NOT NULL,
	period_ends_at timestamptz NOT NULL,
	due_at timestamptz NOT NULL,
	status text NOT NULL DEFAULT 'open',
	paid_at timestamptz,
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	updated_by_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	CONSTRAINT subscription_invoices_invoice_code_not_blank CHECK (btrim(invoice_code) <> ''),
	CONSTRAINT subscription_invoices_amount_non_negative CHECK (amount >= 0),
	CONSTRAINT subscription_invoices_currency_not_blank CHECK (btrim(currency) <> ''),
	CONSTRAINT subscription_invoices_period_check CHECK (period_ends_at >= period_starts_at),
	CONSTRAINT subscription_invoices_status_check CHECK (status IN ('draft', 'open', 'paid', 'past_due', 'void')),
	CONSTRAINT subscription_invoices_metadata_object CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS subscription_invoices_tenant_invoice_code_key ON subscription_invoices (tenant_id, invoice_code);
CREATE UNIQUE INDEX IF NOT EXISTS subscription_invoices_subscription_period_key ON subscription_invoices (subscription_id, period_starts_at, period_ends_at) WHERE status <> 'void';
CREATE INDEX IF NOT EXISTS subscription_invoices_status_due_idx ON subscription_invoices (tenant_id, status, due_at DESC);

CREATE TABLE IF NOT EXISTS subscription_invoice_status_history (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES subscription_invoices(id) ON DELETE CASCADE,
	from_status text NOT NULL DEFAULT '',
	to_status text NOT NULL,
	note text NOT NULL DEFAULT '',
	created_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	CONSTRAINT subscription_invoice_status_history_to_status_check CHECK (to_status IN ('draft', 'open', 'paid', 'past_due', 'void'))
);

CREATE INDEX IF NOT EXISTS subscription_invoice_status_history_invoice_idx ON subscription_invoice_status_history (invoice_id, created_at DESC);

CREATE TABLE IF NOT EXISTS subscription_dunning_runs (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
	invoice_id uuid NOT NULL REFERENCES subscription_invoices(id) ON DELETE CASCADE,
	recipient_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	recipient_email text NOT NULL,
	status text NOT NULL,
	dry_run boolean NOT NULL DEFAULT false,
	error_message text NOT NULL DEFAULT '',
	provider text NOT NULL DEFAULT '',
	provider_message_id text NOT NULL DEFAULT '',
	sent_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT subscription_dunning_runs_status_check CHECK (status IN ('sent', 'error', 'skipped', 'dry_run'))
);

CREATE INDEX IF NOT EXISTS subscription_dunning_runs_invoice_idx ON subscription_dunning_runs (invoice_id, created_at DESC);
CREATE INDEX IF NOT EXISTS subscription_dunning_runs_tenant_status_idx ON subscription_dunning_runs (tenant_id, status, created_at DESC);

DROP TRIGGER IF EXISTS subscription_invoices_set_updated_at ON subscription_invoices;
CREATE TRIGGER subscription_invoices_set_updated_at
	BEFORE UPDATE ON subscription_invoices
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();
