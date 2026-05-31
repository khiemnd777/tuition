-- School tree management for organizing existing school years, classes,
-- students, parents, fee schedules, and student fee adjustments under schools.

CREATE TABLE IF NOT EXISTS schools (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	name text NOT NULL,
	status text NOT NULL DEFAULT 'active',
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT schools_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT schools_name_not_blank CHECK (btrim(name) <> ''),
	CONSTRAINT schools_status_check CHECK (status IN ('active', 'archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS schools_code_key ON schools (code);
CREATE INDEX IF NOT EXISTS schools_status_idx ON schools (status);

INSERT INTO schools (code, name, status)
VALUES ('ABC_SUN', 'ABC SUN', 'active')
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	status = EXCLUDED.status,
	updated_at = now();

ALTER TABLE school_years ADD COLUMN IF NOT EXISTS school_id uuid;

UPDATE school_years
SET school_id = (SELECT id FROM schools WHERE code = 'ABC_SUN')
WHERE school_id IS NULL;

ALTER TABLE school_years ALTER COLUMN school_id SET NOT NULL;
ALTER TABLE school_years DROP CONSTRAINT IF EXISTS school_years_code_key;

CREATE INDEX IF NOT EXISTS school_years_school_id_idx ON school_years (school_id);
CREATE UNIQUE INDEX IF NOT EXISTS school_years_school_code_key ON school_years (school_id, code);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'schools_created_by_user_id_fkey') THEN
		ALTER TABLE schools ADD CONSTRAINT schools_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'schools_updated_by_user_id_fkey') THEN
		ALTER TABLE schools ADD CONSTRAINT schools_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'school_years_school_id_fkey') THEN
		ALTER TABLE school_years ADD CONSTRAINT school_years_school_id_fkey
			FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE RESTRICT;
	END IF;
END $$;

DROP TRIGGER IF EXISTS schools_set_updated_at ON schools;
CREATE TRIGGER schools_set_updated_at
	BEFORE UPDATE ON schools
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();
