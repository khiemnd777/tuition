package main

import "testing"

func TestBuildSchoolTreeGroupsClassesByYearAndGrade(t *testing.T) {
	tree := buildSchoolTree(
		[]schoolTreeSchool{{ID: "school-1", Code: "ABC_SUN", Name: "ABC SUN", Status: "active"}},
		[]schoolTreeSchoolYear{
			{ID: "year-1", SchoolID: "school-1", SchoolCode: "ABC_SUN", Code: "2025-2026", Name: "2025-2026", Status: "active", ClassCount: 2, StudentCount: 33, FeeScheduleCount: 4, AdjustmentCount: 3},
			{ID: "year-2", SchoolID: "school-1", SchoolCode: "ABC_SUN", Code: "2024-2025", Name: "2024-2025", Status: "archived", ClassCount: 1, StudentCount: 10, FeeScheduleCount: 1, AdjustmentCount: 0},
		},
		[]schoolTreeClass{
			{ID: "class-1", SchoolID: "school-1", SchoolCode: "ABC_SUN", SchoolYearID: "year-1", SchoolYearCode: "2025-2026", Grade: "1", Name: "1A", Status: "active", StudentCount: 18, FeeScheduleCount: 1, ActiveFeeScheduleCount: 1, AdjustmentCount: 2},
			{ID: "class-2", SchoolID: "school-1", SchoolCode: "ABC_SUN", SchoolYearID: "year-1", SchoolYearCode: "2025-2026", Grade: "1", Name: "1B", Status: "active", StudentCount: 15, FeeScheduleCount: 1, ActiveFeeScheduleCount: 0, AdjustmentCount: 1},
			{ID: "class-3", SchoolID: "school-1", SchoolCode: "ABC_SUN", SchoolYearID: "year-2", SchoolYearCode: "2024-2025", Grade: "2", Name: "2A", Status: "active", StudentCount: 10, FeeScheduleCount: 1},
		},
	)

	if len(tree) != 1 {
		t.Fatalf("expected one school, got %d", len(tree))
	}
	if got := tree[0].StudentCount; got != 43 {
		t.Fatalf("school student count = %d, want 43", got)
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
	grade := year.Grades[0]
	if grade.ClassCount != 2 || grade.StudentCount != 33 || grade.FeeScheduleCount != 2 || grade.AdjustmentCount != 3 {
		t.Fatalf("unexpected grade aggregates: %+v", grade)
	}
	if len(grade.Classes) != 2 || grade.Classes[0].Name != "1A" || grade.Classes[1].Name != "1B" {
		t.Fatalf("unexpected class order: %+v", grade.Classes)
	}
}

func TestNormalizeSchoolTreeInputs(t *testing.T) {
	school := normalizeSchoolTreeSchoolInput(schoolTreeSchoolInput{Code: "abc sun", Name: " ABC SUN ", Status: ""})
	if school.Code != "ABC_SUN" || school.Name != "ABC SUN" || school.Status != "active" {
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
