-- Tenant data isolation for subscription-mode deployments.
-- Keep shared lookup tables global, but make tenant-owned records explicit
-- where tenant scope cannot be inferred from a direct school join.

ALTER TABLE students ADD COLUMN IF NOT EXISTS tenant_id uuid;
ALTER TABLE parents ADD COLUMN IF NOT EXISTS tenant_id uuid;
ALTER TABLE notification_campaigns ADD COLUMN IF NOT EXISTS tenant_id uuid;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS tenant_id uuid;
ALTER TABLE operation_logs ADD COLUMN IF NOT EXISTS tenant_id uuid;

UPDATE students student
SET tenant_id = school.tenant_id
FROM classes class
JOIN school_years school_year ON school_year.id = class.school_year_id
JOIN schools school ON school.id = school_year.school_id
WHERE student.class_id = class.id
	AND student.tenant_id IS NULL;

UPDATE parents parent
SET tenant_id = tenant.id
FROM tenants tenant
WHERE tenant.code = 'DEKISUGI'
	AND parent.tenant_id IS NULL;

UPDATE notification_campaigns campaign
SET tenant_id = COALESCE(
	(
		SELECT school.tenant_id
		FROM school_years school_year
		JOIN schools school ON school.id = school_year.school_id
		WHERE school_year.id = campaign.school_year_id
		LIMIT 1
	),
	(SELECT id FROM tenants WHERE code = 'DEKISUGI')
)
WHERE campaign.tenant_id IS NULL;

ALTER TABLE audit_logs DISABLE TRIGGER audit_logs_prevent_update;

UPDATE audit_logs audit_log
SET tenant_id = tenant.id
FROM tenants tenant
WHERE tenant.code = 'DEKISUGI'
	AND audit_log.tenant_id IS NULL;

ALTER TABLE audit_logs ENABLE TRIGGER audit_logs_prevent_update;

UPDATE operation_logs operation_log
SET tenant_id = tenant.id
FROM tenants tenant
WHERE tenant.code = 'DEKISUGI'
	AND operation_log.tenant_id IS NULL;

ALTER TABLE students ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE parents ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE notification_campaigns ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE audit_logs ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE operation_logs ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE students DROP CONSTRAINT IF EXISTS students_student_code_key;
DROP INDEX IF EXISTS parents_email_key;
DROP INDEX IF EXISTS notification_campaigns_code_key;

CREATE INDEX IF NOT EXISTS students_tenant_id_idx ON students (tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS students_tenant_code_key ON students (tenant_id, student_code);

CREATE INDEX IF NOT EXISTS parents_tenant_id_idx ON parents (tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS parents_tenant_email_key ON parents (tenant_id, email) WHERE email <> '';

CREATE INDEX IF NOT EXISTS notification_campaigns_tenant_id_idx ON notification_campaigns (tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS notification_campaigns_tenant_code_key ON notification_campaigns (tenant_id, code);

CREATE INDEX IF NOT EXISTS audit_logs_tenant_id_idx ON audit_logs (tenant_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS operation_logs_tenant_id_idx ON operation_logs (tenant_id, occurred_at DESC);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'students_tenant_id_fkey') THEN
		ALTER TABLE students ADD CONSTRAINT students_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'parents_tenant_id_fkey') THEN
		ALTER TABLE parents ADD CONSTRAINT parents_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'notification_campaigns_tenant_id_fkey') THEN
		ALTER TABLE notification_campaigns ADD CONSTRAINT notification_campaigns_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'audit_logs_tenant_id_fkey') THEN
		ALTER TABLE audit_logs ADD CONSTRAINT audit_logs_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'operation_logs_tenant_id_fkey') THEN
		ALTER TABLE operation_logs ADD CONSTRAINT operation_logs_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;
END $$;
