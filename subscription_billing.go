package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	subscriptionInvoiceStatusDraft   = "draft"
	subscriptionInvoiceStatusOpen    = "open"
	subscriptionInvoiceStatusPaid    = "paid"
	subscriptionInvoiceStatusPastDue = "past_due"
	subscriptionInvoiceStatusVoid    = "void"
)

type subscriptionInvoiceSummary struct {
	ID             string `json:"id"`
	InvoiceCode    string `json:"invoiceCode"`
	PlanCode       string `json:"planCode"`
	PlanName       string `json:"planName"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	PeriodStartsAt string `json:"periodStartsAt"`
	PeriodEndsAt   string `json:"periodEndsAt"`
	DueAt          string `json:"dueAt"`
	Status         string `json:"status"`
	PaidAt         string `json:"paidAt,omitempty"`
	DunningCount   int    `json:"dunningCount"`
	LastDunningAt  string `json:"lastDunningAt,omitempty"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type subscriptionBillingSummary struct {
	OpenCount         int    `json:"openCount"`
	PastDueCount      int    `json:"pastDueCount"`
	PaidCount         int    `json:"paidCount"`
	LatestStatus      string `json:"latestStatus,omitempty"`
	LatestDueAt       string `json:"latestDueAt,omitempty"`
	LatestInvoiceCode string `json:"latestInvoiceCode,omitempty"`
}

type subscriptionBillingPeriodSuggestion struct {
	PeriodStartsAt string `json:"periodStartsAt"`
	PeriodEndsAt   string `json:"periodEndsAt"`
	DueAt          string `json:"dueAt"`
	Amount         int    `json:"amount"`
}

type subscriptionBillingResponse struct {
	Invoices        []subscriptionInvoiceSummary        `json:"invoices"`
	Summary         subscriptionBillingSummary          `json:"summary"`
	SuggestedPeriod subscriptionBillingPeriodSuggestion `json:"suggestedPeriod"`
	Config          subscriptionBillingConfig           `json:"config"`
	Tenant          tenantSummary                       `json:"tenant"`
	Tenants         []tenantSummary                     `json:"tenants,omitempty"`
	DunningResults  []subscriptionDunningResult         `json:"dunningResults,omitempty"`
}

type subscriptionInvoiceGenerateInput struct {
	TenantID       string `json:"tenantId"`
	PeriodStartsAt string `json:"periodStartsAt"`
	PeriodEndsAt   string `json:"periodEndsAt"`
	DueAt          string `json:"dueAt"`
	Amount         int    `json:"amount"`
}

type subscriptionInvoiceMarkPaidInput struct {
	TenantID    string `json:"tenantId"`
	InvoiceID   string `json:"invoiceId"`
	PaidAt      string `json:"paidAt"`
	PaymentNote string `json:"paymentNote"`
}

type subscriptionDunningRunInput struct {
	TenantID    string `json:"tenantId"`
	DryRun      bool   `json:"dryRun"`
	ConfirmSend bool   `json:"confirmSend"`
}

type subscriptionDunningResult struct {
	InvoiceID      string            `json:"invoiceId"`
	InvoiceCode    string            `json:"invoiceCode"`
	RecipientCount int               `json:"recipientCount"`
	Results        []emailSendResult `json:"results"`
}

type subscriptionBillingConfig struct {
	Amount         int    `json:"amount"`
	IntervalMonths int    `json:"intervalMonths"`
	DueDays        int    `json:"dueDays"`
	AutoRenew      bool   `json:"autoRenew"`
	RenewalMode    string `json:"renewalMode"`
	FinanceNote    string `json:"financeNote,omitempty"`
}

type subscriptionBillingConfigSaveInput struct {
	TenantID       string `json:"tenantId"`
	Amount         int    `json:"amount"`
	IntervalMonths int    `json:"intervalMonths"`
	DueDays        int    `json:"dueDays"`
	AutoRenew      bool   `json:"autoRenew"`
	RenewalMode    string `json:"renewalMode"`
	FinanceNote    string `json:"financeNote"`
}

type subscriptionBillingExportFilters struct {
	TenantID   string
	Dataset    string
	Status     string
	PeriodFrom string
	PeriodTo   string
	DueFrom    string
	DueTo      string
}

type tenantSubscriptionBillingProfile struct {
	TenantID              string
	TenantCode            string
	TenantName            string
	SubscriptionID        string
	SubscriptionStatus    string
	PlanCode              string
	PlanName              string
	CurrentPeriodStartsAt sql.NullTime
	CurrentPeriodEndsAt   sql.NullTime
	BillingMetadata       map[string]any
}

type subscriptionDunningRecipient struct {
	UserID      string
	DisplayName string
	Email       string
	RoleCode    string
}

func handleSubscriptionInvoices(w http.ResponseWriter, r *http.Request) {
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

	payload, err := loadSubscriptionBillingResponse(r.Context(), db, user, tenantID)
	if err != nil {
		http.Error(w, "cannot load subscription invoices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func handleSubscriptionInvoiceGenerate(w http.ResponseWriter, r *http.Request) {
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
	var input subscriptionInvoiceGenerateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSubscriptionInvoiceGenerateInput(input, activeTenantID)

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if _, err := generateSubscriptionInvoice(r.Context(), db, input, user, auditContextFromRequest(r)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := loadSubscriptionBillingResponse(r.Context(), db, user, input.TenantID)
	if err != nil {
		http.Error(w, "cannot reload subscription invoices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func handleSubscriptionInvoiceMarkPaid(w http.ResponseWriter, r *http.Request) {
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
	var input subscriptionInvoiceMarkPaidInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSubscriptionInvoiceMarkPaidInput(input, activeTenantID)
	if strings.TrimSpace(input.InvoiceID) == "" {
		http.Error(w, "invoiceId is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if _, err := markSubscriptionInvoicePaid(r.Context(), db, input, user, auditContextFromRequest(r)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := loadSubscriptionBillingResponse(r.Context(), db, user, input.TenantID)
	if err != nil {
		http.Error(w, "cannot reload subscription invoices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func handleSubscriptionDunningRun(w http.ResponseWriter, r *http.Request) {
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
	var input subscriptionDunningRunInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSubscriptionDunningRunInput(input, activeTenantID)
	if !input.DryRun && !input.ConfirmSend {
		http.Error(w, "confirmSend is required for real dunning sends", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	cfg, _ := loadEmailConfig()
	results, err := runSubscriptionDunning(r.Context(), db, user, input, auditContextFromRequest(r), appBaseURL(r, cfg))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := loadSubscriptionBillingResponse(r.Context(), db, user, input.TenantID)
	if err != nil {
		http.Error(w, "cannot reload subscription invoices", http.StatusInternalServerError)
		return
	}
	payload.DunningResults = results
	writeJSON(w, http.StatusOK, payload)
}

func handleSubscriptionBillingConfigSave(w http.ResponseWriter, r *http.Request) {
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
	var input subscriptionBillingConfigSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSubscriptionBillingConfigSaveInput(input, activeTenantID)
	if err := validateSubscriptionBillingConfig(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()
	if err := saveSubscriptionBillingConfig(r.Context(), db, input, user, auditContextFromRequest(r)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := loadSubscriptionBillingResponse(r.Context(), db, user, input.TenantID)
	if err != nil {
		http.Error(w, "cannot reload subscription billing", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func handleSubscriptionBillingExport(w http.ResponseWriter, r *http.Request) {
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
	filters := parseSubscriptionBillingExportFilters(r, tenantID)
	filename, data, err := buildSubscriptionBillingCSV(r.Context(), db, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func loadSubscriptionBillingResponse(ctx context.Context, db *sql.DB, user authenticatedUser, tenantID string) (subscriptionBillingResponse, error) {
	if err := syncSubscriptionInvoicePastDueState(ctx, db, tenantID, time.Now()); err != nil {
		return subscriptionBillingResponse{}, err
	}
	invoices, err := listSubscriptionInvoices(ctx, db, tenantID)
	if err != nil {
		return subscriptionBillingResponse{}, err
	}
	profile, err := loadTenantSubscriptionBillingProfile(ctx, db, tenantID)
	if err != nil {
		return subscriptionBillingResponse{}, err
	}
	tenant, err := loadTenantSummaryByID(ctx, db, tenantID)
	if err != nil {
		return subscriptionBillingResponse{}, err
	}
	tenants, err := listUserTenants(ctx, db, user.ID, tenantID)
	if authenticatedUserHasPermission(user, "operation_log.cross_tenant_view") || authenticatedUserHasPermission(user, "audit_log.cross_tenant_view") {
		tenants, err = listAllTenants(ctx, db, tenantID)
	}
	if err != nil {
		return subscriptionBillingResponse{}, err
	}
	return subscriptionBillingResponse{
		Invoices:        invoices,
		Summary:         summarizeSubscriptionInvoices(invoices),
		SuggestedPeriod: suggestSubscriptionBillingPeriod(profile, time.Now()),
		Config:          subscriptionBillingConfigFromProfile(profile),
		Tenant:          tenant,
		Tenants:         tenants,
	}, nil
}

func normalizeSubscriptionInvoiceGenerateInput(input subscriptionInvoiceGenerateInput, activeTenantID string) subscriptionInvoiceGenerateInput {
	input.TenantID = firstNonEmpty(strings.TrimSpace(input.TenantID), strings.TrimSpace(activeTenantID))
	input.PeriodStartsAt = strings.TrimSpace(input.PeriodStartsAt)
	input.PeriodEndsAt = strings.TrimSpace(input.PeriodEndsAt)
	input.DueAt = strings.TrimSpace(input.DueAt)
	return input
}

func normalizeSubscriptionInvoiceMarkPaidInput(input subscriptionInvoiceMarkPaidInput, activeTenantID string) subscriptionInvoiceMarkPaidInput {
	input.TenantID = firstNonEmpty(strings.TrimSpace(input.TenantID), strings.TrimSpace(activeTenantID))
	input.InvoiceID = strings.TrimSpace(input.InvoiceID)
	input.PaidAt = strings.TrimSpace(input.PaidAt)
	input.PaymentNote = strings.TrimSpace(input.PaymentNote)
	return input
}

func normalizeSubscriptionDunningRunInput(input subscriptionDunningRunInput, activeTenantID string) subscriptionDunningRunInput {
	input.TenantID = firstNonEmpty(strings.TrimSpace(input.TenantID), strings.TrimSpace(activeTenantID))
	if !input.ConfirmSend {
		input.DryRun = true
	}
	return input
}

func normalizeSubscriptionBillingConfigSaveInput(input subscriptionBillingConfigSaveInput, activeTenantID string) subscriptionBillingConfigSaveInput {
	input.TenantID = firstNonEmpty(strings.TrimSpace(input.TenantID), strings.TrimSpace(activeTenantID))
	input.RenewalMode = headerKey(input.RenewalMode)
	input.FinanceNote = strings.TrimSpace(input.FinanceNote)
	return input
}

func validateSubscriptionBillingConfig(input subscriptionBillingConfigSaveInput) error {
	if strings.TrimSpace(input.TenantID) == "" {
		return fmt.Errorf("tenantId is required")
	}
	if input.Amount < 0 {
		return fmt.Errorf("amount must be 0 or greater")
	}
	if input.IntervalMonths < 1 || input.IntervalMonths > 12 {
		return fmt.Errorf("intervalMonths must be between 1 and 12")
	}
	if input.DueDays < 1 || input.DueDays > 28 {
		return fmt.Errorf("dueDays must be between 1 and 28")
	}
	if input.RenewalMode == "" {
		input.RenewalMode = "manual"
	}
	if input.RenewalMode != "manual" && input.RenewalMode != "auto_generate" {
		return fmt.Errorf("renewalMode must be manual or auto_generate")
	}
	return nil
}

func loadTenantSubscriptionBillingProfile(ctx context.Context, exec masterDataExecutor, tenantID string) (tenantSubscriptionBillingProfile, error) {
	var profile tenantSubscriptionBillingProfile
	var billingMetadataBytes []byte
	err := exec.QueryRowContext(ctx, `
SELECT tenant.id::text,
	tenant.code,
	tenant.name,
	ts.id::text,
	ts.status,
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.current_period_starts_at,
	ts.current_period_ends_at,
	COALESCE(ts.billing_metadata, '{}'::jsonb)
FROM tenants tenant
JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE tenant.id = $1::uuid`, tenantID).Scan(
		&profile.TenantID,
		&profile.TenantCode,
		&profile.TenantName,
		&profile.SubscriptionID,
		&profile.SubscriptionStatus,
		&profile.PlanCode,
		&profile.PlanName,
		&profile.CurrentPeriodStartsAt,
		&profile.CurrentPeriodEndsAt,
		&billingMetadataBytes,
	)
	if err != nil {
		return profile, err
	}
	profile.BillingMetadata = decodeMetadata(billingMetadataBytes)
	return profile, nil
}

func subscriptionBillingConfigFromProfile(profile tenantSubscriptionBillingProfile) subscriptionBillingConfig {
	cfg := subscriptionBillingConfig{
		Amount:         parseSubscriptionLimitValue(profile.BillingMetadata["amount"]),
		IntervalMonths: parseSubscriptionLimitValue(profile.BillingMetadata["interval_months"]),
		DueDays:        parseSubscriptionLimitValue(profile.BillingMetadata["due_days"]),
		AutoRenew:      boolValueFromAny(profile.BillingMetadata["auto_renew"]),
		RenewalMode:    strings.TrimSpace(fmt.Sprint(profile.BillingMetadata["renewal_mode"])),
		FinanceNote:    strings.TrimSpace(fmt.Sprint(profile.BillingMetadata["finance_note"])),
	}
	if cfg.IntervalMonths <= 0 {
		cfg.IntervalMonths = 1
	}
	if cfg.DueDays <= 0 {
		cfg.DueDays = 10
	}
	if cfg.RenewalMode == "" || cfg.RenewalMode == "<nil>" {
		cfg.RenewalMode = "manual"
	}
	return cfg
}

func suggestSubscriptionBillingPeriod(profile tenantSubscriptionBillingProfile, now time.Time) subscriptionBillingPeriodSuggestion {
	config := subscriptionBillingConfigFromProfile(profile)
	intervalMonths := config.IntervalMonths
	dueDayOffset := config.DueDays
	start := time.Date(now.UTC().Year(), now.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)
	if profile.CurrentPeriodEndsAt.Valid {
		next := profile.CurrentPeriodEndsAt.Time.UTC().AddDate(0, 0, 1)
		start = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, time.UTC)
	}
	end := start.AddDate(0, intervalMonths, 0).AddDate(0, 0, -1)
	dueAt := start.AddDate(0, 0, dueDayOffset-1)
	return subscriptionBillingPeriodSuggestion{
		PeriodStartsAt: start.Format("2006-01-02"),
		PeriodEndsAt:   end.Format("2006-01-02"),
		DueAt:          dueAt.Format("2006-01-02"),
		Amount:         config.Amount,
	}
}

func boolValueFromAny(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return headerKey(typed) == "true"
	case float64:
		return typed != 0
	case int:
		return typed != 0
	default:
		return false
	}
}

func saveSubscriptionBillingConfig(ctx context.Context, db *sql.DB, input subscriptionBillingConfigSaveInput, user authenticatedUser, auditCtx requestAuditContext) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	metadataBytes, err := json.Marshal(map[string]any{
		"amount":          input.Amount,
		"interval_months": input.IntervalMonths,
		"due_days":        input.DueDays,
		"auto_renew":      input.AutoRenew,
		"renewal_mode":    input.RenewalMode,
		"finance_note":    input.FinanceNote,
	})
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
UPDATE tenant_subscriptions
SET billing_metadata = billing_metadata || $2::jsonb,
	updated_by_user_id = nullif($3, '')::uuid,
	updated_at = now()
WHERE tenant_id = $1::uuid`, input.TenantID, string(metadataBytes), user.ID); err != nil {
		return err
	}
	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "subscription.billing_config.update",
		EntityType: "tenant_subscription",
		EntityID:   input.TenantID,
		Metadata: map[string]any{
			"amount":         input.Amount,
			"intervalMonths": input.IntervalMonths,
			"dueDays":        input.DueDays,
			"autoRenew":      input.AutoRenew,
			"renewalMode":    input.RenewalMode,
			"financeNote":    input.FinanceNote,
		},
	})
	return tx.Commit()
}

func generateSubscriptionInvoice(ctx context.Context, db *sql.DB, input subscriptionInvoiceGenerateInput, user authenticatedUser, auditCtx requestAuditContext) (subscriptionInvoiceSummary, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	defer tx.Rollback()

	profile, err := loadTenantSubscriptionBillingProfile(ctx, tx, input.TenantID)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	if profile.PlanCode == subscriptionPlanFreeTrial {
		return subscriptionInvoiceSummary{}, fmt.Errorf("free_trial tenants do not generate subscription invoices")
	}

	suggested := suggestSubscriptionBillingPeriod(profile, time.Now())
	periodStartsAt, err := parseRequiredBillingDate(firstNonEmpty(input.PeriodStartsAt, suggested.PeriodStartsAt), "periodStartsAt")
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	periodEndsAt, err := parseRequiredBillingDate(firstNonEmpty(input.PeriodEndsAt, suggested.PeriodEndsAt), "periodEndsAt")
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	if periodEndsAt.Before(periodStartsAt) {
		return subscriptionInvoiceSummary{}, fmt.Errorf("periodEndsAt must be on or after periodStartsAt")
	}
	dueAt, err := parseRequiredBillingDate(firstNonEmpty(input.DueAt, suggested.DueAt), "dueAt")
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	amount := input.Amount
	if amount <= 0 {
		amount = suggested.Amount
	}
	if amount <= 0 {
		return subscriptionInvoiceSummary{}, fmt.Errorf("amount is required when subscription billing amount is not configured")
	}

	if existing, err := findSubscriptionInvoiceByPeriod(ctx, tx, profile.SubscriptionID, periodStartsAt, periodEndsAt); err != nil {
		return subscriptionInvoiceSummary{}, err
	} else if existing.ID != "" {
		return existing, tx.Commit()
	}

	status := subscriptionInvoiceStatusOpen
	if dueAt.Before(time.Now().UTC().Truncate(24 * time.Hour)) {
		status = subscriptionInvoiceStatusPastDue
	}
	invoiceCode := buildSubscriptionInvoiceCode(profile.TenantCode, periodStartsAt)
	var invoiceID string
	err = tx.QueryRowContext(ctx, `
INSERT INTO subscription_invoices (
	tenant_id, subscription_id, invoice_code, plan_code, plan_name, amount, currency,
	period_starts_at, period_ends_at, due_at, status, metadata, created_by_user_id, updated_by_user_id
)
VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, 'VND', $7, $8, $9, $10, '{}'::jsonb, nullif($11, '')::uuid, nullif($11, '')::uuid)
RETURNING id::text`,
		input.TenantID,
		profile.SubscriptionID,
		invoiceCode,
		profile.PlanCode,
		profile.PlanName,
		amount,
		periodStartsAt,
		periodEndsAt,
		dueAt,
		status,
		user.ID,
	).Scan(&invoiceID)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	if err := insertSubscriptionInvoiceStatusHistory(ctx, tx, invoiceID, "", status, "invoice generated", user.ID); err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "subscription.invoice.generate",
		EntityType: "subscription_invoice",
		EntityID:   invoiceID,
		Metadata: map[string]any{
			"invoiceCode": invoiceCode,
			"planCode":    profile.PlanCode,
			"amount":      amount,
			"dueAt":       dueAt.Format("2006-01-02"),
		},
	})
	if status == subscriptionInvoiceStatusPastDue {
		if err := updateTenantSubscriptionLifecycleStatus(ctx, tx, input.TenantID, subscriptionStatusPastDue, user.ID, periodStartsAt, periodEndsAt); err != nil {
			return subscriptionInvoiceSummary{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	return loadSubscriptionInvoiceByID(ctx, db, invoiceID, input.TenantID)
}

func markSubscriptionInvoicePaid(ctx context.Context, db *sql.DB, input subscriptionInvoiceMarkPaidInput, user authenticatedUser, auditCtx requestAuditContext) (subscriptionInvoiceSummary, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	defer tx.Rollback()

	invoice, err := loadSubscriptionInvoiceByID(ctx, tx, input.InvoiceID, input.TenantID)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	if invoice.ID == "" {
		return subscriptionInvoiceSummary{}, fmt.Errorf("subscription invoice not found")
	}
	if invoice.Status == subscriptionInvoiceStatusPaid {
		return invoice, tx.Commit()
	}
	paidAt, err := parseRequiredBillingDate(firstNonEmpty(input.PaidAt, time.Now().UTC().Format("2006-01-02")), "paidAt")
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE subscription_invoices
SET status = 'paid',
	paid_at = $2,
	updated_by_user_id = nullif($3, '')::uuid,
	updated_at = now()
WHERE id = $1::uuid`, input.InvoiceID, paidAt, user.ID)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	if err := insertSubscriptionInvoiceStatusHistory(ctx, tx, input.InvoiceID, invoice.Status, subscriptionInvoiceStatusPaid, firstNonEmpty(input.PaymentNote, "marked paid"), user.ID); err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	startAt, _ := time.Parse("2006-01-02", invoice.PeriodStartsAt)
	endAt, _ := time.Parse("2006-01-02", invoice.PeriodEndsAt)
	if err := updateTenantSubscriptionLifecycleStatus(ctx, tx, input.TenantID, subscriptionStatusActive, user.ID, startAt, endAt); err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "subscription.invoice.mark_paid",
		EntityType: "subscription_invoice",
		EntityID:   input.InvoiceID,
		Metadata: map[string]any{
			"invoiceCode": invoice.InvoiceCode,
			"paidAt":      paidAt.Format("2006-01-02"),
			"paymentNote": input.PaymentNote,
		},
	})
	if err := tx.Commit(); err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	return loadSubscriptionInvoiceByID(ctx, db, input.InvoiceID, input.TenantID)
}

func runSubscriptionDunning(ctx context.Context, db *sql.DB, user authenticatedUser, input subscriptionDunningRunInput, auditCtx requestAuditContext, baseURL string) ([]subscriptionDunningResult, error) {
	if err := syncSubscriptionInvoicePastDueState(ctx, db, input.TenantID, time.Now()); err != nil {
		return nil, err
	}
	invoices, err := listSubscriptionDunningCandidates(ctx, db, input.TenantID, time.Now())
	if err != nil {
		return nil, err
	}
	recipients, err := listSubscriptionDunningRecipients(ctx, db, input.TenantID)
	if err != nil {
		return nil, err
	}
	results := make([]subscriptionDunningResult, 0, len(invoices))
	if input.DryRun {
		for _, invoice := range invoices {
			dryRuns := make([]emailSendResult, 0, len(recipients))
			for _, recipient := range recipients {
				dryRuns = append(dryRuns, emailSendResult{
					ID:          recipient.UserID,
					Email:       recipient.Email,
					StudentName: recipient.DisplayName,
					Status:      "dry_run",
				})
			}
			results = append(results, subscriptionDunningResult{
				InvoiceID:      invoice.ID,
				InvoiceCode:    invoice.InvoiceCode,
				RecipientCount: len(recipients),
				Results:        dryRuns,
			})
		}
		return results, nil
	}
	cfg, err := loadEmailConfig()
	if err != nil {
		return nil, err
	}
	if err := validateEmailConfigForSend(cfg); err != nil {
		return nil, err
	}
	quota, err := emailSendQuotaStatus(time.Now())
	if err != nil {
		return nil, err
	}
	if quota.Remaining <= 0 {
		return nil, fmt.Errorf("email daily limit reached (%d/%d in 24h)", quota.Sent, quota.Limit)
	}
	sent := 0
	for _, invoice := range invoices {
		invoiceResult := subscriptionDunningResult{
			InvoiceID:      invoice.ID,
			InvoiceCode:    invoice.InvoiceCode,
			RecipientCount: len(recipients),
			Results:        []emailSendResult{},
		}
		for _, recipient := range recipients {
			if sent >= quota.Remaining {
				invoiceResult.Results = append(invoiceResult.Results, emailSendResult{
					ID:          recipient.UserID,
					Email:       recipient.Email,
					StudentName: recipient.DisplayName,
					Status:      "skipped",
					Error:       "daily email limit reached",
				})
				continue
			}
			email := renderSubscriptionDunningEmail(invoice, recipient, baseURL)
			item := qrItem{paymentRow: paymentRow{
				ID:          invoice.ID,
				StudentName: invoice.PlanName,
				BillNumber:  invoice.InvoiceCode,
				Amount:      invoice.Amount,
				Email:       recipient.Email,
			}}
			outcome, sendErr := sendRenderedEmail(ctx, cfg, item, email)
			sendResult := emailSendResult{
				ID:          recipient.UserID,
				Email:       recipient.Email,
				StudentName: recipient.DisplayName,
				Status:      "sent",
				Provider:    outcome.Provider,
				ResendID:    outcome.ResendID,
				MessageID:   outcome.MessageID,
			}
			if sendErr != nil {
				sendResult.Status = "error"
				sendResult.Error = sendErr.Error()
				sendResult.Transient = isTransientEmailError(sendErr)
			}
			invoiceResult.Results = append(invoiceResult.Results, sendResult)
			if err := insertSubscriptionDunningRun(ctx, db, input.TenantID, invoice, recipient, sendResult); err != nil {
				return nil, err
			}
			if sendResult.Status == "sent" {
				sent++
				sleepEmailSendPace(ctx, cfg)
			} else if sendResult.Status == "error" {
				_ = recordOperationLog(ctx, db, operationLogInput{
					Source:     "email",
					Level:      "error",
					Operation:  "subscription.dunning.send",
					Status:     "error",
					Message:    sendResult.Error,
					EntityType: "subscription_invoice",
					EntityID:   invoice.ID,
					Metadata: map[string]any{
						"invoiceCode":     invoice.InvoiceCode,
						"recipientEmail":  recipient.Email,
						"recipientUserId": recipient.UserID,
					},
				})
			}
		}
		results = append(results, invoiceResult)
	}
	if sent > 0 {
		recordEmailCronSent(sent, time.Now())
	}
	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, db, auditLogInput{
		Context:    auditCtx,
		Action:     "subscription.dunning.run",
		EntityType: "tenant_subscription",
		EntityID:   input.TenantID,
		Metadata: map[string]any{
			"invoiceCount": len(invoices),
			"sentCount":    sent,
		},
	})
	return results, nil
}

func renderSubscriptionDunningEmail(invoice subscriptionInvoiceSummary, recipient subscriptionDunningRecipient, baseURL string) renderedEmail {
	subject := fmt.Sprintf("Subscription payment reminder - %s", invoice.InvoiceCode)
	body := `<html><body style="font-family:Arial,sans-serif;color:#111;background:#fff;">
<div style="max-width:640px;margin:0 auto;padding:24px;">
<p style="margin:0 0 12px;">Xin chào ` + html.EscapeString(firstNonEmpty(recipient.DisplayName, recipient.Email)) + `,</p>
<p style="margin:0 0 12px;">Subscription invoice của tenant đang đến hạn hoặc đã quá hạn.</p>
<table role="presentation" cellpadding="0" cellspacing="0" style="width:100%;border-collapse:collapse;border:1px solid #d4d4d4;font-size:14px;">
<tr><td style="border:1px solid #d4d4d4;padding:8px;">Invoice</td><td style="border:1px solid #d4d4d4;padding:8px;">` + html.EscapeString(invoice.InvoiceCode) + `</td></tr>
<tr><td style="border:1px solid #d4d4d4;padding:8px;">Plan</td><td style="border:1px solid #d4d4d4;padding:8px;">` + html.EscapeString(firstNonEmpty(invoice.PlanName, invoice.PlanCode)) + `</td></tr>
<tr><td style="border:1px solid #d4d4d4;padding:8px;">Period</td><td style="border:1px solid #d4d4d4;padding:8px;">` + html.EscapeString(invoice.PeriodStartsAt) + ` -> ` + html.EscapeString(invoice.PeriodEndsAt) + `</td></tr>
<tr><td style="border:1px solid #d4d4d4;padding:8px;">Due date</td><td style="border:1px solid #d4d4d4;padding:8px;">` + html.EscapeString(invoice.DueAt) + `</td></tr>
<tr><td style="border:1px solid #d4d4d4;padding:8px;">Amount</td><td style="border:1px solid #d4d4d4;padding:8px;">` + html.EscapeString(formatVND(invoice.Amount)) + `</td></tr>
<tr><td style="border:1px solid #d4d4d4;padding:8px;">Status</td><td style="border:1px solid #d4d4d4;padding:8px;font-weight:700;">` + html.EscapeString(invoice.Status) + `</td></tr>
</table>
<p style="margin:16px 0 0;">Vui lòng cập nhật thanh toán trong Web Admin.</p>`
	if strings.TrimSpace(baseURL) != "" {
		body += `<p style="margin:12px 0 0;"><a href="` + html.EscapeString(strings.TrimRight(baseURL, "/")) + `">Mở Web Admin</a></p>`
	}
	body += `</div></body></html>`
	text := fmt.Sprintf("Subscription invoice %s (%s - %s) due %s, amount %s, status %s.",
		invoice.InvoiceCode, invoice.PeriodStartsAt, invoice.PeriodEndsAt, invoice.DueAt, formatVND(invoice.Amount), invoice.Status)
	return renderedEmail{Subject: subject, HTML: body, Text: text}
}

func parseRequiredBillingDate(value string, field string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must use YYYY-MM-DD", field)
	}
	return parsed.UTC(), nil
}

func buildSubscriptionInvoiceCode(tenantCode string, periodStartsAt time.Time) string {
	return fmt.Sprintf("SUB-%s-%s", strings.ToUpper(strings.TrimSpace(tenantCode)), periodStartsAt.UTC().Format("200601"))
}

func insertSubscriptionInvoiceStatusHistory(ctx context.Context, exec masterDataExecutor, invoiceID string, fromStatus string, toStatus string, note string, userID string) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO subscription_invoice_status_history (invoice_id, from_status, to_status, note, created_by_user_id)
VALUES ($1::uuid, $2, $3, $4, nullif($5, '')::uuid)`,
		invoiceID,
		strings.TrimSpace(fromStatus),
		strings.TrimSpace(toStatus),
		strings.TrimSpace(note),
		strings.TrimSpace(userID),
	)
	return err
}

func updateTenantSubscriptionLifecycleStatus(ctx context.Context, exec masterDataExecutor, tenantID string, status string, userID string, currentPeriodStartsAt time.Time, currentPeriodEndsAt time.Time) error {
	_, err := exec.ExecContext(ctx, `
UPDATE tenant_subscriptions
SET status = $2,
	current_period_starts_at = CASE WHEN $3::timestamptz IS NOT NULL THEN $3 ELSE current_period_starts_at END,
	current_period_ends_at = CASE WHEN $4::timestamptz IS NOT NULL THEN $4 ELSE current_period_ends_at END,
	updated_by_user_id = nullif($5, '')::uuid,
	updated_at = now()
WHERE tenant_id = $1::uuid`,
		tenantID,
		status,
		nullableTime(currentPeriodStartsAt),
		nullableTime(currentPeriodEndsAt),
		userID,
	)
	return err
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value.UTC()
}

func syncSubscriptionInvoicePastDueState(ctx context.Context, db *sql.DB, tenantID string, now time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
UPDATE subscription_invoices
SET status = 'past_due',
	updated_at = now()
WHERE tenant_id = $1::uuid
	AND status = 'open'
	AND due_at < $2
RETURNING id::text`, tenantID, now.UTC().Truncate(24*time.Hour))
	if err != nil {
		return err
	}
	updatedIDs := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		updatedIDs = append(updatedIDs, id)
	}
	rows.Close()
	for _, id := range updatedIDs {
		if err := insertSubscriptionInvoiceStatusHistory(ctx, tx, id, subscriptionInvoiceStatusOpen, subscriptionInvoiceStatusPastDue, "due date passed", ""); err != nil {
			return err
		}
	}
	if len(updatedIDs) > 0 {
		if err := updateTenantSubscriptionLifecycleStatus(ctx, tx, tenantID, subscriptionStatusPastDue, "", time.Time{}, time.Time{}); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func listSubscriptionInvoices(ctx context.Context, exec masterDataExecutor, tenantID string) ([]subscriptionInvoiceSummary, error) {
	rows, err := exec.QueryContext(ctx, `
SELECT i.id::text,
	i.invoice_code,
	i.plan_code,
	i.plan_name,
	i.amount,
	i.currency,
	i.period_starts_at,
	i.period_ends_at,
	i.due_at,
	i.status,
	i.paid_at,
	i.created_at,
	i.updated_at,
	COUNT(dr.id) FILTER (WHERE dr.status = 'sent')::integer,
	MAX(dr.created_at) FILTER (WHERE dr.status = 'sent')
FROM subscription_invoices i
LEFT JOIN subscription_dunning_runs dr ON dr.invoice_id = i.id
WHERE i.tenant_id = $1::uuid
GROUP BY i.id
ORDER BY i.created_at DESC
LIMIT 24`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []subscriptionInvoiceSummary{}
	for rows.Next() {
		var item subscriptionInvoiceSummary
		var periodStart time.Time
		var periodEnd time.Time
		var dueAt time.Time
		var paidAt sql.NullTime
		var createdAt time.Time
		var updatedAt time.Time
		var lastDunningAt sql.NullTime
		if err := rows.Scan(
			&item.ID,
			&item.InvoiceCode,
			&item.PlanCode,
			&item.PlanName,
			&item.Amount,
			&item.Currency,
			&periodStart,
			&periodEnd,
			&dueAt,
			&item.Status,
			&paidAt,
			&createdAt,
			&updatedAt,
			&item.DunningCount,
			&lastDunningAt,
		); err != nil {
			return nil, err
		}
		item.PeriodStartsAt = periodStart.UTC().Format("2006-01-02")
		item.PeriodEndsAt = periodEnd.UTC().Format("2006-01-02")
		item.DueAt = dueAt.UTC().Format("2006-01-02")
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		item.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
		if paidAt.Valid {
			item.PaidAt = paidAt.Time.UTC().Format("2006-01-02")
		}
		if lastDunningAt.Valid {
			item.LastDunningAt = lastDunningAt.Time.UTC().Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func summarizeSubscriptionInvoices(invoices []subscriptionInvoiceSummary) subscriptionBillingSummary {
	summary := subscriptionBillingSummary{}
	if len(invoices) > 0 {
		summary.LatestStatus = invoices[0].Status
		summary.LatestDueAt = invoices[0].DueAt
		summary.LatestInvoiceCode = invoices[0].InvoiceCode
	}
	for _, invoice := range invoices {
		switch invoice.Status {
		case subscriptionInvoiceStatusOpen:
			summary.OpenCount++
		case subscriptionInvoiceStatusPastDue:
			summary.PastDueCount++
		case subscriptionInvoiceStatusPaid:
			summary.PaidCount++
		}
	}
	return summary
}

func findSubscriptionInvoiceByPeriod(ctx context.Context, exec masterDataExecutor, subscriptionID string, periodStartsAt time.Time, periodEndsAt time.Time) (subscriptionInvoiceSummary, error) {
	var invoiceID string
	err := exec.QueryRowContext(ctx, `
SELECT id::text
FROM subscription_invoices
WHERE subscription_id = $1::uuid
	AND period_starts_at = $2
	AND period_ends_at = $3
	AND status <> 'void'
LIMIT 1`, subscriptionID, periodStartsAt.UTC(), periodEndsAt.UTC()).Scan(&invoiceID)
	if errors.Is(err, sql.ErrNoRows) {
		return subscriptionInvoiceSummary{}, nil
	}
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	tenantID, err := tenantIDForSubscriptionInvoice(ctx, exec, invoiceID)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	return loadSubscriptionInvoiceByID(ctx, exec, invoiceID, tenantID)
}

func tenantIDForSubscriptionInvoice(ctx context.Context, exec masterDataExecutor, invoiceID string) (string, error) {
	var tenantID string
	err := exec.QueryRowContext(ctx, `SELECT tenant_id::text FROM subscription_invoices WHERE id = $1::uuid`, invoiceID).Scan(&tenantID)
	return tenantID, err
}

func loadSubscriptionInvoiceByID(ctx context.Context, exec masterDataExecutor, invoiceID string, tenantID string) (subscriptionInvoiceSummary, error) {
	items, err := listSubscriptionInvoices(ctx, exec, tenantID)
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	for _, item := range items {
		if item.ID == invoiceID {
			return item, nil
		}
	}
	return subscriptionInvoiceSummary{}, fmt.Errorf("subscription invoice not found")
}

func listSubscriptionDunningCandidates(ctx context.Context, exec masterDataExecutor, tenantID string, now time.Time) ([]subscriptionInvoiceSummary, error) {
	invoices, err := listSubscriptionInvoices(ctx, exec, tenantID)
	if err != nil {
		return nil, err
	}
	candidates := []subscriptionInvoiceSummary{}
	today := now.UTC().Format("2006-01-02")
	for _, invoice := range invoices {
		if (invoice.Status == subscriptionInvoiceStatusOpen || invoice.Status == subscriptionInvoiceStatusPastDue) && invoice.DueAt <= today {
			candidates = append(candidates, invoice)
		}
	}
	return candidates, nil
}

func listSubscriptionDunningRecipients(ctx context.Context, exec masterDataExecutor, tenantID string) ([]subscriptionDunningRecipient, error) {
	rows, err := exec.QueryContext(ctx, `
SELECT DISTINCT user_item.id::text,
	COALESCE(user_item.display_name, ''),
	COALESCE(user_item.email, ''),
	COALESCE(role.code, '')
FROM tenant_memberships membership
JOIN app_users user_item ON user_item.id = membership.user_id
LEFT JOIN tenant_user_roles tur ON tur.tenant_id = membership.tenant_id AND tur.user_id = membership.user_id
LEFT JOIN app_roles role ON role.id = tur.role_id
WHERE membership.tenant_id = $1::uuid
	AND membership.status = 'active'
	AND user_item.status = 'active'
	AND COALESCE(user_item.email, '') <> ''
	AND (
		membership.is_owner
		OR role.code IN ('admin', 'accountant')
	)
ORDER BY role.code, user_item.display_name, user_item.email`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []subscriptionDunningRecipient{}
	seen := map[string]bool{}
	for rows.Next() {
		var item subscriptionDunningRecipient
		if err := rows.Scan(&item.UserID, &item.DisplayName, &item.Email, &item.RoleCode); err != nil {
			return nil, err
		}
		key := strings.ToLower(strings.TrimSpace(item.Email))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, item)
	}
	return items, rows.Err()
}

func insertSubscriptionDunningRun(ctx context.Context, exec masterDataExecutor, tenantID string, invoice subscriptionInvoiceSummary, recipient subscriptionDunningRecipient, result emailSendResult) error {
	status := result.Status
	if status == "" {
		status = "error"
	}
	var sentAt any
	if result.Status == "sent" {
		sentAt = time.Now().UTC()
	}
	_, err := exec.ExecContext(ctx, `
INSERT INTO subscription_dunning_runs (
	tenant_id, invoice_id, recipient_user_id, recipient_email, status, dry_run,
	error_message, provider, provider_message_id, sent_at
)
VALUES (
	$1::uuid, $2::uuid, nullif($3, '')::uuid, $4, $5, false, $6, $7, $8, $9
)`,
		tenantID,
		invoice.ID,
		recipient.UserID,
		recipient.Email,
		status,
		result.Error,
		result.Provider,
		firstNonEmpty(result.MessageID, result.ResendID),
		sentAt,
	)
	return err
}

func parseSubscriptionBillingExportFilters(r *http.Request, tenantID string) subscriptionBillingExportFilters {
	query := r.URL.Query()
	return subscriptionBillingExportFilters{
		TenantID:   tenantID,
		Dataset:    headerKey(firstNonEmpty(query.Get("dataset"), "invoices")),
		Status:     headerKey(query.Get("status")),
		PeriodFrom: strings.TrimSpace(query.Get("periodFrom")),
		PeriodTo:   strings.TrimSpace(query.Get("periodTo")),
		DueFrom:    strings.TrimSpace(query.Get("dueFrom")),
		DueTo:      strings.TrimSpace(query.Get("dueTo")),
	}
}

func buildSubscriptionBillingCSV(ctx context.Context, db *sql.DB, filters subscriptionBillingExportFilters) (string, []byte, error) {
	invoices, err := listSubscriptionInvoices(ctx, db, filters.TenantID)
	if err != nil {
		return "", nil, fmt.Errorf("cannot load subscription invoices")
	}
	filteredInvoices, err := filterSubscriptionInvoicesForExport(invoices, filters)
	if err != nil {
		return "", nil, err
	}
	var records [][]string
	switch filters.Dataset {
	case "invoices":
		records = subscriptionInvoiceCSVRecords(filteredInvoices)
	case "overdue":
		records = subscriptionInvoiceCSVRecords(filterSubscriptionInvoicesByStatus(filteredInvoices, subscriptionInvoiceStatusPastDue, subscriptionInvoiceStatusOpen))
	case "paid":
		records = subscriptionInvoiceCSVRecords(filterSubscriptionInvoicesByStatus(filteredInvoices, subscriptionInvoiceStatusPaid))
	case "dunning":
		runs, err := listSubscriptionDunningRuns(ctx, db, filters.TenantID, filters)
		if err != nil {
			return "", nil, fmt.Errorf("cannot load subscription dunning runs")
		}
		records = subscriptionDunningCSVRecords(runs)
	default:
		return "", nil, fmt.Errorf("unsupported subscription billing dataset %q", filters.Dataset)
	}
	data, err := encodeSubscriptionCSVRecords(records)
	if err != nil {
		return "", nil, err
	}
	filename := fmt.Sprintf("abcsun-subscription-%s-%s.csv", filters.Dataset, time.Now().Format("20060102"))
	return filename, data, nil
}

func filterSubscriptionInvoicesForExport(invoices []subscriptionInvoiceSummary, filters subscriptionBillingExportFilters) ([]subscriptionInvoiceSummary, error) {
	var periodFrom time.Time
	var periodTo time.Time
	var dueFrom time.Time
	var dueTo time.Time
	var err error
	if filters.PeriodFrom != "" {
		periodFrom, err = parseRequiredBillingDate(filters.PeriodFrom, "periodFrom")
		if err != nil {
			return nil, err
		}
	}
	if filters.PeriodTo != "" {
		periodTo, err = parseRequiredBillingDate(filters.PeriodTo, "periodTo")
		if err != nil {
			return nil, err
		}
	}
	if filters.DueFrom != "" {
		dueFrom, err = parseRequiredBillingDate(filters.DueFrom, "dueFrom")
		if err != nil {
			return nil, err
		}
	}
	if filters.DueTo != "" {
		dueTo, err = parseRequiredBillingDate(filters.DueTo, "dueTo")
		if err != nil {
			return nil, err
		}
	}
	filtered := []subscriptionInvoiceSummary{}
	for _, invoice := range invoices {
		if filters.Status != "" && invoice.Status != filters.Status {
			continue
		}
		periodStart, _ := time.Parse("2006-01-02", invoice.PeriodStartsAt)
		dueAt, _ := time.Parse("2006-01-02", invoice.DueAt)
		if !periodFrom.IsZero() && periodStart.Before(periodFrom) {
			continue
		}
		if !periodTo.IsZero() && periodStart.After(periodTo) {
			continue
		}
		if !dueFrom.IsZero() && dueAt.Before(dueFrom) {
			continue
		}
		if !dueTo.IsZero() && dueAt.After(dueTo) {
			continue
		}
		filtered = append(filtered, invoice)
	}
	return filtered, nil
}

func filterSubscriptionInvoicesByStatus(invoices []subscriptionInvoiceSummary, statuses ...string) []subscriptionInvoiceSummary {
	if len(statuses) == 0 {
		return invoices
	}
	allowed := map[string]bool{}
	for _, status := range statuses {
		allowed[status] = true
	}
	filtered := []subscriptionInvoiceSummary{}
	for _, invoice := range invoices {
		if allowed[invoice.Status] {
			filtered = append(filtered, invoice)
		}
	}
	return filtered
}

func subscriptionInvoiceCSVRecords(invoices []subscriptionInvoiceSummary) [][]string {
	records := [][]string{{
		"invoice_code",
		"plan_code",
		"plan_name",
		"period_starts_at",
		"period_ends_at",
		"due_at",
		"status",
		"amount",
		"currency",
		"paid_at",
		"dunning_count",
		"last_dunning_at",
	}}
	for _, invoice := range invoices {
		records = append(records, []string{
			invoice.InvoiceCode,
			invoice.PlanCode,
			invoice.PlanName,
			invoice.PeriodStartsAt,
			invoice.PeriodEndsAt,
			invoice.DueAt,
			invoice.Status,
			strconv.Itoa(invoice.Amount),
			invoice.Currency,
			invoice.PaidAt,
			strconv.Itoa(invoice.DunningCount),
			invoice.LastDunningAt,
		})
	}
	return records
}

type subscriptionDunningRunRow struct {
	InvoiceCode       string
	RecipientEmail    string
	Status            string
	DryRun            bool
	ErrorMessage      string
	Provider          string
	ProviderMessageID string
	SentAt            string
	CreatedAt         string
}

func listSubscriptionDunningRuns(ctx context.Context, exec masterDataExecutor, tenantID string, filters subscriptionBillingExportFilters) ([]subscriptionDunningRunRow, error) {
	query := `
SELECT i.invoice_code,
	dr.recipient_email,
	dr.status,
	dr.dry_run,
	dr.error_message,
	dr.provider,
	dr.provider_message_id,
	dr.sent_at,
	dr.created_at
FROM subscription_dunning_runs dr
JOIN subscription_invoices i ON i.id = dr.invoice_id
WHERE dr.tenant_id = $1::uuid`
	args := []any{tenantID}
	if filters.Status != "" {
		args = append(args, filters.Status)
		query += fmt.Sprintf(" AND dr.status = $%d", len(args))
	}
	query += " ORDER BY dr.created_at DESC"
	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []subscriptionDunningRunRow{}
	for rows.Next() {
		var item subscriptionDunningRunRow
		var sentAt sql.NullTime
		var createdAt time.Time
		if err := rows.Scan(&item.InvoiceCode, &item.RecipientEmail, &item.Status, &item.DryRun, &item.ErrorMessage, &item.Provider, &item.ProviderMessageID, &sentAt, &createdAt); err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		if sentAt.Valid {
			item.SentAt = sentAt.Time.UTC().Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func subscriptionDunningCSVRecords(rows []subscriptionDunningRunRow) [][]string {
	records := [][]string{{
		"invoice_code",
		"recipient_email",
		"status",
		"dry_run",
		"error_message",
		"provider",
		"provider_message_id",
		"sent_at",
		"created_at",
	}}
	for _, row := range rows {
		records = append(records, []string{
			row.InvoiceCode,
			row.RecipientEmail,
			row.Status,
			strconv.FormatBool(row.DryRun),
			row.ErrorMessage,
			row.Provider,
			row.ProviderMessageID,
			row.SentAt,
			row.CreatedAt,
		})
	}
	return records
}

func encodeSubscriptionCSVRecords(records [][]string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(&buf)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}
