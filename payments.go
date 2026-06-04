package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	paymentProviderManualVietQR = "manual_vietqr"
	paymentProviderSePay        = "sepay"
	paymentProviderPayOS        = "payos"

	paymentIntentStatusActive = "active"

	paymentTransactionStatusUnmatched    = "unmatched"
	paymentTransactionStatusMatched      = "matched"
	paymentTransactionStatusManualReview = "manual_review"

	reconciliationStatusMatched = "matched"
	paymentDirectionIn          = "in"
	paymentCurrencyVND          = "VND"

	maxPaymentWebhookBytes = 1 << 20
)

var nonAlnumPattern = regexp.MustCompile(`[^A-Z0-9]+`)

type paymentProvider struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	DisplayName  string `json:"displayName"`
	ProviderType string `json:"providerType"`
	Status       string `json:"status"`
	Configured   bool   `json:"configured"`
	WebhookPath  string `json:"webhookPath,omitempty"`
}

type paymentIntentRequest struct {
	InvoiceID string `json:"invoiceId"`
	Provider  string `json:"provider"`
}

type paymentIntentSummary struct {
	ID                string    `json:"id"`
	InvoiceID         string    `json:"invoiceId"`
	InvoiceCode       string    `json:"invoiceCode,omitempty"`
	ProviderCode      string    `json:"provider"`
	IntentCode        string    `json:"intentCode"`
	Status            string    `json:"status"`
	Amount            int       `json:"amount"`
	Currency          string    `json:"currency"`
	ProviderReference string    `json:"providerReference,omitempty"`
	PaymentURL        string    `json:"paymentUrl,omitempty"`
	QRPayload         string    `json:"qrPayload,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
}

type paymentIntentResponse struct {
	Intent paymentIntentSummary `json:"intent"`
	QR     *qrItem              `json:"qr,omitempty"`
}

type paymentTransactionSummary struct {
	ID                    string    `json:"id"`
	ProviderCode          string    `json:"provider"`
	ProviderTransactionID string    `json:"providerTransactionId"`
	InvoiceID             string    `json:"invoiceId,omitempty"`
	InvoiceCode           string    `json:"invoiceCode,omitempty"`
	StudentCode           string    `json:"studentCode,omitempty"`
	StudentName           string    `json:"studentName,omitempty"`
	Direction             string    `json:"direction"`
	Amount                int       `json:"amount"`
	Currency              string    `json:"currency"`
	TransactionTime       time.Time `json:"transactionTime"`
	AccountNumber         string    `json:"accountNumber,omitempty"`
	BankName              string    `json:"bankName,omitempty"`
	Description           string    `json:"description,omitempty"`
	ReferenceCode         string    `json:"referenceCode,omitempty"`
	Status                string    `json:"status"`
	MatchType             string    `json:"matchType,omitempty"`
	MatchStatus           string    `json:"matchStatus,omitempty"`
	MatchScore            int       `json:"matchScore,omitempty"`
	AmountApplied         int       `json:"amountApplied,omitempty"`
	MatchReason           string    `json:"matchReason,omitempty"`
}

type paymentReconciliationSummary struct {
	InvoiceCount      int     `json:"invoiceCount"`
	TotalReceivable   int     `json:"totalReceivable"`
	TotalCollected    int     `json:"totalCollected"`
	OutstandingAmount int     `json:"outstandingAmount"`
	CollectionRate    float64 `json:"collectionRate"`
	UnpaidCount       int     `json:"unpaidCount"`
	PaidCount         int     `json:"paidCount"`
	UnmatchedCount    int     `json:"unmatchedCount"`
	MatchedCount      int     `json:"matchedCount"`
	PartialCount      int     `json:"partialCount"`
	OverpaidCount     int     `json:"overpaidCount"`
	ManualReviewCount int     `json:"manualReviewCount"`
}

type paymentMatchSummary struct {
	ID                    string    `json:"id"`
	InvoiceID             string    `json:"invoiceId"`
	InvoiceCode           string    `json:"invoiceCode"`
	TransactionID         string    `json:"transactionId"`
	ProviderCode          string    `json:"provider"`
	ProviderTransactionID string    `json:"providerTransactionId,omitempty"`
	MatchType             string    `json:"matchType"`
	Status                string    `json:"status"`
	Score                 int       `json:"score"`
	AmountApplied         int       `json:"amountApplied"`
	Reason                string    `json:"reason,omitempty"`
	CreatedAt             time.Time `json:"createdAt"`
}

type paymentReconciliationResponse struct {
	Providers    []paymentProvider                `json:"providers"`
	Schools      []masterDataSchoolOption         `json:"schools,omitempty"`
	SchoolYears  []masterDataSchoolYearOption     `json:"schoolYears,omitempty"`
	Classes      []masterDataClassOption          `json:"classes,omitempty"`
	Summary      paymentReconciliationSummary     `json:"summary"`
	Invoices     []invoiceSummary                 `json:"invoices"`
	Transactions []paymentTransactionSummary      `json:"transactions"`
	Intents      map[string]paymentIntentSummary  `json:"intents,omitempty"`
	Matches      map[string][]paymentMatchSummary `json:"matches,omitempty"`
}

type paymentTransactionListFilters struct {
	TenantID string
	Provider string
	Status   string
	Limit    int
}

type providerEventRecord struct {
	ID        string
	Duplicate bool
	Status    string
}

type normalizedPaymentTransaction struct {
	ProviderCode          string
	ProviderTransactionID string
	ProviderEventID       string
	PaymentIntentID       string
	InvoiceID             string
	Direction             string
	Amount                int
	Currency              string
	TransactionTime       time.Time
	AccountNumber         string
	AccountName           string
	BankBIN               string
	BankName              string
	Description           string
	ReferenceCode         string
	Status                string
	RawPayload            map[string]any
}

type insertedPaymentTransaction struct {
	ID        string
	Duplicate bool
	Summary   paymentTransactionSummary
}

type paymentInvoiceCandidate struct {
	ID                    string
	InvoiceCode           string
	QRBillNumber          string
	CollectionBankAccount string
	Status                string
	TotalAmount           int
	PaidAmount            int
	ProviderReferences    []string
}

type paymentMatchCandidate struct {
	Invoice       paymentInvoiceCandidate
	Score         int
	MatchType     string
	AmountApplied int
	Reason        string
}

type invoicePaymentStatusRefresh struct {
	InvoiceID  string
	OldStatus  string
	NewStatus  string
	PaidAmount int
}

func (item invoicePaymentStatusRefresh) BecamePaid() bool {
	return invoiceStatusBecamePaid(item.OldStatus, item.NewStatus)
}

func invoiceStatusBecamePaid(oldStatus string, newStatus string) bool {
	return oldStatus != invoiceStatusPaid && newStatus == invoiceStatusPaid
}

type manualCashReceiptRequest struct {
	InvoiceID        string `json:"invoiceId"`
	Amount           int    `json:"amount"`
	CollectorUserID  string `json:"collectorUserId,omitempty"`
	CollectorName    string `json:"collectorName"`
	ReceiptReference string `json:"receiptReference"`
	PaidAt           string `json:"paidAt,omitempty"`
	Reason           string `json:"reason"`
	Note             string `json:"note,omitempty"`
}

type manualCashReceiptResponse struct {
	ReceiptID   string                    `json:"receiptId"`
	Transaction paymentTransactionSummary `json:"transaction"`
	Invoice     invoiceSummary            `json:"invoice"`
}

type payOSConfig struct {
	ClientID    string
	APIKey      string
	ChecksumKey string
	ReturnURL   string
	CancelURL   string
	APIBaseURL  string
}

type payOSCreatePaymentResult struct {
	OrderCode       int64          `json:"orderCode"`
	PaymentLinkID   string         `json:"paymentLinkId"`
	CheckoutURL     string         `json:"checkoutUrl"`
	QRCode          string         `json:"qrCode"`
	Status          string         `json:"status"`
	RequestPayload  map[string]any `json:"requestPayload"`
	ResponsePayload map[string]any `json:"responsePayload"`
}

func handlePaymentProviders(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	providers, err := listPaymentProviders(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load payment providers", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"providers": providers})
}

func handlePaymentIntentCreate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input paymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.InvoiceID = strings.TrimSpace(input.InvoiceID)
	input.Provider = headerKey(input.Provider)
	if input.Provider == "" {
		input.Provider = paymentProviderManualVietQR
	}
	if input.InvoiceID == "" {
		http.Error(w, "invoiceId is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	provider, err := loadPaymentProviderByCode(r.Context(), db, input.Provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	invoice, err := loadInvoiceDocument(r.Context(), db, input.InvoiceID, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := createPaymentIntentForInvoice(r.Context(), db, provider, invoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func handlePaymentTransactions(w http.ResponseWriter, r *http.Request) {
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
	transactions, err := listPaymentTransactions(r.Context(), db, paymentTransactionListFilters{
		TenantID: tenantID,
		Provider: headerKey(query.Get("provider")),
		Status:   headerKey(query.Get("status")),
		Limit:    parsePositiveInt(query.Get("limit"), 200),
	})
	if err != nil {
		http.Error(w, "cannot load payment transactions", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"transactions": transactions})
}

func handlePaymentReconciliation(w http.ResponseWriter, r *http.Request) {
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

	providers, err := listPaymentProviders(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load payment providers", http.StatusInternalServerError)
		return
	}
	options, err := listMasterDataOptions(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load master data options", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()
	invoices, err := listInvoiceSummaries(r.Context(), db, invoiceListFilters{
		TenantID:     tenantID,
		SchoolID:     strings.TrimSpace(query.Get("schoolId")),
		SchoolYearID: strings.TrimSpace(query.Get("schoolYearId")),
		ClassID:      strings.TrimSpace(query.Get("classId")),
		Grade:        normalizeGrade(query.Get("grade")),
		PeriodCode:   strings.TrimSpace(query.Get("periodCode")),
		Status:       headerKey(query.Get("invoiceStatus")),
	})
	if err != nil {
		http.Error(w, "cannot load invoices", http.StatusInternalServerError)
		return
	}

	transactions, err := listPaymentTransactions(r.Context(), db, paymentTransactionListFilters{
		TenantID: tenantID,
		Provider: headerKey(query.Get("provider")),
		Status:   headerKey(query.Get("transactionStatus")),
		Limit:    300,
	})
	if err != nil {
		http.Error(w, "cannot load payment transactions", http.StatusInternalServerError)
		return
	}

	intents, err := listLatestPaymentIntentsByInvoice(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load payment intents", http.StatusInternalServerError)
		return
	}
	matches, err := listPaymentMatchesByInvoice(r.Context(), db, invoiceIDsFromSummaries(invoices))
	if err != nil {
		http.Error(w, "cannot load reconciliation matches", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, paymentReconciliationResponse{
		Providers:    providers,
		Schools:      options.Schools,
		SchoolYears:  options.SchoolYears,
		Classes:      options.Classes,
		Summary:      summarizePaymentReconciliation(invoices, transactions),
		Invoices:     invoices,
		Transactions: transactions,
		Intents:      intents,
		Matches:      matches,
	})
}

func handlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	providerCode := strings.TrimPrefix(r.URL.Path, "/api/v1/payments/webhooks/")
	providerCode = headerKey(strings.Trim(providerCode, "/"))
	if providerCode == "" || strings.Contains(providerCode, "/") {
		http.Error(w, "payment provider is required", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPaymentWebhookBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read webhook body", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	provider, err := loadPaymentProviderByCode(r.Context(), db, providerCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawPayload, err := decodeJSONObject(body)
	if err != nil {
		http.Error(w, "webhook body must be a json object", http.StatusBadRequest)
		return
	}

	event, err := insertProviderEvent(r.Context(), db, provider.ID, body, rawPayload, webhookHeaderSnapshot(r.Header))
	if err != nil {
		http.Error(w, "cannot store provider event", http.StatusInternalServerError)
		return
	}
	if event.Duplicate && (event.Status == "processed" || event.Status == "duplicate") {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "duplicate": true, "providerEventId": event.ID})
		return
	}

	normalized, err := normalizeProviderWebhook(provider.Code, rawPayload)
	if err != nil {
		_ = updateProviderEventStatus(r.Context(), db, event.ID, "invalid", "", nil, err.Error())
		_ = recordOperationLog(r.Context(), db, paymentWebhookOperationLog(r, provider.Code, event.ID, "normalize_failed", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	normalized.ProviderCode = provider.Code
	normalized.ProviderEventID = event.ID

	if provider.Code == paymentProviderPayOS {
		if err := verifyPayOSWebhookSignature(rawPayload); err != nil {
			_ = updateProviderEventStatus(r.Context(), db, event.ID, "invalid", normalized.ProviderTransactionID, normalized.RawPayload, err.Error())
			_ = recordOperationLog(r.Context(), db, paymentWebhookOperationLog(r, provider.Code, event.ID, "invalid_signature", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	inserted, match, statusRefresh, err := recordAndReconcilePaymentTransaction(r.Context(), db, provider, normalized)
	if err != nil {
		_ = updateProviderEventStatus(r.Context(), db, event.ID, "invalid", normalized.ProviderTransactionID, normalized.RawPayload, err.Error())
		_ = recordOperationLog(r.Context(), db, paymentWebhookOperationLog(r, provider.Code, event.ID, "reconcile_failed", err))
		http.Error(w, "cannot reconcile payment transaction", http.StatusInternalServerError)
		return
	}

	status := "processed"
	if inserted.Duplicate {
		status = "duplicate"
	}
	_ = updateProviderEventStatus(r.Context(), db, event.ID, status, normalized.ProviderTransactionID, normalized.RawPayload, "")
	if statusRefresh.BecamePaid() {
		sendAutomaticPaidConfirmationBestEffort(r.Context(), db, statusRefresh.InvoiceID, "payment.webhook")
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":         true,
		"duplicate":       inserted.Duplicate,
		"transaction":     inserted.Summary,
		"match":           match,
		"providerEventId": event.ID,
	})
}

func paymentWebhookOperationLog(r *http.Request, providerCode string, eventID string, status string, err error) operationLogInput {
	message := ""
	if err != nil {
		message = err.Error()
	}
	return operationLogInput{
		RequestID:  strings.TrimSpace(r.Header.Get(requestIDHeader)),
		Source:     "webhook",
		Level:      "error",
		Operation:  "payment.webhook",
		Status:     status,
		Message:    message,
		EntityType: "provider_event",
		EntityID:   eventID,
		Metadata: map[string]any{
			"provider": providerCode,
		},
	}
}

func handleManualCashReceipt(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input manualCashReceiptRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.InvoiceID = strings.TrimSpace(input.InvoiceID)
	input.CollectorUserID = strings.TrimSpace(input.CollectorUserID)
	input.CollectorName = strings.TrimSpace(input.CollectorName)
	input.ReceiptReference = normalizeReceiptReference(input.ReceiptReference)
	input.Reason = strings.TrimSpace(input.Reason)
	input.Note = strings.TrimSpace(input.Note)
	if input.Reason == "" && input.Note != "" {
		input.Reason = input.Note
	}
	if input.InvoiceID == "" {
		http.Error(w, "invoiceId is required", http.StatusBadRequest)
		return
	}
	if input.Amount <= 0 {
		http.Error(w, "amount must be greater than zero", http.StatusBadRequest)
		return
	}
	if input.CollectorName == "" {
		http.Error(w, "collectorName is required", http.StatusBadRequest)
		return
	}
	if input.ReceiptReference == "" {
		http.Error(w, "receiptReference is required", http.StatusBadRequest)
		return
	}
	if input.Reason == "" {
		http.Error(w, "reason is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	auditCtx := auditContextFromRequest(r)
	auditCtx.ActorUserID = firstNonEmpty(input.CollectorUserID, auditCtx.ActorUserID)
	auditCtx.ActorName = firstNonEmpty(input.CollectorName, auditCtx.ActorName)
	response, err := recordManualCashReceipt(r.Context(), db, input, tenantID, auditCtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func listPaymentProviders(ctx context.Context, db *sql.DB) ([]paymentProvider, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id::text, code, display_name, provider_type, status
FROM payment_providers
ORDER BY CASE code WHEN 'manual_vietqr' THEN 1 WHEN 'sepay' THEN 2 WHEN 'payos' THEN 3 ELSE 9 END, code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := []paymentProvider{}
	for rows.Next() {
		var provider paymentProvider
		if err := rows.Scan(&provider.ID, &provider.Code, &provider.DisplayName, &provider.ProviderType, &provider.Status); err != nil {
			return nil, err
		}
		provider.Configured = isPaymentProviderConfigured(provider.Code)
		if provider.Code == paymentProviderSePay || provider.Code == paymentProviderPayOS {
			provider.WebhookPath = "/api/v1/payments/webhooks/" + provider.Code
		}
		providers = append(providers, provider)
	}
	return providers, rows.Err()
}

func loadPaymentProviderByCode(ctx context.Context, db *sql.DB, code string) (paymentProvider, error) {
	var provider paymentProvider
	err := db.QueryRowContext(ctx, `
SELECT id::text, code, display_name, provider_type, status
FROM payment_providers
WHERE code = $1
	AND status <> 'inactive'`, code).Scan(
		&provider.ID,
		&provider.Code,
		&provider.DisplayName,
		&provider.ProviderType,
		&provider.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return provider, fmt.Errorf("payment provider %q is not available", code)
	}
	if err != nil {
		return provider, err
	}
	provider.Configured = isPaymentProviderConfigured(provider.Code)
	if provider.Code == paymentProviderSePay || provider.Code == paymentProviderPayOS {
		provider.WebhookPath = "/api/v1/payments/webhooks/" + provider.Code
	}
	return provider, nil
}

func createPaymentIntentForInvoice(ctx context.Context, db *sql.DB, provider paymentProvider, invoice invoiceDocument) (paymentIntentResponse, error) {
	if existing, ok, err := loadActivePaymentIntent(ctx, db, invoice.ID, provider.ID); err != nil {
		return paymentIntentResponse{}, err
	} else if ok {
		return paymentIntentResponse{Intent: existing, QR: qrItemForPaymentIntent(invoice, existing)}, nil
	}

	switch provider.Code {
	case paymentProviderManualVietQR, paymentProviderSePay:
		qr := buildQRItem(paymentRowFromInvoice(invoice), 360)
		if len(qr.Errors) > 0 {
			return paymentIntentResponse{}, errors.New(strings.Join(qr.Errors, "; "))
		}
		intent := paymentIntentSummary{
			InvoiceID:         invoice.ID,
			InvoiceCode:       invoice.InvoiceCode,
			ProviderCode:      provider.Code,
			IntentCode:        stablePaymentIntentCode(invoice.ID, provider.Code),
			Status:            paymentIntentStatusActive,
			Amount:            invoice.TotalAmount,
			Currency:          paymentCurrencyVND,
			ProviderReference: invoice.InvoiceCode,
			QRPayload:         qr.VietQR,
		}
		saved, err := insertPaymentIntent(ctx, db, provider.ID, intent, map[string]any{}, map[string]any{}, map[string]any{"adapter": provider.Code})
		if err != nil {
			return paymentIntentResponse{}, err
		}
		return paymentIntentResponse{Intent: saved, QR: &qr}, nil
	case paymentProviderPayOS:
		result, err := createPayOSPaymentLink(ctx, invoice)
		if err != nil {
			return paymentIntentResponse{}, err
		}
		intent := paymentIntentSummary{
			InvoiceID:         invoice.ID,
			InvoiceCode:       invoice.InvoiceCode,
			ProviderCode:      provider.Code,
			IntentCode:        stablePaymentIntentCode(invoice.ID, provider.Code),
			Status:            paymentIntentStatusActive,
			Amount:            invoice.TotalAmount,
			Currency:          paymentCurrencyVND,
			ProviderReference: strconv.FormatInt(result.OrderCode, 10),
			PaymentURL:        result.CheckoutURL,
			QRPayload:         firstNonEmpty(result.QRCode, buildQRItem(paymentRowFromInvoice(invoice), 360).VietQR),
		}
		metadata := map[string]any{"paymentLinkId": result.PaymentLinkID, "payosStatus": result.Status}
		saved, err := insertPaymentIntent(ctx, db, provider.ID, intent, result.RequestPayload, result.ResponsePayload, metadata)
		if err != nil {
			return paymentIntentResponse{}, err
		}
		return paymentIntentResponse{Intent: saved}, nil
	default:
		return paymentIntentResponse{}, fmt.Errorf("unsupported payment provider %q", provider.Code)
	}
}

func loadActivePaymentIntent(ctx context.Context, db *sql.DB, invoiceID string, providerID string) (paymentIntentSummary, bool, error) {
	var intent paymentIntentSummary
	err := db.QueryRowContext(ctx, `
SELECT pi.id::text, pi.invoice_id::text, i.invoice_code, pp.code, pi.intent_code, pi.status,
	pi.amount, pi.currency, pi.provider_reference, pi.payment_url, pi.qr_payload, pi.created_at
FROM payment_intents pi
JOIN payment_providers pp ON pp.id = pi.provider_id
JOIN invoices i ON i.id = pi.invoice_id
WHERE pi.invoice_id = $1::uuid
	AND pi.provider_id = $2::uuid
	AND pi.status NOT IN ('cancelled', 'expired', 'failed')
ORDER BY pi.created_at DESC
LIMIT 1`, invoiceID, providerID).Scan(
		&intent.ID,
		&intent.InvoiceID,
		&intent.InvoiceCode,
		&intent.ProviderCode,
		&intent.IntentCode,
		&intent.Status,
		&intent.Amount,
		&intent.Currency,
		&intent.ProviderReference,
		&intent.PaymentURL,
		&intent.QRPayload,
		&intent.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return intent, false, nil
	}
	return intent, err == nil, err
}

func insertPaymentIntent(ctx context.Context, exec masterDataExecutor, providerID string, intent paymentIntentSummary, request map[string]any, response map[string]any, metadata map[string]any) (paymentIntentSummary, error) {
	requestJSON, err := jsonObjectString(request)
	if err != nil {
		return intent, err
	}
	responseJSON, err := jsonObjectString(response)
	if err != nil {
		return intent, err
	}
	metadataJSON, err := jsonObjectString(metadata)
	if err != nil {
		return intent, err
	}
	err = exec.QueryRowContext(ctx, `
INSERT INTO payment_intents (
	invoice_id, provider_id, intent_code, status, amount, currency,
	provider_reference, payment_url, qr_payload, provider_request, provider_response, metadata
)
VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11::jsonb, $12::jsonb)
RETURNING id::text, created_at`,
		intent.InvoiceID,
		providerID,
		intent.IntentCode,
		intent.Status,
		intent.Amount,
		firstNonEmpty(intent.Currency, paymentCurrencyVND),
		intent.ProviderReference,
		intent.PaymentURL,
		intent.QRPayload,
		requestJSON,
		responseJSON,
		metadataJSON,
	).Scan(&intent.ID, &intent.CreatedAt)
	if err != nil {
		return intent, err
	}
	return intent, nil
}

func listLatestPaymentIntentsByInvoice(ctx context.Context, db *sql.DB, tenantID string) (map[string]paymentIntentSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT DISTINCT ON (pi.invoice_id)
	pi.id::text, pi.invoice_id::text, i.invoice_code, pp.code, pi.intent_code, pi.status,
	pi.amount, pi.currency, pi.provider_reference, pi.payment_url, pi.qr_payload, pi.created_at
FROM payment_intents pi
JOIN payment_providers pp ON pp.id = pi.provider_id
JOIN invoices i ON i.id = pi.invoice_id
JOIN school_years sy ON sy.id = i.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE pi.status NOT IN ('cancelled', 'expired', 'failed')
	AND sc.tenant_id = $1::uuid
ORDER BY pi.invoice_id, pi.created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	intents := map[string]paymentIntentSummary{}
	for rows.Next() {
		var intent paymentIntentSummary
		if err := rows.Scan(
			&intent.ID,
			&intent.InvoiceID,
			&intent.InvoiceCode,
			&intent.ProviderCode,
			&intent.IntentCode,
			&intent.Status,
			&intent.Amount,
			&intent.Currency,
			&intent.ProviderReference,
			&intent.PaymentURL,
			&intent.QRPayload,
			&intent.CreatedAt,
		); err != nil {
			return nil, err
		}
		intents[intent.InvoiceID] = intent
	}
	return intents, rows.Err()
}

func listPaymentTransactions(ctx context.Context, db *sql.DB, filters paymentTransactionListFilters) ([]paymentTransactionSummary, error) {
	conditions := []string{"1 = 1"}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.Provider != "" {
		conditions = append(conditions, "pp.code = "+addArg(filters.Provider))
	}
	if filters.Status != "" {
		conditions = append(conditions, "pt.status = "+addArg(filters.Status))
	}
	if filters.TenantID != "" {
		conditions = append(conditions, `EXISTS (
			SELECT 1
			FROM invoices tenant_invoice
			JOIN school_years tenant_sy ON tenant_sy.id = tenant_invoice.school_year_id
			JOIN schools tenant_school ON tenant_school.id = tenant_sy.school_id
			WHERE tenant_invoice.id = pt.invoice_id
				AND tenant_school.tenant_id = `+addArg(filters.TenantID)+`::uuid
		)`)
	}
	limit := filters.Limit
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	limitArg := addArg(limit)

	query := `
SELECT pt.id::text,
	pp.code,
	pt.provider_transaction_id,
	COALESCE(pt.invoice_id::text, ''),
	COALESCE(i.invoice_code, ''),
	COALESCE(i.student_code, ''),
	COALESCE(i.student_name, ''),
	pt.direction,
	pt.amount,
	pt.currency,
	pt.transaction_time,
	pt.account_number,
	pt.bank_name,
	pt.description,
	pt.reference_code,
	pt.status,
	COALESCE(match_detail.match_type, ''),
	COALESCE(match_detail.status, ''),
	COALESCE(match_detail.score, 0),
	COALESCE(match_detail.amount_applied, 0),
	COALESCE(match_detail.reason, '')
FROM payment_transactions pt
JOIN payment_providers pp ON pp.id = pt.provider_id
LEFT JOIN invoices i ON i.id = pt.invoice_id
LEFT JOIN LATERAL (
	SELECT rm.match_type, rm.status, rm.score, rm.amount_applied, rm.reason
	FROM reconciliation_matches rm
	WHERE rm.transaction_id = pt.id
		AND rm.status <> 'reversed'
	ORDER BY CASE WHEN rm.status = 'matched' THEN 0 ELSE 1 END, rm.created_at DESC
	LIMIT 1
) match_detail ON true
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY pt.transaction_time DESC, pt.created_at DESC
LIMIT ` + limitArg

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []paymentTransactionSummary{}
	for rows.Next() {
		var item paymentTransactionSummary
		if err := rows.Scan(
			&item.ID,
			&item.ProviderCode,
			&item.ProviderTransactionID,
			&item.InvoiceID,
			&item.InvoiceCode,
			&item.StudentCode,
			&item.StudentName,
			&item.Direction,
			&item.Amount,
			&item.Currency,
			&item.TransactionTime,
			&item.AccountNumber,
			&item.BankName,
			&item.Description,
			&item.ReferenceCode,
			&item.Status,
			&item.MatchType,
			&item.MatchStatus,
			&item.MatchScore,
			&item.AmountApplied,
			&item.MatchReason,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, item)
	}
	return transactions, rows.Err()
}

func invoiceIDsFromSummaries(invoices []invoiceSummary) []string {
	ids := make([]string, 0, len(invoices))
	for _, invoice := range invoices {
		id := strings.TrimSpace(invoice.ID)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func listPaymentMatchesByInvoice(ctx context.Context, db *sql.DB, invoiceIDs []string) (map[string][]paymentMatchSummary, error) {
	matches := map[string][]paymentMatchSummary{}
	invoiceIDs = normalizeStringList(invoiceIDs)
	if len(invoiceIDs) == 0 {
		return matches, nil
	}
	args := make([]any, 0, len(invoiceIDs))
	placeholders := make([]string, 0, len(invoiceIDs))
	for _, id := range invoiceIDs {
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d::uuid", len(args)))
	}
	rows, err := db.QueryContext(ctx, `
SELECT rm.id::text,
	rm.invoice_id::text,
	i.invoice_code,
	rm.transaction_id::text,
	pp.code,
	pt.provider_transaction_id,
	rm.match_type,
	rm.status,
	rm.score,
	rm.amount_applied,
	rm.reason,
	rm.created_at
FROM reconciliation_matches rm
JOIN invoices i ON i.id = rm.invoice_id
JOIN payment_transactions pt ON pt.id = rm.transaction_id
JOIN payment_providers pp ON pp.id = pt.provider_id
WHERE rm.status <> 'reversed'
	AND rm.invoice_id IN (`+strings.Join(placeholders, ", ")+`)
ORDER BY rm.created_at DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item paymentMatchSummary
		if err := rows.Scan(
			&item.ID,
			&item.InvoiceID,
			&item.InvoiceCode,
			&item.TransactionID,
			&item.ProviderCode,
			&item.ProviderTransactionID,
			&item.MatchType,
			&item.Status,
			&item.Score,
			&item.AmountApplied,
			&item.Reason,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		matches[item.InvoiceID] = append(matches[item.InvoiceID], item)
	}
	return matches, rows.Err()
}

func insertProviderEvent(ctx context.Context, db *sql.DB, providerID string, body []byte, rawPayload map[string]any, headers map[string]any) (providerEventRecord, error) {
	payloadHash := sha256Hex(body)
	rawJSON, err := jsonObjectString(rawPayload)
	if err != nil {
		return providerEventRecord{}, err
	}
	headerJSON, err := jsonObjectString(headers)
	if err != nil {
		return providerEventRecord{}, err
	}
	var event providerEventRecord
	err = db.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO provider_events (provider_id, payload_hash, raw_payload, headers)
	VALUES ($1::uuid, $2, $3::jsonb, $4::jsonb)
	ON CONFLICT (provider_id, payload_hash) DO NOTHING
	RETURNING id::text, false AS duplicate, status
)
SELECT id, duplicate, status FROM inserted
UNION ALL
SELECT id::text, true AS duplicate, status
FROM provider_events
WHERE provider_id = $1::uuid
	AND payload_hash = $2
LIMIT 1`,
		providerID,
		payloadHash,
		rawJSON,
		headerJSON,
	).Scan(&event.ID, &event.Duplicate, &event.Status)
	return event, err
}

func updateProviderEventStatus(ctx context.Context, exec masterDataExecutor, eventID string, status string, providerEventID string, normalized map[string]any, errorMessage string) error {
	normalizedJSON, err := jsonObjectString(normalized)
	if err != nil {
		return err
	}
	_, err = exec.ExecContext(ctx, `
UPDATE provider_events
SET status = $2,
	provider_event_id = CASE WHEN $3 = '' THEN provider_event_id ELSE $3 END,
	normalized_payload = $4::jsonb,
	error_message = $5,
	processed_at = now()
WHERE id = $1::uuid`,
		eventID,
		status,
		providerEventID,
		normalizedJSON,
		errorMessage,
	)
	return err
}

func normalizeProviderWebhook(providerCode string, raw map[string]any) (normalizedPaymentTransaction, error) {
	switch providerCode {
	case paymentProviderSePay:
		return normalizeSePayWebhook(raw)
	case paymentProviderPayOS:
		return normalizePayOSWebhook(raw)
	default:
		return normalizedPaymentTransaction{}, fmt.Errorf("provider %q does not accept webhooks", providerCode)
	}
}

func normalizeSePayWebhook(raw map[string]any) (normalizedPaymentTransaction, error) {
	direction := headerKey(jsonStringValue(raw["transferType"]))
	if direction == "" {
		direction = paymentDirectionIn
	}
	if direction != paymentDirectionIn {
		return normalizedPaymentTransaction{}, fmt.Errorf("only inbound SePay transactions are reconciled")
	}
	amount := jsonIntValue(raw["transferAmount"])
	if amount <= 0 {
		return normalizedPaymentTransaction{}, fmt.Errorf("transferAmount must be greater than zero")
	}
	transactionID := jsonStringValue(raw["id"])
	if transactionID == "" {
		transactionID = firstNonEmpty(jsonStringValue(raw["referenceCode"]), jsonStringValue(raw["code"]))
	}
	if transactionID == "" {
		return normalizedPaymentTransaction{}, fmt.Errorf("SePay transaction id is required")
	}
	transactionTime := parseProviderTime(firstNonEmpty(jsonStringValue(raw["transactionDate"]), jsonStringValue(raw["transactionDateTime"])), time.Now())
	content := firstNonEmpty(jsonStringValue(raw["content"]), jsonStringValue(raw["description"]), jsonStringValue(raw["code"]))
	description := strings.TrimSpace(strings.Join(nonEmptyStrings(
		jsonStringValue(raw["code"]),
		content,
		jsonStringValue(raw["description"]),
	), " "))

	return normalizedPaymentTransaction{
		ProviderTransactionID: transactionID,
		Direction:             paymentDirectionIn,
		Amount:                amount,
		Currency:              paymentCurrencyVND,
		TransactionTime:       transactionTime,
		AccountNumber:         cleanAccount(firstNonEmpty(jsonStringValue(raw["subAccount"]), jsonStringValue(raw["accountNumber"]))),
		BankName:              strings.TrimSpace(jsonStringValue(raw["gateway"])),
		Description:           description,
		ReferenceCode:         firstNonEmpty(jsonStringValue(raw["referenceCode"]), jsonStringValue(raw["code"])),
		Status:                paymentTransactionStatusUnmatched,
		RawPayload:            raw,
	}, nil
}

func normalizePayOSWebhook(raw map[string]any) (normalizedPaymentTransaction, error) {
	data, ok := raw["data"].(map[string]any)
	if !ok {
		return normalizedPaymentTransaction{}, fmt.Errorf("payOS webhook data is required")
	}
	if success, ok := raw["success"].(bool); ok && !success {
		return normalizedPaymentTransaction{}, fmt.Errorf("payOS webhook is not successful")
	}
	amount := jsonIntValue(data["amount"])
	if amount <= 0 {
		return normalizedPaymentTransaction{}, fmt.Errorf("payOS amount must be greater than zero")
	}
	transactionID := firstNonEmpty(jsonStringValue(data["reference"]), jsonStringValue(data["paymentLinkId"]), jsonStringValue(data["orderCode"]))
	if transactionID == "" {
		return normalizedPaymentTransaction{}, fmt.Errorf("payOS transaction reference is required")
	}
	return normalizedPaymentTransaction{
		ProviderTransactionID: transactionID,
		Direction:             paymentDirectionIn,
		Amount:                amount,
		Currency:              firstNonEmpty(jsonStringValue(data["currency"]), paymentCurrencyVND),
		TransactionTime:       parseProviderTime(jsonStringValue(data["transactionDateTime"]), time.Now()),
		AccountNumber:         cleanAccount(jsonStringValue(data["accountNumber"])),
		Description:           strings.TrimSpace(jsonStringValue(data["description"])),
		ReferenceCode:         firstNonEmpty(jsonStringValue(data["reference"]), jsonStringValue(data["orderCode"]), jsonStringValue(data["paymentLinkId"])),
		Status:                paymentTransactionStatusUnmatched,
		RawPayload:            raw,
	}, nil
}

func recordAndReconcilePaymentTransaction(ctx context.Context, db *sql.DB, provider paymentProvider, normalized normalizedPaymentTransaction) (insertedPaymentTransaction, *paymentMatchCandidate, invoicePaymentStatusRefresh, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return insertedPaymentTransaction{}, nil, invoicePaymentStatusRefresh{}, err
	}
	defer tx.Rollback()

	inserted, err := insertPaymentTransaction(ctx, tx, provider.ID, normalized)
	if err != nil {
		return insertedPaymentTransaction{}, nil, invoicePaymentStatusRefresh{}, err
	}
	match, statusRefresh, err := reconcilePaymentTransaction(ctx, tx, inserted.Summary)
	if err != nil {
		return insertedPaymentTransaction{}, nil, invoicePaymentStatusRefresh{}, err
	}
	if err := tx.Commit(); err != nil {
		return insertedPaymentTransaction{}, nil, invoicePaymentStatusRefresh{}, err
	}
	return inserted, match, statusRefresh, nil
}

func insertPaymentTransaction(ctx context.Context, exec masterDataExecutor, providerID string, normalized normalizedPaymentTransaction) (insertedPaymentTransaction, error) {
	rawJSON, err := jsonObjectString(normalized.RawPayload)
	if err != nil {
		return insertedPaymentTransaction{}, err
	}
	var item insertedPaymentTransaction
	err = exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO payment_transactions (
		provider_id, provider_event_id, payment_intent_id, invoice_id, provider_transaction_id,
		direction, amount, currency, transaction_time, account_number, account_name,
		bank_bin, bank_name, description, reference_code, status, raw_payload
	)
	VALUES (
		$1::uuid, nullif($2, '')::uuid, nullif($3, '')::uuid, nullif($4, '')::uuid, $5,
		$6, $7, $8, $9, $10, $11,
		$12, $13, $14, $15, $16, $17::jsonb
	)
	ON CONFLICT (provider_id, provider_transaction_id) WHERE provider_transaction_id <> '' DO NOTHING
	RETURNING id::text, false AS duplicate
)
SELECT id, duplicate FROM inserted
UNION ALL
SELECT id::text, true AS duplicate
FROM payment_transactions
WHERE provider_id = $1::uuid
	AND provider_transaction_id = $5
LIMIT 1`,
		providerID,
		normalized.ProviderEventID,
		normalized.PaymentIntentID,
		normalized.InvoiceID,
		normalized.ProviderTransactionID,
		firstNonEmpty(normalized.Direction, paymentDirectionIn),
		normalized.Amount,
		firstNonEmpty(normalized.Currency, paymentCurrencyVND),
		normalized.TransactionTime,
		normalized.AccountNumber,
		normalized.AccountName,
		normalized.BankBIN,
		normalized.BankName,
		normalized.Description,
		normalized.ReferenceCode,
		firstNonEmpty(normalized.Status, paymentTransactionStatusUnmatched),
		rawJSON,
	).Scan(&item.ID, &item.Duplicate)
	if err != nil {
		return item, err
	}
	summary, err := loadPaymentTransactionSummary(ctx, exec, item.ID)
	if err != nil {
		return item, err
	}
	item.Summary = summary
	return item, nil
}

func loadPaymentTransactionSummary(ctx context.Context, exec masterDataExecutor, transactionID string) (paymentTransactionSummary, error) {
	var item paymentTransactionSummary
	err := exec.QueryRowContext(ctx, `
SELECT pt.id::text,
	pp.code,
	pt.provider_transaction_id,
	COALESCE(pt.invoice_id::text, ''),
	COALESCE(i.invoice_code, ''),
	COALESCE(i.student_code, ''),
	COALESCE(i.student_name, ''),
	pt.direction,
	pt.amount,
	pt.currency,
	pt.transaction_time,
	pt.account_number,
	pt.bank_name,
	pt.description,
	pt.reference_code,
	pt.status
FROM payment_transactions pt
JOIN payment_providers pp ON pp.id = pt.provider_id
LEFT JOIN invoices i ON i.id = pt.invoice_id
WHERE pt.id = $1::uuid`, transactionID).Scan(
		&item.ID,
		&item.ProviderCode,
		&item.ProviderTransactionID,
		&item.InvoiceID,
		&item.InvoiceCode,
		&item.StudentCode,
		&item.StudentName,
		&item.Direction,
		&item.Amount,
		&item.Currency,
		&item.TransactionTime,
		&item.AccountNumber,
		&item.BankName,
		&item.Description,
		&item.ReferenceCode,
		&item.Status,
	)
	return item, err
}

func reconcilePaymentTransaction(ctx context.Context, exec masterDataExecutor, transaction paymentTransactionSummary) (*paymentMatchCandidate, invoicePaymentStatusRefresh, error) {
	if transaction.Status == paymentTransactionStatusMatched && transaction.InvoiceID != "" {
		return nil, invoicePaymentStatusRefresh{}, nil
	}
	candidates, err := loadPaymentInvoiceCandidates(ctx, exec, transaction)
	if err != nil {
		return nil, invoicePaymentStatusRefresh{}, err
	}
	match, ok := matchPaymentTransactionToInvoices(transaction, candidates)
	if !ok {
		_, err := exec.ExecContext(ctx, `
UPDATE payment_transactions
SET status = $2
WHERE id = $1::uuid
	AND status <> 'matched'`,
			transaction.ID,
			paymentTransactionStatusManualReview,
		)
		return nil, invoicePaymentStatusRefresh{}, err
	}

	if _, err := exec.ExecContext(ctx, `
UPDATE payment_transactions
SET invoice_id = $2::uuid,
	status = $3
WHERE id = $1::uuid`,
		transaction.ID,
		match.Invoice.ID,
		paymentTransactionStatusMatched,
	); err != nil {
		return nil, invoicePaymentStatusRefresh{}, err
	}
	metadata, err := jsonObjectString(map[string]any{
		"provider":              transaction.ProviderCode,
		"providerTransactionId": transaction.ProviderTransactionID,
		"referenceCode":         transaction.ReferenceCode,
		"description":           transaction.Description,
	})
	if err != nil {
		return nil, invoicePaymentStatusRefresh{}, err
	}
	if _, err := exec.ExecContext(ctx, `
INSERT INTO reconciliation_matches (invoice_id, transaction_id, match_type, status, score, amount_applied, reason, metadata)
VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, $7, $8::jsonb)
ON CONFLICT (transaction_id, invoice_id) WHERE status <> 'reversed' DO NOTHING`,
		match.Invoice.ID,
		transaction.ID,
		match.MatchType,
		reconciliationStatusMatched,
		match.Score,
		match.AmountApplied,
		match.Reason,
		metadata,
	); err != nil {
		return nil, invoicePaymentStatusRefresh{}, err
	}
	statusRefresh, err := refreshInvoicePaymentStatus(ctx, exec, match.Invoice.ID, "reconciled payment transaction")
	if err != nil {
		return nil, invoicePaymentStatusRefresh{}, err
	}
	return &match, statusRefresh, nil
}

func loadPaymentInvoiceCandidates(ctx context.Context, exec masterDataExecutor, transaction paymentTransactionSummary) ([]paymentInvoiceCandidate, error) {
	account := cleanAccount(transaction.AccountNumber)
	rows, err := exec.QueryContext(ctx, `
SELECT i.id::text,
	i.invoice_code,
	i.qr_bill_number,
	i.collection_bank_account,
	i.status,
	i.total_amount,
	i.paid_amount,
	COALESCE(string_agg(pi.provider_reference, ' '), '')
FROM invoices i
LEFT JOIN payment_intents pi ON pi.invoice_id = i.id
	AND pi.status NOT IN ('cancelled', 'expired', 'failed')
WHERE i.status <> 'void'
	AND ($1 = '' OR i.collection_bank_account = $1 OR pi.provider_reference = $2 OR pi.provider_reference = $3)
GROUP BY i.id
ORDER BY i.issued_at DESC
LIMIT 1000`,
		account,
		transaction.ReferenceCode,
		transaction.ProviderTransactionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candidates := []paymentInvoiceCandidate{}
	for rows.Next() {
		var candidate paymentInvoiceCandidate
		var providerRefs string
		if err := rows.Scan(
			&candidate.ID,
			&candidate.InvoiceCode,
			&candidate.QRBillNumber,
			&candidate.CollectionBankAccount,
			&candidate.Status,
			&candidate.TotalAmount,
			&candidate.PaidAmount,
			&providerRefs,
		); err != nil {
			return nil, err
		}
		candidate.ProviderReferences = strings.Fields(providerRefs)
		candidates = append(candidates, candidate)
	}
	return candidates, rows.Err()
}

func matchPaymentTransactionToInvoices(transaction paymentTransactionSummary, candidates []paymentInvoiceCandidate) (paymentMatchCandidate, bool) {
	if transaction.Amount <= 0 || transaction.Direction != paymentDirectionIn {
		return paymentMatchCandidate{}, false
	}
	referenceBlob := normalizedPaymentReference(strings.Join(nonEmptyStrings(
		transaction.ProviderTransactionID,
		transaction.ReferenceCode,
		transaction.Description,
	), " "))
	account := cleanAccount(transaction.AccountNumber)
	best := paymentMatchCandidate{}
	for _, candidate := range candidates {
		score := 0
		reasons := []string{}
		if account != "" && account == cleanAccount(candidate.CollectionBankAccount) {
			score += 30
			reasons = append(reasons, "account")
		}
		if containsPaymentReference(referenceBlob, candidate.InvoiceCode) || containsPaymentReference(referenceBlob, candidate.QRBillNumber) {
			score += 50
			reasons = append(reasons, "invoice_code")
		}
		for _, ref := range candidate.ProviderReferences {
			if containsPaymentReference(referenceBlob, ref) {
				score += 50
				reasons = append(reasons, "provider_reference")
				break
			}
		}
		outstanding := candidate.TotalAmount - candidate.PaidAmount
		if outstanding < 0 {
			outstanding = 0
		}
		switch {
		case transaction.Amount == outstanding && outstanding > 0:
			score += 20
			reasons = append(reasons, "exact_outstanding_amount")
		case transaction.Amount == candidate.TotalAmount && candidate.PaidAmount == 0:
			score += 20
			reasons = append(reasons, "exact_invoice_amount")
		case transaction.Amount > 0 && transaction.Amount < outstanding:
			score += 10
			reasons = append(reasons, "partial_amount")
		case outstanding > 0 && transaction.Amount > outstanding:
			score += 10
			reasons = append(reasons, "overpayment_amount")
		}
		if score > best.Score {
			matchType := "auto"
			if containsString(reasons, "provider_reference") {
				matchType = "provider_reference"
			}
			best = paymentMatchCandidate{
				Invoice:       candidate,
				Score:         score,
				MatchType:     matchType,
				AmountApplied: transaction.Amount,
				Reason:        strings.Join(reasons, "+"),
			}
		}
	}
	if best.Score < 70 {
		return paymentMatchCandidate{}, false
	}
	return best, true
}

func refreshInvoicePaymentStatus(ctx context.Context, exec masterDataExecutor, invoiceID string, reason string) (invoicePaymentStatusRefresh, error) {
	var totalAmount int
	var oldStatus string
	err := exec.QueryRowContext(ctx, `
SELECT total_amount, status
FROM invoices
WHERE id = $1::uuid
FOR UPDATE`, invoiceID).Scan(&totalAmount, &oldStatus)
	if err != nil {
		return invoicePaymentStatusRefresh{}, err
	}
	var paidAmount int
	if err := exec.QueryRowContext(ctx, `
SELECT COALESCE(SUM(amount_applied), 0)
FROM reconciliation_matches
WHERE invoice_id = $1::uuid
	AND status <> 'reversed'`, invoiceID).Scan(&paidAmount); err != nil {
		return invoicePaymentStatusRefresh{}, err
	}
	newStatus := deriveInvoiceStatus(totalAmount, paidAmount)
	if _, err := exec.ExecContext(ctx, `
UPDATE invoices
SET paid_amount = $2,
	status = $3
WHERE id = $1::uuid`,
		invoiceID,
		paidAmount,
		newStatus,
	); err != nil {
		return invoicePaymentStatusRefresh{}, err
	}
	refresh := invoicePaymentStatusRefresh{
		InvoiceID:  invoiceID,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		PaidAmount: paidAmount,
	}
	if newStatus != oldStatus {
		return refresh, insertInvoiceStatusHistory(ctx, exec, invoiceID, oldStatus, newStatus, reason)
	}
	return refresh, nil
}

func recordManualCashReceipt(ctx context.Context, db *sql.DB, input manualCashReceiptRequest, tenantID string, auditCtx requestAuditContext) (manualCashReceiptResponse, error) {
	provider, err := loadPaymentProviderByCode(ctx, db, paymentProviderManualVietQR)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	invoice, err := loadInvoiceDocument(ctx, db, input.InvoiceID, tenantID)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	paidAt, err := parseManualReceiptTime(input.PaidAt)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	defer tx.Rollback()

	normalized := normalizedPaymentTransaction{
		ProviderTransactionID: "cash:" + input.ReceiptReference,
		InvoiceID:             input.InvoiceID,
		Direction:             paymentDirectionIn,
		Amount:                input.Amount,
		Currency:              paymentCurrencyVND,
		TransactionTime:       paidAt,
		AccountNumber:         invoice.CollectionBankAccount,
		BankBIN:               invoice.CollectionBankBIN,
		BankName:              "Cash",
		Description:           "Cash receipt " + input.ReceiptReference + " " + invoice.InvoiceCode,
		ReferenceCode:         input.ReceiptReference,
		Status:                paymentTransactionStatusMatched,
		RawPayload: map[string]any{
			"source":           "manual_cash_receipt",
			"invoiceId":        input.InvoiceID,
			"receiptReference": input.ReceiptReference,
			"collectorName":    input.CollectorName,
			"reason":           input.Reason,
			"note":             input.Note,
		},
	}
	inserted, err := insertPaymentTransaction(ctx, tx, provider.ID, normalized)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	if inserted.Duplicate {
		return manualCashReceiptResponse{}, fmt.Errorf("receiptReference already exists")
	}

	var receiptID string
	err = tx.QueryRowContext(ctx, `
INSERT INTO manual_cash_receipts (
	invoice_id, payment_transaction_id, collector_user_id, collector_name,
	amount, currency, paid_at, receipt_reference, reason, note, created_by_user_id
)
VALUES ($1::uuid, $2::uuid, nullif($3, '')::uuid, $4, $5, $6, $7, $8, $9, $10, nullif($3, '')::uuid)
RETURNING id::text`,
		input.InvoiceID,
		inserted.ID,
		input.CollectorUserID,
		input.CollectorName,
		input.Amount,
		paymentCurrencyVND,
		paidAt,
		input.ReceiptReference,
		input.Reason,
		input.Note,
	).Scan(&receiptID)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	metadata, _ := jsonObjectString(map[string]any{
		"collectorName":    input.CollectorName,
		"receiptReference": input.ReceiptReference,
		"reason":           input.Reason,
	})
	if _, err := tx.ExecContext(ctx, `
INSERT INTO reconciliation_matches (invoice_id, transaction_id, match_type, status, score, amount_applied, reason, metadata)
VALUES ($1::uuid, $2::uuid, 'cash', 'matched', 100, $3, $4, $5::jsonb)`,
		input.InvoiceID,
		inserted.ID,
		input.Amount,
		input.Reason,
		metadata,
	); err != nil {
		return manualCashReceiptResponse{}, err
	}
	statusRefresh, err := refreshInvoicePaymentStatus(ctx, tx, input.InvoiceID, input.Reason)
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	if err := insertAuditLog(ctx, tx, auditLogInput{
		Context:    auditCtx,
		Action:     "manual_cash_receipt.created",
		EntityType: "manual_cash_receipt",
		EntityID:   receiptID,
		Reason:     input.Reason,
		Metadata: map[string]any{
			"invoiceId":            input.InvoiceID,
			"invoiceCode":          invoice.InvoiceCode,
			"paymentTransactionId": inserted.ID,
			"collectorName":        input.CollectorName,
			"amount":               input.Amount,
			"receiptReference":     input.ReceiptReference,
			"paidAt":               paidAt.Format(time.RFC3339),
		},
	}); err != nil {
		return manualCashReceiptResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return manualCashReceiptResponse{}, err
	}
	if statusRefresh.BecamePaid() {
		sendAutomaticPaidConfirmationBestEffort(ctx, db, statusRefresh.InvoiceID, "manual_cash_receipt")
	}

	invoices, err := listInvoiceSummaries(ctx, db, invoiceListFilters{TenantID: tenantID})
	if err != nil {
		return manualCashReceiptResponse{}, err
	}
	var updated invoiceSummary
	for _, item := range invoices {
		if item.ID == input.InvoiceID {
			updated = item
			break
		}
	}
	inserted.Summary.Status = paymentTransactionStatusMatched
	inserted.Summary.InvoiceID = input.InvoiceID
	inserted.Summary.InvoiceCode = invoice.InvoiceCode
	inserted.Summary.StudentCode = invoice.StudentCode
	inserted.Summary.StudentName = invoice.StudentName
	return manualCashReceiptResponse{ReceiptID: receiptID, Transaction: inserted.Summary, Invoice: updated}, nil
}

func summarizePaymentReconciliation(invoices []invoiceSummary, transactions []paymentTransactionSummary) paymentReconciliationSummary {
	summary := paymentReconciliationSummary{InvoiceCount: len(invoices)}
	for _, invoice := range invoices {
		summary.TotalReceivable += invoice.TotalAmount
		summary.TotalCollected += invoice.PaidAmount
		if invoice.TotalAmount > invoice.PaidAmount {
			summary.OutstandingAmount += invoice.TotalAmount - invoice.PaidAmount
		}
		switch invoice.Status {
		case invoiceStatusUnpaid:
			summary.UnpaidCount++
		case invoiceStatusPartial:
			summary.PartialCount++
		case invoiceStatusPaid:
			summary.PaidCount++
		case invoiceStatusOverpaid:
			summary.OverpaidCount++
		case invoiceStatusManualReview:
			summary.ManualReviewCount++
		}
	}
	if summary.TotalReceivable > 0 {
		summary.CollectionRate = float64(summary.TotalCollected) / float64(summary.TotalReceivable)
	}
	for _, transaction := range transactions {
		switch transaction.Status {
		case paymentTransactionStatusMatched:
			summary.MatchedCount++
		case paymentTransactionStatusUnmatched:
			summary.UnmatchedCount++
		case paymentTransactionStatusManualReview:
			summary.ManualReviewCount++
		}
	}
	return summary
}

func createPayOSPaymentLink(ctx context.Context, invoice invoiceDocument) (payOSCreatePaymentResult, error) {
	cfg := loadPayOSConfig()
	if !cfg.configuredForCreate() {
		return payOSCreatePaymentResult{}, fmt.Errorf("payOS is not configured; set ABC_PAYOS_CLIENT_ID, ABC_PAYOS_API_KEY, ABC_PAYOS_CHECKSUM_KEY, ABC_PAYOS_RETURN_URL, and ABC_PAYOS_CANCEL_URL")
	}
	orderCode := payOSOrderCode(invoice.InvoiceCode)
	description := cleanANS(invoice.InvoiceCode, 9)
	if description == "" {
		description = "ABCSUN"
	}
	payload := map[string]any{
		"orderCode":   orderCode,
		"amount":      invoice.TotalAmount,
		"description": description,
		"buyerName":   invoice.StudentName,
		"items": []map[string]any{{
			"name":     invoice.InvoiceCode,
			"quantity": 1,
			"price":    invoice.TotalAmount,
		}},
		"cancelUrl": cfg.CancelURL,
		"returnUrl": cfg.ReturnURL,
	}
	payload["signature"] = payOSSignature(map[string]any{
		"amount":      invoice.TotalAmount,
		"cancelUrl":   cfg.CancelURL,
		"description": description,
		"orderCode":   orderCode,
		"returnUrl":   cfg.ReturnURL,
	}, cfg.ChecksumKey)

	body, err := json.Marshal(payload)
	if err != nil {
		return payOSCreatePaymentResult{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.APIBaseURL, "/")+"/v2/payment-requests", bytes.NewReader(body))
	if err != nil {
		return payOSCreatePaymentResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-client-id", cfg.ClientID)
	req.Header.Set("x-api-key", cfg.APIKey)

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return payOSCreatePaymentResult{}, err
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return payOSCreatePaymentResult{}, err
	}
	responsePayload := map[string]any{}
	_ = json.Unmarshal(responseBody, &responsePayload)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return payOSCreatePaymentResult{}, fmt.Errorf("payOS create payment link failed with status %d", res.StatusCode)
	}
	if code := jsonStringValue(responsePayload["code"]); code != "" && code != "00" {
		return payOSCreatePaymentResult{}, fmt.Errorf("payOS create payment link failed: %s", firstNonEmpty(jsonStringValue(responsePayload["desc"]), code))
	}
	data, _ := responsePayload["data"].(map[string]any)
	return payOSCreatePaymentResult{
		OrderCode:       orderCode,
		PaymentLinkID:   jsonStringValue(data["paymentLinkId"]),
		CheckoutURL:     jsonStringValue(data["checkoutUrl"]),
		QRCode:          jsonStringValue(data["qrCode"]),
		Status:          jsonStringValue(data["status"]),
		RequestPayload:  payload,
		ResponsePayload: responsePayload,
	}, nil
}

func verifyPayOSWebhookSignature(raw map[string]any) error {
	cfg := loadPayOSConfig()
	if strings.TrimSpace(cfg.ChecksumKey) == "" {
		return nil
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("payOS webhook data is required")
	}
	signature := jsonStringValue(raw["signature"])
	if signature == "" {
		return fmt.Errorf("payOS webhook signature is required")
	}
	if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(payOSSignature(data, cfg.ChecksumKey))) {
		return fmt.Errorf("payOS webhook signature is invalid")
	}
	return nil
}

func loadPayOSConfig() payOSConfig {
	return payOSConfig{
		ClientID:    strings.TrimSpace(os.Getenv("ABC_PAYOS_CLIENT_ID")),
		APIKey:      strings.TrimSpace(os.Getenv("ABC_PAYOS_API_KEY")),
		ChecksumKey: strings.TrimSpace(os.Getenv("ABC_PAYOS_CHECKSUM_KEY")),
		ReturnURL:   strings.TrimSpace(os.Getenv("ABC_PAYOS_RETURN_URL")),
		CancelURL:   strings.TrimSpace(os.Getenv("ABC_PAYOS_CANCEL_URL")),
		APIBaseURL:  firstNonEmpty(os.Getenv("ABC_PAYOS_API_BASE_URL"), "https://api-merchant.payos.vn"),
	}
}

func (cfg payOSConfig) configuredForCreate() bool {
	return cfg.ClientID != "" && cfg.APIKey != "" && cfg.ChecksumKey != "" && cfg.ReturnURL != "" && cfg.CancelURL != ""
}

func isPaymentProviderConfigured(code string) bool {
	switch code {
	case paymentProviderManualVietQR:
		return true
	case paymentProviderSePay:
		return true
	case paymentProviderPayOS:
		return loadPayOSConfig().configuredForCreate()
	default:
		return false
	}
}

func payOSOrderCode(invoiceCode string) int64 {
	sum := crc64.Checksum([]byte(invoiceCode), crc64.MakeTable(crc64.ISO))
	return int64(sum%899999999) + 100000000
}

func payOSSignature(data map[string]any, checksumKey string) string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+payOSSignatureValue(data[key]))
	}
	mac := hmac.New(sha256.New, []byte(checksumKey))
	_, _ = mac.Write([]byte(strings.Join(parts, "&")))
	return hex.EncodeToString(mac.Sum(nil))
}

func payOSSignatureValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		if typed == "undefined" || typed == "null" {
			return ""
		}
		return typed
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		return strconv.FormatBool(typed)
	default:
		data, _ := json.Marshal(typed)
		return string(data)
	}
}

func qrItemForPaymentIntent(invoice invoiceDocument, intent paymentIntentSummary) *qrItem {
	if intent.QRPayload == "" {
		return nil
	}
	item := buildQRItem(paymentRowFromInvoice(invoice), 360)
	if len(item.Errors) > 0 {
		return nil
	}
	item.VietQR = intent.QRPayload
	return &item
}

func stablePaymentIntentCode(invoiceID string, providerCode string) string {
	source := invoiceID + "|" + providerCode
	sum := crc64.Checksum([]byte(source), crc64.MakeTable(crc64.ISO))
	return "PINT" + strings.ToUpper(strconv.FormatUint(sum, 36))
}

func decodeJSONObject(body []byte) (map[string]any, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if payload == nil {
		return nil, fmt.Errorf("json object is required")
	}
	return payload, nil
}

func jsonObjectString(value map[string]any) (string, error) {
	if value == nil {
		value = map[string]any{}
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func jsonStringValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return typed.String()
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func jsonIntValue(value any) int {
	switch typed := value.(type) {
	case nil:
		return 0
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		parsed, _ := typed.Int64()
		return int(parsed)
	case string:
		return parseAmount(typed)
	default:
		return parseAmount(fmt.Sprint(typed))
	}
}

func parseProviderTime(value string, fallback time.Time) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback.UTC()
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05Z07:00"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return fallback.UTC()
}

func parseManualReceiptTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().UTC(), nil
	}
	return parseInvoiceDateTime(value, time.Now())
}

func normalizedPaymentReference(value string) string {
	return nonAlnumPattern.ReplaceAllString(strings.ToUpper(value), "")
}

func containsPaymentReference(blob string, reference string) bool {
	reference = normalizedPaymentReference(reference)
	return reference != "" && strings.Contains(blob, reference)
}

func normalizeReceiptReference(value string) string {
	return cleanANS(strings.ToUpper(strings.TrimSpace(value)), 48)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func webhookHeaderSnapshot(header http.Header) map[string]any {
	keep := []string{
		"Content-Type",
		"User-Agent",
		"X-Forwarded-For",
		"X-Real-Ip",
		"Svix-Id",
		"Svix-Timestamp",
		"X-Sepay-Signature",
		"X-Signature",
	}
	out := map[string]any{}
	for _, key := range keep {
		if value := strings.TrimSpace(header.Get(key)); value != "" {
			out[key] = value
		}
	}
	return out
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func nonEmptyStrings(values ...string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, strings.TrimSpace(value))
		}
	}
	return out
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
