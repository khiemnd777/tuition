-- Subscription automation scheduler run log for renewals, dunning, and suspension policy execution.

CREATE TABLE IF NOT EXISTS subscription_automation_runs (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id uuid REFERENCES tenants(id) ON DELETE CASCADE,
	tenant_scope text NOT NULL,
	trigger_source text NOT NULL,
	status text NOT NULL,
	dry_run boolean NOT NULL DEFAULT false,
	summary jsonb NOT NULL DEFAULT '{}'::jsonb,
	error_message text NOT NULL DEFAULT '',
	started_at timestamptz NOT NULL DEFAULT now(),
	finished_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	triggered_by_user_id uuid REFERENCES app_users(id) ON DELETE SET NULL,
	CONSTRAINT subscription_automation_runs_scope_check CHECK (tenant_scope IN ('active', 'all')),
	CONSTRAINT subscription_automation_runs_trigger_check CHECK (trigger_source IN ('manual', 'scheduler')),
	CONSTRAINT subscription_automation_runs_status_check CHECK (status IN ('dry_run', 'success', 'partial', 'error')),
	CONSTRAINT subscription_automation_runs_summary_object CHECK (jsonb_typeof(summary) = 'object')
);

CREATE INDEX IF NOT EXISTS subscription_automation_runs_scope_created_at_idx
	ON subscription_automation_runs (tenant_scope, created_at DESC);

CREATE INDEX IF NOT EXISTS subscription_automation_runs_tenant_created_at_idx
	ON subscription_automation_runs (tenant_id, created_at DESC);
