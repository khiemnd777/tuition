package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	maxMasterDataImportRows = 1000
	defaultTenantCode       = "DEKISUGI"
)

var classGradePattern = regexp.MustCompile(`\d+`)

type masterDataSchoolOption struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type masterDataSchoolYearOption struct {
	ID         string `json:"id"`
	SchoolID   string `json:"schoolId"`
	SchoolCode string `json:"schoolCode"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

type masterDataClassOption struct {
	ID             string `json:"id"`
	SchoolID       string `json:"schoolId"`
	SchoolCode     string `json:"schoolCode"`
	SchoolYearID   string `json:"schoolYearId"`
	SchoolYearCode string `json:"schoolYearCode"`
	Grade          string `json:"grade"`
	Name           string `json:"name"`
	Status         string `json:"status"`
}

type masterDataOptions struct {
	Schools     []masterDataSchoolOption     `json:"schools"`
	SchoolYears []masterDataSchoolYearOption `json:"schoolYears"`
	Classes     []masterDataClassOption      `json:"classes"`
}

type masterDataStudentListFilters struct {
	TenantID        string
	StudentID       string
	SchoolID        string
	SchoolYearID    string
	SchoolYear      string
	ClassID         string
	Grade           string
	Query           string
	IncludeInactive bool
}

type masterDataStudent struct {
	ID                      string                     `json:"id"`
	StudentCode             string                     `json:"studentCode"`
	StudentName             string                     `json:"studentName"`
	Status                  string                     `json:"status"`
	SchoolID                string                     `json:"schoolId"`
	SchoolCode              string                     `json:"schoolCode"`
	ClassID                 string                     `json:"classId"`
	ClassName               string                     `json:"className"`
	Grade                   string                     `json:"grade"`
	SchoolYearID            string                     `json:"schoolYearId"`
	SchoolYearCode          string                     `json:"schoolYearCode"`
	ParentCount             int                        `json:"parentCount"`
	BillingRecipientCount   int                        `json:"billingRecipientCount"`
	MissingBillingRecipient bool                       `json:"missingBillingRecipient"`
	ContactWarning          string                     `json:"contactWarning,omitempty"`
	InvoiceAttentionCount   int                        `json:"invoiceAttentionCount"`
	Parents                 []masterDataParentContact  `json:"parents"`
	Siblings                []masterDataStudentSibling `json:"siblings,omitempty"`
}

type masterDataParentContact struct {
	ID                   string `json:"id"`
	ParentName           string `json:"parentName"`
	Email                string `json:"email"`
	Phone                string `json:"phone,omitempty"`
	Relationship         string `json:"relationship"`
	EmailActive          bool   `json:"emailActive"`
	IsPrimary            bool   `json:"isPrimary"`
	IsActive             bool   `json:"isActive"`
	ReceivesBillingEmail bool   `json:"receivesBillingEmail"`
	BillingReady         bool   `json:"billingReady"`
}

type masterDataStudentSibling struct {
	ID                string   `json:"id"`
	StudentCode       string   `json:"studentCode"`
	StudentName       string   `json:"studentName"`
	ClassName         string   `json:"className"`
	Grade             string   `json:"grade"`
	SchoolYearCode    string   `json:"schoolYearCode"`
	SharedParentNames []string `json:"sharedParentNames"`
}

type masterDataStudentSaveInput struct {
	ID          string                      `json:"id,omitempty"`
	StudentCode string                      `json:"studentCode"`
	StudentName string                      `json:"studentName"`
	ClassID     string                      `json:"classId"`
	Status      string                      `json:"status"`
	Parents     []masterDataParentSaveInput `json:"parents"`
}

type masterDataParentSaveInput struct {
	ID                   string `json:"id,omitempty"`
	ParentName           string `json:"parentName"`
	Email                string `json:"email"`
	Phone                string `json:"phone,omitempty"`
	Relationship         string `json:"relationship,omitempty"`
	EmailActive          *bool  `json:"emailActive,omitempty"`
	IsPrimary            *bool  `json:"isPrimary,omitempty"`
	IsActive             *bool  `json:"isActive,omitempty"`
	ReceivesBillingEmail *bool  `json:"receivesBillingEmail,omitempty"`
}

type masterDataImportRow struct {
	RowNumber            int    `json:"rowNumber"`
	StudentCode          string `json:"studentCode"`
	StudentName          string `json:"studentName"`
	SchoolCode           string `json:"schoolCode"`
	SchoolYearCode       string `json:"schoolYearCode"`
	Grade                string `json:"grade"`
	ClassName            string `json:"className"`
	ParentName           string `json:"parentName"`
	ParentEmail          string `json:"parentEmail"`
	ParentPhone          string `json:"parentPhone"`
	Relationship         string `json:"relationship"`
	IsPrimary            bool   `json:"isPrimary"`
	ParentActive         bool   `json:"parentActive"`
	ReceivesBillingEmail bool   `json:"receivesBillingEmail"`
}

type masterDataImportRowResult struct {
	RowNumber      int    `json:"rowNumber"`
	StudentCode    string `json:"studentCode"`
	StudentName    string `json:"studentName"`
	SchoolCode     string `json:"schoolCode"`
	SchoolYearCode string `json:"schoolYearCode"`
	Grade          string `json:"grade"`
	ClassName      string `json:"className"`
	ParentName     string `json:"parentName"`
	ParentEmail    string `json:"parentEmail"`
	ParentPhone    string `json:"parentPhone"`
	Relationship   string `json:"relationship"`
	Action         string `json:"action"`
}

type masterDataImportIssue struct {
	RowNumber   int    `json:"rowNumber"`
	StudentCode string `json:"studentCode"`
	Type        string `json:"type"`
	Message     string `json:"message"`
	Existing    string `json:"existing,omitempty"`
	Incoming    string `json:"incoming,omitempty"`
}

type masterDataImportSummary struct {
	TotalRows int `json:"totalRows"`
	Ready     int `json:"ready"`
	Created   int `json:"created"`
	Unchanged int `json:"unchanged"`
	Conflicts int `json:"conflicts"`
}

type masterDataImportResponse struct {
	Applied bool                        `json:"applied"`
	Summary masterDataImportSummary     `json:"summary"`
	Rows    []masterDataImportRowResult `json:"rows"`
	Issues  []masterDataImportIssue     `json:"issues"`
	Options *masterDataOptions          `json:"options,omitempty"`
}

type masterDataStudentExisting struct {
	ID             string
	StudentCode    string
	StudentName    string
	SchoolID       string
	SchoolCode     string
	ClassID        string
	ClassName      string
	Grade          string
	SchoolYearID   string
	SchoolYearCode string
}

type masterDataParentExisting struct {
	ID          string
	ParentName  string
	Email       string
	Phone       string
	EmailActive bool
}

type masterDataParentLinkExisting struct {
	IsPrimary            bool
	IsActive             bool
	ReceivesBillingEmail bool
}

type masterDataExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type masterDataSaveError struct {
	Status  int
	Message string
}

func (err *masterDataSaveError) Error() string {
	return err.Message
}

func handleMasterDataOptions(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	options, err := listMasterDataOptions(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load master data options", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, options)
}

func handleMasterDataStudents(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	query := r.URL.Query()
	students, err := listMasterDataStudents(r.Context(), db, masterDataStudentListFilters{
		TenantID:     tenantID,
		SchoolID:     strings.TrimSpace(query.Get("schoolId")),
		SchoolYearID: strings.TrimSpace(query.Get("schoolYearId")),
		SchoolYear:   normalizeSchoolYearCode(query.Get("schoolYear")),
		ClassID:      strings.TrimSpace(query.Get("classId")),
		Grade:        normalizeGrade(query.Get("grade")),
		Query:        strings.TrimSpace(query.Get("q")),
	})
	if err != nil {
		http.Error(w, "cannot load students", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"students": students})
}

func handleMasterDataStudentSave(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input masterDataStudentSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeMasterDataStudentSaveInput(input)
	if err := validateMasterDataStudentSaveInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	student, err := saveMasterDataStudent(r.Context(), db, tenantID, input, auditContextFromRequest(r))
	if err != nil {
		var saveErr *masterDataSaveError
		if errors.As(err, &saveErr) {
			http.Error(w, saveErr.Message, saveErr.Status)
			return
		}
		var usageErr *tenantUsageLimitError
		if errors.As(err, &usageErr) {
			http.Error(w, usageErr.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "cannot save student", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"student": student})
}

func handleMasterDataImportCSV(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)

	req, err := readImportFileRequest(r, "master_data")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rows, err := parseMasterDataRows(req.Table, req.Mapping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	apply := parseBoolWithDefault(r.URL.Query().Get("apply"), false)
	response, err := importMasterDataRows(r.Context(), db, rows, apply, tenantID)
	if err != nil {
		var usageErr *tenantUsageLimitError
		if errors.As(err, &usageErr) {
			http.Error(w, usageErr.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "cannot import master data", http.StatusInternalServerError)
		return
	}
	status := http.StatusOK
	if apply && len(response.Issues) > 0 {
		status = http.StatusConflict
	}
	writeJSON(w, status, response)
}

func openMasterDataDatabase(ctx context.Context) (*sql.DB, error) {
	cfg, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	if err := cfg.requireURL(); err != nil {
		return nil, err
	}
	db, err := openConfiguredDatabase(ctx, cfg)
	if err != nil {
		return nil, errors.New("database connection failed; check DEKISUGI database environment and server")
	}
	return db, nil
}

func writeMasterDataDBError(w http.ResponseWriter, err error) {
	http.Error(w, "master data database unavailable: "+err.Error(), http.StatusServiceUnavailable)
}

func listMasterDataOptions(ctx context.Context, db *sql.DB, tenantID string) (masterDataOptions, error) {
	options := masterDataOptions{
		Schools:     []masterDataSchoolOption{},
		SchoolYears: []masterDataSchoolYearOption{},
		Classes:     []masterDataClassOption{},
	}

	schoolRows, err := db.QueryContext(ctx, `
SELECT id::text, code, name, status
FROM schools
WHERE tenant_id = $1::uuid
ORDER BY code`, tenantID)
	if err != nil {
		return options, err
	}
	defer schoolRows.Close()
	for schoolRows.Next() {
		var item masterDataSchoolOption
		if err := schoolRows.Scan(&item.ID, &item.Code, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.Schools = append(options.Schools, item)
	}
	if err := schoolRows.Err(); err != nil {
		return options, err
	}

	schoolYearRows, err := db.QueryContext(ctx, `
SELECT sy.id::text, sy.school_id::text, sc.code, sy.code, sy.name, sy.status
FROM school_years sy
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
ORDER BY sc.code, sy.code DESC`, tenantID)
	if err != nil {
		return options, err
	}
	defer schoolYearRows.Close()
	for schoolYearRows.Next() {
		var item masterDataSchoolYearOption
		if err := schoolYearRows.Scan(&item.ID, &item.SchoolID, &item.SchoolCode, &item.Code, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.SchoolYears = append(options.SchoolYears, item)
	}
	if err := schoolYearRows.Err(); err != nil {
		return options, err
	}

	classRows, err := db.QueryContext(ctx, `
SELECT c.id::text, sy.school_id::text, sc.code, c.school_year_id::text, sy.code, c.grade, c.name, c.status
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
ORDER BY sc.code, sy.code DESC, c.grade, c.name`, tenantID)
	if err != nil {
		return options, err
	}
	defer classRows.Close()
	for classRows.Next() {
		var item masterDataClassOption
		if err := classRows.Scan(&item.ID, &item.SchoolID, &item.SchoolCode, &item.SchoolYearID, &item.SchoolYearCode, &item.Grade, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.Classes = append(options.Classes, item)
	}
	if err := classRows.Err(); err != nil {
		return options, err
	}

	return options, nil
}

func listMasterDataStudents(ctx context.Context, exec masterDataExecutor, filters masterDataStudentListFilters) ([]masterDataStudent, error) {
	args := []any{}
	conditions := []string{}
	if !filters.IncludeInactive {
		conditions = append(conditions, "s.status <> 'inactive'")
	}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.TenantID != "" {
		conditions = append(conditions, "sc.tenant_id = "+addArg(filters.TenantID)+"::uuid")
	}

	if filters.SchoolYearID != "" {
		conditions = append(conditions, "sy.id = "+addArg(filters.SchoolYearID)+"::uuid")
	}
	if filters.StudentID != "" {
		conditions = append(conditions, "s.id = "+addArg(filters.StudentID)+"::uuid")
	}
	if filters.SchoolID != "" {
		conditions = append(conditions, "sc.id = "+addArg(filters.SchoolID)+"::uuid")
	}
	if filters.SchoolYear != "" {
		conditions = append(conditions, "sy.code = "+addArg(filters.SchoolYear))
	}
	if filters.ClassID != "" {
		conditions = append(conditions, "c.id = "+addArg(filters.ClassID)+"::uuid")
	}
	if filters.Grade != "" {
		conditions = append(conditions, "c.grade = "+addArg(filters.Grade))
	}
	if filters.Query != "" {
		needle := "%" + filters.Query + "%"
		placeholder := addArg(needle)
		conditions = append(conditions, "(s.student_code ILIKE "+placeholder+" OR s.full_name ILIKE "+placeholder+" OR p.full_name ILIKE "+placeholder+" OR p.email ILIKE "+placeholder+" OR p.phone ILIKE "+placeholder+")")
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}

	query := `
SELECT
	s.id::text,
	s.student_code,
	s.full_name,
	s.status,
	sc.id::text,
	sc.code,
	c.id::text,
	c.name,
	c.grade,
	sy.id::text,
	sy.code,
	COALESCE(p.id::text, ''),
	COALESCE(p.full_name, ''),
	COALESCE(p.email, ''),
	COALESCE(p.phone, ''),
	COALESCE(sp.relationship, ''),
	COALESCE(p.email_active, false),
	COALESCE(sp.is_primary, false),
	COALESCE(sp.is_active, false),
	COALESCE(sp.receives_billing_email, false)
FROM students s
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN student_parents sp ON sp.student_id = s.id
LEFT JOIN parents p ON p.id = sp.parent_id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY sc.code, sy.code DESC, c.grade, c.name, s.student_code, sp.is_primary DESC, p.full_name
LIMIT 1000`

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	studentByID := map[string]*masterDataStudent{}
	order := []string{}
	for rows.Next() {
		var student masterDataStudent
		var parent masterDataParentContact
		if err := rows.Scan(
			&student.ID,
			&student.StudentCode,
			&student.StudentName,
			&student.Status,
			&student.SchoolID,
			&student.SchoolCode,
			&student.ClassID,
			&student.ClassName,
			&student.Grade,
			&student.SchoolYearID,
			&student.SchoolYearCode,
			&parent.ID,
			&parent.ParentName,
			&parent.Email,
			&parent.Phone,
			&parent.Relationship,
			&parent.EmailActive,
			&parent.IsPrimary,
			&parent.IsActive,
			&parent.ReceivesBillingEmail,
		); err != nil {
			return nil, err
		}
		existing := studentByID[student.ID]
		if existing == nil {
			student.Parents = []masterDataParentContact{}
			studentByID[student.ID] = &student
			order = append(order, student.ID)
			existing = &student
		}
		if parent.ID != "" {
			existing.Parents = append(existing.Parents, parent)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	students := make([]masterDataStudent, 0, len(order))
	for _, id := range order {
		students = append(students, *studentByID[id])
	}
	return enrichMasterDataStudentRelationships(ctx, exec, filters.TenantID, students)
}

func enrichMasterDataStudentRelationships(ctx context.Context, exec masterDataExecutor, tenantID string, students []masterDataStudent) ([]masterDataStudent, error) {
	studentIDs := make([]string, 0, len(students))
	for idx := range students {
		deriveMasterDataStudentRelationshipState(&students[idx])
		if students[idx].ID != "" {
			studentIDs = append(studentIDs, students[idx].ID)
		}
	}
	if len(studentIDs) == 0 {
		return students, nil
	}

	attentionCounts, err := loadMasterDataInvoiceAttentionCounts(ctx, exec, studentIDs)
	if err != nil {
		return nil, err
	}
	siblingsByStudent, err := loadMasterDataStudentSiblings(ctx, exec, tenantID, studentIDs)
	if err != nil {
		return nil, err
	}
	for idx := range students {
		students[idx].InvoiceAttentionCount = attentionCounts[students[idx].ID]
		students[idx].Siblings = siblingsByStudent[students[idx].ID]
	}
	return students, nil
}

func deriveMasterDataStudentRelationshipState(student *masterDataStudent) {
	student.ParentCount = len(student.Parents)
	student.BillingRecipientCount = 0
	for idx := range student.Parents {
		parent := &student.Parents[idx]
		parent.Relationship = normalizeMasterDataParentRelationship(parent.Relationship)
		parent.BillingReady = masterDataParentBillingReady(*parent)
		if parent.BillingReady {
			student.BillingRecipientCount++
		}
	}
	student.MissingBillingRecipient = student.BillingRecipientCount == 0
	switch {
	case student.ParentCount == 0:
		student.ContactWarning = "missing_parent"
	case student.MissingBillingRecipient:
		student.ContactWarning = "missing_billing_recipient"
	default:
		student.ContactWarning = ""
	}
}

func masterDataParentBillingReady(parent masterDataParentContact) bool {
	return parent.IsActive && parent.ReceivesBillingEmail && parent.EmailActive && strings.TrimSpace(parent.Email) != ""
}

func loadMasterDataInvoiceAttentionCounts(ctx context.Context, exec masterDataExecutor, studentIDs []string) (map[string]int, error) {
	counts := map[string]int{}
	placeholders, args := masterDataUUIDPlaceholders(studentIDs)
	rows, err := exec.QueryContext(ctx, fmt.Sprintf(`
SELECT student_id::text, COUNT(*)::integer
FROM invoices
WHERE student_id IN (%s)
	AND status IN ('unpaid', 'partial', 'overpaid', 'manual_review')
GROUP BY student_id`, placeholders), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var studentID string
		var count int
		if err := rows.Scan(&studentID, &count); err != nil {
			return nil, err
		}
		counts[studentID] = count
	}
	return counts, rows.Err()
}

func loadMasterDataStudentSiblings(ctx context.Context, exec masterDataExecutor, tenantID string, studentIDs []string) (map[string][]masterDataStudentSibling, error) {
	siblings := map[string][]masterDataStudentSibling{}
	placeholders, args := masterDataUUIDPlaceholders(studentIDs)
	tenantCondition := ""
	if tenantID != "" {
		args = append(args, tenantID)
		tenantCondition = fmt.Sprintf("AND sc.tenant_id = $%d::uuid", len(args))
	}
	rows, err := exec.QueryContext(ctx, fmt.Sprintf(`
SELECT
	base.student_id::text,
	sibling.id::text,
	sibling.student_code,
	sibling.full_name,
	c.name,
	c.grade,
	sy.code,
	COALESCE(string_agg(DISTINCT p.full_name, E'\x1f'), '')
FROM student_parents base
JOIN student_parents sibling_link ON sibling_link.parent_id = base.parent_id
JOIN parents p ON p.id = base.parent_id
JOIN students sibling ON sibling.id = sibling_link.student_id
JOIN classes c ON c.id = sibling.class_id
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE base.student_id IN (%s)
	AND sibling_link.student_id <> base.student_id
	AND base.is_active
	AND sibling_link.is_active
	AND p.status = 'active'
	AND sibling.status <> 'inactive'
	%s
GROUP BY base.student_id, sibling.id, sibling.student_code, sibling.full_name, c.name, c.grade, sy.code
ORDER BY sibling.student_code`, placeholders, tenantCondition), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var studentID string
		var sibling masterDataStudentSibling
		var sharedParents string
		if err := rows.Scan(
			&studentID,
			&sibling.ID,
			&sibling.StudentCode,
			&sibling.StudentName,
			&sibling.ClassName,
			&sibling.Grade,
			&sibling.SchoolYearCode,
			&sharedParents,
		); err != nil {
			return nil, err
		}
		if sharedParents != "" {
			sibling.SharedParentNames = strings.Split(sharedParents, "\x1f")
		} else {
			sibling.SharedParentNames = []string{}
		}
		siblings[studentID] = append(siblings[studentID], sibling)
	}
	return siblings, rows.Err()
}

func masterDataUUIDPlaceholders(ids []string) (string, []any) {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for idx, id := range ids {
		placeholders[idx] = fmt.Sprintf("$%d::uuid", idx+1)
		args[idx] = id
	}
	return strings.Join(placeholders, ", "), args
}

func parseMasterDataCSVRows(input io.Reader) ([]masterDataImportRow, error) {
	return parseMasterDataCSVRowsWithMapping(input, nil)
}

func parseMasterDataCSVRowsWithMapping(input io.Reader, mapping map[string]string) ([]masterDataImportRow, error) {
	table, err := readCSVTable(input)
	if err != nil {
		return nil, err
	}
	return parseMasterDataRows(table, normalizeImportMapping("master_data", mapping))
}

func parseMasterDataRows(table importTable, mapping map[string]string) ([]masterDataImportRow, error) {
	rows := make([]masterDataImportRow, 0, len(table.Records))
	for idx, record := range table.Records {
		if isBlankRecord(record) {
			continue
		}
		parentEmail := normalizeEmail(importFieldValue(record, table, mapping, "parent_email", masterDataCSVAliases("parent_email")))
		row := masterDataImportRow{
			RowNumber:            idx + 2,
			StudentCode:          normalizeStudentCode(importFieldValue(record, table, mapping, "student_code", masterDataCSVAliases("student_code"))),
			StudentName:          strings.TrimSpace(importFieldValue(record, table, mapping, "student", masterDataCSVAliases("student"))),
			SchoolCode:           normalizeSchoolCode(importFieldValue(record, table, mapping, "school", masterDataCSVAliases("school"))),
			SchoolYearCode:       normalizeSchoolYearCode(importFieldValue(record, table, mapping, "school_year", masterDataCSVAliases("school_year"))),
			Grade:                normalizeGrade(importFieldValue(record, table, mapping, "grade", masterDataCSVAliases("grade"))),
			ClassName:            strings.TrimSpace(importFieldValue(record, table, mapping, "class_name", masterDataCSVAliases("class_name"))),
			ParentName:           strings.TrimSpace(importFieldValue(record, table, mapping, "parent", masterDataCSVAliases("parent"))),
			ParentEmail:          parentEmail,
			ParentPhone:          normalizeAdminPhone(importFieldValue(record, table, mapping, "parent_phone", masterDataCSVAliases("parent_phone"))),
			Relationship:         normalizeMasterDataParentRelationship(importFieldValue(record, table, mapping, "relationship", masterDataCSVAliases("relationship"))),
			IsPrimary:            parseBoolWithDefault(importFieldValue(record, table, mapping, "parent_primary", masterDataCSVAliases("parent_primary")), true),
			ParentActive:         parseBoolWithDefault(importFieldValue(record, table, mapping, "parent_active", masterDataCSVAliases("parent_active")), true),
			ReceivesBillingEmail: parseBoolWithDefault(importFieldValue(record, table, mapping, "receives_billing_email", masterDataCSVAliases("receives_billing_email")), parentEmail != ""),
		}
		if row.Grade == "" {
			row.Grade = deriveGradeFromClass(row.ClassName)
		}
		rows = append(rows, row)
	}

	if len(rows) > maxMasterDataImportRows {
		return nil, fmt.Errorf("too many rows, max is %d", maxMasterDataImportRows)
	}
	return rows, nil
}

func masterDataCSVAliases(canonical string) []string {
	switch canonical {
	case "student_code":
		return []string{"student_code", "student_id", "ma_hoc_sinh", "ma_hs", "mhs"}
	case "student":
		return []string{"student_name", "student", "ten_hoc_sinh", "hoc_sinh", "ho_va_ten", "ho_ten", "ten_hs"}
	case "school":
		return []string{"school", "school_code", "ma_truong", "truong", "campus", "co_so"}
	case "school_year":
		return []string{"school_year", "academic_year", "nam_hoc", "nien_khoa"}
	case "grade":
		return []string{"grade", "khoi", "khoi_lop"}
	case "class_name":
		return []string{"class_name", "class", "lop", "ten_lop", "lop_hoc"}
	case "parent":
		return []string{"parent_name", "parent", "ten_phu_huynh", "phu_huynh", "ten_ba_me", "ba_me", "ten_bo_me"}
	case "parent_email":
		return []string{"parent_email", "billing_email", "email_phu_huynh", "email", "mail"}
	case "parent_phone":
		return []string{"parent_phone", "phone", "sdt_phu_huynh", "so_dien_thoai", "dien_thoai"}
	case "relationship":
		return []string{"relationship", "parent_relationship", "guardian_relationship", "quan_he", "moi_quan_he", "vai_tro"}
	case "parent_primary":
		return []string{"parent_primary", "is_primary", "primary", "lien_he_chinh", "chinh"}
	case "parent_active":
		return []string{"parent_active", "email_active", "active", "dang_hoat_dong"}
	case "receives_billing_email":
		return []string{"receives_billing_email", "receive_billing_email", "billing_enabled", "nhan_email_hoc_phi", "nhan_email"}
	default:
		return []string{canonical}
	}
}

func validateMasterDataImportRows(rows []masterDataImportRow) []masterDataImportIssue {
	issues := []masterDataImportIssue{}
	studentByCode := map[string]masterDataImportRow{}
	parentByEmail := map[string]masterDataImportRow{}
	linkByStudentParent := map[string]masterDataImportRow{}
	primaryByStudent := map[string]masterDataImportRow{}

	for _, row := range rows {
		if row.StudentCode == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber: row.RowNumber,
				Type:      "missing_student_code",
				Message:   "student_code is required for production master data",
			})
		}
		if row.StudentName == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_student_name",
				Message:     "student_name is required",
			})
		}
		if row.SchoolYearCode == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_school_year",
				Message:     "school_year is required",
			})
		}
		if row.ClassName == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_class",
				Message:     "class_name is required",
			})
		}
		if row.Grade == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_grade",
				Message:     "grade is required or must be derivable from class_name",
			})
		}
		if row.ParentName == "" && row.ParentEmail == "" && row.ParentPhone == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_parent",
				Message:     "parent_name, parent_email, or parent_phone is required",
			})
		}
		if row.ReceivesBillingEmail && row.ParentEmail == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_parent_email_for_billing",
				Message:     "parent_email is required when receives_billing_email is true",
			})
		}

		if row.StudentCode != "" {
			studentKey := strings.ToLower(row.StudentCode)
			if previous, ok := studentByCode[studentKey]; ok && !sameImportStudent(previous, row) {
				issues = append(issues, masterDataImportIssue{
					RowNumber:   row.RowNumber,
					StudentCode: row.StudentCode,
					Type:        "student_code_conflict_in_csv",
					Message:     "same student_code appears with different student, class, or school-year data",
					Existing:    importStudentDescription(previous),
					Incoming:    importStudentDescription(row),
				})
			} else if !ok {
				studentByCode[studentKey] = row
			}
		}

		if row.ParentEmail != "" {
			if previous, ok := parentByEmail[row.ParentEmail]; ok && !equalText(previous.ParentName, row.ParentName) {
				issues = append(issues, masterDataImportIssue{
					RowNumber:   row.RowNumber,
					StudentCode: row.StudentCode,
					Type:        "parent_email_conflict_in_csv",
					Message:     "same parent_email appears with different parent names",
					Existing:    previous.ParentName,
					Incoming:    row.ParentName,
				})
			} else if !ok {
				parentByEmail[row.ParentEmail] = row
			}
		}

		if row.StudentCode != "" && (row.ParentEmail != "" || row.ParentName != "" || row.ParentPhone != "") {
			linkKey := strings.ToLower(row.StudentCode) + "|" + importParentKey(row)
			if previous, ok := linkByStudentParent[linkKey]; ok {
				issues = append(issues, masterDataImportIssue{
					RowNumber:   row.RowNumber,
					StudentCode: row.StudentCode,
					Type:        "duplicate_parent_link_in_csv",
					Message:     "same student and parent contact appears more than once",
					Existing:    fmt.Sprintf("row %d", previous.RowNumber),
				})
			} else {
				linkByStudentParent[linkKey] = row
			}
		}

		if row.StudentCode != "" && row.IsPrimary && (row.ParentEmail != "" || row.ParentName != "" || row.ParentPhone != "") {
			studentKey := strings.ToLower(row.StudentCode)
			if previous, ok := primaryByStudent[studentKey]; ok && importParentKey(previous) != importParentKey(row) {
				issues = append(issues, masterDataImportIssue{
					RowNumber:   row.RowNumber,
					StudentCode: row.StudentCode,
					Type:        "multiple_primary_parents_in_csv",
					Message:     "only one active primary parent is allowed per student",
					Existing:    importParentDescription(previous),
					Incoming:    importParentDescription(row),
				})
			} else if !ok {
				primaryByStudent[studentKey] = row
			}
		}
	}
	return issues
}

func importMasterDataRows(ctx context.Context, db *sql.DB, rows []masterDataImportRow, apply bool, tenantID string) (masterDataImportResponse, error) {
	if apply {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return masterDataImportResponse{}, err
		}
		defer tx.Rollback()

		response, err := buildMasterDataImportResponse(ctx, tx, rows, tenantID)
		if err != nil || len(response.Issues) > 0 {
			response.Applied = false
			return response, err
		}
		usageDelta, err := masterDataImportUsageDelta(ctx, tx, rows, tenantID)
		if err != nil {
			return masterDataImportResponse{}, err
		}
		if err := enforceTenantUsageLimit(ctx, tx, tenantID, subscriptionMetricSchools, usageDelta[subscriptionMetricSchools], time.Now()); err != nil {
			return masterDataImportResponse{}, err
		}
		if err := enforceTenantUsageLimit(ctx, tx, tenantID, subscriptionMetricStudents, usageDelta[subscriptionMetricStudents], time.Now()); err != nil {
			return masterDataImportResponse{}, err
		}
		if err := applyMasterDataImportRows(ctx, tx, rows, tenantID); err != nil {
			return masterDataImportResponse{}, err
		}
		if err := rebuildTenantUsageCounter(ctx, tx, tenantID, subscriptionMetricSchools, time.Now()); err != nil {
			return masterDataImportResponse{}, err
		}
		if err := rebuildTenantUsageCounter(ctx, tx, tenantID, subscriptionMetricStudents, time.Now()); err != nil {
			return masterDataImportResponse{}, err
		}
		options, err := listMasterDataOptionsTx(ctx, tx, tenantID)
		if err != nil {
			return masterDataImportResponse{}, err
		}
		response.Applied = true
		response.Options = &options
		return response, tx.Commit()
	}

	response, err := buildMasterDataImportResponse(ctx, db, rows, tenantID)
	response.Applied = false
	return response, err
}

func buildMasterDataImportResponse(ctx context.Context, exec masterDataExecutor, rows []masterDataImportRow, tenantID string) (masterDataImportResponse, error) {
	response := masterDataImportResponse{
		Summary: masterDataImportSummary{TotalRows: len(rows)},
		Rows:    []masterDataImportRowResult{},
		Issues:  validateMasterDataImportRows(rows),
	}
	rowHasIssue := map[int]bool{}
	for _, issue := range response.Issues {
		rowHasIssue[issue.RowNumber] = true
	}

	for _, row := range rows {
		result := masterDataImportRowResult{
			RowNumber:      row.RowNumber,
			StudentCode:    row.StudentCode,
			StudentName:    row.StudentName,
			SchoolCode:     schoolCodeOrDefault(row.SchoolCode),
			SchoolYearCode: row.SchoolYearCode,
			Grade:          row.Grade,
			ClassName:      row.ClassName,
			ParentName:     row.ParentName,
			ParentEmail:    row.ParentEmail,
			ParentPhone:    row.ParentPhone,
			Relationship:   row.Relationship,
			Action:         "ready",
		}
		if rowHasIssue[row.RowNumber] {
			result.Action = "conflict"
			response.Rows = append(response.Rows, result)
			continue
		}

		action, issues, err := classifyMasterDataImportRow(ctx, exec, row, tenantID)
		if err != nil {
			return response, err
		}
		result.Action = action
		if len(issues) > 0 {
			result.Action = "conflict"
			rowHasIssue[row.RowNumber] = true
			response.Issues = append(response.Issues, issues...)
		}
		response.Rows = append(response.Rows, result)
	}

	for _, result := range response.Rows {
		if rowHasIssue[result.RowNumber] {
			continue
		}
		response.Summary.Ready++
		switch result.Action {
		case "unchanged":
			response.Summary.Unchanged++
		default:
			response.Summary.Created++
		}
	}
	sortMasterDataImportIssues(response.Issues)
	response.Summary.Conflicts = len(response.Issues)
	return response, nil
}

func classifyMasterDataImportRow(ctx context.Context, exec masterDataExecutor, row masterDataImportRow, tenantID string) (string, []masterDataImportIssue, error) {
	issues := []masterDataImportIssue{}
	student, err := findMasterDataStudentByCode(ctx, exec, row.StudentCode, tenantID)
	if err != nil {
		return "", nil, err
	}
	if student != nil {
		if !equalText(student.StudentName, row.StudentName) {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "student_name_conflict",
				Message:     "student_code already exists with a different student name",
				Existing:    student.StudentName,
				Incoming:    row.StudentName,
			})
		}
		if !equalText(student.SchoolCode, schoolCodeOrDefault(row.SchoolCode)) || !equalText(student.SchoolYearCode, row.SchoolYearCode) || !equalText(student.Grade, row.Grade) || !equalText(student.ClassName, row.ClassName) {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "student_class_conflict",
				Message:     "student_code already exists in a different school, class, grade, or school year",
				Existing:    fmt.Sprintf("%s / %s / %s / %s", student.SchoolCode, student.SchoolYearCode, student.Grade, student.ClassName),
				Incoming:    fmt.Sprintf("%s / %s / %s / %s", schoolCodeOrDefault(row.SchoolCode), row.SchoolYearCode, row.Grade, row.ClassName),
			})
		}
	}

	parent, err := findMasterDataParentByEmail(ctx, exec, row.ParentEmail, tenantID)
	if err != nil {
		return "", nil, err
	}
	if parent != nil && row.ParentName != "" && !equalText(parent.ParentName, row.ParentName) {
		issues = append(issues, masterDataImportIssue{
			RowNumber:   row.RowNumber,
			StudentCode: row.StudentCode,
			Type:        "parent_email_conflict",
			Message:     "parent_email already exists with a different parent name",
			Existing:    parent.ParentName,
			Incoming:    row.ParentName,
		})
	}

	if student == nil {
		return "create", issues, nil
	}
	if parent == nil && row.ParentEmail == "" {
		parent, err = findMasterDataLinkedParentByName(ctx, exec, student.ID, row.ParentName, tenantID)
		if err != nil {
			return "", nil, err
		}
	}
	if parent == nil {
		if row.IsPrimary {
			primary, err := findMasterDataPrimaryParent(ctx, exec, student.ID)
			if err != nil {
				return "", nil, err
			}
			if primary != nil {
				issues = append(issues, masterDataImportIssue{
					RowNumber:   row.RowNumber,
					StudentCode: row.StudentCode,
					Type:        "student_primary_parent_conflict",
					Message:     "student already has a different active primary parent",
					Existing:    importExistingParentDescription(*primary),
					Incoming:    importParentDescription(row),
				})
			}
		}
		return "link_parent", issues, nil
	}

	link, err := findMasterDataStudentParentLink(ctx, exec, student.ID, parent.ID)
	if err != nil {
		return "", nil, err
	}
	if link == nil {
		if row.IsPrimary {
			primary, err := findMasterDataPrimaryParent(ctx, exec, student.ID)
			if err != nil {
				return "", nil, err
			}
			if primary != nil {
				issues = append(issues, masterDataImportIssue{
					RowNumber:   row.RowNumber,
					StudentCode: row.StudentCode,
					Type:        "student_primary_parent_conflict",
					Message:     "student already has a different active primary parent",
					Existing:    importExistingParentDescription(*primary),
					Incoming:    importParentDescription(row),
				})
			}
		}
		return "link_parent", issues, nil
	}

	if link.IsPrimary != row.IsPrimary || link.IsActive != row.ParentActive || link.ReceivesBillingEmail != row.ReceivesBillingEmail {
		issues = append(issues, masterDataImportIssue{
			RowNumber:   row.RowNumber,
			StudentCode: row.StudentCode,
			Type:        "parent_delivery_flags_conflict",
			Message:     "student-parent delivery flags differ from existing data",
			Existing:    formatParentFlags(link.IsPrimary, link.IsActive, link.ReceivesBillingEmail),
			Incoming:    formatParentFlags(row.IsPrimary, row.ParentActive, row.ReceivesBillingEmail),
		})
	}

	return "unchanged", issues, nil
}

func masterDataImportUsageDelta(ctx context.Context, exec masterDataExecutor, rows []masterDataImportRow, tenantID string) (map[string]int, error) {
	schoolCodes := []string{}
	studentCodes := []string{}
	schoolSeen := map[string]bool{}
	studentSeen := map[string]bool{}
	for _, row := range rows {
		schoolCode := schoolCodeOrDefault(row.SchoolCode)
		if schoolCode != "" && !schoolSeen[schoolCode] {
			schoolSeen[schoolCode] = true
			schoolCodes = append(schoolCodes, schoolCode)
		}
		if row.StudentCode != "" && !studentSeen[row.StudentCode] {
			studentSeen[row.StudentCode] = true
			studentCodes = append(studentCodes, row.StudentCode)
		}
	}
	existingSchools, err := masterDataExistingCodeSet(ctx, exec, tenantID, "schools", "code", schoolCodes)
	if err != nil {
		return nil, err
	}
	existingStudents, err := masterDataExistingCodeSet(ctx, exec, tenantID, "students", "student_code", studentCodes)
	if err != nil {
		return nil, err
	}
	delta := map[string]int{
		subscriptionMetricSchools:  0,
		subscriptionMetricStudents: 0,
	}
	for _, code := range schoolCodes {
		if !existingSchools[code] {
			delta[subscriptionMetricSchools]++
		}
	}
	for _, code := range studentCodes {
		if !existingStudents[code] {
			delta[subscriptionMetricStudents]++
		}
	}
	return delta, nil
}

func masterDataExistingCodeSet(ctx context.Context, exec masterDataExecutor, tenantID string, tableName string, columnName string, codes []string) (map[string]bool, error) {
	result := map[string]bool{}
	if len(codes) == 0 {
		return result, nil
	}
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE tenant_id = $1::uuid AND %s = ANY($2)`, columnName, tableName, columnName)
	rows, err := exec.QueryContext(ctx, query, tenantID, codes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		result[code] = true
	}
	return result, rows.Err()
}

func applyMasterDataImportRows(ctx context.Context, exec masterDataExecutor, rows []masterDataImportRow, tenantID string) error {
	for _, row := range rows {
		schoolYearID, err := ensureMasterDataSchoolYear(ctx, exec, row.SchoolCode, row.SchoolYearCode, tenantID)
		if err != nil {
			return err
		}
		classID, err := ensureMasterDataClass(ctx, exec, schoolYearID, row.Grade, row.ClassName)
		if err != nil {
			return err
		}
		studentID, err := ensureMasterDataStudent(ctx, exec, row.StudentCode, row.StudentName, classID, tenantID)
		if err != nil {
			return err
		}
		parentID, err := ensureMasterDataParent(ctx, exec, studentID, row.ParentName, row.ParentEmail, row.ParentPhone, row.ParentActive, tenantID)
		if err != nil {
			return err
		}
		if err := ensureMasterDataStudentParent(ctx, exec, studentID, parentID, row); err != nil {
			return err
		}
	}
	return nil
}

func saveMasterDataStudent(ctx context.Context, db *sql.DB, tenantID string, input masterDataStudentSaveInput, auditCtx requestAuditContext) (masterDataStudent, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return masterDataStudent{}, err
	}
	defer tx.Rollback()

	if _, err := findMasterDataClassByID(ctx, tx, input.ClassID, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return masterDataStudent{}, &masterDataSaveError{Status: http.StatusBadRequest, Message: "classId does not exist"}
		}
		return masterDataStudent{}, err
	}

	var existing *masterDataStudentExisting
	if input.ID != "" {
		existing, err = findMasterDataStudentByID(ctx, tx, input.ID, tenantID)
		if err != nil {
			return masterDataStudent{}, err
		}
		if existing == nil {
			return masterDataStudent{}, &masterDataSaveError{Status: http.StatusNotFound, Message: "student not found"}
		}
		byCode, err := findMasterDataStudentByCode(ctx, tx, input.StudentCode, tenantID)
		if err != nil {
			return masterDataStudent{}, err
		}
		if byCode != nil && byCode.ID != input.ID {
			return masterDataStudent{}, &masterDataSaveError{Status: http.StatusConflict, Message: "student_code already belongs to another student"}
		}
	} else {
		existing, err = findMasterDataStudentByCode(ctx, tx, input.StudentCode, tenantID)
		if err != nil {
			return masterDataStudent{}, err
		}
	}

	var studentID string
	if existing != nil {
		if existing.ClassID != input.ClassID {
			locked, err := masterDataStudentHasLockedProductionRefs(ctx, tx, existing.ID)
			if err != nil {
				return masterDataStudent{}, err
			}
			if locked {
				return masterDataStudent{}, &masterDataSaveError{Status: http.StatusConflict, Message: "student has invoices or fee adjustments; class changes are blocked"}
			}
		}
		err = tx.QueryRowContext(ctx, `
UPDATE students
SET student_code = $2,
	full_name = $3,
	class_id = $4::uuid,
	status = $5,
	updated_by_user_id = nullif($6, '')::uuid
WHERE id = $1::uuid
RETURNING id::text`,
			existing.ID, input.StudentCode, input.StudentName, input.ClassID, input.Status, auditCtx.ActorUserID).Scan(&studentID)
		if err != nil {
			return masterDataStudent{}, err
		}
	} else {
		if err := enforceTenantUsageLimit(ctx, tx, tenantID, subscriptionMetricStudents, 1, time.Now()); err != nil {
			return masterDataStudent{}, err
		}
		err = tx.QueryRowContext(ctx, `
INSERT INTO students (tenant_id, student_code, full_name, class_id, status, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2, $3, $4::uuid, $5, nullif($6, '')::uuid, nullif($6, '')::uuid)
RETURNING id::text`,
			tenantID, input.StudentCode, input.StudentName, input.ClassID, input.Status, auditCtx.ActorUserID).Scan(&studentID)
		if err != nil {
			return masterDataStudent{}, err
		}
	}

	if hasActivePrimaryParentInput(input.Parents) {
		if _, err := tx.ExecContext(ctx, `
UPDATE student_parents
SET is_primary = false,
	updated_by_user_id = nullif($2, '')::uuid
WHERE student_id = $1::uuid AND is_primary`,
			studentID, auditCtx.ActorUserID); err != nil {
			return masterDataStudent{}, err
		}
	}

	for _, parent := range input.Parents {
		parentID, err := saveMasterDataParentForStudent(ctx, tx, studentID, parent, tenantID, auditCtx)
		if err != nil {
			return masterDataStudent{}, err
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO student_parents (student_id, parent_id, relationship, is_primary, is_active, receives_billing_email, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, nullif($7, '')::uuid, nullif($7, '')::uuid)
ON CONFLICT (student_id, parent_id) DO UPDATE
SET relationship = EXCLUDED.relationship,
	is_primary = EXCLUDED.is_primary,
	is_active = EXCLUDED.is_active,
	receives_billing_email = EXCLUDED.receives_billing_email,
	updated_by_user_id = EXCLUDED.updated_by_user_id`,
			studentID, parentID, normalizeMasterDataParentRelationship(parent.Relationship), boolValue(parent.IsPrimary), boolValue(parent.IsActive), boolValue(parent.ReceivesBillingEmail), auditCtx.ActorUserID); err != nil {
			return masterDataStudent{}, err
		}
	}

	saved, err := getMasterDataStudentByID(ctx, tx, studentID, tenantID)
	if err != nil {
		return masterDataStudent{}, err
	}
	if err := rebuildTenantUsageCounter(ctx, tx, tenantID, subscriptionMetricStudents, time.Now()); err != nil {
		return masterDataStudent{}, err
	}
	return saved, tx.Commit()
}

func saveMasterDataParentForStudent(ctx context.Context, exec masterDataExecutor, studentID string, input masterDataParentSaveInput, tenantID string, auditCtx requestAuditContext) (string, error) {
	if input.ID != "" {
		existing, err := findMasterDataParentByID(ctx, exec, input.ID, tenantID)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return "", &masterDataSaveError{Status: http.StatusNotFound, Message: "parent not found"}
		}
		if input.Email != "" {
			byEmail, err := findMasterDataParentByEmail(ctx, exec, input.Email, tenantID)
			if err != nil {
				return "", err
			}
			if byEmail != nil && byEmail.ID != input.ID {
				return "", &masterDataSaveError{Status: http.StatusConflict, Message: "parent_email already belongs to another parent"}
			}
		}
		var id string
		err = exec.QueryRowContext(ctx, `
UPDATE parents
SET full_name = $2,
	email = $3,
	phone = $4,
	email_active = $5,
	updated_by_user_id = nullif($6, '')::uuid
WHERE id = $1::uuid
RETURNING id::text`,
			input.ID, input.ParentName, input.Email, input.Phone, boolValue(input.EmailActive), auditCtx.ActorUserID).Scan(&id)
		return id, err
	}

	if input.Email != "" {
		existing, err := findMasterDataParentByEmail(ctx, exec, input.Email, tenantID)
		if err != nil {
			return "", err
		}
		if existing != nil {
			if !equalText(existing.ParentName, input.ParentName) {
				return "", &masterDataSaveError{Status: http.StatusConflict, Message: "parent_email already exists with a different parent name"}
			}
			var id string
			err = exec.QueryRowContext(ctx, `
UPDATE parents
SET full_name = $2,
	phone = $3,
	email_active = $4,
	updated_by_user_id = nullif($5, '')::uuid
WHERE id = $1::uuid
RETURNING id::text`,
				existing.ID, input.ParentName, input.Phone, boolValue(input.EmailActive), auditCtx.ActorUserID).Scan(&id)
			return id, err
		}
	}

	if input.Email == "" {
		existing, err := findMasterDataLinkedParentByName(ctx, exec, studentID, input.ParentName, tenantID)
		if err != nil {
			return "", err
		}
		if existing != nil {
			var id string
			err = exec.QueryRowContext(ctx, `
UPDATE parents
SET full_name = $2,
	phone = $3,
	email_active = $4,
	updated_by_user_id = nullif($5, '')::uuid
WHERE id = $1::uuid
RETURNING id::text`,
				existing.ID, input.ParentName, input.Phone, boolValue(input.EmailActive), auditCtx.ActorUserID).Scan(&id)
			return id, err
		}
	}

	var id string
	err := exec.QueryRowContext(ctx, `
INSERT INTO parents (tenant_id, full_name, email, phone, email_active, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2, $3, $4, $5, nullif($6, '')::uuid, nullif($6, '')::uuid)
RETURNING id::text`,
		tenantID, input.ParentName, input.Email, input.Phone, boolValue(input.EmailActive), auditCtx.ActorUserID).Scan(&id)
	return id, err
}

func getMasterDataStudentByID(ctx context.Context, exec masterDataExecutor, id string, tenantID string) (masterDataStudent, error) {
	students, err := listMasterDataStudents(ctx, exec, masterDataStudentListFilters{TenantID: tenantID, StudentID: id, IncludeInactive: true})
	if err != nil {
		return masterDataStudent{}, err
	}
	if len(students) == 0 {
		return masterDataStudent{}, sql.ErrNoRows
	}
	return students[0], nil
}

func listMasterDataOptionsTx(ctx context.Context, tx *sql.Tx, tenantID string) (masterDataOptions, error) {
	options := masterDataOptions{
		Schools:     []masterDataSchoolOption{},
		SchoolYears: []masterDataSchoolYearOption{},
		Classes:     []masterDataClassOption{},
	}

	schoolRows, err := tx.QueryContext(ctx, `
SELECT id::text, code, name, status
FROM schools
WHERE tenant_id = $1::uuid
ORDER BY code`, tenantID)
	if err != nil {
		return options, err
	}
	defer schoolRows.Close()
	for schoolRows.Next() {
		var item masterDataSchoolOption
		if err := schoolRows.Scan(&item.ID, &item.Code, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.Schools = append(options.Schools, item)
	}
	if err := schoolRows.Err(); err != nil {
		return options, err
	}

	schoolYearRows, err := tx.QueryContext(ctx, `
SELECT sy.id::text, sy.school_id::text, sc.code, sy.code, sy.name, sy.status
FROM school_years sy
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
ORDER BY sc.code, sy.code DESC`, tenantID)
	if err != nil {
		return options, err
	}
	defer schoolYearRows.Close()
	for schoolYearRows.Next() {
		var item masterDataSchoolYearOption
		if err := schoolYearRows.Scan(&item.ID, &item.SchoolID, &item.SchoolCode, &item.Code, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.SchoolYears = append(options.SchoolYears, item)
	}
	if err := schoolYearRows.Err(); err != nil {
		return options, err
	}

	classRows, err := tx.QueryContext(ctx, `
SELECT c.id::text, sy.school_id::text, sc.code, c.school_year_id::text, sy.code, c.grade, c.name, c.status
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
ORDER BY sc.code, sy.code DESC, c.grade, c.name`, tenantID)
	if err != nil {
		return options, err
	}
	defer classRows.Close()
	for classRows.Next() {
		var item masterDataClassOption
		if err := classRows.Scan(&item.ID, &item.SchoolID, &item.SchoolCode, &item.SchoolYearID, &item.SchoolYearCode, &item.Grade, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.Classes = append(options.Classes, item)
	}
	if err := classRows.Err(); err != nil {
		return options, err
	}
	return options, nil
}

func findMasterDataStudentByCode(ctx context.Context, exec masterDataExecutor, code string, tenantID string) (*masterDataStudentExisting, error) {
	var student masterDataStudentExisting
	err := exec.QueryRowContext(ctx, `
SELECT s.id::text, s.student_code, s.full_name, sc.id::text, sc.code, c.id::text, c.name, c.grade, sy.id::text, sy.code
FROM students s
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE s.student_code = $1
	AND s.tenant_id = $2::uuid`, code, tenantID).Scan(
		&student.ID,
		&student.StudentCode,
		&student.StudentName,
		&student.SchoolID,
		&student.SchoolCode,
		&student.ClassID,
		&student.ClassName,
		&student.Grade,
		&student.SchoolYearID,
		&student.SchoolYearCode,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func findMasterDataStudentByID(ctx context.Context, exec masterDataExecutor, id string, tenantID string) (*masterDataStudentExisting, error) {
	var student masterDataStudentExisting
	err := exec.QueryRowContext(ctx, `
SELECT s.id::text, s.student_code, s.full_name, sc.id::text, sc.code, c.id::text, c.name, c.grade, sy.id::text, sy.code
FROM students s
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE s.id = $1::uuid
	AND s.tenant_id = $2::uuid
	AND sc.tenant_id = $2::uuid`, id, tenantID).Scan(
		&student.ID,
		&student.StudentCode,
		&student.StudentName,
		&student.SchoolID,
		&student.SchoolCode,
		&student.ClassID,
		&student.ClassName,
		&student.Grade,
		&student.SchoolYearID,
		&student.SchoolYearCode,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func findMasterDataClassByID(ctx context.Context, exec masterDataExecutor, id string, tenantID string) (*masterDataClassOption, error) {
	var class masterDataClassOption
	err := exec.QueryRowContext(ctx, `
SELECT c.id::text, sy.school_id::text, sc.code, c.school_year_id::text, sy.code, c.grade, c.name, c.status
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE c.id = $1::uuid
	AND sc.tenant_id = $2::uuid`, id, tenantID).Scan(
		&class.ID,
		&class.SchoolID,
		&class.SchoolCode,
		&class.SchoolYearID,
		&class.SchoolYearCode,
		&class.Grade,
		&class.Name,
		&class.Status,
	)
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func findMasterDataParentByEmail(ctx context.Context, exec masterDataExecutor, email string, tenantID string) (*masterDataParentExisting, error) {
	if email == "" {
		return nil, nil
	}
	var parent masterDataParentExisting
	err := exec.QueryRowContext(ctx, `
SELECT id::text, full_name, email, phone, email_active
FROM parents
WHERE email = $1
	AND tenant_id = $2::uuid`, email, tenantID).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.Phone, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func findMasterDataParentByID(ctx context.Context, exec masterDataExecutor, id string, tenantID string) (*masterDataParentExisting, error) {
	var parent masterDataParentExisting
	err := exec.QueryRowContext(ctx, `
SELECT id::text, full_name, email, phone, email_active
FROM parents
WHERE id = $1::uuid
	AND tenant_id = $2::uuid`, id, tenantID).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.Phone, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func findMasterDataLinkedParentByName(ctx context.Context, exec masterDataExecutor, studentID string, parentName string, tenantID string) (*masterDataParentExisting, error) {
	if parentName == "" {
		return nil, nil
	}
	var parent masterDataParentExisting
	err := exec.QueryRowContext(ctx, `
SELECT p.id::text, p.full_name, p.email, p.phone, p.email_active
FROM student_parents sp
JOIN parents p ON p.id = sp.parent_id
WHERE sp.student_id = $1::uuid
	AND p.email = ''
	AND p.tenant_id = $3::uuid
	AND lower(p.full_name) = lower($2)
LIMIT 1`, studentID, strings.TrimSpace(parentName), tenantID).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.Phone, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func masterDataStudentHasLockedProductionRefs(ctx context.Context, exec masterDataExecutor, studentID string) (bool, error) {
	var locked bool
	err := exec.QueryRowContext(ctx, `
SELECT EXISTS (
	SELECT 1 FROM invoices WHERE student_id = $1::uuid AND status <> 'void'
	UNION ALL
	SELECT 1 FROM student_fee_adjustments WHERE student_id = $1::uuid AND status = 'active'
)`, studentID).Scan(&locked)
	return locked, err
}

func findMasterDataStudentParentLink(ctx context.Context, exec masterDataExecutor, studentID string, parentID string) (*masterDataParentLinkExisting, error) {
	var link masterDataParentLinkExisting
	err := exec.QueryRowContext(ctx, `
SELECT is_primary, is_active, receives_billing_email
FROM student_parents
WHERE student_id = $1::uuid AND parent_id = $2::uuid`, studentID, parentID).Scan(
		&link.IsPrimary,
		&link.IsActive,
		&link.ReceivesBillingEmail,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func findMasterDataPrimaryParent(ctx context.Context, exec masterDataExecutor, studentID string) (*masterDataParentExisting, error) {
	var parent masterDataParentExisting
	err := exec.QueryRowContext(ctx, `
SELECT p.id::text, p.full_name, p.email, p.phone, p.email_active
FROM student_parents sp
JOIN parents p ON p.id = sp.parent_id
WHERE sp.student_id = $1::uuid
	AND sp.is_primary
	AND sp.is_active
LIMIT 1`, studentID).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.Phone, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func ensureMasterDataSchool(ctx context.Context, exec masterDataExecutor, code string, tenantID string) (string, error) {
	code = schoolCodeOrDefault(code)
	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO schools (tenant_id, code, name)
	VALUES ($2::uuid, $1, $1)
	ON CONFLICT (tenant_id, code) DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT sc.id::text
FROM schools sc
WHERE sc.code = $1
	AND sc.tenant_id = $2::uuid
LIMIT 1`, code, tenantID).Scan(&id)
	return id, err
}

func ensureMasterDataSchoolYear(ctx context.Context, exec masterDataExecutor, schoolCode string, code string, tenantID string) (string, error) {
	schoolID, err := ensureMasterDataSchool(ctx, exec, schoolCode, tenantID)
	if err != nil {
		return "", err
	}
	var id string
	err = exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO school_years (school_id, code, name)
	VALUES ($1::uuid, $2, $2)
	ON CONFLICT (school_id, code) DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM school_years WHERE school_id = $1::uuid AND code = $2
LIMIT 1`, schoolID, code).Scan(&id)
	return id, err
}

func ensureMasterDataClass(ctx context.Context, exec masterDataExecutor, schoolYearID string, grade string, className string) (string, error) {
	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO classes (school_year_id, grade, name)
	VALUES ($1::uuid, $2, $3)
	ON CONFLICT (school_year_id, grade, name) DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM classes WHERE school_year_id = $1::uuid AND grade = $2 AND name = $3
LIMIT 1`, schoolYearID, grade, className).Scan(&id)
	return id, err
}

func ensureMasterDataStudent(ctx context.Context, exec masterDataExecutor, code string, name string, classID string, tenantID string) (string, error) {
	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO students (tenant_id, student_code, full_name, class_id)
	VALUES ($4::uuid, $1, $2, $3::uuid)
	ON CONFLICT (tenant_id, student_code) DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM students WHERE tenant_id = $4::uuid AND student_code = $1
LIMIT 1`, code, name, classID, tenantID).Scan(&id)
	return id, err
}

func ensureMasterDataParent(ctx context.Context, exec masterDataExecutor, studentID string, parentName string, email string, phone string, active bool, tenantID string) (string, error) {
	if email == "" {
		if existing, err := findMasterDataLinkedParentByName(ctx, exec, studentID, parentName, tenantID); err != nil {
			return "", err
		} else if existing != nil {
			if phone != "" && existing.Phone != phone {
				_, err := exec.ExecContext(ctx, `UPDATE parents SET phone = $2 WHERE id = $1::uuid`, existing.ID, phone)
				if err != nil {
					return "", err
				}
			}
			return existing.ID, nil
		}
		var id string
		err := exec.QueryRowContext(ctx, `
INSERT INTO parents (tenant_id, full_name, email, phone, email_active)
VALUES ($1::uuid, $2, '', $3, $4)
RETURNING id::text`, tenantID, parentName, phone, active).Scan(&id)
		return id, err
	}

	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO parents (tenant_id, full_name, email, phone, email_active)
	VALUES ($5::uuid, $1, $2, $3, $4)
	ON CONFLICT (tenant_id, email) WHERE email <> '' DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM parents WHERE tenant_id = $5::uuid AND email = $2
LIMIT 1`, parentName, email, phone, active, tenantID).Scan(&id)
	return id, err
}

func ensureMasterDataStudentParent(ctx context.Context, exec masterDataExecutor, studentID string, parentID string, row masterDataImportRow) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO student_parents (student_id, parent_id, relationship, is_primary, is_active, receives_billing_email)
VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6)
ON CONFLICT (student_id, parent_id) DO NOTHING`, studentID, parentID, normalizeMasterDataParentRelationship(row.Relationship), row.IsPrimary, row.ParentActive, row.ReceivesBillingEmail)
	return err
}

func normalizeMasterDataStudentSaveInput(input masterDataStudentSaveInput) masterDataStudentSaveInput {
	input.ID = strings.TrimSpace(input.ID)
	input.StudentCode = normalizeStudentCode(input.StudentCode)
	input.StudentName = strings.TrimSpace(input.StudentName)
	input.ClassID = strings.TrimSpace(input.ClassID)
	input.Status = normalizeMasterDataStudentStatus(input.Status)

	hasPrimaryValue := false
	for idx := range input.Parents {
		parent := &input.Parents[idx]
		parent.ID = strings.TrimSpace(parent.ID)
		parent.ParentName = strings.TrimSpace(parent.ParentName)
		parent.Email = normalizeEmail(parent.Email)
		parent.Phone = normalizeAdminPhone(parent.Phone)
		parent.Relationship = normalizeMasterDataParentRelationship(parent.Relationship)
		parent.EmailActive = boolDefault(parent.EmailActive, true)
		parent.IsActive = boolDefault(parent.IsActive, true)
		parent.ReceivesBillingEmail = boolDefault(parent.ReceivesBillingEmail, parent.Email != "")
		if parent.IsPrimary != nil {
			hasPrimaryValue = true
		}
	}
	if len(input.Parents) > 0 && !hasPrimaryValue {
		input.Parents[0].IsPrimary = boolDefault(input.Parents[0].IsPrimary, true)
	}
	for idx := range input.Parents {
		input.Parents[idx].IsPrimary = boolDefault(input.Parents[idx].IsPrimary, false)
	}
	return input
}

func normalizeMasterDataParentRelationship(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "guardian"
	}
	return headerKey(value)
}

func validateMasterDataStudentSaveInput(input masterDataStudentSaveInput) error {
	if input.StudentCode == "" {
		return errors.New("studentCode is required")
	}
	if input.StudentName == "" {
		return errors.New("studentName is required")
	}
	if input.ClassID == "" {
		return errors.New("classId is required")
	}
	if err := validateMasterDataStudentStatus(input.Status); err != nil {
		return err
	}
	if len(input.Parents) == 0 {
		return errors.New("at least one parent contact is required")
	}

	activePrimaryCount := 0
	parentKeys := map[string]bool{}
	for idx, parent := range input.Parents {
		if parent.ParentName == "" {
			return fmt.Errorf("parents[%d].parentName is required", idx)
		}
		if boolValue(parent.ReceivesBillingEmail) && parent.Email == "" {
			return fmt.Errorf("parents[%d].email is required when receivesBillingEmail is true", idx)
		}
		if boolValue(parent.IsPrimary) && boolValue(parent.IsActive) {
			activePrimaryCount++
		}
		key := "name:" + strings.ToLower(parent.ParentName)
		if parent.Email != "" {
			key = "email:" + parent.Email
		} else if parent.Phone != "" {
			key = "phone:" + parent.Phone
		}
		if parentKeys[key] {
			return fmt.Errorf("duplicate parent contact: %s", parent.ParentName)
		}
		parentKeys[key] = true
	}
	if activePrimaryCount > 1 {
		return errors.New("only one active primary parent is allowed")
	}
	return nil
}

func hasActivePrimaryParentInput(parents []masterDataParentSaveInput) bool {
	for _, parent := range parents {
		if boolValue(parent.IsPrimary) && boolValue(parent.IsActive) {
			return true
		}
	}
	return false
}

func normalizeMasterDataStudentStatus(status string) string {
	status = headerKey(status)
	if status == "" {
		return "active"
	}
	return status
}

func validateMasterDataStudentStatus(status string) error {
	switch status {
	case "active", "inactive", "graduated":
		return nil
	default:
		return errors.New("status must be active, inactive, or graduated")
	}
}

func boolDefault(value *bool, fallback bool) *bool {
	if value != nil {
		return value
	}
	return &fallback
}

func boolValue(value *bool) bool {
	return value != nil && *value
}

func sameImportStudent(a masterDataImportRow, b masterDataImportRow) bool {
	return equalText(schoolCodeOrDefault(a.SchoolCode), schoolCodeOrDefault(b.SchoolCode)) &&
		equalText(a.StudentName, b.StudentName) &&
		equalText(a.SchoolYearCode, b.SchoolYearCode) &&
		equalText(a.Grade, b.Grade) &&
		equalText(a.ClassName, b.ClassName)
}

func importStudentDescription(row masterDataImportRow) string {
	return fmt.Sprintf("%s / %s / %s / %s / %s", row.StudentName, schoolCodeOrDefault(row.SchoolCode), row.SchoolYearCode, row.Grade, row.ClassName)
}

func importParentDescription(row masterDataImportRow) string {
	if row.ParentEmail != "" {
		return row.ParentName + " <" + row.ParentEmail + ">"
	}
	if row.ParentPhone != "" {
		return row.ParentName + " <" + row.ParentPhone + ">"
	}
	return row.ParentName
}

func importExistingParentDescription(parent masterDataParentExisting) string {
	if parent.Email != "" {
		return parent.ParentName + " <" + parent.Email + ">"
	}
	if parent.Phone != "" {
		return parent.ParentName + " <" + parent.Phone + ">"
	}
	return parent.ParentName
}

func importParentKey(row masterDataImportRow) string {
	if row.ParentEmail != "" {
		return "email:" + row.ParentEmail
	}
	if row.ParentPhone != "" {
		return "phone:" + row.ParentPhone
	}
	return "name:" + strings.ToLower(strings.TrimSpace(row.ParentName))
}

func formatParentFlags(primary bool, active bool, receivesBilling bool) string {
	return fmt.Sprintf("primary=%t active=%t receives_billing_email=%t", primary, active, receivesBilling)
}

func normalizeStudentCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeSchoolCode(value string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(value)))
	return strings.Join(parts, "_")
}

func schoolCodeOrDefault(value string) string {
	if code := normalizeSchoolCode(value); code != "" {
		return code
	}
	return defaultTenantCode
}

func normalizeSchoolYearCode(value string) string {
	return strings.TrimSpace(value)
}

func normalizeGrade(value string) string {
	return strings.TrimSpace(value)
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func deriveGradeFromClass(className string) string {
	match := classGradePattern.FindString(strings.TrimSpace(className))
	return match
}

func parseBoolWithDefault(value string, fallback bool) bool {
	normalized := headerKey(value)
	switch normalized {
	case "1", "true", "t", "yes", "y", "co", "c", "x", "bat", "active":
		return true
	case "0", "false", "f", "no", "n", "khong", "k", "tat", "inactive":
		return false
	default:
		return fallback
	}
}

func equalText(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func sortMasterDataImportIssues(issues []masterDataImportIssue) {
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].RowNumber == issues[j].RowNumber {
			return issues[i].Type < issues[j].Type
		}
		return issues[i].RowNumber < issues[j].RowNumber
	})
}
