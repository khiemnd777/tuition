-- Automatic paid-confirmation notifications.
-- Reuse the existing payment_paid email renderer and notification log table.

ALTER TABLE notification_campaigns
	DROP CONSTRAINT IF EXISTS notification_campaigns_type_check;

ALTER TABLE notification_campaigns
	ADD CONSTRAINT notification_campaigns_type_check
	CHECK (campaign_type IN ('first_notice', 'reminder', 'payment_confirmation'));

INSERT INTO notification_templates (code, version, name, subject, email_template, status)
VALUES
	('payment_confirmation', 1, 'Xác nhận đã thanh toán học phí', 'Xác nhận đã thanh toán học phí - {{student_name}}', 'payment_paid', 'active')
ON CONFLICT (code, version) DO UPDATE
SET name = EXCLUDED.name,
	subject = EXCLUDED.subject,
	email_template = EXCLUDED.email_template,
	status = EXCLUDED.status,
	updated_at = now();

INSERT INTO notification_campaigns (
	code, name, campaign_type, template_id, invoice_status, status, target_filter
)
SELECT
	'auto_paid_confirmation',
	'Tự động xác nhận đã thanh toán',
	'payment_confirmation',
	t.id,
	'paid',
	'draft',
	'{"auto": true, "trigger": "invoice_paid"}'::jsonb
FROM notification_templates t
WHERE t.code = 'payment_confirmation'
	AND t.status = 'active'
ORDER BY t.version DESC
LIMIT 1
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	campaign_type = EXCLUDED.campaign_type,
	template_id = EXCLUDED.template_id,
	invoice_status = EXCLUDED.invoice_status,
	target_filter = EXCLUDED.target_filter,
	status = CASE
		WHEN notification_campaigns.status = 'archived' THEN 'draft'
		ELSE notification_campaigns.status
	END,
	updated_at = now();
