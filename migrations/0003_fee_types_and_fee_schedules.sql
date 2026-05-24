-- Production fee setup for reusable fee types, period fee schedules, schedule
-- line items, and audited per-student adjustments.

CREATE TABLE IF NOT EXISTS fee_types (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	label_vi text NOT NULL,
	label_en text NOT NULL,
	category text NOT NULL DEFAULT 'standard',
	default_display_order integer NOT NULL DEFAULT 0,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT fee_types_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT fee_types_code_format CHECK (code ~ '^[a-z][a-z0-9_:-]*$'),
	CONSTRAINT fee_types_label_vi_not_blank CHECK (btrim(label_vi) <> ''),
	CONSTRAINT fee_types_label_en_not_blank CHECK (btrim(label_en) <> ''),
	CONSTRAINT fee_types_status_check CHECK (status IN ('active', 'archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS fee_types_code_key ON fee_types (code);
CREATE INDEX IF NOT EXISTS fee_types_status_idx ON fee_types (status);
CREATE INDEX IF NOT EXISTS fee_types_display_order_idx ON fee_types (default_display_order, code);

CREATE TABLE IF NOT EXISTS fee_schedules (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	school_year_id uuid NOT NULL REFERENCES school_years(id) ON DELETE RESTRICT,
	scope_type text NOT NULL DEFAULT 'school_year',
	class_id uuid REFERENCES classes(id) ON DELETE RESTRICT,
	grade text NOT NULL DEFAULT '',
	period_code text NOT NULL,
	month smallint,
	name text NOT NULL DEFAULT '',
	notes text NOT NULL DEFAULT '',
	status text NOT NULL DEFAULT 'draft',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT fee_schedules_period_not_blank CHECK (btrim(period_code) <> ''),
	CONSTRAINT fee_schedules_status_check CHECK (status IN ('draft', 'active', 'archived')),
	CONSTRAINT fee_schedules_month_check CHECK (month IS NULL OR (month >= 1 AND month <= 12)),
	CONSTRAINT fee_schedules_scope_check CHECK (
		(scope_type = 'school_year' AND class_id IS NULL AND btrim(grade) = '') OR
		(scope_type = 'grade' AND class_id IS NULL AND btrim(grade) <> '') OR
		(scope_type = 'class' AND class_id IS NOT NULL)
	)
);

CREATE INDEX IF NOT EXISTS fee_schedules_school_year_idx ON fee_schedules (school_year_id);
CREATE INDEX IF NOT EXISTS fee_schedules_class_idx ON fee_schedules (class_id);
CREATE INDEX IF NOT EXISTS fee_schedules_grade_idx ON fee_schedules (grade);
CREATE INDEX IF NOT EXISTS fee_schedules_status_idx ON fee_schedules (status);
CREATE UNIQUE INDEX IF NOT EXISTS fee_schedules_active_scope_period_key
	ON fee_schedules (
		school_year_id,
		scope_type,
		COALESCE(class_id::text, ''),
		lower(grade),
		lower(period_code),
		COALESCE(month, 0)
	)
	WHERE status <> 'archived';

CREATE TABLE IF NOT EXISTS fee_schedule_items (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	schedule_id uuid NOT NULL REFERENCES fee_schedules(id) ON DELETE CASCADE,
	fee_type_id uuid NOT NULL REFERENCES fee_types(id) ON DELETE RESTRICT,
	label_vi text NOT NULL,
	label_en text NOT NULL,
	amount integer NOT NULL DEFAULT 0,
	display_order integer NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT fee_schedule_items_label_vi_not_blank CHECK (btrim(label_vi) <> ''),
	CONSTRAINT fee_schedule_items_label_en_not_blank CHECK (btrim(label_en) <> ''),
	CONSTRAINT fee_schedule_items_amount_check CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS fee_schedule_items_schedule_idx ON fee_schedule_items (schedule_id, display_order, id);
CREATE INDEX IF NOT EXISTS fee_schedule_items_fee_type_idx ON fee_schedule_items (fee_type_id);

CREATE TABLE IF NOT EXISTS student_fee_adjustments (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	schedule_id uuid NOT NULL REFERENCES fee_schedules(id) ON DELETE CASCADE,
	student_id uuid NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
	fee_type_id uuid REFERENCES fee_types(id) ON DELETE RESTRICT,
	adjustment_type text NOT NULL,
	label_vi text NOT NULL,
	label_en text NOT NULL,
	amount integer NOT NULL DEFAULT 0,
	reason text NOT NULL,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT student_fee_adjustments_type_check CHECK (adjustment_type IN ('discount', 'surcharge', 'waiver', 'carry_over')),
	CONSTRAINT student_fee_adjustments_label_vi_not_blank CHECK (btrim(label_vi) <> ''),
	CONSTRAINT student_fee_adjustments_label_en_not_blank CHECK (btrim(label_en) <> ''),
	CONSTRAINT student_fee_adjustments_amount_check CHECK (amount >= 0),
	CONSTRAINT student_fee_adjustments_reason_not_blank CHECK (btrim(reason) <> ''),
	CONSTRAINT student_fee_adjustments_status_check CHECK (status IN ('active', 'void'))
);

CREATE INDEX IF NOT EXISTS student_fee_adjustments_schedule_idx ON student_fee_adjustments (schedule_id);
CREATE INDEX IF NOT EXISTS student_fee_adjustments_student_idx ON student_fee_adjustments (student_id);
CREATE INDEX IF NOT EXISTS student_fee_adjustments_type_idx ON student_fee_adjustments (adjustment_type);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fee_types_created_by_user_id_fkey') THEN
		ALTER TABLE fee_types ADD CONSTRAINT fee_types_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fee_types_updated_by_user_id_fkey') THEN
		ALTER TABLE fee_types ADD CONSTRAINT fee_types_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fee_schedules_created_by_user_id_fkey') THEN
		ALTER TABLE fee_schedules ADD CONSTRAINT fee_schedules_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fee_schedules_updated_by_user_id_fkey') THEN
		ALTER TABLE fee_schedules ADD CONSTRAINT fee_schedules_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fee_schedule_items_created_by_user_id_fkey') THEN
		ALTER TABLE fee_schedule_items ADD CONSTRAINT fee_schedule_items_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fee_schedule_items_updated_by_user_id_fkey') THEN
		ALTER TABLE fee_schedule_items ADD CONSTRAINT fee_schedule_items_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'student_fee_adjustments_created_by_user_id_fkey') THEN
		ALTER TABLE student_fee_adjustments ADD CONSTRAINT student_fee_adjustments_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'student_fee_adjustments_updated_by_user_id_fkey') THEN
		ALTER TABLE student_fee_adjustments ADD CONSTRAINT student_fee_adjustments_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS fee_types_set_updated_at ON fee_types;
CREATE TRIGGER fee_types_set_updated_at
	BEFORE UPDATE ON fee_types
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS fee_schedules_set_updated_at ON fee_schedules;
CREATE TRIGGER fee_schedules_set_updated_at
	BEFORE UPDATE ON fee_schedules
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS fee_schedule_items_set_updated_at ON fee_schedule_items;
CREATE TRIGGER fee_schedule_items_set_updated_at
	BEFORE UPDATE ON fee_schedule_items
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS student_fee_adjustments_set_updated_at ON student_fee_adjustments;
CREATE TRIGGER student_fee_adjustments_set_updated_at
	BEFORE UPDATE ON student_fee_adjustments
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

INSERT INTO fee_types (code, label_vi, label_en, category, default_display_order)
VALUES
	('tuition', 'Học phí', 'Tuition', 'tuition', 10),
	('lunch', 'Tiền ăn', 'Lunch', 'meal', 20),
	('shuttle', 'Phí xe đưa rước', 'Shuttle', 'transport', 30),
	('uniform', 'Đồng phục', 'Uniform', 'merchandise', 40),
	('insurance', 'Bảo hiểm', 'Insurance', 'insurance', 50),
	('materials', 'Học liệu', 'Learning materials', 'materials', 60),
	('previous_fees', 'Phí kỳ trước', 'Previous fees', 'carry_over', 70),
	('custom', 'Khoản phí khác', 'Custom fee', 'custom', 100)
ON CONFLICT (code) DO UPDATE
SET label_vi = EXCLUDED.label_vi,
	label_en = EXCLUDED.label_en,
	category = EXCLUDED.category,
	default_display_order = EXCLUDED.default_display_order,
	status = 'active',
	updated_at = now();

INSERT INTO app_permissions (code, description)
VALUES
	('fee_schedules.read', 'Read fee types, period fee schedules, and fee previews'),
	('fee_schedules.write', 'Create and update period fee schedules and student adjustments')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('fee_schedules.read', 'fee_schedules.write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('fee_schedules.read', 'fee_schedules.write')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code = 'fee_schedules.read'
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;
