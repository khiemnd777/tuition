-- Route-level RBAC permission seeds for protected legacy email operations.

INSERT INTO app_permissions (code, description)
VALUES
	('email.config.read', 'Read local email provider configuration state'),
	('email.config.write', 'Update local email provider configuration'),
	('email.send', 'Preview, dry-run, and send legacy payment emails'),
	('email.cron.manage', 'Read, update, and run the local email cron queue')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('email.config.read', 'email.config.write', 'email.send', 'email.cron.manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('email.config.read', 'email.config.write', 'email.send', 'email.cron.manage')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;
