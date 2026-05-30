package main

import "testing"

func TestBuildFeeSchedulePreviewAppliesStudentAdjustments(t *testing.T) {
	schedule := feeScheduleInput{
		Items: []feeScheduleItemInput{
			{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000, DisplayOrder: 1},
			{FeeTypeCode: "shuttle", LabelVI: "Phi xe", LabelEN: "Shuttle", Amount: 500, DisplayOrder: 2},
		},
		Adjustments: []studentFeeAdjustmentInput{
			{StudentCode: "S001", AdjustmentType: "discount", FeeTypeCode: "tuition", LabelVI: "Giam hoc phi", LabelEN: "Tuition discount", Amount: 100, Reason: "Sibling discount"},
			{StudentCode: "S001", AdjustmentType: "surcharge", FeeTypeCode: "custom", LabelVI: "Phi bo sung", LabelEN: "Surcharge", Amount: 50, Reason: "Late registration"},
			{StudentCode: "S001", AdjustmentType: "waiver", FeeTypeCode: "shuttle", LabelVI: "Mien phi xe", LabelEN: "Shuttle waiver", Amount: 0, Reason: "No shuttle"},
			{StudentCode: "S001", AdjustmentType: "carry_over", FeeTypeCode: "previous_fees", LabelVI: "Phi thang truoc", LabelEN: "Previous fees", Amount: 200, Reason: "Previous balance"},
		},
	}
	students := []feeScheduleStudent{
		{ID: "student-1", StudentCode: "S001", StudentName: "Nguyen An", ClassName: "3.02", Grade: "3", SchoolYearCode: "2025-2026"},
		{ID: "student-2", StudentCode: "S002", StudentName: "Tran Binh", ClassName: "3.02", Grade: "3", SchoolYearCode: "2025-2026"},
	}

	preview, issues := buildFeeSchedulePreview(schedule, students)
	if len(issues) > 0 {
		t.Fatalf("expected no preview issues, got %+v", issues)
	}
	if len(preview.Rows) != 2 {
		t.Fatalf("expected two preview rows, got %d", len(preview.Rows))
	}

	first := preview.Rows[0]
	if first.BaseAmount != 1500 {
		t.Fatalf("expected base amount 1500, got %d", first.BaseAmount)
	}
	if first.AdjustmentAmount != -350 {
		t.Fatalf("expected adjustment amount -350, got %d", first.AdjustmentAmount)
	}
	if first.TotalAmount != 1150 {
		t.Fatalf("expected total amount 1150, got %d", first.TotalAmount)
	}
	if got := paymentItemsTotal(first.PaymentItems); got != first.TotalAmount {
		t.Fatalf("expected payment items to sum to preview total %d, got %d", first.TotalAmount, got)
	}

	second := preview.Rows[1]
	if second.AdjustmentAmount != 0 || second.TotalAmount != 1500 {
		t.Fatalf("expected second student to use base total only, got adjustment=%d total=%d", second.AdjustmentAmount, second.TotalAmount)
	}
	if len(second.PaymentItems) != 2 {
		t.Fatalf("expected unadjusted student to keep base payment items, got %+v", second.PaymentItems)
	}
}

func TestBuildFeeSchedulePreviewReportsAdjustmentWithoutReason(t *testing.T) {
	schedule := feeScheduleInput{
		Items: []feeScheduleItemInput{
			{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
		},
		Adjustments: []studentFeeAdjustmentInput{
			{StudentCode: "S001", AdjustmentType: "discount", Amount: 100},
		},
	}
	students := []feeScheduleStudent{
		{ID: "student-1", StudentCode: "S001", StudentName: "Nguyen An"},
	}

	_, issues := buildFeeSchedulePreview(schedule, students)
	if !hasFeeScheduleIssue(issues, "missing_adjustment_reason") {
		t.Fatalf("expected missing adjustment reason issue, got %+v", issues)
	}
}

func TestValidateFeeScheduleInputRequiresOperatorForSavedAdjustments(t *testing.T) {
	input := normalizeFeeScheduleInput(feeScheduleInput{
		SchoolYearID: "11111111-1111-1111-1111-111111111111",
		PeriodCode:   "2026-04",
		Items: []feeScheduleItemInput{
			{FeeTypeCode: "tuition", LabelVI: "Hoc phi", LabelEN: "Tuition", Amount: 1000},
		},
		Adjustments: []studentFeeAdjustmentInput{
			{StudentCode: "S001", AdjustmentType: "discount", Amount: 100, Reason: "Sibling discount"},
		},
	})

	issues := validateFeeScheduleInput(input, true)
	if !hasFeeScheduleIssue(issues, "missing_operator_name") {
		t.Fatalf("expected missing operator issue, got %+v", issues)
	}

	input.OperatorName = "Ke toan"
	if issues := validateFeeScheduleInput(input, true); hasFeeScheduleIssue(issues, "missing_operator_name") {
		t.Fatalf("did not expect missing operator issue after operator name, got %+v", issues)
	}
}

func hasFeeScheduleIssue(issues []feeSchedulePreviewIssue, issueType string) bool {
	for _, issue := range issues {
		if issue.Type == issueType {
			return true
		}
	}
	return false
}
