-- Invoice-based notification campaigns and delivery logs.
-- Campaigns target durable invoices and parent billing contacts instead of
-- temporary payment rows.

CREATE TABLE IF NOT EXISTS notification_templates (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	version integer NOT NULL DEFAULT 1,
	name text NOT NULL,
	subject text NOT NULL,
	email_template text NOT NULL DEFAULT 'payment_due',
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT notification_templates_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT notification_templates_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT notification_templates_subject_not_blank CHECK (btrim(subject) <> ''),
	CONSTRAINT notification_templates_version_check CHECK (version > 0),
	CONSTRAINT notification_templates_email_template_check CHECK (email_template IN ('payment_due', 'payment_paid')),
	CONSTRAINT notification_templates_status_check CHECK (status IN ('active', 'archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS notification_templates_code_version_key
	ON notification_templates (code, version);
CREATE INDEX IF NOT EXISTS notification_templates_status_idx ON notification_templates (status);

CREATE TABLE IF NOT EXISTS notification_campaigns (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	name text NOT NULL,
	campaign_type text NOT NULL,
	template_id uuid NOT NULL REFERENCES notification_templates(id) ON DELETE RESTRICT,
	school_year_id uuid REFERENCES school_years(id) ON DELETE RESTRICT,
	class_id uuid REFERENCES classes(id) ON DELETE RESTRICT,
	grade text NOT NULL DEFAULT '',
	period_code text NOT NULL DEFAULT '',
	invoice_status text NOT NULL DEFAULT '',
	due_on_or_before date,
	status text NOT NULL DEFAULT 'draft',
	target_filter jsonb NOT NULL DEFAULT '{}'::jsonb,
	last_dry_run_at timestamptz,
	sent_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT notification_campaigns_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT notification_campaigns_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT notification_campaigns_type_check CHECK (campaign_type IN ('first_notice', 'reminder')),
	CONSTRAINT notification_campaigns_invoice_status_check CHECK (invoice_status = '' OR invoice_status IN ('draft', 'unpaid', 'partial', 'paid', 'overpaid', 'manual_review', 'void')),
	CONSTRAINT notification_campaigns_status_check CHECK (status IN ('draft', 'dry_run', 'sending', 'sent', 'partial', 'archived')),
	CONSTRAINT notification_campaigns_target_filter_object CHECK (jsonb_typeof(target_filter) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS notification_campaigns_code_key ON notification_campaigns (code);
CREATE INDEX IF NOT EXISTS notification_campaigns_template_idx ON notification_campaigns (template_id);
CREATE INDEX IF NOT EXISTS notification_campaigns_filters_idx
	ON notification_campaigns (school_year_id, class_id, grade, period_code, invoice_status);
CREATE INDEX IF NOT EXISTS notification_campaigns_status_idx ON notification_campaigns (status, created_at DESC);

CREATE TABLE IF NOT EXISTS notification_recipients (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	campaign_id uuid NOT NULL REFERENCES notification_campaigns(id) ON DELETE CASCADE,
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
	parent_id uuid REFERENCES parents(id) ON DELETE SET NULL,
	recipient_name text NOT NULL DEFAULT '',
	recipient_email text NOT NULL,
	invoice_code text NOT NULL,
	student_code text NOT NULL,
	student_name text NOT NULL,
	class_name text NOT NULL,
	period_code text NOT NULL,
	amount integer NOT NULL DEFAULT 0,
	status text NOT NULL DEFAULT 'pending',
	last_error text NOT NULL DEFAULT '',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT notification_recipients_email_not_blank CHECK (btrim(recipient_email) <> ''),
	CONSTRAINT notification_recipients_invoice_code_not_blank CHECK (btrim(invoice_code) <> ''),
	CONSTRAINT notification_recipients_amount_check CHECK (amount >= 0),
	CONSTRAINT notification_recipients_status_check CHECK (status IN ('pending', 'dry_run', 'sent', 'skipped', 'error'))
);

CREATE UNIQUE INDEX IF NOT EXISTS notification_recipients_campaign_invoice_email_key
	ON notification_recipients (campaign_id, invoice_id, lower(recipient_email));
CREATE INDEX IF NOT EXISTS notification_recipients_campaign_idx
	ON notification_recipients (campaign_id, status, created_at);
CREATE INDEX IF NOT EXISTS notification_recipients_invoice_idx ON notification_recipients (invoice_id);

CREATE TABLE IF NOT EXISTS notification_logs (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	campaign_id uuid NOT NULL REFERENCES notification_campaigns(id) ON DELETE CASCADE,
	template_id uuid NOT NULL REFERENCES notification_templates(id) ON DELETE RESTRICT,
	template_code text NOT NULL,
	template_version integer NOT NULL,
	recipient_id uuid REFERENCES notification_recipients(id) ON DELETE SET NULL,
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
	recipient_email text NOT NULL,
	provider text NOT NULL DEFAULT '',
	status text NOT NULL,
	provider_message_id text NOT NULL DEFAULT '',
	error_message text NOT NULL DEFAULT '',
	dry_run boolean NOT NULL DEFAULT false,
	sent_at timestamptz NOT NULL DEFAULT now(),
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT notification_logs_template_code_not_blank CHECK (btrim(template_code) <> ''),
	CONSTRAINT notification_logs_template_version_check CHECK (template_version > 0),
	CONSTRAINT notification_logs_recipient_email_not_blank CHECK (btrim(recipient_email) <> ''),
	CONSTRAINT notification_logs_status_check CHECK (status IN ('dry_run', 'sent', 'skipped', 'error'))
);

CREATE UNIQUE INDEX IF NOT EXISTS notification_logs_send_idempotency_key
	ON notification_logs (campaign_id, template_id, invoice_id, lower(recipient_email))
	WHERE status = 'sent';
CREATE INDEX IF NOT EXISTS notification_logs_campaign_idx ON notification_logs (campaign_id, sent_at DESC);
CREATE INDEX IF NOT EXISTS notification_logs_invoice_idx ON notification_logs (invoice_id, sent_at DESC);
CREATE INDEX IF NOT EXISTS notification_logs_status_idx ON notification_logs (status, sent_at DESC);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'notification_templates_created_by_user_id_fkey') THEN
		ALTER TABLE notification_templates ADD CONSTRAINT notification_templates_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'notification_templates_updated_by_user_id_fkey') THEN
		ALTER TABLE notification_templates ADD CONSTRAINT notification_templates_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'notification_campaigns_created_by_user_id_fkey') THEN
		ALTER TABLE notification_campaigns ADD CONSTRAINT notification_campaigns_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'notification_campaigns_updated_by_user_id_fkey') THEN
		ALTER TABLE notification_campaigns ADD CONSTRAINT notification_campaigns_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS notification_templates_set_updated_at ON notification_templates;
CREATE TRIGGER notification_templates_set_updated_at
	BEFORE UPDATE ON notification_templates
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS notification_campaigns_set_updated_at ON notification_campaigns;
CREATE TRIGGER notification_campaigns_set_updated_at
	BEFORE UPDATE ON notification_campaigns
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS notification_recipients_set_updated_at ON notification_recipients;
CREATE TRIGGER notification_recipients_set_updated_at
	BEFORE UPDATE ON notification_recipients
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

INSERT INTO notification_templates (code, version, name, subject, email_template, status)
VALUES
	('first_notice', 1, 'Thông báo thanh toán lần đầu', 'Thông báo thanh toán học phí - {{student_name}}', 'payment_due', 'active'),
	('reminder', 1, 'Nhắc thanh toán học phí', 'Nhắc thanh toán học phí - {{student_name}}', 'payment_due', 'active')
ON CONFLICT (code, version) DO UPDATE
SET name = EXCLUDED.name,
	subject = EXCLUDED.subject,
	email_template = EXCLUDED.email_template,
	status = EXCLUDED.status,
	updated_at = now();

INSERT INTO app_permissions (code, description)
VALUES
	('notifications.read', 'Read notification templates, campaigns, recipients, and logs'),
	('notifications.write', 'Create and update invoice-based notification campaigns'),
	('notifications.send', 'Dry-run and send invoice-based notification campaigns')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('notifications.read', 'notifications.write', 'notifications.send')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('notifications.read', 'notifications.write', 'notifications.send')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code = 'notifications.read'
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;
