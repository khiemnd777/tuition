-- Production invoices and receipt documents.
-- Invoices snapshot generated fee data so payment requests do not change when
-- fee schedules are edited later.

CREATE TABLE IF NOT EXISTS invoices (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	fee_schedule_id uuid NOT NULL REFERENCES fee_schedules(id) ON DELETE RESTRICT,
	student_id uuid NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
	school_year_id uuid NOT NULL REFERENCES school_years(id) ON DELETE RESTRICT,
	class_id uuid NOT NULL REFERENCES classes(id) ON DELETE RESTRICT,
	invoice_code text NOT NULL,
	student_code text NOT NULL,
	student_name text NOT NULL,
	class_name text NOT NULL,
	grade text NOT NULL,
	school_year_code text NOT NULL,
	period_code text NOT NULL,
	month smallint,
	issued_at timestamptz NOT NULL DEFAULT now(),
	due_date date,
	status text NOT NULL DEFAULT 'unpaid',
	base_amount integer NOT NULL DEFAULT 0,
	adjustment_amount integer NOT NULL DEFAULT 0,
	total_amount integer NOT NULL DEFAULT 0,
	paid_amount integer NOT NULL DEFAULT 0,
	currency text NOT NULL DEFAULT 'VND',
	collection_bank_bin text NOT NULL,
	collection_bank_account text NOT NULL,
	qr_bill_number text NOT NULL,
	qr_note text NOT NULL,
	source_snapshot jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT invoices_invoice_code_not_blank CHECK (btrim(invoice_code) <> ''),
	CONSTRAINT invoices_student_code_not_blank CHECK (btrim(student_code) <> ''),
	CONSTRAINT invoices_student_name_not_blank CHECK (btrim(student_name) <> ''),
	CONSTRAINT invoices_class_name_not_blank CHECK (btrim(class_name) <> ''),
	CONSTRAINT invoices_period_not_blank CHECK (btrim(period_code) <> ''),
	CONSTRAINT invoices_status_check CHECK (status IN ('draft', 'unpaid', 'partial', 'paid', 'overpaid', 'manual_review', 'void')),
	CONSTRAINT invoices_month_check CHECK (month IS NULL OR (month >= 1 AND month <= 12)),
	CONSTRAINT invoices_amounts_check CHECK (
		base_amount >= 0 AND total_amount >= 0 AND paid_amount >= 0
	),
	CONSTRAINT invoices_bank_bin_check CHECK (collection_bank_bin ~ '^[0-9]{6}$'),
	CONSTRAINT invoices_bank_account_not_blank CHECK (btrim(collection_bank_account) <> ''),
	CONSTRAINT invoices_qr_bill_number_not_blank CHECK (btrim(qr_bill_number) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS invoices_invoice_code_key ON invoices (invoice_code);
CREATE UNIQUE INDEX IF NOT EXISTS invoices_schedule_student_active_key
	ON invoices (fee_schedule_id, student_id)
	WHERE status <> 'void';
CREATE INDEX IF NOT EXISTS invoices_fee_schedule_idx ON invoices (fee_schedule_id);
CREATE INDEX IF NOT EXISTS invoices_student_idx ON invoices (student_id);
CREATE INDEX IF NOT EXISTS invoices_school_year_idx ON invoices (school_year_id);
CREATE INDEX IF NOT EXISTS invoices_class_idx ON invoices (class_id);
CREATE INDEX IF NOT EXISTS invoices_period_idx ON invoices (period_code);
CREATE INDEX IF NOT EXISTS invoices_status_idx ON invoices (status);

CREATE TABLE IF NOT EXISTS invoice_items (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
	fee_type_code text NOT NULL DEFAULT '',
	label_vi text NOT NULL,
	label_en text NOT NULL,
	amount integer NOT NULL DEFAULT 0,
	display_order integer NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	CONSTRAINT invoice_items_label_vi_not_blank CHECK (btrim(label_vi) <> ''),
	CONSTRAINT invoice_items_label_en_not_blank CHECK (btrim(label_en) <> ''),
	CONSTRAINT invoice_items_amount_check CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS invoice_items_invoice_idx ON invoice_items (invoice_id, display_order, id);

CREATE TABLE IF NOT EXISTS invoice_adjustments (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
	adjustment_type text NOT NULL,
	fee_type_code text NOT NULL DEFAULT '',
	label_vi text NOT NULL,
	label_en text NOT NULL,
	amount integer NOT NULL DEFAULT 0,
	delta integer NOT NULL DEFAULT 0,
	reason text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	CONSTRAINT invoice_adjustments_type_check CHECK (adjustment_type IN ('discount', 'surcharge', 'waiver', 'carry_over')),
	CONSTRAINT invoice_adjustments_label_vi_not_blank CHECK (btrim(label_vi) <> ''),
	CONSTRAINT invoice_adjustments_label_en_not_blank CHECK (btrim(label_en) <> ''),
	CONSTRAINT invoice_adjustments_reason_not_blank CHECK (btrim(reason) <> '')
);

CREATE INDEX IF NOT EXISTS invoice_adjustments_invoice_idx ON invoice_adjustments (invoice_id);

CREATE TABLE IF NOT EXISTS invoice_status_history (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
	from_status text NOT NULL DEFAULT '',
	to_status text NOT NULL,
	reason text NOT NULL DEFAULT '',
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	CONSTRAINT invoice_status_history_to_status_check CHECK (to_status IN ('draft', 'unpaid', 'partial', 'paid', 'overpaid', 'manual_review', 'void'))
);

CREATE INDEX IF NOT EXISTS invoice_status_history_invoice_idx ON invoice_status_history (invoice_id, created_at);

CREATE TABLE IF NOT EXISTS receipt_documents (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
	document_type text NOT NULL DEFAULT 'pdf_receipt',
	content_type text NOT NULL DEFAULT 'application/pdf',
	storage_kind text NOT NULL DEFAULT 'generated',
	storage_key text NOT NULL DEFAULT '',
	checksum text NOT NULL DEFAULT '',
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	generated_at timestamptz NOT NULL DEFAULT now(),
	generated_by_user_id uuid,
	CONSTRAINT receipt_documents_type_not_blank CHECK (btrim(document_type) <> ''),
	CONSTRAINT receipt_documents_content_type_not_blank CHECK (btrim(content_type) <> ''),
	CONSTRAINT receipt_documents_storage_kind_check CHECK (storage_kind IN ('generated', 'local_file', 'object_storage'))
);

CREATE INDEX IF NOT EXISTS receipt_documents_invoice_idx ON receipt_documents (invoice_id, generated_at DESC);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'invoices_created_by_user_id_fkey') THEN
		ALTER TABLE invoices ADD CONSTRAINT invoices_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'invoices_updated_by_user_id_fkey') THEN
		ALTER TABLE invoices ADD CONSTRAINT invoices_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'invoice_items_created_by_user_id_fkey') THEN
		ALTER TABLE invoice_items ADD CONSTRAINT invoice_items_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'invoice_adjustments_created_by_user_id_fkey') THEN
		ALTER TABLE invoice_adjustments ADD CONSTRAINT invoice_adjustments_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'invoice_status_history_created_by_user_id_fkey') THEN
		ALTER TABLE invoice_status_history ADD CONSTRAINT invoice_status_history_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'receipt_documents_generated_by_user_id_fkey') THEN
		ALTER TABLE receipt_documents ADD CONSTRAINT receipt_documents_generated_by_user_id_fkey
			FOREIGN KEY (generated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS invoices_set_updated_at ON invoices;
CREATE TRIGGER invoices_set_updated_at
	BEFORE UPDATE ON invoices
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

INSERT INTO app_permissions (code, description)
VALUES
	('invoices.read', 'Read invoices, invoice QR payment data, and PDF receipts'),
	('invoices.write', 'Generate and regenerate invoices from fee schedules')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('invoices.read', 'invoices.write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('invoices.read', 'invoices.write')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code = 'invoices.read'
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;
