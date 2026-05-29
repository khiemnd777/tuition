-- Payment provider adapters, transaction ledger, and invoice reconciliation.
-- Provider events are stored before normalized transactions are derived so
-- webhook retries and raw provider payloads remain auditable.

CREATE TABLE IF NOT EXISTS payment_providers (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	code text NOT NULL,
	display_name text NOT NULL,
	provider_type text NOT NULL,
	status text NOT NULL DEFAULT 'active',
	config jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT payment_providers_code_not_blank CHECK (btrim(code) <> ''),
	CONSTRAINT payment_providers_code_format CHECK (code ~ '^[a-z][a-z0-9_:-]*$'),
	CONSTRAINT payment_providers_display_name_not_blank CHECK (btrim(display_name) <> ''),
	CONSTRAINT payment_providers_type_check CHECK (provider_type IN ('manual_vietqr', 'bank_webhook', 'payment_link')),
	CONSTRAINT payment_providers_status_check CHECK (status IN ('active', 'inactive', 'sandbox')),
	CONSTRAINT payment_providers_config_object CHECK (jsonb_typeof(config) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS payment_providers_code_key ON payment_providers (code);
CREATE INDEX IF NOT EXISTS payment_providers_status_idx ON payment_providers (status);

CREATE TABLE IF NOT EXISTS payment_intents (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
	provider_id uuid NOT NULL REFERENCES payment_providers(id) ON DELETE RESTRICT,
	intent_code text NOT NULL,
	status text NOT NULL DEFAULT 'pending',
	amount integer NOT NULL,
	currency text NOT NULL DEFAULT 'VND',
	provider_reference text NOT NULL DEFAULT '',
	payment_url text NOT NULL DEFAULT '',
	qr_payload text NOT NULL DEFAULT '',
	provider_request jsonb NOT NULL DEFAULT '{}'::jsonb,
	provider_response jsonb NOT NULL DEFAULT '{}'::jsonb,
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT payment_intents_code_not_blank CHECK (btrim(intent_code) <> ''),
	CONSTRAINT payment_intents_status_check CHECK (status IN ('pending', 'active', 'completed', 'failed', 'cancelled', 'expired')),
	CONSTRAINT payment_intents_amount_check CHECK (amount >= 0),
	CONSTRAINT payment_intents_currency_not_blank CHECK (btrim(currency) <> ''),
	CONSTRAINT payment_intents_request_object CHECK (jsonb_typeof(provider_request) = 'object'),
	CONSTRAINT payment_intents_response_object CHECK (jsonb_typeof(provider_response) = 'object'),
	CONSTRAINT payment_intents_metadata_object CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS payment_intents_code_key ON payment_intents (intent_code);
CREATE UNIQUE INDEX IF NOT EXISTS payment_intents_invoice_provider_active_key
	ON payment_intents (invoice_id, provider_id)
	WHERE status NOT IN ('cancelled', 'expired', 'failed');
CREATE INDEX IF NOT EXISTS payment_intents_invoice_idx ON payment_intents (invoice_id, created_at DESC);
CREATE INDEX IF NOT EXISTS payment_intents_provider_reference_idx ON payment_intents (provider_id, provider_reference)
	WHERE provider_reference <> '';
CREATE INDEX IF NOT EXISTS payment_intents_status_idx ON payment_intents (status);

CREATE TABLE IF NOT EXISTS provider_events (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	provider_id uuid NOT NULL REFERENCES payment_providers(id) ON DELETE RESTRICT,
	provider_event_id text NOT NULL DEFAULT '',
	payload_hash text NOT NULL,
	raw_payload jsonb NOT NULL,
	headers jsonb NOT NULL DEFAULT '{}'::jsonb,
	normalized_payload jsonb NOT NULL DEFAULT '{}'::jsonb,
	status text NOT NULL DEFAULT 'received',
	error_message text NOT NULL DEFAULT '',
	received_at timestamptz NOT NULL DEFAULT now(),
	processed_at timestamptz,
	CONSTRAINT provider_events_payload_hash_not_blank CHECK (btrim(payload_hash) <> ''),
	CONSTRAINT provider_events_status_check CHECK (status IN ('received', 'processed', 'duplicate', 'invalid', 'ignored')),
	CONSTRAINT provider_events_raw_payload_object CHECK (jsonb_typeof(raw_payload) = 'object'),
	CONSTRAINT provider_events_headers_object CHECK (jsonb_typeof(headers) = 'object'),
	CONSTRAINT provider_events_normalized_payload_object CHECK (jsonb_typeof(normalized_payload) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS provider_events_payload_hash_key ON provider_events (provider_id, payload_hash);
CREATE UNIQUE INDEX IF NOT EXISTS provider_events_provider_event_id_key
	ON provider_events (provider_id, provider_event_id)
	WHERE provider_event_id <> '';
CREATE INDEX IF NOT EXISTS provider_events_received_at_idx ON provider_events (received_at DESC);
CREATE INDEX IF NOT EXISTS provider_events_status_idx ON provider_events (status);

CREATE TABLE IF NOT EXISTS payment_transactions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	provider_id uuid NOT NULL REFERENCES payment_providers(id) ON DELETE RESTRICT,
	provider_event_id uuid REFERENCES provider_events(id) ON DELETE SET NULL,
	payment_intent_id uuid REFERENCES payment_intents(id) ON DELETE SET NULL,
	invoice_id uuid REFERENCES invoices(id) ON DELETE SET NULL,
	provider_transaction_id text NOT NULL DEFAULT '',
	direction text NOT NULL DEFAULT 'in',
	amount integer NOT NULL,
	currency text NOT NULL DEFAULT 'VND',
	transaction_time timestamptz NOT NULL DEFAULT now(),
	account_number text NOT NULL DEFAULT '',
	account_name text NOT NULL DEFAULT '',
	bank_bin text NOT NULL DEFAULT '',
	bank_name text NOT NULL DEFAULT '',
	description text NOT NULL DEFAULT '',
	reference_code text NOT NULL DEFAULT '',
	status text NOT NULL DEFAULT 'unmatched',
	raw_payload jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	updated_by_user_id uuid,
	CONSTRAINT payment_transactions_direction_check CHECK (direction IN ('in', 'out')),
	CONSTRAINT payment_transactions_amount_check CHECK (amount >= 0),
	CONSTRAINT payment_transactions_currency_not_blank CHECK (btrim(currency) <> ''),
	CONSTRAINT payment_transactions_status_check CHECK (status IN ('unmatched', 'matched', 'manual_review', 'ignored', 'reversed')),
	CONSTRAINT payment_transactions_raw_payload_object CHECK (jsonb_typeof(raw_payload) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS payment_transactions_provider_txn_key
	ON payment_transactions (provider_id, provider_transaction_id)
	WHERE provider_transaction_id <> '';
CREATE INDEX IF NOT EXISTS payment_transactions_invoice_idx ON payment_transactions (invoice_id);
CREATE INDEX IF NOT EXISTS payment_transactions_time_idx ON payment_transactions (transaction_time DESC);
CREATE INDEX IF NOT EXISTS payment_transactions_status_idx ON payment_transactions (status);
CREATE INDEX IF NOT EXISTS payment_transactions_account_idx ON payment_transactions (account_number)
	WHERE account_number <> '';

CREATE TABLE IF NOT EXISTS reconciliation_matches (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
	transaction_id uuid NOT NULL REFERENCES payment_transactions(id) ON DELETE RESTRICT,
	match_type text NOT NULL,
	status text NOT NULL DEFAULT 'matched',
	score integer NOT NULL DEFAULT 0,
	amount_applied integer NOT NULL,
	reason text NOT NULL DEFAULT '',
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	CONSTRAINT reconciliation_matches_type_check CHECK (match_type IN ('auto', 'manual', 'cash', 'provider_reference')),
	CONSTRAINT reconciliation_matches_status_check CHECK (status IN ('matched', 'manual_review', 'reversed')),
	CONSTRAINT reconciliation_matches_score_check CHECK (score >= 0 AND score <= 100),
	CONSTRAINT reconciliation_matches_amount_check CHECK (amount_applied >= 0),
	CONSTRAINT reconciliation_matches_metadata_object CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS reconciliation_matches_transaction_invoice_key
	ON reconciliation_matches (transaction_id, invoice_id)
	WHERE status <> 'reversed';
CREATE INDEX IF NOT EXISTS reconciliation_matches_invoice_idx ON reconciliation_matches (invoice_id, created_at DESC);

CREATE TABLE IF NOT EXISTS manual_cash_receipts (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
	payment_transaction_id uuid NOT NULL REFERENCES payment_transactions(id) ON DELETE RESTRICT,
	collector_user_id uuid,
	collector_name text NOT NULL,
	amount integer NOT NULL,
	currency text NOT NULL DEFAULT 'VND',
	paid_at timestamptz NOT NULL DEFAULT now(),
	receipt_reference text NOT NULL,
	note text NOT NULL DEFAULT '',
	created_at timestamptz NOT NULL DEFAULT now(),
	created_by_user_id uuid,
	CONSTRAINT manual_cash_receipts_collector_not_blank CHECK (btrim(collector_name) <> ''),
	CONSTRAINT manual_cash_receipts_amount_check CHECK (amount > 0),
	CONSTRAINT manual_cash_receipts_currency_not_blank CHECK (btrim(currency) <> ''),
	CONSTRAINT manual_cash_receipts_reference_not_blank CHECK (btrim(receipt_reference) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS manual_cash_receipts_reference_key ON manual_cash_receipts (receipt_reference);
CREATE INDEX IF NOT EXISTS manual_cash_receipts_invoice_idx ON manual_cash_receipts (invoice_id, paid_at DESC);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_providers_created_by_user_id_fkey') THEN
		ALTER TABLE payment_providers ADD CONSTRAINT payment_providers_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_providers_updated_by_user_id_fkey') THEN
		ALTER TABLE payment_providers ADD CONSTRAINT payment_providers_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_intents_created_by_user_id_fkey') THEN
		ALTER TABLE payment_intents ADD CONSTRAINT payment_intents_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_intents_updated_by_user_id_fkey') THEN
		ALTER TABLE payment_intents ADD CONSTRAINT payment_intents_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_transactions_created_by_user_id_fkey') THEN
		ALTER TABLE payment_transactions ADD CONSTRAINT payment_transactions_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_transactions_updated_by_user_id_fkey') THEN
		ALTER TABLE payment_transactions ADD CONSTRAINT payment_transactions_updated_by_user_id_fkey
			FOREIGN KEY (updated_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'reconciliation_matches_created_by_user_id_fkey') THEN
		ALTER TABLE reconciliation_matches ADD CONSTRAINT reconciliation_matches_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'manual_cash_receipts_collector_user_id_fkey') THEN
		ALTER TABLE manual_cash_receipts ADD CONSTRAINT manual_cash_receipts_collector_user_id_fkey
			FOREIGN KEY (collector_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'manual_cash_receipts_created_by_user_id_fkey') THEN
		ALTER TABLE manual_cash_receipts ADD CONSTRAINT manual_cash_receipts_created_by_user_id_fkey
			FOREIGN KEY (created_by_user_id) REFERENCES app_users(id) ON DELETE SET NULL;
	END IF;
END $$;

DROP TRIGGER IF EXISTS payment_providers_set_updated_at ON payment_providers;
CREATE TRIGGER payment_providers_set_updated_at
	BEFORE UPDATE ON payment_providers
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS payment_intents_set_updated_at ON payment_intents;
CREATE TRIGGER payment_intents_set_updated_at
	BEFORE UPDATE ON payment_intents
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

DROP TRIGGER IF EXISTS payment_transactions_set_updated_at ON payment_transactions;
CREATE TRIGGER payment_transactions_set_updated_at
	BEFORE UPDATE ON payment_transactions
	FOR EACH ROW
	EXECUTE FUNCTION abc_set_updated_at();

INSERT INTO payment_providers (code, display_name, provider_type, status, config)
VALUES
	('manual_vietqr', 'Manual VietQR', 'manual_vietqr', 'active', '{"reconciliation":"manual"}'::jsonb),
	('sepay', 'SePay', 'bank_webhook', 'active', '{"webhook":"bank_transaction"}'::jsonb),
	('payos', 'payOS', 'payment_link', 'active', '{"webhook":"payment_link"}'::jsonb)
ON CONFLICT (code) DO UPDATE
SET display_name = EXCLUDED.display_name,
	provider_type = EXCLUDED.provider_type,
	status = EXCLUDED.status,
	config = payment_providers.config || EXCLUDED.config,
	updated_at = now();

INSERT INTO app_permissions (code, description)
VALUES
	('payments.read', 'Read payment providers, intents, transactions, and reconciliation state'),
	('payments.write', 'Create payment intents and manual cash receipts'),
	('payments.reconcile', 'Process provider webhooks and reconcile transactions to invoices')
ON CONFLICT (code) DO UPDATE
SET description = EXCLUDED.description,
	updated_at = now();

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
CROSS JOIN app_permissions permission
WHERE role.code = 'super_admin'
	AND permission.code IN ('payments.read', 'payments.write', 'payments.reconcile')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN ('payments.read', 'payments.write', 'payments.reconcile')
WHERE role.code = 'billing_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code = 'payments.read'
WHERE role.code = 'viewer'
ON CONFLICT (role_id, permission_id) DO NOTHING;
