-- Foundation schema for production persistence.
-- Primary IDs are UUIDs generated in PostgreSQL. Operational tables use
-- timestamptz audit columns so later modules can preserve actor and time data.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS app_users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	email text NOT NULL,
	display_name text NOT NULL DEFAULT '',
	status text NOT NULL DEFAULT 'active',
	last_login_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT app_users_email_not_blank CHECK (btrim(email) <> ''),
	CONSTRAINT app_users_status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE UNIQUE INDEX IF NOT EXISTS app_users_email_lower_key ON app_users (lower(email));
CREATE INDEX IF NOT EXISTS app_users_status_idx ON app_users (status);

CREATE TABLE IF NOT EXISTS app_roles (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	name text NOT NULL,
	description text NOT NULL DEFAULT '',
	is_system boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT app_roles_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT app_roles_code_format CHECK (code ~ '^[a-z][a-z0-9_.:-]*$'),
	CONSTRAINT app_roles_name_not_blank CHECK (btrim(name) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS app_roles_code_key ON app_roles (code);

CREATE TABLE IF NOT EXISTS app_permissions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	description text NOT NULL DEFAULT '',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT app_permissions_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT app_permissions_code_format CHECK (code ~ '^[a-z][a-z0-9_.:-]*$')
);

CREATE UNIQUE INDEX IF NOT EXISTS app_permissions_code_key ON app_permissions (code);

CREATE TABLE IF NOT EXISTS app_user_roles (
	user_id uuid NOT NULL,
	role_id uuid NOT NULL,
	assigned_at timestamptz NOT NULL DEFAULT now(),
	assigned_by_user_id uuid,
	PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS app_user_roles_role_id_idx ON app_user_roles (role_id);

CREATE TABLE IF NOT EXISTS app_role_permissions (
	role_id uuid NOT NULL,
	permission_id uuid NOT NULL,
	granted_at timestamptz NOT NULL DEFAULT now(),
	granted_by_user_id uuid,
	PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS app_role_permissions_permission_id_idx ON app_role_permissions (permission_id);

CREATE TABLE IF NOT EXISTS audit_logs (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	occurred_at timestamptz NOT NULL DEFAULT now(),
	actor_user_id uuid,
	action text NOT NULL,
	entity_type text NOT NULL,
	entity_id uuid,
	request_id text NOT NULL DEFAULT '',
	ip_address inet,
	user_agent text NOT NULL DEFAULT '',
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT audit_logs_action_not_blank CHECK (btrim(action) <> ''),
	CONSTRAINT audit_logs_entity_type_not_blank CHECK (btrim(entity_type) <> ''),
	CONSTRAINT audit_logs_metadata_object CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS audit_logs_occurred_at_idx ON audit_logs (occurred_at DESC);
CREATE INDEX IF NOT EXISTS audit_logs_actor_user_id_idx ON audit_logs (actor_user_id);
CREATE INDEX IF NOT EXISTS audit_logs_entity_idx ON audit_logs (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS audit_logs_action_idx ON audit_logs (action);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_users_created_by_user_id_fkey') THEN
		ALTER TABLE app_users ADD CONSTRAINT app_users_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_users_updated_by_user_id_fkey') THEN
		ALTER TABLE app_users ADD CONSTRAINT app_users_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_roles_created_by_user_id_fkey') THEN
		ALTER TABLE app_roles ADD CONSTRAINT app_roles_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_roles_updated_by_user_id_fkey') THEN
		ALTER TABLE app_roles ADD CONSTRAINT app_roles_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_permissions_created_by_user_id_fkey') THEN
		ALTER TABLE app_permissions ADD CONSTRAINT app_permissions_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_permissions_updated_by_user_id_fkey') THEN
		ALTER TABLE app_permissions ADD CONSTRAINT app_permissions_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_user_roles_user_id_fkey') THEN
		ALTER TABLE app_user_roles ADD CONSTRAINT app_user_roles_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES app_users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_user_roles_role_id_fkey') THEN
		ALTER TABLE app_user_roles ADD CONSTRAINT app_user_roles_role_id_fkey
			FOREIGN KEY (role_id) REFERENCES app_roles(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_user_roles_assigned_by_user_id_fkey') THEN
		ALTER TABLE app_user_roles ADD CONSTRAINT app_user_roles_assigned_by_user_id_fkey
			FOREIGN KEY (assigned_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_role_permissions_role_id_fkey') THEN
		ALTER TABLE app_role_permissions ADD CONSTRAINT app_role_permissions_role_id_fkey
			FOREIGN KEY (role_id) REFERENCES app_roles(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_role_permissions_permission_id_fkey') THEN
		ALTER TABLE app_role_permissions ADD CONSTRAINT app_role_permissions_permission_id_fkey
			FOREIGN KEY (permission_id) REFERENCES app_permissions(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_role_permissions_granted_by_user_id_fkey') THEN
		ALTER TABLE app_role_permissions ADD CONSTRAINT app_role_permissions_granted_by_user_id_fkey
			FOREIGN KEY (granted_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'audit_logs_actor_user_id_fkey') THEN
		ALTER TABLE audit_logs ADD CONSTRAINT audit_logs_actor_user_id_fkey
			FOREIGN KEY (actor_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

CREATE OR REPLACE FUNCTION abc_set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS app_users_set_updated_at ON app_users;
CREATE TRIGGER app_users_set_updated_at
	BEFORE UPDATE ON app_users
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS app_roles_set_updated_at ON app_roles;
CREATE TRIGGER app_roles_set_updated_at
	BEFORE UPDATE ON app_roles
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS app_permissions_set_updated_at ON app_permissions;
CREATE TRIGGER app_permissions_set_updated_at
	BEFORE UPDATE ON app_permissions
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

CREATE OR REPLACE FUNCTION abc_prevent_audit_log_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
	RAISE EXCEPTION 'audit_logs entries are immutable';
END;
$$;

DROP TRIGGER IF EXISTS audit_logs_prevent_update ON audit_logs;
CREATE TRIGGER audit_logs_prevent_update
	BEFORE UPDATE ON audit_logs
	FOR EACH ROW
	EXECUTE FUNCTION abc_prevent_audit_log_mutation();

DROP TRIGGER IF EXISTS audit_logs_prevent_delete ON audit_logs;
CREATE TRIGGER audit_logs_prevent_delete
	BEFORE DELETE ON audit_logs
	FOR EACH ROW
	EXECUTE FUNCTION abc_prevent_audit_log_mutation();

INSERT INTO app_permissions (code, description)
VALUES
	('system.users.read', 'Read users and role assignments'),
	('system.users.write', 'Create and update users and role assignments'),
	('system.roles.read', 'Read roles and permissions'),
	('system.roles.write', 'Create and update roles and permissions'),
	('audit.read', 'Read audit logs')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_roles (code, name, description, is_system)
VALUES
	('super_admin', 'Super Admin', 'Full administrative access', true),
	('billing_admin', 'Billing Admin', 'Manage billing operations', true),
	('viewer', 'Viewer', 'Read-only operational access', true)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	description = EXCLUDED.description,
	is_system = EXCLUDED.is_system,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('system.users.read', 'system.roles.read', 'audit.read')
WHERE role.code IN ('billing_admin', 'viewer')
ON CONFLICT (role_id, permission_id) DO NOTHING;
