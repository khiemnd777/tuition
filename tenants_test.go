package main

import "testing"

func TestNormalizeTenantSaveInputCreatesStableTenantAndInitialSchool(t *testing.T) {
	input := normalizeTenantSaveInput(tenantSaveInput{
		Code:              " school b ",
		Name:              " School B ",
		Status:            "",
		InitialSchoolCode: "",
		InitialSchoolName: "",
	})

	if input.Code != "SCHOOL_B" || input.Name != "School B" || input.Status != "active" {
		t.Fatalf("unexpected tenant normalization: %+v", input)
	}
	if input.InitialSchoolCode != "SCHOOL_B" || input.InitialSchoolName != "School B" {
		t.Fatalf("expected initial school to default from tenant, got %+v", input)
	}
	if err := validateTenantSaveInput(input); err != nil {
		t.Fatalf("expected normalized tenant input to be valid: %v", err)
	}
}

func TestValidateTenantSaveInputRejectsUnsafeTenantData(t *testing.T) {
	cases := map[string]tenantSaveInput{
		"missing code":             {Name: "School B", InitialSchoolCode: "SCHOOL_B", InitialSchoolName: "School B"},
		"missing name":             {Code: "SCHOOL_B", InitialSchoolCode: "SCHOOL_B", InitialSchoolName: "School B"},
		"bad status":               {Code: "SCHOOL_B", Name: "School B", Status: "deleted", InitialSchoolCode: "SCHOOL_B", InitialSchoolName: "School B"},
		"missing school on create": {Code: "SCHOOL_B", Name: "School B", Status: "active"},
	}
	for name, input := range cases {
		if err := validateTenantSaveInput(input); err == nil {
			t.Fatalf("expected %s to be rejected", name)
		}
	}
}

func TestValidateTenantSwitchInputRequiresTenantID(t *testing.T) {
	input := normalizeTenantSwitchInput(tenantSwitchInput{TenantID: " tenant-1 "})
	if input.TenantID != "tenant-1" {
		t.Fatalf("unexpected switch normalization: %+v", input)
	}
	if err := validateTenantSwitchInput(input); err != nil {
		t.Fatalf("expected switch input to be valid: %v", err)
	}
	if err := validateTenantSwitchInput(tenantSwitchInput{}); err == nil {
		t.Fatal("expected missing tenant id to be rejected")
	}
}
