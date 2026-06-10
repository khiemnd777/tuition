package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSamplesDirectoryDoesNotCarryCSVFixtures(t *testing.T) {
	err := filepath.WalkDir("samples", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".csv") {
			t.Fatalf("unexpected CSV sample fixture %s; use seeded DB data or Excel import workbooks", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFinanceHubDemoSeedMasterDataRowsAreImportable(t *testing.T) {
	rows := financeHubDemoMasterDataRows()
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
		if row.SchoolCode != demoSchoolCode {
			t.Fatalf("expected seeded row %s to target %s, got %s", row.StudentCode, demoSchoolCode, row.SchoolCode)
		}
		if row.ParentEmail != "" && !hasExampleDomain(row.ParentEmail) {
			t.Fatalf("expected fictional parent email for %s, got %q", row.StudentCode, row.ParentEmail)
		}
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

func TestFinanceHubDemoSeedFeeProfilesAndAdjustmentsAreReady(t *testing.T) {
	profiles := demoFeeProfiles()
	if len(profiles) != 20 {
		t.Fatalf("expected 20 finance hub fee profile rows, got %d", len(profiles))
	}
	for idx, row := range profiles {
		if row.Profile == "" || row.ClassName == "" || row.FeeTypeCode == "" {
			t.Fatalf("profile row %d: expected profile, class, and fee type: %+v", idx+1, row)
		}
		if row.Amount <= 0 {
			t.Fatalf("profile row %d: expected positive fee amount, got %d", idx+1, row.Amount)
		}
	}

	adjustments := demoFeeAdjustments()
	if len(adjustments) != 11 {
		t.Fatalf("expected 11 finance hub adjustment rows, got %d", len(adjustments))
	}
	allowedTypes := map[string]bool{
		adjustmentTypeDiscount:  true,
		adjustmentTypeSurcharge: true,
		adjustmentTypeWaiver:    true,
		adjustmentTypeCarryOver: true,
	}
	for idx, adjustment := range adjustments {
		if !strings.HasPrefix(adjustment.StudentCode, "FH-S") {
			t.Fatalf("adjustment row %d: expected finance hub student code, got %q", idx+1, adjustment.StudentCode)
		}
		if !allowedTypes[adjustment.AdjustmentType] {
			t.Fatalf("adjustment row %d: unexpected adjustment type %q", idx+1, adjustment.AdjustmentType)
		}
		if adjustment.Reason == "" {
			t.Fatalf("adjustment row %d: expected reason", idx+1)
		}
		if adjustment.AdjustmentType != adjustmentTypeWaiver && adjustment.Amount <= 0 {
			t.Fatalf("adjustment row %d: expected positive amount, got %d", idx+1, adjustment.Amount)
		}
		if adjustment.AdjustmentType == adjustmentTypeWaiver && adjustment.Amount == 0 && adjustment.FeeTypeCode == "" {
			t.Fatalf("adjustment row %d: expected waiver fee type target", idx+1)
		}
	}
}

func TestFinanceHubDemoImportMoreStudentsWorkbookIsImportable(t *testing.T) {
	data, err := os.ReadFile("samples/finance_hub_demo/import_more_students.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	table, err := parseImportTableBytes(data, "import_more_students.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := parseMasterDataRows(table, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 6 {
		t.Fatalf("expected 6 additional student rows, got %d", len(rows))
	}
	if issues := validateMasterDataImportRows(rows); len(issues) > 0 {
		t.Fatalf("expected additional student workbook to be importable, got issues %+v", issues)
	}
	for _, row := range rows {
		if !strings.HasPrefix(row.StudentCode, "FH-N") {
			t.Fatalf("expected additional student code, got %q", row.StudentCode)
		}
		if row.SchoolCode != demoSchoolCode {
			t.Fatalf("expected workbook row %s to target %s, got %s", row.StudentCode, demoSchoolCode, row.SchoolCode)
		}
	}
}

func TestFinanceHubDemoImportMorePaymentsWorkbookIsQRReady(t *testing.T) {
	data, err := os.ReadFile("samples/finance_hub_demo/import_more_payments.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	table, err := parseImportTableBytes(data, "import_more_payments.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := parsePaymentRows(table, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("expected 4 additional payment rows, got %d", len(rows))
	}
	for _, row := range rows {
		if row.Amount <= 0 {
			t.Fatalf("expected fee columns to drive a positive amount for %s, got %d", row.StudentName, row.Amount)
		}
		if len(row.PaymentItems) == 0 {
			t.Fatalf("expected payment items for %s", row.StudentName)
		}
		if !strings.HasPrefix(row.BillNumber, "FHIMP") {
			t.Fatalf("expected additional import bill number for %s, got %q", row.StudentName, row.BillNumber)
		}

		item := buildQRItem(row, 128)
		if len(item.Errors) > 0 {
			t.Fatalf("expected QR-ready additional payment row for %s, got errors %+v", row.StudentName, item.Errors)
		}
		if item.VietQR == "" || item.QRData == "" {
			t.Fatalf("expected VietQR payload and PNG data for %s", row.StudentName)
		}
	}
}

func TestFinanceHubDemoNotificationAndWebhookTemplatesAreValid(t *testing.T) {
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
