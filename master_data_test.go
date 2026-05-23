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

func hasMasterDataIssue(issues []masterDataImportIssue, issueType string) bool {
	for _, issue := range issues {
		if issue.Type == issueType {
			return true
		}
	}
	return false
}
