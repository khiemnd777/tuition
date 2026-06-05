-- Tenant feature entitlements, usage counters, and metering support.

CREATE TABLE IF NOT EXISTS tenant_usage_counters (
	tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
	metric_code text NOT NULL,
	period_key text NOT NULL DEFAULT '',
	used_count integer NOT NULL DEFAULT 0,
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, metric_code, period_key),
	CONSTRAINT tenant_usage_counters_metric_code_check CHECK (metric_code IN ('schools', 'operators', 'students', 'monthly_notifications')),
	CONSTRAINT tenant_usage_counters_period_key_not_blank CHECK (period_key IS NOT NULL),
	CONSTRAINT tenant_usage_counters_used_count_non_negative CHECK (used_count >= 0)
);

CREATE INDEX IF NOT EXISTS tenant_usage_counters_metric_idx ON tenant_usage_counters (metric_code, period_key, tenant_id);

UPDATE subscription_plans
SET limits = limits || CASE code
	WHEN 'free_trial' THEN '{"schools":1,"operators":5,"students":200,"monthly_notifications":500}'::jsonb
	WHEN 'standard' THEN '{"schools":10,"operators":100,"students":5000,"monthly_notifications":10000}'::jsonb
	ELSE '{}'::jsonb
END,
	updated_at = now()
WHERE code IN ('free_trial', 'standard');

INSERT INTO tenant_usage_counters (tenant_id, metric_code, period_key, used_count)
SELECT tenant.id, 'schools', '', COUNT(school.id)::integer
FROM tenants tenant
LEFT JOIN schools school ON school.tenant_id = tenant.id
GROUP BY tenant.id
ON CONFLICT (tenant_id, metric_code, period_key) DO UPDATE
SET used_count = EXCLUDED.used_count,
	updated_at = now();

INSERT INTO tenant_usage_counters (tenant_id, metric_code, period_key, used_count)
SELECT tenant.id, 'operators', '', COUNT(DISTINCT tur.user_id)::integer
FROM tenants tenant
LEFT JOIN tenant_user_roles tur ON tur.tenant_id = tenant.id
LEFT JOIN tenant_memberships tm
	ON tm.tenant_id = tur.tenant_id
	AND tm.user_id = tur.user_id
	AND tm.status = 'active'
WHERE tur.user_id IS NULL OR tm.user_id IS NOT NULL
GROUP BY tenant.id
ON CONFLICT (tenant_id, metric_code, period_key) DO UPDATE
SET used_count = EXCLUDED.used_count,
	updated_at = now();

INSERT INTO tenant_usage_counters (tenant_id, metric_code, period_key, used_count)
SELECT tenant.id, 'students', '', COUNT(student.id)::integer
FROM tenants tenant
LEFT JOIN students student ON student.tenant_id = tenant.id AND student.status <> 'inactive'
GROUP BY tenant.id
ON CONFLICT (tenant_id, metric_code, period_key) DO UPDATE
SET used_count = EXCLUDED.used_count,
	updated_at = now();

INSERT INTO tenant_usage_counters (tenant_id, metric_code, period_key, used_count)
SELECT campaign.tenant_id, 'monthly_notifications',
	to_char(COALESCE(log.sent_at, log.created_at) AT TIME ZONE 'UTC', 'YYYY-MM'),
	COUNT(log.id)::integer
FROM notification_logs log
JOIN notification_campaigns campaign ON campaign.id = log.campaign_id
WHERE log.status = 'sent'
	AND NOT log.dry_run
GROUP BY campaign.tenant_id, to_char(COALESCE(log.sent_at, log.created_at) AT TIME ZONE 'UTC', 'YYYY-MM')
ON CONFLICT (tenant_id, metric_code, period_key) DO UPDATE
SET used_count = EXCLUDED.used_count,
	updated_at = now();
