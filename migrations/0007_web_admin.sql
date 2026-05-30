-- Web admin read models and permission seeds.
-- Initiative 7 does not add business tables; it exposes dashboard, reports,
-- and role administration over the durable production tables.

INSERT INTO app_permissions (code, description)
VALUES
	('admin.dashboard.read', 'Read production admin dashboard metrics'),
	('admin.reports.read', 'Read production admin reports'),
	('system.users.assign_roles', 'Assign roles to application users')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('admin.dashboard.read', 'admin.reports.read', 'system.users.assign_roles')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('admin.dashboard.read', 'admin.reports.read')
WHERE role.code IN ('billing_admin', 'viewer')
ON CONFLICT (role_id, permission_id) DO NOTHING;
