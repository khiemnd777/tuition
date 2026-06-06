package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	subscriptionAutomationScopeActive      = "active"
	subscriptionAutomationScopeAll         = "all"
	subscriptionAutomationTriggerManual    = "manual"
	subscriptionAutomationTriggerScheduler = "scheduler"

	subscriptionAutomationStatusDryRun  = "dry_run"
	subscriptionAutomationStatusSuccess = "success"
	subscriptionAutomationStatusPartial = "partial"
	subscriptionAutomationStatusError   = "error"

	defaultSubscriptionAutomationInterval = 5 * time.Minute
)

type subscriptionAutomationStatusResponse struct {
	Scope     string                                `json:"scope"`
	Scheduler subscriptionAutomationSchedulerStatus `json:"scheduler"`
	LatestRun subscriptionAutomationRunRecord       `json:"latestRun"`
}

type subscriptionAutomationSchedulerStatus struct {
	Enabled         bool   `json:"enabled"`
	Interval        string `json:"interval"`
	IntervalSeconds int    `json:"intervalSeconds"`
}

type subscriptionAutomationRunRecord struct {
	ID            string                           `json:"id,omitempty"`
	TenantID      string                           `json:"tenantId,omitempty"`
	TenantCode    string                           `json:"tenantCode,omitempty"`
	Scope         string                           `json:"scope"`
	TriggerSource string                           `json:"triggerSource"`
	Status        string                           `json:"status"`
	DryRun        bool                             `json:"dryRun"`
	Summary       subscriptionAutomationRunSummary `json:"summary"`
	ErrorMessage  string                           `json:"errorMessage,omitempty"`
	StartedAt     string                           `json:"startedAt,omitempty"`
	FinishedAt    string                           `json:"finishedAt,omitempty"`
}

type subscriptionAutomationRunSummary struct {
	TenantsEvaluated      int `json:"tenantsEvaluated"`
	TenantsEnabled        int `json:"tenantsEnabled"`
	RenewalPreviewCount   int `json:"renewalPreviewCount"`
	RenewalGeneratedCount int `json:"renewalGeneratedCount"`
	DunningInvoiceCount   int `json:"dunningInvoiceCount"`
	DunningRecipientCount int `json:"dunningRecipientCount"`
	DunningSentCount      int `json:"dunningSentCount"`
	SuspendPreviewCount   int `json:"suspendPreviewCount"`
	SuspendedTenantCount  int `json:"suspendedTenantCount"`
	ErrorCount            int `json:"errorCount"`
}

func handleSubscriptionAutomationStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	activeTenantID := activeTenantIDFromRequest(r)
	scope, ok := resolveSubscriptionFinanceScope(w, r, user, activeTenantID)
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	payload, err := loadSubscriptionAutomationStatus(r.Context(), db, scope, activeTenantID)
	if err != nil {
		http.Error(w, "cannot load subscription automation status", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func handleSubscriptionAutomationRun(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	activeTenantID := activeTenantIDFromRequest(r)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input subscriptionBatchRunInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSubscriptionBatchRunInput(input)
	scope, ok := resolveSubscriptionFinanceScopeFromValue(w, user, activeTenantID, input.Scope)
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	emailCfg, _ := loadEmailConfig()
	run, err := runSubscriptionAutomationCycle(
		r.Context(),
		db,
		user,
		scope,
		input.DryRun,
		subscriptionAutomationTriggerManual,
		auditContextFromRequest(r),
		schedulerBaseURL(emailCfg),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, subscriptionAutomationStatusResponse{
		Scope:     scopeValueForResponse(scope),
		Scheduler: subscriptionAutomationSchedulerStatusFromEnv(),
		LatestRun: run,
	})
}

func startSubscriptionAutomationScheduler(ctx context.Context) {
	settings := subscriptionAutomationSchedulerStatusFromEnv()
	if !settings.Enabled {
		return
	}
	ticker := time.NewTicker(time.Duration(settings.IntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := runSubscriptionAutomationSchedulerOnce(ctx); err != nil {
				recordOperationLogBestEffort(ctx, operationLogInput{
					Source:    "background_job",
					Level:     "error",
					Operation: "subscription.automation.run",
					Status:    "error",
					Message:   err.Error(),
					Metadata: map[string]any{
						"trigger": subscriptionAutomationTriggerScheduler,
					},
				})
			}
		}
	}
}

func runSubscriptionAutomationSchedulerOnce(ctx context.Context) error {
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return err
	}
	defer db.Close()
	emailCfg, _ := loadEmailConfig()
	_, err = runSubscriptionAutomationCycle(
		ctx,
		db,
		authenticatedUser{},
		subscriptionAutomationScopeAll,
		false,
		subscriptionAutomationTriggerScheduler,
		requestAuditContext{},
		schedulerBaseURL(emailCfg),
	)
	return err
}

func subscriptionAutomationSchedulerStatusFromEnv() subscriptionAutomationSchedulerStatus {
	interval := defaultSubscriptionAutomationInterval
	if raw := strings.TrimSpace(os.Getenv("ABC_SUBSCRIPTION_AUTOMATION_INTERVAL")); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed >= time.Minute {
			interval = parsed
		}
	}
	return subscriptionAutomationSchedulerStatus{
		Enabled:         parseBoolWithDefault(os.Getenv("ABC_SUBSCRIPTION_AUTOMATION_ENABLED"), false),
		Interval:        interval.String(),
		IntervalSeconds: int(interval / time.Second),
	}
}

func loadSubscriptionAutomationStatus(ctx context.Context, db *sql.DB, scope string, activeTenantID string) (subscriptionAutomationStatusResponse, error) {
	run, err := loadLatestSubscriptionAutomationRun(ctx, db, scope, activeTenantID)
	if err != nil {
		return subscriptionAutomationStatusResponse{}, err
	}
	return subscriptionAutomationStatusResponse{
		Scope:     scopeValueForResponse(scope),
		Scheduler: subscriptionAutomationSchedulerStatusFromEnv(),
		LatestRun: run,
	}, nil
}

func runSubscriptionAutomationCycle(ctx context.Context, db *sql.DB, user authenticatedUser, scope string, dryRun bool, trigger string, auditCtx requestAuditContext, baseURL string) (subscriptionAutomationRunRecord, error) {
	startedAt := time.Now().UTC()
	rows, err := listSubscriptionFinanceConsoleRows(ctx, db, subscriptionFinanceConsoleFilters{Scope: scope})
	if err != nil {
		return subscriptionAutomationRunRecord{}, err
	}
	run := subscriptionAutomationRunRecord{
		Scope:         scopeValueForResponse(scope),
		TriggerSource: trigger,
		DryRun:        dryRun,
		StartedAt:     startedAt.Format(time.RFC3339),
		Status:        subscriptionAutomationStatusDryRun,
	}
	var runErrors []string
	now := time.Now().UTC()

	for _, row := range rows {
		run.Summary.TenantsEvaluated++
		if err := syncSubscriptionInvoicePastDueState(ctx, db, row.TenantID, now); err != nil {
			runErrors = append(runErrors, fmt.Sprintf("%s sync: %v", row.TenantCode, err))
			run.Summary.ErrorCount++
			continue
		}
		profile, err := loadTenantSubscriptionBillingProfile(ctx, db, row.TenantID)
		if err != nil {
			runErrors = append(runErrors, fmt.Sprintf("%s profile: %v", row.TenantCode, err))
			run.Summary.ErrorCount++
			continue
		}
		config := subscriptionBillingConfigFromProfile(profile)
		if !config.AutomationEnabled {
			continue
		}
		run.Summary.TenantsEnabled++

		if shouldGenerateSubscriptionRenewal(row, config, now) {
			if dryRun {
				run.Summary.RenewalPreviewCount++
			} else {
				invoice, generateErr := generateSubscriptionInvoice(ctx, db, subscriptionInvoiceGenerateInput{
					TenantID:       row.TenantID,
					PeriodStartsAt: row.NextPeriodStartsAt,
					PeriodEndsAt:   row.NextPeriodEndsAt,
					DueAt:          row.NextDueAt,
					Amount:         row.Amount,
				}, user, auditCtx)
				if generateErr != nil {
					runErrors = append(runErrors, fmt.Sprintf("%s renewal: %v", row.TenantCode, generateErr))
					run.Summary.ErrorCount++
				} else if invoice.ID != "" {
					run.Summary.RenewalGeneratedCount++
				}
			}
		}

		if config.DunningEnabled {
			candidates, candidateErr := listSubscriptionAutomationDunningCandidates(ctx, db, row.TenantID, config.DunningIntervalDays, now)
			if candidateErr != nil {
				runErrors = append(runErrors, fmt.Sprintf("%s dunning candidates: %v", row.TenantCode, candidateErr))
				run.Summary.ErrorCount++
			} else if len(candidates) > 0 {
				run.Summary.DunningInvoiceCount += len(candidates)
				results, dunningErr := runSubscriptionDunningForInvoices(ctx, db, user, subscriptionDunningRunInput{
					TenantID:    row.TenantID,
					DryRun:      dryRun,
					ConfirmSend: !dryRun,
				}, candidates, auditCtx, baseURL)
				if dunningErr != nil {
					runErrors = append(runErrors, fmt.Sprintf("%s dunning: %v", row.TenantCode, dunningErr))
					run.Summary.ErrorCount++
				} else {
					run.Summary.DunningRecipientCount += subscriptionDunningRecipientCount(results)
					run.Summary.DunningSentCount += subscriptionDunningSentCount(results)
				}
			}
		}

		if config.SuspendEnabled && row.SubscriptionStatus != subscriptionStatusSuspended {
			shouldSuspend, invoiceCode, suspendErr := shouldSuspendSubscriptionTenant(ctx, db, row.TenantID, config.SuspendAfterDays, now)
			if suspendErr != nil {
				runErrors = append(runErrors, fmt.Sprintf("%s suspend-check: %v", row.TenantCode, suspendErr))
				run.Summary.ErrorCount++
			} else if shouldSuspend {
				if dryRun {
					run.Summary.SuspendPreviewCount++
				} else {
					if err := updateTenantSubscriptionLifecycleStatus(ctx, db, row.TenantID, subscriptionStatusSuspended, user.ID, time.Time{}, time.Time{}); err != nil {
						runErrors = append(runErrors, fmt.Sprintf("%s suspend: %v", row.TenantCode, err))
						run.Summary.ErrorCount++
					} else {
						run.Summary.SuspendedTenantCount++
						auditCtx.TenantID = row.TenantID
						_ = insertAuditLog(ctx, db, auditLogInput{
							Context:    auditCtx,
							Action:     "subscription.automation.suspend",
							EntityType: "tenant_subscription",
							EntityID:   row.TenantID,
							Metadata: map[string]any{
								"invoiceCode":      invoiceCode,
								"suspendAfterDays": config.SuspendAfterDays,
								"triggerSource":    trigger,
							},
						})
					}
				}
			}
		}
	}

	run.Status = subscriptionAutomationStatusFromSummary(run.Summary, dryRun)
	if len(runErrors) > 0 {
		run.ErrorMessage = strings.Join(runErrors, "; ")
	}
	finishedAt := time.Now().UTC()
	run.FinishedAt = finishedAt.Format(time.RFC3339)
	runID, err := insertSubscriptionAutomationRun(ctx, db, run, scope, activeTenantIDForScope(scope), user.ID, startedAt, finishedAt)
	if err != nil {
		return subscriptionAutomationRunRecord{}, err
	}
	run.ID = runID
	return run, nil
}

func shouldGenerateSubscriptionRenewal(row subscriptionFinanceConsoleRow, config subscriptionBillingConfig, now time.Time) bool {
	if !config.AutoRenew || row.RenewalMode != "auto_generate" || row.MissingConfig {
		return false
	}
	switch row.SubscriptionStatus {
	case subscriptionStatusActive, subscriptionStatusTrial, subscriptionStatusPastDue:
	default:
		return false
	}
	startAt, err := parseRequiredBillingDate(row.NextPeriodStartsAt, "nextPeriodStartsAt")
	if err != nil {
		return false
	}
	leadStart := startAt.AddDate(0, 0, -config.RenewalLeadDays)
	return !now.Truncate(24 * time.Hour).Before(leadStart)
}

func listSubscriptionAutomationDunningCandidates(ctx context.Context, exec masterDataExecutor, tenantID string, intervalDays int, now time.Time) ([]subscriptionInvoiceSummary, error) {
	invoices, err := listSubscriptionInvoices(ctx, exec, tenantID)
	if err != nil {
		return nil, err
	}
	candidates := []subscriptionInvoiceSummary{}
	today := now.UTC().Format("2006-01-02")
	cutoff := now.UTC().AddDate(0, 0, -intervalDays)
	for _, invoice := range invoices {
		if invoice.Status != subscriptionInvoiceStatusOpen && invoice.Status != subscriptionInvoiceStatusPastDue {
			continue
		}
		if invoice.DueAt == "" || invoice.DueAt > today {
			continue
		}
		if invoice.LastDunningAt != "" {
			lastRunAt, err := time.Parse(time.RFC3339, invoice.LastDunningAt)
			if err == nil && lastRunAt.After(cutoff) {
				continue
			}
		}
		candidates = append(candidates, invoice)
	}
	return candidates, nil
}

func shouldSuspendSubscriptionTenant(ctx context.Context, exec masterDataExecutor, tenantID string, suspendAfterDays int, now time.Time) (bool, string, error) {
	invoices, err := listSubscriptionInvoices(ctx, exec, tenantID)
	if err != nil {
		return false, "", err
	}
	cutoff := now.UTC().AddDate(0, 0, -suspendAfterDays).Format("2006-01-02")
	for _, invoice := range invoices {
		if invoice.Status != subscriptionInvoiceStatusOpen && invoice.Status != subscriptionInvoiceStatusPastDue {
			continue
		}
		if invoice.DueAt != "" && invoice.DueAt <= cutoff {
			return true, invoice.InvoiceCode, nil
		}
	}
	return false, "", nil
}

func subscriptionDunningRecipientCount(results []subscriptionDunningResult) int {
	count := 0
	for _, result := range results {
		count += result.RecipientCount
	}
	return count
}

func subscriptionDunningSentCount(results []subscriptionDunningResult) int {
	count := 0
	for _, result := range results {
		for _, item := range result.Results {
			if item.Status == "sent" {
				count++
			}
		}
	}
	return count
}

func subscriptionAutomationStatusFromSummary(summary subscriptionAutomationRunSummary, dryRun bool) string {
	if dryRun && summary.ErrorCount == 0 {
		return subscriptionAutomationStatusDryRun
	}
	if summary.ErrorCount == 0 {
		return subscriptionAutomationStatusSuccess
	}
	workDone := summary.RenewalGeneratedCount > 0 || summary.DunningSentCount > 0 || summary.SuspendedTenantCount > 0 || summary.RenewalPreviewCount > 0 || summary.SuspendPreviewCount > 0 || summary.DunningRecipientCount > 0
	if workDone {
		return subscriptionAutomationStatusPartial
	}
	return subscriptionAutomationStatusError
}

func insertSubscriptionAutomationRun(ctx context.Context, exec masterDataExecutor, run subscriptionAutomationRunRecord, scope string, tenantID string, userID string, startedAt time.Time, finishedAt time.Time) (string, error) {
	summaryBytes, err := json.Marshal(run.Summary)
	if err != nil {
		return "", err
	}
	var runID string
	err = exec.QueryRowContext(ctx, `
INSERT INTO subscription_automation_runs (
	tenant_id,
	tenant_scope,
	trigger_source,
	status,
	dry_run,
	summary,
	error_message,
	started_at,
	finished_at,
	triggered_by_user_id
)
VALUES (
	nullif($1, '')::uuid,
	$2,
	$3,
	$4,
	$5,
	$6::jsonb,
	$7,
	$8,
	$9,
	nullif($10, '')::uuid
)
RETURNING id::text`,
		tenantID,
		scopeValueForResponse(scope),
		run.TriggerSource,
		run.Status,
		run.DryRun,
		string(summaryBytes),
		run.ErrorMessage,
		startedAt,
		finishedAt,
		userID,
	).Scan(&runID)
	return runID, err
}

func loadLatestSubscriptionAutomationRun(ctx context.Context, exec masterDataExecutor, scope string, activeTenantID string) (subscriptionAutomationRunRecord, error) {
	query := `
SELECT run.id::text,
	COALESCE(run.tenant_id::text, ''),
	COALESCE(tenant.code, ''),
	run.tenant_scope,
	run.trigger_source,
	run.status,
	run.dry_run,
	COALESCE(run.summary, '{}'::jsonb),
	COALESCE(run.error_message, ''),
	run.started_at,
	run.finished_at
FROM subscription_automation_runs run
LEFT JOIN tenants tenant ON tenant.id = run.tenant_id`
	args := []any{}
	if scope == subscriptionAutomationScopeAll {
		query += " WHERE run.tenant_scope = 'all'"
	} else {
		args = append(args, activeTenantID)
		query += fmt.Sprintf(" WHERE run.tenant_id = $%d::uuid", len(args))
	}
	query += " ORDER BY run.created_at DESC LIMIT 1"
	var (
		run          subscriptionAutomationRunRecord
		summaryBytes []byte
		startedAt    time.Time
		finishedAt   sql.NullTime
	)
	err := exec.QueryRowContext(ctx, query, args...).Scan(
		&run.ID,
		&run.TenantID,
		&run.TenantCode,
		&run.Scope,
		&run.TriggerSource,
		&run.Status,
		&run.DryRun,
		&summaryBytes,
		&run.ErrorMessage,
		&startedAt,
		&finishedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return subscriptionAutomationRunRecord{}, nil
	}
	if err != nil {
		return subscriptionAutomationRunRecord{}, err
	}
	run.StartedAt = startedAt.UTC().Format(time.RFC3339)
	if finishedAt.Valid {
		run.FinishedAt = finishedAt.Time.UTC().Format(time.RFC3339)
	}
	if len(summaryBytes) > 0 {
		_ = json.Unmarshal(summaryBytes, &run.Summary)
	}
	return run, nil
}

func scopeValueForResponse(scope string) string {
	if strings.TrimSpace(scope) == subscriptionAutomationScopeAll {
		return subscriptionAutomationScopeAll
	}
	return subscriptionAutomationScopeActive
}

func activeTenantIDForScope(scope string) string {
	if scopeValueForResponse(scope) == subscriptionAutomationScopeAll {
		return ""
	}
	return strings.TrimSpace(scope)
}
