-- Subscription hardening: cross-tenant operations monitoring permissions.

INSERT INTO app_permissions (code, description)
VALUES
	('operation_log.cross_tenant_view', 'View operation logs across tenant boundaries'),
	('audit_log.cross_tenant_view', 'View audit logs across tenant boundaries'),
	('operations.cross_tenant.read', 'Read operation logs across tenant boundaries'),
	('audit.cross_tenant.read', 'Read audit logs across tenant boundaries')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'operation_log.cross_tenant_view',
	'audit_log.cross_tenant_view',
	'operations.cross_tenant.read',
	'audit.cross_tenant.read'
)
WHERE role.code IN ('admin', 'super_admin')
ON CONFLICT (role_id, permission_id) DO NOTHING;
