-- Access/refresh token authentication for production Web Admin sessions.
-- Tokens are opaque browser credentials; only SHA-256 hashes are stored.

ALTER TABLE app_users
	ADD COLUMN IF NOT EXISTS password_hash text NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS password_updated_at timestamptz;

CREATE TABLE IF NOT EXISTS app_auth_sessions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	expires_at timestamptz NOT NULL,
	last_used_at timestamptz,
	revoked_at timestamptz,
	revoked_reason text NOT NULL DEFAULT '',
	ip_address inet,
	user_agent text NOT NULL DEFAULT '',
	CONSTRAINT app_auth_sessions_expiry_check CHECK (expires_at > created_at)
);

CREATE INDEX IF NOT EXISTS app_auth_sessions_user_id_idx ON app_auth_sessions (user_id);
CREATE INDEX IF NOT EXISTS app_auth_sessions_active_idx ON app_auth_sessions (user_id, expires_at)
	WHERE revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS app_auth_access_tokens (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	session_id uuid NOT NULL,
	user_id uuid NOT NULL,
	token_hash text NOT NULL,
	issued_at timestamptz NOT NULL DEFAULT now(),
	expires_at timestamptz NOT NULL,
	last_used_at timestamptz,
	revoked_at timestamptz,
	CONSTRAINT app_auth_access_tokens_hash_not_blank CHECK (btrim(token_hash) <> ''),
	CONSTRAINT app_auth_access_tokens_expiry_check CHECK (expires_at > issued_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS app_auth_access_tokens_hash_key ON app_auth_access_tokens (token_hash);
CREATE INDEX IF NOT EXISTS app_auth_access_tokens_session_id_idx ON app_auth_access_tokens (session_id);
CREATE INDEX IF NOT EXISTS app_auth_access_tokens_active_idx ON app_auth_access_tokens (token_hash, expires_at)
	WHERE revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS app_auth_refresh_tokens (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	session_id uuid NOT NULL,
	user_id uuid NOT NULL,
	token_hash text NOT NULL,
	issued_at timestamptz NOT NULL DEFAULT now(),
	expires_at timestamptz NOT NULL,
	used_at timestamptz,
	revoked_at timestamptz,
	replaced_by_token_id uuid,
	CONSTRAINT app_auth_refresh_tokens_hash_not_blank CHECK (btrim(token_hash) <> ''),
	CONSTRAINT app_auth_refresh_tokens_expiry_check CHECK (expires_at > issued_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS app_auth_refresh_tokens_hash_key ON app_auth_refresh_tokens (token_hash);
CREATE INDEX IF NOT EXISTS app_auth_refresh_tokens_session_id_idx ON app_auth_refresh_tokens (session_id);
CREATE INDEX IF NOT EXISTS app_auth_refresh_tokens_active_idx ON app_auth_refresh_tokens (token_hash, expires_at)
	WHERE used_at IS NULL AND revoked_at IS NULL;

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_sessions_user_id_fkey') THEN
		ALTER TABLE app_auth_sessions ADD CONSTRAINT app_auth_sessions_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES app_users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_access_tokens_session_id_fkey') THEN
		ALTER TABLE app_auth_access_tokens ADD CONSTRAINT app_auth_access_tokens_session_id_fkey
			FOREIGN KEY (session_id) REFERENCES app_auth_sessions(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_access_tokens_user_id_fkey') THEN
		ALTER TABLE app_auth_access_tokens ADD CONSTRAINT app_auth_access_tokens_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES app_users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_refresh_tokens_session_id_fkey') THEN
		ALTER TABLE app_auth_refresh_tokens ADD CONSTRAINT app_auth_refresh_tokens_session_id_fkey
			FOREIGN KEY (session_id) REFERENCES app_auth_sessions(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_refresh_tokens_user_id_fkey') THEN
		ALTER TABLE app_auth_refresh_tokens ADD CONSTRAINT app_auth_refresh_tokens_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES app_users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_auth_refresh_tokens_replaced_by_token_id_fkey') THEN
		ALTER TABLE app_auth_refresh_tokens ADD CONSTRAINT app_auth_refresh_tokens_replaced_by_token_id_fkey
			FOREIGN KEY (replaced_by_token_id) REFERENCES app_auth_refresh_tokens(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS app_auth_sessions_set_updated_at ON app_auth_sessions;
CREATE TRIGGER app_auth_sessions_set_updated_at
	BEFORE UPDATE ON app_auth_sessions
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();
