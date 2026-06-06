-- Compatibility shim for legacy trigger migrations that still call set_updated_at().

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN abc_set_updated_at();
END;
$$;
