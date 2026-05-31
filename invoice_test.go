package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestStableInvoiceCodeIsBillNumberSafe(t *testing.T) {
	first := stableInvoiceCode("schedule-1", "student-1", "2026-04", "S001")
	second := stableInvoiceCode("schedule-1", "student-1", "2026-04", "S001")
	if first != second {
		t.Fatalf("expected stable invoice code, got %q and %q", first, second)
	}
	if len(first) > 25 {
		t.Fatalf("expected invoice code to fit VietQR BillNumber, got %q length %d", first, len(first))
	}
	if !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(first) {
		t.Fatalf("expected ANS-safe uppercase invoice code, got %q", first)
	}
	if cleanANS(first, 25) != first {
		t.Fatalf("expected invoice code to survive cleanANS, got %q", cleanANS(first, 25))
	}
}

func TestBuildInvoicePreviewMapsInvoiceCodeToPaymentBillNumber(t *testing.T) {
	meta := invoiceScheduleMeta{
		ID:             "schedule-1",
		SchoolYearID:   "year-1",
		SchoolYearCode: "2025-2026",
		PeriodCode:     "2026-04",
		Month:          4,
	}
	feePreview := feeSchedulePreview{
		Rows: []feeSchedulePreviewRow{{
			StudentID:        "student-1",
			StudentCode:      "S001",
			StudentName:      "Nguyen An",
			ClassID:          "class-1",
			ClassName:        "3.02",
			Grade:            "3",
			SchoolYearCode:   "2025-2026",
			BaseAmount:       1500,
			AdjustmentAmount: -300,
			TotalAmount:      1200,
			Items: []feeSchedulePreviewItem{
				{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
				{FeeTypeCode: "shuttle", LabelVI: "Phi xe", LabelEN: "Shuttle", Amount: 500},
			},
			Adjustments: []feeSchedulePreviewAdjustment{
				{AdjustmentType: "discount", FeeTypeCode: "tuition", LabelVI: "Giam tru", LabelEN: "Discount", Amount: 300, Delta: -300, Reason: "Sibling"},
			},
		}},
	}
	input := invoiceGenerateInput{BankBIN: "970415", BankAccount: "0011001932418", IssueDate: "2026-05-24"}

	preview := buildInvoicePreview(meta, input, feePreview, nil, time.Date(2026, 5, 24, 9, 0, 0, 0, time.UTC))
	if len(preview.Rows) != 1 {
		t.Fatalf("expected one invoice row, got %d", len(preview.Rows))
	}
	row := preview.Rows[0]
	if row.InvoiceCode == "" || row.QRBillNumber != row.InvoiceCode {
		t.Fatalf("expected invoice code to map to QR bill number, got row %+v", row)
	}
	payment := invoicePaymentRow(row)
	if payment.BillNumber != row.InvoiceCode {
		t.Fatalf("expected payment bill number %q, got %q", row.InvoiceCode, payment.BillNumber)
	}
	if got := paymentItemsTotal(payment.PaymentItems); got != row.TotalAmount {
		t.Fatalf("expected payment items total %d, got %d", row.TotalAmount, got)
	}
}

func TestBuildInvoicePreviewMatchesFeePreviewTotals(t *testing.T) {
	meta := invoiceScheduleMeta{
		ID:             "schedule-1",
		SchoolYearID:   "year-1",
		SchoolYearCode: "2025-2026",
		PeriodCode:     "2026-04",
		Month:          4,
	}
	feePreview := feeSchedulePreview{
		Summary: feeSchedulePreviewSummary{
			StudentCount: 2,
			BaseAmount:   3000,
			Adjustments:  -350,
			TotalAmount:  2650,
		},
		Rows: []feeSchedulePreviewRow{
			{
				StudentID:        "student-1",
				StudentCode:      "S001",
				StudentName:      "Nguyen An",
				ClassID:          "class-1",
				ClassName:        "3.02",
				Grade:            "3",
				BaseAmount:       1500,
				AdjustmentAmount: -350,
				TotalAmount:      1150,
				Items: []feeSchedulePreviewItem{
					{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
					{FeeTypeCode: "shuttle", LabelVI: "Phi xe", LabelEN: "Shuttle", Amount: 500},
				},
				Adjustments: []feeSchedulePreviewAdjustment{
					{AdjustmentType: "discount", FeeTypeCode: "tuition", LabelVI: "Giam tru", LabelEN: "Discount", Amount: 350, Delta: -350, Reason: "Sibling"},
				},
			},
			{
				StudentID:        "student-2",
				StudentCode:      "S002",
				StudentName:      "Tran Binh",
				ClassID:          "class-1",
				ClassName:        "3.02",
				Grade:            "3",
				BaseAmount:       1500,
				AdjustmentAmount: 0,
				TotalAmount:      1500,
				Items: []feeSchedulePreviewItem{
					{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
					{FeeTypeCode: "shuttle", LabelVI: "Phi xe", LabelEN: "Shuttle", Amount: 500},
				},
			},
		},
	}

	preview := buildInvoicePreview(meta, invoiceGenerateInput{}, feePreview, nil, time.Date(2026, 5, 31, 9, 0, 0, 0, time.UTC))
	if preview.Summary.StudentCount != feePreview.Summary.StudentCount {
		t.Fatalf("expected student count %d, got %d", feePreview.Summary.StudentCount, preview.Summary.StudentCount)
	}
	if preview.Summary.BaseAmount != feePreview.Summary.BaseAmount {
		t.Fatalf("expected base total %d, got %d", feePreview.Summary.BaseAmount, preview.Summary.BaseAmount)
	}
	if preview.Summary.AdjustmentAmount != feePreview.Summary.Adjustments {
		t.Fatalf("expected adjustment total %d, got %d", feePreview.Summary.Adjustments, preview.Summary.AdjustmentAmount)
	}
	if preview.Summary.TotalAmount != feePreview.Summary.TotalAmount {
		t.Fatalf("expected invoice total %d, got %d", feePreview.Summary.TotalAmount, preview.Summary.TotalAmount)
	}
	for idx, row := range preview.Rows {
		if row.TotalAmount != feePreview.Rows[idx].TotalAmount {
			t.Fatalf("expected row %d total %d, got %d", idx, feePreview.Rows[idx].TotalAmount, row.TotalAmount)
		}
		if got := paymentItemsTotal(row.PaymentItems); got != row.TotalAmount {
			t.Fatalf("expected row %d payment items total %d, got %d", idx, row.TotalAmount, got)
		}
	}
}

func TestBuildInvoicePreviewClassifiesIdempotencyState(t *testing.T) {
	meta := invoiceScheduleMeta{
		ID:             "schedule-1",
		SchoolYearID:   "year-1",
		SchoolYearCode: "2025-2026",
		PeriodCode:     "2026-04",
	}
	feePreview := feeSchedulePreview{
		Rows: []feeSchedulePreviewRow{
			{
				StudentID:    "student-1",
				StudentCode:  "S001",
				StudentName:  "Nguyen An",
				ClassID:      "class-1",
				ClassName:    "3.02",
				BillingReady: true,
				BaseAmount:   1000,
				TotalAmount:  1000,
				Items: []feeSchedulePreviewItem{
					{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
				},
			},
			{
				StudentID:    "student-2",
				StudentCode:  "S002",
				StudentName:  "Tran Binh",
				ClassID:      "class-1",
				ClassName:    "3.02",
				BillingReady: true,
				BaseAmount:   1000,
				TotalAmount:  1000,
				Items: []feeSchedulePreviewItem{
					{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
				},
			},
		},
	}
	existing := map[string]invoiceExistingRef{
		"student-1": {ID: "invoice-1", InvoiceCode: "SUN001", Status: invoiceStatusUnpaid, PaidAmount: 0},
		"student-2": {ID: "invoice-2", InvoiceCode: "SUN002", Status: invoiceStatusPartial, PaidAmount: 500},
	}

	preview := buildInvoicePreview(meta, invoiceGenerateInput{Regenerate: true}, feePreview, existing, time.Date(2026, 5, 31, 9, 0, 0, 0, time.UTC))
	if len(preview.Rows) != 2 {
		t.Fatalf("expected two preview rows, got %d", len(preview.Rows))
	}
	if preview.Rows[0].GenerationState != "ready_to_regenerate" || !preview.Rows[0].Regenerable || preview.Rows[0].BlocksGeneration {
		t.Fatalf("expected first row to be regenerable, got %+v", preview.Rows[0])
	}
	if preview.Rows[1].GenerationState != "blocked_paid_regenerate" || !preview.Rows[1].BlocksGeneration {
		t.Fatalf("expected second row to block paid regeneration, got %+v", preview.Rows[1])
	}
	if preview.Summary.RegenerableCount != 1 || preview.Summary.BlockedCount != 1 {
		t.Fatalf("expected one regenerable and one blocked row, got %+v", preview.Summary)
	}
	if !hasInvoicePreviewIssue(preview.Issues, "cannot_regenerate_paid_invoice") {
		t.Fatalf("expected paid regeneration issue, got %+v", preview.Issues)
	}
}

func TestDeriveInvoiceStatusFromPaymentTotals(t *testing.T) {
	cases := []struct {
		total int
		paid  int
		want  string
	}{
		{total: 1000, paid: 0, want: "unpaid"},
		{total: 1000, paid: 500, want: "partial"},
		{total: 1000, paid: 1000, want: "paid"},
		{total: 1000, paid: 1200, want: "overpaid"},
	}
	for _, tc := range cases {
		if got := deriveInvoiceStatus(tc.total, tc.paid); got != tc.want {
			t.Fatalf("deriveInvoiceStatus(%d, %d) = %q, want %q", tc.total, tc.paid, got, tc.want)
		}
	}
}

func hasInvoicePreviewIssue(issues []invoicePreviewIssue, issueType string) bool {
	for _, issue := range issues {
		if issue.Type == issueType {
			return true
		}
	}
	return false
}

func TestRenderInvoicePDFContainsInvoiceDataAndQR(t *testing.T) {
	invoice := invoiceDocument{
		invoiceSummary: invoiceSummary{
			ID:                    "invoice-1",
			InvoiceCode:           "SUNTEST001",
			StudentCode:           "S001",
			StudentName:           "Nguyen An",
			ClassName:             "3.02",
			SchoolYearCode:        "2025-2026",
			PeriodCode:            "2026-04",
			IssuedAt:              time.Date(2026, 5, 24, 9, 0, 0, 0, time.UTC),
			Status:                "unpaid",
			BaseAmount:            1500,
			AdjustmentAmount:      0,
			TotalAmount:           1500,
			CollectionBankBIN:     "970415",
			CollectionBankAccount: "0011001932418",
			QRBillNumber:          "SUNTEST001",
			QRNote:                "HP S001 2026-04",
		},
		Items: []invoiceDocumentItem{
			{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1500, DisplayOrder: 1},
		},
	}
	qr := buildQRItem(paymentRowFromInvoice(invoice), 256)
	if len(qr.Errors) > 0 {
		t.Fatalf("unexpected QR errors: %v", qr.Errors)
	}

	pdf, err := renderInvoicePDF(invoice, qr)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range [][]byte{
		[]byte("%PDF-1.4"),
		[]byte("ABC SUN Invoice Receipt"),
		[]byte("SUNTEST001"),
		[]byte("Nguyen An"),
		[]byte("VietQR BillNumber"),
		[]byte(" re f"),
	} {
		if !bytes.Contains(pdf, want) {
			t.Fatalf("expected PDF to contain %q", string(want))
		}
	}
	if !strings.Contains(string(pdf), "Total due") {
		t.Fatalf("expected PDF total section, got %s", string(pdf[:min(len(pdf), 400)]))
	}
}
