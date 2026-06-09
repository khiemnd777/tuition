-- Tenant foundation for subscription-mode deployments.
-- Keep existing production data under the default DEKISUGI tenant while
-- preserving schools as the business-level school/campus entity.

CREATE TABLE IF NOT EXISTS tenants (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	name text NOT NULL,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT tenants_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT tenants_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT tenants_status_check CHECK (status IN ('active', 'trial', 'suspended', 'archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS tenants_code_key ON tenants (code);
CREATE INDEX IF NOT EXISTS tenants_status_idx ON tenants (status);

INSERT INTO tenants (code, name, status)
VALUES ('DEKISUGI', 'DEKISUGI', 'active')
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	status = EXCLUDED.status,
	updated_at = now();

ALTER TABLE schools ADD COLUMN IF NOT EXISTS tenant_id uuid;

UPDATE schools
SET tenant_id = (SELECT id FROM tenants WHERE code = 'DEKISUGI')
WHERE tenant_id IS NULL;

ALTER TABLE schools ALTER COLUMN tenant_id SET NOT NULL;

DROP INDEX IF EXISTS schools_code_key;
CREATE INDEX IF NOT EXISTS schools_tenant_id_idx ON schools (tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS schools_tenant_code_key ON schools (tenant_id, code);

CREATE TABLE IF NOT EXISTS tenant_memberships (
	tenant_id uuid NOT NULL,
	user_id uuid NOT NULL,
	status text NOT NULL DEFAULT 'active',
	is_owner boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	PRIMARY KEY (tenant_id, user_id),
	CONSTRAINT tenant_memberships_status_check CHECK (status IN ('active', 'invited', 'suspended', 'removed'))
);

CREATE INDEX IF NOT EXISTS tenant_memberships_user_id_idx ON tenant_memberships (user_id);
CREATE INDEX IF NOT EXISTS tenant_memberships_status_idx ON tenant_memberships (status);

INSERT INTO tenant_memberships (tenant_id, user_id, status, is_owner)
SELECT tenant.id, app_user.id, 'active', true
FROM tenants tenant
CROSS JOIN app_users app_user
WHERE tenant.code = 'DEKISUGI'
ON CONFLICT (tenant_id, user_id) DO NOTHING;

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenants_created_by_user_id_fkey') THEN
		ALTER TABLE tenants ADD CONSTRAINT tenants_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenants_updated_by_user_id_fkey') THEN
		ALTER TABLE tenants ADD CONSTRAINT tenants_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'schools_tenant_id_fkey') THEN
		ALTER TABLE schools ADD CONSTRAINT schools_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_memberships_tenant_id_fkey') THEN
		ALTER TABLE tenant_memberships ADD CONSTRAINT tenant_memberships_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_memberships_user_id_fkey') THEN
		ALTER TABLE tenant_memberships ADD CONSTRAINT tenant_memberships_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES app_users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_memberships_created_by_user_id_fkey') THEN
		ALTER TABLE tenant_memberships ADD CONSTRAINT tenant_memberships_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_memberships_updated_by_user_id_fkey') THEN
		ALTER TABLE tenant_memberships ADD CONSTRAINT tenant_memberships_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS tenants_set_updated_at ON tenants;
CREATE TRIGGER tenants_set_updated_at
	BEFORE UPDATE ON tenants
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS tenant_memberships_set_updated_at ON tenant_memberships;
CREATE TRIGGER tenant_memberships_set_updated_at
	BEFORE UPDATE ON tenant_memberships
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();
