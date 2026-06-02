package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestDemoPaymentSampleRowsAreImportableAndQRReady(t *testing.T) {
	file, err := os.Open("samples/demo_payments.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rows, err := parseCSVRows(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 12 {
		t.Fatalf("expected 12 demo payment rows, got %d", len(rows))
	}

	for _, row := range rows {
		if row.Amount <= 0 {
			t.Fatalf("expected fee columns to drive a positive amount for %s, got %d", row.StudentName, row.Amount)
		}
		if len(row.PaymentItems) == 0 {
			t.Fatalf("expected payment items for %s", row.StudentName)
		}
		if !strings.HasPrefix(row.BillNumber, "DEMO2504") {
			t.Fatalf("expected demo bill number for %s, got %q", row.StudentName, row.BillNumber)
		}

		item := buildQRItem(row, 128)
		if len(item.Errors) > 0 {
			t.Fatalf("expected QR-ready demo row for %s, got errors %+v", row.StudentName, item.Errors)
		}
		if item.VietQR == "" || item.QRData == "" {
			t.Fatalf("expected VietQR payload and PNG data for %s", row.StudentName)
		}
	}
}

func TestDemoMasterDataSampleRowsAreImportable(t *testing.T) {
	file, err := os.Open("samples/demo_master_data.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rows, err := parseMasterDataCSVRows(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 15 {
		t.Fatalf("expected 15 demo master data rows, got %d", len(rows))
	}
	if issues := validateMasterDataImportRows(rows); len(issues) > 0 {
		t.Fatalf("expected demo master data to be importable, got issues %+v", issues)
	}

	missingBillingReady := 0
	duplicateNameCount := 0
	for _, row := range rows {
		if !row.ParentActive || !row.ReceivesBillingEmail {
			missingBillingReady++
		}
		if row.StudentName == "Minh Nguyen An" {
			duplicateNameCount++
		}
	}
	if missingBillingReady < 3 {
		t.Fatalf("expected demo rows that can surface billing readiness warnings, got %d", missingBillingReady)
	}
	if duplicateNameCount < 2 {
		t.Fatalf("expected duplicate display-name demo case, got %d", duplicateNameCount)
	}
}

func TestDemoFeeAdjustmentsSampleRowsArePasteReady(t *testing.T) {
	file, err := os.Open("samples/demo_fee_adjustments.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 9 {
		t.Fatalf("expected header plus 8 demo adjustment rows, got %d records", len(records))
	}

	wantHeader := []string{"student_code", "adjustment_type", "fee_type_code", "amount", "reason"}
	for idx, want := range wantHeader {
		if records[0][idx] != want {
			t.Fatalf("expected adjustment header %q at column %d, got %q", want, idx, records[0][idx])
		}
	}

	allowedTypes := map[string]bool{
		adjustmentTypeDiscount:  true,
		adjustmentTypeSurcharge: true,
		adjustmentTypeWaiver:    true,
		adjustmentTypeCarryOver: true,
	}
	for idx, record := range records[1:] {
		rowNumber := idx + 2
		studentCode := strings.ToUpper(strings.TrimSpace(record[0]))
		adjustmentType := headerKey(record[1])
		feeTypeCode := headerKey(record[2])
		amount := parseAmount(record[3])
		reason := strings.TrimSpace(record[4])

		if !strings.HasPrefix(studentCode, "DEMO-S") {
			t.Fatalf("row %d: expected DEMO-S student code, got %q", rowNumber, studentCode)
		}
		if !allowedTypes[adjustmentType] {
			t.Fatalf("row %d: unexpected adjustment type %q", rowNumber, adjustmentType)
		}
		if reason == "" {
			t.Fatalf("row %d: expected reason", rowNumber)
		}
		if adjustmentType != adjustmentTypeWaiver && amount <= 0 {
			t.Fatalf("row %d: expected positive amount for %s, got %d", rowNumber, adjustmentType, amount)
		}
		if adjustmentType == adjustmentTypeWaiver && amount == 0 && feeTypeCode == "" {
			t.Fatalf("row %d: expected waiver fee type target", rowNumber)
		}
	}
}

func TestFinanceHubDemoMasterDataRowsAreImportable(t *testing.T) {
	file, err := os.Open("samples/finance_hub_demo/master_data.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rows, err := parseMasterDataCSVRows(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 17 {
		t.Fatalf("expected 17 finance hub master data rows, got %d", len(rows))
	}
	if issues := validateMasterDataImportRows(rows); len(issues) > 0 {
		t.Fatalf("expected finance hub master data to be importable, got issues %+v", issues)
	}

	students := map[string]bool{}
	missingBillingReady := 0
	duplicateNameCount := 0
	for _, row := range rows {
		students[row.StudentCode] = true
		if !row.ParentActive || !row.ReceivesBillingEmail {
			missingBillingReady++
		}
		if row.StudentName == "An Nguyen Minh" {
			duplicateNameCount++
		}
	}
	if len(students) != 16 {
		t.Fatalf("expected 16 distinct students, got %d", len(students))
	}
	if missingBillingReady < 4 {
		t.Fatalf("expected readiness-warning cases, got %d", missingBillingReady)
	}
	if duplicateNameCount < 2 {
		t.Fatalf("expected duplicate display-name cases, got %d", duplicateNameCount)
	}
}

func TestFinanceHubDemoLegacyPaymentsAreQRReady(t *testing.T) {
	file, err := os.Open("samples/finance_hub_demo/legacy_qr_payments.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rows, err := parseCSVRows(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 16 {
		t.Fatalf("expected 16 finance hub payment rows, got %d", len(rows))
	}

	for _, row := range rows {
		if row.Amount <= 0 {
			t.Fatalf("expected fee columns to drive a positive amount for %s, got %d", row.StudentName, row.Amount)
		}
		if len(row.PaymentItems) == 0 {
			t.Fatalf("expected payment items for %s", row.StudentName)
		}
		if !strings.HasPrefix(row.BillNumber, "FH2504") {
			t.Fatalf("expected finance hub bill number for %s, got %q", row.StudentName, row.BillNumber)
		}

		item := buildQRItem(row, 128)
		if len(item.Errors) > 0 {
			t.Fatalf("expected QR-ready finance hub row for %s, got errors %+v", row.StudentName, item.Errors)
		}
		if item.VietQR == "" || item.QRData == "" {
			t.Fatalf("expected VietQR payload and PNG data for %s", row.StudentName)
		}
	}
}

func TestFinanceHubDemoFeeProfilesAndCashReceiptsArePasteReady(t *testing.T) {
	profiles := readDemoCSV(t, "samples/finance_hub_demo/fee_schedule_profiles.csv")
	if len(profiles) < 2 {
		t.Fatal("expected fee profile rows")
	}
	wantProfileHeader := []string{"profile", "grade", "class_name", "fee_type_code", "label_vi", "label_en", "amount", "display_order"}
	for idx, want := range wantProfileHeader {
		if profiles[0][idx] != want {
			t.Fatalf("expected fee profile header %q at column %d, got %q", want, idx, profiles[0][idx])
		}
	}
	for idx, record := range profiles[1:] {
		rowNumber := idx + 2
		if record[0] == "" || record[2] == "" || record[3] == "" {
			t.Fatalf("row %d: expected profile, class, and fee type", rowNumber)
		}
		if amount := parseAmount(record[6]); amount <= 0 {
			t.Fatalf("row %d: expected positive fee amount, got %d", rowNumber, amount)
		}
	}

	receipts := readDemoCSV(t, "samples/finance_hub_demo/manual_cash_receipts.csv")
	if len(receipts) != 4 {
		t.Fatalf("expected header plus 3 cash receipt rows, got %d records", len(receipts))
	}
	for idx, record := range receipts[1:] {
		rowNumber := idx + 2
		if !strings.HasPrefix(record[1], "FH-S") {
			t.Fatalf("row %d: expected finance hub student code, got %q", rowNumber, record[1])
		}
		if amount := parseAmount(record[3]); amount <= 0 {
			t.Fatalf("row %d: expected positive receipt amount, got %d", rowNumber, amount)
		}
		if !strings.HasPrefix(record[5], "FH-CASH-") {
			t.Fatalf("row %d: expected demo cash receipt reference, got %q", rowNumber, record[5])
		}
	}
}

func TestFinanceHubDemoFeeAdjustmentsArePasteReady(t *testing.T) {
	records := readDemoCSV(t, "samples/finance_hub_demo/fee_adjustments.csv")
	if len(records) != 12 {
		t.Fatalf("expected header plus 11 finance hub adjustment rows, got %d records", len(records))
	}

	allowedTypes := map[string]bool{
		adjustmentTypeDiscount:  true,
		adjustmentTypeSurcharge: true,
		adjustmentTypeWaiver:    true,
		adjustmentTypeCarryOver: true,
	}
	for idx, record := range records[1:] {
		rowNumber := idx + 2
		adjustmentType := headerKey(record[1])
		amount := parseAmount(record[3])

		if !strings.HasPrefix(record[0], "FH-S") {
			t.Fatalf("row %d: expected finance hub student code, got %q", rowNumber, record[0])
		}
		if !allowedTypes[adjustmentType] {
			t.Fatalf("row %d: unexpected adjustment type %q", rowNumber, adjustmentType)
		}
		if strings.TrimSpace(record[4]) == "" {
			t.Fatalf("row %d: expected reason", rowNumber)
		}
		if adjustmentType != adjustmentTypeWaiver && amount <= 0 {
			t.Fatalf("row %d: expected positive adjustment amount, got %d", rowNumber, amount)
		}
		if adjustmentType == adjustmentTypeWaiver && amount == 0 && strings.TrimSpace(record[2]) == "" {
			t.Fatalf("row %d: expected waiver fee type target", rowNumber)
		}
	}
}

func TestFinanceHubDemoInvoiceNotificationAndWebhookTemplatesAreValid(t *testing.T) {
	var invoiceRequests struct {
		Preview  invoiceGenerateInput `json:"preview"`
		Generate invoiceGenerateInput `json:"generate"`
	}
	readDemoJSON(t, "samples/finance_hub_demo/invoice_generation_request.json", &invoiceRequests)
	if issues := validateInvoiceGenerateInput(invoiceRequests.Preview, true); len(issues) > 0 {
		t.Fatalf("expected invoice preview request template to validate, got %+v", issues)
	}
	if issues := validateInvoiceGenerateInput(invoiceRequests.Generate, true); len(issues) > 0 {
		t.Fatalf("expected invoice generate request template to validate, got %+v", issues)
	}

	var notificationFile struct {
		Campaigns    []notificationCampaignInput   `json:"campaigns"`
		EmailPreview notificationEmailPreviewInput `json:"emailPreview"`
	}
	readDemoJSON(t, "samples/finance_hub_demo/notification_campaigns.json", &notificationFile)
	if len(notificationFile.Campaigns) != 3 {
		t.Fatalf("expected 3 notification campaign templates, got %d", len(notificationFile.Campaigns))
	}
	for _, campaign := range notificationFile.Campaigns {
		normalized := normalizeNotificationCampaignInput(campaign)
		if issues := validateNotificationCampaignInput(normalized); len(issues) > 0 {
			t.Fatalf("expected notification campaign template to validate, got %+v", issues)
		}
		if !normalized.DryRun {
			t.Fatalf("expected demo notification campaign %q to default to dry-run", normalized.Name)
		}
	}
	preview := normalizeNotificationCampaignInput(notificationFile.EmailPreview.notificationCampaignInput)
	if preview.Name == "" || notificationFile.EmailPreview.RecipientID == "" {
		t.Fatalf("expected email-preview template to include a campaign name and recipient placeholder")
	}

	data, err := os.ReadFile("samples/finance_hub_demo/reconciliation_webhooks.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 webhook templates, got %d", len(lines))
	}
	for idx, line := range lines {
		var item struct {
			Case           string         `json:"case"`
			Provider       string         `json:"provider"`
			Path           string         `json:"path"`
			ExpectedStatus string         `json:"expectedStatus"`
			Payload        map[string]any `json:"payload"`
		}
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			t.Fatalf("line %d: invalid webhook json: %v", idx+1, err)
		}
		if item.Case == "" || item.Provider == "" || item.Path == "" || item.ExpectedStatus == "" {
			t.Fatalf("line %d: expected case metadata, got %+v", idx+1, item)
		}
		normalized, err := normalizeProviderWebhook(item.Provider, item.Payload)
		if err != nil {
			t.Fatalf("line %d: expected webhook payload to normalize: %v", idx+1, err)
		}
		if normalized.Amount <= 0 || normalized.Direction != paymentDirectionIn {
			t.Fatalf("line %d: expected inbound positive transaction, got %+v", idx+1, normalized)
		}
	}
}

func readDemoCSV(t *testing.T, path string) [][]string {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func readDemoJSON(t *testing.T, path string, target any) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatal(err)
	}
}
