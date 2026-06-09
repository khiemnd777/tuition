package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildSchoolTreeGroupsClassesByYearAndGrade(t *testing.T) {
	tree := buildSchoolTree(
		[]schoolTreeSchool{{ID: "school-1", Code: "DEKISUGI", Name: "DEKISUGI", Status: "active"}},
		[]schoolTreeSchoolYear{
			{ID: "year-1", SchoolID: "school-1", SchoolCode: "DEKISUGI", Code: "2025-2026", Name: "2025-2026", Status: "active", ClassCount: 2, StudentCount: 33, FeeScheduleCount: 4, AdjustmentCount: 3},
			{ID: "year-2", SchoolID: "school-1", SchoolCode: "DEKISUGI", Code: "2024-2025", Name: "2024-2025", Status: "archived", ClassCount: 1, StudentCount: 10, FeeScheduleCount: 1, AdjustmentCount: 0},
		},
		[]schoolTreeClass{
			{ID: "class-1", SchoolID: "school-1", SchoolCode: "DEKISUGI", SchoolYearID: "year-1", SchoolYearCode: "2025-2026", Grade: "1", Name: "1A", Status: "active", StudentCount: 18, FeeScheduleCount: 1, ActiveFeeScheduleCount: 1, AdjustmentCount: 2, BillingReadyStudentCount: 17, MissingBillingRecipientCount: 1, CurrentFeeScheduleCount: 1, CurrentActiveScheduleCount: 1, CurrentInvoiceCount: 18, OpenInvoiceCount: 3},
			{ID: "class-2", SchoolID: "school-1", SchoolCode: "DEKISUGI", SchoolYearID: "year-1", SchoolYearCode: "2025-2026", Grade: "1", Name: "1B", Status: "active", StudentCount: 15, FeeScheduleCount: 1, ActiveFeeScheduleCount: 0, AdjustmentCount: 1, BillingReadyStudentCount: 15, MissingBillingRecipientCount: 0, CurrentFeeScheduleCount: 0, CurrentActiveScheduleCount: 0, CurrentInvoiceCount: 0, OpenInvoiceCount: 0},
			{ID: "class-3", SchoolID: "school-1", SchoolCode: "DEKISUGI", SchoolYearID: "year-2", SchoolYearCode: "2024-2025", Grade: "2", Name: "2A", Status: "active", StudentCount: 10, FeeScheduleCount: 1, BillingReadyStudentCount: 8, MissingBillingRecipientCount: 2, CurrentFeeScheduleCount: 1, CurrentActiveScheduleCount: 0, CurrentInvoiceCount: 9, OpenInvoiceCount: 2},
		},
	)

	if len(tree) != 1 {
		t.Fatalf("expected one school, got %d", len(tree))
	}
	if got := tree[0].StudentCount; got != 43 {
		t.Fatalf("school student count = %d, want 43", got)
	}
	if tree[0].BillingReadyStudentCount != 40 || tree[0].MissingBillingRecipientCount != 3 || tree[0].CurrentInvoiceCount != 27 || tree[0].OpenInvoiceCount != 5 || tree[0].IssueCount != 21 {
		t.Fatalf("unexpected school readiness aggregates: %+v", tree[0])
	}
	if len(tree[0].SchoolYears) != 2 {
		t.Fatalf("expected two school years, got %d", len(tree[0].SchoolYears))
	}
	year := tree[0].SchoolYears[0]
	if year.Code != "2025-2026" {
		t.Fatalf("expected newest year first, got %s", year.Code)
	}
	if len(year.Grades) != 1 || year.Grades[0].Grade != "1" {
		t.Fatalf("expected grade 1 grouping, got %+v", year.Grades)
	}
	if year.BillingReadyStudentCount != 32 || year.MissingBillingRecipientCount != 1 || year.CurrentFeeScheduleCount != 1 || year.CurrentActiveScheduleCount != 1 || year.CurrentInvoiceCount != 18 || year.OpenInvoiceCount != 3 || year.IssueCount != 17 {
		t.Fatalf("unexpected year readiness aggregates: %+v", year)
	}
	grade := year.Grades[0]
	if grade.ClassCount != 2 || grade.StudentCount != 33 || grade.FeeScheduleCount != 2 || grade.AdjustmentCount != 3 {
		t.Fatalf("unexpected grade aggregates: %+v", grade)
	}
	if grade.BillingReadyStudentCount != 32 || grade.MissingBillingRecipientCount != 1 || grade.CurrentFeeScheduleCount != 1 || grade.CurrentActiveScheduleCount != 1 || grade.CurrentInvoiceCount != 18 || grade.OpenInvoiceCount != 3 || grade.IssueCount != 17 {
		t.Fatalf("unexpected grade readiness aggregates: %+v", grade)
	}
	if len(grade.Classes) != 2 || grade.Classes[0].Name != "1A" || grade.Classes[1].Name != "1B" {
		t.Fatalf("unexpected class order: %+v", grade.Classes)
	}
	if grade.Classes[0].IssueCount != 1 || grade.Classes[1].IssueCount != 16 {
		t.Fatalf("unexpected class readiness issues: %+v", grade.Classes)
	}
}

func TestSchoolTreeReadinessScopeFromRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/school-tree?periodCode=2026-05&month=5", nil)
	scope, err := schoolTreeReadinessScopeFromRequest(request)
	if err != nil {
		t.Fatalf("expected valid readiness scope: %v", err)
	}
	if scope.PeriodCode != "2026-05" || scope.Month != 5 || !scope.HasMonth {
		t.Fatalf("unexpected readiness scope: %+v", scope)
	}

	request = httptest.NewRequest(http.MethodGet, "/api/v1/school-tree?month=13", nil)
	if _, err := schoolTreeReadinessScopeFromRequest(request); err == nil {
		t.Fatal("expected invalid month to fail")
	}
}

func TestNormalizeSchoolTreeInputs(t *testing.T) {
	school := normalizeSchoolTreeSchoolInput(schoolTreeSchoolInput{Code: "dekisugi", Name: " DEKISUGI ", Status: ""})
	if school.Code != "DEKISUGI" || school.Name != "DEKISUGI" || school.Status != "active" {
		t.Fatalf("unexpected normalized school: %+v", school)
	}
	if err := validateSchoolTreeSchoolInput(school); err != nil {
		t.Fatalf("expected normalized school to be valid: %v", err)
	}

	year := normalizeSchoolTreeSchoolYearInput(schoolTreeSchoolYearInput{SchoolID: " school-1 ", Code: "2025-2026", Status: "active"})
	if year.SchoolID != "school-1" || year.Name != "2025-2026" {
		t.Fatalf("unexpected normalized year: %+v", year)
	}
	if err := validateSchoolTreeClassInput(normalizeSchoolTreeClassInput(schoolTreeClassInput{SchoolYearID: "year-1", Grade: " 1 ", Name: " 1A "})); err != nil {
		t.Fatalf("expected class input to be valid: %v", err)
	}
}
