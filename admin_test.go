package main

import (
	"strings"
	"testing"
	"time"
)

func TestBuildAdminDashboardSummarySeparatesReceivablesAndReviewWork(t *testing.T) {
	invoices := []adminInvoiceReportRow{
		{invoiceSummary: invoiceSummary{StudentCode: "S001", ClassName: "3.02", Status: invoiceStatusUnpaid, TotalAmount: 1000, PaidAmount: 0}},
		{invoiceSummary: invoiceSummary{StudentCode: "S002", ClassName: "3.02", Status: invoiceStatusPartial, TotalAmount: 2000, PaidAmount: 750}},
		{invoiceSummary: invoiceSummary{StudentCode: "S003", ClassName: "4.01", Status: invoiceStatusPaid, TotalAmount: 3000, PaidAmount: 3000}},
		{invoiceSummary: invoiceSummary{StudentCode: "S004", ClassName: "4.01", Status: invoiceStatusOverpaid, TotalAmount: 4000, PaidAmount: 4500}},
		{invoiceSummary: invoiceSummary{StudentCode: "S005", ClassName: "5.01", Status: invoiceStatusManualReview, TotalAmount: 5000, PaidAmount: 0}},
	}
	transactions := []paymentTransactionSummary{
		{Status: paymentTransactionStatusUnmatched},
		{Status: paymentTransactionStatusManualReview},
		{Status: paymentTransactionStatusMatched},
	}

	summary := buildAdminDashboardSummary(invoices, transactions)

	if summary.InvoiceCount != 5 || summary.StudentCount != 5 {
		t.Fatalf("unexpected count summary: %+v", summary)
	}
	if summary.TotalReceivable != 15000 || summary.TotalCollected != 8250 || summary.OutstandingAmount != 7250 {
		t.Fatalf("unexpected money summary: %+v", summary)
	}
	if summary.UnpaidStudentCount != 1 || summary.PartialPaymentCount != 1 || summary.PaidInvoiceCount != 1 || summary.OverpaidManualReviewCount != 2 {
		t.Fatalf("unexpected status summary: %+v", summary)
	}
	if summary.UnmatchedTransactionCount != 1 || summary.ManualReviewCount != 1 {
		t.Fatalf("unexpected transaction summary: %+v", summary)
	}
}

func TestBuildAdminClassReportRowsSortsByOutstandingAmount(t *testing.T) {
	invoices := []adminInvoiceReportRow{
		{invoiceSummary: invoiceSummary{SchoolYearCode: "2025-2026", Grade: "3", ClassName: "3.02", StudentCode: "S001", Status: invoiceStatusUnpaid, TotalAmount: 1000, PaidAmount: 0}},
		{invoiceSummary: invoiceSummary{SchoolYearCode: "2025-2026", Grade: "3", ClassName: "3.02", StudentCode: "S002", Status: invoiceStatusPartial, TotalAmount: 2000, PaidAmount: 500}},
		{invoiceSummary: invoiceSummary{SchoolYearCode: "2025-2026", Grade: "4", ClassName: "4.01", StudentCode: "S003", Status: invoiceStatusPaid, TotalAmount: 3000, PaidAmount: 3000}},
		{invoiceSummary: invoiceSummary{SchoolYearCode: "2025-2026", Grade: "5", ClassName: "5.01", StudentCode: "S004", Status: invoiceStatusManualReview, TotalAmount: 5000, PaidAmount: 1000}},
	}

	rows := buildAdminClassReportRows(invoices)

	if len(rows) != 3 {
		t.Fatalf("expected three class rows, got %+v", rows)
	}
	if rows[0].ClassName != "5.01" || rows[0].Outstanding != 4000 || rows[0].ReviewCount != 1 {
		t.Fatalf("expected highest outstanding class first, got %+v", rows[0])
	}
	if rows[1].ClassName != "3.02" || rows[1].StudentCount != 2 || rows[1].UnpaidCount != 1 || rows[1].PartialCount != 1 {
		t.Fatalf("unexpected 3.02 aggregation: %+v", rows[1])
	}
}

func TestAdminAttentionInvoicesKeepsUnpaidPartialAndReviewOnly(t *testing.T) {
	invoices := []adminInvoiceReportRow{
		{invoiceSummary: invoiceSummary{InvoiceCode: "INV-PAID", Status: invoiceStatusPaid, TotalAmount: 1000, PaidAmount: 1000}, OutstandingAmount: 0},
		{invoiceSummary: invoiceSummary{InvoiceCode: "INV-LOW", Status: invoiceStatusPartial, TotalAmount: 1000, PaidAmount: 750}, OutstandingAmount: 250},
		{invoiceSummary: invoiceSummary{InvoiceCode: "INV-HIGH", Status: invoiceStatusUnpaid, TotalAmount: 2000, PaidAmount: 0}, OutstandingAmount: 2000},
		{invoiceSummary: invoiceSummary{InvoiceCode: "INV-REVIEW", Status: invoiceStatusManualReview, TotalAmount: 1500, PaidAmount: 0}, OutstandingAmount: 1500},
	}

	rows := adminAttentionInvoices(invoices, 2)

	if len(rows) != 2 {
		t.Fatalf("expected limited attention rows, got %+v", rows)
	}
	if rows[0].InvoiceCode != "INV-HIGH" || rows[1].InvoiceCode != "INV-REVIEW" {
		t.Fatalf("expected attention rows sorted by outstanding amount, got %+v", rows)
	}
}

func TestValidateAdminUserSaveInputNormalizesFields(t *testing.T) {
	input := adminUserSaveInput{
		Email:       " USER@Example.COM ",
		DisplayName: "",
		Status:      "",
	}
	if err := validateAdminUserSaveInput(&input); err != nil {
		t.Fatal(err)
	}
	if input.Email != "user@example.com" || input.DisplayName != "user@example.com" || input.Status != "active" {
		t.Fatalf("unexpected normalized input: %+v", input)
	}
}

func TestNormalizeAdminRoleCodesDeduplicatesAndSorts(t *testing.T) {
	got := normalizeAdminRoleCodes([]string{"Billing Admin", "viewer", "billing-admin", "", "Super.Admin"})
	want := []string{"billing_admin", "super_admin", "viewer"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for idx := range want {
		if got[idx] != want[idx] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestAdminInvoiceReportCSVIncludesAccountingAmounts(t *testing.T) {
	records := adminInvoiceReportCSVRecords([]adminInvoiceReportRow{{
		invoiceSummary: invoiceSummary{
			InvoiceCode:           "INV-001",
			StudentCode:           "S001",
			StudentName:           "Nguyen Van A",
			ClassName:             "3.02",
			Grade:                 "3",
			SchoolYearCode:        "2025-2026",
			PeriodCode:            "2026-04",
			Month:                 4,
			IssuedAt:              time.Date(2026, 4, 1, 8, 0, 0, 0, time.UTC),
			DueDate:               "2026-04-15",
			Status:                invoiceStatusPartial,
			TotalAmount:           1000,
			PaidAmount:            750,
			CollectionBankBIN:     "970441",
			CollectionBankAccount: "123456789",
			QRBillNumber:          "INV-001",
		},
		OutstandingAmount: 250,
	}})

	if records[0][0] != "invoice_code" {
		t.Fatalf("expected csv header, got %v", records[0])
	}
	row := strings.Join(records[1], "|")
	for _, want := range []string{"INV-001", "S001", "2026-04", "partial", "1000", "750", "250", "970441"} {
		if !strings.Contains(row, want) {
			t.Fatalf("expected csv row to contain %q, got %s", want, row)
		}
	}
}

func TestFilterAdminReportTransactionsUsesFilteredInvoiceScope(t *testing.T) {
	invoices := []adminInvoiceReportRow{
		{invoiceSummary: invoiceSummary{ID: "invoice-1"}},
	}
	transactions := []paymentTransactionSummary{
		{ID: "txn-1", InvoiceID: "invoice-1"},
		{ID: "txn-2", InvoiceID: "invoice-2"},
		{ID: "txn-unmatched"},
	}

	filtered := filterAdminReportTransactions(invoices, transactions, adminFilters{PeriodCode: "2026-04"})

	if len(filtered) != 1 || filtered[0].ID != "txn-1" {
		t.Fatalf("expected only transactions matched to filtered invoices, got %+v", filtered)
	}
	all := filterAdminReportTransactions(invoices, transactions, adminFilters{})
	if len(all) != 3 {
		t.Fatalf("expected unconstrained filters to keep all transactions, got %+v", all)
	}
}
