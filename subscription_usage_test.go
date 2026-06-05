package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseSubscriptionLimitValue(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  int
	}{
		{name: "int", value: 5, want: 5},
		{name: "float", value: float64(10), want: 10},
		{name: "json number", value: json.Number("42"), want: 42},
		{name: "string", value: "7", want: 7},
		{name: "invalid", value: struct{}{}, want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseSubscriptionLimitValue(tc.value); got != tc.want {
				t.Fatalf("parseSubscriptionLimitValue(%v) = %d, want %d", tc.value, got, tc.want)
			}
		})
	}
}

func TestSubscriptionUsagePeriodKey(t *testing.T) {
	now := time.Date(2026, time.June, 5, 10, 0, 0, 0, time.FixedZone("+07", 7*3600))
	if got := subscriptionUsagePeriodKey(subscriptionMetricMonthlyNotifications, now); got != "2026-06" {
		t.Fatalf("expected monthly notification period key 2026-06, got %q", got)
	}
	if got := subscriptionUsagePeriodKey(subscriptionMetricStudents, now); got != "" {
		t.Fatalf("expected non-period metric to use empty period key, got %q", got)
	}
}

func TestTenantUsageLimitErrorMessage(t *testing.T) {
	err := (&tenantUsageLimitError{
		Label:     "Notifications / month",
		Limit:     500,
		Used:      500,
		Requested: 2,
		PeriodKey: "2026-06",
	}).Error()
	const want = "subscription limit reached for Notifications / month for 2026-06 (500/500 used)"
	if err != want {
		t.Fatalf("unexpected error message %q, want %q", err, want)
	}
}
