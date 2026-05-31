package main

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	adminReadinessSeverityBlocking = "blocking"
	adminReadinessSeverityWarning  = "warning"
	adminReadinessSeverityInfo     = "info"
)

type adminReadinessCenter struct {
	Summary adminReadinessSummary `json:"summary"`
	Issues  []adminReadinessIssue `json:"issues"`
}

type adminReadinessSummary struct {
	BlockingCount int `json:"blockingCount"`
	WarningCount  int `json:"warningCount"`
	InfoCount     int `json:"infoCount"`
	TotalCount    int `json:"totalCount"`
}

type adminReadinessIssue struct {
	ID             string `json:"id"`
	Severity       string `json:"severity"`
	Type           string `json:"type"`
	EntityType     string `json:"entityType"`
	EntityID       string `json:"entityId,omitempty"`
	EntityLabel    string `json:"entityLabel"`
	Scope          string `json:"scope"`
	Message        string `json:"message"`
	Action         string `json:"action,omitempty"`
	SchoolID       string `json:"schoolId,omitempty"`
	SchoolYearID   string `json:"schoolYearId,omitempty"`
	ClassID        string `json:"classId,omitempty"`
	Grade          string `json:"grade,omitempty"`
	PeriodCode     string `json:"periodCode,omitempty"`
	Month          int    `json:"month,omitempty"`
	ReferenceCode  string `json:"referenceCode,omitempty"`
	ReferenceCount int    `json:"referenceCount,omitempty"`
}

func buildAdminReadinessCenter(ctx context.Context, db *sql.DB, filters adminFilters, invoices []adminInvoiceReportRow, transactions []paymentTransactionSummary, now time.Time) (adminReadinessCenter, error) {
	issues := []adminReadinessIssue{}
	loaders := []func(context.Context, *sql.DB, adminFilters) ([]adminReadinessIssue, error){
		listAdminStudentReadinessIssues,
		listAdminClassReadinessIssues,
		listAdminFeeScheduleReadinessIssues,
		listAdminAdjustmentReadinessIssues,
		listAdminNotificationReadinessIssues,
		listAdminOperationReadinessIssues,
	}
	for _, load := range loaders {
		rows, err := load(ctx, db, filters)
		if err != nil {
			return adminReadinessCenter{}, err
		}
		issues = append(issues, rows...)
	}
	issues = appendAdminInvoiceReadinessIssues(issues, invoices)
	issues = appendAdminTransactionReadinessIssues(issues, filterAdminReportTransactions(invoices, transactions, filters))
	issues = append(issues, adminLocalReadinessIssues(now)...)
	return buildAdminReadinessCenterFromIssues(issues), nil
}

func buildAdminReadinessCenterFromIssues(issues []adminReadinessIssue) adminReadinessCenter {
	center := adminReadinessCenter{
		Issues: make([]adminReadinessIssue, 0, len(issues)),
	}
	for _, issue := range issues {
		issue.Severity = normalizeAdminReadinessSeverity(issue.Severity)
		issue.Type = headerKey(issue.Type)
		issue.EntityType = headerKey(issue.EntityType)
		issue.Action = headerKey(issue.Action)
		issue.Grade = normalizeGrade(issue.Grade)
		issue.PeriodCode = strings.TrimSpace(issue.PeriodCode)
		if issue.Scope == "" {
			issue.Scope = "Toàn hệ thống"
		}
		if issue.EntityLabel == "" {
			issue.EntityLabel = firstNonEmpty(issue.ReferenceCode, issue.EntityID, issue.EntityType)
		}
		if issue.Message == "" || issue.Type == "" {
			continue
		}
		center.Issues = append(center.Issues, issue)
	}
	sort.SliceStable(center.Issues, func(i, j int) bool {
		left := center.Issues[i]
		right := center.Issues[j]
		if adminReadinessSeverityRank(left.Severity) != adminReadinessSeverityRank(right.Severity) {
			return adminReadinessSeverityRank(left.Severity) < adminReadinessSeverityRank(right.Severity)
		}
		if left.Scope != right.Scope {
			return left.Scope < right.Scope
		}
		if left.EntityType != right.EntityType {
			return left.EntityType < right.EntityType
		}
		if left.EntityLabel != right.EntityLabel {
			return left.EntityLabel < right.EntityLabel
		}
		return left.Type < right.Type
	})
	for idx := range center.Issues {
		if center.Issues[idx].ID == "" {
			center.Issues[idx].ID = adminReadinessIssueID(center.Issues[idx], idx)
		}
		switch center.Issues[idx].Severity {
		case adminReadinessSeverityBlocking:
			center.Summary.BlockingCount++
		case adminReadinessSeverityWarning:
			center.Summary.WarningCount++
		default:
			center.Summary.InfoCount++
		}
	}
	center.Summary.TotalCount = len(center.Issues)
	return center
}

func normalizeAdminReadinessSeverity(value string) string {
	switch headerKey(value) {
	case adminReadinessSeverityBlocking:
		return adminReadinessSeverityBlocking
	case adminReadinessSeverityInfo, "informational":
		return adminReadinessSeverityInfo
	default:
		return adminReadinessSeverityWarning
	}
}

func adminReadinessSeverityRank(value string) int {
	switch normalizeAdminReadinessSeverity(value) {
	case adminReadinessSeverityBlocking:
		return 0
	case adminReadinessSeverityWarning:
		return 1
	default:
		return 2
	}
}

func adminReadinessIssueID(issue adminReadinessIssue, index int) string {
	parts := []string{issue.Severity, issue.Type, issue.EntityType, issue.EntityID, issue.ReferenceCode, issue.PeriodCode}
	if issue.Month > 0 {
		parts = append(parts, strconv.Itoa(issue.Month))
	}
	id := headerKey(strings.Join(parts, "_"))
	if id == "" {
		return fmt.Sprintf("readiness_%03d", index+1)
	}
	return id
}

func listAdminStudentReadinessIssues(ctx context.Context, db *sql.DB, filters adminFilters) ([]adminReadinessIssue, error) {
	args := []any{}
	addArg := adminReadinessArgFunc(&args)
	conditions := []string{"s.status <> 'inactive'"}
	if filters.SchoolID != "" {
		conditions = append(conditions, "sy.school_id = "+addArg(filters.SchoolID)+"::uuid")
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "c.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.ClassID != "" {
		conditions = append(conditions, "s.class_id = "+addArg(filters.ClassID)+"::uuid")
	}
	if filters.Grade != "" {
		conditions = append(conditions, "lower(c.grade) = lower("+addArg(filters.Grade)+")")
	}
	query := `
SELECT
	s.id::text,
	s.student_code,
	s.full_name,
	COALESCE(c.id::text, ''),
	COALESCE(c.name, ''),
	COALESCE(c.grade, ''),
	COALESCE(sy.id::text, ''),
	COALESCE(sy.code, ''),
	COALESCE(sy.school_id::text, ''),
	COALESCE(sc.code, ''),
	COUNT(sp.parent_id) FILTER (WHERE sp.is_active AND p.status = 'active')::integer,
	COUNT(sp.parent_id) FILTER (WHERE sp.receives_billing_email)::integer,
	COUNT(sp.parent_id) FILTER (
		WHERE sp.is_active
			AND sp.receives_billing_email
			AND p.status = 'active'
			AND p.email_active
			AND p.email <> ''
	)::integer,
	COUNT(sp.parent_id) FILTER (
		WHERE sp.is_active
			AND sp.receives_billing_email
			AND p.status = 'active'
			AND p.email = ''
	)::integer,
	COUNT(sp.parent_id) FILTER (
		WHERE sp.is_active
			AND sp.receives_billing_email
			AND p.status = 'active'
			AND p.email <> ''
			AND NOT p.email_active
	)::integer,
	COUNT(sp.parent_id) FILTER (
		WHERE NOT sp.is_active
			AND sp.receives_billing_email
	)::integer
FROM students s
LEFT JOIN classes c ON c.id = s.class_id
LEFT JOIN school_years sy ON sy.id = c.school_year_id
LEFT JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN student_parents sp ON sp.student_id = s.id
LEFT JOIN parents p ON p.id = sp.parent_id
WHERE ` + strings.Join(conditions, " AND ") + `
GROUP BY s.id, s.student_code, s.full_name, c.id, c.name, c.grade, sy.id, sy.code, sy.school_id, sc.code
ORDER BY sc.code, sy.code DESC, c.grade, c.name, s.student_code`
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := []adminReadinessIssue{}
	for rows.Next() {
		var row struct {
			StudentID             string
			StudentCode           string
			StudentName           string
			ClassID               string
			ClassName             string
			Grade                 string
			SchoolYearID          string
			SchoolYearCode        string
			SchoolID              string
			SchoolCode            string
			ActiveParentCount     int
			BillingSelectedCount  int
			BillingReadyCount     int
			MissingEmailCount     int
			InactiveEmailCount    int
			InactiveSelectedCount int
		}
		if err := rows.Scan(
			&row.StudentID,
			&row.StudentCode,
			&row.StudentName,
			&row.ClassID,
			&row.ClassName,
			&row.Grade,
			&row.SchoolYearID,
			&row.SchoolYearCode,
			&row.SchoolID,
			&row.SchoolCode,
			&row.ActiveParentCount,
			&row.BillingSelectedCount,
			&row.BillingReadyCount,
			&row.MissingEmailCount,
			&row.InactiveEmailCount,
			&row.InactiveSelectedCount,
		); err != nil {
			return nil, err
		}
		base := adminReadinessIssue{
			EntityType:    "student",
			EntityID:      row.StudentID,
			EntityLabel:   strings.TrimSpace(row.StudentCode + " · " + row.StudentName),
			Scope:         adminReadinessScope(row.SchoolCode, row.SchoolYearCode, row.Grade, row.ClassName, filters.PeriodCode, filters.Month),
			Action:        "students",
			SchoolID:      row.SchoolID,
			SchoolYearID:  row.SchoolYearID,
			ClassID:       row.ClassID,
			Grade:         row.Grade,
			PeriodCode:    filters.PeriodCode,
			Month:         filters.Month,
			ReferenceCode: row.StudentCode,
		}
		if row.ClassID == "" {
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "student_missing_class"
			issue.Message = "Học sinh chưa gắn lớp nên không thể vào workflow hóa đơn."
			issues = append(issues, issue)
		}
		if row.ActiveParentCount == 0 {
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "student_missing_parent"
			issue.Message = "Học sinh chưa có phụ huynh/người giám hộ đang hoạt động."
			issues = append(issues, issue)
			continue
		}
		if row.BillingReadyCount == 0 {
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "student_missing_billing_recipient"
			if row.BillingSelectedCount == 0 {
				issue.Message = "Chưa chọn phụ huynh nhận email học phí."
			} else {
				issue.Message = "Người nhận học phí chưa hợp lệ hoặc chưa có email active."
			}
			issues = append(issues, issue)
		}
		if row.MissingEmailCount > 0 {
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "billing_recipient_missing_email"
			issue.Message = "Có người nhận học phí chưa có email."
			issue.ReferenceCount = row.MissingEmailCount
			issues = append(issues, issue)
		}
		if row.InactiveEmailCount > 0 {
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "billing_recipient_email_inactive"
			issue.Message = "Có người nhận học phí có email đang inactive."
			issue.ReferenceCount = row.InactiveEmailCount
			issues = append(issues, issue)
		}
		if row.InactiveSelectedCount > 0 {
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "inactive_parent_selected_for_billing"
			issue.Message = "Quan hệ phụ huynh inactive vẫn đang được chọn nhận học phí."
			issue.ReferenceCount = row.InactiveSelectedCount
			issues = append(issues, issue)
		}
	}
	return issues, rows.Err()
}

func listAdminClassReadinessIssues(ctx context.Context, db *sql.DB, filters adminFilters) ([]adminReadinessIssue, error) {
	args := []any{}
	addArg := adminReadinessArgFunc(&args)
	conditions := []string{"c.status <> 'archived'"}
	if filters.SchoolID != "" {
		conditions = append(conditions, "sy.school_id = "+addArg(filters.SchoolID)+"::uuid")
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "c.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.ClassID != "" {
		conditions = append(conditions, "c.id = "+addArg(filters.ClassID)+"::uuid")
	}
	if filters.Grade != "" {
		conditions = append(conditions, "lower(c.grade) = lower("+addArg(filters.Grade)+")")
	}
	scheduleConditions := []string{
		"fs.school_year_id = c.school_year_id",
		"fs.status <> 'archived'",
		`(
			fs.class_id = c.id
			OR (fs.scope_type = 'grade' AND fs.class_id IS NULL AND lower(fs.grade) = lower(c.grade))
			OR (fs.scope_type = 'school_year' AND fs.class_id IS NULL AND btrim(fs.grade) = '')
		)`,
	}
	if filters.PeriodCode != "" {
		scheduleConditions = append(scheduleConditions, "lower(fs.period_code) = lower("+addArg(filters.PeriodCode)+")")
	}
	if filters.Month > 0 {
		scheduleConditions = append(scheduleConditions, "fs.month = "+addArg(filters.Month)+"::integer")
	}
	query := `
SELECT
	c.id::text,
	c.name,
	c.grade,
	sy.id::text,
	sy.code,
	sy.school_id::text,
	sc.code,
	COALESCE(student_counts.student_count, 0),
	COALESCE(schedule_counts.current_schedule_count, 0),
	COALESCE(schedule_counts.current_active_schedule_count, 0)
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN LATERAL (
	SELECT COUNT(*)::integer AS student_count
	FROM students s
	WHERE s.class_id = c.id
		AND s.status <> 'inactive'
) student_counts ON true
LEFT JOIN LATERAL (
	SELECT
		COUNT(*)::integer AS current_schedule_count,
		COUNT(*) FILTER (WHERE fs.status = 'active')::integer AS current_active_schedule_count
	FROM fee_schedules fs
	WHERE ` + strings.Join(scheduleConditions, " AND ") + `
) schedule_counts ON true
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY sc.code, sy.code DESC, c.grade, c.name`
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hasSelectedPeriod := filters.PeriodCode != "" || filters.Month > 0
	issues := []adminReadinessIssue{}
	for rows.Next() {
		var row struct {
			ClassID             string
			ClassName           string
			Grade               string
			SchoolYearID        string
			SchoolYearCode      string
			SchoolID            string
			SchoolCode          string
			StudentCount        int
			ScheduleCount       int
			ActiveScheduleCount int
		}
		if err := rows.Scan(
			&row.ClassID,
			&row.ClassName,
			&row.Grade,
			&row.SchoolYearID,
			&row.SchoolYearCode,
			&row.SchoolID,
			&row.SchoolCode,
			&row.StudentCount,
			&row.ScheduleCount,
			&row.ActiveScheduleCount,
		); err != nil {
			return nil, err
		}
		base := adminReadinessIssue{
			EntityType:   "class",
			EntityID:     row.ClassID,
			EntityLabel:  row.ClassName,
			Scope:        adminReadinessScope(row.SchoolCode, row.SchoolYearCode, row.Grade, row.ClassName, filters.PeriodCode, filters.Month),
			SchoolID:     row.SchoolID,
			SchoolYearID: row.SchoolYearID,
			ClassID:      row.ClassID,
			Grade:        row.Grade,
			PeriodCode:   filters.PeriodCode,
			Month:        filters.Month,
		}
		if row.StudentCount == 0 {
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "class_has_no_students"
			issue.Message = "Lớp chưa có học sinh active."
			issue.Action = "students"
			issues = append(issues, issue)
		}
		if hasSelectedPeriod && row.StudentCount > 0 && row.ActiveScheduleCount == 0 {
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "class_missing_fee_schedule"
			if row.ScheduleCount > 0 {
				issue.Message = "Có bảng phí trong kỳ nhưng chưa active cho lớp này."
			} else {
				issue.Message = "Lớp chưa có bảng phí active cho kỳ đang chọn."
			}
			issue.Action = "fees"
			issues = append(issues, issue)
		}
	}
	return issues, rows.Err()
}

func listAdminFeeScheduleReadinessIssues(ctx context.Context, db *sql.DB, filters adminFilters) ([]adminReadinessIssue, error) {
	args := []any{}
	addArg := adminReadinessArgFunc(&args)
	conditions := []string{"fs.status = 'active'"}
	if filters.SchoolID != "" {
		conditions = append(conditions, "sy.school_id = "+addArg(filters.SchoolID)+"::uuid")
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "fs.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.ClassID != "" {
		classArg := addArg(filters.ClassID)
		conditions = append(conditions, `(fs.class_id = `+classArg+`::uuid
			OR (fs.scope_type = 'grade' AND EXISTS (
				SELECT 1 FROM classes fc
				WHERE fc.id = `+classArg+`::uuid
					AND fc.school_year_id = fs.school_year_id
					AND lower(fc.grade) = lower(fs.grade)
			))
			OR (fs.scope_type = 'school_year' AND fs.class_id IS NULL AND btrim(fs.grade) = ''))`)
	}
	if filters.Grade != "" {
		gradeArg := addArg(filters.Grade)
		conditions = append(conditions, `(lower(COALESCE(c.grade, fs.grade)) = lower(`+gradeArg+`)
			OR (fs.scope_type = 'school_year' AND fs.class_id IS NULL AND btrim(fs.grade) = ''))`)
	}
	if filters.PeriodCode != "" {
		conditions = append(conditions, "lower(fs.period_code) = lower("+addArg(filters.PeriodCode)+")")
	}
	if filters.Month > 0 {
		conditions = append(conditions, "fs.month = "+addArg(filters.Month)+"::integer")
	}
	query := `
SELECT
	fs.id::text,
	fs.name,
	fs.period_code,
	COALESCE(fs.month, 0),
	fs.scope_type,
	COALESCE(c.id::text, ''),
	COALESCE(c.name, ''),
	COALESCE(c.grade, fs.grade),
	sy.id::text,
	sy.code,
	sy.school_id::text,
	sc.code,
	COUNT(fsi.id)::integer,
	COUNT(fsi.id) FILTER (WHERE fsi.amount = 0)::integer
FROM fee_schedules fs
JOIN school_years sy ON sy.id = fs.school_year_id
JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN classes c ON c.id = fs.class_id
LEFT JOIN fee_schedule_items fsi ON fsi.schedule_id = fs.id
WHERE ` + strings.Join(conditions, " AND ") + `
GROUP BY fs.id, fs.name, fs.period_code, fs.month, fs.scope_type, c.id, c.name, c.grade, sy.id, sy.code, sy.school_id, sc.code
ORDER BY sc.code, sy.code DESC, fs.period_code DESC, fs.name`
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := []adminReadinessIssue{}
	for rows.Next() {
		var row struct {
			ScheduleID     string
			ScheduleName   string
			PeriodCode     string
			Month          int
			ScopeType      string
			ClassID        string
			ClassName      string
			Grade          string
			SchoolYearID   string
			SchoolYearCode string
			SchoolID       string
			SchoolCode     string
			ItemCount      int
			ZeroItemCount  int
		}
		if err := rows.Scan(
			&row.ScheduleID,
			&row.ScheduleName,
			&row.PeriodCode,
			&row.Month,
			&row.ScopeType,
			&row.ClassID,
			&row.ClassName,
			&row.Grade,
			&row.SchoolYearID,
			&row.SchoolYearCode,
			&row.SchoolID,
			&row.SchoolCode,
			&row.ItemCount,
			&row.ZeroItemCount,
		); err != nil {
			return nil, err
		}
		scopeClass := row.ClassName
		if scopeClass == "" && row.ScopeType == "school_year" {
			scopeClass = "Toàn năm học"
		}
		base := adminReadinessIssue{
			EntityType:    "fee_schedule",
			EntityID:      row.ScheduleID,
			EntityLabel:   firstNonEmpty(row.ScheduleName, row.PeriodCode, row.ScheduleID),
			Scope:         adminReadinessScope(row.SchoolCode, row.SchoolYearCode, row.Grade, scopeClass, row.PeriodCode, row.Month),
			Action:        "fees",
			SchoolID:      row.SchoolID,
			SchoolYearID:  row.SchoolYearID,
			ClassID:       row.ClassID,
			Grade:         row.Grade,
			PeriodCode:    row.PeriodCode,
			Month:         row.Month,
			ReferenceCode: row.PeriodCode,
		}
		if row.ItemCount == 0 {
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "fee_schedule_empty_items"
			issue.Message = "Bảng phí active chưa có dòng phí."
			issues = append(issues, issue)
			continue
		}
		if row.ZeroItemCount > 0 {
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "fee_schedule_zero_amount_items"
			issue.Message = "Bảng phí có dòng phí amount = 0, cần kiểm tra trước khi sinh hóa đơn."
			issue.ReferenceCount = row.ZeroItemCount
			issues = append(issues, issue)
		}
	}
	return issues, rows.Err()
}

func listAdminAdjustmentReadinessIssues(ctx context.Context, db *sql.DB, filters adminFilters) ([]adminReadinessIssue, error) {
	args := []any{}
	addArg := adminReadinessArgFunc(&args)
	conditions := []string{"sfa.status = 'active'", "btrim(sfa.reason) = ''"}
	if filters.SchoolID != "" {
		conditions = append(conditions, "sy.school_id = "+addArg(filters.SchoolID)+"::uuid")
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "fs.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.ClassID != "" {
		conditions = append(conditions, "s.class_id = "+addArg(filters.ClassID)+"::uuid")
	}
	if filters.Grade != "" {
		conditions = append(conditions, "lower(c.grade) = lower("+addArg(filters.Grade)+")")
	}
	if filters.PeriodCode != "" {
		conditions = append(conditions, "lower(fs.period_code) = lower("+addArg(filters.PeriodCode)+")")
	}
	if filters.Month > 0 {
		conditions = append(conditions, "fs.month = "+addArg(filters.Month)+"::integer")
	}
	query := `
SELECT
	sfa.id::text,
	s.student_code,
	s.full_name,
	c.id::text,
	c.name,
	c.grade,
	sy.id::text,
	sy.code,
	sy.school_id::text,
	sc.code,
	fs.period_code,
	COALESCE(fs.month, 0)
FROM student_fee_adjustments sfa
JOIN fee_schedules fs ON fs.id = sfa.schedule_id
JOIN students s ON s.id = sfa.student_id
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY sc.code, sy.code DESC, c.grade, c.name, s.student_code`
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := []adminReadinessIssue{}
	for rows.Next() {
		var row struct {
			AdjustmentID   string
			StudentCode    string
			StudentName    string
			ClassID        string
			ClassName      string
			Grade          string
			SchoolYearID   string
			SchoolYearCode string
			SchoolID       string
			SchoolCode     string
			PeriodCode     string
			Month          int
		}
		if err := rows.Scan(
			&row.AdjustmentID,
			&row.StudentCode,
			&row.StudentName,
			&row.ClassID,
			&row.ClassName,
			&row.Grade,
			&row.SchoolYearID,
			&row.SchoolYearCode,
			&row.SchoolID,
			&row.SchoolCode,
			&row.PeriodCode,
			&row.Month,
		); err != nil {
			return nil, err
		}
		issues = append(issues, adminReadinessIssue{
			Severity:      adminReadinessSeverityWarning,
			Type:          "student_adjustment_missing_reason",
			EntityType:    "student_adjustment",
			EntityID:      row.AdjustmentID,
			EntityLabel:   strings.TrimSpace(row.StudentCode + " · " + row.StudentName),
			Scope:         adminReadinessScope(row.SchoolCode, row.SchoolYearCode, row.Grade, row.ClassName, row.PeriodCode, row.Month),
			Message:       "Điều chỉnh học phí thiếu lý do.",
			Action:        "fees",
			SchoolID:      row.SchoolID,
			SchoolYearID:  row.SchoolYearID,
			ClassID:       row.ClassID,
			Grade:         row.Grade,
			PeriodCode:    row.PeriodCode,
			Month:         row.Month,
			ReferenceCode: row.StudentCode,
		})
	}
	return issues, rows.Err()
}

func listAdminNotificationReadinessIssues(ctx context.Context, db *sql.DB, filters adminFilters) ([]adminReadinessIssue, error) {
	args := []any{}
	addArg := adminReadinessArgFunc(&args)
	conditions := []string{"nl.status = 'error'"}
	if filters.SchoolID != "" {
		conditions = append(conditions, "sy.school_id = "+addArg(filters.SchoolID)+"::uuid")
	}
	if filters.SchoolYearID != "" {
		conditions = append(conditions, "i.school_year_id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.ClassID != "" {
		conditions = append(conditions, "i.class_id = "+addArg(filters.ClassID)+"::uuid")
	}
	if filters.Grade != "" {
		conditions = append(conditions, "lower(i.grade) = lower("+addArg(filters.Grade)+")")
	}
	if filters.PeriodCode != "" {
		conditions = append(conditions, "lower(i.period_code) = lower("+addArg(filters.PeriodCode)+")")
	}
	if filters.Month > 0 {
		conditions = append(conditions, "i.month = "+addArg(filters.Month)+"::integer")
	}
	if filters.Status != "" {
		conditions = append(conditions, "i.status = "+addArg(filters.Status))
	}
	query := `
SELECT
	nc.id::text,
	nc.code,
	COALESCE(nc.period_code, ''),
	COUNT(*)::integer,
	MAX(nl.sent_at)
FROM notification_logs nl
JOIN notification_campaigns nc ON nc.id = nl.campaign_id
JOIN invoices i ON i.id = nl.invoice_id
JOIN school_years sy ON sy.id = i.school_year_id
WHERE ` + strings.Join(conditions, " AND ") + `
GROUP BY nc.id, nc.code, nc.period_code
ORDER BY MAX(nl.sent_at) DESC
LIMIT 30`
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := []adminReadinessIssue{}
	for rows.Next() {
		var campaignID, campaignCode, periodCode string
		var count int
		var lastSent time.Time
		if err := rows.Scan(&campaignID, &campaignCode, &periodCode, &count, &lastSent); err != nil {
			return nil, err
		}
		issues = append(issues, adminReadinessIssue{
			Severity:       adminReadinessSeverityWarning,
			Type:           "notification_failed_recipients",
			EntityType:     "notification_campaign",
			EntityID:       campaignID,
			EntityLabel:    campaignCode,
			Scope:          adminReadinessScope("", "", filters.Grade, "", firstNonEmpty(filters.PeriodCode, periodCode), filters.Month),
			Message:        fmt.Sprintf("Campaign có %d người nhận gửi lỗi.", count),
			Action:         "notify",
			SchoolID:       filters.SchoolID,
			SchoolYearID:   filters.SchoolYearID,
			ClassID:        filters.ClassID,
			Grade:          filters.Grade,
			PeriodCode:     firstNonEmpty(filters.PeriodCode, periodCode),
			Month:          filters.Month,
			ReferenceCode:  campaignCode,
			ReferenceCount: count,
		})
	}
	return issues, rows.Err()
}

func listAdminOperationReadinessIssues(ctx context.Context, db *sql.DB, _ adminFilters) ([]adminReadinessIssue, error) {
	var count int
	err := db.QueryRowContext(ctx, `
SELECT COUNT(*)::integer
FROM operation_logs
WHERE source = 'background_job'
	AND operation LIKE 'email.cron.%'
	AND level = 'error'
	AND occurred_at >= now() - interval '7 days'`).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	return []adminReadinessIssue{{
		Severity:       adminReadinessSeverityWarning,
		Type:           "cron_recent_errors",
		EntityType:     "email_cron",
		EntityLabel:    "Email cron",
		Scope:          "7 ngày gần nhất",
		Message:        fmt.Sprintf("Cron email có %d lỗi gần đây trong operation logs.", count),
		Action:         "operations",
		ReferenceCount: count,
	}}, nil
}

func appendAdminInvoiceReadinessIssues(issues []adminReadinessIssue, invoices []adminInvoiceReportRow) []adminReadinessIssue {
	for _, invoice := range invoices {
		base := adminReadinessIssue{
			EntityType:    "invoice",
			EntityID:      invoice.ID,
			EntityLabel:   firstNonEmpty(invoice.InvoiceCode, invoice.StudentCode),
			Scope:         adminReadinessScope("", invoice.SchoolYearCode, invoice.Grade, invoice.ClassName, invoice.PeriodCode, invoice.Month),
			Action:        "invoices",
			SchoolYearID:  invoice.SchoolYearID,
			ClassID:       invoice.ClassID,
			Grade:         invoice.Grade,
			PeriodCode:    invoice.PeriodCode,
			Month:         invoice.Month,
			ReferenceCode: invoice.InvoiceCode,
		}
		if strings.TrimSpace(invoice.CollectionBankBIN) == "" ||
			strings.TrimSpace(invoice.CollectionBankAccount) == "" ||
			strings.TrimSpace(invoice.QRBillNumber) == "" {
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "invoice_missing_payment_data"
			issue.Message = "Invoice thiếu dữ liệu QR/thanh toán."
			issues = append(issues, issue)
		}
		switch invoice.Status {
		case invoiceStatusUnpaid:
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "invoice_unpaid"
			issue.Message = "Invoice chưa thanh toán."
			issues = append(issues, issue)
		case invoiceStatusPartial:
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "invoice_partial"
			issue.Message = "Invoice mới thanh toán một phần."
			issues = append(issues, issue)
		case invoiceStatusOverpaid:
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "invoice_overpaid"
			issue.Message = "Invoice overpaid, cần kiểm tra trước khi chốt."
			issues = append(issues, issue)
		case invoiceStatusManualReview:
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "invoice_manual_review"
			issue.Message = "Invoice đang cần review thủ công."
			issues = append(issues, issue)
		}
	}
	return issues
}

func appendAdminTransactionReadinessIssues(issues []adminReadinessIssue, transactions []paymentTransactionSummary) []adminReadinessIssue {
	for _, transaction := range transactions {
		if transaction.Direction != "" && transaction.Direction != paymentDirectionIn {
			continue
		}
		base := adminReadinessIssue{
			EntityType:    "payment_transaction",
			EntityID:      transaction.ID,
			EntityLabel:   firstNonEmpty(transaction.ProviderTransactionID, transaction.ReferenceCode, transaction.ID),
			Scope:         firstNonEmpty(transaction.ProviderCode, "Giao dịch tiền vào"),
			Action:        "reconcile",
			ReferenceCode: firstNonEmpty(transaction.ProviderTransactionID, transaction.ReferenceCode),
		}
		switch transaction.Status {
		case paymentTransactionStatusUnmatched:
			issue := base
			issue.Severity = adminReadinessSeverityWarning
			issue.Type = "incoming_transaction_unmatched"
			issue.Message = "Giao dịch tiền vào chưa được khớp invoice."
			issues = append(issues, issue)
		case paymentTransactionStatusManualReview:
			issue := base
			issue.Severity = adminReadinessSeverityBlocking
			issue.Type = "transaction_manual_review"
			issue.Message = "Giao dịch cần review thủ công."
			issues = append(issues, issue)
		}
	}
	return issues
}

func adminLocalReadinessIssues(now time.Time) []adminReadinessIssue {
	issues := []adminReadinessIssue{}
	cfg, cfgErr := loadEmailConfig()
	if cfgErr != nil {
		issues = append(issues, adminReadinessIssue{
			Severity:    adminReadinessSeverityBlocking,
			Type:        "email_provider_not_configured",
			EntityType:  "email_config",
			EntityLabel: "Email provider",
			Scope:       "Email & Cron",
			Message:     "Không đọc được cấu hình email.",
			Action:      "email_config",
		})
	} else if err := validateEmailConfigForSend(cfg); err != nil {
		issues = append(issues, adminReadinessIssue{
			Severity:    adminReadinessSeverityBlocking,
			Type:        "email_provider_not_configured",
			EntityType:  "email_config",
			EntityLabel: "Email provider",
			Scope:       "Email & Cron",
			Message:     "Email provider chưa sẵn sàng để gửi: " + err.Error(),
			Action:      "email_config",
		})
	}

	state, cronErr := loadEmailCronState()
	if cronErr != nil {
		issues = append(issues, adminReadinessIssue{
			Severity:    adminReadinessSeverityWarning,
			Type:        "cron_queue_errors",
			EntityType:  "email_cron",
			EntityLabel: "Email cron",
			Scope:       "Email & Cron",
			Message:     "Không đọc được trạng thái cron email.",
			Action:      "email_config",
		})
		return issues
	}
	cron := emailCronPublic(state, now)
	if cron.Errors > 0 {
		issues = append(issues, adminReadinessIssue{
			Severity:       adminReadinessSeverityWarning,
			Type:           "cron_queue_errors",
			EntityType:     "email_cron",
			EntityLabel:    "Email cron",
			Scope:          "Email & Cron",
			Message:        fmt.Sprintf("Cron queue có %d job lỗi.", cron.Errors),
			Action:         "email_config",
			ReferenceCount: cron.Errors,
		})
	}
	if cron.Queued > 0 && cron.DailyLimit > 0 && cron.SentLast24H >= cron.DailyLimit {
		issues = append(issues, adminReadinessIssue{
			Severity:       adminReadinessSeverityWarning,
			Type:           "cron_over_limit_sends",
			EntityType:     "email_cron",
			EntityLabel:    "Email cron",
			Scope:          "Email & Cron",
			Message:        "Cron queue còn job chờ nhưng quota gửi 24 giờ đã hết.",
			Action:         "email_config",
			ReferenceCount: cron.Queued,
		})
	}
	return issues
}

func adminReadinessScope(schoolCode string, schoolYearCode string, grade string, className string, periodCode string, month int) string {
	parts := []string{}
	if strings.TrimSpace(schoolCode) != "" {
		parts = append(parts, strings.TrimSpace(schoolCode))
	}
	if strings.TrimSpace(schoolYearCode) != "" {
		parts = append(parts, strings.TrimSpace(schoolYearCode))
	}
	if strings.TrimSpace(grade) != "" {
		parts = append(parts, "Khối "+strings.TrimSpace(grade))
	}
	if strings.TrimSpace(className) != "" {
		parts = append(parts, strings.TrimSpace(className))
	}
	if strings.TrimSpace(periodCode) != "" {
		parts = append(parts, "Kỳ "+strings.TrimSpace(periodCode))
	}
	if month > 0 {
		parts = append(parts, fmt.Sprintf("T%d", month))
	}
	if len(parts) == 0 {
		return "Toàn hệ thống"
	}
	return strings.Join(parts, " · ")
}

func adminReadinessArgFunc(args *[]any) func(any) string {
	return func(value any) string {
		*args = append(*args, value)
		return fmt.Sprintf("$%d", len(*args))
	}
}
