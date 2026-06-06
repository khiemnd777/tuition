package main

import (
	"os"
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

func TestMatchPaymentTransactionToSubscriptionInvoicesMatchesByInvoiceCodeAndAmount(t *testing.T) {
	transaction := paymentTransactionSummary{
		ProviderCode:          paymentProviderSePay,
		ProviderTransactionID: "txn-1",
		Direction:             paymentDirectionIn,
		Amount:                250000,
		Description:           "Thanh toan SUB-SCHOOL_B-202606",
		ReferenceCode:         "SUB-SCHOOL_B-202606",
	}
	candidates := []subscriptionPaymentCandidate{
		{ID: "sub-1", InvoiceCode: "SUB-SCHOOL_B-202606", Amount: 250000},
		{ID: "sub-2", InvoiceCode: "SUB-SCHOOL_B-202607", Amount: 250000},
	}
	match, ok := matchPaymentTransactionToSubscriptionInvoices(transaction, candidates)
	if !ok {
		t.Fatal("expected subscription invoice match")
	}
	if match.ID != "sub-1" {
		t.Fatalf("expected sub-1 match, got %+v", match)
	}
}

func TestMatchPaymentTransactionToSubscriptionInvoicesRejectsAmountMismatch(t *testing.T) {
	transaction := paymentTransactionSummary{
		ProviderCode: paymentProviderSePay,
		Direction:    paymentDirectionIn,
		Amount:       150000,
		Description:  "SUB-SCHOOL_B-202606",
	}
	candidates := []subscriptionPaymentCandidate{
		{ID: "sub-1", InvoiceCode: "SUB-SCHOOL_B-202606", Amount: 250000},
	}
	if _, ok := matchPaymentTransactionToSubscriptionInvoices(transaction, candidates); ok {
		t.Fatal("expected amount mismatch to reject subscription auto-confirm")
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

func TestResolvePaymentWebhookTenant(t *testing.T) {
	cases := []struct {
		name         string
		path         string
		wantTenant   string
		wantProvider string
		wantErr      bool
	}{
		{name: "legacy provider path", path: "/api/v1/payments/webhooks/sepay", wantTenant: "", wantProvider: "sepay"},
		{name: "tenant-scoped provider path", path: "/api/v1/payments/webhooks/abc_sun/sepay", wantTenant: "ABC_SUN", wantProvider: "sepay"},
		{name: "tenant path supports mixed case", path: "/api/v1/payments/webhooks/AbC_SuN/payos", wantTenant: "ABC_SUN", wantProvider: "payos"},
		{name: "missing provider should fail", path: "/api/v1/payments/webhooks/", wantErr: true},
		{name: "invalid extra segment should fail", path: "/api/v1/payments/webhooks/abc_sun/sepay/extra", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tenantCode, providerCode, err := resolvePaymentWebhookTenant(tc.path)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for path %q", tc.path)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for path %q: %v", tc.path, err)
			}
			if tenantCode != tc.wantTenant {
				t.Fatalf("expected tenantCode=%q, got=%q", tc.wantTenant, tenantCode)
			}
			if providerCode != tc.wantProvider {
				t.Fatalf("expected provider=%q, got=%q", tc.wantProvider, providerCode)
			}
		})
	}
}

func TestLoadPayOSConfigPrefersTenantProviderConfig(t *testing.T) {
	t.Setenv("ABC_PAYOS_CLIENT_ID", "env-client")
	t.Setenv("ABC_PAYOS_API_KEY", "env-api")
	t.Setenv("ABC_PAYOS_CHECKSUM_KEY", "env-checksum")
	t.Setenv("ABC_PAYOS_RETURN_URL", "https://env-return.example")
	t.Setenv("ABC_PAYOS_CANCEL_URL", "https://env-cancel.example")
	t.Setenv("ABC_PAYOS_API_BASE_URL", "https://env-api-base.example")

	cfg := loadPayOSConfig(paymentProvider{
		Code: paymentProviderPayOS,
		config: map[string]any{
			"clientId":    "tenant-client",
			"apiKey":      "tenant-api",
			"checksumKey": "tenant-checksum",
			"returnUrl":   "https://tenant-return.example",
			"cancelUrl":   "https://tenant-cancel.example",
			"apiBaseUrl":  "https://tenant-api-base.example",
		},
	})

	if cfg.ClientID != "tenant-client" || cfg.APIKey != "tenant-api" || cfg.ChecksumKey != "tenant-checksum" {
		t.Fatalf("expected tenant payOS credentials, got %+v", cfg)
	}
	if cfg.ReturnURL != "https://tenant-return.example" || cfg.CancelURL != "https://tenant-cancel.example" || cfg.APIBaseURL != "https://tenant-api-base.example" {
		t.Fatalf("expected tenant payOS URLs, got %+v", cfg)
	}
}

func TestLoadPayOSConfigFallsBackToEnvForMissingTenantFields(t *testing.T) {
	for key, value := range map[string]string{
		"ABC_PAYOS_CLIENT_ID":    "env-client",
		"ABC_PAYOS_API_KEY":      "env-api",
		"ABC_PAYOS_CHECKSUM_KEY": "env-checksum",
		"ABC_PAYOS_RETURN_URL":   "https://env-return.example",
		"ABC_PAYOS_CANCEL_URL":   "https://env-cancel.example",
	} {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		for _, key := range []string{"ABC_PAYOS_CLIENT_ID", "ABC_PAYOS_API_KEY", "ABC_PAYOS_CHECKSUM_KEY", "ABC_PAYOS_RETURN_URL", "ABC_PAYOS_CANCEL_URL"} {
			_ = os.Unsetenv(key)
		}
	})

	cfg := loadPayOSConfig(paymentProvider{
		Code:   paymentProviderPayOS,
		config: map[string]any{"clientId": "tenant-client"},
	})

	if cfg.ClientID != "tenant-client" {
		t.Fatalf("expected tenant client id, got %+v", cfg)
	}
	if cfg.APIKey != "env-api" || cfg.ChecksumKey != "env-checksum" || cfg.ReturnURL != "https://env-return.example" || cfg.CancelURL != "https://env-cancel.example" {
		t.Fatalf("expected env fallback for missing fields, got %+v", cfg)
	}
}
