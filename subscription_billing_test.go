package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
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
	got := buildSubscriptionInvoiceCode("dekisugi", time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC))
	if got != "SUB-DEKISUGI-202607" {
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

func TestSubscriptionBillingConfigFromProfileDefaults(t *testing.T) {
	got := subscriptionBillingConfigFromProfile(tenantSubscriptionBillingProfile{})
	if got.IntervalMonths != 1 || got.DueDays != 10 || got.RenewalMode != "manual" {
		t.Fatalf("unexpected config defaults %+v", got)
	}
	if got.RenewalLeadDays != 7 || got.DunningIntervalDays != 3 || got.SuspendAfterDays != 14 {
		t.Fatalf("unexpected automation defaults %+v", got)
	}
	if !got.DunningEnabled || !got.SuspendEnabled {
		t.Fatalf("expected automation booleans enabled by default, got %+v", got)
	}
}

func TestValidateSubscriptionBillingConfig(t *testing.T) {
	err := validateSubscriptionBillingConfig(subscriptionBillingConfigSaveInput{
		TenantID:            "tenant-1",
		Amount:              1000,
		IntervalMonths:      1,
		DueDays:             10,
		RenewalMode:         "auto_generate",
		RenewalLeadDays:     7,
		DunningIntervalDays: 3,
		SuspendAfterDays:    14,
	})
	if err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
	if err := validateSubscriptionBillingConfig(subscriptionBillingConfigSaveInput{
		TenantID:            "tenant-1",
		Amount:              1000,
		IntervalMonths:      0,
		DueDays:             10,
		RenewalMode:         "manual",
		RenewalLeadDays:     7,
		DunningIntervalDays: 3,
		SuspendAfterDays:    14,
	}); err == nil {
		t.Fatal("expected interval validation error")
	}
	if err := validateSubscriptionBillingConfig(subscriptionBillingConfigSaveInput{
		TenantID:            "tenant-1",
		Amount:              1000,
		IntervalMonths:      1,
		DueDays:             10,
		RenewalMode:         "manual",
		RenewalLeadDays:     31,
		DunningIntervalDays: 3,
		SuspendAfterDays:    14,
	}); err == nil {
		t.Fatal("expected renewal lead validation error")
	}
}

func TestFilterSubscriptionInvoicesByStatus(t *testing.T) {
	invoices := []subscriptionInvoiceSummary{
		{InvoiceCode: "A", Status: subscriptionInvoiceStatusOpen},
		{InvoiceCode: "B", Status: subscriptionInvoiceStatusPastDue},
		{InvoiceCode: "C", Status: subscriptionInvoiceStatusPaid},
	}
	got := filterSubscriptionInvoicesByStatus(invoices, subscriptionInvoiceStatusPaid, subscriptionInvoiceStatusPastDue)
	if len(got) != 2 || got[0].InvoiceCode != "B" || got[1].InvoiceCode != "C" {
		t.Fatalf("unexpected filtered invoices %+v", got)
	}
}

func TestNormalizeSubscriptionBatchRunInput(t *testing.T) {
	got := normalizeSubscriptionBatchRunInput(subscriptionBatchRunInput{Scope: "ALL", DryRun: false, ConfirmRun: false})
	if got.Scope != "all" || !got.DryRun {
		t.Fatalf("unexpected normalized batch input %+v", got)
	}
}

func TestResolveSubscriptionFinanceScopeDefaultsToAllForPlatformOnlySession(t *testing.T) {
	rec := httptest.NewRecorder()
	user := authenticatedUser{
		IsPlatformAdmin: true,
		PermissionSet: map[string]bool{
			"tenant.view": true,
		},
	}
	got, ok := resolveSubscriptionFinanceScopeFromValue(rec, user, "", "")
	if !ok {
		t.Fatalf("expected platform-only scope resolution to pass, got %d %s", rec.Code, rec.Body.String())
	}
	if got != "all" {
		t.Fatalf("expected all scope, got %q", got)
	}
}

func TestResolveSubscriptionFinanceScopeRejectsActiveWithoutTenant(t *testing.T) {
	rec := httptest.NewRecorder()
	user := authenticatedUser{
		IsPlatformAdmin: true,
		PermissionSet: map[string]bool{
			"tenant.view": true,
		},
	}
	_, ok := resolveSubscriptionFinanceScopeFromValue(rec, user, "", "active")
	if ok {
		t.Fatal("expected active scope without tenant to be rejected")
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden, got %d", rec.Code)
	}
}
