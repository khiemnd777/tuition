package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc64"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	invoiceStatusDraft        = "draft"
	invoiceStatusUnpaid       = "unpaid"
	invoiceStatusPartial      = "partial"
	invoiceStatusPaid         = "paid"
	invoiceStatusOverpaid     = "overpaid"
	invoiceStatusManualReview = "manual_review"
	invoiceStatusVoid         = "void"
	maxInvoiceGenerationRows  = 2000
)

var invoiceCodePattern = regexp.MustCompile(`^[A-Z0-9]+$`)

type invoiceOptionsResponse struct {
	Schedules   []feeScheduleSummary         `json:"schedules"`
	SchoolYears []masterDataSchoolYearOption `json:"schoolYears"`
	Classes     []masterDataClassOption      `json:"classes"`
}

type invoiceGenerateInput struct {
	TenantID      string `json:"-"`
	FeeScheduleID string `json:"feeScheduleId"`
	BankBIN       string `json:"bankBin"`
	BankAccount   string `json:"bankAccount"`
	IssueDate     string `json:"issueDate,omitempty"`
	DueDate       string `json:"dueDate,omitempty"`
	Regenerate    bool   `json:"regenerate,omitempty"`
}

type invoiceListFilters struct {
	TenantID      string
	FeeScheduleID string
	StudentID     string
	StudentCode   string
	SchoolID      string
	SchoolYearID  string
	ClassID       string
	Grade         string
	PeriodCode    string
	Status        string
}

type invoiceScheduleMeta struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	SchoolYearID   string `json:"schoolYearId"`
	SchoolYearCode string `json:"schoolYearCode"`
	ScopeType      string `json:"scopeType"`
	ClassID        string `json:"classId,omitempty"`
	ClassName      string `json:"className,omitempty"`
	Grade          string `json:"grade,omitempty"`
	PeriodCode     string `json:"periodCode"`
	Month          int    `json:"month,omitempty"`
	Status         string `json:"status"`
}

type invoicePreviewSummary struct {
	StudentCount        int `json:"studentCount"`
	ExistingCount       int `json:"existingCount"`
	ReadyCount          int `json:"readyCount"`
	RegenerableCount    int `json:"regenerableCount"`
	BlockedCount        int `json:"blockedCount"`
	MissingBillingCount int `json:"missingBillingCount"`
	GeneratedCount      int `json:"generatedCount"`
	BaseAmount          int `json:"baseAmount"`
	AdjustmentAmount    int `json:"adjustmentAmount"`
	TotalAmount         int `json:"totalAmount"`
}

type invoicePreview struct {
	Schedule invoiceScheduleMeta   `json:"schedule"`
	Summary  invoicePreviewSummary `json:"summary"`
	Rows     []invoicePreviewRow   `json:"rows"`
	Issues   []invoicePreviewIssue `json:"issues,omitempty"`
}

type invoicePreviewIssue struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	StudentCode string `json:"studentCode,omitempty"`
	Field       string `json:"field,omitempty"`
}

type invoicePreviewRow struct {
	InvoiceID             string                      `json:"invoiceId,omitempty"`
	InvoiceCode           string                      `json:"invoiceCode"`
	Existing              bool                        `json:"existing"`
	GenerationState       string                      `json:"generationState"`
	Regenerable           bool                        `json:"regenerable"`
	BlocksGeneration      bool                        `json:"blocksGeneration"`
	IssueState            string                      `json:"issueState,omitempty"`
	StudentID             string                      `json:"studentId"`
	StudentCode           string                      `json:"studentCode"`
	StudentName           string                      `json:"studentName"`
	ClassID               string                      `json:"classId"`
	ClassName             string                      `json:"className"`
	Grade                 string                      `json:"grade"`
	SchoolYearID          string                      `json:"schoolYearId"`
	SchoolYearCode        string                      `json:"schoolYearCode"`
	FeeScheduleID         string                      `json:"feeScheduleId"`
	PeriodCode            string                      `json:"periodCode"`
	Month                 int                         `json:"month,omitempty"`
	IssuedAt              time.Time                   `json:"issuedAt"`
	DueDate               string                      `json:"dueDate,omitempty"`
	Status                string                      `json:"status"`
	BaseAmount            int                         `json:"baseAmount"`
	AdjustmentAmount      int                         `json:"adjustmentAmount"`
	TotalAmount           int                         `json:"totalAmount"`
	PaidAmount            int                         `json:"paidAmount"`
	OutstandingAmount     int                         `json:"outstandingAmount"`
	CollectionBankBIN     string                      `json:"bankBin"`
	CollectionBankAccount string                      `json:"bankAccount"`
	QRBillNumber          string                      `json:"qrBillNumber"`
	QRNote                string                      `json:"qrNote"`
	ItemCount             int                         `json:"itemCount"`
	AdjustmentCount       int                         `json:"adjustmentCount"`
	BillingRecipientReady bool                        `json:"billingRecipientReady"`
	QRReady               bool                        `json:"qrReady"`
	PDFReady              bool                        `json:"pdfReady"`
	Items                 []invoiceDocumentItem       `json:"items"`
	Adjustments           []invoiceDocumentAdjustment `json:"adjustments,omitempty"`
	PaymentItems          []paymentItem               `json:"paymentItems"`
}

type invoiceSummary struct {
	ID                    string    `json:"id"`
	InvoiceCode           string    `json:"invoiceCode"`
	StudentID             string    `json:"studentId"`
	StudentCode           string    `json:"studentCode"`
	StudentName           string    `json:"studentName"`
	ClassID               string    `json:"classId"`
	ClassName             string    `json:"className"`
	Grade                 string    `json:"grade"`
	SchoolYearID          string    `json:"schoolYearId"`
	SchoolYearCode        string    `json:"schoolYearCode"`
	FeeScheduleID         string    `json:"feeScheduleId"`
	PeriodCode            string    `json:"periodCode"`
	Month                 int       `json:"month,omitempty"`
	IssuedAt              time.Time `json:"issuedAt"`
	DueDate               string    `json:"dueDate,omitempty"`
	Status                string    `json:"status"`
	BaseAmount            int       `json:"baseAmount"`
	AdjustmentAmount      int       `json:"adjustmentAmount"`
	TotalAmount           int       `json:"totalAmount"`
	PaidAmount            int       `json:"paidAmount"`
	OutstandingAmount     int       `json:"outstandingAmount"`
	CollectionBankBIN     string    `json:"bankBin"`
	CollectionBankAccount string    `json:"bankAccount"`
	QRBillNumber          string    `json:"qrBillNumber"`
	QRNote                string    `json:"qrNote"`
	ItemCount             int       `json:"itemCount"`
	AdjustmentCount       int       `json:"adjustmentCount"`
	PaymentIntentCount    int       `json:"paymentIntentCount"`
	MatchedPaymentCount   int       `json:"matchedPaymentCount"`
	SentCount             int       `json:"sentCount"`
	LastSentAt            string    `json:"lastSentAt,omitempty"`
	QRReady               bool      `json:"qrReady"`
	PDFReady              bool      `json:"pdfReady"`
	IssueState            string    `json:"issueState,omitempty"`
}

type invoiceDocument struct {
	invoiceSummary
	StudentID     string                      `json:"studentId"`
	Items         []invoiceDocumentItem       `json:"items"`
	Adjustments   []invoiceDocumentAdjustment `json:"adjustments,omitempty"`
	StatusHistory []invoiceStatusHistoryEntry `json:"statusHistory,omitempty"`
}

type invoiceDocumentItem struct {
	FeeTypeCode  string `json:"feeTypeCode"`
	LabelVI      string `json:"labelVi"`
	LabelEN      string `json:"labelEn"`
	Amount       int    `json:"amount"`
	DisplayOrder int    `json:"displayOrder"`
}

type invoiceDocumentAdjustment struct {
	AdjustmentType string `json:"adjustmentType"`
	FeeTypeCode    string `json:"feeTypeCode,omitempty"`
	LabelVI        string `json:"labelVi"`
	LabelEN        string `json:"labelEn"`
	Amount         int    `json:"amount"`
	Delta          int    `json:"delta"`
	Reason         string `json:"reason"`
}

type invoiceExistingRef struct {
	ID          string
	InvoiceCode string
	Status      string
	PaidAmount  int
}

type invoiceStatusHistoryEntry struct {
	FromStatus string    `json:"fromStatus"`
	ToStatus   string    `json:"toStatus"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"createdAt"`
}

func handleInvoiceOptions(w http.ResponseWriter, r *http.Request) {
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

	schedules, err := listFeeScheduleSummaries(r.Context(), db, feeScheduleListFilters{TenantID: tenantID})
	if err != nil {
		http.Error(w, "cannot load fee schedules", http.StatusInternalServerError)
		return
	}
	options, err := listMasterDataOptions(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load master data options", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, invoiceOptionsResponse{
		Schedules:   schedules,
		SchoolYears: options.SchoolYears,
		Classes:     options.Classes,
	})
}

func handleInvoiceList(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query()
	invoices, err := listInvoiceSummaries(r.Context(), db, invoiceListFilters{
		TenantID:      tenantID,
		FeeScheduleID: strings.TrimSpace(query.Get("feeScheduleId")),
		StudentID:     strings.TrimSpace(query.Get("studentId")),
		StudentCode:   normalizeStudentCode(query.Get("studentCode")),
		SchoolYearID:  strings.TrimSpace(query.Get("schoolYearId")),
		ClassID:       strings.TrimSpace(query.Get("classId")),
		Grade:         normalizeGrade(query.Get("grade")),
		PeriodCode:    strings.TrimSpace(query.Get("periodCode")),
		Status:        headerKey(query.Get("status")),
	})
	if err != nil {
		http.Error(w, "cannot load invoices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"invoices": invoices})
}

func handleInvoicePreview(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	input, issues, ok := decodeInvoiceGenerateInput(w, r, false)
	if !ok {
		return
	}
	if len(issues) > 0 {
		writeJSON(w, http.StatusOK, invoicePreview{Issues: issues})
		return
	}

	input.TenantID = tenantID
	preview, err := loadInvoicePreview(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, preview)
}

func handleInvoiceGenerate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	input, issues, ok := decodeInvoiceGenerateInput(w, r, true)
	if !ok {
		return
	}
	if len(issues) > 0 {
		writeJSON(w, http.StatusBadRequest, invoicePreview{Issues: issues})
		return
	}
	input.TenantID = tenantID

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	preview, err := buildInvoicePreviewFromDB(r.Context(), db, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(preview.Issues) > 0 {
		writeJSON(w, http.StatusBadRequest, preview)
		return
	}

	saved, err := saveGeneratedInvoices(r.Context(), db, input, preview)
	if err != nil {
		http.Error(w, "cannot generate invoices: "+err.Error(), http.StatusInternalServerError)
		return
	}
	invoices, err := listInvoiceSummaries(r.Context(), db, invoiceListFilters{
		TenantID:      tenantID,
		FeeScheduleID: input.FeeScheduleID,
		SchoolYearID:  saved.Schedule.SchoolYearID,
		PeriodCode:    saved.Schedule.PeriodCode,
	})
	if err != nil {
		http.Error(w, "cannot load invoices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"preview":  saved,
		"invoices": invoices,
	})
}

func handleInvoicePayment(w http.ResponseWriter, r *http.Request) {
	invoice, err := loadInvoiceFromRequest(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item := buildQRItem(paymentRowFromInvoice(invoice), 360)
	if len(item.Errors) > 0 {
		writeJSON(w, http.StatusBadRequest, item)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func handleInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	invoice, err := loadInvoiceFromRequest(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, invoice)
}

func handleInvoicePDF(w http.ResponseWriter, r *http.Request) {
	invoice, err := loadInvoiceFromRequest(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item := buildQRItem(paymentRowFromInvoice(invoice), 512)
	if len(item.Errors) > 0 {
		http.Error(w, strings.Join(item.Errors, "; "), http.StatusBadRequest)
		return
	}
	pdf, err := renderInvoicePDF(invoice, item)
	if err != nil {
		http.Error(w, "cannot render invoice pdf", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.pdf"`, invoice.InvoiceCode))
	_, _ = w.Write(pdf)
}

func decodeInvoiceGenerateInput(w http.ResponseWriter, r *http.Request, requireBank bool) (invoiceGenerateInput, []invoicePreviewIssue, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input invoiceGenerateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return input, nil, false
	}
	input = normalizeInvoiceGenerateInput(input)
	return input, validateInvoiceGenerateInput(input, requireBank), true
}

func normalizeInvoiceGenerateInput(input invoiceGenerateInput) invoiceGenerateInput {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.FeeScheduleID = strings.TrimSpace(input.FeeScheduleID)
	input.BankBIN = onlyDigits(input.BankBIN)
	input.BankAccount = cleanAccount(input.BankAccount)
	input.IssueDate = strings.TrimSpace(input.IssueDate)
	input.DueDate = strings.TrimSpace(input.DueDate)
	return input
}

func validateInvoiceGenerateInput(input invoiceGenerateInput, requireBank bool) []invoicePreviewIssue {
	issues := []invoicePreviewIssue{}
	if input.FeeScheduleID == "" {
		issues = append(issues, invoicePreviewIssue{Type: "missing_fee_schedule", Field: "feeScheduleId", Message: "feeScheduleId is required"})
	}
	if requireBank {
		if len(input.BankBIN) != 6 {
			issues = append(issues, invoicePreviewIssue{Type: "invalid_bank_bin", Field: "bankBin", Message: "bankBin must be 6 digits"})
		}
		if input.BankAccount == "" {
			issues = append(issues, invoicePreviewIssue{Type: "missing_bank_account", Field: "bankAccount", Message: "bankAccount is required"})
		}
		if len(input.BankAccount) > 19 {
			issues = append(issues, invoicePreviewIssue{Type: "invalid_bank_account", Field: "bankAccount", Message: "bankAccount max length is 19 characters"})
		}
		for _, message := range validateRow(paymentRow{BankBIN: input.BankBIN, BankAccount: input.BankAccount}) {
			if strings.Contains(message, "not in VietQR bank list") {
				issues = append(issues, invoicePreviewIssue{Type: "unknown_bank_bin", Field: "bankBin", Message: message})
			}
		}
	}
	if input.IssueDate != "" {
		if _, err := parseInvoiceDateTime(input.IssueDate, time.Now()); err != nil {
			issues = append(issues, invoicePreviewIssue{Type: "invalid_issue_date", Field: "issueDate", Message: "issueDate must be YYYY-MM-DD or RFC3339"})
		}
	}
	if input.DueDate != "" {
		if _, err := parseInvoiceDate(input.DueDate); err != nil {
			issues = append(issues, invoicePreviewIssue{Type: "invalid_due_date", Field: "dueDate", Message: "dueDate must be YYYY-MM-DD"})
		}
	}
	return issues
}

func loadInvoicePreview(ctx context.Context, input invoiceGenerateInput) (invoicePreview, error) {
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return invoicePreview{}, err
	}
	defer db.Close()
	return buildInvoicePreviewFromDB(ctx, db, input)
}

func buildInvoicePreviewFromDB(ctx context.Context, db *sql.DB, input invoiceGenerateInput) (invoicePreview, error) {
	schedule, meta, err := loadFeeScheduleForInvoice(ctx, db, input.FeeScheduleID, input.TenantID)
	if err != nil {
		return invoicePreview{}, err
	}
	schedule.TenantID = input.TenantID
	students, err := loadFeeScheduleStudents(ctx, db, schedule)
	if err != nil {
		return invoicePreview{}, errors.New("cannot load students for invoice generation")
	}
	if len(students) > maxInvoiceGenerationRows {
		return invoicePreview{}, fmt.Errorf("too many students, max is %d", maxInvoiceGenerationRows)
	}
	feePreview, feeIssues := buildFeeSchedulePreview(schedule, students)
	existing, err := loadExistingInvoiceRefs(ctx, db, input.FeeScheduleID, input.TenantID)
	if err != nil {
		return invoicePreview{}, err
	}
	preview := buildInvoicePreview(meta, input, feePreview, existing, time.Now())
	for _, issue := range feeIssues {
		preview.Issues = append(preview.Issues, invoicePreviewIssue{
			Type:        issue.Type,
			Message:     issue.Message,
			StudentCode: issue.StudentCode,
			Field:       issue.Field,
		})
	}
	return preview, nil
}

func buildInvoicePreview(meta invoiceScheduleMeta, input invoiceGenerateInput, feePreview feeSchedulePreview, existing map[string]invoiceExistingRef, now time.Time) invoicePreview {
	issuedAt, err := parseInvoiceDateTime(input.IssueDate, now)
	if err != nil {
		issuedAt = now
	}
	dueDate := ""
	if input.DueDate != "" {
		if parsed, err := parseInvoiceDate(input.DueDate); err == nil {
			dueDate = parsed.Format("2006-01-02")
		}
	}

	preview := invoicePreview{
		Schedule: meta,
		Rows:     []invoicePreviewRow{},
	}
	for _, feeRow := range feePreview.Rows {
		ref := existing[strings.ToLower(feeRow.StudentID)]
		invoiceCode := stableInvoiceCode(meta.ID, feeRow.StudentID, meta.PeriodCode, feeRow.StudentCode)
		if ref.InvoiceCode != "" {
			invoiceCode = ref.InvoiceCode
		}
		row := invoicePreviewRow{
			InvoiceID:             ref.ID,
			InvoiceCode:           invoiceCode,
			Existing:              ref.ID != "",
			StudentID:             feeRow.StudentID,
			StudentCode:           feeRow.StudentCode,
			StudentName:           feeRow.StudentName,
			ClassID:               feeRowClassID(feeRow),
			ClassName:             feeRow.ClassName,
			Grade:                 feeRow.Grade,
			SchoolYearID:          meta.SchoolYearID,
			SchoolYearCode:        firstNonEmpty(feeRow.SchoolYearCode, meta.SchoolYearCode),
			FeeScheduleID:         meta.ID,
			PeriodCode:            meta.PeriodCode,
			Month:                 meta.Month,
			IssuedAt:              issuedAt.UTC(),
			DueDate:               dueDate,
			Status:                firstNonEmpty(ref.Status, deriveInvoiceStatus(feeRow.TotalAmount, 0)),
			BaseAmount:            feeRow.BaseAmount,
			AdjustmentAmount:      feeRow.AdjustmentAmount,
			TotalAmount:           feeRow.TotalAmount,
			PaidAmount:            ref.PaidAmount,
			OutstandingAmount:     maxInt(feeRow.TotalAmount-ref.PaidAmount, 0),
			CollectionBankBIN:     input.BankBIN,
			CollectionBankAccount: input.BankAccount,
			QRBillNumber:          invoiceCode,
			ItemCount:             len(feeRow.Items),
			AdjustmentCount:       len(feeRow.Adjustments),
			BillingRecipientReady: feeRow.BillingReady,
			Items:                 invoiceItemsFromFeePreview(feeRow.Items),
			Adjustments:           invoiceAdjustmentsFromFeePreview(feeRow.Adjustments),
		}
		row.QRNote = invoiceQRNote(row)
		row.PaymentItems = invoicePaymentItems(row)
		classifyInvoicePreviewRow(&row, input)
		preview.Rows = append(preview.Rows, row)
		preview.Summary.StudentCount++
		if row.Existing {
			preview.Summary.ExistingCount++
		}
		if row.Regenerable {
			preview.Summary.RegenerableCount++
		}
		if row.BlocksGeneration {
			preview.Summary.BlockedCount++
		} else if row.GenerationState == "ready_to_generate" || row.GenerationState == "ready_to_regenerate" {
			preview.Summary.ReadyCount++
		}
		if !row.BillingRecipientReady {
			preview.Summary.MissingBillingCount++
		}
		preview.Summary.BaseAmount += row.BaseAmount
		preview.Summary.AdjustmentAmount += row.AdjustmentAmount
		preview.Summary.TotalAmount += row.TotalAmount
		if row.GenerationState == "blocked_paid_regenerate" {
			preview.Issues = append(preview.Issues, invoicePreviewIssue{
				Type:        "cannot_regenerate_paid_invoice",
				StudentCode: row.StudentCode,
				Message:     "invoice has recorded payment amount and cannot be regenerated",
			})
		}
	}
	return preview
}

func classifyInvoicePreviewRow(row *invoicePreviewRow, input invoiceGenerateInput) {
	row.Regenerable = row.Existing && row.PaidAmount == 0
	row.QRReady = row.QRBillNumber != "" && row.CollectionBankBIN != "" && row.CollectionBankAccount != ""
	row.PDFReady = row.Existing && row.InvoiceID != ""
	row.IssueState = "ready"
	switch {
	case row.Existing && input.Regenerate && row.PaidAmount > 0:
		row.GenerationState = "blocked_paid_regenerate"
		row.BlocksGeneration = true
		row.IssueState = "blocking"
	case row.Existing && input.Regenerate:
		row.GenerationState = "ready_to_regenerate"
	case row.Existing:
		row.GenerationState = "already_generated"
	case !row.BillingRecipientReady:
		row.GenerationState = "ready_to_generate"
		row.IssueState = "warning"
	default:
		row.GenerationState = "ready_to_generate"
	}
}

func feeRowClassID(row feeSchedulePreviewRow) string {
	return strings.TrimSpace(row.ClassID)
}

func invoiceItemsFromFeePreview(items []feeSchedulePreviewItem) []invoiceDocumentItem {
	out := make([]invoiceDocumentItem, 0, len(items))
	for idx, item := range items {
		out = append(out, invoiceDocumentItem{
			FeeTypeCode:  item.FeeTypeCode,
			LabelVI:      item.LabelVI,
			LabelEN:      item.LabelEN,
			Amount:       item.Amount,
			DisplayOrder: idx + 1,
		})
	}
	return out
}

func invoiceAdjustmentsFromFeePreview(adjustments []feeSchedulePreviewAdjustment) []invoiceDocumentAdjustment {
	out := make([]invoiceDocumentAdjustment, 0, len(adjustments))
	for _, item := range adjustments {
		out = append(out, invoiceDocumentAdjustment{
			AdjustmentType: item.AdjustmentType,
			FeeTypeCode:    item.FeeTypeCode,
			LabelVI:        item.LabelVI,
			LabelEN:        item.LabelEN,
			Amount:         item.Amount,
			Delta:          item.Delta,
			Reason:         item.Reason,
		})
	}
	return out
}

func saveGeneratedInvoices(ctx context.Context, db *sql.DB, input invoiceGenerateInput, preview invoicePreview) (invoicePreview, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return invoicePreview{}, err
	}
	defer tx.Rollback()

	saved := preview
	saved.Summary.GeneratedCount = 0
	for idx := range saved.Rows {
		row := &saved.Rows[idx]
		if row.Existing && !input.Regenerate {
			continue
		}
		if row.Existing && input.Regenerate {
			if row.PaidAmount > 0 {
				saved.Issues = append(saved.Issues, invoicePreviewIssue{
					Type:        "cannot_regenerate_paid_invoice",
					StudentCode: row.StudentCode,
					Message:     "invoice has recorded payment amount and cannot be regenerated",
				})
				continue
			}
			oldStatus := row.Status
			row.Status = deriveInvoiceStatus(row.TotalAmount, row.PaidAmount)
			if err := updateInvoice(ctx, tx, *row); err != nil {
				return invoicePreview{}, err
			}
			if err := replaceInvoiceSnapshots(ctx, tx, row.InvoiceID, *row); err != nil {
				return invoicePreview{}, err
			}
			if err := insertInvoiceStatusHistory(ctx, tx, row.InvoiceID, oldStatus, row.Status, "regenerated from fee schedule"); err != nil {
				return invoicePreview{}, err
			}
			saved.Summary.GeneratedCount++
			continue
		}

		invoiceID, err := insertInvoice(ctx, tx, *row)
		if err != nil {
			return invoicePreview{}, err
		}
		row.InvoiceID = invoiceID
		if err := replaceInvoiceSnapshots(ctx, tx, invoiceID, *row); err != nil {
			return invoicePreview{}, err
		}
		if err := insertInvoiceStatusHistory(ctx, tx, invoiceID, "", row.Status, "generated from fee schedule"); err != nil {
			return invoicePreview{}, err
		}
		saved.Summary.GeneratedCount++
	}

	if err := tx.Commit(); err != nil {
		return invoicePreview{}, err
	}
	return saved, nil
}

func insertInvoice(ctx context.Context, exec masterDataExecutor, row invoicePreviewRow) (string, error) {
	snapshot, err := invoiceSourceSnapshot(row)
	if err != nil {
		return "", err
	}
	var invoiceID string
	err = exec.QueryRowContext(ctx, `
INSERT INTO invoices (
	fee_schedule_id, student_id, school_year_id, class_id, invoice_code,
	student_code, student_name, class_name, grade, school_year_code,
	period_code, month, issued_at, due_date, status,
	base_amount, adjustment_amount, total_amount, paid_amount,
	collection_bank_bin, collection_bank_account, qr_bill_number, qr_note,
	source_snapshot
)
VALUES (
	$1::uuid, $2::uuid, $3::uuid, $4::uuid, $5,
	$6, $7, $8, $9, $10,
	$11, $12, $13, $14, $15,
	$16, $17, $18, $19,
	$20, $21, $22, $23,
	$24::jsonb
)
RETURNING id::text`,
		row.FeeScheduleID,
		row.StudentID,
		row.SchoolYearID,
		row.ClassID,
		row.InvoiceCode,
		row.StudentCode,
		row.StudentName,
		row.ClassName,
		row.Grade,
		row.SchoolYearCode,
		row.PeriodCode,
		nullableInt(row.Month),
		row.IssuedAt,
		nullableDateString(row.DueDate),
		row.Status,
		row.BaseAmount,
		row.AdjustmentAmount,
		row.TotalAmount,
		row.PaidAmount,
		row.CollectionBankBIN,
		row.CollectionBankAccount,
		row.QRBillNumber,
		row.QRNote,
		snapshot,
	).Scan(&invoiceID)
	return invoiceID, err
}

func updateInvoice(ctx context.Context, exec masterDataExecutor, row invoicePreviewRow) error {
	snapshot, err := invoiceSourceSnapshot(row)
	if err != nil {
		return err
	}
	_, err = exec.ExecContext(ctx, `
UPDATE invoices
SET school_year_id = $2::uuid,
	class_id = $3::uuid,
	student_code = $4,
	student_name = $5,
	class_name = $6,
	grade = $7,
	school_year_code = $8,
	period_code = $9,
	month = $10,
	issued_at = $11,
	due_date = $12,
	status = $13,
	base_amount = $14,
	adjustment_amount = $15,
	total_amount = $16,
	collection_bank_bin = $17,
	collection_bank_account = $18,
	qr_bill_number = $19,
	qr_note = $20,
	source_snapshot = $21::jsonb
WHERE id = $1::uuid`,
		row.InvoiceID,
		row.SchoolYearID,
		row.ClassID,
		row.StudentCode,
		row.StudentName,
		row.ClassName,
		row.Grade,
		row.SchoolYearCode,
		row.PeriodCode,
		nullableInt(row.Month),
		row.IssuedAt,
		nullableDateString(row.DueDate),
		row.Status,
		row.BaseAmount,
		row.AdjustmentAmount,
		row.TotalAmount,
		row.CollectionBankBIN,
		row.CollectionBankAccount,
		row.QRBillNumber,
		row.QRNote,
		snapshot,
	)
	return err
}

func replaceInvoiceSnapshots(ctx context.Context, exec masterDataExecutor, invoiceID string, row invoicePreviewRow) error {
	if _, err := exec.ExecContext(ctx, `DELETE FROM invoice_items WHERE invoice_id = $1::uuid`, invoiceID); err != nil {
		return err
	}
	if _, err := exec.ExecContext(ctx, `DELETE FROM invoice_adjustments WHERE invoice_id = $1::uuid`, invoiceID); err != nil {
		return err
	}
	for _, item := range row.Items {
		if _, err := exec.ExecContext(ctx, `
INSERT INTO invoice_items (invoice_id, fee_type_code, label_vi, label_en, amount, display_order)
VALUES ($1::uuid, $2, $3, $4, $5, $6)`,
			invoiceID,
			item.FeeTypeCode,
			item.LabelVI,
			item.LabelEN,
			item.Amount,
			item.DisplayOrder,
		); err != nil {
			return err
		}
	}
	for _, adjustment := range row.Adjustments {
		if _, err := exec.ExecContext(ctx, `
INSERT INTO invoice_adjustments (invoice_id, adjustment_type, fee_type_code, label_vi, label_en, amount, delta, reason)
VALUES ($1::uuid, $2, $3, $4, $5, $6, $7, $8)`,
			invoiceID,
			adjustment.AdjustmentType,
			adjustment.FeeTypeCode,
			adjustment.LabelVI,
			adjustment.LabelEN,
			adjustment.Amount,
			adjustment.Delta,
			adjustment.Reason,
		); err != nil {
			return err
		}
	}
	return nil
}

func insertInvoiceStatusHistory(ctx context.Context, exec masterDataExecutor, invoiceID string, fromStatus string, toStatus string, reason string) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO invoice_status_history (invoice_id, from_status, to_status, reason)
VALUES ($1::uuid, $2, $3, $4)`,
		invoiceID,
		fromStatus,
		toStatus,
		reason,
	)
	return err
}

func loadFeeScheduleForInvoice(ctx context.Context, db *sql.DB, scheduleID string, tenantID string) (feeScheduleInput, invoiceScheduleMeta, error) {
	var input feeScheduleInput
	var meta invoiceScheduleMeta
	var classID sql.NullString
	var className sql.NullString
	var month sql.NullInt64
	err := db.QueryRowContext(ctx, `
SELECT fs.id::text,
	fs.school_year_id::text,
	sy.code,
	fs.scope_type,
	COALESCE(fs.class_id::text, ''),
	COALESCE(c.name, ''),
	fs.grade,
	fs.period_code,
	fs.month,
	fs.name,
	fs.status
FROM fee_schedules fs
JOIN school_years sy ON sy.id = fs.school_year_id
JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN classes c ON c.id = fs.class_id
WHERE fs.id = $1::uuid
	AND sc.tenant_id = $2::uuid
	AND fs.status <> 'archived'`, scheduleID, tenantID).Scan(
		&meta.ID,
		&meta.SchoolYearID,
		&meta.SchoolYearCode,
		&meta.ScopeType,
		&classID,
		&className,
		&meta.Grade,
		&meta.PeriodCode,
		&month,
		&meta.Name,
		&meta.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return input, meta, errors.New("fee schedule not found")
	}
	if err != nil {
		return input, meta, err
	}
	if classID.Valid {
		meta.ClassID = classID.String
	}
	if className.Valid {
		meta.ClassName = className.String
	}
	if month.Valid {
		meta.Month = int(month.Int64)
	}

	input = feeScheduleInput{
		TenantID:     tenantID,
		ID:           meta.ID,
		SchoolYearID: meta.SchoolYearID,
		ClassID:      meta.ClassID,
		Grade:        meta.Grade,
		PeriodCode:   meta.PeriodCode,
		Month:        meta.Month,
		Name:         meta.Name,
		Status:       meta.Status,
		Items:        []feeScheduleItemInput{},
		Adjustments:  []studentFeeAdjustmentInput{},
	}

	itemRows, err := db.QueryContext(ctx, `
SELECT fsi.fee_type_id::text, ft.code, fsi.label_vi, fsi.label_en, fsi.amount, fsi.display_order
FROM fee_schedule_items fsi
JOIN fee_types ft ON ft.id = fsi.fee_type_id
WHERE fsi.schedule_id = $1::uuid
ORDER BY fsi.display_order, fsi.id`, scheduleID)
	if err != nil {
		return input, meta, err
	}
	defer itemRows.Close()
	for itemRows.Next() {
		var item feeScheduleItemInput
		if err := itemRows.Scan(&item.FeeTypeID, &item.FeeTypeCode, &item.LabelVI, &item.LabelEN, &item.Amount, &item.DisplayOrder); err != nil {
			return input, meta, err
		}
		input.Items = append(input.Items, item)
	}
	if err := itemRows.Err(); err != nil {
		return input, meta, err
	}

	adjustmentRows, err := db.QueryContext(ctx, `
SELECT sfa.student_id::text,
	s.student_code,
	sfa.adjustment_type,
	COALESCE(sfa.fee_type_id::text, ''),
	COALESCE(ft.code, ''),
	sfa.label_vi,
	sfa.label_en,
	sfa.amount,
	sfa.reason
FROM student_fee_adjustments sfa
JOIN students s ON s.id = sfa.student_id
LEFT JOIN fee_types ft ON ft.id = sfa.fee_type_id
WHERE sfa.schedule_id = $1::uuid
	AND sfa.status = 'active'
ORDER BY s.student_code, sfa.created_at, sfa.id`, scheduleID)
	if err != nil {
		return input, meta, err
	}
	defer adjustmentRows.Close()
	for adjustmentRows.Next() {
		var adjustment studentFeeAdjustmentInput
		if err := adjustmentRows.Scan(
			&adjustment.StudentID,
			&adjustment.StudentCode,
			&adjustment.AdjustmentType,
			&adjustment.FeeTypeID,
			&adjustment.FeeTypeCode,
			&adjustment.LabelVI,
			&adjustment.LabelEN,
			&adjustment.Amount,
			&adjustment.Reason,
		); err != nil {
			return input, meta, err
		}
		input.Adjustments = append(input.Adjustments, adjustment)
	}
	if err := adjustmentRows.Err(); err != nil {
		return input, meta, err
	}

	return normalizeFeeScheduleInput(input), meta, nil
}

func loadExistingInvoiceRefs(ctx context.Context, db *sql.DB, feeScheduleID string, tenantID string) (map[string]invoiceExistingRef, error) {
	rows, err := db.QueryContext(ctx, `
SELECT student_id::text, id::text, invoice_code, status, paid_amount
FROM invoices i
JOIN school_years sy ON sy.id = i.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE i.fee_schedule_id = $1::uuid
	AND sc.tenant_id = $2::uuid
	AND i.status <> 'void'`, feeScheduleID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := map[string]invoiceExistingRef{}
	for rows.Next() {
		var studentID string
		var ref invoiceExistingRef
		if err := rows.Scan(&studentID, &ref.ID, &ref.InvoiceCode, &ref.Status, &ref.PaidAmount); err != nil {
			return nil, err
		}
		existing[strings.ToLower(studentID)] = ref
	}
	return existing, rows.Err()
}

func listInvoiceSummaries(ctx context.Context, db *sql.DB, filters invoiceListFilters) ([]invoiceSummary, error) {
	conditions := []string{"1 = 1"}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.TenantID != "" {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM school_years tenant_sy
			JOIN schools tenant_school ON tenant_school.id = tenant_sy.school_id
			WHERE tenant_sy.id = i.school_year_id
				AND tenant_school.tenant_id = `+addArg(filters.TenantID)+`::uuid
		)`)
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "i.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.FeeScheduleID != "" {
		conditions = append(conditions, "i.fee_schedule_id = "+addArg(filters.FeeScheduleID)+"::uuid")
	}
	if filters.StudentID != "" {
		conditions = append(conditions, "i.student_id = "+addArg(filters.StudentID)+"::uuid")
	}
	if filters.StudentCode != "" {
		conditions = append(conditions, "i.student_code = "+addArg(filters.StudentCode))
	}
	if filters.SchoolID != "" {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM school_years sy_filter
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
	if filters.Status != "" {
		conditions = append(conditions, "i.status = "+addArg(filters.Status))
	}

	query := `
SELECT i.id::text,
	i.invoice_code,
	i.student_id::text,
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
ORDER BY i.issued_at DESC, i.period_code DESC, i.class_name, i.student_code
LIMIT 1000`

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []invoiceSummary{}
	for rows.Next() {
		var invoice invoiceSummary
		var month sql.NullInt64
		var dueDate sql.NullTime
		var lastSentAt sql.NullTime
		if err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceCode,
			&invoice.StudentID,
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
		invoice.IssueState = invoiceSummaryIssueState(invoice)
		invoices = append(invoices, invoice)
	}
	return invoices, rows.Err()
}

func invoiceSummaryIssueState(invoice invoiceSummary) string {
	switch {
	case invoice.Status == invoiceStatusManualReview || invoice.Status == invoiceStatusOverpaid:
		return "review"
	case invoice.OutstandingAmount > 0:
		return "open"
	default:
		return "ready"
	}
}

func loadInvoiceFromRequest(ctx context.Context, r *http.Request) (invoiceDocument, error) {
	tenantID := activeTenantIDFromRequest(r)
	if tenantID == "" {
		return invoiceDocument{}, errors.New("active tenant required")
	}
	invoiceID := strings.TrimSpace(r.URL.Query().Get("id"))
	if invoiceID == "" {
		return invoiceDocument{}, errors.New("invoice id is required")
	}
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return invoiceDocument{}, err
	}
	defer db.Close()
	return loadInvoiceDocument(ctx, db, invoiceID, tenantID)
}

func loadInvoiceDocument(ctx context.Context, db *sql.DB, invoiceID string, tenantID string) (invoiceDocument, error) {
	var doc invoiceDocument
	var month sql.NullInt64
	var dueDate sql.NullTime
	err := db.QueryRowContext(ctx, `
SELECT i.id::text,
	i.invoice_code,
	i.student_id::text,
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
	i.collection_bank_bin,
	i.collection_bank_account,
	i.qr_bill_number,
	i.qr_note
FROM invoices i
JOIN school_years sy ON sy.id = i.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE i.id = $1::uuid
	AND sc.tenant_id = $2::uuid`, invoiceID, tenantID).Scan(
		&doc.ID,
		&doc.InvoiceCode,
		&doc.StudentID,
		&doc.StudentCode,
		&doc.StudentName,
		&doc.ClassID,
		&doc.ClassName,
		&doc.Grade,
		&doc.SchoolYearID,
		&doc.SchoolYearCode,
		&doc.FeeScheduleID,
		&doc.PeriodCode,
		&month,
		&doc.IssuedAt,
		&dueDate,
		&doc.Status,
		&doc.BaseAmount,
		&doc.AdjustmentAmount,
		&doc.TotalAmount,
		&doc.PaidAmount,
		&doc.CollectionBankBIN,
		&doc.CollectionBankAccount,
		&doc.QRBillNumber,
		&doc.QRNote,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return doc, errors.New("invoice not found")
	}
	if err != nil {
		return doc, err
	}
	if month.Valid {
		doc.Month = int(month.Int64)
	}
	if dueDate.Valid {
		doc.DueDate = dueDate.Time.Format("2006-01-02")
	}

	items, err := loadInvoiceDocumentItems(ctx, db, invoiceID)
	if err != nil {
		return doc, err
	}
	adjustments, err := loadInvoiceDocumentAdjustments(ctx, db, invoiceID)
	if err != nil {
		return doc, err
	}
	doc.Items = items
	doc.Adjustments = adjustments
	doc.ItemCount = len(items)
	doc.AdjustmentCount = len(adjustments)
	doc.OutstandingAmount = maxInt(doc.TotalAmount-doc.PaidAmount, 0)
	doc.QRReady = doc.QRBillNumber != "" && doc.CollectionBankBIN != "" && doc.CollectionBankAccount != ""
	doc.PDFReady = doc.ID != ""
	doc.IssueState = invoiceSummaryIssueState(doc.invoiceSummary)
	if err := loadInvoiceOperationalCounts(ctx, db, &doc.invoiceSummary); err != nil {
		return doc, err
	}
	statusHistory, err := loadInvoiceStatusHistory(ctx, db, invoiceID)
	if err != nil {
		return doc, err
	}
	doc.StatusHistory = statusHistory
	return doc, nil
}

func loadInvoiceOperationalCounts(ctx context.Context, db *sql.DB, invoice *invoiceSummary) error {
	var lastSentAt sql.NullTime
	err := db.QueryRowContext(ctx, `
SELECT
	COALESCE(intent_counts.intent_count, 0),
	COALESCE(match_counts.match_count, 0),
	COALESCE(notification_counts.sent_count, 0),
	notification_counts.last_sent_at
FROM invoices i
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
WHERE i.id = $1::uuid`, invoice.ID).Scan(
		&invoice.PaymentIntentCount,
		&invoice.MatchedPaymentCount,
		&invoice.SentCount,
		&lastSentAt,
	)
	if err != nil {
		return err
	}
	if lastSentAt.Valid {
		invoice.LastSentAt = lastSentAt.Time.UTC().Format(time.RFC3339)
	}
	return nil
}

func loadInvoiceStatusHistory(ctx context.Context, db *sql.DB, invoiceID string) ([]invoiceStatusHistoryEntry, error) {
	rows, err := db.QueryContext(ctx, `
SELECT from_status, to_status, reason, created_at
FROM invoice_status_history
WHERE invoice_id = $1::uuid
ORDER BY created_at DESC, id DESC
LIMIT 20`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []invoiceStatusHistoryEntry{}
	for rows.Next() {
		var item invoiceStatusHistoryEntry
		if err := rows.Scan(&item.FromStatus, &item.ToStatus, &item.Reason, &item.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, item)
	}
	return history, rows.Err()
}

func loadInvoiceDocumentItems(ctx context.Context, db *sql.DB, invoiceID string) ([]invoiceDocumentItem, error) {
	rows, err := db.QueryContext(ctx, `
SELECT fee_type_code, label_vi, label_en, amount, display_order
FROM invoice_items
WHERE invoice_id = $1::uuid
ORDER BY display_order, id`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []invoiceDocumentItem{}
	for rows.Next() {
		var item invoiceDocumentItem
		if err := rows.Scan(&item.FeeTypeCode, &item.LabelVI, &item.LabelEN, &item.Amount, &item.DisplayOrder); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func loadInvoiceDocumentAdjustments(ctx context.Context, db *sql.DB, invoiceID string) ([]invoiceDocumentAdjustment, error) {
	rows, err := db.QueryContext(ctx, `
SELECT adjustment_type, fee_type_code, label_vi, label_en, amount, delta, reason
FROM invoice_adjustments
WHERE invoice_id = $1::uuid
ORDER BY created_at, id`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	adjustments := []invoiceDocumentAdjustment{}
	for rows.Next() {
		var adjustment invoiceDocumentAdjustment
		if err := rows.Scan(
			&adjustment.AdjustmentType,
			&adjustment.FeeTypeCode,
			&adjustment.LabelVI,
			&adjustment.LabelEN,
			&adjustment.Amount,
			&adjustment.Delta,
			&adjustment.Reason,
		); err != nil {
			return nil, err
		}
		adjustments = append(adjustments, adjustment)
	}
	return adjustments, rows.Err()
}

func invoiceSourceSnapshot(row invoicePreviewRow) (string, error) {
	payload := map[string]any{
		"invoiceCode":      row.InvoiceCode,
		"feeScheduleId":    row.FeeScheduleID,
		"periodCode":       row.PeriodCode,
		"studentCode":      row.StudentCode,
		"items":            row.Items,
		"adjustments":      row.Adjustments,
		"baseAmount":       row.BaseAmount,
		"adjustmentAmount": row.AdjustmentAmount,
		"totalAmount":      row.TotalAmount,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func paymentRowFromInvoice(invoice invoiceDocument) paymentRow {
	row := invoicePreviewRow{
		InvoiceCode:           invoice.InvoiceCode,
		StudentCode:           invoice.StudentCode,
		StudentName:           invoice.StudentName,
		ClassName:             invoice.ClassName,
		TotalAmount:           invoice.TotalAmount,
		CollectionBankBIN:     invoice.CollectionBankBIN,
		CollectionBankAccount: invoice.CollectionBankAccount,
		QRBillNumber:          invoice.QRBillNumber,
		QRNote:                invoice.QRNote,
		Items:                 invoice.Items,
		Adjustments:           invoice.Adjustments,
	}
	return invoicePaymentRow(row)
}

func invoicePaymentRow(row invoicePreviewRow) paymentRow {
	return paymentRow{
		ID:           firstNonEmpty(row.InvoiceID, row.InvoiceCode),
		StudentName:  row.StudentName,
		ClassName:    row.ClassName,
		BankBIN:      row.CollectionBankBIN,
		BankAccount:  row.CollectionBankAccount,
		Amount:       row.TotalAmount,
		PaymentItems: invoicePaymentItems(row),
		BillNumber:   row.QRBillNumber,
		Note:         row.QRNote,
	}
}

func invoicePaymentItems(row invoicePreviewRow) []paymentItem {
	if len(row.Adjustments) == 0 {
		items := make([]paymentItem, 0, len(row.Items))
		for _, item := range row.Items {
			items = append(items, paymentItem{
				Label:   item.LabelVI,
				LabelEN: item.LabelEN,
				Amount:  item.Amount,
			})
		}
		return items
	}
	return []paymentItem{{
		Label:   "Tổng hóa đơn " + row.InvoiceCode,
		LabelEN: "Invoice total " + row.InvoiceCode,
		Amount:  row.TotalAmount,
	}}
}

func invoiceQRNote(row invoicePreviewRow) string {
	return cleanANS(strings.TrimSpace("HP "+row.StudentCode+" "+row.PeriodCode), 25)
}

func deriveInvoiceStatus(totalAmount int, paidAmount int) string {
	switch {
	case paidAmount <= 0:
		return invoiceStatusUnpaid
	case paidAmount < totalAmount:
		return invoiceStatusPartial
	case paidAmount == totalAmount:
		return invoiceStatusPaid
	default:
		return invoiceStatusOverpaid
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func stableInvoiceCode(scheduleID string, studentID string, periodCode string, studentCode string) string {
	source := strings.Join([]string{scheduleID, studentID, periodCode, studentCode}, "|")
	sum := crc64.Checksum([]byte(source), crc64.MakeTable(crc64.ISO))
	code := "SUN" + strings.ToUpper(strconv.FormatUint(sum, 36))
	code = cleanANS(code, 25)
	if !invoiceCodePattern.MatchString(code) {
		code = strings.ToUpper(regexp.MustCompile(`[^A-Z0-9]`).ReplaceAllString(code, ""))
	}
	if code == "" {
		return "SUN00000000"
	}
	return code
}

func parseInvoiceDateTime(value string, fallback time.Time) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback.UTC(), nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.UTC(), nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func parseInvoiceDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func nullableDateString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}
