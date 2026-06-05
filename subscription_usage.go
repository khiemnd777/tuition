package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	subscriptionMetricSchools              = "schools"
	subscriptionMetricOperators            = "operators"
	subscriptionMetricStudents             = "students"
	subscriptionMetricMonthlyNotifications = "monthly_notifications"
)

type tenantUsageMetricSummary struct {
	MetricCode string `json:"metricCode"`
	Label      string `json:"label"`
	Limit      int    `json:"limit"`
	Used       int    `json:"used"`
	Remaining  int    `json:"remaining"`
	PeriodKey  string `json:"periodKey,omitempty"`
	Unlimited  bool   `json:"unlimited,omitempty"`
	Exhausted  bool   `json:"exhausted,omitempty"`
}

type tenantUsageLimitError struct {
	MetricCode string
	Label      string
	Limit      int
	Used       int
	Requested  int
	PeriodKey  string
}

func (err *tenantUsageLimitError) Error() string {
	if err == nil {
		return "subscription limit reached"
	}
	scope := ""
	if strings.TrimSpace(err.PeriodKey) != "" {
		scope = " for " + strings.TrimSpace(err.PeriodKey)
	}
	return fmt.Sprintf("subscription limit reached for %s%s (%d/%d used)", err.Label, scope, err.Used, err.Limit)
}

func subscriptionUsageMetricCodes() []string {
	return []string{
		subscriptionMetricSchools,
		subscriptionMetricOperators,
		subscriptionMetricStudents,
		subscriptionMetricMonthlyNotifications,
	}
}

func subscriptionUsageMetricLabel(metricCode string) string {
	switch strings.TrimSpace(metricCode) {
	case subscriptionMetricSchools:
		return "Schools"
	case subscriptionMetricOperators:
		return "Operators"
	case subscriptionMetricStudents:
		return "Students"
	case subscriptionMetricMonthlyNotifications:
		return "Notifications / month"
	default:
		return strings.TrimSpace(metricCode)
	}
}

func subscriptionUsagePeriodKey(metricCode string, now time.Time) string {
	if strings.TrimSpace(metricCode) == subscriptionMetricMonthlyNotifications {
		return now.UTC().Format("2006-01")
	}
	return ""
}

func parseSubscriptionLimitValue(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	case json.Number:
		number, _ := typed.Int64()
		return int(number)
	case string:
		number, _ := json.Number(strings.TrimSpace(typed)).Int64()
		return int(number)
	default:
		return 0
	}
}

func loadTenantUsageSummaries(ctx context.Context, db *sql.DB, tenantID string, now time.Time) ([]tenantUsageMetricSummary, error) {
	metrics := make([]tenantUsageMetricSummary, 0, len(subscriptionUsageMetricCodes()))
	for _, metricCode := range subscriptionUsageMetricCodes() {
		limit, limited, err := loadTenantPlanLimit(ctx, db, tenantID, metricCode)
		if err != nil {
			return nil, err
		}
		used, err := loadOrRebuildTenantUsageCounter(ctx, db, tenantID, metricCode, now)
		if err != nil {
			return nil, err
		}
		item := tenantUsageMetricSummary{
			MetricCode: metricCode,
			Label:      subscriptionUsageMetricLabel(metricCode),
			Limit:      limit,
			Used:       used,
			PeriodKey:  subscriptionUsagePeriodKey(metricCode, now),
			Unlimited:  !limited,
		}
		if limited {
			item.Remaining = limit - used
			if item.Remaining < 0 {
				item.Remaining = 0
			}
			item.Exhausted = used >= limit
		}
		metrics = append(metrics, item)
	}
	return metrics, nil
}

func loadTenantPlanLimit(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string) (int, bool, error) {
	var limitsBytes []byte
	err := exec.QueryRowContext(ctx, `
SELECT COALESCE(plan.limits, '{}'::jsonb)
FROM tenant_subscriptions ts
JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE ts.tenant_id = $1::uuid`, tenantID).Scan(&limitsBytes)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	limits := decodeMetadata(limitsBytes)
	limit := parseSubscriptionLimitValue(limits[metricCode])
	if limit <= 0 {
		return 0, false, nil
	}
	return limit, true, nil
}

func loadOrRebuildTenantUsageCounter(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, now time.Time) (int, error) {
	periodKey := subscriptionUsagePeriodKey(metricCode, now)
	count, ok, err := loadTenantUsageCounter(ctx, exec, tenantID, metricCode, periodKey)
	if err != nil {
		return 0, err
	}
	if ok {
		return count, nil
	}
	count, err = computeTenantUsageCount(ctx, exec, tenantID, metricCode, now)
	if err != nil {
		return 0, err
	}
	if err := upsertTenantUsageCounter(ctx, exec, tenantID, metricCode, periodKey, count); err != nil {
		return 0, err
	}
	return count, nil
}

func loadTenantUsageCounter(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, periodKey string) (int, bool, error) {
	var used int
	err := exec.QueryRowContext(ctx, `
SELECT used_count
FROM tenant_usage_counters
WHERE tenant_id = $1::uuid
	AND metric_code = $2
	AND period_key = $3`, tenantID, metricCode, periodKey).Scan(&used)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return used, true, nil
}

func computeTenantUsageCount(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, now time.Time) (int, error) {
	var query string
	var args []any
	switch strings.TrimSpace(metricCode) {
	case subscriptionMetricSchools:
		query = `SELECT COUNT(*)::integer FROM schools WHERE tenant_id = $1::uuid`
		args = []any{tenantID}
	case subscriptionMetricOperators:
		query = `
SELECT COUNT(DISTINCT tur.user_id)::integer
FROM tenant_user_roles tur
JOIN tenant_memberships tm
	ON tm.tenant_id = tur.tenant_id
	AND tm.user_id = tur.user_id
WHERE tur.tenant_id = $1::uuid
	AND tm.status = 'active'`
		args = []any{tenantID}
	case subscriptionMetricStudents:
		query = `SELECT COUNT(*)::integer FROM students WHERE tenant_id = $1::uuid AND status <> 'inactive'`
		args = []any{tenantID}
	case subscriptionMetricMonthlyNotifications:
		query = `
SELECT COUNT(*)::integer
FROM notification_logs nl
JOIN notification_campaigns c ON c.id = nl.campaign_id
WHERE c.tenant_id = $1::uuid
	AND nl.status = 'sent'
	AND NOT nl.dry_run
	AND to_char(COALESCE(nl.sent_at, nl.created_at) AT TIME ZONE 'UTC', 'YYYY-MM') = $2`
		args = []any{tenantID, subscriptionUsagePeriodKey(metricCode, now)}
	default:
		return 0, nil
	}
	var used int
	if err := exec.QueryRowContext(ctx, query, args...).Scan(&used); err != nil {
		return 0, err
	}
	return used, nil
}

func upsertTenantUsageCounter(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, periodKey string, used int) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO tenant_usage_counters (tenant_id, metric_code, period_key, used_count, updated_at)
VALUES ($1::uuid, $2, $3, $4, now())
ON CONFLICT (tenant_id, metric_code, period_key) DO UPDATE
SET used_count = EXCLUDED.used_count,
	updated_at = now()`, tenantID, metricCode, periodKey, used)
	return err
}

func incrementTenantUsageCounter(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, periodKey string, delta int) error {
	if delta <= 0 {
		return nil
	}
	_, err := exec.ExecContext(ctx, `
INSERT INTO tenant_usage_counters (tenant_id, metric_code, period_key, used_count, updated_at)
VALUES ($1::uuid, $2, $3, $4, now())
ON CONFLICT (tenant_id, metric_code, period_key) DO UPDATE
SET used_count = tenant_usage_counters.used_count + EXCLUDED.used_count,
	updated_at = now()`, tenantID, metricCode, periodKey, delta)
	return err
}

func rebuildTenantUsageCounter(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, now time.Time) error {
	used, err := computeTenantUsageCount(ctx, exec, tenantID, metricCode, now)
	if err != nil {
		return err
	}
	return upsertTenantUsageCounter(ctx, exec, tenantID, metricCode, subscriptionUsagePeriodKey(metricCode, now), used)
}

func enforceTenantUsageLimit(ctx context.Context, exec masterDataExecutor, tenantID string, metricCode string, delta int, now time.Time) error {
	if delta <= 0 {
		return nil
	}
	limit, limited, err := loadTenantPlanLimit(ctx, exec, tenantID, metricCode)
	if err != nil || !limited {
		return err
	}
	used, err := loadOrRebuildTenantUsageCounter(ctx, exec, tenantID, metricCode, now)
	if err != nil {
		return err
	}
	if used+delta > limit {
		return &tenantUsageLimitError{
			MetricCode: metricCode,
			Label:      subscriptionUsageMetricLabel(metricCode),
			Limit:      limit,
			Used:       used,
			Requested:  delta,
			PeriodKey:  subscriptionUsagePeriodKey(metricCode, now),
		}
	}
	return nil
}
