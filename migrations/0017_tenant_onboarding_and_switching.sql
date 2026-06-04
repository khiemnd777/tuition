-- Tenant onboarding and switching permissions for subscription-mode deployments.
-- Tenant-owned data was isolated in earlier migrations; this adds the operator
-- controls needed to create a new tenant and switch the active tenant context.

CREATE INDEX IF NOT EXISTS tenant_memberships_user_status_idx
	ON tenant_memberships (user_id, status, tenant_id);

INSERT INTO app_permissions (code, description)
VALUES
	('tenant.view', 'View available subscription tenants'),
	('tenant.create', 'Create subscription tenants and initial schools'),
	('tenant.update', 'Update the active subscription tenant'),
	('tenant.switch', 'Switch active subscription tenant context')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code IN ('admin', 'super_admin')
	AND permission.code IN ('tenant.view', 'tenant.create', 'tenant.update', 'tenant.switch')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('tenant.view', 'tenant.switch')
WHERE role.code IN ('staff', 'accountant', 'billing_admin', 'viewer')
ON CONFLICT (role_id, permission_id) DO NOTHING;
