package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	subscriptionPlanFreeTrial = "free_trial"
	subscriptionPlanStandard  = "standard"

	subscriptionStatusTrial     = "trial"
	subscriptionStatusActive    = "active"
	subscriptionStatusPastDue   = "past_due"
	subscriptionStatusSuspended = "suspended"
	subscriptionStatusCancelled = "cancelled"
)

type subscriptionPlanSummary struct {
	ID          string         `json:"id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Status      string         `json:"status"`
	Description string         `json:"description,omitempty"`
	Limits      map[string]any `json:"limits,omitempty"`
}

type tenantSubscriptionSummary struct {
	ID                  string `json:"id"`
	Status              string `json:"status"`
	PlanCode            string `json:"planCode,omitempty"`
	PlanName            string `json:"planName,omitempty"`
	TrialEndsAt         string `json:"trialEndsAt,omitempty"`
	CurrentPeriodEndsAt string `json:"currentPeriodEndsAt,omitempty"`
}

type tenantSubscriptionSaveInput struct {
	TenantID            string `json:"tenantId"`
	PlanCode            string `json:"planCode"`
	Status              string `json:"status"`
	TrialEndsAt         string `json:"trialEndsAt"`
	CurrentPeriodEndsAt string `json:"currentPeriodEndsAt"`
}

func handleSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	plans, err := listSubscriptionPlans(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load subscription plans", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"plans": plans})
}

func handleTenantSubscriptionSave(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	activeTenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input tenantSubscriptionSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeTenantSubscriptionSaveInput(input, activeTenantID)
	if err := validateTenantSubscriptionSaveInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if input.TenantID != activeTenantID {
		http.Error(w, "subscription update is limited to the active tenant", http.StatusForbidden)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	tenant, err := saveTenantSubscription(r.Context(), db, input, user, auditContextFromRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenants, err := listUserTenants(r.Context(), db, user.ID, tenant.ID)
	if authenticatedUserHasPermission(user, "operation_log.cross_tenant_view") || authenticatedUserHasPermission(user, "audit_log.cross_tenant_view") {
		tenants, err = listAllTenants(r.Context(), db, tenant.ID)
	}
	if err != nil {
		http.Error(w, "cannot reload tenants", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant": tenant, "tenants": tenants})
}

func listSubscriptionPlans(ctx context.Context, db *sql.DB) ([]subscriptionPlanSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id::text, code, name, status, description, limits
FROM subscription_plans
WHERE status <> 'archived'
ORDER BY CASE code WHEN 'free_trial' THEN 0 WHEN 'standard' THEN 1 ELSE 9 END, code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []subscriptionPlanSummary{}
	for rows.Next() {
		var plan subscriptionPlanSummary
		var limitsBytes []byte
		if err := rows.Scan(&plan.ID, &plan.Code, &plan.Name, &plan.Status, &plan.Description, &limitsBytes); err != nil {
			return nil, err
		}
		plan.Limits = decodeMetadata(limitsBytes)
		plans = append(plans, plan)
	}
	return plans, rows.Err()
}

func ensureTenantSubscription(ctx context.Context, exec masterDataExecutor, tenantID string, tenantStatus string, userID string) error {
	planCode := subscriptionPlanStandard
	subscriptionStatus := subscriptionStatusActive
	switch tenantStatus {
	case "trial":
		planCode = subscriptionPlanFreeTrial
		subscriptionStatus = subscriptionStatusTrial
	case "suspended":
		subscriptionStatus = subscriptionStatusSuspended
	case "archived":
		subscriptionStatus = subscriptionStatusCancelled
	}
	var trialEndsAt any
	if subscriptionStatus == subscriptionStatusTrial {
		trialEndsAt = time.Now().UTC().Add(30 * 24 * time.Hour)
	}
	_, err := exec.ExecContext(ctx, `
INSERT INTO tenant_subscriptions (
	tenant_id, plan_id, status, trial_ends_at, created_by_user_id, updated_by_user_id
)
SELECT $1::uuid, plan.id, $3, $4, nullif($2, '')::uuid, nullif($2, '')::uuid
FROM subscription_plans plan
WHERE plan.code = $5
ON CONFLICT (tenant_id) DO NOTHING`,
		tenantID,
		strings.TrimSpace(userID),
		subscriptionStatus,
		trialEndsAt,
		planCode,
	)
	return err
}

func saveTenantSubscription(ctx context.Context, db *sql.DB, input tenantSubscriptionSaveInput, user authenticatedUser, auditCtx requestAuditContext) (tenantSummary, error) {
	trialEndsAt, err := parseOptionalDate(input.TrialEndsAt)
	if err != nil {
		return tenantSummary{}, fmt.Errorf("trialEndsAt must use YYYY-MM-DD")
	}
	currentPeriodEndsAt, err := parseOptionalDate(input.CurrentPeriodEndsAt)
	if err != nil {
		return tenantSummary{}, fmt.Errorf("currentPeriodEndsAt must use YYYY-MM-DD")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return tenantSummary{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
UPDATE tenant_subscriptions ts
SET plan_id = plan.id,
	status = $3,
	trial_ends_at = $4,
	current_period_ends_at = $5,
	updated_by_user_id = nullif($6, '')::uuid,
	updated_at = now()
FROM subscription_plans plan
WHERE ts.tenant_id = $1::uuid
	AND plan.code = $2
	AND plan.status = 'active'`,
		input.TenantID,
		input.PlanCode,
		input.Status,
		trialEndsAt,
		currentPeriodEndsAt,
		user.ID,
	)
	if err != nil {
		return tenantSummary{}, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return tenantSummary{}, fmt.Errorf("subscription plan %q is not available", input.PlanCode)
	}

	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "tenant.subscription.update",
		EntityType: "tenant_subscription",
		EntityID:   input.TenantID,
		Metadata: map[string]any{
			"planCode":            input.PlanCode,
			"subscriptionStatus":  input.Status,
			"trialEndsAt":         input.TrialEndsAt,
			"currentPeriodEndsAt": input.CurrentPeriodEndsAt,
		},
	})

	if err := tx.Commit(); err != nil {
		return tenantSummary{}, err
	}

	tenants, err := listUserTenants(ctx, db, user.ID, input.TenantID)
	if err == nil {
		for _, tenant := range tenants {
			if tenant.ID == input.TenantID {
				return tenant, nil
			}
		}
	}
	return loadTenantSummaryByID(ctx, db, input.TenantID)
}

func loadTenantSummaryByID(ctx context.Context, db *sql.DB, tenantID string) (tenantSummary, error) {
	var tenant tenantSummary
	var trialEndsAt sql.NullTime
	var currentPeriodEndsAt sql.NullTime
	err := db.QueryRowContext(ctx, `
SELECT tenant.id::text,
	tenant.code,
	tenant.name,
	tenant.status,
	COUNT(school.id)::integer,
	COALESCE(ts.status, ''),
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.trial_ends_at,
	ts.current_period_ends_at
FROM tenants tenant
LEFT JOIN schools school ON school.tenant_id = tenant.id
LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE tenant.id = $1::uuid
GROUP BY tenant.id, tenant.code, tenant.name, tenant.status, ts.status, plan.code, plan.name, ts.trial_ends_at, ts.current_period_ends_at`,
		tenantID,
	).Scan(
		&tenant.ID,
		&tenant.Code,
		&tenant.Name,
		&tenant.Status,
		&tenant.SchoolCount,
		&tenant.SubscriptionStatus,
		&tenant.PlanCode,
		&tenant.PlanName,
		&trialEndsAt,
		&currentPeriodEndsAt,
	)
	if err != nil {
		return tenantSummary{}, err
	}
	tenant.TrialEndsAt = formatNullDate(trialEndsAt)
	tenant.CurrentPeriodEndsAt = formatNullDate(currentPeriodEndsAt)
	return tenant, nil
}

func normalizeTenantSubscriptionSaveInput(input tenantSubscriptionSaveInput, activeTenantID string) tenantSubscriptionSaveInput {
	input.TenantID = firstNonEmpty(strings.TrimSpace(input.TenantID), strings.TrimSpace(activeTenantID))
	input.PlanCode = headerKey(input.PlanCode)
	input.Status = headerKey(input.Status)
	input.TrialEndsAt = strings.TrimSpace(input.TrialEndsAt)
	input.CurrentPeriodEndsAt = strings.TrimSpace(input.CurrentPeriodEndsAt)
	return input
}

func validateTenantSubscriptionSaveInput(input tenantSubscriptionSaveInput) error {
	if strings.TrimSpace(input.TenantID) == "" {
		return fmt.Errorf("tenantId is required")
	}
	if strings.TrimSpace(input.PlanCode) == "" {
		return fmt.Errorf("planCode is required")
	}
	if !isTenantSubscriptionStatus(input.Status) {
		return fmt.Errorf("subscription status must be trial, active, past_due, suspended, or cancelled")
	}
	return nil
}

func parseOptionalDate(value string) (any, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return parsed.UTC(), nil
}

func formatNullDate(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format("2006-01-02")
}

func isTenantSubscriptionStatus(value string) bool {
	switch strings.TrimSpace(value) {
	case subscriptionStatusTrial, subscriptionStatusActive, subscriptionStatusPastDue, subscriptionStatusSuspended, subscriptionStatusCancelled:
		return true
	default:
		return false
	}
}

func subscriptionWritePermissionBlocked(status string) bool {
	switch strings.TrimSpace(status) {
	case "", subscriptionStatusTrial, subscriptionStatusActive:
		return false
	default:
		return true
	}
}

func permissionRequiresActiveTenantSubscription(permission string) bool {
	switch permission {
	case "student.create",
		"student.update",
		"school_tree.update",
		"fee.create",
		"fee.update",
		"invoice.create",
		"invoice.update",
		"payment.create",
		"payment.reconcile",
		"notification.create",
		"notification.send",
		"email_config.update",
		"email_cron.update",
		"tenant.update":
		return true
	default:
		return false
	}
}

func defaultSubscriptionPlanCodeForTenantStatus(tenantStatus string) string {
	if strings.TrimSpace(tenantStatus) == "trial" {
		return subscriptionPlanFreeTrial
	}
	return subscriptionPlanStandard
}

func defaultSubscriptionStatusForTenantStatus(tenantStatus string) string {
	switch strings.TrimSpace(tenantStatus) {
	case "trial":
		return subscriptionStatusTrial
	case "suspended":
		return subscriptionStatusSuspended
	case "archived":
		return subscriptionStatusCancelled
	default:
		return subscriptionStatusActive
	}
}

func defaultSubscriptionPlanName(planCode string) string {
	switch strings.TrimSpace(planCode) {
	case subscriptionPlanFreeTrial:
		return "Free Trial"
	case subscriptionPlanStandard:
		return "Standard"
	default:
		return ""
	}
}
