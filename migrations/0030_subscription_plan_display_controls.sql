ALTER TABLE subscription_plans
	ADD COLUMN IF NOT EXISTS contact_price boolean NOT NULL DEFAULT false,
	ADD COLUMN IF NOT EXISTS display_order integer NOT NULL DEFAULT 100;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'subscription_plans_display_order_non_negative'
	) THEN
		ALTER TABLE subscription_plans
			ADD CONSTRAINT subscription_plans_display_order_non_negative CHECK (display_order >= 0);
	END IF;
END $$;

CREATE INDEX IF NOT EXISTS subscription_plans_display_order_idx
	ON subscription_plans (display_order, code);
