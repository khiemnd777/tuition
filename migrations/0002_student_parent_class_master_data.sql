-- Production master data for students, parents, school years, and classes.
-- This module intentionally uses student_code as the durable identifier, so
-- duplicate student names are safe and never drive production identity.

CREATE TABLE IF NOT EXISTS school_years (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	name text NOT NULL,
	starts_on date,
	ends_on date,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT school_years_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT school_years_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT school_years_status_check CHECK (status IN ('active', 'archived')),
	CONSTRAINT school_years_date_order CHECK (starts_on IS NULL OR ends_on IS NULL OR starts_on <= ends_on),
	CONSTRAINT school_years_code_key UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS school_years_status_idx ON school_years (status);

CREATE TABLE IF NOT EXISTS classes (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	school_year_id uuid NOT NULL REFERENCES school_years(id) ON DELETE RESTRICT,
	grade text NOT NULL,
	name text NOT NULL,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT classes_grade_not_blank CHECK (btrim(grade) <> ''),
	CONSTRAINT classes_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT classes_status_check CHECK (status IN ('active', 'archived')),
	CONSTRAINT classes_school_year_grade_name_key UNIQUE (school_year_id, grade, name)
);

CREATE INDEX IF NOT EXISTS classes_school_year_idx ON classes (school_year_id);
CREATE INDEX IF NOT EXISTS classes_grade_idx ON classes (grade);
CREATE INDEX IF NOT EXISTS classes_status_idx ON classes (status);

CREATE TABLE IF NOT EXISTS students (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	student_code text NOT NULL,
	full_name text NOT NULL,
	class_id uuid NOT NULL REFERENCES classes(id) ON DELETE RESTRICT,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT students_student_code_not_blank CHECK (btrim(student_code) <> ''),
	CONSTRAINT students_full_name_not_blank CHECK (btrim(full_name) <> ''),
	CONSTRAINT students_status_check CHECK (status IN ('active', 'inactive', 'graduated')),
	CONSTRAINT students_student_code_key UNIQUE (student_code)
);

CREATE INDEX IF NOT EXISTS students_class_id_idx ON students (class_id);
CREATE INDEX IF NOT EXISTS students_status_idx ON students (status);
CREATE INDEX IF NOT EXISTS students_full_name_idx ON students (full_name);

CREATE TABLE IF NOT EXISTS parents (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	full_name text NOT NULL,
	email text NOT NULL DEFAULT '',
	phone text NOT NULL DEFAULT '',
	email_active boolean NOT NULL DEFAULT true,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT parents_full_name_not_blank CHECK (btrim(full_name) <> ''),
	CONSTRAINT parents_email_lowercase CHECK (email = lower(email)),
	CONSTRAINT parents_status_check CHECK (status IN ('active', 'inactive'))
);

CREATE UNIQUE INDEX IF NOT EXISTS parents_email_key ON parents (email) WHERE email <> '';
CREATE INDEX IF NOT EXISTS parents_status_idx ON parents (status);

CREATE TABLE IF NOT EXISTS student_parents (
	student_id uuid NOT NULL REFERENCES students(id) ON DELETE CASCADE,
	parent_id uuid NOT NULL REFERENCES parents(id) ON DELETE RESTRICT,
	relationship text NOT NULL DEFAULT 'guardian',
	is_primary boolean NOT NULL DEFAULT false,
	is_active boolean NOT NULL DEFAULT true,
	receives_billing_email boolean NOT NULL DEFAULT true,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	PRIMARY KEY (student_id, parent_id),
	CONSTRAINT student_parents_relationship_not_blank CHECK (btrim(relationship) <> '')
);

CREATE INDEX IF NOT EXISTS student_parents_parent_id_idx ON student_parents (parent_id);
CREATE INDEX IF NOT EXISTS student_parents_billing_idx ON student_parents (receives_billing_email, is_active);
CREATE UNIQUE INDEX IF NOT EXISTS student_parents_one_primary_per_student_idx
	ON student_parents (student_id)
	WHERE is_primary AND is_active;

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'school_years_created_by_user_id_fkey') THEN
		ALTER TABLE school_years ADD CONSTRAINT school_years_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'school_years_updated_by_user_id_fkey') THEN
		ALTER TABLE school_years ADD CONSTRAINT school_years_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'classes_created_by_user_id_fkey') THEN
		ALTER TABLE classes ADD CONSTRAINT classes_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'classes_updated_by_user_id_fkey') THEN
		ALTER TABLE classes ADD CONSTRAINT classes_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'students_created_by_user_id_fkey') THEN
		ALTER TABLE students ADD CONSTRAINT students_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'students_updated_by_user_id_fkey') THEN
		ALTER TABLE students ADD CONSTRAINT students_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'parents_created_by_user_id_fkey') THEN
		ALTER TABLE parents ADD CONSTRAINT parents_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'parents_updated_by_user_id_fkey') THEN
		ALTER TABLE parents ADD CONSTRAINT parents_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'student_parents_created_by_user_id_fkey') THEN
		ALTER TABLE student_parents ADD CONSTRAINT student_parents_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'student_parents_updated_by_user_id_fkey') THEN
		ALTER TABLE student_parents ADD CONSTRAINT student_parents_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS school_years_set_updated_at ON school_years;
CREATE TRIGGER school_years_set_updated_at
	BEFORE UPDATE ON school_years
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS classes_set_updated_at ON classes;
CREATE TRIGGER classes_set_updated_at
	BEFORE UPDATE ON classes
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS students_set_updated_at ON students;
CREATE TRIGGER students_set_updated_at
	BEFORE UPDATE ON students
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS parents_set_updated_at ON parents;
CREATE TRIGGER parents_set_updated_at
	BEFORE UPDATE ON parents
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS student_parents_set_updated_at ON student_parents;
CREATE TRIGGER student_parents_set_updated_at
	BEFORE UPDATE ON student_parents
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

INSERT INTO app_permissions (code, description)
VALUES
	('master_data.read', 'Read student, parent, class, and school-year master data'),
	('master_data.write', 'Import and maintain student, parent, class, and school-year master data')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('master_data.read', 'master_data.write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('master_data.read', 'master_data.write')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code = 'master_data.read'
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;
