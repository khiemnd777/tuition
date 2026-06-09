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

type tenantIntakeRequestSummary struct {
	ID              string    `json:"id"`
	SchoolName      string    `json:"schoolName"`
	ContactName     string    `json:"contactName"`
	ContactEmail    string    `json:"contactEmail"`
	ContactPhone    string    `json:"contactPhone"`
	DesiredPlanCode string    `json:"desiredPlanCode"`
	Note            string    `json:"note"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	HandledAt       string    `json:"handledAt,omitempty"`
}

type tenantIntakeSubmitInput struct {
	SchoolName      string `json:"schoolName"`
	ContactName     string `json:"contactName"`
	ContactEmail    string `json:"contactEmail"`
	ContactPhone    string `json:"contactPhone"`
	DesiredPlanCode string `json:"desiredPlanCode"`
	Note            string `json:"note"`
}

type tenantIntakeStatusInput struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type platformTenantOnboardInput struct {
	IntakeRequestID     string `json:"intakeRequestId"`
	TenantCode          string `json:"tenantCode"`
	TenantName          string `json:"tenantName"`
	TenantStatus        string `json:"tenantStatus"`
	InitialSchoolCode   string `json:"initialSchoolCode"`
	InitialSchoolName   string `json:"initialSchoolName"`
	OwnerEmail          string `json:"ownerEmail"`
	OwnerPhone          string `json:"ownerPhone"`
	OwnerDisplayName    string `json:"ownerDisplayName"`
	OwnerPassword       string `json:"ownerPassword"`
	PlanCode            string `json:"planCode"`
	SubscriptionStatus  string `json:"subscriptionStatus"`
	TrialEndsAt         string `json:"trialEndsAt"`
	CurrentPeriodEndsAt string `json:"currentPeriodEndsAt"`
}

func handleTenantIntakeSubmit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input tenantIntakeSubmitInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeTenantIntakeSubmitInput(input)
	if err := validateTenantIntakeSubmitInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	intake, err := createTenantIntakeRequest(r.Context(), db, input)
	if err != nil {
		http.Error(w, "cannot save intake request", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"ok": true, "intake": intake})
}

func handlePlatformIntakeRequests(w http.ResponseWriter, r *http.Request) {
	status := headerKey(r.URL.Query().Get("status"))
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	requests, err := listTenantIntakeRequests(r.Context(), db, status)
	if err != nil {
		http.Error(w, "cannot load intake requests", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"requests": requests})
}

func handlePlatformIntakeStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input tenantIntakeStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.ID = strings.TrimSpace(input.ID)
	input.Status = headerKey(input.Status)
	if input.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if !isTenantIntakeStatus(input.Status) {
		http.Error(w, "status must be new, contacted, converted, or closed", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	request, err := updateTenantIntakeStatus(r.Context(), db, input.ID, input.Status, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"request": request})
}

func handlePlatformTenantOnboard(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input platformTenantOnboardInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizePlatformTenantOnboardInput(input)
	if err := validatePlatformTenantOnboardInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	tenant, owner, err := createPlatformOnboardedTenant(r.Context(), db, input, user, auditContextFromRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenants, err := listAllTenants(r.Context(), db, tenant.ID)
	if err != nil {
		http.Error(w, "cannot reload tenants", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"tenant":  tenant,
		"owner":   owner,
		"tenants": tenants,
	})
}

func normalizeTenantIntakeSubmitInput(input tenantIntakeSubmitInput) tenantIntakeSubmitInput {
	input.SchoolName = strings.TrimSpace(input.SchoolName)
	input.ContactName = strings.TrimSpace(input.ContactName)
	input.ContactEmail = normalizeAuthEmail(input.ContactEmail)
	input.ContactPhone = normalizeAdminPhone(input.ContactPhone)
	input.DesiredPlanCode = headerKey(input.DesiredPlanCode)
	input.Note = strings.TrimSpace(input.Note)
	return input
}

func validateTenantIntakeSubmitInput(input tenantIntakeSubmitInput) error {
	if input.SchoolName == "" {
		return fmt.Errorf("schoolName is required")
	}
	if input.ContactEmail == "" && input.ContactPhone == "" {
		return fmt.Errorf("contact email or phone is required")
	}
	if input.ContactPhone != "" {
		if err := validateAdminPhone(input.ContactPhone); err != nil {
			return err
		}
	}
	return nil
}

func createTenantIntakeRequest(ctx context.Context, db *sql.DB, input tenantIntakeSubmitInput) (tenantIntakeRequestSummary, error) {
	var request tenantIntakeRequestSummary
	err := db.QueryRowContext(ctx, `
INSERT INTO tenant_intake_requests (
	school_name, contact_name, contact_email, contact_phone, desired_plan_code, note
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id::text, school_name, contact_name, contact_email, contact_phone, desired_plan_code, note, status, created_at, updated_at`,
		input.SchoolName,
		input.ContactName,
		input.ContactEmail,
		input.ContactPhone,
		input.DesiredPlanCode,
		input.Note,
	).Scan(
		&request.ID,
		&request.SchoolName,
		&request.ContactName,
		&request.ContactEmail,
		&request.ContactPhone,
		&request.DesiredPlanCode,
		&request.Note,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	return request, err
}

func listTenantIntakeRequests(ctx context.Context, db *sql.DB, status string) ([]tenantIntakeRequestSummary, error) {
	status = headerKey(status)
	rows, err := db.QueryContext(ctx, `
SELECT id::text,
	school_name,
	contact_name,
	contact_email,
	contact_phone,
	desired_plan_code,
	note,
	status,
	created_at,
	updated_at,
	handled_at
FROM tenant_intake_requests
WHERE ($1 = '' OR status = $1)
ORDER BY created_at DESC
LIMIT 100`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []tenantIntakeRequestSummary{}
	for rows.Next() {
		var request tenantIntakeRequestSummary
		var handledAt sql.NullTime
		if err := rows.Scan(
			&request.ID,
			&request.SchoolName,
			&request.ContactName,
			&request.ContactEmail,
			&request.ContactPhone,
			&request.DesiredPlanCode,
			&request.Note,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
			&handledAt,
		); err != nil {
			return nil, err
		}
		request.HandledAt = formatNullTimestamp(handledAt)
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

func updateTenantIntakeStatus(ctx context.Context, db *sql.DB, id string, status string, userID string) (tenantIntakeRequestSummary, error) {
	var request tenantIntakeRequestSummary
	var handledAt sql.NullTime
	err := db.QueryRowContext(ctx, `
UPDATE tenant_intake_requests
SET status = $2,
	handled_by_user_id = CASE WHEN $2 IN ('converted', 'closed') THEN nullif($3, '')::uuid ELSE handled_by_user_id END,
	handled_at = CASE WHEN $2 IN ('converted', 'closed') THEN now() ELSE handled_at END
WHERE id = $1::uuid
RETURNING id::text, school_name, contact_name, contact_email, contact_phone, desired_plan_code, note, status, created_at, updated_at, handled_at`,
		id,
		status,
		userID,
	).Scan(
		&request.ID,
		&request.SchoolName,
		&request.ContactName,
		&request.ContactEmail,
		&request.ContactPhone,
		&request.DesiredPlanCode,
		&request.Note,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
		&handledAt,
	)
	if err == sql.ErrNoRows {
		return request, fmt.Errorf("intake request not found")
	}
	request.HandledAt = formatNullTimestamp(handledAt)
	return request, err
}

func normalizePlatformTenantOnboardInput(input platformTenantOnboardInput) platformTenantOnboardInput {
	input.IntakeRequestID = strings.TrimSpace(input.IntakeRequestID)
	input.TenantCode = normalizeSchoolCode(input.TenantCode)
	input.TenantName = strings.TrimSpace(input.TenantName)
	input.TenantStatus = headerKey(input.TenantStatus)
	if input.TenantStatus == "" {
		input.TenantStatus = "active"
	}
	input.InitialSchoolCode = normalizeSchoolCode(input.InitialSchoolCode)
	if input.InitialSchoolCode == "" {
		input.InitialSchoolCode = input.TenantCode
	}
	input.InitialSchoolName = strings.TrimSpace(input.InitialSchoolName)
	if input.InitialSchoolName == "" {
		input.InitialSchoolName = input.TenantName
	}
	input.OwnerEmail = normalizeAuthEmail(input.OwnerEmail)
	input.OwnerPhone = normalizeAdminPhone(input.OwnerPhone)
	input.OwnerDisplayName = strings.TrimSpace(input.OwnerDisplayName)
	input.OwnerPassword = strings.TrimSpace(input.OwnerPassword)
	input.PlanCode = headerKey(input.PlanCode)
	if input.PlanCode == "" {
		input.PlanCode = subscriptionPlanStandard
	}
	input.SubscriptionStatus = headerKey(input.SubscriptionStatus)
	if input.SubscriptionStatus == "" {
		input.SubscriptionStatus = defaultSubscriptionStatusForTenantStatus(input.TenantStatus)
	}
	input.TrialEndsAt = strings.TrimSpace(input.TrialEndsAt)
	input.CurrentPeriodEndsAt = strings.TrimSpace(input.CurrentPeriodEndsAt)
	return input
}

func validatePlatformTenantOnboardInput(input platformTenantOnboardInput) error {
	tenantInput := tenantSaveInput{
		Code:              input.TenantCode,
		Name:              input.TenantName,
		Status:            input.TenantStatus,
		InitialSchoolCode: input.InitialSchoolCode,
		InitialSchoolName: input.InitialSchoolName,
	}
	if err := validateTenantSaveInput(tenantInput); err != nil {
		return err
	}
	userInput := adminUserSaveInput{
		Email:       input.OwnerEmail,
		Phone:       input.OwnerPhone,
		DisplayName: input.OwnerDisplayName,
		Status:      "active",
		Password:    input.OwnerPassword,
	}
	if err := validateAdminUserSaveInput(&userInput); err != nil {
		return err
	}
	if !isTenantSubscriptionStatus(input.SubscriptionStatus) {
		return fmt.Errorf("subscription status must be trial, active, past_due, suspended, or cancelled")
	}
	if _, err := parseOptionalDate(input.TrialEndsAt); err != nil {
		return fmt.Errorf("trialEndsAt must use YYYY-MM-DD")
	}
	if _, err := parseOptionalDate(input.CurrentPeriodEndsAt); err != nil {
		return fmt.Errorf("currentPeriodEndsAt must use YYYY-MM-DD")
	}
	return nil
}

func createPlatformOnboardedTenant(ctx context.Context, db *sql.DB, input platformTenantOnboardInput, platformUser authenticatedUser, auditCtx requestAuditContext) (tenantSummary, adminUserSummary, error) {
	passwordHash, err := hashPassword(input.OwnerPassword)
	if err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}
	trialEndsAt, err := parseOptionalDate(input.TrialEndsAt)
	if err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}
	currentPeriodEndsAt, err := parseOptionalDate(input.CurrentPeriodEndsAt)
	if err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}
	defer tx.Rollback()

	var existingUserID string
	err = tx.QueryRowContext(ctx, `
SELECT id::text
FROM app_users
WHERE ($1 <> '' AND lower(COALESCE(email, '')) = lower($1))
	OR ($2 <> '' AND phone = $2)
LIMIT 1`, input.OwnerEmail, input.OwnerPhone).Scan(&existingUserID)
	if err == nil && existingUserID != "" {
		return tenantSummary{}, adminUserSummary{}, fmt.Errorf("owner email or phone already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return tenantSummary{}, adminUserSummary{}, err
	}

	var owner adminUserSummary
	if err := tx.QueryRowContext(ctx, `
INSERT INTO app_users (email, phone, display_name, status, password_hash, password_updated_at)
VALUES ($1, $2, $3, 'active', $4, now())
RETURNING id::text, COALESCE(email, ''), phone, display_name, status, password_hash <> '', created_at, updated_at`,
		input.OwnerEmail,
		input.OwnerPhone,
		input.OwnerDisplayName,
		passwordHash,
	).Scan(
		&owner.ID,
		&owner.Email,
		&owner.Phone,
		&owner.DisplayName,
		&owner.Status,
		&owner.HasPassword,
		&owner.CreatedAt,
		&owner.UpdatedAt,
	); err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}

	tenantInput := tenantSaveInput{
		Code:              input.TenantCode,
		Name:              input.TenantName,
		Status:            input.TenantStatus,
		InitialSchoolCode: input.InitialSchoolCode,
		InitialSchoolName: input.InitialSchoolName,
	}
	ownerAuth := authenticatedUser{
		ID:          owner.ID,
		Email:       owner.Email,
		Phone:       owner.Phone,
		DisplayName: owner.DisplayName,
		Status:      owner.Status,
	}
	tenant, err := createTenantWithInitialSchool(ctx, tx, tenantInput, ownerAuth, auditCtx)
	if err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}

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
		tenant.ID,
		input.PlanCode,
		input.SubscriptionStatus,
		trialEndsAt,
		currentPeriodEndsAt,
		platformUser.ID,
	)
	if err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return tenantSummary{}, adminUserSummary{}, fmt.Errorf("subscription plan %q is not available", input.PlanCode)
	}

	if input.IntakeRequestID != "" {
		if _, err := tx.ExecContext(ctx, `
UPDATE tenant_intake_requests
SET status = 'converted',
	handled_by_user_id = nullif($2, '')::uuid,
	handled_at = now()
WHERE id = $1::uuid`, input.IntakeRequestID, platformUser.ID); err != nil {
			return tenantSummary{}, adminUserSummary{}, err
		}
	}
	auditCtx.TenantID = tenant.ID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "platform.tenant.onboard",
		EntityType: "tenant",
		EntityID:   tenant.ID,
		Metadata: map[string]any{
			"tenantCode":         tenant.Code,
			"ownerUserId":        owner.ID,
			"planCode":           input.PlanCode,
			"subscriptionStatus": input.SubscriptionStatus,
		},
	})

	if err := tx.Commit(); err != nil {
		return tenantSummary{}, adminUserSummary{}, err
	}

	if tenants, err := listAllTenants(ctx, db, tenant.ID); err == nil {
		for _, item := range tenants {
			if item.ID == tenant.ID {
				tenant = item
				break
			}
		}
	}
	owner.Roles = []adminRoleSummary{{Code: "tenant_owner", Name: "Tenant Owner"}}
	return tenant, owner, nil
}

func isTenantIntakeStatus(status string) bool {
	switch status {
	case "new", "contacted", "converted", "closed":
		return true
	default:
		return false
	}
}

func formatNullTimestamp(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format(time.RFC3339)
}
