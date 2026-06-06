package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type adminFilters struct {
	TenantID     string
	SchoolID     string
	SchoolYearID string
	ClassID      string
	Grade        string
	PeriodCode   string
	Month        int
	Status       string
	Provider     string
}

type adminDashboardResponse struct {
	Options           masterDataOptions       `json:"options"`
	Filters           adminFiltersPublic      `json:"filters"`
	Summary           adminDashboardSummary   `json:"summary"`
	Readiness         adminReadinessCenter    `json:"readiness"`
	TopClasses        []adminClassReportRow   `json:"topClasses"`
	AttentionInvoices []adminInvoiceReportRow `json:"attentionInvoices"`
}

type adminReportsResponse struct {
	Options        masterDataOptions           `json:"options"`
	Providers      []paymentProvider           `json:"providers"`
	Filters        adminFiltersPublic          `json:"filters"`
	Summary        adminDashboardSummary       `json:"summary"`
	ClassRows      []adminClassReportRow       `json:"classRows"`
	InvoiceRows    []adminInvoiceReportRow     `json:"invoiceRows"`
	TransactionRow []adminTransactionRow       `json:"transactionRows"`
	Transactions   []paymentTransactionSummary `json:"transactions"`
}

type adminFiltersPublic struct {
	SchoolID     string `json:"schoolId,omitempty"`
	SchoolYearID string `json:"schoolYearId,omitempty"`
	ClassID      string `json:"classId,omitempty"`
	Grade        string `json:"grade,omitempty"`
	PeriodCode   string `json:"periodCode,omitempty"`
	Month        int    `json:"month,omitempty"`
	Status       string `json:"status,omitempty"`
	Provider     string `json:"provider,omitempty"`
}

type adminDashboardSummary struct {
	InvoiceCount              int     `json:"invoiceCount"`
	StudentCount              int     `json:"studentCount"`
	TotalReceivable           int     `json:"totalReceivable"`
	TotalCollected            int     `json:"totalCollected"`
	OutstandingAmount         int     `json:"outstandingAmount"`
	CollectionRate            float64 `json:"collectionRate"`
	UnpaidStudentCount        int     `json:"unpaidStudentCount"`
	PartialPaymentCount       int     `json:"partialPaymentCount"`
	PaidInvoiceCount          int     `json:"paidInvoiceCount"`
	OverpaidManualReviewCount int     `json:"overpaidManualReviewCount"`
	UnmatchedTransactionCount int     `json:"unmatchedTransactionCount"`
	ManualReviewCount         int     `json:"manualReviewCount"`
}

type adminClassReportRow struct {
	SchoolYearCode string  `json:"schoolYearCode"`
	Grade          string  `json:"grade"`
	ClassName      string  `json:"className"`
	InvoiceCount   int     `json:"invoiceCount"`
	StudentCount   int     `json:"studentCount"`
	TotalAmount    int     `json:"totalAmount"`
	PaidAmount     int     `json:"paidAmount"`
	Outstanding    int     `json:"outstandingAmount"`
	CollectionRate float64 `json:"collectionRate"`
	UnpaidCount    int     `json:"unpaidCount"`
	PartialCount   int     `json:"partialCount"`
	PaidCount      int     `json:"paidCount"`
	OverpaidCount  int     `json:"overpaidCount"`
	ReviewCount    int     `json:"manualReviewCount"`
}

type adminInvoiceReportRow struct {
	invoiceSummary
	ClassID           string `json:"classId,omitempty"`
	OutstandingAmount int    `json:"outstandingAmount"`
}

type adminTransactionRow struct {
	Provider       string `json:"provider"`
	Status         string `json:"status"`
	Count          int    `json:"count"`
	TotalAmount    int    `json:"totalAmount"`
	MatchedCount   int    `json:"matchedCount"`
	UnmatchedCount int    `json:"unmatchedCount"`
	ReviewCount    int    `json:"manualReviewCount"`
}

type adminUsersResponse struct {
	Users       []adminUserSummary       `json:"users"`
	Roles       []adminRoleSummary       `json:"roles"`
	Permissions []adminPermissionSummary `json:"permissions"`
}

type adminUserSummary struct {
	ID          string             `json:"id"`
	Email       string             `json:"email"`
	Phone       string             `json:"phone"`
	DisplayName string             `json:"displayName"`
	Status      string             `json:"status"`
	HasPassword bool               `json:"hasPassword"`
	LastLoginAt string             `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Roles       []adminRoleSummary `json:"roles"`
}

type adminRoleSummary struct {
	ID          string                   `json:"id"`
	Code        string                   `json:"code"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	IsSystem    bool                     `json:"isSystem"`
	Permissions []adminPermissionSummary `json:"permissions,omitempty"`
}

type adminPermissionSummary struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type adminUserSaveInput struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	Password    string `json:"password,omitempty"`
}

type adminUserRoleInput struct {
	UserID    string   `json:"userId"`
	RoleCodes []string `json:"roleCodes"`
}

func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
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

	options, err := listMasterDataOptions(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load admin options", http.StatusInternalServerError)
		return
	}
	filters := parseAdminFilters(r)
	filters.TenantID = tenantID
	invoices, err := listAdminInvoiceRows(r.Context(), db, filters, 5000)
	if err != nil {
		http.Error(w, "cannot load dashboard invoices", http.StatusInternalServerError)
		return
	}
	transactions, err := listPaymentTransactions(r.Context(), db, paymentTransactionListFilters{TenantID: tenantID, Limit: 1000})
	if err != nil {
		http.Error(w, "cannot load dashboard transactions", http.StatusInternalServerError)
		return
	}
	readiness, err := buildAdminReadinessCenter(r.Context(), db, filters, invoices, transactions, time.Now())
	if err != nil {
		http.Error(w, "cannot load dashboard readiness", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, adminDashboardResponse{
		Options:           options,
		Filters:           filters.public(),
		Summary:           buildAdminDashboardSummary(invoices, transactions),
		Readiness:         readiness,
		TopClasses:        topAdminClasses(buildAdminClassReportRows(invoices), 8),
		AttentionInvoices: adminAttentionInvoices(invoices, 12),
	})
}

func handleAdminReports(w http.ResponseWriter, r *http.Request) {
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

	options, err := listMasterDataOptions(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load admin options", http.StatusInternalServerError)
		return
	}
	providers, err := listPaymentProviders(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load report providers", http.StatusInternalServerError)
		return
	}
	filters := parseAdminFilters(r)
	filters.TenantID = tenantID
	invoices, err := listAdminInvoiceRows(r.Context(), db, filters, 1000)
	if err != nil {
		http.Error(w, "cannot load report invoices", http.StatusInternalServerError)
		return
	}
	transactions, err := listPaymentTransactions(r.Context(), db, paymentTransactionListFilters{TenantID: tenantID, Provider: filters.Provider, Limit: 1000})
	if err != nil {
		http.Error(w, "cannot load report transactions", http.StatusInternalServerError)
		return
	}
	transactions = filterAdminReportTransactions(invoices, transactions, filters)

	writeJSON(w, http.StatusOK, adminReportsResponse{
		Options:        options,
		Providers:      providers,
		Filters:        filters.public(),
		Summary:        buildAdminDashboardSummary(invoices, transactions),
		ClassRows:      buildAdminClassReportRows(invoices),
		InvoiceRows:    invoices,
		TransactionRow: buildAdminTransactionRows(transactions),
		Transactions:   transactions,
	})
}

func handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	tenantID := activeTenantIDFromRequest(r)
	if tenantID == "" {
		http.Error(w, "active tenant required", http.StatusForbidden)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	users, roles, permissions, err := loadAdminUsersAndRoles(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load users and roles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, adminUsersResponse{
		Users:       users,
		Roles:       roles,
		Permissions: permissions,
	})
}

func handlePlatformUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok || !user.IsPlatformAdmin {
		http.Error(w, "platform admin required", http.StatusForbidden)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	users, roles, permissions, err := loadPlatformUsersAndRoles(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load platform users and roles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, adminUsersResponse{
		Users:       users,
		Roles:       roles,
		Permissions: permissions,
	})
}

func handleAdminRoles(w http.ResponseWriter, r *http.Request) {
	tenantID := activeTenantIDFromRequest(r)
	if tenantID == "" {
		http.Error(w, "active tenant required", http.StatusForbidden)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	_, roles, permissions, err := loadAdminUsersAndRoles(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load roles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"roles": roles, "permissions": permissions})
}

func handlePlatformRoles(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok || !user.IsPlatformAdmin {
		http.Error(w, "platform admin required", http.StatusForbidden)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	_, roles, permissions, err := loadPlatformUsersAndRoles(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load platform roles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"roles": roles, "permissions": permissions})
}

func handleAdminUserSave(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input adminUserSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if err := validateAdminUserSaveInput(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	permission := "user.create"
	if input.ID != "" {
		permission = "user.update"
	}
	if !requireAdminAPIPermission(w, r, permission) {
		return
	}
	tenantID := activeTenantIDFromRequest(r)
	if tenantID == "" {
		http.Error(w, "active tenant required", http.StatusForbidden)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	user, err := saveAdminUser(r.Context(), db, input, tenantID)
	if err != nil {
		http.Error(w, "cannot save user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

func handlePlatformUserSave(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok || !user.IsPlatformAdmin {
		http.Error(w, "platform admin required", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input adminUserSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if err := validateAdminUserSaveInput(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	permission := "user.create"
	if input.ID != "" {
		permission = "user.update"
	}
	if !authenticatedUserHasPermission(user, permission) {
		http.Error(w, "missing required API permission: "+permission, http.StatusForbidden)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	saved, err := saveAdminUser(r.Context(), db, input, "")
	if err != nil {
		http.Error(w, "cannot save platform user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": saved})
}

func handleAdminUserRoles(w http.ResponseWriter, r *http.Request) {
	if !requireAdminAPIPermission(w, r, "user.assign_role") {
		return
	}
	tenantID := activeTenantIDFromRequest(r)
	if tenantID == "" {
		http.Error(w, "active tenant required", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input adminUserRoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.UserID = strings.TrimSpace(input.UserID)
	input.RoleCodes = normalizeAdminRoleCodes(input.RoleCodes)
	if input.UserID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if err := assignAdminUserRoles(r.Context(), db, input, tenantID); err != nil {
		var usageErr *tenantUsageLimitError
		if errors.As(err, &usageErr) {
			http.Error(w, usageErr.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "cannot assign user roles", http.StatusInternalServerError)
		return
	}
	users, _, _, err := loadAdminUsersAndRoles(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot reload users", http.StatusInternalServerError)
		return
	}
	for _, user := range users {
		if user.ID == input.UserID {
			writeJSON(w, http.StatusOK, map[string]any{"user": user})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"userId": input.UserID})
}

func handlePlatformUserRoles(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok || !user.IsPlatformAdmin {
		http.Error(w, "platform admin required", http.StatusForbidden)
		return
	}
	if !authenticatedUserHasPermission(user, "user.assign_role") {
		http.Error(w, "missing required API permission: user.assign_role", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input adminUserRoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.UserID = strings.TrimSpace(input.UserID)
	input.RoleCodes = normalizePlatformRoleCodes(input.RoleCodes)
	if input.UserID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if err := assignPlatformUserRoles(r.Context(), db, input); err != nil {
		http.Error(w, "cannot assign platform user roles", http.StatusInternalServerError)
		return
	}
	users, _, _, err := loadPlatformUsersAndRoles(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot reload platform users", http.StatusInternalServerError)
		return
	}
	for _, platformUser := range users {
		if platformUser.ID == input.UserID {
			writeJSON(w, http.StatusOK, map[string]any{"user": platformUser})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"userId": input.UserID})
}

func requireAdminAPIPermission(w http.ResponseWriter, r *http.Request, permission string) bool {
	if user, ok := authenticatedUserFromRequest(r); ok {
		if !authenticatedUserHasActiveTenant(user) {
			http.Error(w, "active tenant required", http.StatusForbidden)
			return false
		}
		if authenticatedUserHasPermission(user, permission) {
			return true
		}
	}
	http.Error(w, "missing required API permission: "+permission, http.StatusForbidden)
	return false
}

func parseAdminFilters(r *http.Request) adminFilters {
	query := r.URL.Query()
	return adminFilters{
		SchoolID:     strings.TrimSpace(query.Get("schoolId")),
		SchoolYearID: strings.TrimSpace(query.Get("schoolYearId")),
		ClassID:      strings.TrimSpace(query.Get("classId")),
		Grade:        normalizeGrade(query.Get("grade")),
		PeriodCode:   strings.TrimSpace(query.Get("periodCode")),
		Month:        parsePositiveInt(query.Get("month"), 0),
		Status:       headerKey(query.Get("status")),
		Provider:     headerKey(query.Get("provider")),
	}
}

func (filters adminFilters) public() adminFiltersPublic {
	return adminFiltersPublic{
		SchoolID:     filters.SchoolID,
		SchoolYearID: filters.SchoolYearID,
		ClassID:      filters.ClassID,
		Grade:        filters.Grade,
		PeriodCode:   filters.PeriodCode,
		Month:        filters.Month,
		Status:       filters.Status,
		Provider:     filters.Provider,
	}
}

func listAdminInvoiceRows(ctx context.Context, db *sql.DB, filters adminFilters, limit int) ([]adminInvoiceReportRow, error) {
	conditions := []string{"i.status <> 'void'"}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.TenantID != "" {
		conditions = append(conditions, `EXISTS (
			SELECT 1
			FROM school_years tenant_sy
			JOIN schools tenant_school ON tenant_school.id = tenant_sy.school_id
			WHERE tenant_sy.id = i.school_year_id
				AND tenant_school.tenant_id = `+addArg(filters.TenantID)+`::uuid
		)`)
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "i.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.SchoolID != "" {
		conditions = append(conditions, `EXISTS (
			SELECT 1
			FROM school_years sy_filter
			WHERE sy_filter.id = i.school_year_id
				AND sy_filter.school_id = `+addArg(filters.SchoolID)+`::uuid
		)`)
	}
	if filters.ClassID != "" {
		conditions = append(conditions, "i.class_id = "+addArg(filters.ClassID)+"::uuid")
	}
	if filters.Grade != "" {
		conditions = append(conditions, "i.grade = "+addArg(filters.Grade))
	}
	if filters.PeriodCode != "" {
		conditions = append(conditions, "i.period_code = "+addArg(filters.PeriodCode))
	}
	if filters.Month > 0 {
		conditions = append(conditions, "i.month = "+addArg(filters.Month))
	}
	if filters.Status != "" {
		conditions = append(conditions, "i.status = "+addArg(filters.Status))
	}
	if limit <= 0 || limit > 5000 {
		limit = 1000
	}
	limitArg := addArg(limit)

	query := `
SELECT i.id::text,
	i.invoice_code,
	i.student_code,
	i.student_name,
	i.class_id::text,
	i.class_name,
	i.grade,
	i.school_year_id::text,
	i.school_year_code,
	i.fee_schedule_id::text,
	i.period_code,
	i.month,
	i.issued_at,
	i.due_date,
	i.status,
	i.base_amount,
	i.adjustment_amount,
	i.total_amount,
	i.paid_amount,
	GREATEST(i.total_amount - i.paid_amount, 0),
	i.collection_bank_bin,
	i.collection_bank_account,
	i.qr_bill_number,
	i.qr_note,
	COALESCE(item_counts.item_count, 0),
	COALESCE(adjustment_counts.adjustment_count, 0),
	COALESCE(intent_counts.intent_count, 0),
	COALESCE(match_counts.match_count, 0),
	COALESCE(notification_counts.sent_count, 0),
	notification_counts.last_sent_at,
	(i.qr_bill_number <> '' AND i.collection_bank_bin <> '' AND i.collection_bank_account <> ''),
	true
FROM invoices i
LEFT JOIN (
	SELECT invoice_id, COUNT(*)::integer AS item_count
	FROM invoice_items
	GROUP BY invoice_id
) item_counts ON item_counts.invoice_id = i.id
LEFT JOIN (
	SELECT invoice_id, COUNT(*)::integer AS adjustment_count
	FROM invoice_adjustments
	GROUP BY invoice_id
) adjustment_counts ON adjustment_counts.invoice_id = i.id
LEFT JOIN (
	SELECT invoice_id, COUNT(*)::integer AS intent_count
	FROM payment_intents
	WHERE status NOT IN ('cancelled', 'expired', 'failed')
	GROUP BY invoice_id
) intent_counts ON intent_counts.invoice_id = i.id
LEFT JOIN (
	SELECT invoice_id, COUNT(*)::integer AS match_count
	FROM reconciliation_matches
	WHERE status = 'matched'
	GROUP BY invoice_id
) match_counts ON match_counts.invoice_id = i.id
LEFT JOIN (
	SELECT invoice_id, COUNT(*)::integer AS sent_count, MAX(sent_at) AS last_sent_at
	FROM notification_logs
	WHERE status = 'sent'
		AND NOT dry_run
	GROUP BY invoice_id
) notification_counts ON notification_counts.invoice_id = i.id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY i.period_code DESC, i.class_name, i.student_code
LIMIT ` + limitArg

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []adminInvoiceReportRow{}
	for rows.Next() {
		var invoice adminInvoiceReportRow
		var month sql.NullInt64
		var dueDate sql.NullTime
		var lastSentAt sql.NullTime
		if err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceCode,
			&invoice.StudentCode,
			&invoice.StudentName,
			&invoice.ClassID,
			&invoice.ClassName,
			&invoice.Grade,
			&invoice.SchoolYearID,
			&invoice.SchoolYearCode,
			&invoice.FeeScheduleID,
			&invoice.PeriodCode,
			&month,
			&invoice.IssuedAt,
			&dueDate,
			&invoice.Status,
			&invoice.BaseAmount,
			&invoice.AdjustmentAmount,
			&invoice.TotalAmount,
			&invoice.PaidAmount,
			&invoice.OutstandingAmount,
			&invoice.CollectionBankBIN,
			&invoice.CollectionBankAccount,
			&invoice.QRBillNumber,
			&invoice.QRNote,
			&invoice.ItemCount,
			&invoice.AdjustmentCount,
			&invoice.PaymentIntentCount,
			&invoice.MatchedPaymentCount,
			&invoice.SentCount,
			&lastSentAt,
			&invoice.QRReady,
			&invoice.PDFReady,
		); err != nil {
			return nil, err
		}
		if month.Valid {
			invoice.Month = int(month.Int64)
		}
		if dueDate.Valid {
			invoice.DueDate = dueDate.Time.Format("2006-01-02")
		}
		if lastSentAt.Valid {
			invoice.LastSentAt = lastSentAt.Time.UTC().Format(time.RFC3339)
		}
		invoice.IssueState = invoiceSummaryIssueState(invoice.invoiceSummary)
		invoices = append(invoices, invoice)
	}
	return invoices, rows.Err()
}

func buildAdminDashboardSummary(invoices []adminInvoiceReportRow, transactions []paymentTransactionSummary) adminDashboardSummary {
	summary := adminDashboardSummary{InvoiceCount: len(invoices)}
	studentCodes := map[string]bool{}
	unpaidStudents := map[string]bool{}
	for _, invoice := range invoices {
		if invoice.StudentCode != "" {
			studentCodes[invoice.StudentCode] = true
		}
		summary.TotalReceivable += invoice.TotalAmount
		summary.TotalCollected += invoice.PaidAmount
		summary.OutstandingAmount += outstandingAmount(invoice.TotalAmount, invoice.PaidAmount)
		switch invoice.Status {
		case invoiceStatusUnpaid:
			if invoice.StudentCode != "" {
				unpaidStudents[invoice.StudentCode] = true
			}
		case invoiceStatusPartial:
			summary.PartialPaymentCount++
		case invoiceStatusPaid:
			summary.PaidInvoiceCount++
		case invoiceStatusOverpaid, invoiceStatusManualReview:
			summary.OverpaidManualReviewCount++
		}
	}
	summary.StudentCount = len(studentCodes)
	summary.UnpaidStudentCount = len(unpaidStudents)
	if summary.TotalReceivable > 0 {
		summary.CollectionRate = float64(summary.TotalCollected) / float64(summary.TotalReceivable)
	}
	for _, transaction := range transactions {
		switch transaction.Status {
		case paymentTransactionStatusUnmatched:
			summary.UnmatchedTransactionCount++
		case paymentTransactionStatusManualReview:
			summary.ManualReviewCount++
		}
	}
	return summary
}

func buildAdminClassReportRows(invoices []adminInvoiceReportRow) []adminClassReportRow {
	type classAccumulator struct {
		row      adminClassReportRow
		students map[string]bool
	}
	classes := map[string]*classAccumulator{}
	for _, invoice := range invoices {
		key := strings.Join([]string{invoice.SchoolYearCode, invoice.Grade, invoice.ClassName}, "\x00")
		acc, ok := classes[key]
		if !ok {
			acc = &classAccumulator{
				row: adminClassReportRow{
					SchoolYearCode: invoice.SchoolYearCode,
					Grade:          invoice.Grade,
					ClassName:      invoice.ClassName,
				},
				students: map[string]bool{},
			}
			classes[key] = acc
		}
		acc.row.InvoiceCount++
		acc.row.TotalAmount += invoice.TotalAmount
		acc.row.PaidAmount += invoice.PaidAmount
		acc.row.Outstanding += outstandingAmount(invoice.TotalAmount, invoice.PaidAmount)
		if invoice.StudentCode != "" {
			acc.students[invoice.StudentCode] = true
		}
		switch invoice.Status {
		case invoiceStatusUnpaid:
			acc.row.UnpaidCount++
		case invoiceStatusPartial:
			acc.row.PartialCount++
		case invoiceStatusPaid:
			acc.row.PaidCount++
		case invoiceStatusOverpaid:
			acc.row.OverpaidCount++
		case invoiceStatusManualReview:
			acc.row.ReviewCount++
		}
	}
	rows := make([]adminClassReportRow, 0, len(classes))
	for _, acc := range classes {
		acc.row.StudentCount = len(acc.students)
		if acc.row.TotalAmount > 0 {
			acc.row.CollectionRate = float64(acc.row.PaidAmount) / float64(acc.row.TotalAmount)
		}
		rows = append(rows, acc.row)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Outstanding != rows[j].Outstanding {
			return rows[i].Outstanding > rows[j].Outstanding
		}
		if rows[i].SchoolYearCode != rows[j].SchoolYearCode {
			return rows[i].SchoolYearCode > rows[j].SchoolYearCode
		}
		if rows[i].Grade != rows[j].Grade {
			return rows[i].Grade < rows[j].Grade
		}
		return rows[i].ClassName < rows[j].ClassName
	})
	return rows
}

func topAdminClasses(rows []adminClassReportRow, limit int) []adminClassReportRow {
	if limit <= 0 || len(rows) <= limit {
		return rows
	}
	return rows[:limit]
}

func adminAttentionInvoices(invoices []adminInvoiceReportRow, limit int) []adminInvoiceReportRow {
	rows := make([]adminInvoiceReportRow, 0, len(invoices))
	for _, invoice := range invoices {
		if invoice.Status == invoiceStatusUnpaid ||
			invoice.Status == invoiceStatusPartial ||
			invoice.Status == invoiceStatusOverpaid ||
			invoice.Status == invoiceStatusManualReview {
			rows = append(rows, invoice)
		}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].OutstandingAmount != rows[j].OutstandingAmount {
			return rows[i].OutstandingAmount > rows[j].OutstandingAmount
		}
		return rows[i].InvoiceCode < rows[j].InvoiceCode
	})
	if limit > 0 && len(rows) > limit {
		return rows[:limit]
	}
	return rows
}

func buildAdminTransactionRows(transactions []paymentTransactionSummary) []adminTransactionRow {
	rowsByKey := map[string]*adminTransactionRow{}
	for _, transaction := range transactions {
		key := transaction.ProviderCode + "\x00" + transaction.Status
		row, ok := rowsByKey[key]
		if !ok {
			row = &adminTransactionRow{Provider: transaction.ProviderCode, Status: transaction.Status}
			rowsByKey[key] = row
		}
		row.Count++
		if transaction.Direction == paymentDirectionIn {
			row.TotalAmount += transaction.Amount
		}
		switch transaction.Status {
		case paymentTransactionStatusMatched:
			row.MatchedCount++
		case paymentTransactionStatusUnmatched:
			row.UnmatchedCount++
		case paymentTransactionStatusManualReview:
			row.ReviewCount++
		}
	}
	rows := make([]adminTransactionRow, 0, len(rowsByKey))
	for _, row := range rowsByKey {
		rows = append(rows, *row)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Provider != rows[j].Provider {
			return rows[i].Provider < rows[j].Provider
		}
		return rows[i].Status < rows[j].Status
	})
	return rows
}

func outstandingAmount(total int, paid int) int {
	if total <= paid {
		return 0
	}
	return total - paid
}

func loadAdminUsersAndRoles(ctx context.Context, db *sql.DB, tenantID string) ([]adminUserSummary, []adminRoleSummary, []adminPermissionSummary, error) {
	permissions, err := listAdminPermissions(ctx, db)
	if err != nil {
		return nil, nil, nil, err
	}
	roles, err := listAdminRoles(ctx, db)
	if err != nil {
		return nil, nil, nil, err
	}
	users, err := listAdminUsers(ctx, db, tenantID)
	if err != nil {
		return nil, nil, nil, err
	}
	return users, roles, permissions, nil
}

func loadPlatformUsersAndRoles(ctx context.Context, db *sql.DB) ([]adminUserSummary, []adminRoleSummary, []adminPermissionSummary, error) {
	permissions, err := listAdminPermissions(ctx, db)
	if err != nil {
		return nil, nil, nil, err
	}
	roles, err := listPlatformRoles(ctx, db)
	if err != nil {
		return nil, nil, nil, err
	}
	users, err := listPlatformUsers(ctx, db)
	if err != nil {
		return nil, nil, nil, err
	}
	return users, roles, permissions, nil
}

func listAdminPermissions(ctx context.Context, db *sql.DB) ([]adminPermissionSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id::text, code, description
FROM app_permissions
ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []adminPermissionSummary{}
	for rows.Next() {
		var item adminPermissionSummary
		if err := rows.Scan(&item.ID, &item.Code, &item.Description); err != nil {
			return nil, err
		}
		if !isCanonicalAdminPermissionCode(item.Code) {
			continue
		}
		permissions = append(permissions, item)
	}
	return permissions, rows.Err()
}

func listAdminRoles(ctx context.Context, db *sql.DB) ([]adminRoleSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT r.id::text,
	r.code,
	r.name,
	r.description,
	r.is_system,
	COALESCE(p.id::text, ''),
	COALESCE(p.code, ''),
	COALESCE(p.description, '')
FROM app_roles r
LEFT JOIN app_role_permissions rp ON rp.role_id = r.id
LEFT JOIN app_permissions p ON p.id = rp.permission_id
WHERE r.code IN ('tenant_owner', 'tenant_admin', 'tenant_staff', 'tenant_accountant')
ORDER BY CASE r.code
	WHEN 'tenant_owner' THEN 1
	WHEN 'tenant_admin' THEN 2
	WHEN 'tenant_staff' THEN 3
	WHEN 'tenant_accountant' THEN 4
	ELSE 99
END, p.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleOrder := []string{}
	roleByID := map[string]*adminRoleSummary{}
	for rows.Next() {
		var roleID string
		var role adminRoleSummary
		var permission adminPermissionSummary
		if err := rows.Scan(
			&roleID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
			&permission.ID,
			&permission.Code,
			&permission.Description,
		); err != nil {
			return nil, err
		}
		existing, ok := roleByID[roleID]
		if !ok {
			role.ID = roleID
			role.Permissions = []adminPermissionSummary{}
			existing = &role
			roleByID[roleID] = existing
			roleOrder = append(roleOrder, roleID)
		}
		if permission.Code != "" && isCanonicalAdminPermissionCode(permission.Code) {
			existing.Permissions = append(existing.Permissions, permission)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	roles := make([]adminRoleSummary, 0, len(roleOrder))
	for _, roleID := range roleOrder {
		roles = append(roles, *roleByID[roleID])
	}
	return roles, nil
}

func listPlatformRoles(ctx context.Context, db *sql.DB) ([]adminRoleSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT r.id::text,
	r.code,
	r.name,
	r.description,
	r.is_system,
	COALESCE(p.id::text, ''),
	COALESCE(p.code, ''),
	COALESCE(p.description, '')
FROM app_roles r
LEFT JOIN app_role_permissions rp ON rp.role_id = r.id
LEFT JOIN app_permissions p ON p.id = rp.permission_id
WHERE r.code IN ('platform_admin')
ORDER BY r.code, p.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleOrder := []string{}
	roleByID := map[string]*adminRoleSummary{}
	for rows.Next() {
		var roleID string
		var role adminRoleSummary
		var permission adminPermissionSummary
		if err := rows.Scan(
			&roleID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
			&permission.ID,
			&permission.Code,
			&permission.Description,
		); err != nil {
			return nil, err
		}
		existing := roleByID[roleID]
		if existing == nil {
			role.ID = roleID
			role.Permissions = []adminPermissionSummary{}
			existing = &role
			roleByID[roleID] = existing
			roleOrder = append(roleOrder, roleID)
		}
		if permission.Code != "" && isCanonicalAdminPermissionCode(permission.Code) {
			existing.Permissions = append(existing.Permissions, permission)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	roles := make([]adminRoleSummary, 0, len(roleOrder))
	for _, roleID := range roleOrder {
		roles = append(roles, *roleByID[roleID])
	}
	return roles, nil
}

func listAdminUsers(ctx context.Context, db *sql.DB, tenantID string) ([]adminUserSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT u.id::text,
	COALESCE(u.email, ''),
	u.phone,
	u.display_name,
	u.status,
	u.password_hash <> '',
	u.last_login_at,
	u.created_at,
	u.updated_at,
	COALESCE(r.id::text, ''),
	COALESCE(r.code, ''),
	COALESCE(r.name, ''),
	COALESCE(r.description, ''),
	COALESCE(r.is_system, false)
FROM tenant_memberships membership
JOIN app_users u ON u.id = membership.user_id
LEFT JOIN tenant_user_roles ur ON ur.user_id = u.id AND ur.tenant_id = membership.tenant_id
LEFT JOIN app_roles r ON r.id = ur.role_id
	AND r.code IN ('tenant_owner', 'tenant_admin', 'tenant_staff', 'tenant_accountant')
WHERE membership.tenant_id = $1::uuid
	AND membership.status <> 'removed'
ORDER BY lower(u.email), r.code`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userOrder := []string{}
	userByID := map[string]*adminUserSummary{}
	for rows.Next() {
		var userID string
		var user adminUserSummary
		var lastLogin sql.NullTime
		var role adminRoleSummary
		if err := rows.Scan(
			&userID,
			&user.Email,
			&user.Phone,
			&user.DisplayName,
			&user.Status,
			&user.HasPassword,
			&lastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
		); err != nil {
			return nil, err
		}
		existing, ok := userByID[userID]
		if !ok {
			user.ID = userID
			user.Roles = []adminRoleSummary{}
			if lastLogin.Valid {
				user.LastLoginAt = lastLogin.Time.Format(time.RFC3339)
			}
			existing = &user
			userByID[userID] = existing
			userOrder = append(userOrder, userID)
		}
		if role.Code != "" {
			existing.Roles = append(existing.Roles, role)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	users := make([]adminUserSummary, 0, len(userOrder))
	for _, userID := range userOrder {
		users = append(users, *userByID[userID])
	}
	return users, nil
}

func listPlatformUsers(ctx context.Context, db *sql.DB) ([]adminUserSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT u.id::text,
	COALESCE(u.email, ''),
	u.phone,
	u.display_name,
	u.status,
	u.password_hash <> '',
	u.last_login_at,
	u.created_at,
	u.updated_at,
	COALESCE(r.id::text, ''),
	COALESCE(r.code, ''),
	COALESCE(r.name, ''),
	COALESCE(r.description, ''),
	COALESCE(r.is_system, false)
FROM app_users u
LEFT JOIN app_user_roles ur ON ur.user_id = u.id
LEFT JOIN app_roles r ON r.id = ur.role_id
	AND r.code IN ('platform_admin')
ORDER BY lower(COALESCE(u.email, '')), u.phone, r.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userOrder := []string{}
	userByID := map[string]*adminUserSummary{}
	for rows.Next() {
		var userID string
		var platformUser adminUserSummary
		var lastLogin sql.NullTime
		var role adminRoleSummary
		if err := rows.Scan(
			&userID,
			&platformUser.Email,
			&platformUser.Phone,
			&platformUser.DisplayName,
			&platformUser.Status,
			&platformUser.HasPassword,
			&lastLogin,
			&platformUser.CreatedAt,
			&platformUser.UpdatedAt,
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
		); err != nil {
			return nil, err
		}
		existing := userByID[userID]
		if existing == nil {
			platformUser.ID = userID
			platformUser.Roles = []adminRoleSummary{}
			if lastLogin.Valid {
				platformUser.LastLoginAt = lastLogin.Time.Format(time.RFC3339)
			}
			existing = &platformUser
			userByID[userID] = existing
			userOrder = append(userOrder, userID)
		}
		if role.Code != "" {
			existing.Roles = append(existing.Roles, role)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	users := make([]adminUserSummary, 0, len(userOrder))
	for _, userID := range userOrder {
		users = append(users, *userByID[userID])
	}
	return users, nil
}

func validateAdminUserSaveInput(input *adminUserSaveInput) error {
	input.ID = strings.TrimSpace(input.ID)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Phone = normalizeAdminPhone(input.Phone)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Status = headerKey(input.Status)
	input.Password = strings.TrimSpace(input.Password)
	if input.Status == "" {
		input.Status = "active"
	}
	if input.Email == "" && input.Phone == "" {
		return fmt.Errorf("email or phone is required")
	}
	if input.Email != "" && !strings.Contains(input.Email, "@") {
		return fmt.Errorf("valid email is required")
	}
	if input.Phone != "" {
		if err := validateAdminPhone(input.Phone); err != nil {
			return err
		}
	}
	if input.DisplayName == "" {
		input.DisplayName = firstNonEmpty(input.Email, input.Phone)
	}
	if input.Password != "" {
		if err := validateAuthPassword(input.Password); err != nil {
			return err
		}
	}
	switch input.Status {
	case "active", "inactive", "suspended":
		return nil
	default:
		return fmt.Errorf("unsupported user status %q", input.Status)
	}
}

func normalizeAdminPhone(value string) string {
	value = strings.TrimSpace(value)
	var builder strings.Builder
	for idx, char := range value {
		switch {
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
		case char == '+' && idx == 0:
			builder.WriteRune(char)
		case char == ' ' || char == '-' || char == '.' || char == '(' || char == ')':
			continue
		default:
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func validateAdminPhone(phone string) error {
	digits := phone
	if strings.HasPrefix(digits, "+") {
		digits = strings.TrimPrefix(digits, "+")
	}
	if len(digits) < 8 || len(digits) > 15 {
		return fmt.Errorf("valid phone is required")
	}
	for _, char := range digits {
		if char < '0' || char > '9' {
			return fmt.Errorf("valid phone is required")
		}
	}
	return nil
}

func saveAdminUser(ctx context.Context, db *sql.DB, input adminUserSaveInput, tenantID string) (adminUserSummary, error) {
	passwordHash := ""
	if input.Password != "" {
		var err error
		passwordHash, err = hashPassword(input.Password)
		if err != nil {
			return adminUserSummary{}, err
		}
	}
	if input.ID == "" {
		var existingID string
		err := db.QueryRowContext(ctx, `
SELECT id::text
FROM app_users
WHERE ($1 <> '' AND lower(COALESCE(email, '')) = lower($1))
	OR ($2 <> '' AND phone = $2)
LIMIT 1`, input.Email, input.Phone).Scan(&existingID)
		if err == nil && existingID != "" {
			input.ID = existingID
			return saveAdminUser(ctx, db, input, tenantID)
		}
		if err != nil && err != sql.ErrNoRows {
			return adminUserSummary{}, err
		}
		var user adminUserSummary
		err = db.QueryRowContext(ctx, `
INSERT INTO app_users (email, phone, display_name, status, password_hash, password_updated_at)
VALUES ($1, $2, $3, $4, $5, CASE WHEN $5 <> '' THEN now() ELSE NULL END)
RETURNING id::text, COALESCE(email, ''), phone, display_name, status, password_hash <> '', created_at, updated_at`,
			input.Email,
			input.Phone,
			input.DisplayName,
			input.Status,
			passwordHash,
		).Scan(&user.ID, &user.Email, &user.Phone, &user.DisplayName, &user.Status, &user.HasPassword, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return adminUserSummary{}, err
		}
		if tenantID != "" {
			if err := ensureTenantMembership(ctx, db, tenantID, user.ID, false); err != nil {
				return adminUserSummary{}, err
			}
		}
		user.Roles = []adminRoleSummary{}
		return user, nil
	}
	var user adminUserSummary
	err := db.QueryRowContext(ctx, `
UPDATE app_users
SET email = $2,
	phone = $3,
	display_name = $4,
	status = $5,
	password_hash = CASE WHEN $6 <> '' THEN $6 ELSE password_hash END,
	password_updated_at = CASE WHEN $6 <> '' THEN now() ELSE password_updated_at END
WHERE id = $1::uuid
RETURNING id::text, COALESCE(email, ''), phone, display_name, status, password_hash <> '', created_at, updated_at`,
		input.ID,
		input.Email,
		input.Phone,
		input.DisplayName,
		input.Status,
		passwordHash,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.DisplayName, &user.Status, &user.HasPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return adminUserSummary{}, err
	}
	if tenantID != "" {
		if err := ensureTenantMembership(ctx, db, tenantID, user.ID, false); err != nil {
			return adminUserSummary{}, err
		}
	}
	user.Roles = []adminRoleSummary{}
	return user, nil
}

func assignAdminUserRoles(ctx context.Context, db *sql.DB, input adminUserRoleInput, tenantID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	hadRoles, err := tenantUserHasAssignedRoles(ctx, tx, tenantID, input.UserID)
	if err != nil {
		return err
	}
	if !hadRoles && len(input.RoleCodes) > 0 {
		if err := enforceTenantUsageLimit(ctx, tx, tenantID, subscriptionMetricOperators, 1, time.Now()); err != nil {
			return err
		}
	}
	if err := ensureTenantMembership(ctx, tx, tenantID, input.UserID, false); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM tenant_user_roles WHERE tenant_id = $1::uuid AND user_id = $2::uuid`, tenantID, input.UserID); err != nil {
		return err
	}
	for _, roleCode := range input.RoleCodes {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO tenant_user_roles (tenant_id, user_id, role_id)
SELECT $1::uuid, $2::uuid, id
FROM app_roles
WHERE code = $3
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING`,
			tenantID,
			input.UserID,
			roleCode,
		); err != nil {
			return err
		}
	}
	if err := rebuildTenantUsageCounter(ctx, tx, tenantID, subscriptionMetricOperators, time.Now()); err != nil {
		return err
	}
	return tx.Commit()
}

func tenantUserHasAssignedRoles(ctx context.Context, exec masterDataExecutor, tenantID string, userID string) (bool, error) {
	var hasRoles bool
	err := exec.QueryRowContext(ctx, `
SELECT EXISTS (
	SELECT 1
	FROM tenant_user_roles
	WHERE tenant_id = $1::uuid
		AND user_id = $2::uuid
)`, tenantID, userID).Scan(&hasRoles)
	return hasRoles, err
}

func normalizeAdminRoleCodes(values []string) []string {
	seen := map[string]bool{}
	normalized := []string{}
	for _, value := range values {
		code := headerKey(value)
		if code == "" || seen[code] || !isCanonicalAdminRoleCode(code) {
			continue
		}
		seen[code] = true
		normalized = append(normalized, code)
	}
	sort.Strings(normalized)
	return normalized
}

func normalizePlatformRoleCodes(values []string) []string {
	seen := map[string]bool{}
	normalized := []string{}
	for _, value := range values {
		code := headerKey(value)
		if code == "" || seen[code] || !isCanonicalPlatformRoleCode(code) {
			continue
		}
		seen[code] = true
		normalized = append(normalized, code)
	}
	sort.Strings(normalized)
	return normalized
}

func isCanonicalAdminRoleCode(code string) bool {
	switch code {
	case "tenant_owner", "tenant_admin", "tenant_staff", "tenant_accountant":
		return true
	default:
		return false
	}
}

func isCanonicalPlatformRoleCode(code string) bool {
	return code == "platform_admin"
}

func assignPlatformUserRoles(ctx context.Context, db *sql.DB, input adminUserRoleInput) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM app_user_roles WHERE user_id = $1::uuid`, input.UserID); err != nil {
		return err
	}
	for _, roleCode := range input.RoleCodes {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO app_user_roles (user_id, role_id)
SELECT $1::uuid, id
FROM app_roles
WHERE code = $2
ON CONFLICT (user_id, role_id) DO NOTHING`,
			input.UserID,
			roleCode,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func isCanonicalAdminPermissionCode(code string) bool {
	switch code {
	case "user.view",
		"user.create",
		"user.update",
		"user.assign_role",
		"role.view",
		"role.update",
		"tenant.view",
		"tenant.create",
		"tenant.update",
		"tenant.switch",
		"subscription.view",
		"subscription.update",
		"student.view",
		"student.create",
		"student.update",
		"school_tree.view",
		"school_tree.update",
		"fee.view",
		"fee.create",
		"fee.update",
		"invoice.view",
		"invoice.create",
		"invoice.update",
		"payment.view",
		"payment.create",
		"payment.reconcile",
		"notification.view",
		"notification.create",
		"notification.send",
		"email_config.view",
		"email_config.update",
		"email_cron.view",
		"email_cron.update",
		"report.view",
		"report.export",
		"dashboard.view",
		"operation_log.view",
		"operation_log.cross_tenant_view",
		"audit_log.view",
		"audit_log.cross_tenant_view":
		return true
	default:
		return false
	}
}
