-- Reports, audit review, and operational readiness.
-- Keep prior migration checksums stable; add reporting/export permissions and
-- durable operation logs for production failure review.

ALTER TABLE manual_cash_receipts
	ADD COLUMN IF NOT EXISTS reason text NOT NULL DEFAULT '';

UPDATE manual_cash_receipts
SET reason = COALESCE(NULLIF(btrim(reason), ''), NULLIF(btrim(note), ''), 'legacy manual cash receipt')
WHERE reason IS NULL OR btrim(reason) = '';

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'manual_cash_receipts_reason_not_blank') THEN
		ALTER TABLE manual_cash_receipts ADD CONSTRAINT manual_cash_receipts_reason_not_blank
			CHECK (btrim(reason) <> '');
	END IF;
END $$;

CREATE TABLE IF NOT EXISTS operation_logs (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	occurred_at timestamptz NOT NULL DEFAULT now(),
	source text NOT NULL,
	level text NOT NULL DEFAULT 'error',
	operation text NOT NULL,
	status text NOT NULL DEFAULT 'error',
	message text NOT NULL DEFAULT '',
	entity_type text NOT NULL DEFAULT '',
	entity_id uuid,
	request_id text NOT NULL DEFAULT '',
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT operation_logs_source_not_blank CHECK (btrim(source) <> ''),
	CONSTRAINT operation_logs_operation_not_blank CHECK (btrim(operation) <> ''),
	CONSTRAINT operation_logs_level_check CHECK (level IN ('info', 'warn', 'error')),
	CONSTRAINT operation_logs_status_not_blank CHECK (btrim(status) <> ''),
	CONSTRAINT operation_logs_metadata_object CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS operation_logs_occurred_at_idx ON operation_logs (occurred_at DESC);
CREATE INDEX IF NOT EXISTS operation_logs_source_idx ON operation_logs (source, occurred_at DESC);
CREATE INDEX IF NOT EXISTS operation_logs_level_idx ON operation_logs (level, occurred_at DESC);
CREATE INDEX IF NOT EXISTS operation_logs_entity_idx ON operation_logs (entity_type, entity_id);

INSERT INTO app_permissions (code, description)
VALUES
	('admin.reports.export', 'Export production reports to CSV'),
	('operations.read', 'Read production operational logs')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('admin.reports.export', 'operations.read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('admin.reports.export', 'operations.read')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code = 'admin.reports.export'
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;
