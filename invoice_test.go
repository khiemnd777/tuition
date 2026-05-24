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
			TotalAmount:           1500,
			CollectionBankBIN:     "970415",
			CollectionBankAccount: "0011001932418",
			QRBillNumber:          "SUNTEST001",
			QRNote:                "HP S001 2026-04",
		},
		BaseAmount:       1500,
		AdjustmentAmount: 0,
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
