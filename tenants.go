package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type tenantSummary struct {
	ID                  string                     `json:"id"`
	Code                string                     `json:"code"`
	Name                string                     `json:"name"`
	Status              string                     `json:"status"`
	MembershipStatus    string                     `json:"membershipStatus,omitempty"`
	IsOwner             bool                       `json:"isOwner,omitempty"`
	SchoolCount         int                        `json:"schoolCount"`
	IsActive            bool                       `json:"isActive,omitempty"`
	SubscriptionStatus  string                     `json:"subscriptionStatus,omitempty"`
	PlanCode            string                     `json:"planCode,omitempty"`
	PlanName            string                     `json:"planName,omitempty"`
	TrialEndsAt         string                     `json:"trialEndsAt,omitempty"`
	CurrentPeriodEndsAt string                     `json:"currentPeriodEndsAt,omitempty"`
	UsageMetrics        []tenantUsageMetricSummary `json:"usageMetrics,omitempty"`
}

type tenantListResponse struct {
	ActiveTenantID string          `json:"activeTenantId"`
	Tenants        []tenantSummary `json:"tenants"`
}

type tenantSaveInput struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	InitialSchoolCode string `json:"initialSchoolCode"`
	InitialSchoolName string `json:"initialSchoolName"`
}

type tenantSwitchInput struct {
	TenantID string `json:"tenantId"`
}

func handleTenants(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	tenants, err := listUserTenants(r.Context(), db, user.ID, tenantID)
	if authenticatedUserHasPermission(user, "operation_log.cross_tenant_view") || authenticatedUserHasPermission(user, "audit_log.cross_tenant_view") {
		tenants, err = listAllTenants(r.Context(), db, tenantID)
	}
	if err != nil {
		http.Error(w, "cannot load tenants", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tenantListResponse{ActiveTenantID: tenantID, Tenants: tenants})
}

func handleTenantSave(w http.ResponseWriter, r *http.Request) {
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
	var input tenantSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeTenantSaveInput(input)
	if err := validateTenantSaveInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	var tenant tenantSummary
	if input.ID == "" {
		tenant, err = createTenantWithInitialSchool(r.Context(), db, input, user, auditContextFromRequest(r))
	} else {
		if input.ID != activeTenantID {
			http.Error(w, "tenant update is limited to the active tenant", http.StatusForbidden)
			return
		}
		tenant, err = updateActiveTenant(r.Context(), db, input, user, auditContextFromRequest(r))
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenants, err := listUserTenants(r.Context(), db, user.ID, firstNonEmpty(tenant.ID, activeTenantID))
	if authenticatedUserHasPermission(user, "operation_log.cross_tenant_view") || authenticatedUserHasPermission(user, "audit_log.cross_tenant_view") {
		tenants, err = listAllTenants(r.Context(), db, firstNonEmpty(tenant.ID, activeTenantID))
	}
	if err != nil {
		http.Error(w, "cannot reload tenants", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant": tenant, "tenants": tenants})
}

func handleAuthTenantSwitch(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input tenantSwitchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeTenantSwitchInput(input)
	if err := validateTenantSwitchInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	oldAccessToken, _ := readCookieValue(r, authAccessCookieName)
	oldRefreshToken, _ := readCookieValue(r, authRefreshCookieName)
	now := time.Now().UTC()
	session, tokens, err := switchAuthTenantSession(r.Context(), db, user.ID, input.TenantID, now, loadAuthConfig(), auditContextFromRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	setAuthCookies(w, r, tokens, loadAuthConfig())
	if oldAccessToken != "" || oldRefreshToken != "" {
		_ = revokeSessionByTokenHashes(r.Context(), db, hashAuthToken(oldAccessToken), hashAuthToken(oldRefreshToken), "tenant switch")
	}
	writeJSON(w, http.StatusOK, session)
}

func tenantSavePermission(r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		return permissionAuthenticated, nil
	}
	body, err := readAndRestoreTenantRequestBody(r, 1<<20)
	if err != nil {
		return "", err
	}
	var input tenantSaveInput
	if strings.TrimSpace(string(body)) != "" {
		_ = json.Unmarshal(body, &input)
	}
	if strings.TrimSpace(input.ID) == "" {
		return "tenant.create", nil
	}
	return "tenant.update", nil
}

func readAndRestoreTenantRequestBody(r *http.Request, limit int64) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func normalizeTenantSaveInput(input tenantSaveInput) tenantSaveInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Code = normalizeSchoolCode(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.Status = headerKey(input.Status)
	if input.Status == "" {
		input.Status = "active"
	}
	input.InitialSchoolCode = normalizeSchoolCode(input.InitialSchoolCode)
	input.InitialSchoolName = strings.TrimSpace(input.InitialSchoolName)
	if input.ID == "" {
		if input.InitialSchoolCode == "" {
			input.InitialSchoolCode = input.Code
		}
		if input.InitialSchoolName == "" {
			input.InitialSchoolName = input.Name
		}
	}
	return input
}

func validateTenantSaveInput(input tenantSaveInput) error {
	if strings.TrimSpace(input.Code) == "" {
		return fmt.Errorf("tenant code is required")
	}
	if !isSafeTenantCode(input.Code) {
		return fmt.Errorf("tenant code may contain only letters, numbers, underscore, or hyphen")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("tenant name is required")
	}
	if !isTenantStatus(input.Status) {
		return fmt.Errorf("tenant status must be active, trial, suspended, or archived")
	}
	if strings.TrimSpace(input.ID) == "" {
		if strings.TrimSpace(input.InitialSchoolCode) == "" {
			return fmt.Errorf("initial school code is required")
		}
		if !isSafeTenantCode(input.InitialSchoolCode) {
			return fmt.Errorf("initial school code may contain only letters, numbers, underscore, or hyphen")
		}
		if strings.TrimSpace(input.InitialSchoolName) == "" {
			return fmt.Errorf("initial school name is required")
		}
	}
	return nil
}

func normalizeTenantSwitchInput(input tenantSwitchInput) tenantSwitchInput {
	input.TenantID = strings.TrimSpace(input.TenantID)
	return input
}

func validateTenantSwitchInput(input tenantSwitchInput) error {
	if strings.TrimSpace(input.TenantID) == "" {
		return fmt.Errorf("tenantId is required")
	}
	return nil
}

func isTenantStatus(status string) bool {
	switch status {
	case "active", "trial", "suspended", "archived":
		return true
	default:
		return false
	}
}

func isSafeTenantCode(code string) bool {
	if code == "" {
		return false
	}
	for _, ch := range code {
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if ch >= '0' && ch <= '9' {
			continue
		}
		if ch == '_' || ch == '-' {
			continue
		}
		return false
	}
	return true
}

func listUserTenants(ctx context.Context, db *sql.DB, userID string, activeTenantID string) ([]tenantSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT tenant.id::text,
	tenant.code,
	tenant.name,
	tenant.status,
	membership.status,
	membership.is_owner,
	COUNT(school.id)::integer,
	COALESCE(ts.status, ''),
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.trial_ends_at,
	ts.current_period_ends_at
FROM tenant_memberships membership
JOIN tenants tenant ON tenant.id = membership.tenant_id
LEFT JOIN schools school ON school.tenant_id = tenant.id
LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE membership.user_id = $1::uuid
	AND membership.status <> 'removed'
GROUP BY tenant.id, tenant.code, tenant.name, tenant.status, membership.status, membership.is_owner, ts.status, plan.code, plan.name, ts.trial_ends_at, ts.current_period_ends_at
ORDER BY CASE WHEN tenant.id::text = $2 THEN 0 WHEN tenant.code = $3 THEN 1 ELSE 2 END,
	tenant.code`, userID, strings.TrimSpace(activeTenantID), defaultTenantCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []tenantSummary{}
	for rows.Next() {
		var tenant tenantSummary
		var trialEndsAt sql.NullTime
		var currentPeriodEndsAt sql.NullTime
		if err := rows.Scan(
			&tenant.ID,
			&tenant.Code,
			&tenant.Name,
			&tenant.Status,
			&tenant.MembershipStatus,
			&tenant.IsOwner,
			&tenant.SchoolCount,
			&tenant.SubscriptionStatus,
			&tenant.PlanCode,
			&tenant.PlanName,
			&trialEndsAt,
			&currentPeriodEndsAt,
		); err != nil {
			return nil, err
		}
		tenant.IsActive = tenant.ID == strings.TrimSpace(activeTenantID)
		tenant.TrialEndsAt = formatNullDate(trialEndsAt)
		tenant.CurrentPeriodEndsAt = formatNullDate(currentPeriodEndsAt)
		usageMetrics, err := loadTenantUsageSummaries(ctx, db, tenant.ID, time.Now())
		if err != nil {
			return nil, err
		}
		tenant.UsageMetrics = usageMetrics
		tenants = append(tenants, tenant)
	}
	return tenants, rows.Err()
}

func listAllTenants(ctx context.Context, db *sql.DB, activeTenantID string) ([]tenantSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT tenant.id::text,
	tenant.code,
	tenant.name,
	tenant.status,
	'' AS membership_status,
	false AS is_owner,
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
GROUP BY tenant.id, tenant.code, tenant.name, tenant.status, ts.status, plan.code, plan.name, ts.trial_ends_at, ts.current_period_ends_at
ORDER BY CASE WHEN tenant.id::text = $1 THEN 0 WHEN tenant.code = $2 THEN 1 ELSE 2 END,
	tenant.code`, strings.TrimSpace(activeTenantID), defaultTenantCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []tenantSummary{}
	for rows.Next() {
		var tenant tenantSummary
		var trialEndsAt sql.NullTime
		var currentPeriodEndsAt sql.NullTime
		if err := rows.Scan(
			&tenant.ID,
			&tenant.Code,
			&tenant.Name,
			&tenant.Status,
			&tenant.MembershipStatus,
			&tenant.IsOwner,
			&tenant.SchoolCount,
			&tenant.SubscriptionStatus,
			&tenant.PlanCode,
			&tenant.PlanName,
			&trialEndsAt,
			&currentPeriodEndsAt,
		); err != nil {
			return nil, err
		}
		tenant.IsActive = tenant.ID == strings.TrimSpace(activeTenantID)
		tenant.TrialEndsAt = formatNullDate(trialEndsAt)
		tenant.CurrentPeriodEndsAt = formatNullDate(currentPeriodEndsAt)
		usageMetrics, err := loadTenantUsageSummaries(ctx, db, tenant.ID, time.Now())
		if err != nil {
			return nil, err
		}
		tenant.UsageMetrics = usageMetrics
		tenants = append(tenants, tenant)
	}
	return tenants, rows.Err()
}

func createTenantWithInitialSchool(ctx context.Context, db *sql.DB, input tenantSaveInput, user authenticatedUser, auditCtx requestAuditContext) (tenantSummary, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return tenantSummary{}, err
	}
	defer tx.Rollback()

	var tenant tenantSummary
	err = tx.QueryRowContext(ctx, `
INSERT INTO tenants (code, name, status, created_by_user_id, updated_by_user_id)
VALUES ($1, $2, $3, nullif($4, '')::uuid, nullif($4, '')::uuid)
ON CONFLICT (code) DO NOTHING
RETURNING id::text, code, name, status`, input.Code, input.Name, input.Status, user.ID).Scan(
		&tenant.ID,
		&tenant.Code,
		&tenant.Name,
		&tenant.Status,
	)
	if err == sql.ErrNoRows {
		return tenantSummary{}, fmt.Errorf("tenant code already exists")
	}
	if err != nil {
		return tenantSummary{}, err
	}
	if err := ensureTenantMembership(ctx, tx, tenant.ID, user.ID, true); err != nil {
		return tenantSummary{}, err
	}
	if err := ensureTenantSubscription(ctx, tx, tenant.ID, tenant.Status, user.ID); err != nil {
		return tenantSummary{}, err
	}
	if err := ensureTenantUserRole(ctx, tx, tenant.ID, user.ID, "admin", user.ID); err != nil {
		return tenantSummary{}, err
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO schools (tenant_id, code, name, status, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2, $3, 'active', nullif($4, '')::uuid, nullif($4, '')::uuid)
ON CONFLICT (tenant_id, code) DO NOTHING`, tenant.ID, input.InitialSchoolCode, input.InitialSchoolName, user.ID); err != nil {
		return tenantSummary{}, err
	}
	auditCtx.TenantID = tenant.ID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "tenant.create",
		EntityType: "tenant",
		EntityID:   tenant.ID,
		Metadata: map[string]any{
			"tenantCode":        tenant.Code,
			"tenantName":        tenant.Name,
			"initialSchoolCode": input.InitialSchoolCode,
		},
	})
	if err := tx.Commit(); err != nil {
		return tenantSummary{}, err
	}
	tenant.MembershipStatus = "active"
	tenant.IsOwner = true
	tenant.SchoolCount = 1
	tenant.SubscriptionStatus = defaultSubscriptionStatusForTenantStatus(tenant.Status)
	tenant.PlanCode = defaultSubscriptionPlanCodeForTenantStatus(tenant.Status)
	tenant.PlanName = defaultSubscriptionPlanName(tenant.PlanCode)
	if tenant.SubscriptionStatus == subscriptionStatusTrial {
		tenant.TrialEndsAt = time.Now().UTC().Add(30 * 24 * time.Hour).Format("2006-01-02")
	}
	tenant.UsageMetrics = []tenantUsageMetricSummary{}
	return tenant, nil
}

func updateActiveTenant(ctx context.Context, db *sql.DB, input tenantSaveInput, user authenticatedUser, auditCtx requestAuditContext) (tenantSummary, error) {
	var tenant tenantSummary
	err := db.QueryRowContext(ctx, `
UPDATE tenants
SET code = $2,
	name = $3,
	status = $4,
	updated_by_user_id = nullif($5, '')::uuid,
	updated_at = now()
WHERE id = $1::uuid
RETURNING id::text, code, name, status`, input.ID, input.Code, input.Name, input.Status, user.ID).Scan(
		&tenant.ID,
		&tenant.Code,
		&tenant.Name,
		&tenant.Status,
	)
	if err != nil {
		return tenantSummary{}, err
	}
	auditCtx.TenantID = tenant.ID
	_ = insertAuditLog(ctx, db, auditLogInput{
		Context:    auditCtx,
		Action:     "tenant.update",
		EntityType: "tenant",
		EntityID:   tenant.ID,
		Metadata: map[string]any{
			"tenantCode": tenant.Code,
			"tenantName": tenant.Name,
			"status":     tenant.Status,
		},
	})
	if tenants, err := listUserTenants(ctx, db, user.ID, tenant.ID); err == nil {
		for _, item := range tenants {
			if item.ID == tenant.ID {
				return item, nil
			}
		}
	}
	return tenant, nil
}

func switchAuthTenantSession(ctx context.Context, db *sql.DB, userID string, tenantID string, now time.Time, cfg authConfig, auditCtx requestAuditContext) (authSessionResponse, authIssuedTokens, error) {
	tenant, err := loadSwitchableTenantForUser(ctx, db, userID, tenantID)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, fmt.Errorf("active tenant membership required")
	}
	tokens, err := issueAuthSession(ctx, db, userID, tenant.ID, now, cfg, auditCtx)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	session, err := loadAuthSessionByAccessToken(ctx, db, tokens.AccessTokenHash, now)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	auditCtx.TenantID = tenant.ID
	_ = insertAuditLog(ctx, db, auditLogInput{
		Context:    auditCtx,
		Action:     "tenant.switch",
		EntityType: "tenant",
		EntityID:   tenant.ID,
		Metadata: map[string]any{
			"tenantCode": tenant.Code,
			"tenantName": tenant.Name,
		},
	})
	return session, tokens, nil
}

func loadSwitchableTenantForUser(ctx context.Context, db *sql.DB, userID string, tenantID string) (authTenantSummary, error) {
	var tenant authTenantSummary
	var trialEndsAt sql.NullTime
	var currentPeriodEndsAt sql.NullTime
	err := db.QueryRowContext(ctx, `
SELECT tenant.id::text,
	tenant.code,
	tenant.name,
	tenant.status,
	membership.status,
	membership.is_owner,
	COALESCE(ts.status, ''),
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.trial_ends_at,
	ts.current_period_ends_at
FROM tenant_memberships membership
JOIN tenants tenant ON tenant.id = membership.tenant_id
LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE membership.user_id = $1::uuid
	AND tenant.id = $2::uuid
	AND membership.status = 'active'
	AND tenant.status IN ('active', 'trial')`, userID, strings.TrimSpace(tenantID)).Scan(
		&tenant.ID,
		&tenant.Code,
		&tenant.Name,
		&tenant.Status,
		&tenant.MembershipStatus,
		&tenant.IsOwner,
		&tenant.SubscriptionStatus,
		&tenant.PlanCode,
		&tenant.PlanName,
		&trialEndsAt,
		&currentPeriodEndsAt,
	)
	tenant.TrialEndsAt = formatNullDate(trialEndsAt)
	tenant.CurrentPeriodEndsAt = formatNullDate(currentPeriodEndsAt)
	return tenant, err
}
