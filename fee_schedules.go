package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	maxFeeScheduleItems       = 100
	maxStudentFeeAdjustments  = 2000
	defaultFeeScheduleStatus  = "draft"
	adjustmentTypeDiscount    = "discount"
	adjustmentTypeSurcharge   = "surcharge"
	adjustmentTypeWaiver      = "waiver"
	adjustmentTypeCarryOver   = "carry_over"
	feeScheduleScopeSchool    = "school_year"
	feeScheduleScopeGrade     = "grade"
	feeScheduleScopeClass     = "class"
	feeScheduleStatusDraft    = "draft"
	feeScheduleStatusActive   = "active"
	feeScheduleStatusArchived = "archived"
)

type feeTypeOption struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	LabelVI      string `json:"labelVi"`
	LabelEN      string `json:"labelEn"`
	Category     string `json:"category"`
	DisplayOrder int    `json:"displayOrder"`
	Status       string `json:"status"`
}

type feeScheduleOptionsResponse struct {
	FeeTypes    []feeTypeOption              `json:"feeTypes"`
	SchoolYears []masterDataSchoolYearOption `json:"schoolYears"`
	Classes     []masterDataClassOption      `json:"classes"`
}

type feeScheduleListFilters struct {
	SchoolYearID string
	ClassID      string
	Grade        string
	Status       string
}

type feeScheduleSummary struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	SchoolYearID    string `json:"schoolYearId"`
	SchoolYearCode  string `json:"schoolYearCode"`
	ScopeType       string `json:"scopeType"`
	ClassID         string `json:"classId,omitempty"`
	ClassName       string `json:"className,omitempty"`
	Grade           string `json:"grade,omitempty"`
	PeriodCode      string `json:"periodCode"`
	Month           int    `json:"month,omitempty"`
	Status          string `json:"status"`
	ItemTotal       int    `json:"itemTotal"`
	AdjustmentCount int    `json:"adjustmentCount"`
}

type feeScheduleInput struct {
	ID           string                      `json:"id,omitempty"`
	SchoolYearID string                      `json:"schoolYearId"`
	ClassID      string                      `json:"classId,omitempty"`
	Grade        string                      `json:"grade,omitempty"`
	PeriodCode   string                      `json:"periodCode"`
	Month        int                         `json:"month,omitempty"`
	Name         string                      `json:"name"`
	Notes        string                      `json:"notes,omitempty"`
	Status       string                      `json:"status"`
	Items        []feeScheduleItemInput      `json:"items"`
	Adjustments  []studentFeeAdjustmentInput `json:"adjustments,omitempty"`
}

type feeScheduleItemInput struct {
	FeeTypeID    string `json:"feeTypeId,omitempty"`
	FeeTypeCode  string `json:"feeTypeCode,omitempty"`
	LabelVI      string `json:"labelVi"`
	LabelEN      string `json:"labelEn"`
	Amount       int    `json:"amount"`
	DisplayOrder int    `json:"displayOrder"`
}

type studentFeeAdjustmentInput struct {
	StudentID      string `json:"studentId,omitempty"`
	StudentCode    string `json:"studentCode,omitempty"`
	AdjustmentType string `json:"adjustmentType"`
	FeeTypeID      string `json:"feeTypeId,omitempty"`
	FeeTypeCode    string `json:"feeTypeCode,omitempty"`
	LabelVI        string `json:"labelVi"`
	LabelEN        string `json:"labelEn"`
	Amount         int    `json:"amount"`
	Reason         string `json:"reason"`
}

type feeScheduleStudent struct {
	ID             string `json:"id"`
	StudentCode    string `json:"studentCode"`
	StudentName    string `json:"studentName"`
	ClassID        string `json:"classId"`
	ClassName      string `json:"className"`
	Grade          string `json:"grade"`
	SchoolYearCode string `json:"schoolYearCode"`
}

type feeSchedulePreview struct {
	Summary feeSchedulePreviewSummary `json:"summary"`
	Rows    []feeSchedulePreviewRow   `json:"rows"`
	Issues  []feeSchedulePreviewIssue `json:"issues,omitempty"`
}

type feeSchedulePreviewSummary struct {
	StudentCount int `json:"studentCount"`
	BaseAmount   int `json:"baseAmount"`
	Adjustments  int `json:"adjustments"`
	TotalAmount  int `json:"totalAmount"`
}

type feeSchedulePreviewRow struct {
	StudentID        string                         `json:"studentId"`
	StudentCode      string                         `json:"studentCode"`
	StudentName      string                         `json:"studentName"`
	ClassID          string                         `json:"classId"`
	ClassName        string                         `json:"className"`
	Grade            string                         `json:"grade"`
	SchoolYearCode   string                         `json:"schoolYearCode"`
	BaseAmount       int                            `json:"baseAmount"`
	AdjustmentAmount int                            `json:"adjustmentAmount"`
	TotalAmount      int                            `json:"totalAmount"`
	Items            []feeSchedulePreviewItem       `json:"items"`
	Adjustments      []feeSchedulePreviewAdjustment `json:"adjustments,omitempty"`
	PaymentItems     []paymentItem                  `json:"paymentItems"`
}

type feeSchedulePreviewItem struct {
	FeeTypeCode string `json:"feeTypeCode"`
	LabelVI     string `json:"labelVi"`
	LabelEN     string `json:"labelEn"`
	Amount      int    `json:"amount"`
}

type feeSchedulePreviewAdjustment struct {
	AdjustmentType string `json:"adjustmentType"`
	FeeTypeCode    string `json:"feeTypeCode,omitempty"`
	LabelVI        string `json:"labelVi"`
	LabelEN        string `json:"labelEn"`
	Amount         int    `json:"amount"`
	Delta          int    `json:"delta"`
	Reason         string `json:"reason"`
}

type feeSchedulePreviewIssue struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	StudentCode string `json:"studentCode,omitempty"`
	Field       string `json:"field,omitempty"`
}

type feeScheduleLine struct {
	FeeTypeID    string
	FeeTypeCode  string
	LabelVI      string
	LabelEN      string
	Amount       int
	DisplayOrder int
}

type feeTypeMaps struct {
	byID   map[string]feeTypeOption
	byCode map[string]feeTypeOption
}

func handleFeeScheduleOptions(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	feeTypes, err := listFeeTypes(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load fee types", http.StatusInternalServerError)
		return
	}
	masterOptions, err := listMasterDataOptions(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load master data options", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, feeScheduleOptionsResponse{
		FeeTypes:    feeTypes,
		SchoolYears: masterOptions.SchoolYears,
		Classes:     masterOptions.Classes,
	})
}

func handleFeeScheduleList(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	query := r.URL.Query()
	schedules, err := listFeeScheduleSummaries(r.Context(), db, feeScheduleListFilters{
		SchoolYearID: strings.TrimSpace(query.Get("schoolYearId")),
		ClassID:      strings.TrimSpace(query.Get("classId")),
		Grade:        normalizeGrade(query.Get("grade")),
		Status:       strings.TrimSpace(query.Get("status")),
	})
	if err != nil {
		http.Error(w, "cannot load fee schedules", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"schedules": schedules})
}

func handleFeeSchedulePreview(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input feeScheduleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeFeeScheduleInput(input)
	if issues := validateFeeScheduleInput(input, true); len(issues) > 0 {
		writeJSON(w, http.StatusOK, feeSchedulePreview{Issues: issues})
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	students, err := loadFeeScheduleStudents(r.Context(), db, input)
	if err != nil {
		http.Error(w, "cannot load fee schedule students", http.StatusInternalServerError)
		return
	}
	preview, issues := buildFeeSchedulePreview(input, students)
	preview.Issues = issues
	writeJSON(w, http.StatusOK, preview)
}

func handleFeeScheduleSave(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input feeScheduleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeFeeScheduleInput(input)
	if issues := validateFeeScheduleInput(input, true); len(issues) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"issues": issues})
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	savedID, err := saveFeeSchedule(r.Context(), db, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = savedID

	students, err := loadFeeScheduleStudents(r.Context(), db, input)
	if err != nil {
		http.Error(w, "cannot load fee schedule students", http.StatusInternalServerError)
		return
	}
	preview, issues := buildFeeSchedulePreview(input, students)
	preview.Issues = issues

	schedules, err := listFeeScheduleSummaries(r.Context(), db, feeScheduleListFilters{SchoolYearID: input.SchoolYearID})
	if err != nil {
		http.Error(w, "cannot load fee schedules", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":        savedID,
		"preview":   preview,
		"schedules": schedules,
	})
}

func buildFeeSchedulePreview(input feeScheduleInput, students []feeScheduleStudent) (feeSchedulePreview, []feeSchedulePreviewIssue) {
	input = normalizeFeeScheduleInput(input)
	issues := validateFeeScheduleInput(input, false)
	lines := normalizedFeeScheduleLines(input.Items)
	if len(lines) == 0 {
		issues = append(issues, feeSchedulePreviewIssue{
			Type:    "missing_fee_items",
			Field:   "items",
			Message: "at least one fee schedule item with a positive amount is required",
		})
	}

	studentByKey := map[string]feeScheduleStudent{}
	for _, student := range students {
		if student.ID != "" {
			studentByKey["id:"+strings.ToLower(student.ID)] = student
		}
		if student.StudentCode != "" {
			studentByKey["code:"+normalizeStudentCode(student.StudentCode)] = student
		}
	}

	adjustmentsByStudent := map[string][]studentFeeAdjustmentInput{}
	for _, adjustment := range input.Adjustments {
		key := feeAdjustmentStudentKey(adjustment)
		if key == "" {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:    "missing_adjustment_student",
				Field:   "adjustments.studentCode",
				Message: "student_code or student_id is required for every adjustment",
			})
			continue
		}
		student, ok := studentByKey[key]
		if !ok {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:        "unknown_adjustment_student",
				Field:       "adjustments.studentCode",
				StudentCode: adjustment.StudentCode,
				Message:     "adjustment references a student outside the selected fee schedule scope",
			})
			continue
		}
		adjustmentsByStudent["id:"+strings.ToLower(student.ID)] = append(adjustmentsByStudent["id:"+strings.ToLower(student.ID)], adjustment)
	}

	sort.SliceStable(students, func(i, j int) bool {
		if students[i].Grade != students[j].Grade {
			return students[i].Grade < students[j].Grade
		}
		if students[i].ClassName != students[j].ClassName {
			return students[i].ClassName < students[j].ClassName
		}
		return students[i].StudentCode < students[j].StudentCode
	})

	preview := feeSchedulePreview{
		Rows: []feeSchedulePreviewRow{},
	}
	for _, student := range students {
		baseItems := feePreviewItems(lines)
		baseAmount := feePreviewItemsTotal(baseItems)
		baseByFeeTypeCode := feeBaseAmountByFeeType(lines)
		rowAdjustments := []feeSchedulePreviewAdjustment{}
		adjustmentAmount := 0

		for _, adjustment := range adjustmentsByStudent["id:"+strings.ToLower(student.ID)] {
			delta, issue := feeAdjustmentDelta(adjustment, baseByFeeTypeCode)
			if issue != nil {
				issue.StudentCode = student.StudentCode
				issues = append(issues, *issue)
				continue
			}
			labelVI, labelEN := feeAdjustmentLabels(adjustment)
			rowAdjustments = append(rowAdjustments, feeSchedulePreviewAdjustment{
				AdjustmentType: adjustment.AdjustmentType,
				FeeTypeCode:    adjustment.FeeTypeCode,
				LabelVI:        labelVI,
				LabelEN:        labelEN,
				Amount:         adjustment.Amount,
				Delta:          delta,
				Reason:         adjustment.Reason,
			})
			adjustmentAmount += delta
		}

		totalAmount := baseAmount + adjustmentAmount
		if totalAmount < 0 {
			totalAmount = 0
		}
		row := feeSchedulePreviewRow{
			StudentID:        student.ID,
			StudentCode:      student.StudentCode,
			StudentName:      student.StudentName,
			ClassID:          student.ClassID,
			ClassName:        student.ClassName,
			Grade:            student.Grade,
			SchoolYearCode:   student.SchoolYearCode,
			BaseAmount:       baseAmount,
			AdjustmentAmount: totalAmount - baseAmount,
			TotalAmount:      totalAmount,
			Items:            baseItems,
			Adjustments:      rowAdjustments,
			PaymentItems:     feePreviewPaymentItems(baseItems, totalAmount),
		}
		preview.Rows = append(preview.Rows, row)
		preview.Summary.StudentCount++
		preview.Summary.BaseAmount += row.BaseAmount
		preview.Summary.Adjustments += row.AdjustmentAmount
		preview.Summary.TotalAmount += row.TotalAmount
	}

	return preview, issues
}

func validateFeeScheduleInput(input feeScheduleInput, requireScope bool) []feeSchedulePreviewIssue {
	issues := []feeSchedulePreviewIssue{}
	if requireScope {
		if input.SchoolYearID == "" {
			issues = append(issues, feeSchedulePreviewIssue{Type: "missing_school_year", Field: "schoolYearId", Message: "schoolYearId is required"})
		}
		if input.PeriodCode == "" {
			issues = append(issues, feeSchedulePreviewIssue{Type: "missing_period", Field: "periodCode", Message: "periodCode is required"})
		}
	}
	if input.Month < 0 || input.Month > 12 {
		issues = append(issues, feeSchedulePreviewIssue{Type: "invalid_month", Field: "month", Message: "month must be from 1 to 12 when provided"})
	}
	if input.Status != "" && !validFeeScheduleStatus(input.Status) {
		issues = append(issues, feeSchedulePreviewIssue{Type: "invalid_schedule_status", Field: "status", Message: "status must be draft, active, or archived"})
	}
	if len(input.Items) > maxFeeScheduleItems {
		issues = append(issues, feeSchedulePreviewIssue{Type: "too_many_fee_items", Field: "items", Message: fmt.Sprintf("at most %d fee items are allowed", maxFeeScheduleItems)})
	}
	if len(input.Adjustments) > maxStudentFeeAdjustments {
		issues = append(issues, feeSchedulePreviewIssue{Type: "too_many_adjustments", Field: "adjustments", Message: fmt.Sprintf("at most %d student adjustments are allowed", maxStudentFeeAdjustments)})
	}
	for _, item := range input.Items {
		if item.Amount < 0 {
			issues = append(issues, feeSchedulePreviewIssue{Type: "negative_fee_item_amount", Field: "items.amount", Message: "fee schedule item amount must not be negative"})
		}
	}
	for _, adjustment := range input.Adjustments {
		if adjustment.AdjustmentType == "" || !validFeeAdjustmentType(adjustment.AdjustmentType) {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:        "invalid_adjustment_type",
				Field:       "adjustments.adjustmentType",
				StudentCode: adjustment.StudentCode,
				Message:     "adjustment_type must be discount, surcharge, waiver, or carry_over",
			})
		}
		if adjustment.Reason == "" {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:        "missing_adjustment_reason",
				Field:       "adjustments.reason",
				StudentCode: adjustment.StudentCode,
				Message:     "every student fee adjustment requires a reason",
			})
		}
		if adjustment.Amount < 0 {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:        "negative_adjustment_amount",
				Field:       "adjustments.amount",
				StudentCode: adjustment.StudentCode,
				Message:     "adjustment amount must not be negative",
			})
		}
		if adjustment.AdjustmentType != adjustmentTypeWaiver && adjustment.Amount <= 0 {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:        "missing_adjustment_amount",
				Field:       "adjustments.amount",
				StudentCode: adjustment.StudentCode,
				Message:     "discount, surcharge, and carry-over adjustments require a positive amount",
			})
		}
		if adjustment.AdjustmentType == adjustmentTypeWaiver && adjustment.Amount == 0 && adjustment.FeeTypeCode == "" && adjustment.FeeTypeID == "" {
			issues = append(issues, feeSchedulePreviewIssue{
				Type:        "missing_waiver_target",
				Field:       "adjustments.feeTypeCode",
				StudentCode: adjustment.StudentCode,
				Message:     "zero-amount waiver requires a fee type target",
			})
		}
	}
	return issues
}

func normalizedFeeScheduleLines(items []feeScheduleItemInput) []feeScheduleLine {
	lines := []feeScheduleLine{}
	for _, item := range items {
		if item.Amount <= 0 {
			continue
		}
		line := feeScheduleLine{
			FeeTypeID:    strings.TrimSpace(item.FeeTypeID),
			FeeTypeCode:  headerKey(item.FeeTypeCode),
			LabelVI:      strings.TrimSpace(item.LabelVI),
			LabelEN:      strings.TrimSpace(item.LabelEN),
			Amount:       item.Amount,
			DisplayOrder: item.DisplayOrder,
		}
		if line.LabelVI == "" {
			line.LabelVI = line.LabelEN
		}
		if line.LabelEN == "" {
			line.LabelEN = line.LabelVI
		}
		if line.LabelVI == "" && line.FeeTypeCode != "" {
			line.LabelVI = line.FeeTypeCode
			line.LabelEN = line.FeeTypeCode
		}
		lines = append(lines, line)
	}
	sort.SliceStable(lines, func(i, j int) bool {
		if lines[i].DisplayOrder != lines[j].DisplayOrder {
			return lines[i].DisplayOrder < lines[j].DisplayOrder
		}
		return lines[i].FeeTypeCode < lines[j].FeeTypeCode
	})
	return lines
}

func feePreviewItems(lines []feeScheduleLine) []feeSchedulePreviewItem {
	items := make([]feeSchedulePreviewItem, 0, len(lines))
	for _, line := range lines {
		items = append(items, feeSchedulePreviewItem{
			FeeTypeCode: line.FeeTypeCode,
			LabelVI:     line.LabelVI,
			LabelEN:     line.LabelEN,
			Amount:      line.Amount,
		})
	}
	return items
}

func feePreviewItemsTotal(items []feeSchedulePreviewItem) int {
	total := 0
	for _, item := range items {
		total += item.Amount
	}
	return total
}

func feeBaseAmountByFeeType(lines []feeScheduleLine) map[string]int {
	amounts := map[string]int{}
	for _, line := range lines {
		if line.FeeTypeCode != "" {
			amounts[line.FeeTypeCode] += line.Amount
		}
	}
	return amounts
}

func feePreviewPaymentItems(items []feeSchedulePreviewItem, totalAmount int) []paymentItem {
	baseTotal := feePreviewItemsTotal(items)
	if baseTotal == totalAmount {
		paymentItems := make([]paymentItem, 0, len(items))
		for _, item := range items {
			paymentItems = append(paymentItems, paymentItem{
				Label:   item.LabelVI,
				LabelEN: item.LabelEN,
				Amount:  item.Amount,
			})
		}
		return paymentItems
	}
	return []paymentItem{{
		Label:   "Tổng phí sau điều chỉnh",
		LabelEN: "Total after adjustments",
		Amount:  totalAmount,
	}}
}

func feeAdjustmentDelta(adjustment studentFeeAdjustmentInput, baseByFeeTypeCode map[string]int) (int, *feeSchedulePreviewIssue) {
	switch adjustment.AdjustmentType {
	case adjustmentTypeDiscount:
		return -adjustment.Amount, nil
	case adjustmentTypeSurcharge, adjustmentTypeCarryOver:
		return adjustment.Amount, nil
	case adjustmentTypeWaiver:
		amount := adjustment.Amount
		if amount == 0 {
			amount = baseByFeeTypeCode[adjustment.FeeTypeCode]
			if amount == 0 {
				return 0, &feeSchedulePreviewIssue{
					Type:    "waiver_fee_type_not_found",
					Field:   "adjustments.feeTypeCode",
					Message: "waiver fee type does not match a positive schedule item",
				}
			}
		}
		return -amount, nil
	default:
		return 0, &feeSchedulePreviewIssue{
			Type:    "invalid_adjustment_type",
			Field:   "adjustments.adjustmentType",
			Message: "adjustment_type must be discount, surcharge, waiver, or carry_over",
		}
	}
}

func feeAdjustmentLabels(adjustment studentFeeAdjustmentInput) (string, string) {
	labelVI := strings.TrimSpace(adjustment.LabelVI)
	labelEN := strings.TrimSpace(adjustment.LabelEN)
	if labelVI == "" || labelEN == "" {
		defaultVI, defaultEN := feeAdjustmentDefaultLabels(adjustment.AdjustmentType)
		if labelVI == "" {
			labelVI = defaultVI
		}
		if labelEN == "" {
			labelEN = defaultEN
		}
	}
	return labelVI, labelEN
}

func feeAdjustmentDefaultLabels(adjustmentType string) (string, string) {
	switch adjustmentType {
	case adjustmentTypeDiscount:
		return "Giảm trừ", "Discount"
	case adjustmentTypeSurcharge:
		return "Phụ thu", "Surcharge"
	case adjustmentTypeWaiver:
		return "Miễn giảm", "Waiver"
	case adjustmentTypeCarryOver:
		return "Chuyển kỳ trước", "Carry-over"
	default:
		return "Điều chỉnh", "Adjustment"
	}
}

func feeAdjustmentStudentKey(adjustment studentFeeAdjustmentInput) string {
	if adjustment.StudentID != "" {
		return "id:" + strings.ToLower(adjustment.StudentID)
	}
	if adjustment.StudentCode != "" {
		return "code:" + normalizeStudentCode(adjustment.StudentCode)
	}
	return ""
}

func normalizeFeeScheduleInput(input feeScheduleInput) feeScheduleInput {
	input.ID = strings.TrimSpace(input.ID)
	input.SchoolYearID = strings.TrimSpace(input.SchoolYearID)
	input.ClassID = strings.TrimSpace(input.ClassID)
	input.Grade = normalizeGrade(input.Grade)
	input.PeriodCode = strings.TrimSpace(input.PeriodCode)
	input.Name = strings.TrimSpace(input.Name)
	input.Notes = strings.TrimSpace(input.Notes)
	input.Status = feeScheduleStatus(input.Status)
	if input.ClassID != "" {
		input.Grade = ""
	}
	for idx := range input.Items {
		input.Items[idx].FeeTypeID = strings.TrimSpace(input.Items[idx].FeeTypeID)
		input.Items[idx].FeeTypeCode = headerKey(input.Items[idx].FeeTypeCode)
		input.Items[idx].LabelVI = strings.TrimSpace(input.Items[idx].LabelVI)
		input.Items[idx].LabelEN = strings.TrimSpace(input.Items[idx].LabelEN)
	}
	for idx := range input.Adjustments {
		input.Adjustments[idx].StudentID = strings.TrimSpace(input.Adjustments[idx].StudentID)
		input.Adjustments[idx].StudentCode = normalizeStudentCode(input.Adjustments[idx].StudentCode)
		input.Adjustments[idx].AdjustmentType = headerKey(input.Adjustments[idx].AdjustmentType)
		input.Adjustments[idx].FeeTypeID = strings.TrimSpace(input.Adjustments[idx].FeeTypeID)
		input.Adjustments[idx].FeeTypeCode = headerKey(input.Adjustments[idx].FeeTypeCode)
		input.Adjustments[idx].LabelVI = strings.TrimSpace(input.Adjustments[idx].LabelVI)
		input.Adjustments[idx].LabelEN = strings.TrimSpace(input.Adjustments[idx].LabelEN)
		input.Adjustments[idx].Reason = strings.TrimSpace(input.Adjustments[idx].Reason)
	}
	return input
}

func feeScheduleStatus(value string) string {
	value = headerKey(value)
	if value == "" {
		return defaultFeeScheduleStatus
	}
	return value
}

func validFeeScheduleStatus(value string) bool {
	switch value {
	case feeScheduleStatusDraft, feeScheduleStatusActive, feeScheduleStatusArchived:
		return true
	default:
		return false
	}
}

func validFeeAdjustmentType(value string) bool {
	switch value {
	case adjustmentTypeDiscount, adjustmentTypeSurcharge, adjustmentTypeWaiver, adjustmentTypeCarryOver:
		return true
	default:
		return false
	}
}

func feeScheduleScopeType(input feeScheduleInput) string {
	if input.ClassID != "" {
		return feeScheduleScopeClass
	}
	if input.Grade != "" {
		return feeScheduleScopeGrade
	}
	return feeScheduleScopeSchool
}

func listFeeTypes(ctx context.Context, exec masterDataExecutor) ([]feeTypeOption, error) {
	rows, err := exec.QueryContext(ctx, `
SELECT id::text, code, label_vi, label_en, category, default_display_order, status
FROM fee_types
WHERE status = 'active'
ORDER BY default_display_order, code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeTypes := []feeTypeOption{}
	for rows.Next() {
		var item feeTypeOption
		if err := rows.Scan(&item.ID, &item.Code, &item.LabelVI, &item.LabelEN, &item.Category, &item.DisplayOrder, &item.Status); err != nil {
			return nil, err
		}
		feeTypes = append(feeTypes, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeTypes, nil
}

func listFeeScheduleSummaries(ctx context.Context, db *sql.DB, filters feeScheduleListFilters) ([]feeScheduleSummary, error) {
	conditions := []string{"1 = 1"}
	args := []any{}
	if filters.SchoolYearID != "" {
		args = append(args, filters.SchoolYearID)
		conditions = append(conditions, fmt.Sprintf("fs.school_year_id = $%d::uuid", len(args)))
	}
	if filters.ClassID != "" {
		args = append(args, filters.ClassID)
		conditions = append(conditions, fmt.Sprintf("fs.class_id = $%d::uuid", len(args)))
	}
	if filters.Grade != "" {
		args = append(args, filters.Grade)
		conditions = append(conditions, fmt.Sprintf("(fs.grade = $%d OR c.grade = $%d)", len(args), len(args)))
	}
	if filters.Status != "" {
		args = append(args, filters.Status)
		conditions = append(conditions, fmt.Sprintf("fs.status = $%d", len(args)))
	}

	query := `
SELECT fs.id::text,
	fs.name,
	fs.school_year_id::text,
	sy.code,
	fs.scope_type,
	COALESCE(fs.class_id::text, ''),
	COALESCE(c.name, ''),
	fs.grade,
	fs.period_code,
	fs.month,
	fs.status,
	COALESCE(item_totals.item_total, 0),
	COALESCE(adjustment_counts.adjustment_count, 0)
FROM fee_schedules fs
JOIN school_years sy ON sy.id = fs.school_year_id
LEFT JOIN classes c ON c.id = fs.class_id
LEFT JOIN (
	SELECT schedule_id, SUM(amount)::integer AS item_total
	FROM fee_schedule_items
	GROUP BY schedule_id
) item_totals ON item_totals.schedule_id = fs.id
LEFT JOIN (
	SELECT schedule_id, COUNT(*)::integer AS adjustment_count
	FROM student_fee_adjustments
	WHERE status = 'active'
	GROUP BY schedule_id
) adjustment_counts ON adjustment_counts.schedule_id = fs.id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY sy.code DESC, fs.period_code DESC, fs.created_at DESC
LIMIT 500`

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := []feeScheduleSummary{}
	for rows.Next() {
		var schedule feeScheduleSummary
		var month sql.NullInt64
		if err := rows.Scan(
			&schedule.ID,
			&schedule.Name,
			&schedule.SchoolYearID,
			&schedule.SchoolYearCode,
			&schedule.ScopeType,
			&schedule.ClassID,
			&schedule.ClassName,
			&schedule.Grade,
			&schedule.PeriodCode,
			&month,
			&schedule.Status,
			&schedule.ItemTotal,
			&schedule.AdjustmentCount,
		); err != nil {
			return nil, err
		}
		if month.Valid {
			schedule.Month = int(month.Int64)
		}
		schedules = append(schedules, schedule)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}

func loadFeeScheduleStudents(ctx context.Context, exec masterDataExecutor, input feeScheduleInput) ([]feeScheduleStudent, error) {
	conditions := []string{"sy.id = $1::uuid", "s.status = 'active'", "c.status = 'active'"}
	args := []any{input.SchoolYearID}
	if input.ClassID != "" {
		args = append(args, input.ClassID)
		conditions = append(conditions, fmt.Sprintf("c.id = $%d::uuid", len(args)))
	} else if input.Grade != "" {
		args = append(args, input.Grade)
		conditions = append(conditions, fmt.Sprintf("c.grade = $%d", len(args)))
	}

	query := `
SELECT s.id::text, s.student_code, s.full_name, c.id::text, c.name, c.grade, sy.code
FROM students s
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY c.grade, c.name, s.student_code
LIMIT 2000`

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []feeScheduleStudent{}
	for rows.Next() {
		var student feeScheduleStudent
		if err := rows.Scan(&student.ID, &student.StudentCode, &student.StudentName, &student.ClassID, &student.ClassName, &student.Grade, &student.SchoolYearCode); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func saveFeeSchedule(ctx context.Context, db *sql.DB, input feeScheduleInput) (string, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	feeTypes, err := loadFeeTypeMaps(ctx, tx)
	if err != nil {
		return "", err
	}
	lines, err := resolveFeeScheduleLines(input.Items, feeTypes)
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", errors.New("at least one fee schedule item with a positive amount is required")
	}

	scopeType := feeScheduleScopeType(input)
	classID := nullableString(input.ClassID)
	month := nullableInt(input.Month)
	var scheduleID string
	if input.ID == "" {
		err = tx.QueryRowContext(ctx, `
INSERT INTO fee_schedules (school_year_id, scope_type, class_id, grade, period_code, month, name, notes, status)
VALUES ($1::uuid, $2, $3::uuid, $4, $5, $6, $7, $8, $9)
RETURNING id::text`,
			input.SchoolYearID,
			scopeType,
			classID,
			input.Grade,
			input.PeriodCode,
			month,
			input.Name,
			input.Notes,
			input.Status,
		).Scan(&scheduleID)
	} else {
		err = tx.QueryRowContext(ctx, `
UPDATE fee_schedules
SET school_year_id = $2::uuid,
	scope_type = $3,
	class_id = $4::uuid,
	grade = $5,
	period_code = $6,
	month = $7,
	name = $8,
	notes = $9,
	status = $10
WHERE id = $1::uuid
RETURNING id::text`,
			input.ID,
			input.SchoolYearID,
			scopeType,
			classID,
			input.Grade,
			input.PeriodCode,
			month,
			input.Name,
			input.Notes,
			input.Status,
		).Scan(&scheduleID)
	}
	if err != nil {
		return "", err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM fee_schedule_items WHERE schedule_id = $1::uuid`, scheduleID); err != nil {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM student_fee_adjustments WHERE schedule_id = $1::uuid`, scheduleID); err != nil {
		return "", err
	}

	for _, line := range lines {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO fee_schedule_items (schedule_id, fee_type_id, label_vi, label_en, amount, display_order)
VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6)`,
			scheduleID,
			line.FeeTypeID,
			line.LabelVI,
			line.LabelEN,
			line.Amount,
			line.DisplayOrder,
		); err != nil {
			return "", err
		}
	}

	for _, adjustment := range input.Adjustments {
		studentID, err := resolveFeeAdjustmentStudentID(ctx, tx, adjustment)
		if err != nil {
			return "", err
		}
		feeTypeID, feeType, err := resolveOptionalFeeType(adjustment.FeeTypeID, adjustment.FeeTypeCode, feeTypes)
		if err != nil {
			return "", err
		}
		labelVI, labelEN := feeAdjustmentLabels(adjustment)
		if labelVI == "" && feeType.Code != "" {
			labelVI = feeType.LabelVI
		}
		if labelEN == "" && feeType.Code != "" {
			labelEN = feeType.LabelEN
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO student_fee_adjustments (schedule_id, student_id, fee_type_id, adjustment_type, label_vi, label_en, amount, reason)
VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6, $7, $8)`,
			scheduleID,
			studentID,
			nullableString(feeTypeID),
			adjustment.AdjustmentType,
			labelVI,
			labelEN,
			adjustment.Amount,
			adjustment.Reason,
		); err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return scheduleID, nil
}

func loadFeeTypeMaps(ctx context.Context, exec masterDataExecutor) (feeTypeMaps, error) {
	feeTypes, err := listFeeTypes(ctx, exec)
	if err != nil {
		return feeTypeMaps{}, err
	}
	maps := feeTypeMaps{byID: map[string]feeTypeOption{}, byCode: map[string]feeTypeOption{}}
	for _, feeType := range feeTypes {
		maps.byID[feeType.ID] = feeType
		maps.byCode[feeType.Code] = feeType
	}
	return maps, nil
}

func resolveFeeScheduleLines(items []feeScheduleItemInput, feeTypes feeTypeMaps) ([]feeScheduleLine, error) {
	lines := []feeScheduleLine{}
	for _, item := range normalizedFeeScheduleLines(items) {
		feeType, ok := resolveFeeType(item.FeeTypeID, item.FeeTypeCode, feeTypes)
		if !ok {
			return nil, fmt.Errorf("unknown fee type %q", firstNonEmpty(item.FeeTypeCode, item.FeeTypeID))
		}
		if item.LabelVI == "" {
			item.LabelVI = feeType.LabelVI
		}
		if item.LabelEN == "" {
			item.LabelEN = feeType.LabelEN
		}
		if item.DisplayOrder == 0 {
			item.DisplayOrder = feeType.DisplayOrder
		}
		item.FeeTypeID = feeType.ID
		item.FeeTypeCode = feeType.Code
		lines = append(lines, item)
	}
	return lines, nil
}

func resolveFeeType(id string, code string, feeTypes feeTypeMaps) (feeTypeOption, bool) {
	if id != "" {
		if feeType, ok := feeTypes.byID[id]; ok {
			return feeType, true
		}
	}
	if code != "" {
		if feeType, ok := feeTypes.byCode[code]; ok {
			return feeType, true
		}
	}
	return feeTypeOption{}, false
}

func resolveOptionalFeeType(id string, code string, feeTypes feeTypeMaps) (string, feeTypeOption, error) {
	if id == "" && code == "" {
		return "", feeTypeOption{}, nil
	}
	feeType, ok := resolveFeeType(id, code, feeTypes)
	if !ok {
		return "", feeTypeOption{}, fmt.Errorf("unknown fee type %q", firstNonEmpty(code, id))
	}
	return feeType.ID, feeType, nil
}

func resolveFeeAdjustmentStudentID(ctx context.Context, exec masterDataExecutor, adjustment studentFeeAdjustmentInput) (string, error) {
	if adjustment.StudentID != "" {
		return adjustment.StudentID, nil
	}
	if adjustment.StudentCode == "" {
		return "", errors.New("student adjustment requires studentCode or studentId")
	}
	var studentID string
	err := exec.QueryRowContext(ctx, `
SELECT id::text
FROM students
WHERE student_code = $1`, adjustment.StudentCode).Scan(&studentID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("student_code %s was not found", adjustment.StudentCode)
	}
	return studentID, err
}

func nullableString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func nullableInt(value int) any {
	if value == 0 {
		return nil
	}
	return value
}
