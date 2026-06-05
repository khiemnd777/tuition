package main

import (
	"testing"
	"time"
)

func TestShouldGenerateSubscriptionRenewalWithinLeadWindow(t *testing.T) {
	row := subscriptionFinanceConsoleRow{
		SubscriptionStatus: subscriptionStatusActive,
		AutoRenew:          true,
		RenewalMode:        "auto_generate",
		NextPeriodStartsAt: "2026-07-01",
	}
	config := subscriptionBillingConfig{
		AutoRenew:       true,
		RenewalMode:     "auto_generate",
		RenewalLeadDays: 7,
	}
	if !shouldGenerateSubscriptionRenewal(row, config, time.Date(2026, time.June, 25, 9, 0, 0, 0, time.UTC)) {
		t.Fatal("expected renewal generation inside lead window")
	}
	if shouldGenerateSubscriptionRenewal(row, config, time.Date(2026, time.June, 20, 9, 0, 0, 0, time.UTC)) {
		t.Fatal("expected renewal generation to stay blocked before lead window")
	}
}

func TestShouldGenerateSubscriptionRenewalSkipsSuspendedTenant(t *testing.T) {
	row := subscriptionFinanceConsoleRow{
		SubscriptionStatus: subscriptionStatusSuspended,
		AutoRenew:          true,
		RenewalMode:        "auto_generate",
		NextPeriodStartsAt: "2026-07-01",
	}
	config := subscriptionBillingConfig{
		AutoRenew:       true,
		RenewalMode:     "auto_generate",
		RenewalLeadDays: 7,
	}
	if shouldGenerateSubscriptionRenewal(row, config, time.Date(2026, time.June, 29, 9, 0, 0, 0, time.UTC)) {
		t.Fatal("expected suspended tenant to skip renewal generation")
	}
}

func TestSubscriptionAutomationStatusFromSummary(t *testing.T) {
	if got := subscriptionAutomationStatusFromSummary(subscriptionAutomationRunSummary{}, true); got != subscriptionAutomationStatusDryRun {
		t.Fatalf("expected dry_run status, got %q", got)
	}
	if got := subscriptionAutomationStatusFromSummary(subscriptionAutomationRunSummary{ErrorCount: 1}, false); got != subscriptionAutomationStatusError {
		t.Fatalf("expected error status, got %q", got)
	}
	if got := subscriptionAutomationStatusFromSummary(subscriptionAutomationRunSummary{ErrorCount: 1, DunningSentCount: 2}, false); got != subscriptionAutomationStatusPartial {
		t.Fatalf("expected partial status, got %q", got)
	}
}

func TestScopeValueForResponse(t *testing.T) {
	if got := scopeValueForResponse("all"); got != "all" {
		t.Fatalf("expected all scope, got %q", got)
	}
	if got := scopeValueForResponse("tenant-id"); got != "active" {
		t.Fatalf("expected active scope, got %q", got)
	}
}
