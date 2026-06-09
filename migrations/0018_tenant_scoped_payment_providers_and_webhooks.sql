-- Tenant-scoped payment providers and webhook ownership.
-- Keep legacy provider behavior while isolating provider credentials and
-- webhook routing by tenant.

ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS tenant_id uuid;

UPDATE payment_providers
SET tenant_id = (SELECT id FROM tenants WHERE code = 'DEKISUGI')
WHERE tenant_id IS NULL;

ALTER TABLE payment_providers ALTER COLUMN tenant_id SET NOT NULL;

DROP INDEX IF EXISTS payment_providers_code_key;
DROP INDEX IF EXISTS payment_providers_status_idx;

CREATE UNIQUE INDEX IF NOT EXISTS payment_providers_tenant_code_key
	ON payment_providers (tenant_id, code);
CREATE INDEX IF NOT EXISTS payment_providers_tenant_id_idx ON payment_providers (tenant_id);
CREATE INDEX IF NOT EXISTS payment_providers_tenant_status_idx ON payment_providers (tenant_id, status);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_providers_tenant_id_fkey') THEN
		ALTER TABLE payment_providers ADD CONSTRAINT payment_providers_tenant_id_fkey
			FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
	END IF;
END $$;

INSERT INTO payment_providers (tenant_id, code, display_name, provider_type, status, config)
SELECT tenant.id,
	template.code,
	template.display_name,
	template.provider_type,
	template.status,
	template.config
FROM tenants tenant
CROSS JOIN (
	SELECT code, display_name, provider_type, status, config
	FROM payment_providers
	WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'DEKISUGI')
) template
WHERE tenant.code <> 'DEKISUGI'
ON CONFLICT (tenant_id, code) DO UPDATE
SET display_name = EXCLUDED.display_name,
	provider_type = EXCLUDED.provider_type,
	status = EXCLUDED.status,
	config = payment_providers.config || EXCLUDED.config,
	updated_at = now();
