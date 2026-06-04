-- Tenant-aware auth sessions and RBAC assignments.
-- Keep legacy app_user_roles in place while moving runtime authorization to
-- tenant-scoped role assignments.

ALTER TABLE app_auth_sessions ADD COLUMN IF NOT EXISTS tenant_id uuid;

UPDATE app_auth_sessions
SET tenant_id = (SELECT id FROM tenants WHERE code = 'ABC_SUN')
WHERE tenant_id IS NULL;

ALTER TABLE app_auth_sessions ALTER COLUMN tenant_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS app_auth_sessions_tenant_id_idx ON app_auth_sessions (tenant_id);
CREATE INDEX IF NOT EXISTS app_auth_sessions_tenant_user_idx ON app_auth_sessions (tenant_id, user_id);

CREATE TABLE IF NOT EXISTS tenant_user_roles (
	tenant_id uuid NOT NULL,
	user_id uuid NOT NULL,
	role_id uuid NOT NULL,
	assigned_at timestamptz NOT NULL DEFAULT now(),
	assigned_by_user_id uuid,
	PRIMARY KEY (tenant_id, user_id, role_id)
);

CREATE INDEX IF NOT EXISTS tenant_user_roles_user_id_idx ON tenant_user_roles (user_id);
CREATE INDEX IF NOT EXISTS tenant_user_roles_role_id_idx ON tenant_user_roles (role_id);

INSERT INTO tenant_user_roles (tenant_id, user_id, role_id, assigned_at, assigned_by_user_id)
SELECT tenant.id, ur.user_id, ur.role_id, ur.assigned_at, ur.assigned_by_user_id
FROM tenants tenant
CROSS JOIN app_user_roles ur
WHERE tenant.code = 'ABC_SUN'
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_sessions_tenant_id_fkey') THEN
		ALTER TABLE app_auth_sessions ADD CONSTRAINT app_auth_sessions_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_user_roles_tenant_id_fkey') THEN
		ALTER TABLE tenant_user_roles ADD CONSTRAINT tenant_user_roles_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_user_roles_user_id_fkey') THEN
		ALTER TABLE tenant_user_roles ADD CONSTRAINT tenant_user_roles_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES app_users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_user_roles_role_id_fkey') THEN
		ALTER TABLE tenant_user_roles ADD CONSTRAINT tenant_user_roles_role_id_fkey
			FOREIGN KEY (role_id) REFERENCES app_roles(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_user_roles_assigned_by_user_id_fkey') THEN
		ALTER TABLE tenant_user_roles ADD CONSTRAINT tenant_user_roles_assigned_by_user_id_fkey
			FOREIGN KEY (assigned_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;
