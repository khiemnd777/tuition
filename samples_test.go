package main

import (
	"encoding/csv"
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
