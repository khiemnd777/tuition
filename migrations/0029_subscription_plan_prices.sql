ALTER TABLE subscription_plans
	ADD COLUMN IF NOT EXISTS base_price_vnd integer NOT NULL DEFAULT 0,
	ADD COLUMN IF NOT EXISTS partner_price_vnd integer,
	ADD COLUMN IF NOT EXISTS promotional_price_vnd integer;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'subscription_plans_base_price_non_negative'
	) THEN
		ALTER TABLE subscription_plans
			ADD CONSTRAINT subscription_plans_base_price_non_negative CHECK (base_price_vnd >= 0);
	END IF;

	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'subscription_plans_partner_price_non_negative'
	) THEN
		ALTER TABLE subscription_plans
			ADD CONSTRAINT subscription_plans_partner_price_non_negative CHECK (partner_price_vnd IS NULL OR partner_price_vnd >= 0);
	END IF;

	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'subscription_plans_promotional_price_non_negative'
	) THEN
		ALTER TABLE subscription_plans
			ADD CONSTRAINT subscription_plans_promotional_price_non_negative CHECK (promotional_price_vnd IS NULL OR promotional_price_vnd >= 0);
	END IF;

	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'subscription_plans_promotional_price_not_above_base'
	) THEN
		ALTER TABLE subscription_plans
			ADD CONSTRAINT subscription_plans_promotional_price_not_above_base CHECK (
				promotional_price_vnd IS NULL OR base_price_vnd = 0 OR promotional_price_vnd <= base_price_vnd
			);
	END IF;
END $$;
