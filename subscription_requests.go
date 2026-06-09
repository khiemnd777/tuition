package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

const (
	subscriptionChangeRequestUpgrade   = "upgrade"
	subscriptionChangeRequestDowngrade = "downgrade"
	subscriptionChangeRequestCancel    = "cancel"
	subscriptionChangeRequestRefund    = "refund"

	subscriptionChangeRequestStatusNew       = "new"
	subscriptionChangeRequestStatusApproved  = "approved"
	subscriptionChangeRequestStatusRejected  = "rejected"
	subscriptionChangeRequestStatusProcessed = "processed"
)

type subscriptionChangeRequestSummary struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenantId"`
	TenantCode        string `json:"tenantCode"`
	TenantName        string `json:"tenantName"`
	RequestType       string `json:"requestType"`
	DesiredPlanCode   string `json:"desiredPlanCode,omitempty"`
	DesiredPlanName   string `json:"desiredPlanName,omitempty"`
	Reason            string `json:"reason,omitempty"`
	EffectiveAt       string `json:"effectiveAt,omitempty"`
	RefundAmount      int    `json:"refundAmount,omitempty"`
	Status            string `json:"status"`
	RequestedByUserID string `json:"requestedByUserId,omitempty"`
	RequestedByName   string `json:"requestedByName,omitempty"`
	RequestedByEmail  string `json:"requestedByEmail,omitempty"`
	ProcessedByUserID string `json:"processedByUserId,omitempty"`
	ProcessedByName   string `json:"processedByName,omitempty"`
	ProcessedByEmail  string `json:"processedByEmail,omitempty"`
	ProcessedAt       string `json:"processedAt,omitempty"`
	AdminNote         string `json:"adminNote,omitempty"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type subscriptionChangeRequestInput struct {
	TenantID        string `json:"tenantId"`
	RequestType     string `json:"requestType"`
	DesiredPlanCode string `json:"desiredPlanCode"`
	Reason          string `json:"reason"`
	EffectiveAt     string `json:"effectiveAt"`
	RefundAmount    int    `json:"refundAmount"`
}

type subscriptionChangeRequestProcessInput struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Status    string `json:"status"`
	AdminNote string `json:"adminNote"`
	Apply     bool   `json:"apply"`
}

type subscriptionChangeRequestFilters struct {
	TenantID string
	Status   string
	Limit    int
}

type subscriptionRequestQueryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func handleSubscriptionChangeRequests(w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodGet:
		requests, err := listSubscriptionChangeRequests(r.Context(), db, subscriptionChangeRequestFilters{
			TenantID: tenantID,
			Status:   headerKey(r.URL.Query().Get("status")),
			Limit:    parsePositiveInt(r.URL.Query().Get("limit"), 50),
		})
		if err != nil {
			http.Error(w, "cannot load subscription requests", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"requests": requests})
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var input subscriptionChangeRequestInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		input.TenantID = tenantID
		request, err := createSubscriptionChangeRequest(r.Context(), db, input, user, auditContextFromRequest(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"request": request})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePlatformTenantSubscriptionRequests(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
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
		requests, err := listSubscriptionChangeRequests(r.Context(), db, subscriptionChangeRequestFilters{
			TenantID: strings.TrimSpace(r.URL.Query().Get("tenantId")),
			Status:   headerKey(r.URL.Query().Get("status")),
			Limit:    parsePositiveInt(r.URL.Query().Get("limit"), 100),
		})
		if err != nil {
			http.Error(w, "cannot load subscription requests", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"requests": requests})
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var input subscriptionChangeRequestProcessInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		request, err := processSubscriptionChangeRequest(r.Context(), db, input, user, auditContextFromRequest(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"request": request})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createSubscriptionChangeRequest(ctx context.Context, db *sql.DB, input subscriptionChangeRequestInput, user authenticatedUser, auditCtx requestAuditContext) (subscriptionChangeRequestSummary, error) {
	input = normalizeSubscriptionChangeRequestInput(input)
	if err := validateSubscriptionChangeRequestInput(ctx, db, input); err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	effectiveAt, err := parseOptionalSubscriptionRequestDate(input.EffectiveAt, "effectiveAt")
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	var requestID string
	err = db.QueryRowContext(ctx, `
INSERT INTO tenant_subscription_change_requests (
	tenant_id, request_type, desired_plan_code, reason, effective_at, refund_amount, requested_by_user_id
)
VALUES ($1::uuid, $2, $3, $4, $5, $6, nullif($7, '')::uuid)
RETURNING id::text`,
		input.TenantID,
		input.RequestType,
		input.DesiredPlanCode,
		input.Reason,
		effectiveAt,
		input.RefundAmount,
		user.ID,
	).Scan(&requestID)
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, db, auditLogInput{
		Context:    auditCtx,
		Action:     "tenant.subscription_change.request",
		EntityType: "tenant_subscription_change_request",
		EntityID:   requestID,
		Metadata: map[string]any{
			"requestType":     input.RequestType,
			"desiredPlanCode": input.DesiredPlanCode,
			"effectiveAt":     input.EffectiveAt,
			"refundAmount":    input.RefundAmount,
		},
	})
	return loadSubscriptionChangeRequestByID(ctx, db, requestID, input.TenantID, false)
}

func processSubscriptionChangeRequest(ctx context.Context, db *sql.DB, input subscriptionChangeRequestProcessInput, user authenticatedUser, auditCtx requestAuditContext) (subscriptionChangeRequestSummary, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Status = headerKey(input.Status)
	input.AdminNote = strings.TrimSpace(input.AdminNote)
	if input.ID == "" {
		return subscriptionChangeRequestSummary{}, fmt.Errorf("id is required")
	}
	switch input.Status {
	case subscriptionChangeRequestStatusApproved, subscriptionChangeRequestStatusRejected, subscriptionChangeRequestStatusProcessed:
	default:
		return subscriptionChangeRequestSummary{}, fmt.Errorf("status must be approved, rejected, or processed")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	defer tx.Rollback()

	request, err := loadSubscriptionChangeRequestByID(ctx, tx, input.ID, input.TenantID, true)
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	if request.ID == "" {
		return subscriptionChangeRequestSummary{}, fmt.Errorf("subscription request not found")
	}
	if request.Status == subscriptionChangeRequestStatusProcessed && input.Status != subscriptionChangeRequestStatusProcessed {
		return subscriptionChangeRequestSummary{}, fmt.Errorf("processed request cannot be moved back")
	}
	if request.Status == subscriptionChangeRequestStatusRejected && input.Status != subscriptionChangeRequestStatusRejected {
		return subscriptionChangeRequestSummary{}, fmt.Errorf("rejected request cannot be changed")
	}

	finalStatus := input.Status
	shouldApply := input.Apply || input.Status == subscriptionChangeRequestStatusProcessed
	if shouldApply {
		if input.Status == subscriptionChangeRequestStatusRejected {
			return subscriptionChangeRequestSummary{}, fmt.Errorf("rejected request cannot be applied")
		}
		if err := applySubscriptionChangeRequest(ctx, tx, request, user.ID); err != nil {
			return subscriptionChangeRequestSummary{}, err
		}
		finalStatus = subscriptionChangeRequestStatusProcessed
	}

	_, err = tx.ExecContext(ctx, `
UPDATE tenant_subscription_change_requests
SET status = $2,
	processed_by_user_id = nullif($3, '')::uuid,
	processed_at = now(),
	admin_note = $4,
	updated_at = now()
WHERE id = $1::uuid`,
		input.ID,
		finalStatus,
		user.ID,
		input.AdminNote,
	)
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}

	auditCtx.TenantID = request.TenantID
	_ = insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "platform.tenant.subscription_change.process",
		EntityType: "tenant_subscription_change_request",
		EntityID:   request.ID,
		Metadata: map[string]any{
			"requestType": request.RequestType,
			"status":      finalStatus,
			"applied":     shouldApply,
		},
	})
	if err := tx.Commit(); err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	return loadSubscriptionChangeRequestByID(ctx, db, input.ID, request.TenantID, false)
}

func applySubscriptionChangeRequest(ctx context.Context, tx *sql.Tx, request subscriptionChangeRequestSummary, userID string) error {
	switch request.RequestType {
	case subscriptionChangeRequestUpgrade, subscriptionChangeRequestDowngrade:
		if request.DesiredPlanCode == "" {
			return fmt.Errorf("desiredPlanCode is required for %s request", request.RequestType)
		}
		effectiveAt, err := parseOptionalSubscriptionRequestDate(request.EffectiveAt, "effectiveAt")
		if err != nil {
			return err
		}
		result, err := tx.ExecContext(ctx, `
UPDATE tenant_subscriptions ts
SET plan_id = plan.id,
	status = 'active',
	current_period_starts_at = COALESCE($3::date, ts.current_period_starts_at),
	updated_by_user_id = nullif($4, '')::uuid,
	updated_at = now()
FROM subscription_plans plan
WHERE ts.tenant_id = $1::uuid
	AND plan.code = $2
	AND plan.status = 'active'`,
			request.TenantID,
			request.DesiredPlanCode,
			effectiveAt,
			userID,
		)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("subscription plan %q is not available", request.DesiredPlanCode)
		}
	case subscriptionChangeRequestCancel:
		effectiveAt, err := parseOptionalSubscriptionRequestDate(request.EffectiveAt, "effectiveAt")
		if err != nil {
			return err
		}
		if effectiveAt == nil {
			effectiveAt = time.Now().UTC()
		}
		result, err := tx.ExecContext(ctx, `
UPDATE tenant_subscriptions
SET status = 'cancelled',
	current_period_ends_at = $2,
	updated_by_user_id = nullif($3, '')::uuid,
	updated_at = now()
WHERE tenant_id = $1::uuid`,
			request.TenantID,
			effectiveAt,
			userID,
		)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("tenant subscription not found")
		}
	case subscriptionChangeRequestRefund:
		return nil
	default:
		return fmt.Errorf("unsupported subscription request type %q", request.RequestType)
	}
	return nil
}

func normalizeSubscriptionChangeRequestInput(input subscriptionChangeRequestInput) subscriptionChangeRequestInput {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.RequestType = headerKey(input.RequestType)
	input.DesiredPlanCode = headerKey(input.DesiredPlanCode)
	input.Reason = strings.TrimSpace(input.Reason)
	input.EffectiveAt = strings.TrimSpace(input.EffectiveAt)
	if input.RefundAmount < 0 {
		input.RefundAmount = 0
	}
	return input
}

func validateSubscriptionChangeRequestInput(ctx context.Context, db *sql.DB, input subscriptionChangeRequestInput) error {
	if input.TenantID == "" {
		return fmt.Errorf("tenantId is required")
	}
	switch input.RequestType {
	case subscriptionChangeRequestUpgrade, subscriptionChangeRequestDowngrade, subscriptionChangeRequestCancel, subscriptionChangeRequestRefund:
	default:
		return fmt.Errorf("requestType must be upgrade, downgrade, cancel, or refund")
	}
	if input.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	if _, err := parseOptionalSubscriptionRequestDate(input.EffectiveAt, "effectiveAt"); err != nil {
		return err
	}
	if input.RequestType == subscriptionChangeRequestUpgrade || input.RequestType == subscriptionChangeRequestDowngrade {
		if input.DesiredPlanCode == "" {
			return fmt.Errorf("desiredPlanCode is required for %s request", input.RequestType)
		}
		var exists bool
		err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM subscription_plans WHERE code = $1 AND status = 'active')`, input.DesiredPlanCode).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("subscription plan %q is not available", input.DesiredPlanCode)
		}
	}
	return nil
}

func parseOptionalSubscriptionRequestDate(value string, field string) (any, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("%s must use YYYY-MM-DD", field)
	}
	return parsed.UTC(), nil
}

func listSubscriptionChangeRequests(ctx context.Context, queryer subscriptionRequestQueryer, filters subscriptionChangeRequestFilters) ([]subscriptionChangeRequestSummary, error) {
	limit := filters.Limit
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := queryer.QueryContext(ctx, `
SELECT req.id::text,
	req.tenant_id::text,
	tenant.code,
	tenant.name,
	req.request_type,
	req.desired_plan_code,
	COALESCE(plan.name, ''),
	req.reason,
	req.effective_at,
	req.refund_amount,
	req.status,
	COALESCE(requested.id::text, ''),
	COALESCE(NULLIF(requested.display_name, ''), requested.email, ''),
	COALESCE(requested.email, ''),
	COALESCE(processed.id::text, ''),
	COALESCE(NULLIF(processed.display_name, ''), processed.email, ''),
	COALESCE(processed.email, ''),
	req.processed_at,
	req.admin_note,
	req.created_at,
	req.updated_at
FROM tenant_subscription_change_requests req
JOIN tenants tenant ON tenant.id = req.tenant_id
LEFT JOIN subscription_plans plan ON plan.code = req.desired_plan_code
LEFT JOIN app_users requested ON requested.id = req.requested_by_user_id
LEFT JOIN app_users processed ON processed.id = req.processed_by_user_id
WHERE (nullif($1, '') IS NULL OR req.tenant_id = nullif($1, '')::uuid)
	AND (nullif($2, '') IS NULL OR req.status = $2)
ORDER BY req.created_at DESC
LIMIT $3`,
		strings.TrimSpace(filters.TenantID),
		headerKey(filters.Status),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSubscriptionChangeRequestRows(rows)
}

func loadSubscriptionChangeRequestByID(ctx context.Context, queryer subscriptionRequestQueryer, id string, tenantID string, forUpdate bool) (subscriptionChangeRequestSummary, error) {
	query := `
SELECT req.id::text,
	req.tenant_id::text,
	tenant.code,
	tenant.name,
	req.request_type,
	req.desired_plan_code,
	COALESCE(plan.name, ''),
	req.reason,
	req.effective_at,
	req.refund_amount,
	req.status,
	COALESCE(requested.id::text, ''),
	COALESCE(NULLIF(requested.display_name, ''), requested.email, ''),
	COALESCE(requested.email, ''),
	COALESCE(processed.id::text, ''),
	COALESCE(NULLIF(processed.display_name, ''), processed.email, ''),
	COALESCE(processed.email, ''),
	req.processed_at,
	req.admin_note,
	req.created_at,
	req.updated_at
FROM tenant_subscription_change_requests req
JOIN tenants tenant ON tenant.id = req.tenant_id
LEFT JOIN subscription_plans plan ON plan.code = req.desired_plan_code
LEFT JOIN app_users requested ON requested.id = req.requested_by_user_id
LEFT JOIN app_users processed ON processed.id = req.processed_by_user_id
WHERE req.id = $1::uuid
	AND (nullif($2, '') IS NULL OR req.tenant_id = nullif($2, '')::uuid)`
	if forUpdate {
		query += `
FOR UPDATE OF req`
	}
	rows, err := queryer.QueryContext(ctx, query, strings.TrimSpace(id), strings.TrimSpace(tenantID))
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	defer rows.Close()
	items, err := scanSubscriptionChangeRequestRows(rows)
	if err != nil {
		return subscriptionChangeRequestSummary{}, err
	}
	if len(items) == 0 {
		return subscriptionChangeRequestSummary{}, fmt.Errorf("subscription request not found")
	}
	return items[0], nil
}

func scanSubscriptionChangeRequestRows(rows *sql.Rows) ([]subscriptionChangeRequestSummary, error) {
	items := []subscriptionChangeRequestSummary{}
	for rows.Next() {
		var item subscriptionChangeRequestSummary
		var effectiveAt sql.NullTime
		var processedAt sql.NullTime
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(
			&item.ID,
			&item.TenantID,
			&item.TenantCode,
			&item.TenantName,
			&item.RequestType,
			&item.DesiredPlanCode,
			&item.DesiredPlanName,
			&item.Reason,
			&effectiveAt,
			&item.RefundAmount,
			&item.Status,
			&item.RequestedByUserID,
			&item.RequestedByName,
			&item.RequestedByEmail,
			&item.ProcessedByUserID,
			&item.ProcessedByName,
			&item.ProcessedByEmail,
			&processedAt,
			&item.AdminNote,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		if effectiveAt.Valid {
			item.EffectiveAt = effectiveAt.Time.UTC().Format("2006-01-02")
		}
		if processedAt.Valid {
			item.ProcessedAt = processedAt.Time.UTC().Format(time.RFC3339)
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		item.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}
	return items, rows.Err()
}

func handleSubscriptionInvoiceReceipt(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	query := r.URL.Query()
	tenantID := activeTenantIDFromRequest(r)
	requestedTenantID := strings.TrimSpace(query.Get("tenantId"))
	if requestedTenantID != "" {
		if !user.IsPlatformAdmin && requestedTenantID != tenantID {
			http.Error(w, "receipt tenant is outside the active tenant", http.StatusForbidden)
			return
		}
		tenantID = requestedTenantID
	}
	if tenantID == "" {
		http.Error(w, "tenantId is required", http.StatusBadRequest)
		return
	}
	invoiceID := strings.TrimSpace(query.Get("invoiceId"))
	if invoiceID == "" {
		http.Error(w, "invoiceId is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	invoice, err := loadSubscriptionInvoiceForReceipt(r.Context(), db, invoiceID, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if invoice.Status != subscriptionInvoiceStatusPaid {
		http.Error(w, "receipt is available after subscription invoice is paid", http.StatusBadRequest)
		return
	}
	tenant, err := loadTenantSummaryByID(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load tenant", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderSubscriptionInvoiceReceiptHTML(tenant, invoice)))
}

func loadSubscriptionInvoiceForReceipt(ctx context.Context, db *sql.DB, invoiceID string, tenantID string) (subscriptionInvoiceSummary, error) {
	var invoice subscriptionInvoiceSummary
	var periodStart time.Time
	var periodEnd time.Time
	var dueAt time.Time
	var paidAt sql.NullTime
	var createdAt time.Time
	var updatedAt time.Time
	err := db.QueryRowContext(ctx, `
SELECT id::text,
	invoice_code,
	plan_code,
	plan_name,
	amount,
	currency,
	period_starts_at,
	period_ends_at,
	due_at,
	status,
	paid_at,
	created_at,
	updated_at
FROM subscription_invoices
WHERE id = $1::uuid
	AND tenant_id = $2::uuid`,
		strings.TrimSpace(invoiceID),
		strings.TrimSpace(tenantID),
	).Scan(
		&invoice.ID,
		&invoice.InvoiceCode,
		&invoice.PlanCode,
		&invoice.PlanName,
		&invoice.Amount,
		&invoice.Currency,
		&periodStart,
		&periodEnd,
		&dueAt,
		&invoice.Status,
		&paidAt,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return subscriptionInvoiceSummary{}, fmt.Errorf("subscription invoice not found")
	}
	if err != nil {
		return subscriptionInvoiceSummary{}, err
	}
	invoice.PeriodStartsAt = periodStart.UTC().Format("2006-01-02")
	invoice.PeriodEndsAt = periodEnd.UTC().Format("2006-01-02")
	invoice.DueAt = dueAt.UTC().Format("2006-01-02")
	if paidAt.Valid {
		invoice.PaidAt = paidAt.Time.UTC().Format("2006-01-02")
	}
	invoice.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	invoice.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
	return invoice, nil
}

func renderSubscriptionInvoiceReceiptHTML(tenant tenantSummary, invoice subscriptionInvoiceSummary) string {
	return `<!doctype html>
<html lang="vi">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Biên nhận phí subscription ` + html.EscapeString(invoice.InvoiceCode) + `</title>
  <style>
    body { margin: 0; font-family: Arial, Helvetica, sans-serif; color: #172033; background: #f6f7f9; }
    main { max-width: 760px; margin: 32px auto; background: #fff; border: 1px solid #d8dee8; padding: 32px; }
    h1 { margin: 0 0 8px; font-size: 26px; }
    .muted { color: #657082; }
    .topline { display: flex; justify-content: space-between; gap: 16px; border-bottom: 2px solid #172033; padding-bottom: 18px; }
    .status { display: inline-block; padding: 6px 10px; border: 1px solid #14532d; color: #14532d; font-weight: 700; border-radius: 4px; }
    table { width: 100%; border-collapse: collapse; margin-top: 24px; }
    th, td { text-align: left; border: 1px solid #d8dee8; padding: 10px 12px; vertical-align: top; }
    th { width: 34%; background: #f2f5f8; }
    .amount { font-size: 22px; font-weight: 700; }
    @media print { body { background: #fff; } main { border: 0; margin: 0; max-width: none; } }
  </style>
</head>
<body>
  <main>
    <div class="topline">
      <div>
        <h1>Biên nhận phí subscription</h1>
        <div class="muted">DEKISUGI Finance Hub</div>
      </div>
      <div><span class="status">Đã thanh toán</span></div>
    </div>
    <table>
      <tr><th>Mã biên nhận</th><td>` + html.EscapeString(invoice.InvoiceCode) + `</td></tr>
      <tr><th>Trường / tenant</th><td>` + html.EscapeString(firstNonEmpty(tenant.Name, tenant.Code)) + `</td></tr>
      <tr><th>Gói subscription</th><td>` + html.EscapeString(firstNonEmpty(invoice.PlanName, invoice.PlanCode)) + `</td></tr>
      <tr><th>Kỳ phí</th><td>` + html.EscapeString(invoice.PeriodStartsAt) + ` đến ` + html.EscapeString(invoice.PeriodEndsAt) + `</td></tr>
      <tr><th>Ngày đến hạn</th><td>` + html.EscapeString(invoice.DueAt) + `</td></tr>
      <tr><th>Ngày thanh toán</th><td>` + html.EscapeString(firstNonEmpty(invoice.PaidAt, "-")) + `</td></tr>
      <tr><th>Số tiền</th><td class="amount">` + html.EscapeString(formatVND(invoice.Amount)) + `</td></tr>
    </table>
    <p class="muted">Biên nhận này xác nhận phí subscription của tenant đã được ghi nhận trên hệ thống.</p>
  </main>
</body>
</html>`
}
