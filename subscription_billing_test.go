package main

import (
	"database/sql"
	"testing"
	"time"
)

func TestSuggestSubscriptionBillingPeriodFromCurrentPeriodEnd(t *testing.T) {
	profile := tenantSubscriptionBillingProfile{
		BillingMetadata: map[string]any{
			"interval_months": 1,
			"due_days":        7,
			"amount":          1500000,
		},
		CurrentPeriodEndsAt: sql.NullTime{Valid: true, Time: time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC)},
	}

	got := suggestSubscriptionBillingPeriod(profile, time.Date(2026, time.June, 5, 0, 0, 0, 0, time.UTC))
	if got.PeriodStartsAt != "2026-07-01" || got.PeriodEndsAt != "2026-07-31" {
		t.Fatalf("unexpected next period %+v", got)
	}
	if got.DueAt != "2026-07-07" {
		t.Fatalf("expected due date 2026-07-07, got %+v", got)
	}
	if got.Amount != 1500000 {
		t.Fatalf("expected amount 1500000, got %+v", got)
	}
}

func TestSuggestSubscriptionBillingPeriodDefaultsFromCurrentMonth(t *testing.T) {
	profile := tenantSubscriptionBillingProfile{}
	got := suggestSubscriptionBillingPeriod(profile, time.Date(2026, time.June, 5, 10, 0, 0, 0, time.UTC))
	if got.PeriodStartsAt != "2026-06-01" || got.PeriodEndsAt != "2026-06-30" {
		t.Fatalf("unexpected default period %+v", got)
	}
	if got.DueAt != "2026-06-10" {
		t.Fatalf("expected default due date 2026-06-10, got %+v", got)
	}
}

func TestBuildSubscriptionInvoiceCode(t *testing.T) {
	got := buildSubscriptionInvoiceCode("abc_sun", time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC))
	if got != "SUB-ABC_SUN-202607" {
		t.Fatalf("unexpected invoice code %q", got)
	}
}

func TestNormalizeSubscriptionDunningRunInputKeepsRealSendOnlyWithConfirm(t *testing.T) {
	got := normalizeSubscriptionDunningRunInput(subscriptionDunningRunInput{
		TenantID:    "tenant-1",
		DryRun:      false,
		ConfirmSend: false,
	}, "tenant-1")
	if !got.DryRun {
		t.Fatalf("expected dry run fallback when confirmSend=false, got %+v", got)
	}
}
