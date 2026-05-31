package main

import (
	"strings"
	"testing"
)

func TestParseMasterDataCSVAllowsDuplicateStudentNamesWithDistinctCodes(t *testing.T) {
	body := strings.NewReader(`student_code,student_name,school_year,class_name,parent_name,parent_email
S001,Nguyen An,2025-2026,3.02,Nguyen Van A,a@example.com
S002,Nguyen An,2025-2026,3.03,Nguyen Van B,b@example.com
`)
	rows, err := parseMasterDataCSVRows(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected two rows, got %d", len(rows))
	}
	if rows[0].StudentCode != "S001" || rows[1].StudentCode != "S002" {
		t.Fatalf("expected normalized student codes, got %+v", rows)
	}
	if rows[0].Grade != "3" || rows[1].Grade != "3" {
		t.Fatalf("expected grade derived from class name, got %+v", rows)
	}
	if issues := validateMasterDataImportRows(rows); len(issues) > 0 {
		t.Fatalf("duplicate names with distinct codes must be valid, got %+v", issues)
	}
}

func TestValidateMasterDataImportRequiresStudentCode(t *testing.T) {
	body := strings.NewReader(`student_name,school_year,class_name,parent_name,parent_email
Nguyen An,2025-2026,3.02,Nguyen Van A,a@example.com
`)
	rows, err := parseMasterDataCSVRows(body)
	if err != nil {
		t.Fatal(err)
	}
	issues := validateMasterDataImportRows(rows)
	if !hasMasterDataIssue(issues, "missing_student_code") {
		t.Fatalf("expected missing student code issue, got %+v", issues)
	}
}

func TestParseMasterDataCSVRowsWithFieldMapping(t *testing.T) {
	body := strings.NewReader(`Mã HS,Họ và tên,Năm học,Lớp,Phụ huynh,Email
S001,Nguyen An,2025-2026,3.02,Nguyen Van A,a@example.com
`)
	rows, err := parseMasterDataCSVRowsWithMapping(body, map[string]string{
		"Mã HS":     "student_code",
		"Họ và tên": "student",
		"Năm học":   "school_year",
		"Lớp":       "class_name",
		"Phụ huynh": "parent",
		"Email":     "parent_email",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}
	row := rows[0]
	if row.StudentCode != "S001" || row.StudentName != "Nguyen An" || row.ClassName != "3.02" {
		t.Fatalf("unexpected mapped student fields: %+v", row)
	}
	if row.Grade != "3" || row.ParentName != "Nguyen Van A" || row.ParentEmail != "a@example.com" {
		t.Fatalf("unexpected mapped parent fields: %+v", row)
	}
}

func TestValidateMasterDataImportDetectsStudentCodeConflicts(t *testing.T) {
	body := strings.NewReader(`student_code,student_name,school_year,class_name,parent_name,parent_email
S001,Nguyen An,2025-2026,3.02,Nguyen Van A,a@example.com
S001,Tran Binh,2025-2026,3.02,Tran Van B,b@example.com
`)
	rows, err := parseMasterDataCSVRows(body)
	if err != nil {
		t.Fatal(err)
	}
	issues := validateMasterDataImportRows(rows)
	if !hasMasterDataIssue(issues, "student_code_conflict_in_csv") {
		t.Fatalf("expected student code conflict issue, got %+v", issues)
	}
}

func TestValidateMasterDataImportRejectsMultiplePrimaryParents(t *testing.T) {
	body := strings.NewReader(`student_code,student_name,school_year,class_name,parent_name,parent_email,parent_primary
S001,Nguyen An,2025-2026,3.02,Nguyen Van A,a@example.com,true
S001,Nguyen An,2025-2026,3.02,Nguyen Van B,b@example.com,true
`)
	rows, err := parseMasterDataCSVRows(body)
	if err != nil {
		t.Fatal(err)
	}
	issues := validateMasterDataImportRows(rows)
	if !hasMasterDataIssue(issues, "multiple_primary_parents_in_csv") {
		t.Fatalf("expected multiple primary parent issue, got %+v", issues)
	}
}

func TestNormalizeMasterDataStudentSaveInputDefaultsParentFlags(t *testing.T) {
	input := normalizeMasterDataStudentSaveInput(masterDataStudentSaveInput{
		StudentCode: " s001 ",
		StudentName: " Nguyen An ",
		ClassID:     " class-1 ",
		Parents: []masterDataParentSaveInput{
			{ParentName: " Nguyen Van A ", Email: " A@Example.COM ", Phone: " 090 123 4567 ", Relationship: " Mẹ "},
		},
	})

	if input.StudentCode != "S001" || input.StudentName != "Nguyen An" || input.ClassID != "class-1" {
		t.Fatalf("unexpected normalized student input: %+v", input)
	}
	if input.Status != "active" {
		t.Fatalf("expected default active status, got %q", input.Status)
	}
	parent := input.Parents[0]
	if parent.ParentName != "Nguyen Van A" || parent.Email != "a@example.com" || parent.Phone != "0901234567" || parent.Relationship != "me" {
		t.Fatalf("unexpected normalized parent input: %+v", parent)
	}
	if !boolValue(parent.EmailActive) || !boolValue(parent.IsActive) || !boolValue(parent.IsPrimary) || !boolValue(parent.ReceivesBillingEmail) {
		t.Fatalf("expected first parent defaults to active primary billing contact, got %+v", parent)
	}
}

func TestDeriveMasterDataStudentRelationshipState(t *testing.T) {
	student := masterDataStudent{
		Parents: []masterDataParentContact{
			{ParentName: "Nguyen Van A", Email: "a@example.com", Relationship: "mother", EmailActive: true, IsActive: true, ReceivesBillingEmail: true},
			{ParentName: "Nguyen Van B", Email: "", Relationship: "", EmailActive: true, IsActive: true, ReceivesBillingEmail: false},
		},
	}

	deriveMasterDataStudentRelationshipState(&student)

	if student.ParentCount != 2 || student.BillingRecipientCount != 1 || student.MissingBillingRecipient || student.ContactWarning != "" {
		t.Fatalf("unexpected relationship state: %+v", student)
	}
	if !student.Parents[0].BillingReady {
		t.Fatalf("expected first parent to be billing ready: %+v", student.Parents[0])
	}
	if student.Parents[1].Relationship != "guardian" {
		t.Fatalf("expected default guardian relationship, got %+v", student.Parents[1])
	}
}

func TestDeriveMasterDataStudentRelationshipStateWarnsWhenMissingBilling(t *testing.T) {
	student := masterDataStudent{
		Parents: []masterDataParentContact{
			{ParentName: "Nguyen Van A", Email: "a@example.com", EmailActive: false, IsActive: true, ReceivesBillingEmail: true},
		},
	}

	deriveMasterDataStudentRelationshipState(&student)

	if student.BillingRecipientCount != 0 || !student.MissingBillingRecipient || student.ContactWarning != "missing_billing_recipient" {
		t.Fatalf("expected missing billing warning, got %+v", student)
	}
}

func TestMasterDataUUIDPlaceholders(t *testing.T) {
	placeholders, args := masterDataUUIDPlaceholders([]string{"student-1", "student-2"})
	if placeholders != "$1::uuid, $2::uuid" {
		t.Fatalf("unexpected placeholders: %s", placeholders)
	}
	if len(args) != 2 || args[0] != "student-1" || args[1] != "student-2" {
		t.Fatalf("unexpected args: %+v", args)
	}
}

func TestValidateMasterDataStudentSaveInputRejectsMultipleActivePrimaryParents(t *testing.T) {
	input := normalizeMasterDataStudentSaveInput(masterDataStudentSaveInput{
		StudentCode: "S001",
		StudentName: "Nguyen An",
		ClassID:     "class-1",
		Parents: []masterDataParentSaveInput{
			{ParentName: "Nguyen Van A", Email: "a@example.com", IsPrimary: boolPtr(true)},
			{ParentName: "Nguyen Van B", Email: "b@example.com", IsPrimary: boolPtr(true)},
		},
	})

	if err := validateMasterDataStudentSaveInput(input); err == nil || !strings.Contains(err.Error(), "only one active primary parent") {
		t.Fatalf("expected multiple primary parent validation error, got %v", err)
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func hasMasterDataIssue(issues []masterDataImportIssue, issueType string) bool {
	for _, issue := range issues {
		if issue.Type == issueType {
			return true
		}
	}
	return false
}
