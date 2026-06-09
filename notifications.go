package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	notificationCampaignFirstNotice              = "first_notice"
	notificationCampaignReminder                 = "reminder"
	notificationCampaignPaymentConfirmation      = "payment_confirmation"
	notificationAutoPaidConfirmationCampaignCode = "auto_paid_confirmation"

	notificationStatusDraft   = "draft"
	notificationStatusDryRun  = "dry_run"
	notificationStatusSent    = "sent"
	notificationStatusPartial = "partial"
)

type notificationTemplate struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Version       int       `json:"version"`
	Name          string    `json:"name"`
	Subject       string    `json:"subject"`
	EmailTemplate string    `json:"emailTemplate"`
	Status        string    `json:"status"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type notificationCampaignInput struct {
	TenantID      string   `json:"-"`
	ID            string   `json:"id,omitempty"`
	CampaignID    string   `json:"campaignId,omitempty"`
	Name          string   `json:"name"`
	CampaignType  string   `json:"campaignType"`
	TemplateID    string   `json:"templateId"`
	SchoolYearID  string   `json:"schoolYearId,omitempty"`
	ClassID       string   `json:"classId,omitempty"`
	Grade         string   `json:"grade,omitempty"`
	PeriodCode    string   `json:"periodCode,omitempty"`
	InvoiceStatus string   `json:"invoiceStatus,omitempty"`
	DueOnOrBefore string   `json:"dueOnOrBefore,omitempty"`
	DryRun        bool     `json:"dryRun,omitempty"`
	ForceResend   bool     `json:"forceResend,omitempty"`
	ConfirmSend   bool     `json:"confirmSend,omitempty"`
	RecipientIDs  []string `json:"recipientIds,omitempty"`
}

type notificationCampaignSummary struct {
	ID              string    `json:"id,omitempty"`
	Code            string    `json:"code,omitempty"`
	Name            string    `json:"name"`
	CampaignType    string    `json:"campaignType"`
	Status          string    `json:"status"`
	Template        string    `json:"template"`
	TemplateID      string    `json:"templateId"`
	TemplateVersion int       `json:"templateVersion"`
	SchoolYearID    string    `json:"schoolYearId,omitempty"`
	SchoolYearCode  string    `json:"schoolYearCode,omitempty"`
	ClassID         string    `json:"classId,omitempty"`
	ClassName       string    `json:"className,omitempty"`
	Grade           string    `json:"grade,omitempty"`
	PeriodCode      string    `json:"periodCode,omitempty"`
	InvoiceStatus   string    `json:"invoiceStatus,omitempty"`
	DueOnOrBefore   string    `json:"dueOnOrBefore,omitempty"`
	RecipientCount  int       `json:"recipientCount"`
	SentCount       int       `json:"sentCount"`
	ErrorCount      int       `json:"errorCount"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}

type notificationRecipientCandidate struct {
	ID                string `json:"id,omitempty"`
	CampaignID        string `json:"campaignId,omitempty"`
	InvoiceID         string `json:"invoiceId"`
	ParentID          string `json:"parentId,omitempty"`
	RecipientName     string `json:"recipientName"`
	RecipientEmail    string `json:"recipientEmail"`
	InvoiceCode       string `json:"invoiceCode"`
	StudentCode       string `json:"studentCode"`
	StudentName       string `json:"studentName"`
	ClassName         string `json:"className"`
	Grade             string `json:"grade"`
	SchoolYearID      string `json:"schoolYearId"`
	SchoolYearCode    string `json:"schoolYearCode"`
	PeriodCode        string `json:"periodCode"`
	DueDate           string `json:"dueDate,omitempty"`
	InvoiceStatus     string `json:"invoiceStatus"`
	Amount            int    `json:"amount"`
	PaidAmount        int    `json:"paidAmount"`
	OutstandingAmount int    `json:"outstandingAmount"`
	Status            string `json:"status,omitempty"`
	LastError         string `json:"lastError,omitempty"`
	AlreadySent       bool   `json:"alreadySent,omitempty"`
	QRReady           bool   `json:"qrReady"`
	SendCount         int    `json:"sendCount"`
	LastSentAt        string `json:"lastSentAt,omitempty"`
	LastLogStatus     string `json:"lastLogStatus,omitempty"`
	RetryEligible     bool   `json:"retryEligible,omitempty"`
}

type notificationRecipientSummary struct {
	InvoiceCount       int `json:"invoiceCount"`
	RecipientCount     int `json:"recipientCount"`
	TotalAmount        int `json:"totalAmount"`
	UnpaidAmount       int `json:"unpaidAmount"`
	AlreadySent        int `json:"alreadySent"`
	QRMissingCount     int `json:"qrMissingCount"`
	ErrorCount         int `json:"errorCount"`
	RetryEligibleCount int `json:"retryEligibleCount"`
}

type notificationIssue struct {
	Type    string `json:"type"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type notificationPreviewResponse struct {
	Template   notificationTemplate             `json:"template"`
	Campaign   notificationCampaignSummary      `json:"campaign"`
	Summary    notificationRecipientSummary     `json:"summary"`
	Recipients []notificationRecipientCandidate `json:"recipients"`
	Issues     []notificationIssue              `json:"issues,omitempty"`
}

type notificationOptionsResponse struct {
	Templates   []notificationTemplate        `json:"templates"`
	Campaigns   []notificationCampaignSummary `json:"campaigns"`
	SchoolYears []masterDataSchoolYearOption  `json:"schoolYears"`
	Classes     []masterDataClassOption       `json:"classes"`
}

type notificationLogSummary struct {
	ID                string    `json:"id"`
	CampaignID        string    `json:"campaignId"`
	CampaignName      string    `json:"campaignName"`
	TemplateCode      string    `json:"templateCode"`
	TemplateVersion   int       `json:"templateVersion"`
	InvoiceID         string    `json:"invoiceId"`
	InvoiceCode       string    `json:"invoiceCode"`
	RecipientEmail    string    `json:"recipientEmail"`
	Provider          string    `json:"provider,omitempty"`
	Status            string    `json:"status"`
	ProviderMessageID string    `json:"providerMessageId,omitempty"`
	Error             string    `json:"error,omitempty"`
	DryRun            bool      `json:"dryRun"`
	SentAt            time.Time `json:"sentAt"`
}

type notificationSendResponse struct {
	Campaign notificationCampaignSummary  `json:"campaign"`
	Summary  notificationRecipientSummary `json:"summary"`
	Results  []emailSendResult            `json:"results"`
	Logs     []notificationLogSummary     `json:"logs"`
}

type notificationEmailPreviewInput struct {
	notificationCampaignInput
	RecipientID    string `json:"recipientId,omitempty"`
	InvoiceID      string `json:"invoiceId,omitempty"`
	RecipientEmail string `json:"recipientEmail,omitempty"`
}

type notificationEmailPreviewResponse struct {
	Subject   string                         `json:"subject"`
	HTML      string                         `json:"html"`
	Text      string                         `json:"text"`
	To        string                         `json:"to"`
	Template  notificationTemplate           `json:"template"`
	Recipient notificationRecipientCandidate `json:"recipient"`
	QRReady   bool                           `json:"qrReady"`
}

type paidConfirmationSendInput struct {
	InvoiceID   string `json:"invoiceId"`
	ForceResend bool   `json:"forceResend,omitempty"`
	ConfirmSend bool   `json:"confirmSend,omitempty"`
}

type paidConfirmationSendResponse struct {
	InvoiceID string            `json:"invoiceId"`
	Results   []emailSendResult `json:"results"`
}

type paidConfirmationSendOptions struct {
	ForceResend  bool
	CampaignCode string
	CampaignName string
	Trigger      string
	Operation    string
}

func handleNotificationOptions(w http.ResponseWriter, r *http.Request) {
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

	templates, err := listNotificationTemplates(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load notification templates", http.StatusInternalServerError)
		return
	}
	campaigns, err := listNotificationCampaigns(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load notification campaigns", http.StatusInternalServerError)
		return
	}
	options, err := listMasterDataOptions(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load master data options", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, notificationOptionsResponse{
		Templates:   templates,
		Campaigns:   campaigns,
		SchoolYears: options.SchoolYears,
		Classes:     options.Classes,
	})
}

func handleNotificationTemplates(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	templates, err := listNotificationTemplates(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load notification templates", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"templates": templates})
}

func handleNotificationCampaigns(w http.ResponseWriter, r *http.Request) {
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

	campaigns, err := listNotificationCampaigns(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load notification campaigns", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"campaigns": campaigns})
}

func handleNotificationCampaignPreview(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	input, ok := decodeNotificationCampaignInput(w, r)
	if !ok {
		return
	}
	input.TenantID = tenantID
	preview, err := buildNotificationPreview(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, preview)
}

func handleNotificationCampaignEmailPreview(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	req, ok := decodeNotificationEmailPreviewInput(w, r)
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	input := req.notificationCampaignInput
	input.TenantID = tenantID
	if input.CampaignID != "" {
		stored, err := loadNotificationCampaignInput(r.Context(), db, input.CampaignID, tenantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stored.ForceResend = input.ForceResend
		input = stored
	}
	template, err := loadNotificationTemplateForInput(r.Context(), db, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	recipients, err := notificationPreviewRecipientsForEmailPreview(r.Context(), db, input, template)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	recipient, ok := selectNotificationPreviewRecipient(recipients, req.RecipientID, req.InvoiceID, req.RecipientEmail)
	if !ok {
		http.Error(w, "notification recipient not found", http.StatusBadRequest)
		return
	}
	invoice, err := loadInvoiceDocument(r.Context(), db, recipient.InvoiceID, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := loadEmailConfigForTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	emailCfg := notificationEmailConfig(cfg, template, invoice)
	item := buildQRItem(notificationPaymentRow(invoice, recipient), 512)
	if len(item.Errors) > 0 {
		http.Error(w, strings.Join(item.Errors, "; "), http.StatusBadRequest)
		return
	}
	email, err := renderPaymentEmail(emailCfg, item, template.EmailTemplate, appBaseURL(r, emailCfg), "data")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, notificationEmailPreviewResponse{
		Subject:   email.Subject,
		HTML:      email.HTML,
		Text:      email.Text,
		To:        item.Email,
		Template:  template,
		Recipient: recipient,
		QRReady:   recipient.QRReady,
	})
}

func handleNotificationCampaignSave(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	input, ok := decodeNotificationCampaignInput(w, r)
	if !ok {
		return
	}
	input.TenantID = tenantID
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	preview, err := buildNotificationPreviewFromDB(r.Context(), db, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(preview.Issues) > 0 {
		writeJSON(w, http.StatusBadRequest, preview)
		return
	}
	campaign, err := saveNotificationCampaign(r.Context(), db, input, preview.Template, preview.Recipients)
	if err != nil {
		http.Error(w, "cannot save notification campaign: "+err.Error(), http.StatusInternalServerError)
		return
	}
	campaigns, _ := listNotificationCampaigns(r.Context(), db, tenantID)
	writeJSON(w, http.StatusOK, map[string]any{
		"campaign":   campaign,
		"summary":    summarizeNotificationRecipients(preview.Recipients),
		"recipients": preview.Recipients,
		"campaigns":  campaigns,
	})
}

func handleNotificationCampaignSend(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	input, ok := decodeNotificationCampaignInput(w, r)
	if !ok {
		return
	}
	input.TenantID = tenantID
	recipientIDs := input.RecipientIDs
	if !input.DryRun && !input.ConfirmSend {
		http.Error(w, "confirmSend is required for real notification sends", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if input.CampaignID != "" {
		stored, err := loadNotificationCampaignInput(r.Context(), db, input.CampaignID, tenantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stored.DryRun = input.DryRun
		stored.ForceResend = input.ForceResend
		stored.ConfirmSend = input.ConfirmSend
		stored.RecipientIDs = recipientIDs
		input = stored
	}

	preview, err := buildNotificationPreviewFromDB(r.Context(), db, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(preview.Issues) > 0 {
		writeJSON(w, http.StatusBadRequest, preview)
		return
	}
	if input.CampaignID == "" {
		campaign, err := saveNotificationCampaign(r.Context(), db, input, preview.Template, preview.Recipients)
		if err != nil {
			http.Error(w, "cannot save notification campaign: "+err.Error(), http.StatusInternalServerError)
			return
		}
		input.CampaignID = campaign.ID
	}
	recipients, err := loadNotificationRecipients(r.Context(), db, input.CampaignID, tenantID)
	if err != nil {
		http.Error(w, "cannot load notification recipients", http.StatusInternalServerError)
		return
	}
	recipients = filterNotificationRecipientsForSend(recipients, input.RecipientIDs)
	if len(recipients) == 0 {
		if len(input.RecipientIDs) > 0 {
			http.Error(w, "selected notification recipients were not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "campaign has no recipients", http.StatusBadRequest)
		return
	}

	cfg, err := loadEmailConfigForTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !input.DryRun {
		if err := validateEmailConfigForSend(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	sentLimit := 0
	if !input.DryRun {
		quota, err := emailSendQuotaStatusForTenant(r.Context(), tenantID, time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if quota.Remaining <= 0 {
			http.Error(w, fmt.Sprintf("email daily limit reached (%d/%d in 24h)", quota.Sent, quota.Limit), http.StatusTooManyRequests)
			return
		}
		sentLimit = quota.Remaining
	}
	if !input.DryRun {
		plannedSends, err := countNotificationRecipientsPlannedForSend(r.Context(), db, input.CampaignID, preview.Template, recipients, input.ForceResend, sentLimit)
		if err != nil {
			http.Error(w, "cannot inspect notification send quota", http.StatusInternalServerError)
			return
		}
		if err := enforceTenantUsageLimit(r.Context(), db, tenantID, subscriptionMetricMonthlyNotifications, plannedSends, time.Now()); err != nil {
			var usageErr *tenantUsageLimitError
			if errors.As(err, &usageErr) {
				http.Error(w, usageErr.Error(), http.StatusForbidden)
				return
			}
			http.Error(w, "cannot inspect subscription usage", http.StatusInternalServerError)
			return
		}
	}

	response, err := sendNotificationCampaign(r.Context(), db, cfg, preview.Template, input, recipients, appBaseURL(r, cfg), sentLimit)
	if err != nil {
		http.Error(w, "cannot send notification campaign: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func handleNotificationLogs(w http.ResponseWriter, r *http.Request) {
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

	logs, err := listNotificationLogs(r.Context(), db, strings.TrimSpace(r.URL.Query().Get("campaignId")), tenantID, parsePositiveInt(r.URL.Query().Get("limit"), 100))
	if err != nil {
		http.Error(w, "cannot load notification logs", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"logs": logs})
}

func handleNotificationPaidConfirmationSend(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input paidConfirmationSendInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.InvoiceID = strings.TrimSpace(input.InvoiceID)
	if input.InvoiceID == "" {
		http.Error(w, "invoiceId is required", http.StatusBadRequest)
		return
	}
	if !input.ConfirmSend {
		http.Error(w, "confirmSend is required for paid confirmation sends", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	results, err := sendPaidConfirmationForInvoice(r.Context(), db, input.InvoiceID, tenantID, paidConfirmationSendOptions{
		ForceResend: input.ForceResend,
		Trigger:     "manual",
		Operation:   "notification.paid_confirmation.manual",
	})
	if err != nil {
		var usageErr *tenantUsageLimitError
		if errors.As(err, &usageErr) {
			http.Error(w, usageErr.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, paidConfirmationSendResponse{
		InvoiceID: input.InvoiceID,
		Results:   results,
	})
}

func decodeNotificationCampaignInput(w http.ResponseWriter, r *http.Request) (notificationCampaignInput, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	var input notificationCampaignInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return input, false
	}
	return normalizeNotificationCampaignInput(input), true
}

func decodeNotificationEmailPreviewInput(w http.ResponseWriter, r *http.Request) (notificationEmailPreviewInput, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	var input notificationEmailPreviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return input, false
	}
	input.notificationCampaignInput = normalizeNotificationCampaignInput(input.notificationCampaignInput)
	input.RecipientID = strings.TrimSpace(input.RecipientID)
	input.InvoiceID = strings.TrimSpace(input.InvoiceID)
	input.RecipientEmail = strings.ToLower(strings.TrimSpace(input.RecipientEmail))
	return input, true
}

func buildNotificationPreview(ctx context.Context, input notificationCampaignInput) (notificationPreviewResponse, error) {
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return notificationPreviewResponse{}, err
	}
	defer db.Close()
	return buildNotificationPreviewFromDB(ctx, db, input)
}

func buildNotificationPreviewFromDB(ctx context.Context, db *sql.DB, input notificationCampaignInput) (notificationPreviewResponse, error) {
	input = normalizeNotificationCampaignInput(input)
	issues := validateNotificationCampaignInput(input)
	template, err := loadNotificationTemplateForInput(ctx, db, input)
	if err != nil {
		return notificationPreviewResponse{}, err
	}
	recipients := []notificationRecipientCandidate{}
	if len(issues) == 0 {
		recipients, err = listNotificationRecipientCandidates(ctx, db, input, template)
		if err != nil {
			return notificationPreviewResponse{}, err
		}
	}
	return notificationPreviewResponse{
		Template:   template,
		Campaign:   notificationCampaignSummaryFromInput(input, template, notificationStatusDraft, recipients),
		Summary:    summarizeNotificationRecipients(recipients),
		Recipients: recipients,
		Issues:     issues,
	}, nil
}

func normalizeNotificationCampaignInput(input notificationCampaignInput) notificationCampaignInput {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.ID = strings.TrimSpace(input.ID)
	input.CampaignID = strings.TrimSpace(firstNonEmpty(input.CampaignID, input.ID))
	input.Name = strings.TrimSpace(input.Name)
	input.CampaignType = headerKey(input.CampaignType)
	if input.CampaignType == "" {
		input.CampaignType = notificationCampaignFirstNotice
	}
	if input.CampaignType != notificationCampaignFirstNotice &&
		input.CampaignType != notificationCampaignReminder &&
		input.CampaignType != notificationCampaignPaymentConfirmation {
		input.CampaignType = notificationCampaignFirstNotice
	}
	input.TemplateID = strings.TrimSpace(input.TemplateID)
	input.SchoolYearID = strings.TrimSpace(input.SchoolYearID)
	input.ClassID = strings.TrimSpace(input.ClassID)
	input.Grade = normalizeGrade(input.Grade)
	input.PeriodCode = strings.TrimSpace(input.PeriodCode)
	input.InvoiceStatus = headerKey(input.InvoiceStatus)
	input.DueOnOrBefore = strings.TrimSpace(input.DueOnOrBefore)
	input.RecipientIDs = normalizeStringList(input.RecipientIDs)
	if input.Name == "" {
		label := "Thông báo thanh toán"
		if input.CampaignType == notificationCampaignReminder {
			label = "Nhắc thanh toán"
		} else if input.CampaignType == notificationCampaignPaymentConfirmation {
			label = "Xác nhận đã thanh toán"
		}
		input.Name = strings.TrimSpace(label + " " + firstNonEmpty(input.PeriodCode, time.Now().Format("2006-01-02")))
	}
	if input.CampaignType == notificationCampaignPaymentConfirmation && input.InvoiceStatus == "" {
		input.InvoiceStatus = invoiceStatusPaid
	}
	return input
}

func notificationPreviewRecipientsForEmailPreview(ctx context.Context, db *sql.DB, input notificationCampaignInput, template notificationTemplate) ([]notificationRecipientCandidate, error) {
	if input.CampaignID != "" {
		return loadNotificationRecipients(ctx, db, input.CampaignID, input.TenantID)
	}
	issues := validateNotificationCampaignInput(input)
	if len(issues) > 0 {
		return nil, errors.New(issues[0].Message)
	}
	return listNotificationRecipientCandidates(ctx, db, input, template)
}

func selectNotificationPreviewRecipient(recipients []notificationRecipientCandidate, recipientID string, invoiceID string, recipientEmail string) (notificationRecipientCandidate, bool) {
	recipientEmail = strings.ToLower(strings.TrimSpace(recipientEmail))
	for _, recipient := range recipients {
		if recipientID != "" && recipient.ID == recipientID {
			return recipient, true
		}
		if invoiceID != "" && recipient.InvoiceID == invoiceID {
			if recipientEmail == "" || strings.EqualFold(recipient.RecipientEmail, recipientEmail) {
				return recipient, true
			}
		}
	}
	if recipientID == "" && invoiceID == "" && recipientEmail == "" && len(recipients) > 0 {
		return recipients[0], true
	}
	return notificationRecipientCandidate{}, false
}

func normalizeStringList(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func filterNotificationRecipientsForSend(recipients []notificationRecipientCandidate, recipientIDs []string) []notificationRecipientCandidate {
	recipientIDs = normalizeStringList(recipientIDs)
	if len(recipientIDs) == 0 {
		return recipients
	}
	allowed := make(map[string]bool, len(recipientIDs))
	for _, id := range recipientIDs {
		allowed[id] = true
	}
	filtered := make([]notificationRecipientCandidate, 0, len(recipientIDs))
	for _, recipient := range recipients {
		if allowed[recipient.ID] {
			filtered = append(filtered, recipient)
		}
	}
	return filtered
}

func validateNotificationCampaignInput(input notificationCampaignInput) []notificationIssue {
	issues := []notificationIssue{}
	if input.CampaignType == notificationCampaignReminder && input.InvoiceStatus != "" {
		if input.InvoiceStatus != invoiceStatusUnpaid && input.InvoiceStatus != invoiceStatusPartial {
			issues = append(issues, notificationIssue{
				Type:    "invalid_reminder_status",
				Field:   "invoiceStatus",
				Message: "reminder campaigns can only target unpaid or partial invoices",
			})
		}
	}
	if input.CampaignType == notificationCampaignPaymentConfirmation && input.InvoiceStatus != "" && input.InvoiceStatus != invoiceStatusPaid {
		issues = append(issues, notificationIssue{
			Type:    "invalid_payment_confirmation_status",
			Field:   "invoiceStatus",
			Message: "payment confirmation campaigns can only target paid invoices",
		})
	}
	if input.DueOnOrBefore != "" {
		if _, err := parseInvoiceDate(input.DueOnOrBefore); err != nil {
			issues = append(issues, notificationIssue{
				Type:    "invalid_due_date",
				Field:   "dueOnOrBefore",
				Message: "dueOnOrBefore must be YYYY-MM-DD",
			})
		}
	}
	return issues
}

func listNotificationTemplates(ctx context.Context, db *sql.DB) ([]notificationTemplate, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id::text, code, version, name, subject, email_template, status, updated_at
FROM notification_templates
ORDER BY code, version DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := []notificationTemplate{}
	for rows.Next() {
		var item notificationTemplate
		if err := rows.Scan(&item.ID, &item.Code, &item.Version, &item.Name, &item.Subject, &item.EmailTemplate, &item.Status, &item.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, item)
	}
	return templates, rows.Err()
}

func loadNotificationTemplateForInput(ctx context.Context, db *sql.DB, input notificationCampaignInput) (notificationTemplate, error) {
	if input.TemplateID != "" {
		return loadNotificationTemplateByID(ctx, db, input.TemplateID)
	}
	code := input.CampaignType
	if code == "" {
		code = notificationCampaignFirstNotice
	}
	return loadLatestNotificationTemplateByCode(ctx, db, code)
}

func loadNotificationTemplateByID(ctx context.Context, db *sql.DB, templateID string) (notificationTemplate, error) {
	var item notificationTemplate
	err := db.QueryRowContext(ctx, `
SELECT id::text, code, version, name, subject, email_template, status, updated_at
FROM notification_templates
WHERE id = $1::uuid`, templateID).Scan(
		&item.ID,
		&item.Code,
		&item.Version,
		&item.Name,
		&item.Subject,
		&item.EmailTemplate,
		&item.Status,
		&item.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return item, errors.New("notification template not found")
	}
	return item, err
}

func loadLatestNotificationTemplateByCode(ctx context.Context, db *sql.DB, code string) (notificationTemplate, error) {
	var item notificationTemplate
	err := db.QueryRowContext(ctx, `
SELECT id::text, code, version, name, subject, email_template, status, updated_at
FROM notification_templates
WHERE code = $1
	AND status = 'active'
ORDER BY version DESC
LIMIT 1`, code).Scan(
		&item.ID,
		&item.Code,
		&item.Version,
		&item.Name,
		&item.Subject,
		&item.EmailTemplate,
		&item.Status,
		&item.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return item, fmt.Errorf("active notification template %q not found", code)
	}
	return item, err
}

func listNotificationCampaigns(ctx context.Context, db *sql.DB, tenantID string) ([]notificationCampaignSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT c.id::text,
	c.code,
	c.name,
	c.campaign_type,
	c.status,
	c.template_id::text,
	t.code,
	t.version,
	COALESCE(c.school_year_id::text, ''),
	COALESCE(sy.code, ''),
	COALESCE(c.class_id::text, ''),
	COALESCE(cls.name, ''),
	c.grade,
	c.period_code,
	c.invoice_status,
	c.due_on_or_before,
	c.created_at,
	c.updated_at,
	COUNT(r.id)::int,
	COUNT(r.id) FILTER (WHERE r.status = 'sent')::int,
	COUNT(r.id) FILTER (WHERE r.status = 'error')::int
FROM notification_campaigns c
JOIN notification_templates t ON t.id = c.template_id
LEFT JOIN school_years sy ON sy.id = c.school_year_id
LEFT JOIN classes cls ON cls.id = c.class_id
LEFT JOIN notification_recipients r ON r.campaign_id = c.id
WHERE c.status <> 'archived'
	AND c.tenant_id = $1::uuid
GROUP BY c.id, t.id, sy.id, cls.id
ORDER BY c.created_at DESC, c.id DESC
LIMIT 200`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := []notificationCampaignSummary{}
	for rows.Next() {
		var item notificationCampaignSummary
		var dueDate sql.NullTime
		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.CampaignType,
			&item.Status,
			&item.TemplateID,
			&item.Template,
			&item.TemplateVersion,
			&item.SchoolYearID,
			&item.SchoolYearCode,
			&item.ClassID,
			&item.ClassName,
			&item.Grade,
			&item.PeriodCode,
			&item.InvoiceStatus,
			&dueDate,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.RecipientCount,
			&item.SentCount,
			&item.ErrorCount,
		); err != nil {
			return nil, err
		}
		if dueDate.Valid {
			item.DueOnOrBefore = dueDate.Time.Format("2006-01-02")
		}
		campaigns = append(campaigns, item)
	}
	return campaigns, rows.Err()
}

func loadNotificationCampaignInput(ctx context.Context, db *sql.DB, campaignID string, tenantID string) (notificationCampaignInput, error) {
	var input notificationCampaignInput
	var dueDate sql.NullTime
	err := db.QueryRowContext(ctx, `
SELECT id::text,
	name,
	campaign_type,
	template_id::text,
	COALESCE(school_year_id::text, ''),
	COALESCE(class_id::text, ''),
	grade,
	period_code,
	invoice_status,
	due_on_or_before
FROM notification_campaigns
WHERE id = $1::uuid
	AND tenant_id = $2::uuid
	AND status <> 'archived'`, campaignID, tenantID).Scan(
		&input.CampaignID,
		&input.Name,
		&input.CampaignType,
		&input.TemplateID,
		&input.SchoolYearID,
		&input.ClassID,
		&input.Grade,
		&input.PeriodCode,
		&input.InvoiceStatus,
		&dueDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return input, errors.New("notification campaign not found")
	}
	if err != nil {
		return input, err
	}
	if dueDate.Valid {
		input.DueOnOrBefore = dueDate.Time.Format("2006-01-02")
	}
	input.TenantID = tenantID
	return normalizeNotificationCampaignInput(input), nil
}

func listNotificationRecipientCandidates(ctx context.Context, db *sql.DB, input notificationCampaignInput, template notificationTemplate) ([]notificationRecipientCandidate, error) {
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	templateArg := addArg(template.ID)
	tenantArg := addArg(input.TenantID)
	conditions := []string{
		"sc.tenant_id = " + tenantArg + "::uuid",
		"i.status <> 'void'",
		"sp.is_active",
		"sp.receives_billing_email",
		"p.email_active",
		"p.status = 'active'",
		"p.email <> ''",
	}
	if input.SchoolYearID != "" {
		conditions = append(conditions, "i.school_year_id = "+addArg(input.SchoolYearID)+"::uuid")
	}
	if input.ClassID != "" {
		conditions = append(conditions, "i.class_id = "+addArg(input.ClassID)+"::uuid")
	}
	if input.Grade != "" {
		conditions = append(conditions, "i.grade = "+addArg(input.Grade))
	}
	if input.PeriodCode != "" {
		conditions = append(conditions, "i.period_code = "+addArg(input.PeriodCode))
	}
	if input.DueOnOrBefore != "" {
		conditions = append(conditions, "i.due_date IS NOT NULL AND i.due_date <= "+addArg(input.DueOnOrBefore)+"::date")
	}
	switch {
	case input.CampaignType == notificationCampaignReminder && input.InvoiceStatus == "":
		conditions = append(conditions, "i.status IN ('unpaid', 'partial')")
	case input.CampaignType == notificationCampaignPaymentConfirmation && input.InvoiceStatus == "":
		conditions = append(conditions, "i.status = 'paid'")
	case input.InvoiceStatus != "":
		conditions = append(conditions, "i.status = "+addArg(input.InvoiceStatus))
	default:
		conditions = append(conditions, "i.status = 'unpaid'")
	}

	query := `
SELECT i.id::text,
	COALESCE(p.id::text, ''),
	p.full_name,
	p.email,
	i.invoice_code,
	i.student_code,
	i.student_name,
	i.class_name,
	i.grade,
	i.school_year_id::text,
	i.school_year_code,
	i.period_code,
	i.due_date,
	i.status,
	i.total_amount,
	i.paid_amount,
	GREATEST(i.total_amount - i.paid_amount, 0),
	(i.qr_bill_number <> '' AND i.collection_bank_bin <> '' AND i.collection_bank_account <> ''),
	COALESCE(log_counts.send_count, 0),
	log_counts.last_sent_at,
	COALESCE(log_counts.last_status, ''),
	COALESCE(log_counts.last_error, '')
FROM invoices i
JOIN school_years sy ON sy.id = i.school_year_id
JOIN schools sc ON sc.id = sy.school_id
JOIN students s ON s.id = i.student_id
JOIN student_parents sp ON sp.student_id = s.id
JOIN parents p ON p.id = sp.parent_id
LEFT JOIN LATERAL (
	SELECT
		COUNT(*) FILTER (WHERE nl.status = 'sent')::integer AS send_count,
		MAX(nl.sent_at) FILTER (WHERE nl.status = 'sent') AS last_sent_at,
		(ARRAY_AGG(nl.status ORDER BY nl.sent_at DESC, nl.id DESC))[1] AS last_status,
		(ARRAY_AGG(nl.error_message ORDER BY nl.sent_at DESC, nl.id DESC))[1] AS last_error
	FROM notification_logs nl
	WHERE nl.template_id = ` + templateArg + `::uuid
		AND nl.invoice_id = i.id
		AND lower(nl.recipient_email) = lower(p.email)
) log_counts ON true
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY i.period_code DESC, i.class_name, i.student_code, sp.is_primary DESC, p.full_name
LIMIT 2000`

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipients := []notificationRecipientCandidate{}
	seen := map[string]bool{}
	for rows.Next() {
		var item notificationRecipientCandidate
		var dueDate sql.NullTime
		var lastSentAt sql.NullTime
		if err := rows.Scan(
			&item.InvoiceID,
			&item.ParentID,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.InvoiceCode,
			&item.StudentCode,
			&item.StudentName,
			&item.ClassName,
			&item.Grade,
			&item.SchoolYearID,
			&item.SchoolYearCode,
			&item.PeriodCode,
			&dueDate,
			&item.InvoiceStatus,
			&item.Amount,
			&item.PaidAmount,
			&item.OutstandingAmount,
			&item.QRReady,
			&item.SendCount,
			&lastSentAt,
			&item.LastLogStatus,
			&item.LastError,
		); err != nil {
			return nil, err
		}
		key := strings.ToLower(item.InvoiceID + "|" + item.RecipientEmail)
		if seen[key] {
			continue
		}
		seen[key] = true
		if dueDate.Valid {
			item.DueDate = dueDate.Time.Format("2006-01-02")
		}
		if lastSentAt.Valid {
			item.LastSentAt = lastSentAt.Time.UTC().Format(time.RFC3339)
		}
		item = finalizeNotificationRecipientState(item)
		recipients = append(recipients, item)
	}
	return recipients, rows.Err()
}

func saveNotificationCampaign(ctx context.Context, db *sql.DB, input notificationCampaignInput, template notificationTemplate, recipients []notificationRecipientCandidate) (notificationCampaignSummary, error) {
	input = normalizeNotificationCampaignInput(input)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return notificationCampaignSummary{}, err
	}
	defer tx.Rollback()

	filterJSON, err := notificationTargetFilterJSON(input)
	if err != nil {
		return notificationCampaignSummary{}, err
	}
	campaignID := input.CampaignID
	if campaignID == "" {
		err = tx.QueryRowContext(ctx, `
INSERT INTO notification_campaigns (
	tenant_id, code, name, campaign_type, template_id, school_year_id, class_id,
	grade, period_code, invoice_status, due_on_or_before, status, target_filter
)
VALUES ($1::uuid, $2, $3, $4, $5::uuid, $6::uuid, $7::uuid, $8, $9, $10, $11, 'draft', $12::jsonb)
RETURNING id::text`,
			input.TenantID,
			notificationCampaignCode(input, time.Now()),
			input.Name,
			input.CampaignType,
			template.ID,
			nullableString(input.SchoolYearID),
			nullableString(input.ClassID),
			input.Grade,
			input.PeriodCode,
			input.InvoiceStatus,
			nullableDateString(input.DueOnOrBefore),
			filterJSON,
		).Scan(&campaignID)
		if err != nil {
			return notificationCampaignSummary{}, err
		}
	} else {
		err = tx.QueryRowContext(ctx, `
UPDATE notification_campaigns
SET name = $2,
	campaign_type = $3,
	template_id = $4::uuid,
	school_year_id = $5::uuid,
	class_id = $6::uuid,
	grade = $7,
	period_code = $8,
	invoice_status = $9,
	due_on_or_before = $10,
	target_filter = $11::jsonb
WHERE id = $1::uuid
	AND tenant_id = $12::uuid
RETURNING id::text`,
			campaignID,
			input.Name,
			input.CampaignType,
			template.ID,
			nullableString(input.SchoolYearID),
			nullableString(input.ClassID),
			input.Grade,
			input.PeriodCode,
			input.InvoiceStatus,
			nullableDateString(input.DueOnOrBefore),
			filterJSON,
			input.TenantID,
		).Scan(&campaignID)
		if err != nil {
			return notificationCampaignSummary{}, err
		}
		if _, err := tx.ExecContext(ctx, `
DELETE FROM notification_recipients
WHERE campaign_id = $1::uuid
	AND EXISTS (
		SELECT 1
		FROM notification_campaigns c
		WHERE c.id = $1::uuid
			AND c.tenant_id = $2::uuid
	)`, campaignID, input.TenantID); err != nil {
			return notificationCampaignSummary{}, err
		}
	}
	for _, recipient := range recipients {
		if err := insertNotificationRecipient(ctx, tx, campaignID, recipient); err != nil {
			return notificationCampaignSummary{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return notificationCampaignSummary{}, err
	}
	campaign, err := loadNotificationCampaignSummary(ctx, db, campaignID, input.TenantID)
	if err != nil {
		return notificationCampaignSummary{}, err
	}
	return campaign, nil
}

func finalizeNotificationRecipientState(item notificationRecipientCandidate) notificationRecipientCandidate {
	item.RecipientEmail = strings.ToLower(strings.TrimSpace(item.RecipientEmail))
	if item.OutstandingAmount == 0 && item.Amount > item.PaidAmount {
		item.OutstandingAmount = item.Amount - item.PaidAmount
	}
	item.AlreadySent = item.SendCount > 0 || item.AlreadySent
	if item.Status == "" {
		if item.AlreadySent {
			item.Status = "already_sent"
		} else {
			item.Status = "pending"
		}
	}
	item.RetryEligible = item.Status == "error" || item.Status == "skipped" || item.LastLogStatus == "error" || item.LastError != ""
	if item.Status == "sent" || item.Status == "already_sent" {
		item.RetryEligible = false
	}
	return item
}

func insertNotificationRecipient(ctx context.Context, exec masterDataExecutor, campaignID string, recipient notificationRecipientCandidate) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO notification_recipients (
	campaign_id, invoice_id, parent_id, recipient_name, recipient_email,
	invoice_code, student_code, student_name, class_name, period_code, amount, status
)
VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11, 'pending')
ON CONFLICT (campaign_id, invoice_id, lower(recipient_email)) DO UPDATE
SET parent_id = EXCLUDED.parent_id,
	recipient_name = EXCLUDED.recipient_name,
	invoice_code = EXCLUDED.invoice_code,
	student_code = EXCLUDED.student_code,
	student_name = EXCLUDED.student_name,
	class_name = EXCLUDED.class_name,
	period_code = EXCLUDED.period_code,
	amount = EXCLUDED.amount,
	status = 'pending',
	last_error = ''`,
		campaignID,
		recipient.InvoiceID,
		nullableString(recipient.ParentID),
		recipient.RecipientName,
		strings.ToLower(strings.TrimSpace(recipient.RecipientEmail)),
		recipient.InvoiceCode,
		recipient.StudentCode,
		recipient.StudentName,
		recipient.ClassName,
		recipient.PeriodCode,
		recipient.Amount,
	)
	return err
}

func loadNotificationCampaignSummary(ctx context.Context, db *sql.DB, campaignID string, tenantID string) (notificationCampaignSummary, error) {
	campaigns, err := listNotificationCampaigns(ctx, db, tenantID)
	if err != nil {
		return notificationCampaignSummary{}, err
	}
	for _, campaign := range campaigns {
		if campaign.ID == campaignID {
			return campaign, nil
		}
	}
	return notificationCampaignSummary{}, errors.New("notification campaign not found")
}

func loadNotificationRecipients(ctx context.Context, db *sql.DB, campaignID string, tenantID string) ([]notificationRecipientCandidate, error) {
	rows, err := db.QueryContext(ctx, `
SELECT r.id::text,
	r.campaign_id::text,
	r.invoice_id::text,
	COALESCE(parent_id::text, ''),
	r.recipient_name,
	r.recipient_email,
	r.invoice_code,
	r.student_code,
	r.student_name,
	r.class_name,
	i.grade,
	i.school_year_id::text,
	i.school_year_code,
	r.period_code,
	i.due_date,
	i.status,
	r.amount,
	i.paid_amount,
	GREATEST(r.amount - i.paid_amount, 0),
	r.status,
	r.last_error,
	(i.qr_bill_number <> '' AND i.collection_bank_bin <> '' AND i.collection_bank_account <> ''),
	COALESCE(log_counts.send_count, 0),
	log_counts.last_sent_at,
	COALESCE(log_counts.last_status, ''),
	COALESCE(log_counts.last_error, '')
FROM notification_recipients r
JOIN notification_campaigns c ON c.id = r.campaign_id
JOIN invoices i ON i.id = r.invoice_id
LEFT JOIN LATERAL (
	SELECT
		COUNT(*) FILTER (WHERE nl.status = 'sent')::integer AS send_count,
		MAX(nl.sent_at) FILTER (WHERE nl.status = 'sent') AS last_sent_at,
		(ARRAY_AGG(nl.status ORDER BY nl.sent_at DESC, nl.id DESC))[1] AS last_status,
		(ARRAY_AGG(nl.error_message ORDER BY nl.sent_at DESC, nl.id DESC))[1] AS last_error
	FROM notification_logs nl
	WHERE nl.campaign_id = r.campaign_id
		AND nl.invoice_id = r.invoice_id
		AND lower(nl.recipient_email) = lower(r.recipient_email)
) log_counts ON true
WHERE r.campaign_id = $1::uuid
	AND c.tenant_id = $2::uuid
ORDER BY r.class_name, r.student_code, r.recipient_name`, campaignID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipients := []notificationRecipientCandidate{}
	for rows.Next() {
		var item notificationRecipientCandidate
		var dueDate sql.NullTime
		var lastSentAt sql.NullTime
		var recipientLastError string
		var logLastError string
		if err := rows.Scan(
			&item.ID,
			&item.CampaignID,
			&item.InvoiceID,
			&item.ParentID,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.InvoiceCode,
			&item.StudentCode,
			&item.StudentName,
			&item.ClassName,
			&item.Grade,
			&item.SchoolYearID,
			&item.SchoolYearCode,
			&item.PeriodCode,
			&dueDate,
			&item.InvoiceStatus,
			&item.Amount,
			&item.PaidAmount,
			&item.OutstandingAmount,
			&item.Status,
			&recipientLastError,
			&item.QRReady,
			&item.SendCount,
			&lastSentAt,
			&item.LastLogStatus,
			&logLastError,
		); err != nil {
			return nil, err
		}
		item.LastError = firstNonEmpty(recipientLastError, logLastError)
		if dueDate.Valid {
			item.DueDate = dueDate.Time.Format("2006-01-02")
		}
		if lastSentAt.Valid {
			item.LastSentAt = lastSentAt.Time.UTC().Format(time.RFC3339)
		}
		recipients = append(recipients, finalizeNotificationRecipientState(item))
	}
	return recipients, rows.Err()
}

func sendNotificationCampaign(ctx context.Context, db *sql.DB, cfg emailConfig, template notificationTemplate, input notificationCampaignInput, recipients []notificationRecipientCandidate, baseURL string, sentLimit int) (notificationSendResponse, error) {
	if !input.DryRun {
		if _, err := db.ExecContext(ctx, `UPDATE notification_campaigns SET status = 'sending' WHERE id = $1::uuid AND tenant_id = $2::uuid`, input.CampaignID, input.TenantID); err != nil {
			return notificationSendResponse{}, err
		}
	}

	results := make([]emailSendResult, 0, len(recipients))
	sent := 0
	for _, recipient := range recipients {
		if !input.DryRun && sentLimit > 0 && sent >= sentLimit {
			result := notificationSkippedResult(recipient, "daily email limit reached")
			results = append(results, result)
			_ = insertNotificationLog(ctx, db, template, recipient, input.CampaignID, result, input.DryRun)
			_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
			continue
		}
		if !input.DryRun && !input.ForceResend {
			alreadySent, err := notificationAlreadySent(ctx, db, input.CampaignID, template, recipient.InvoiceID, recipient.RecipientEmail)
			if err != nil {
				return notificationSendResponse{}, err
			}
			if alreadySent {
				result := notificationSkippedResult(recipient, "already sent for campaign/template/invoice/recipient")
				results = append(results, result)
				_ = insertNotificationLog(ctx, db, template, recipient, input.CampaignID, result, input.DryRun)
				_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
				continue
			}
		}
		invoice, err := loadInvoiceDocument(ctx, db, recipient.InvoiceID, input.TenantID)
		if err != nil {
			result := notificationSkippedResult(recipient, "invoice not found")
			result.Status = "error"
			results = append(results, result)
			_ = insertNotificationLog(ctx, db, template, recipient, input.CampaignID, result, input.DryRun)
			_ = recordNotificationOperationLog(ctx, db, template, recipient, input, result)
			_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
			continue
		}
		row := notificationPaymentRow(invoice, recipient)
		emailCfg := notificationEmailConfig(cfg, template, invoice)
		result := sendPaymentEmailRow(ctx, emailCfg, row, template.EmailTemplate, baseURL, input.DryRun)
		result.ID = recipient.ID
		results = append(results, result)
		_ = insertNotificationLog(ctx, db, template, recipient, input.CampaignID, result, input.DryRun)
		_ = recordNotificationOperationLog(ctx, db, template, recipient, input, result)
		_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
		if result.Status == "sent" {
			sent++
			sleepEmailSendPace(ctx, emailCfg)
		}
		if result.Transient {
			break
		}
	}
	if !input.DryRun {
		recordEmailCronSentForTenant(ctx, input.TenantID, countSentEmails(results), time.Now())
		if err := incrementTenantUsageCounter(ctx, db, input.TenantID, subscriptionMetricMonthlyNotifications, subscriptionUsagePeriodKey(subscriptionMetricMonthlyNotifications, time.Now()), countSentEmails(results)); err != nil {
			return notificationSendResponse{}, err
		}
	}
	status := notificationCampaignStatusFromResults(results, input.DryRun)
	updateQuery := `UPDATE notification_campaigns SET status = $2`
	args := []any{input.CampaignID, status}
	if input.DryRun {
		updateQuery += `, last_dry_run_at = now()`
	} else if status == notificationStatusSent || status == notificationStatusPartial {
		updateQuery += `, sent_at = now()`
	}
	updateQuery += ` WHERE id = $1::uuid AND tenant_id = $3::uuid`
	args = append(args, input.TenantID)
	if _, err := db.ExecContext(ctx, updateQuery, args...); err != nil {
		return notificationSendResponse{}, err
	}

	campaign, err := loadNotificationCampaignSummary(ctx, db, input.CampaignID, input.TenantID)
	if err != nil {
		return notificationSendResponse{}, err
	}
	logs, _ := listNotificationLogs(ctx, db, input.CampaignID, input.TenantID, 50)
	updatedRecipients, err := loadNotificationRecipients(ctx, db, input.CampaignID, input.TenantID)
	if err != nil {
		updatedRecipients = recipients
	}
	return notificationSendResponse{
		Campaign: campaign,
		Summary:  summarizeNotificationRecipients(updatedRecipients),
		Results:  results,
		Logs:     logs,
	}, nil
}

func recordNotificationOperationLog(ctx context.Context, db *sql.DB, template notificationTemplate, recipient notificationRecipientCandidate, input notificationCampaignInput, result emailSendResult) error {
	if result.Status != "error" {
		return nil
	}
	return recordOperationLog(ctx, db, operationLogInput{
		Source:     "email",
		Level:      "error",
		Operation:  "notification.campaign.send",
		Status:     "error",
		Message:    result.Error,
		EntityType: "notification_recipient",
		EntityID:   recipient.ID,
		Metadata: map[string]any{
			"campaignId":      input.CampaignID,
			"templateCode":    template.Code,
			"templateVersion": template.Version,
			"invoiceId":       recipient.InvoiceID,
			"invoiceCode":     recipient.InvoiceCode,
			"recipientEmail":  recipient.RecipientEmail,
			"dryRun":          input.DryRun,
			"transient":       result.Transient,
		},
	})
}

func notificationPaymentRow(invoice invoiceDocument, recipient notificationRecipientCandidate) paymentRow {
	row := paymentRowFromInvoice(invoice)
	row.ID = firstNonEmpty(recipient.ID, invoice.ID, invoice.InvoiceCode)
	row.ParentName = recipient.RecipientName
	row.Email = recipient.RecipientEmail
	return row
}

func notificationEmailConfig(cfg emailConfig, template notificationTemplate, invoice invoiceDocument) emailConfig {
	cfg = normalizeEmailConfig(cfg)
	cfg.Subject = template.Subject
	cfg.PaymentPeriod = notificationPaymentPeriod(invoice)
	return cfg
}

func notificationPaymentPeriod(invoice invoiceDocument) string {
	parts := []string{}
	if invoice.PeriodCode != "" {
		parts = append(parts, "Kỳ "+invoice.PeriodCode)
	}
	if invoice.DueDate != "" {
		parts = append(parts, "hạn "+invoice.DueDate)
	}
	if len(parts) == 0 {
		return defaultEmailConfig().PaymentPeriod
	}
	return strings.Join(parts, " - ")
}

func notificationAlreadySent(ctx context.Context, db *sql.DB, campaignID string, template notificationTemplate, invoiceID string, email string) (bool, error) {
	var exists bool
	if template.Code == notificationCampaignPaymentConfirmation {
		err := db.QueryRowContext(ctx, `
SELECT EXISTS (
	SELECT 1
	FROM notification_logs
	WHERE template_code = $1
		AND invoice_id = $2::uuid
		AND lower(recipient_email) = lower($3)
		AND status = 'sent'
)`, template.Code, invoiceID, email).Scan(&exists)
		return exists, err
	}
	err := db.QueryRowContext(ctx, `
SELECT EXISTS (
	SELECT 1
	FROM notification_logs
	WHERE campaign_id = $1::uuid
		AND template_id = $2::uuid
		AND invoice_id = $3::uuid
		AND lower(recipient_email) = lower($4)
		AND status = 'sent'
)`, campaignID, template.ID, invoiceID, email).Scan(&exists)
	return exists, err
}

func countNotificationRecipientsPlannedForSend(ctx context.Context, db *sql.DB, campaignID string, template notificationTemplate, recipients []notificationRecipientCandidate, forceResend bool, sendLimit int) (int, error) {
	if sendLimit <= 0 {
		return 0, nil
	}
	planned := 0
	for _, recipient := range recipients {
		if sendLimit > 0 && planned >= sendLimit {
			break
		}
		if !forceResend {
			alreadySent, err := notificationAlreadySent(ctx, db, campaignID, template, recipient.InvoiceID, recipient.RecipientEmail)
			if err != nil {
				return 0, err
			}
			if alreadySent {
				continue
			}
		}
		planned++
	}
	return planned, nil
}

func sendAutomaticPaidConfirmationBestEffort(ctx context.Context, db *sql.DB, invoiceID string, trigger string) {
	if strings.TrimSpace(invoiceID) == "" {
		return
	}
	tenantID, err := tenantIDForInvoice(ctx, db, invoiceID)
	if err != nil {
		_ = recordOperationLog(ctx, db, operationLogInput{
			Source:     "email",
			Level:      "error",
			Operation:  "notification.paid_confirmation.auto",
			Status:     "error",
			Message:    err.Error(),
			EntityType: "invoice",
			EntityID:   invoiceID,
			Metadata: map[string]any{
				"trigger": trigger,
			},
		})
		return
	}
	if _, err := sendPaidConfirmationForInvoice(ctx, db, invoiceID, tenantID, paidConfirmationSendOptions{
		CampaignCode: notificationAutoPaidConfirmationCampaignCode,
		CampaignName: "Tự động xác nhận đã thanh toán",
		Trigger:      trigger,
		Operation:    "notification.paid_confirmation.auto",
	}); err != nil {
		_ = recordOperationLog(ctx, db, operationLogInput{
			Source:     "email",
			Level:      "error",
			Operation:  "notification.paid_confirmation.auto",
			Status:     "error",
			Message:    err.Error(),
			EntityType: "invoice",
			EntityID:   invoiceID,
			Metadata: map[string]any{
				"trigger": trigger,
			},
		})
	}
}

func tenantIDForInvoice(ctx context.Context, db *sql.DB, invoiceID string) (string, error) {
	var tenantID string
	err := db.QueryRowContext(ctx, `
SELECT sc.tenant_id::text
FROM invoices i
JOIN school_years sy ON sy.id = i.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE i.id = $1::uuid`, invoiceID).Scan(&tenantID)
	return tenantID, err
}

func sendPaidConfirmationForInvoice(ctx context.Context, db *sql.DB, invoiceID string, tenantID string, options paidConfirmationSendOptions) ([]emailSendResult, error) {
	invoiceID = strings.TrimSpace(invoiceID)
	if invoiceID == "" {
		return nil, errors.New("invoiceId is required")
	}
	invoice, err := loadInvoiceDocument(ctx, db, invoiceID, tenantID)
	if err != nil {
		return nil, err
	}
	if invoice.Status != invoiceStatusPaid {
		return nil, fmt.Errorf("invoice %s is not paid", invoice.InvoiceCode)
	}

	template, err := loadLatestNotificationTemplateByCode(ctx, db, notificationCampaignPaymentConfirmation)
	if err != nil {
		return nil, err
	}
	campaignID, err := ensurePaidConfirmationCampaign(ctx, db, template, invoice, tenantID, options)
	if err != nil {
		return nil, err
	}
	recipients, err := listPaidConfirmationRecipients(ctx, db, invoiceID, tenantID)
	if err != nil {
		return nil, err
	}
	if len(recipients) == 0 {
		_ = recordPaidConfirmationOperationLog(ctx, db, options, invoiceID, notificationRecipientCandidate{}, notificationSkippedResult(notificationRecipientCandidate{InvoiceID: invoiceID, InvoiceCode: invoice.InvoiceCode}, "paid invoice has no active billing recipients"), "warn")
		return []emailSendResult{}, nil
	}
	for idx := range recipients {
		recipientID, err := upsertNotificationRecipient(ctx, db, campaignID, recipients[idx])
		if err != nil {
			return nil, err
		}
		recipients[idx].ID = recipientID
		recipients[idx].CampaignID = campaignID
	}

	cfg, err := loadEmailConfigForTenant(ctx, tenantID)
	if err != nil {
		return recordPaidConfirmationFailure(ctx, db, template, campaignID, recipients, options, "", err.Error()), nil
	}
	cfg = normalizeEmailConfig(cfg)
	if err := validateEmailConfigForSend(cfg); err != nil {
		return recordPaidConfirmationFailure(ctx, db, template, campaignID, recipients, options, cfg.Provider, err.Error()), nil
	}
	quota, err := emailSendQuotaStatusForTenant(ctx, tenantID, time.Now())
	if err != nil {
		return recordPaidConfirmationFailure(ctx, db, template, campaignID, recipients, options, cfg.Provider, err.Error()), nil
	}
	plannedSends, err := countPaidConfirmationRecipientsPlannedForSend(ctx, db, campaignID, template, recipients, options.ForceResend, quota.Remaining)
	if err != nil {
		return nil, err
	}
	if err := enforceTenantUsageLimit(ctx, db, tenantID, subscriptionMetricMonthlyNotifications, plannedSends, time.Now()); err != nil {
		return nil, err
	}

	results := make([]emailSendResult, 0, len(recipients))
	sent := 0
	for _, recipient := range recipients {
		if !options.ForceResend {
			alreadySent, err := notificationAlreadySent(ctx, db, campaignID, template, recipient.InvoiceID, recipient.RecipientEmail)
			if err != nil {
				return results, err
			}
			if alreadySent {
				continue
			}
		}
		if sent >= quota.Remaining {
			result := notificationSkippedResult(recipient, "daily email limit reached")
			result.Provider = cfg.Provider
			results = append(results, result)
			_ = insertNotificationLog(ctx, db, template, recipient, campaignID, result, false)
			_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
			_ = recordPaidConfirmationOperationLog(ctx, db, options, invoiceID, recipient, result, "warn")
			continue
		}

		emailCfg := notificationEmailConfig(cfg, template, invoice)
		result := sendPaymentEmailRow(ctx, emailCfg, notificationPaymentRow(invoice, recipient), template.EmailTemplate, schedulerBaseURL(emailCfg), false)
		result.ID = recipient.ID
		results = append(results, result)
		_ = insertNotificationLog(ctx, db, template, recipient, campaignID, result, false)
		_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
		_ = recordPaidConfirmationOperationLog(ctx, db, options, invoiceID, recipient, result, "error")
		if result.Status == "sent" {
			sent++
			sleepEmailSendPace(ctx, emailCfg)
		}
		if result.Transient {
			break
		}
	}
	if sent > 0 {
		recordEmailCronSentForTenant(ctx, tenantID, sent, time.Now())
		if err := incrementTenantUsageCounter(ctx, db, tenantID, subscriptionMetricMonthlyNotifications, subscriptionUsagePeriodKey(subscriptionMetricMonthlyNotifications, time.Now()), sent); err != nil {
			return nil, err
		}
	}
	status := notificationCampaignStatusFromResults(results, false)
	if len(results) == 0 {
		status = notificationStatusSent
	}
	_, _ = db.ExecContext(ctx, `
UPDATE notification_campaigns
SET status = $2,
	sent_at = CASE WHEN $2 IN ('sent', 'partial') THEN now() ELSE sent_at END
WHERE id = $1::uuid
	AND tenant_id = $3::uuid`, campaignID, status, tenantID)
	return results, nil
}

func countPaidConfirmationRecipientsPlannedForSend(ctx context.Context, db *sql.DB, campaignID string, template notificationTemplate, recipients []notificationRecipientCandidate, forceResend bool, sendLimit int) (int, error) {
	if sendLimit <= 0 {
		return 0, nil
	}
	planned := 0
	for _, recipient := range recipients {
		if sendLimit > 0 && planned >= sendLimit {
			break
		}
		if !forceResend {
			alreadySent, err := notificationAlreadySent(ctx, db, campaignID, template, recipient.InvoiceID, recipient.RecipientEmail)
			if err != nil {
				return 0, err
			}
			if alreadySent {
				continue
			}
		}
		planned++
	}
	return planned, nil
}

func ensurePaidConfirmationCampaign(ctx context.Context, db *sql.DB, template notificationTemplate, invoice invoiceDocument, tenantID string, options paidConfirmationSendOptions) (string, error) {
	code := strings.TrimSpace(options.CampaignCode)
	name := strings.TrimSpace(options.CampaignName)
	if code == "" {
		if options.ForceResend {
			now := time.Now()
			code = cleanANS("NC"+now.Format("20060102150405")+fmt.Sprintf("%09d", now.Nanosecond())+"PAID", 32)
		} else {
			code = notificationAutoPaidConfirmationCampaignCode
		}
	}
	if name == "" {
		name = "Xác nhận đã thanh toán " + firstNonEmpty(invoice.InvoiceCode, invoice.StudentCode, time.Now().Format("2006-01-02"))
	}
	autoCampaign := code == notificationAutoPaidConfirmationCampaignCode
	filter := map[string]any{
		"auto":    autoCampaign,
		"trigger": strings.TrimSpace(options.Trigger),
	}
	schoolYearID := nullableString(invoice.SchoolYearID)
	classID := nullableString(invoice.ClassID)
	grade := invoice.Grade
	periodCode := invoice.PeriodCode
	if !autoCampaign {
		filter["invoiceId"] = invoice.ID
		filter["invoiceCode"] = invoice.InvoiceCode
	} else {
		schoolYearID = nil
		classID = nil
		grade = ""
		periodCode = ""
	}
	filterJSON, err := jsonObjectString(filter)
	if err != nil {
		return "", err
	}
	var campaignID string
	err = db.QueryRowContext(ctx, `
INSERT INTO notification_campaigns (
	tenant_id, code, name, campaign_type, template_id, school_year_id, class_id,
	grade, period_code, invoice_status, status, target_filter
)
VALUES ($1::uuid, $2, $3, $4, $5::uuid, $6::uuid, $7::uuid, $8, $9, 'paid', 'draft', $10::jsonb)
ON CONFLICT (tenant_id, code) DO UPDATE
SET name = EXCLUDED.name,
	campaign_type = EXCLUDED.campaign_type,
	template_id = EXCLUDED.template_id,
	school_year_id = EXCLUDED.school_year_id,
	class_id = EXCLUDED.class_id,
	grade = EXCLUDED.grade,
	period_code = EXCLUDED.period_code,
	invoice_status = EXCLUDED.invoice_status,
	target_filter = EXCLUDED.target_filter,
	status = CASE
		WHEN notification_campaigns.status = 'archived' THEN 'draft'
		ELSE notification_campaigns.status
END
RETURNING id::text`,
		tenantID,
		code,
		name,
		notificationCampaignPaymentConfirmation,
		template.ID,
		schoolYearID,
		classID,
		grade,
		periodCode,
		filterJSON,
	).Scan(&campaignID)
	return campaignID, err
}

func listPaidConfirmationRecipients(ctx context.Context, db *sql.DB, invoiceID string, tenantID string) ([]notificationRecipientCandidate, error) {
	rows, err := db.QueryContext(ctx, `
SELECT i.id::text,
	COALESCE(p.id::text, ''),
	p.full_name,
	p.email,
	i.invoice_code,
	i.student_code,
	i.student_name,
	i.class_name,
	i.grade,
	i.school_year_id::text,
	i.school_year_code,
	i.period_code,
	i.due_date,
	i.status,
	i.total_amount,
	i.paid_amount,
	GREATEST(i.total_amount - i.paid_amount, 0),
	(i.qr_bill_number <> '' AND i.collection_bank_bin <> '' AND i.collection_bank_account <> ''),
	COALESCE(log_counts.send_count, 0),
	log_counts.last_sent_at,
	COALESCE(log_counts.last_status, ''),
	COALESCE(log_counts.last_error, '')
FROM invoices i
JOIN school_years sy ON sy.id = i.school_year_id
JOIN schools sc ON sc.id = sy.school_id
JOIN students s ON s.id = i.student_id
JOIN student_parents sp ON sp.student_id = s.id
JOIN parents p ON p.id = sp.parent_id
LEFT JOIN LATERAL (
	SELECT
		COUNT(*) FILTER (WHERE nl.status = 'sent')::integer AS send_count,
		MAX(nl.sent_at) FILTER (WHERE nl.status = 'sent') AS last_sent_at,
		(ARRAY_AGG(nl.status ORDER BY nl.sent_at DESC, nl.id DESC))[1] AS last_status,
		(ARRAY_AGG(nl.error_message ORDER BY nl.sent_at DESC, nl.id DESC))[1] AS last_error
	FROM notification_logs nl
	WHERE nl.template_code = 'payment_confirmation'
		AND nl.invoice_id = i.id
		AND lower(nl.recipient_email) = lower(p.email)
) log_counts ON true
WHERE i.id = $1::uuid
	AND sc.tenant_id = $2::uuid
	AND i.status = 'paid'
	AND sp.is_active
	AND sp.receives_billing_email
	AND p.email_active
	AND p.status = 'active'
	AND p.email <> ''
ORDER BY sp.is_primary DESC, p.full_name
LIMIT 20`, invoiceID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipients := []notificationRecipientCandidate{}
	seen := map[string]bool{}
	for rows.Next() {
		var item notificationRecipientCandidate
		var dueDate sql.NullTime
		var lastSentAt sql.NullTime
		if err := rows.Scan(
			&item.InvoiceID,
			&item.ParentID,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.InvoiceCode,
			&item.StudentCode,
			&item.StudentName,
			&item.ClassName,
			&item.Grade,
			&item.SchoolYearID,
			&item.SchoolYearCode,
			&item.PeriodCode,
			&dueDate,
			&item.InvoiceStatus,
			&item.Amount,
			&item.PaidAmount,
			&item.OutstandingAmount,
			&item.QRReady,
			&item.SendCount,
			&lastSentAt,
			&item.LastLogStatus,
			&item.LastError,
		); err != nil {
			return nil, err
		}
		item.RecipientEmail = strings.ToLower(strings.TrimSpace(item.RecipientEmail))
		if seen[item.RecipientEmail] {
			continue
		}
		seen[item.RecipientEmail] = true
		if dueDate.Valid {
			item.DueDate = dueDate.Time.Format("2006-01-02")
		}
		if lastSentAt.Valid {
			item.LastSentAt = lastSentAt.Time.UTC().Format(time.RFC3339)
		}
		recipients = append(recipients, finalizeNotificationRecipientState(item))
	}
	return recipients, rows.Err()
}

func upsertNotificationRecipient(ctx context.Context, db *sql.DB, campaignID string, recipient notificationRecipientCandidate) (string, error) {
	if err := insertNotificationRecipient(ctx, db, campaignID, recipient); err != nil {
		return "", err
	}
	var recipientID string
	err := db.QueryRowContext(ctx, `
SELECT id::text
FROM notification_recipients
WHERE campaign_id = $1::uuid
	AND invoice_id = $2::uuid
	AND lower(recipient_email) = lower($3)
LIMIT 1`, campaignID, recipient.InvoiceID, recipient.RecipientEmail).Scan(&recipientID)
	return recipientID, err
}

func recordPaidConfirmationFailure(ctx context.Context, db *sql.DB, template notificationTemplate, campaignID string, recipients []notificationRecipientCandidate, options paidConfirmationSendOptions, provider string, message string) []emailSendResult {
	results := make([]emailSendResult, 0, len(recipients))
	for _, recipient := range recipients {
		result := notificationErrorResult(recipient, provider, message)
		results = append(results, result)
		_ = insertNotificationLog(ctx, db, template, recipient, campaignID, result, false)
		_ = updateNotificationRecipientStatus(ctx, db, recipient.ID, result.Status, result.Error)
		_ = recordPaidConfirmationOperationLog(ctx, db, options, recipient.InvoiceID, recipient, result, "error")
	}
	return results
}

func notificationErrorResult(recipient notificationRecipientCandidate, provider string, message string) emailSendResult {
	return emailSendResult{
		ID:          recipient.ID,
		Email:       recipient.RecipientEmail,
		StudentName: recipient.StudentName,
		Provider:    provider,
		Status:      "error",
		Error:       message,
	}
}

func recordPaidConfirmationOperationLog(ctx context.Context, db *sql.DB, options paidConfirmationSendOptions, invoiceID string, recipient notificationRecipientCandidate, result emailSendResult, defaultLevel string) error {
	if result.Status != "error" && result.Status != "skipped" {
		return nil
	}
	level := defaultLevel
	if level == "" {
		level = "error"
	}
	operation := strings.TrimSpace(options.Operation)
	if operation == "" {
		operation = "notification.paid_confirmation.send"
	}
	status := result.Status
	if status == "" {
		status = "error"
	}
	return recordOperationLog(ctx, db, operationLogInput{
		Source:     "email",
		Level:      level,
		Operation:  operation,
		Status:     status,
		Message:    result.Error,
		EntityType: "invoice",
		EntityID:   invoiceID,
		Metadata: map[string]any{
			"trigger":        strings.TrimSpace(options.Trigger),
			"recipientId":    recipient.ID,
			"recipientEmail": recipient.RecipientEmail,
			"invoiceCode":    firstNonEmpty(recipient.InvoiceCode, result.ID),
			"forceResend":    options.ForceResend,
			"provider":       result.Provider,
			"transient":      result.Transient,
		},
	})
}

func insertNotificationLog(ctx context.Context, db *sql.DB, template notificationTemplate, recipient notificationRecipientCandidate, campaignID string, result emailSendResult, dryRun bool) error {
	status := result.Status
	if status == "" {
		status = "error"
	}
	if status != "sent" && status != "dry_run" && status != "skipped" && status != "error" {
		status = "error"
	}
	messageID := firstNonEmpty(result.MessageID, result.ResendID)
	_, err := db.ExecContext(ctx, `
INSERT INTO notification_logs (
	campaign_id, template_id, template_code, template_version, recipient_id,
	invoice_id, recipient_email, provider, status, provider_message_id,
	error_message, dry_run
)
VALUES ($1::uuid, $2::uuid, $3, $4, $5::uuid, $6::uuid, $7, $8, $9, $10, $11, $12)
ON CONFLICT DO NOTHING`,
		campaignID,
		template.ID,
		template.Code,
		template.Version,
		nullableString(recipient.ID),
		recipient.InvoiceID,
		recipient.RecipientEmail,
		result.Provider,
		status,
		messageID,
		result.Error,
		dryRun,
	)
	return err
}

func updateNotificationRecipientStatus(ctx context.Context, db *sql.DB, recipientID string, status string, lastError string) error {
	if recipientID == "" {
		return nil
	}
	if status != "sent" && status != "dry_run" && status != "skipped" && status != "error" {
		status = "error"
	}
	_, err := db.ExecContext(ctx, `
UPDATE notification_recipients
SET status = $2,
	last_error = $3
WHERE id = $1::uuid`, recipientID, status, lastError)
	return err
}

func listNotificationLogs(ctx context.Context, db *sql.DB, campaignID string, tenantID string, limit int) ([]notificationLogSummary, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	conditions := []string{}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if tenantID != "" {
		conditions = append(conditions, "c.tenant_id = "+addArg(tenantID)+"::uuid")
	}
	if campaignID != "" {
		conditions = append(conditions, "nl.campaign_id = "+addArg(campaignID)+"::uuid")
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	query := `
SELECT nl.id::text,
	nl.campaign_id::text,
	c.name,
	nl.template_code,
	nl.template_version,
	nl.invoice_id::text,
	i.invoice_code,
	nl.recipient_email,
	nl.provider,
	nl.status,
	nl.provider_message_id,
	nl.error_message,
	nl.dry_run,
	nl.sent_at
FROM notification_logs nl
JOIN notification_campaigns c ON c.id = nl.campaign_id
JOIN invoices i ON i.id = nl.invoice_id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY nl.sent_at DESC, nl.id DESC
LIMIT ` + addArg(limit)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	logs := []notificationLogSummary{}
	for rows.Next() {
		var item notificationLogSummary
		if err := rows.Scan(
			&item.ID,
			&item.CampaignID,
			&item.CampaignName,
			&item.TemplateCode,
			&item.TemplateVersion,
			&item.InvoiceID,
			&item.InvoiceCode,
			&item.RecipientEmail,
			&item.Provider,
			&item.Status,
			&item.ProviderMessageID,
			&item.Error,
			&item.DryRun,
			&item.SentAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, item)
	}
	return logs, rows.Err()
}

func notificationSkippedResult(recipient notificationRecipientCandidate, message string) emailSendResult {
	return emailSendResult{
		ID:          recipient.ID,
		Email:       recipient.RecipientEmail,
		StudentName: recipient.StudentName,
		Status:      "skipped",
		Error:       message,
	}
}

func notificationCampaignStatusFromResults(results []emailSendResult, dryRun bool) string {
	if dryRun {
		return notificationStatusDryRun
	}
	hasError := false
	hasSent := false
	for _, result := range results {
		switch result.Status {
		case "sent":
			hasSent = true
		case "error", "skipped":
			hasError = true
		}
	}
	if hasError || !hasSent {
		return notificationStatusPartial
	}
	return notificationStatusSent
}

func notificationCampaignSummaryFromInput(input notificationCampaignInput, template notificationTemplate, status string, recipients []notificationRecipientCandidate) notificationCampaignSummary {
	return notificationCampaignSummary{
		ID:              input.CampaignID,
		Name:            input.Name,
		CampaignType:    input.CampaignType,
		Status:          status,
		Template:        template.Code,
		TemplateID:      template.ID,
		TemplateVersion: template.Version,
		SchoolYearID:    input.SchoolYearID,
		ClassID:         input.ClassID,
		Grade:           input.Grade,
		PeriodCode:      input.PeriodCode,
		InvoiceStatus:   input.InvoiceStatus,
		DueOnOrBefore:   input.DueOnOrBefore,
		RecipientCount:  len(recipients),
	}
}

func summarizeNotificationRecipients(recipients []notificationRecipientCandidate) notificationRecipientSummary {
	summary := notificationRecipientSummary{RecipientCount: len(recipients)}
	seenInvoices := map[string]bool{}
	for _, recipient := range recipients {
		if recipient.AlreadySent {
			summary.AlreadySent++
		}
		if !recipient.QRReady {
			summary.QRMissingCount++
		}
		if recipient.Status == "error" || recipient.LastLogStatus == "error" {
			summary.ErrorCount++
		}
		if recipient.RetryEligible {
			summary.RetryEligibleCount++
		}
		if seenInvoices[recipient.InvoiceID] {
			continue
		}
		seenInvoices[recipient.InvoiceID] = true
		summary.InvoiceCount++
		summary.TotalAmount += recipient.Amount
		outstanding := recipient.OutstandingAmount
		if outstanding == 0 && recipient.Amount > recipient.PaidAmount {
			outstanding = recipient.Amount - recipient.PaidAmount
		}
		if outstanding > 0 {
			summary.UnpaidAmount += outstanding
		}
	}
	return summary
}

func notificationTargetFilterJSON(input notificationCampaignInput) (string, error) {
	payload := map[string]any{
		"campaignType":  input.CampaignType,
		"schoolYearId":  input.SchoolYearID,
		"classId":       input.ClassID,
		"grade":         input.Grade,
		"periodCode":    input.PeriodCode,
		"invoiceStatus": input.InvoiceStatus,
		"dueOnOrBefore": input.DueOnOrBefore,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func notificationCampaignCode(input notificationCampaignInput, now time.Time) string {
	tag := strings.ToUpper(safeTagValue(input.CampaignType))
	if tag == "" || tag == "UNKNOWN" {
		tag = "NOTICE"
	}
	return cleanANS("NC"+now.Format("20060102150405")+tag, 32)
}
