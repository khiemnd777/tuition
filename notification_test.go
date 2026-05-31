package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNotificationReminderRejectsPaidInvoiceStatus(t *testing.T) {
	issues := validateNotificationCampaignInput(normalizeNotificationCampaignInput(notificationCampaignInput{
		CampaignType:  notificationCampaignReminder,
		InvoiceStatus: invoiceStatusPaid,
	}))
	if len(issues) != 1 {
		t.Fatalf("expected one issue, got %+v", issues)
	}
	if issues[0].Type != "invalid_reminder_status" {
		t.Fatalf("expected invalid_reminder_status, got %+v", issues[0])
	}
}

func TestNotificationPaymentRowKeepsInvoicePaymentContract(t *testing.T) {
	invoice := invoiceDocument{
		invoiceSummary: invoiceSummary{
			ID:                    "invoice-1",
			InvoiceCode:           "SUNTEST001",
			StudentCode:           "S001",
			StudentName:           "Nguyen An",
			ClassName:             "3.02",
			PeriodCode:            "2026-04",
			IssuedAt:              time.Date(2026, 5, 29, 9, 0, 0, 0, time.UTC),
			Status:                invoiceStatusUnpaid,
			TotalAmount:           1200,
			CollectionBankBIN:     "970415",
			CollectionBankAccount: "0011001932418",
			QRBillNumber:          "SUNTEST001",
			QRNote:                "HP S001 2026-04",
		},
		Items: []invoiceDocumentItem{
			{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1200, DisplayOrder: 1},
		},
	}
	recipient := notificationRecipientCandidate{
		ID:             "recipient-1",
		RecipientName:  "Nguyen Van Binh",
		RecipientEmail: "billing@example.com",
	}

	row := notificationPaymentRow(invoice, recipient)
	if row.Email != recipient.RecipientEmail {
		t.Fatalf("expected recipient email %q, got %q", recipient.RecipientEmail, row.Email)
	}
	if row.ParentName != recipient.RecipientName {
		t.Fatalf("expected recipient name %q, got %q", recipient.RecipientName, row.ParentName)
	}
	if row.BillNumber != invoice.InvoiceCode {
		t.Fatalf("expected invoice code bill number %q, got %q", invoice.InvoiceCode, row.BillNumber)
	}
	if got := paymentItemsTotal(row.PaymentItems); got != invoice.TotalAmount {
		t.Fatalf("expected payment item total %d, got %d", invoice.TotalAmount, got)
	}
}

func TestNotificationRecipientSummaryCountsInvoiceTotalsOnce(t *testing.T) {
	recipients := []notificationRecipientCandidate{
		{InvoiceID: "invoice-1", Amount: 1200, PaidAmount: 0, RecipientEmail: "one@example.com"},
		{InvoiceID: "invoice-1", Amount: 1200, PaidAmount: 0, RecipientEmail: "two@example.com"},
		{InvoiceID: "invoice-2", Amount: 2000, PaidAmount: 500, RecipientEmail: "three@example.com", AlreadySent: true},
	}
	summary := summarizeNotificationRecipients(recipients)
	if summary.RecipientCount != 3 {
		t.Fatalf("expected 3 recipients, got %d", summary.RecipientCount)
	}
	if summary.InvoiceCount != 2 {
		t.Fatalf("expected 2 invoices, got %d", summary.InvoiceCount)
	}
	if summary.TotalAmount != 3200 {
		t.Fatalf("expected unique invoice total 3200, got %d", summary.TotalAmount)
	}
	if summary.UnpaidAmount != 2700 {
		t.Fatalf("expected unpaid amount 2700, got %d", summary.UnpaidAmount)
	}
	if summary.AlreadySent != 1 {
		t.Fatalf("expected one already-sent recipient, got %d", summary.AlreadySent)
	}
}

func TestNotificationRecipientSummaryTracksRetryAndQRState(t *testing.T) {
	recipients := []notificationRecipientCandidate{
		{InvoiceID: "invoice-1", Amount: 1200, OutstandingAmount: 1000, RecipientEmail: "one@example.com", QRReady: true},
		{InvoiceID: "invoice-2", Amount: 2000, PaidAmount: 500, RecipientEmail: "two@example.com", Status: "error", RetryEligible: true},
		{InvoiceID: "invoice-3", Amount: 3000, PaidAmount: 3000, RecipientEmail: "three@example.com", LastLogStatus: "error", RetryEligible: true},
	}

	summary := summarizeNotificationRecipients(recipients)
	if summary.QRMissingCount != 2 {
		t.Fatalf("expected two missing QR recipients, got %d", summary.QRMissingCount)
	}
	if summary.ErrorCount != 2 {
		t.Fatalf("expected two error recipients, got %d", summary.ErrorCount)
	}
	if summary.RetryEligibleCount != 2 {
		t.Fatalf("expected two retry-eligible recipients, got %d", summary.RetryEligibleCount)
	}
	if summary.UnpaidAmount != 2500 {
		t.Fatalf("expected unpaid amount 2500, got %d", summary.UnpaidAmount)
	}
}

func TestFilterNotificationRecipientsForSend(t *testing.T) {
	recipients := []notificationRecipientCandidate{
		{ID: "recipient-1", RecipientEmail: "one@example.com"},
		{ID: "recipient-2", RecipientEmail: "two@example.com"},
		{ID: "recipient-3", RecipientEmail: "three@example.com"},
	}

	filtered := filterNotificationRecipientsForSend(recipients, []string{"recipient-2", "recipient-2", "missing", ""})
	if len(filtered) != 1 {
		t.Fatalf("expected one filtered recipient, got %+v", filtered)
	}
	if filtered[0].ID != "recipient-2" {
		t.Fatalf("expected recipient-2, got %+v", filtered[0])
	}

	unfiltered := filterNotificationRecipientsForSend(recipients, nil)
	if len(unfiltered) != len(recipients) {
		t.Fatalf("expected unfiltered recipients, got %+v", unfiltered)
	}
}

func TestSelectNotificationPreviewRecipient(t *testing.T) {
	recipients := []notificationRecipientCandidate{
		{ID: "recipient-1", InvoiceID: "invoice-1", RecipientEmail: "one@example.com"},
		{ID: "recipient-2", InvoiceID: "invoice-2", RecipientEmail: "two@example.com"},
	}

	selected, ok := selectNotificationPreviewRecipient(recipients, "", "invoice-2", "TWO@example.com")
	if !ok {
		t.Fatal("expected recipient selection by invoice and email")
	}
	if selected.ID != "recipient-2" {
		t.Fatalf("expected recipient-2, got %+v", selected)
	}

	selected, ok = selectNotificationPreviewRecipient(recipients, "", "", "")
	if !ok || selected.ID != "recipient-1" {
		t.Fatalf("expected first recipient fallback, got %+v ok=%v", selected, ok)
	}
}

func TestDecodeNotificationEmailPreviewInputNormalizesCampaignFields(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/notifications/campaigns/email-preview", strings.NewReader(`{
		"campaignId": " campaign-1 ",
		"campaignType": "reminder",
		"invoiceStatus": "partial",
		"forceResend": true,
		"recipientIds": [" recipient-1 ", "recipient-1", ""],
		"recipientId": " recipient-2 ",
		"invoiceId": " invoice-1 ",
		"recipientEmail": "Parent@Example.COM "
	}`))
	rr := httptest.NewRecorder()

	input, ok := decodeNotificationEmailPreviewInput(rr, req)
	if !ok {
		t.Fatalf("expected decode to succeed, status %d body %q", rr.Code, rr.Body.String())
	}
	if input.CampaignID != "campaign-1" {
		t.Fatalf("expected campaign id to be normalized, got %q", input.CampaignID)
	}
	if input.CampaignType != notificationCampaignReminder {
		t.Fatalf("expected reminder campaign type, got %q", input.CampaignType)
	}
	if !input.ForceResend {
		t.Fatal("expected forceResend to decode")
	}
	if len(input.RecipientIDs) != 1 || input.RecipientIDs[0] != "recipient-1" {
		t.Fatalf("expected recipient IDs to be normalized, got %+v", input.RecipientIDs)
	}
	if input.RecipientID != "recipient-2" || input.InvoiceID != "invoice-1" || input.RecipientEmail != "parent@example.com" {
		t.Fatalf("expected recipient selectors to be normalized, got %+v", input)
	}
}
