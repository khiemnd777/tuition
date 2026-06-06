-- Link payment transactions to subscription invoices for self-serve checkout confirmation.

ALTER TABLE payment_transactions
	ADD COLUMN IF NOT EXISTS subscription_invoice_id uuid;

CREATE INDEX IF NOT EXISTS payment_transactions_subscription_invoice_idx
	ON payment_transactions (subscription_invoice_id);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_transactions_subscription_invoice_id_fkey') THEN
		ALTER TABLE payment_transactions ADD CONSTRAINT payment_transactions_subscription_invoice_id_fkey
			FOREIGN KEY (subscription_invoice_id) REFERENCES subscription_invoices(id) ON DELETE SET NULL;
	END IF;
END $$;
