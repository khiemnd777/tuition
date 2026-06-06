-- Split platform admin and tenant operator roles for subscription SaaS actor model.

INSERT INTO app_roles (code, name, description, is_system)
VALUES
	('platform_admin', 'Platform Admin', 'Manage platform-wide tenants, subscriptions, finance, and internal operations', true),
	('tenant_owner', 'Tenant Owner', 'Own a tenant workspace, subscription, and operator access', true),
	('tenant_admin', 'Tenant Admin', 'Manage tenant operators and day-to-day tenant administration', true),
	('tenant_staff', 'Tenant Staff', 'Operate tenant master data and notification workflows', true),
	('tenant_accountant', 'Tenant Accountant', 'Operate tenant billing, reconciliation, and reporting workflows', true)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	description = EXCLUDED.description,
	is_system = EXCLUDED.is_system,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON true
WHERE role.code = 'platform_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'dashboard.view',
	'user.view',
	'user.create',
	'user.update',
	'user.assign_role',
	'role.view',
	'student.view',
	'student.create',
	'student.update',
	'school_tree.view',
	'school_tree.update',
	'fee.view',
	'fee.create',
	'fee.update',
	'invoice.view',
	'invoice.create',
	'invoice.update',
	'payment.view',
	'payment.create',
	'payment.reconcile',
	'notification.view',
	'notification.create',
	'notification.send',
	'email_config.view',
	'email_config.update',
	'email_cron.view',
	'email_cron.update',
	'report.view',
	'report.export',
	'operation_log.view',
	'audit_log.view',
	'tenant.switch',
	'subscription.view',
	'subscription.update'
)
WHERE role.code = 'tenant_owner'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'dashboard.view',
	'user.view',
	'user.create',
	'user.update',
	'user.assign_role',
	'role.view',
	'student.view',
	'student.create',
	'student.update',
	'school_tree.view',
	'school_tree.update',
	'fee.view',
	'fee.create',
	'fee.update',
	'invoice.view',
	'invoice.create',
	'invoice.update',
	'payment.view',
	'payment.create',
	'payment.reconcile',
	'notification.view',
	'notification.create',
	'notification.send',
	'email_config.view',
	'email_config.update',
	'email_cron.view',
	'email_cron.update',
	'report.view',
	'report.export',
	'operation_log.view',
	'audit_log.view',
	'tenant.switch',
	'subscription.view'
)
WHERE role.code = 'tenant_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'dashboard.view',
	'student.view',
	'student.create',
	'student.update',
	'school_tree.view',
	'school_tree.update',
	'notification.view',
	'notification.create',
	'notification.send',
	'tenant.switch'
)
WHERE role.code = 'tenant_staff'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'dashboard.view',
	'student.view',
	'school_tree.view',
	'fee.view',
	'fee.create',
	'fee.update',
	'invoice.view',
	'invoice.create',
	'invoice.update',
	'payment.view',
	'payment.create',
	'payment.reconcile',
	'notification.view',
	'notification.create',
	'notification.send',
	'email_config.view',
	'email_config.update',
	'email_cron.view',
	'email_cron.update',
	'report.view',
	'report.export',
	'audit_log.view',
	'tenant.switch',
	'subscription.view'
)
WHERE role.code = 'tenant_accountant'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_user_roles (user_id, role_id, assigned_at, assigned_by_user_id)
SELECT ur.user_id, new_role.id, ur.assigned_at, ur.assigned_by_user_id
FROM app_user_roles ur
JOIN app_roles legacy_role ON legacy_role.id = ur.role_id
JOIN app_roles new_role ON new_role.code = 'platform_admin'
WHERE legacy_role.code IN ('admin', 'super_admin')
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO tenant_user_roles (tenant_id, user_id, role_id, assigned_at, assigned_by_user_id)
SELECT tur.tenant_id, tur.user_id, new_role.id, tur.assigned_at, tur.assigned_by_user_id
FROM tenant_user_roles tur
JOIN app_roles legacy_role ON legacy_role.id = tur.role_id
JOIN app_roles new_role ON new_role.code = CASE legacy_role.code
	WHEN 'admin' THEN 'tenant_admin'
	WHEN 'staff' THEN 'tenant_staff'
	WHEN 'accountant' THEN 'tenant_accountant'
END
WHERE legacy_role.code IN ('admin', 'staff', 'accountant')
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;

INSERT INTO tenant_user_roles (tenant_id, user_id, role_id)
SELECT membership.tenant_id, membership.user_id, role.id
FROM tenant_memberships membership
JOIN app_roles role ON role.code = 'tenant_owner'
WHERE membership.is_owner
	AND membership.status = 'active'
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;

DELETE FROM tenant_user_roles tur
USING app_roles role
WHERE role.id = tur.role_id
	AND role.code IN ('admin', 'staff', 'accountant');
