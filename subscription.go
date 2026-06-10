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

	defaultSubscriptionPlanDisplayOrder = 100

	subscriptionStatusTrial     = "trial"
	subscriptionStatusActive    = "active"
	subscriptionStatusPastDue   = "past_due"
	subscriptionStatusSuspended = "suspended"
	subscriptionStatusCancelled = "cancelled"
)

type subscriptionPlanSummary struct {
	ID                  string         `json:"id"`
	Code                string         `json:"code"`
	Name                string         `json:"name"`
	Status              string         `json:"status"`
	Description         string         `json:"description,omitempty"`
	ContactPrice        bool           `json:"contactPrice"`
	DisplayOrder        int            `json:"displayOrder"`
	BasePriceVND        int            `json:"basePriceVnd"`
	PartnerPriceVND     *int           `json:"partnerPriceVnd,omitempty"`
	PromotionalPriceVND *int           `json:"promotionalPriceVnd,omitempty"`
	DisplayPriceVND     int            `json:"displayPriceVnd"`
	Limits              map[string]any `json:"limits,omitempty"`
	UpdatedAt           string         `json:"updatedAt,omitempty"`
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

type subscriptionPlanSaveInput struct {
	ID                  string         `json:"id"`
	Code                string         `json:"code"`
	Name                string         `json:"name"`
	Status              string         `json:"status"`
	Description         string         `json:"description"`
	ContactPrice        bool           `json:"contactPrice"`
	DisplayOrder        *int           `json:"displayOrder"`
	BasePriceVND        int            `json:"basePriceVnd"`
	PartnerPriceVND     *int           `json:"partnerPriceVnd"`
	PromotionalPriceVND *int           `json:"promotionalPriceVnd"`
	Limits              map[string]any `json:"limits"`
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

func handlePublicSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	handleSubscriptionPlans(w, r)
}

func handlePlatformSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	if !user.IsPlatformAdmin {
		http.Error(w, "platform admin required", http.StatusForbidden)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		plans, err := listPlatformSubscriptionPlans(r.Context(), db)
		if err != nil {
			http.Error(w, "cannot load subscription plans", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"plans": plans})
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var input subscriptionPlanSaveInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		plan, err := saveSubscriptionPlan(r.Context(), db, input, user, auditContextFromRequest(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		plans, err := listPlatformSubscriptionPlans(r.Context(), db)
		if err != nil {
			http.Error(w, "cannot reload subscription plans", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"plan": plan, "plans": plans})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTenantSubscriptionSave(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	if !user.IsPlatformAdmin {
		http.Error(w, "platform admin required", http.StatusForbidden)
		return
	}
	activeTenantID := strings.TrimSpace(user.ActiveTenant.ID)
	if activeTenantID == "" && !user.IsPlatformAdmin {
		http.Error(w, "active tenant required", http.StatusForbidden)
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
	if input.TenantID != activeTenantID && !user.IsPlatformAdmin {
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
	if user.IsPlatformAdmin && (authenticatedUserHasPermission(user, "operation_log.cross_tenant_view") || authenticatedUserHasPermission(user, "audit_log.cross_tenant_view") || authenticatedUserHasPermission(user, "tenant.view")) {
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
		SELECT id::text,
			code,
			name,
			status,
			description,
			contact_price,
			display_order,
			base_price_vnd,
			partner_price_vnd,
			promotional_price_vnd,
			limits,
			updated_at
		FROM subscription_plans
		WHERE status <> 'archived'
		ORDER BY display_order,
			CASE code
				WHEN 'free' THEN 0
				WHEN 'go' THEN 1
				WHEN 'plus' THEN 2
				WHEN 'pro' THEN 3
				WHEN 'free_trial' THEN 10
				WHEN 'standard' THEN 11
				ELSE 99
			END,
			code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []subscriptionPlanSummary{}
	for rows.Next() {
		plan, err := scanSubscriptionPlanSummary(rows, false)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, rows.Err()
}

func listPlatformSubscriptionPlans(ctx context.Context, db *sql.DB) ([]subscriptionPlanSummary, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT id::text,
		code,
		name,
		status,
		description,
		contact_price,
		display_order,
		base_price_vnd,
		partner_price_vnd,
		promotional_price_vnd,
		limits,
		updated_at
	FROM subscription_plans
	ORDER BY CASE status WHEN 'active' THEN 0 ELSE 1 END,
	display_order,
	CASE code
		WHEN 'free' THEN 0
		WHEN 'go' THEN 1
		WHEN 'plus' THEN 2
		WHEN 'pro' THEN 3
		WHEN 'free_trial' THEN 10
		WHEN 'standard' THEN 11
		ELSE 99
	END, code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []subscriptionPlanSummary{}
	for rows.Next() {
		plan, err := scanSubscriptionPlanSummary(rows, true)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, rows.Err()
}

func scanSubscriptionPlanSummary(rows *sql.Rows, includePartnerPrice bool) (subscriptionPlanSummary, error) {
	var plan subscriptionPlanSummary
	var partnerPrice sql.NullInt64
	var promotionalPrice sql.NullInt64
	var limitsBytes []byte
	var updatedAt time.Time
	if err := rows.Scan(
		&plan.ID,
		&plan.Code,
		&plan.Name,
		&plan.Status,
		&plan.Description,
		&plan.ContactPrice,
		&plan.DisplayOrder,
		&plan.BasePriceVND,
		&partnerPrice,
		&promotionalPrice,
		&limitsBytes,
		&updatedAt,
	); err != nil {
		return subscriptionPlanSummary{}, err
	}
	if includePartnerPrice {
		plan.PartnerPriceVND = intPtrFromNullInt64(partnerPrice)
	}
	plan.PromotionalPriceVND = intPtrFromNullInt64(promotionalPrice)
	plan.DisplayPriceVND = subscriptionPlanDisplayPrice(plan.BasePriceVND, plan.PromotionalPriceVND)
	plan.Limits = decodeMetadata(limitsBytes)
	plan.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
	return plan, nil
}

func saveSubscriptionPlan(ctx context.Context, db *sql.DB, input subscriptionPlanSaveInput, user authenticatedUser, auditCtx requestAuditContext) (subscriptionPlanSummary, error) {
	input = normalizeSubscriptionPlanSaveInput(input)
	if err := validateSubscriptionPlanSaveInput(input); err != nil {
		return subscriptionPlanSummary{}, err
	}
	limits, err := json.Marshal(input.Limits)
	if err != nil {
		return subscriptionPlanSummary{}, fmt.Errorf("limits must be a JSON object")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return subscriptionPlanSummary{}, err
	}
	defer tx.Rollback()

	var plan subscriptionPlanSummary
	var partnerPrice sql.NullInt64
	var promotionalPrice sql.NullInt64
	var limitsBytes []byte
	var updatedAt time.Time
	if input.ID != "" {
		if err := ensureSubscriptionPlanCodeAvailable(ctx, tx, input.ID, input.Code); err != nil {
			return subscriptionPlanSummary{}, err
		}
		err = tx.QueryRowContext(ctx, `
			UPDATE subscription_plans
			SET code = $2,
				name = $3,
				status = $4,
				description = $5,
				contact_price = $6,
				display_order = $7,
				base_price_vnd = $8,
				partner_price_vnd = $9,
				promotional_price_vnd = $10,
				limits = $11::jsonb,
				updated_at = now()
			WHERE id = $1::uuid
			RETURNING id::text,
				code,
				name,
				status,
				description,
				contact_price,
				display_order,
				base_price_vnd,
				partner_price_vnd,
				promotional_price_vnd,
				limits,
				updated_at`,
			input.ID,
			input.Code,
			input.Name,
			input.Status,
			input.Description,
			input.ContactPrice,
			subscriptionPlanDisplayOrder(input),
			input.BasePriceVND,
			nullableIntArg(input.PartnerPriceVND),
			nullableIntArg(input.PromotionalPriceVND),
			string(limits),
		).Scan(
			&plan.ID,
			&plan.Code,
			&plan.Name,
			&plan.Status,
			&plan.Description,
			&plan.ContactPrice,
			&plan.DisplayOrder,
			&plan.BasePriceVND,
			&partnerPrice,
			&promotionalPrice,
			&limitsBytes,
			&updatedAt,
		)
	} else {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO subscription_plans (
				code,
				name,
				status,
				description,
				contact_price,
				display_order,
				base_price_vnd,
				partner_price_vnd,
				promotional_price_vnd,
				limits
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb)
			ON CONFLICT (code) DO UPDATE
			SET name = EXCLUDED.name,
				status = EXCLUDED.status,
				description = EXCLUDED.description,
				contact_price = EXCLUDED.contact_price,
				display_order = EXCLUDED.display_order,
				base_price_vnd = EXCLUDED.base_price_vnd,
				partner_price_vnd = EXCLUDED.partner_price_vnd,
				promotional_price_vnd = EXCLUDED.promotional_price_vnd,
				limits = EXCLUDED.limits,
				updated_at = now()
			RETURNING id::text,
				code,
				name,
				status,
				description,
				contact_price,
				display_order,
				base_price_vnd,
				partner_price_vnd,
				promotional_price_vnd,
				limits,
				updated_at`,
			input.Code,
			input.Name,
			input.Status,
			input.Description,
			input.ContactPrice,
			subscriptionPlanDisplayOrder(input),
			input.BasePriceVND,
			nullableIntArg(input.PartnerPriceVND),
			nullableIntArg(input.PromotionalPriceVND),
			string(limits),
		).Scan(
			&plan.ID,
			&plan.Code,
			&plan.Name,
			&plan.Status,
			&plan.Description,
			&plan.ContactPrice,
			&plan.DisplayOrder,
			&plan.BasePriceVND,
			&partnerPrice,
			&promotionalPrice,
			&limitsBytes,
			&updatedAt,
		)
	}
	if err == sql.ErrNoRows && input.ID != "" {
		return subscriptionPlanSummary{}, fmt.Errorf("subscription plan not found")
	}
	if err != nil {
		return subscriptionPlanSummary{}, err
	}
	plan.PartnerPriceVND = intPtrFromNullInt64(partnerPrice)
	plan.PromotionalPriceVND = intPtrFromNullInt64(promotionalPrice)
	plan.DisplayPriceVND = subscriptionPlanDisplayPrice(plan.BasePriceVND, plan.PromotionalPriceVND)
	plan.Limits = decodeMetadata(limitsBytes)
	plan.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)

	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "subscription_plan.save",
		EntityType: "subscription_plan",
		EntityID:   plan.ID,
		Metadata: map[string]any{
			"code":                plan.Code,
			"name":                plan.Name,
			"status":              plan.Status,
			"contactPrice":        plan.ContactPrice,
			"displayOrder":        plan.DisplayOrder,
			"basePriceVnd":        plan.BasePriceVND,
			"partnerPriceVnd":     plan.PartnerPriceVND,
			"promotionalPriceVnd": plan.PromotionalPriceVND,
			"displayPriceVnd":     plan.DisplayPriceVND,
			"limits":              plan.Limits,
		},
	})

	if err := tx.Commit(); err != nil {
		return subscriptionPlanSummary{}, err
	}
	return plan, nil
}

func ensureSubscriptionPlanCodeAvailable(ctx context.Context, exec masterDataExecutor, planID string, code string) error {
	var existingID string
	err := exec.QueryRowContext(ctx, `
	SELECT id::text
	FROM subscription_plans
	WHERE code = $1
		AND id <> $2::uuid
	LIMIT 1`, code, planID).Scan(&existingID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("plan code already exists")
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
	usageMetrics, err := loadTenantUsageSummaries(ctx, db, tenant.ID, time.Now())
	if err != nil {
		return tenantSummary{}, err
	}
	tenant.UsageMetrics = usageMetrics
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

func normalizeSubscriptionPlanSaveInput(input subscriptionPlanSaveInput) subscriptionPlanSaveInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Code = headerKey(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.Status = headerKey(firstNonEmpty(input.Status, "active"))
	input.Description = strings.TrimSpace(input.Description)
	if input.DisplayOrder == nil {
		input.DisplayOrder = intPtr(defaultSubscriptionPlanDisplayOrder)
	}
	if input.Limits == nil {
		input.Limits = map[string]any{}
	}
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

func validateSubscriptionPlanSaveInput(input subscriptionPlanSaveInput) error {
	if !isSubscriptionPlanCode(input.Code) {
		return fmt.Errorf("plan code must start with a letter and contain only lowercase letters, numbers, underscore, dash, or colon")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("plan name is required")
	}
	if !isSubscriptionPlanStatus(input.Status) {
		return fmt.Errorf("plan status must be active or archived")
	}
	if input.DisplayOrder == nil {
		return fmt.Errorf("displayOrder is required")
	}
	if *input.DisplayOrder < 0 {
		return fmt.Errorf("displayOrder must be greater than or equal to 0")
	}
	if input.BasePriceVND < 0 {
		return fmt.Errorf("basePriceVnd must be greater than or equal to 0")
	}
	if input.PartnerPriceVND != nil && *input.PartnerPriceVND < 0 {
		return fmt.Errorf("partnerPriceVnd must be greater than or equal to 0")
	}
	if input.PromotionalPriceVND != nil && *input.PromotionalPriceVND < 0 {
		return fmt.Errorf("promotionalPriceVnd must be greater than or equal to 0")
	}
	if input.PromotionalPriceVND != nil && input.BasePriceVND > 0 && *input.PromotionalPriceVND > input.BasePriceVND {
		return fmt.Errorf("promotionalPriceVnd must be less than or equal to basePriceVnd")
	}
	if input.Limits == nil {
		return fmt.Errorf("limits must be a JSON object")
	}
	for key := range input.Limits {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("limits cannot contain blank keys")
		}
	}
	return nil
}

func subscriptionPlanDisplayPrice(basePriceVND int, promotionalPriceVND *int) int {
	if promotionalPriceVND != nil {
		return *promotionalPriceVND
	}
	return basePriceVND
}

func subscriptionPlanDisplayOrder(input subscriptionPlanSaveInput) int {
	if input.DisplayOrder == nil {
		return defaultSubscriptionPlanDisplayOrder
	}
	return *input.DisplayOrder
}

func intPtr(value int) *int {
	return &value
}

func intPtrFromNullInt64(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	converted := int(value.Int64)
	return &converted
}

func nullableIntArg(value *int) any {
	if value == nil {
		return nil
	}
	return *value
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

func isSubscriptionPlanStatus(value string) bool {
	switch strings.TrimSpace(value) {
	case "active", "archived":
		return true
	default:
		return false
	}
}

func isSubscriptionPlanCode(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	for idx, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
		case idx > 0 && r >= '0' && r <= '9':
		case idx > 0 && (r == '_' || r == '-' || r == ':'):
		default:
			return false
		}
	}
	return true
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
