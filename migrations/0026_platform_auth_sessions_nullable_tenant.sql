-- Allow platform-only auth sessions without forcing a tenant binding.

ALTER TABLE app_auth_sessions
	ALTER COLUMN tenant_id DROP NOT NULL;
