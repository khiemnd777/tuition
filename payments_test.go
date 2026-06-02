package main

import (
	"testing"
)

func TestMatchPaymentTransactionToInvoicesMatchesByReferenceAccountAndAmount(t *testing.T) {
	transaction := paymentTransactionSummary{
		ProviderCode:          "sepay",
		ProviderTransactionID: "92704",
		Direction:             paymentDirectionIn,
		Amount:                1200,
		AccountNumber:         "0011001932418",
		Description:           "SUNABC123 chuyen tien hoc phi",
		ReferenceCode:         "FT24012345678",
	}
	candidates := []paymentInvoiceCandidate{
		{ID: "invoice-1", InvoiceCode: "SUNABC123", QRBillNumber: "SUNABC123", CollectionBankAccount: "0011001932418", TotalAmount: 1200, PaidAmount: 0},
		{ID: "invoice-2", InvoiceCode: "SUNOTHER", QRBillNumber: "SUNOTHER", CollectionBankAccount: "0011001932418", TotalAmount: 1200, PaidAmount: 0},
	}

	match, ok := matchPaymentTransactionToInvoices(transaction, candidates)
	if !ok {
		t.Fatal("expected transaction to match an invoice")
	}
	if match.Invoice.ID != "invoice-1" {
		t.Fatalf("expected invoice-1 match, got %+v", match.Invoice)
	}
	if match.Score < 100 {
		t.Fatalf("expected strong match score, got %d", match.Score)
	}
	if match.AmountApplied != transaction.Amount {
		t.Fatalf("expected amount applied %d, got %d", transaction.Amount, match.AmountApplied)
	}
}

func TestInvoiceStatusBecamePaidOnlyOnPaidTransition(t *testing.T) {
	cases := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{name: "unpaid to paid", old: invoiceStatusUnpaid, new: invoiceStatusPaid, want: true},
		{name: "partial to paid", old: invoiceStatusPartial, new: invoiceStatusPaid, want: true},
		{name: "paid unchanged", old: invoiceStatusPaid, new: invoiceStatusPaid, want: false},
		{name: "paid to overpaid", old: invoiceStatusPaid, new: invoiceStatusOverpaid, want: false},
		{name: "partial to overpaid", old: invoiceStatusPartial, new: invoiceStatusOverpaid, want: false},
		{name: "unpaid to manual review", old: invoiceStatusUnpaid, new: invoiceStatusManualReview, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := invoiceStatusBecamePaid(tc.old, tc.new); got != tc.want {
				t.Fatalf("invoiceStatusBecamePaid(%q, %q) = %v, want %v", tc.old, tc.new, got, tc.want)
			}
		})
	}
}

func TestMatchPaymentTransactionToInvoicesSupportsPartialAndOverpayment(t *testing.T) {
	candidate := paymentInvoiceCandidate{
		ID:                    "invoice-1",
		InvoiceCode:           "SUNABC123",
		QRBillNumber:          "SUNABC123",
		CollectionBankAccount: "0011001932418",
		TotalAmount:           1200,
		PaidAmount:            0,
	}
	for _, tc := range []struct {
		name   string
		amount int
		want   string
	}{
		{name: "partial", amount: 500, want: "partial_amount"},
		{name: "overpaid", amount: 1500, want: "overpayment_amount"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			match, ok := matchPaymentTransactionToInvoices(paymentTransactionSummary{
				Direction:     paymentDirectionIn,
				Amount:        tc.amount,
				AccountNumber: "0011001932418",
				Description:   "Thanh toan SUNABC123",
			}, []paymentInvoiceCandidate{candidate})
			if !ok {
				t.Fatalf("expected %s payment to match", tc.name)
			}
			if !containsPaymentReference(normalizedPaymentReference(match.Reason), tc.want) {
				t.Fatalf("expected reason to include %q, got %q", tc.want, match.Reason)
			}
		})
	}
}

func TestSummarizePaymentReconciliationTracksCollectionQueues(t *testing.T) {
	invoices := []invoiceSummary{
		{ID: "invoice-1", Status: invoiceStatusUnpaid, TotalAmount: 1000, PaidAmount: 0},
		{ID: "invoice-2", Status: invoiceStatusPartial, TotalAmount: 2000, PaidAmount: 750},
		{ID: "invoice-3", Status: invoiceStatusPaid, TotalAmount: 3000, PaidAmount: 3000},
		{ID: "invoice-4", Status: invoiceStatusOverpaid, TotalAmount: 4000, PaidAmount: 4500},
		{ID: "invoice-5", Status: invoiceStatusManualReview, TotalAmount: 5000, PaidAmount: 0},
	}
	transactions := []paymentTransactionSummary{
		{Status: paymentTransactionStatusMatched},
		{Status: paymentTransactionStatusUnmatched},
		{Status: paymentTransactionStatusManualReview},
	}

	summary := summarizePaymentReconciliation(invoices, transactions)
	if summary.InvoiceCount != 5 {
		t.Fatalf("expected 5 invoices, got %d", summary.InvoiceCount)
	}
	if summary.TotalReceivable != 15000 {
		t.Fatalf("expected receivable 15000, got %d", summary.TotalReceivable)
	}
	if summary.TotalCollected != 8250 {
		t.Fatalf("expected collected 8250, got %d", summary.TotalCollected)
	}
	if summary.OutstandingAmount != 7250 {
		t.Fatalf("expected outstanding 7250, got %d", summary.OutstandingAmount)
	}
	if summary.CollectionRate != 0.55 {
		t.Fatalf("expected collection rate 0.55, got %f", summary.CollectionRate)
	}
	if summary.UnpaidCount != 1 || summary.PartialCount != 1 || summary.PaidCount != 1 || summary.OverpaidCount != 1 {
		t.Fatalf("unexpected invoice status counts: %+v", summary)
	}
	if summary.UnmatchedCount != 1 || summary.MatchedCount != 1 || summary.ManualReviewCount != 2 {
		t.Fatalf("unexpected work queue counts: %+v", summary)
	}
}

func TestInvoiceIDsFromSummariesKeepsStableInvoiceScope(t *testing.T) {
	ids := invoiceIDsFromSummaries([]invoiceSummary{
		{ID: "invoice-1"},
		{},
		{ID: " invoice-2 "},
	})
	if len(ids) != 2 {
		t.Fatalf("expected two invoice IDs, got %+v", ids)
	}
	if ids[0] != "invoice-1" || ids[1] != "invoice-2" {
		t.Fatalf("expected IDs to preserve list order for query scope, got %+v", ids)
	}
}

func TestNormalizeSePayWebhookTransaction(t *testing.T) {
	transaction, err := normalizeSePayWebhook(map[string]any{
		"id":              float64(92704),
		"gateway":         "Vietcombank",
		"transactionDate": "2024-07-02 11:08:33",
		"accountNumber":   "1017588888",
		"code":            "SEVN63DC8E5C",
		"content":         "SUNABC123 chuyen tien",
		"transferType":    "in",
		"transferAmount":  float64(5000000),
		"referenceCode":   "FT24012345678",
	})
	if err != nil {
		t.Fatal(err)
	}
	if transaction.ProviderTransactionID != "92704" {
		t.Fatalf("unexpected transaction id %q", transaction.ProviderTransactionID)
	}
	if transaction.Amount != 5000000 {
		t.Fatalf("unexpected amount %d", transaction.Amount)
	}
	if transaction.AccountNumber != "1017588888" {
		t.Fatalf("unexpected account %q", transaction.AccountNumber)
	}
	if transaction.ReferenceCode != "FT24012345678" {
		t.Fatalf("unexpected reference %q", transaction.ReferenceCode)
	}
}

func TestPayOSSignatureUsesAlphabeticalPaymentRequestFields(t *testing.T) {
	signature := payOSSignature(map[string]any{
		"returnUrl":   "https://return.example",
		"orderCode":   123,
		"description": "SUN123",
		"cancelUrl":   "https://cancel.example",
		"amount":      3000,
	}, "secret")
	const want = "abe371109d0ccbe17ba782d9e1c08189fa543a127b61192313285c5d783f830f"
	if signature != want {
		t.Fatalf("unexpected payOS signature %q, want %q", signature, want)
	}
}
