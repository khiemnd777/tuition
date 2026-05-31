-- User contact identifiers and canonical production roles/permissions.

ALTER TABLE app_users ADD COLUMN IF NOT EXISTS phone text NOT NULL DEFAULT '';
ALTER TABLE app_users ALTER COLUMN email DROP NOT NULL;
ALTER TABLE app_users ALTER COLUMN email SET DEFAULT '';

UPDATE app_users
SET email = COALESCE(email, ''),
	phone = COALESCE(phone, '');

ALTER TABLE app_users DROP CONSTRAINT IF EXISTS app_users_email_not_blank;
ALTER TABLE app_users DROP CONSTRAINT IF EXISTS app_users_contact_required;
ALTER TABLE app_users ADD CONSTRAINT app_users_contact_required
	CHECK (btrim(COALESCE(email, '')) <> '' OR btrim(COALESCE(phone, '')) <> '');

DROP INDEX IF EXISTS app_users_email_lower_key;
CREATE UNIQUE INDEX IF NOT EXISTS app_users_email_lower_key
	ON app_users (lower(email))
	WHERE btrim(COALESCE(email, '')) <> '';

CREATE UNIQUE INDEX IF NOT EXISTS app_users_phone_key
	ON app_users (phone)
	WHERE btrim(phone) <> '';

INSERT INTO app_roles (code, name, description, is_system)
VALUES
	('admin', 'Admin / Quản trị viên', 'Full application administration', true),
	('staff', 'Staff / Nhân sự', 'Manage students, school tree, and billing notifications', true),
	('accountant', 'Accountant / Kế toán', 'Manage fees, invoices, payments, reports, and email operations', true)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	description = EXCLUDED.description,
	is_system = EXCLUDED.is_system,
	updated_at = now();

INSERT INTO app_permissions (code, description)
VALUES
	('user.view', 'View application users'),
	('user.create', 'Create application users'),
	('user.update', 'Update application users'),
	('user.assign_role', 'Assign roles to application users'),
	('role.view', 'View roles and permissions'),
	('role.update', 'Update roles and permissions'),
	('student.view', 'View students, parents, and class master data'),
	('student.create', 'Create students, parents, and class master data'),
	('student.update', 'Update students, parents, and class master data'),
	('school_tree.view', 'View school, school year, grade, and class tree'),
	('school_tree.update', 'Update school, school year, grade, and class tree'),
	('fee.view', 'View fee schedules and fee previews'),
	('fee.create', 'Create fee schedules'),
	('fee.update', 'Update fee schedules'),
	('invoice.view', 'View invoices and invoice payment details'),
	('invoice.create', 'Create invoices'),
	('invoice.update', 'Update invoices'),
	('payment.view', 'View payment providers, transactions, and reconciliation'),
	('payment.create', 'Create payment intents and manual payment records'),
	('payment.reconcile', 'Reconcile payments'),
	('notification.view', 'View notification templates, campaigns, and logs'),
	('notification.create', 'Create notification campaigns'),
	('notification.send', 'Preview, dry-run, and send notifications'),
	('email_config.view', 'View local email provider configuration state'),
	('email_config.update', 'Update local email provider configuration'),
	('email_cron.view', 'View local email cron queue state'),
	('email_cron.update', 'Update and run local email cron queue'),
	('report.view', 'View accounting reports'),
	('report.export', 'Export accounting reports'),
	('dashboard.view', 'View production dashboard'),
	('operation_log.view', 'View operation logs'),
	('audit_log.view', 'View audit logs')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code IN ('admin', 'super_admin')
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
	'notification.send'
)
WHERE role.code = 'staff'
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
	'audit_log.view'
)
WHERE role.code = 'accountant'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'dashboard.view',
	'student.view',
	'school_tree.view',
	'notification.view',
	'report.view'
)
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_user_roles (user_id, role_id)
SELECT legacy_user.user_id, admin_role.id
FROM app_user_roles legacy_user
JOIN app_roles legacy_role ON legacy_role.id = legacy_user.role_id
JOIN app_roles admin_role ON admin_role.code = 'admin'
WHERE legacy_role.code = 'super_admin'
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO app_user_roles (user_id, role_id)
SELECT legacy_user.user_id, accountant_role.id
FROM app_user_roles legacy_user
JOIN app_roles legacy_role ON legacy_role.id = legacy_user.role_id
JOIN app_roles accountant_role ON accountant_role.code = 'accountant'
WHERE legacy_role.code = 'billing_admin'
ON CONFLICT (user_id, role_id) DO NOTHING;
