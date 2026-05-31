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

type schoolTreeResponse struct {
	Schools []schoolTreeSchool `json:"schools"`
}

type schoolTreeSchool struct {
	ID               string                 `json:"id"`
	Code             string                 `json:"code"`
	Name             string                 `json:"name"`
	Status           string                 `json:"status"`
	StudentCount     int                    `json:"studentCount"`
	ClassCount       int                    `json:"classCount"`
	FeeScheduleCount int                    `json:"feeScheduleCount"`
	AdjustmentCount  int                    `json:"adjustmentCount"`
	SchoolYears      []schoolTreeSchoolYear `json:"schoolYears"`
}

type schoolTreeSchoolYear struct {
	ID               string            `json:"id"`
	SchoolID         string            `json:"schoolId"`
	SchoolCode       string            `json:"schoolCode"`
	Code             string            `json:"code"`
	Name             string            `json:"name"`
	Status           string            `json:"status"`
	StudentCount     int               `json:"studentCount"`
	ClassCount       int               `json:"classCount"`
	FeeScheduleCount int               `json:"feeScheduleCount"`
	AdjustmentCount  int               `json:"adjustmentCount"`
	Grades           []schoolTreeGrade `json:"grades"`
}

type schoolTreeGrade struct {
	SchoolYearID     string            `json:"schoolYearId"`
	SchoolYearCode   string            `json:"schoolYearCode"`
	Grade            string            `json:"grade"`
	StudentCount     int               `json:"studentCount"`
	ClassCount       int               `json:"classCount"`
	FeeScheduleCount int               `json:"feeScheduleCount"`
	AdjustmentCount  int               `json:"adjustmentCount"`
	Classes          []schoolTreeClass `json:"classes"`
}

type schoolTreeClass struct {
	ID                     string `json:"id"`
	SchoolID               string `json:"schoolId"`
	SchoolCode             string `json:"schoolCode"`
	SchoolYearID           string `json:"schoolYearId"`
	SchoolYearCode         string `json:"schoolYearCode"`
	Grade                  string `json:"grade"`
	Name                   string `json:"name"`
	Status                 string `json:"status"`
	StudentCount           int    `json:"studentCount"`
	FeeScheduleCount       int    `json:"feeScheduleCount"`
	ActiveFeeScheduleCount int    `json:"activeFeeScheduleCount"`
	AdjustmentCount        int    `json:"adjustmentCount"`
	LatestFeeScheduleID    string `json:"latestFeeScheduleId,omitempty"`
	LatestFeeScheduleName  string `json:"latestFeeScheduleName,omitempty"`
	LatestPeriodCode       string `json:"latestPeriodCode,omitempty"`
	LatestScheduleStatus   string `json:"latestScheduleStatus,omitempty"`
}

type schoolTreeSchoolInput struct {
	ID     string `json:"id,omitempty"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type schoolTreeSchoolYearInput struct {
	ID       string `json:"id,omitempty"`
	SchoolID string `json:"schoolId"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	StartsOn string `json:"startsOn,omitempty"`
	EndsOn   string `json:"endsOn,omitempty"`
}

type schoolTreeClassInput struct {
	ID           string `json:"id,omitempty"`
	SchoolYearID string `json:"schoolYearId"`
	Grade        string `json:"grade"`
	Name         string `json:"name"`
	Status       string `json:"status"`
}

func handleSchoolTree(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	tree, err := loadSchoolTree(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot load school tree", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, schoolTreeResponse{Schools: tree})
}

func handleSchoolTreeSchoolSave(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input schoolTreeSchoolInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSchoolTreeSchoolInput(input)
	if err := validateSchoolTreeSchoolInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	saved, err := saveSchoolTreeSchool(r.Context(), db, input, auditContextFromRequest(r))
	if err != nil {
		http.Error(w, "cannot save school", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"school": saved})
}

func handleSchoolTreeSchoolYearSave(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input schoolTreeSchoolYearInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSchoolTreeSchoolYearInput(input)
	if err := validateSchoolTreeSchoolYearInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	saved, err := saveSchoolTreeSchoolYear(r.Context(), db, input, auditContextFromRequest(r))
	if err != nil {
		http.Error(w, "cannot save school year", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"schoolYear": saved})
}

func handleSchoolTreeClassSave(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input schoolTreeClassInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeSchoolTreeClassInput(input)
	if err := validateSchoolTreeClassInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	saved, err := saveSchoolTreeClass(r.Context(), db, input, auditContextFromRequest(r))
	if err != nil {
		http.Error(w, "cannot save class", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"class": saved})
}

func loadSchoolTree(ctx context.Context, db *sql.DB) ([]schoolTreeSchool, error) {
	schools, err := listSchoolTreeSchools(ctx, db)
	if err != nil {
		return nil, err
	}
	years, err := listSchoolTreeSchoolYears(ctx, db)
	if err != nil {
		return nil, err
	}
	classes, err := listSchoolTreeClasses(ctx, db)
	if err != nil {
		return nil, err
	}
	return buildSchoolTree(schools, years, classes), nil
}

func listSchoolTreeSchools(ctx context.Context, db *sql.DB) ([]schoolTreeSchool, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id::text, code, name, status
FROM schools
ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []schoolTreeSchool{}
	for rows.Next() {
		var item schoolTreeSchool
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Status); err != nil {
			return nil, err
		}
		item.SchoolYears = []schoolTreeSchoolYear{}
		items = append(items, item)
	}
	return items, rows.Err()
}

func listSchoolTreeSchoolYears(ctx context.Context, db *sql.DB) ([]schoolTreeSchoolYear, error) {
	rows, err := db.QueryContext(ctx, `
SELECT
	sy.id::text,
	sy.school_id::text,
	sc.code,
	sy.code,
	sy.name,
	sy.status,
	COALESCE(class_counts.class_count, 0),
	COALESCE(student_counts.student_count, 0),
	COALESCE(schedule_counts.schedule_count, 0),
	COALESCE(adjustment_counts.adjustment_count, 0)
FROM school_years sy
JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN (
	SELECT school_year_id, COUNT(*)::integer AS class_count
	FROM classes
	WHERE status = 'active'
	GROUP BY school_year_id
) class_counts ON class_counts.school_year_id = sy.id
LEFT JOIN (
	SELECT c.school_year_id, COUNT(s.id)::integer AS student_count
	FROM classes c
	JOIN students s ON s.class_id = c.id
	WHERE s.status <> 'inactive'
	GROUP BY c.school_year_id
) student_counts ON student_counts.school_year_id = sy.id
LEFT JOIN (
	SELECT school_year_id, COUNT(*)::integer AS schedule_count
	FROM fee_schedules
	WHERE status <> 'archived'
	GROUP BY school_year_id
) schedule_counts ON schedule_counts.school_year_id = sy.id
LEFT JOIN (
	SELECT fs.school_year_id, COUNT(sfa.id)::integer AS adjustment_count
	FROM fee_schedules fs
	JOIN student_fee_adjustments sfa ON sfa.schedule_id = fs.id
	WHERE fs.status <> 'archived' AND sfa.status = 'active'
	GROUP BY fs.school_year_id
) adjustment_counts ON adjustment_counts.school_year_id = sy.id
ORDER BY sc.code, sy.code DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []schoolTreeSchoolYear{}
	for rows.Next() {
		var item schoolTreeSchoolYear
		if err := rows.Scan(
			&item.ID,
			&item.SchoolID,
			&item.SchoolCode,
			&item.Code,
			&item.Name,
			&item.Status,
			&item.ClassCount,
			&item.StudentCount,
			&item.FeeScheduleCount,
			&item.AdjustmentCount,
		); err != nil {
			return nil, err
		}
		item.Grades = []schoolTreeGrade{}
		items = append(items, item)
	}
	return items, rows.Err()
}

func listSchoolTreeClasses(ctx context.Context, db *sql.DB) ([]schoolTreeClass, error) {
	rows, err := db.QueryContext(ctx, `
SELECT
	c.id::text,
	sy.school_id::text,
	sc.code,
	c.school_year_id::text,
	sy.code,
	c.grade,
	c.name,
	c.status,
	COALESCE(student_counts.student_count, 0),
	COALESCE(schedule_counts.schedule_count, 0),
	COALESCE(schedule_counts.active_schedule_count, 0),
	COALESCE(adjustment_counts.adjustment_count, 0),
	COALESCE(latest_schedule.id::text, ''),
	COALESCE(latest_schedule.name, ''),
	COALESCE(latest_schedule.period_code, ''),
	COALESCE(latest_schedule.status, '')
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
LEFT JOIN (
	SELECT class_id, COUNT(*)::integer AS student_count
	FROM students
	WHERE status <> 'inactive'
	GROUP BY class_id
) student_counts ON student_counts.class_id = c.id
LEFT JOIN (
	SELECT
		class_id,
		COUNT(*)::integer AS schedule_count,
		COUNT(*) FILTER (WHERE status = 'active')::integer AS active_schedule_count
	FROM fee_schedules
	WHERE class_id IS NOT NULL AND status <> 'archived'
	GROUP BY class_id
) schedule_counts ON schedule_counts.class_id = c.id
LEFT JOIN (
	SELECT fs.class_id, COUNT(sfa.id)::integer AS adjustment_count
	FROM fee_schedules fs
	JOIN student_fee_adjustments sfa ON sfa.schedule_id = fs.id
	WHERE fs.class_id IS NOT NULL AND fs.status <> 'archived' AND sfa.status = 'active'
	GROUP BY fs.class_id
) adjustment_counts ON adjustment_counts.class_id = c.id
LEFT JOIN LATERAL (
	SELECT id, name, period_code, status
	FROM fee_schedules fs
	WHERE fs.class_id = c.id AND fs.status <> 'archived'
	ORDER BY (fs.status = 'active') DESC, fs.created_at DESC
	LIMIT 1
) latest_schedule ON true
ORDER BY sc.code, sy.code DESC, c.grade, c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []schoolTreeClass{}
	for rows.Next() {
		var item schoolTreeClass
		if err := rows.Scan(
			&item.ID,
			&item.SchoolID,
			&item.SchoolCode,
			&item.SchoolYearID,
			&item.SchoolYearCode,
			&item.Grade,
			&item.Name,
			&item.Status,
			&item.StudentCount,
			&item.FeeScheduleCount,
			&item.ActiveFeeScheduleCount,
			&item.AdjustmentCount,
			&item.LatestFeeScheduleID,
			&item.LatestFeeScheduleName,
			&item.LatestPeriodCode,
			&item.LatestScheduleStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func buildSchoolTree(schools []schoolTreeSchool, years []schoolTreeSchoolYear, classes []schoolTreeClass) []schoolTreeSchool {
	schoolByID := map[string]*schoolTreeSchool{}
	for idx := range schools {
		schools[idx].SchoolYears = []schoolTreeSchoolYear{}
		schoolByID[schools[idx].ID] = &schools[idx]
	}

	for _, year := range years {
		year.Grades = []schoolTreeGrade{}
		school := schoolByID[year.SchoolID]
		if school == nil {
			continue
		}
		school.StudentCount += year.StudentCount
		school.ClassCount += year.ClassCount
		school.FeeScheduleCount += year.FeeScheduleCount
		school.AdjustmentCount += year.AdjustmentCount
		school.SchoolYears = append(school.SchoolYears, year)
	}

	yearByID := map[string]*schoolTreeSchoolYear{}
	for schoolIdx := range schools {
		for yearIdx := range schools[schoolIdx].SchoolYears {
			year := &schools[schoolIdx].SchoolYears[yearIdx]
			yearByID[year.ID] = year
		}
	}

	gradeByKey := map[string]*schoolTreeGrade{}
	for _, class := range classes {
		year := yearByID[class.SchoolYearID]
		if year == nil {
			continue
		}
		key := class.SchoolYearID + "|" + class.Grade
		grade := gradeByKey[key]
		if grade == nil {
			year.Grades = append(year.Grades, schoolTreeGrade{
				SchoolYearID:   class.SchoolYearID,
				SchoolYearCode: class.SchoolYearCode,
				Grade:          class.Grade,
				Classes:        []schoolTreeClass{},
			})
			grade = &year.Grades[len(year.Grades)-1]
			gradeByKey[key] = grade
		}
		grade.ClassCount++
		grade.StudentCount += class.StudentCount
		grade.FeeScheduleCount += class.FeeScheduleCount
		grade.AdjustmentCount += class.AdjustmentCount
		grade.Classes = append(grade.Classes, class)
	}

	for idx := range schools {
		sort.SliceStable(schools[idx].SchoolYears, func(i, j int) bool {
			return schools[idx].SchoolYears[i].Code > schools[idx].SchoolYears[j].Code
		})
		for yearIdx := range schools[idx].SchoolYears {
			year := &schools[idx].SchoolYears[yearIdx]
			sort.SliceStable(year.Grades, func(i, j int) bool {
				return naturalLess(year.Grades[i].Grade, year.Grades[j].Grade)
			})
			for gradeIdx := range year.Grades {
				grade := &year.Grades[gradeIdx]
				sort.SliceStable(grade.Classes, func(i, j int) bool {
					return naturalLess(grade.Classes[i].Name, grade.Classes[j].Name)
				})
			}
		}
	}
	return schools
}

func saveSchoolTreeSchool(ctx context.Context, db *sql.DB, input schoolTreeSchoolInput, auditCtx requestAuditContext) (masterDataSchoolOption, error) {
	var saved masterDataSchoolOption
	if input.ID == "" {
		err := db.QueryRowContext(ctx, `
INSERT INTO schools (code, name, status, created_by_user_id, updated_by_user_id)
VALUES ($1, $2, $3, nullif($4, '')::uuid, nullif($4, '')::uuid)
RETURNING id::text, code, name, status`,
			input.Code, input.Name, input.Status, auditCtx.ActorUserID).Scan(&saved.ID, &saved.Code, &saved.Name, &saved.Status)
		return saved, err
	}
	err := db.QueryRowContext(ctx, `
UPDATE schools
SET code = $2,
	name = $3,
	status = $4,
	updated_by_user_id = nullif($5, '')::uuid
WHERE id = $1::uuid
RETURNING id::text, code, name, status`,
		input.ID, input.Code, input.Name, input.Status, auditCtx.ActorUserID).Scan(&saved.ID, &saved.Code, &saved.Name, &saved.Status)
	return saved, err
}

func saveSchoolTreeSchoolYear(ctx context.Context, db *sql.DB, input schoolTreeSchoolYearInput, auditCtx requestAuditContext) (masterDataSchoolYearOption, error) {
	var saved masterDataSchoolYearOption
	if input.ID == "" {
		err := db.QueryRowContext(ctx, `
INSERT INTO school_years (school_id, code, name, starts_on, ends_on, status, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2, $3, nullif($4, '')::date, nullif($5, '')::date, $6, nullif($7, '')::uuid, nullif($7, '')::uuid)
RETURNING id::text, school_id::text, (SELECT code FROM schools WHERE id = school_id), code, name, status`,
			input.SchoolID, input.Code, input.Name, input.StartsOn, input.EndsOn, input.Status, auditCtx.ActorUserID).Scan(
			&saved.ID, &saved.SchoolID, &saved.SchoolCode, &saved.Code, &saved.Name, &saved.Status,
		)
		return saved, err
	}
	err := db.QueryRowContext(ctx, `
UPDATE school_years
SET school_id = $2::uuid,
	code = $3,
	name = $4,
	starts_on = nullif($5, '')::date,
	ends_on = nullif($6, '')::date,
	status = $7,
	updated_by_user_id = nullif($8, '')::uuid
WHERE id = $1::uuid
RETURNING id::text, school_id::text, (SELECT code FROM schools WHERE id = school_id), code, name, status`,
		input.ID, input.SchoolID, input.Code, input.Name, input.StartsOn, input.EndsOn, input.Status, auditCtx.ActorUserID).Scan(
		&saved.ID, &saved.SchoolID, &saved.SchoolCode, &saved.Code, &saved.Name, &saved.Status,
	)
	return saved, err
}

func saveSchoolTreeClass(ctx context.Context, db *sql.DB, input schoolTreeClassInput, auditCtx requestAuditContext) (masterDataClassOption, error) {
	var saved masterDataClassOption
	if input.ID == "" {
		err := db.QueryRowContext(ctx, `
INSERT INTO classes (school_year_id, grade, name, status, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2, $3, $4, nullif($5, '')::uuid, nullif($5, '')::uuid)
RETURNING id::text,
	(SELECT school_id::text FROM school_years WHERE id = school_year_id),
	(SELECT sc.code FROM school_years sy JOIN schools sc ON sc.id = sy.school_id WHERE sy.id = school_year_id),
	school_year_id::text,
	(SELECT code FROM school_years WHERE id = school_year_id),
	grade,
	name,
	status`,
			input.SchoolYearID, input.Grade, input.Name, input.Status, auditCtx.ActorUserID).Scan(
			&saved.ID,
			&saved.SchoolID,
			&saved.SchoolCode,
			&saved.SchoolYearID,
			&saved.SchoolYearCode,
			&saved.Grade,
			&saved.Name,
			&saved.Status,
		)
		return saved, err
	}
	err := db.QueryRowContext(ctx, `
UPDATE classes
SET school_year_id = $2::uuid,
	grade = $3,
	name = $4,
	status = $5,
	updated_by_user_id = nullif($6, '')::uuid
WHERE id = $1::uuid
RETURNING id::text,
	(SELECT school_id::text FROM school_years WHERE id = school_year_id),
	(SELECT sc.code FROM school_years sy JOIN schools sc ON sc.id = sy.school_id WHERE sy.id = school_year_id),
	school_year_id::text,
	(SELECT code FROM school_years WHERE id = school_year_id),
	grade,
	name,
	status`,
		input.ID, input.SchoolYearID, input.Grade, input.Name, input.Status, auditCtx.ActorUserID).Scan(
		&saved.ID,
		&saved.SchoolID,
		&saved.SchoolCode,
		&saved.SchoolYearID,
		&saved.SchoolYearCode,
		&saved.Grade,
		&saved.Name,
		&saved.Status,
	)
	return saved, err
}

func normalizeSchoolTreeSchoolInput(input schoolTreeSchoolInput) schoolTreeSchoolInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Code = schoolCodeOrDefault(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.Status = normalizeSchoolTreeStatus(input.Status)
	return input
}

func normalizeSchoolTreeSchoolYearInput(input schoolTreeSchoolYearInput) schoolTreeSchoolYearInput {
	input.ID = strings.TrimSpace(input.ID)
	input.SchoolID = strings.TrimSpace(input.SchoolID)
	input.Code = normalizeSchoolYearCode(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		input.Name = input.Code
	}
	input.Status = normalizeSchoolTreeStatus(input.Status)
	input.StartsOn = strings.TrimSpace(input.StartsOn)
	input.EndsOn = strings.TrimSpace(input.EndsOn)
	return input
}

func normalizeSchoolTreeClassInput(input schoolTreeClassInput) schoolTreeClassInput {
	input.ID = strings.TrimSpace(input.ID)
	input.SchoolYearID = strings.TrimSpace(input.SchoolYearID)
	input.Grade = normalizeGrade(input.Grade)
	input.Name = strings.TrimSpace(input.Name)
	input.Status = normalizeSchoolTreeStatus(input.Status)
	return input
}

func validateSchoolTreeSchoolInput(input schoolTreeSchoolInput) error {
	if input.Code == "" {
		return errors.New("school code is required")
	}
	if input.Name == "" {
		return errors.New("school name is required")
	}
	return validateSchoolTreeStatus(input.Status)
}

func validateSchoolTreeSchoolYearInput(input schoolTreeSchoolYearInput) error {
	if input.SchoolID == "" {
		return errors.New("schoolId is required")
	}
	if input.Code == "" {
		return errors.New("school year code is required")
	}
	if input.Name == "" {
		return errors.New("school year name is required")
	}
	return validateSchoolTreeStatus(input.Status)
}

func validateSchoolTreeClassInput(input schoolTreeClassInput) error {
	if input.SchoolYearID == "" {
		return errors.New("schoolYearId is required")
	}
	if input.Grade == "" {
		return errors.New("grade is required")
	}
	if input.Name == "" {
		return errors.New("class name is required")
	}
	return validateSchoolTreeStatus(input.Status)
}

func normalizeSchoolTreeStatus(status string) string {
	status = headerKey(status)
	if status == "" {
		return "active"
	}
	return status
}

func validateSchoolTreeStatus(status string) error {
	switch status {
	case "active", "archived":
		return nil
	default:
		return fmt.Errorf("status must be active or archived")
	}
}

func naturalLess(a string, b string) bool {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b)) < 0
}
