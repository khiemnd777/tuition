package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

const maxMasterDataImportRows = 1000

var classGradePattern = regexp.MustCompile(`\d+`)

type masterDataSchoolYearOption struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type masterDataClassOption struct {
	ID             string `json:"id"`
	SchoolYearID   string `json:"schoolYearId"`
	SchoolYearCode string `json:"schoolYearCode"`
	Grade          string `json:"grade"`
	Name           string `json:"name"`
	Status         string `json:"status"`
}

type masterDataOptions struct {
	SchoolYears []masterDataSchoolYearOption `json:"schoolYears"`
	Classes     []masterDataClassOption      `json:"classes"`
}

type masterDataStudentListFilters struct {
	SchoolYearID string
	SchoolYear   string
	ClassID      string
	Grade        string
	Query        string
}

type masterDataStudent struct {
	ID             string                    `json:"id"`
	StudentCode    string                    `json:"studentCode"`
	StudentName    string                    `json:"studentName"`
	Status         string                    `json:"status"`
	ClassID        string                    `json:"classId"`
	ClassName      string                    `json:"className"`
	Grade          string                    `json:"grade"`
	SchoolYearID   string                    `json:"schoolYearId"`
	SchoolYearCode string                    `json:"schoolYearCode"`
	Parents        []masterDataParentContact `json:"parents"`
}

type masterDataParentContact struct {
	ID                   string `json:"id"`
	ParentName           string `json:"parentName"`
	Email                string `json:"email"`
	EmailActive          bool   `json:"emailActive"`
	IsPrimary            bool   `json:"isPrimary"`
	IsActive             bool   `json:"isActive"`
	ReceivesBillingEmail bool   `json:"receivesBillingEmail"`
}

type masterDataImportRow struct {
	RowNumber            int    `json:"rowNumber"`
	StudentCode          string `json:"studentCode"`
	StudentName          string `json:"studentName"`
	SchoolYearCode       string `json:"schoolYearCode"`
	Grade                string `json:"grade"`
	ClassName            string `json:"className"`
	ParentName           string `json:"parentName"`
	ParentEmail          string `json:"parentEmail"`
	IsPrimary            bool   `json:"isPrimary"`
	ParentActive         bool   `json:"parentActive"`
	ReceivesBillingEmail bool   `json:"receivesBillingEmail"`
}

type masterDataImportRowResult struct {
	RowNumber      int    `json:"rowNumber"`
	StudentCode    string `json:"studentCode"`
	StudentName    string `json:"studentName"`
	SchoolYearCode string `json:"schoolYearCode"`
	Grade          string `json:"grade"`
	ClassName      string `json:"className"`
	ParentName     string `json:"parentName"`
	ParentEmail    string `json:"parentEmail"`
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

func handleMasterDataOptions(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	options, err := listMasterDataOptions(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load master data options", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, options)
}

func handleMasterDataStudents(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	query := r.URL.Query()
	students, err := listMasterDataStudents(r.Context(), db, masterDataStudentListFilters{
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

func handleMasterDataImportCSV(w http.ResponseWriter, r *http.Request) {
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
	response, err := importMasterDataRows(r.Context(), db, rows, apply)
	if err != nil {
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
		return nil, errors.New("database connection failed; check ABC database environment and server")
	}
	return db, nil
}

func writeMasterDataDBError(w http.ResponseWriter, err error) {
	http.Error(w, "master data database unavailable: "+err.Error(), http.StatusServiceUnavailable)
}

func listMasterDataOptions(ctx context.Context, db *sql.DB) (masterDataOptions, error) {
	options := masterDataOptions{
		SchoolYears: []masterDataSchoolYearOption{},
		Classes:     []masterDataClassOption{},
	}

	schoolYearRows, err := db.QueryContext(ctx, `
SELECT id::text, code, name, status
FROM school_years
ORDER BY code DESC`)
	if err != nil {
		return options, err
	}
	defer schoolYearRows.Close()
	for schoolYearRows.Next() {
		var item masterDataSchoolYearOption
		if err := schoolYearRows.Scan(&item.ID, &item.Code, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.SchoolYears = append(options.SchoolYears, item)
	}
	if err := schoolYearRows.Err(); err != nil {
		return options, err
	}

	classRows, err := db.QueryContext(ctx, `
SELECT c.id::text, c.school_year_id::text, sy.code, c.grade, c.name, c.status
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
ORDER BY sy.code DESC, c.grade, c.name`)
	if err != nil {
		return options, err
	}
	defer classRows.Close()
	for classRows.Next() {
		var item masterDataClassOption
		if err := classRows.Scan(&item.ID, &item.SchoolYearID, &item.SchoolYearCode, &item.Grade, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.Classes = append(options.Classes, item)
	}
	if err := classRows.Err(); err != nil {
		return options, err
	}

	return options, nil
}

func listMasterDataStudents(ctx context.Context, db *sql.DB, filters masterDataStudentListFilters) ([]masterDataStudent, error) {
	args := []any{}
	conditions := []string{"s.status <> 'inactive'"}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	if filters.SchoolYearID != "" {
		conditions = append(conditions, "sy.id = "+addArg(filters.SchoolYearID)+"::uuid")
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
		conditions = append(conditions, "(s.student_code ILIKE "+placeholder+" OR s.full_name ILIKE "+placeholder+" OR p.full_name ILIKE "+placeholder+" OR p.email ILIKE "+placeholder+")")
	}

	query := `
SELECT
	s.id::text,
	s.student_code,
	s.full_name,
	s.status,
	c.id::text,
	c.name,
	c.grade,
	sy.id::text,
	sy.code,
	COALESCE(p.id::text, ''),
	COALESCE(p.full_name, ''),
	COALESCE(p.email, ''),
	COALESCE(p.email_active, false),
	COALESCE(sp.is_primary, false),
	COALESCE(sp.is_active, false),
	COALESCE(sp.receives_billing_email, false)
FROM students s
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
LEFT JOIN student_parents sp ON sp.student_id = s.id
LEFT JOIN parents p ON p.id = sp.parent_id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY sy.code DESC, c.grade, c.name, s.student_code, sp.is_primary DESC, p.full_name
LIMIT 1000`

	rows, err := db.QueryContext(ctx, query, args...)
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
			&student.ClassID,
			&student.ClassName,
			&student.Grade,
			&student.SchoolYearID,
			&student.SchoolYearCode,
			&parent.ID,
			&parent.ParentName,
			&parent.Email,
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
	return students, nil
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
			SchoolYearCode:       normalizeSchoolYearCode(importFieldValue(record, table, mapping, "school_year", masterDataCSVAliases("school_year"))),
			Grade:                normalizeGrade(importFieldValue(record, table, mapping, "grade", masterDataCSVAliases("grade"))),
			ClassName:            strings.TrimSpace(importFieldValue(record, table, mapping, "class_name", masterDataCSVAliases("class_name"))),
			ParentName:           strings.TrimSpace(importFieldValue(record, table, mapping, "parent", masterDataCSVAliases("parent"))),
			ParentEmail:          parentEmail,
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
		if row.ParentName == "" && row.ParentEmail == "" {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "missing_parent",
				Message:     "parent_name or parent_email is required",
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

		if row.StudentCode != "" && (row.ParentEmail != "" || row.ParentName != "") {
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

		if row.StudentCode != "" && row.IsPrimary && (row.ParentEmail != "" || row.ParentName != "") {
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

func importMasterDataRows(ctx context.Context, db *sql.DB, rows []masterDataImportRow, apply bool) (masterDataImportResponse, error) {
	if apply {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return masterDataImportResponse{}, err
		}
		defer tx.Rollback()

		response, err := buildMasterDataImportResponse(ctx, tx, rows)
		if err != nil || len(response.Issues) > 0 {
			response.Applied = false
			return response, err
		}
		if err := applyMasterDataImportRows(ctx, tx, rows); err != nil {
			return masterDataImportResponse{}, err
		}
		options, err := listMasterDataOptionsTx(ctx, tx)
		if err != nil {
			return masterDataImportResponse{}, err
		}
		response.Applied = true
		response.Options = &options
		return response, tx.Commit()
	}

	response, err := buildMasterDataImportResponse(ctx, db, rows)
	response.Applied = false
	return response, err
}

func buildMasterDataImportResponse(ctx context.Context, exec masterDataExecutor, rows []masterDataImportRow) (masterDataImportResponse, error) {
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
			SchoolYearCode: row.SchoolYearCode,
			Grade:          row.Grade,
			ClassName:      row.ClassName,
			ParentName:     row.ParentName,
			ParentEmail:    row.ParentEmail,
			Action:         "ready",
		}
		if rowHasIssue[row.RowNumber] {
			result.Action = "conflict"
			response.Rows = append(response.Rows, result)
			continue
		}

		action, issues, err := classifyMasterDataImportRow(ctx, exec, row)
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

func classifyMasterDataImportRow(ctx context.Context, exec masterDataExecutor, row masterDataImportRow) (string, []masterDataImportIssue, error) {
	issues := []masterDataImportIssue{}
	student, err := findMasterDataStudentByCode(ctx, exec, row.StudentCode)
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
		if !equalText(student.SchoolYearCode, row.SchoolYearCode) || !equalText(student.Grade, row.Grade) || !equalText(student.ClassName, row.ClassName) {
			issues = append(issues, masterDataImportIssue{
				RowNumber:   row.RowNumber,
				StudentCode: row.StudentCode,
				Type:        "student_class_conflict",
				Message:     "student_code already exists in a different class, grade, or school year",
				Existing:    fmt.Sprintf("%s / %s / %s", student.SchoolYearCode, student.Grade, student.ClassName),
				Incoming:    fmt.Sprintf("%s / %s / %s", row.SchoolYearCode, row.Grade, row.ClassName),
			})
		}
	}

	parent, err := findMasterDataParentByEmail(ctx, exec, row.ParentEmail)
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
		parent, err = findMasterDataLinkedParentByName(ctx, exec, student.ID, row.ParentName)
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

func applyMasterDataImportRows(ctx context.Context, exec masterDataExecutor, rows []masterDataImportRow) error {
	for _, row := range rows {
		schoolYearID, err := ensureMasterDataSchoolYear(ctx, exec, row.SchoolYearCode)
		if err != nil {
			return err
		}
		classID, err := ensureMasterDataClass(ctx, exec, schoolYearID, row.Grade, row.ClassName)
		if err != nil {
			return err
		}
		studentID, err := ensureMasterDataStudent(ctx, exec, row.StudentCode, row.StudentName, classID)
		if err != nil {
			return err
		}
		parentID, err := ensureMasterDataParent(ctx, exec, studentID, row.ParentName, row.ParentEmail, row.ParentActive)
		if err != nil {
			return err
		}
		if err := ensureMasterDataStudentParent(ctx, exec, studentID, parentID, row); err != nil {
			return err
		}
	}
	return nil
}

func listMasterDataOptionsTx(ctx context.Context, tx *sql.Tx) (masterDataOptions, error) {
	options := masterDataOptions{
		SchoolYears: []masterDataSchoolYearOption{},
		Classes:     []masterDataClassOption{},
	}
	schoolYearRows, err := tx.QueryContext(ctx, `
SELECT id::text, code, name, status
FROM school_years
ORDER BY code DESC`)
	if err != nil {
		return options, err
	}
	defer schoolYearRows.Close()
	for schoolYearRows.Next() {
		var item masterDataSchoolYearOption
		if err := schoolYearRows.Scan(&item.ID, &item.Code, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.SchoolYears = append(options.SchoolYears, item)
	}
	if err := schoolYearRows.Err(); err != nil {
		return options, err
	}

	classRows, err := tx.QueryContext(ctx, `
SELECT c.id::text, c.school_year_id::text, sy.code, c.grade, c.name, c.status
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
ORDER BY sy.code DESC, c.grade, c.name`)
	if err != nil {
		return options, err
	}
	defer classRows.Close()
	for classRows.Next() {
		var item masterDataClassOption
		if err := classRows.Scan(&item.ID, &item.SchoolYearID, &item.SchoolYearCode, &item.Grade, &item.Name, &item.Status); err != nil {
			return options, err
		}
		options.Classes = append(options.Classes, item)
	}
	if err := classRows.Err(); err != nil {
		return options, err
	}
	return options, nil
}

func findMasterDataStudentByCode(ctx context.Context, exec masterDataExecutor, code string) (*masterDataStudentExisting, error) {
	var student masterDataStudentExisting
	err := exec.QueryRowContext(ctx, `
SELECT s.id::text, s.student_code, s.full_name, c.id::text, c.name, c.grade, sy.id::text, sy.code
FROM students s
JOIN classes c ON c.id = s.class_id
JOIN school_years sy ON sy.id = c.school_year_id
WHERE s.student_code = $1`, code).Scan(
		&student.ID,
		&student.StudentCode,
		&student.StudentName,
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

func findMasterDataParentByEmail(ctx context.Context, exec masterDataExecutor, email string) (*masterDataParentExisting, error) {
	if email == "" {
		return nil, nil
	}
	var parent masterDataParentExisting
	err := exec.QueryRowContext(ctx, `
SELECT id::text, full_name, email, email_active
FROM parents
WHERE email = $1`, email).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func findMasterDataLinkedParentByName(ctx context.Context, exec masterDataExecutor, studentID string, parentName string) (*masterDataParentExisting, error) {
	if parentName == "" {
		return nil, nil
	}
	var parent masterDataParentExisting
	err := exec.QueryRowContext(ctx, `
SELECT p.id::text, p.full_name, p.email, p.email_active
FROM student_parents sp
JOIN parents p ON p.id = sp.parent_id
WHERE sp.student_id = $1::uuid
	AND p.email = ''
	AND lower(p.full_name) = lower($2)
LIMIT 1`, studentID, strings.TrimSpace(parentName)).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
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
SELECT p.id::text, p.full_name, p.email, p.email_active
FROM student_parents sp
JOIN parents p ON p.id = sp.parent_id
WHERE sp.student_id = $1::uuid
	AND sp.is_primary
	AND sp.is_active
LIMIT 1`, studentID).Scan(&parent.ID, &parent.ParentName, &parent.Email, &parent.EmailActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func ensureMasterDataSchoolYear(ctx context.Context, exec masterDataExecutor, code string) (string, error) {
	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO school_years (code, name)
	VALUES ($1, $1)
	ON CONFLICT (code) DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM school_years WHERE code = $1
LIMIT 1`, code).Scan(&id)
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

func ensureMasterDataStudent(ctx context.Context, exec masterDataExecutor, code string, name string, classID string) (string, error) {
	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO students (student_code, full_name, class_id)
	VALUES ($1, $2, $3::uuid)
	ON CONFLICT (student_code) DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM students WHERE student_code = $1
LIMIT 1`, code, name, classID).Scan(&id)
	return id, err
}

func ensureMasterDataParent(ctx context.Context, exec masterDataExecutor, studentID string, parentName string, email string, active bool) (string, error) {
	if email == "" {
		if existing, err := findMasterDataLinkedParentByName(ctx, exec, studentID, parentName); err != nil {
			return "", err
		} else if existing != nil {
			return existing.ID, nil
		}
		var id string
		err := exec.QueryRowContext(ctx, `
INSERT INTO parents (full_name, email, email_active)
VALUES ($1, '', $2)
RETURNING id::text`, parentName, active).Scan(&id)
		return id, err
	}

	var id string
	err := exec.QueryRowContext(ctx, `
WITH inserted AS (
	INSERT INTO parents (full_name, email, email_active)
	VALUES ($1, $2, $3)
	ON CONFLICT (email) WHERE email <> '' DO NOTHING
	RETURNING id::text
)
SELECT id FROM inserted
UNION ALL
SELECT id::text FROM parents WHERE email = $2
LIMIT 1`, parentName, email, active).Scan(&id)
	return id, err
}

func ensureMasterDataStudentParent(ctx context.Context, exec masterDataExecutor, studentID string, parentID string, row masterDataImportRow) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO student_parents (student_id, parent_id, relationship, is_primary, is_active, receives_billing_email)
VALUES ($1::uuid, $2::uuid, 'guardian', $3, $4, $5)
ON CONFLICT (student_id, parent_id) DO NOTHING`, studentID, parentID, row.IsPrimary, row.ParentActive, row.ReceivesBillingEmail)
	return err
}

func sameImportStudent(a masterDataImportRow, b masterDataImportRow) bool {
	return equalText(a.StudentName, b.StudentName) &&
		equalText(a.SchoolYearCode, b.SchoolYearCode) &&
		equalText(a.Grade, b.Grade) &&
		equalText(a.ClassName, b.ClassName)
}

func importStudentDescription(row masterDataImportRow) string {
	return fmt.Sprintf("%s / %s / %s / %s", row.StudentName, row.SchoolYearCode, row.Grade, row.ClassName)
}

func importParentDescription(row masterDataImportRow) string {
	if row.ParentEmail != "" {
		return row.ParentName + " <" + row.ParentEmail + ">"
	}
	return row.ParentName
}

func importExistingParentDescription(parent masterDataParentExisting) string {
	if parent.Email != "" {
		return parent.ParentName + " <" + parent.Email + ">"
	}
	return parent.ParentName
}

func importParentKey(row masterDataImportRow) string {
	if row.ParentEmail != "" {
		return "email:" + row.ParentEmail
	}
	return "name:" + strings.ToLower(strings.TrimSpace(row.ParentName))
}

func formatParentFlags(primary bool, active bool, receivesBilling bool) string {
	return fmt.Sprintf("primary=%t active=%t receives_billing_email=%t", primary, active, receivesBilling)
}

func normalizeStudentCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
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
