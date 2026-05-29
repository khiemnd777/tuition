package main

import (
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
