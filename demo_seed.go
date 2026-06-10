package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

const (
	demoSeedName                   = "finance-hub"
	demoSeedTimeout                = 3 * time.Minute
	demoDataMigrationTable         = "demo_data_migrations"
	financeHubDemoMigrationName    = "finance_hub_demo"
	financeHubDemoMigrationVersion = "0001"
	financeHubDemoMigrationSalt    = "finance_hub_demo_v1"
	demoTenantCode                 = "SUNRISE_DEMO"
	demoTenantName                 = "Sunrise International School Demo"
	demoSchoolCode                 = "SUNRISE"
	demoSchoolName                 = "Sunrise International School"
	demoOwnerEmail                 = "owner.demo@example.com"
	demoOwnerPhone                 = "0909000001"
	demoOwnerName                  = "Demo Tenant Owner"
	demoOwnerPassword              = "DemoOwner@2026!"
	demoSchoolYearCode             = "2025-2026"
	demoPeriodCode                 = "2025-04"
	demoCollectionBIN              = "970415"
	demoCollectionAcct             = "FHCOLLECT001"
)

type demoSubscriptionPlanSeed struct {
	Code        string
	Name        string
	Description string
	Limits      map[string]int
}

type demoSeedSummary struct {
	TenantID                string
	TenantCode              string
	TenantName              string
	SchoolCode              string
	OwnerID                 string
	OwnerEmail              string
	OwnerPhone              string
	OwnerPassword           string
	PlanCode                string
	StudentRows             int
	FeeScheduleCount        int
	InvoiceCount            int
	PaymentTransactionCount int
	NotificationCount       int
}

type demoFeeProfileRow struct {
	Profile      string
	Grade        string
	ClassName    string
	FeeTypeCode  string
	LabelVI      string
	LabelEN      string
	Amount       int
	DisplayOrder int
}

type demoSeedCommandOptions struct {
	Refresh bool
}

type demoDataMigration struct {
	Version  string
	Name     string
	Checksum string
	Apply    func(context.Context, *sql.DB) (demoSeedSummary, error)
}

type appliedDemoDataMigration struct {
	Version     string
	Name        string
	Checksum    string
	AppliedAt   time.Time
	RefreshedAt time.Time
}

type demoDataMigrationRunResult struct {
	Migration demoDataMigration
	Summary   *demoSeedSummary
	Status    string
}

func runDemoCommand(ctx context.Context, args []string, stdout io.Writer) error {
	options, err := parseDemoSeedCommandOptions(args)
	if err != nil {
		return err
	}

	cfg, err := loadDatabaseConfig()
	if err != nil {
		return err
	}
	if err := cfg.requireURL(); err != nil {
		return err
	}
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		return err
	}

	commandCtx, cancel := context.WithTimeout(ctx, demoSeedTimeout)
	defer cancel()

	db, err := openConfiguredDatabase(commandCtx, cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := databaseMigrationsReady(commandCtx, db, cfg.MigrationsTable, migrations); err != nil {
		return err
	}

	migration, err := financeHubDemoDataMigration()
	if err != nil {
		return err
	}
	result, err := runDemoDataMigration(commandCtx, db, migration, options.Refresh)
	if err != nil {
		return err
	}
	printDemoDataMigrationResult(stdout, result)
	return nil
}

func parseDemoSeedCommandOptions(args []string) (demoSeedCommandOptions, error) {
	if len(args) < 2 || len(args) > 3 || args[0] != "seed" || args[1] != demoSeedName {
		return demoSeedCommandOptions{}, fmt.Errorf("usage: go run . demo seed finance-hub [--refresh]")
	}
	options := demoSeedCommandOptions{}
	if len(args) == 3 {
		if args[2] != "--refresh" {
			return demoSeedCommandOptions{}, fmt.Errorf("usage: go run . demo seed finance-hub [--refresh]")
		}
		options.Refresh = true
	}
	return options, nil
}

func runDefaultDemoDataMigrations(ctx context.Context, db *sql.DB, stdout io.Writer) error {
	migration, err := financeHubDemoDataMigration()
	if err != nil {
		return err
	}
	result, err := runDemoDataMigration(ctx, db, migration, false)
	if err != nil {
		return err
	}
	printDemoDataMigrationResult(stdout, result)
	return nil
}

func financeHubDemoDataMigration() (demoDataMigration, error) {
	checksum, err := financeHubDemoDataMigrationChecksum()
	if err != nil {
		return demoDataMigration{}, err
	}
	return demoDataMigration{
		Version:  financeHubDemoMigrationVersion,
		Name:     financeHubDemoMigrationName,
		Checksum: checksum,
		Apply:    seedFinanceHubDemo,
	}, nil
}

func financeHubDemoDataMigrationChecksum() (string, error) {
	hash := sha256.New()
	hash.Write([]byte(financeHubDemoMigrationSalt + "\n"))
	for _, item := range []struct {
		Name string
		Data any
	}{
		{Name: "subscription_plans", Data: demoSubscriptionPlans()},
		{Name: "master_data", Data: financeHubDemoMasterDataRows()},
		{Name: "fee_profiles", Data: demoFeeProfiles()},
		{Name: "fee_adjustments", Data: demoFeeAdjustments()},
	} {
		payload, err := json.Marshal(item.Data)
		if err != nil {
			return "", err
		}
		hash.Write([]byte(item.Name + "\n"))
		hash.Write(payload)
		hash.Write([]byte("\n"))
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func runDemoDataMigration(ctx context.Context, db *sql.DB, migration demoDataMigration, refresh bool) (demoDataMigrationRunResult, error) {
	if migration.Version == "" || migration.Name == "" || migration.Checksum == "" || migration.Apply == nil {
		return demoDataMigrationRunResult{}, fmt.Errorf("invalid demo data migration")
	}
	if err := ensureDemoDataMigrationTable(ctx, db); err != nil {
		return demoDataMigrationRunResult{}, err
	}
	applied, err := loadAppliedDemoDataMigrations(ctx, db)
	if err != nil {
		return demoDataMigrationRunResult{}, err
	}

	if existing, ok := applied[migration.Version]; ok {
		if existing.Checksum != migration.Checksum {
			return demoDataMigrationRunResult{}, fmt.Errorf("demo data migration %s was already applied with a different checksum", migration.Version)
		}
		if !refresh {
			return demoDataMigrationRunResult{
				Migration: migration,
				Status:    "skipped",
			}, nil
		}
		summary, err := migration.Apply(ctx, db)
		if err != nil {
			return demoDataMigrationRunResult{}, err
		}
		if err := markDemoDataMigrationRefreshed(ctx, db, migration.Version); err != nil {
			return demoDataMigrationRunResult{}, err
		}
		return demoDataMigrationRunResult{
			Migration: migration,
			Summary:   &summary,
			Status:    "refreshed",
		}, nil
	}

	summary, err := migration.Apply(ctx, db)
	if err != nil {
		return demoDataMigrationRunResult{}, err
	}
	if err := recordAppliedDemoDataMigration(ctx, db, migration); err != nil {
		return demoDataMigrationRunResult{}, err
	}
	return demoDataMigrationRunResult{
		Migration: migration,
		Summary:   &summary,
		Status:    "applied",
	}, nil
}

func ensureDemoDataMigrationTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS demo_data_migrations (
	version text PRIMARY KEY,
	name text NOT NULL,
	checksum text NOT NULL,
	applied_at timestamptz NOT NULL DEFAULT now(),
	refreshed_at timestamptz NOT NULL DEFAULT now()
)`)
	return err
}

func loadAppliedDemoDataMigrations(ctx context.Context, db *sql.DB) (map[string]appliedDemoDataMigration, error) {
	rows, err := db.QueryContext(ctx, `
SELECT version, name, checksum, applied_at, refreshed_at
FROM demo_data_migrations
ORDER BY version`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := map[string]appliedDemoDataMigration{}
	for rows.Next() {
		var item appliedDemoDataMigration
		if err := rows.Scan(&item.Version, &item.Name, &item.Checksum, &item.AppliedAt, &item.RefreshedAt); err != nil {
			return nil, err
		}
		applied[item.Version] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return applied, nil
}

func recordAppliedDemoDataMigration(ctx context.Context, db *sql.DB, migration demoDataMigration) error {
	_, err := db.ExecContext(ctx, `
INSERT INTO demo_data_migrations (version, name, checksum)
VALUES ($1, $2, $3)`, migration.Version, migration.Name, migration.Checksum)
	return err
}

func markDemoDataMigrationRefreshed(ctx context.Context, db *sql.DB, version string) error {
	_, err := db.ExecContext(ctx, `
UPDATE demo_data_migrations
SET refreshed_at = now()
WHERE version = $1`, version)
	return err
}

func printDemoDataMigrationResult(out io.Writer, result demoDataMigrationRunResult) {
	fmt.Fprintf(out, "demo_data_migration: %s %s %s\n", result.Status, result.Migration.Version, result.Migration.Name)
	if result.Status == "skipped" {
		fmt.Fprintln(out, "refresh_command: go run . demo seed finance-hub --refresh")
		return
	}
	if result.Summary != nil {
		printDemoSeedSummary(out, *result.Summary)
	}
}

func printDemoSeedSummary(out io.Writer, summary demoSeedSummary) {
	fmt.Fprintln(out, "finance hub demo seed complete")
	fmt.Fprintf(out, "tenant: %s (%s)\n", summary.TenantName, summary.TenantCode)
	fmt.Fprintf(out, "school: %s\n", summary.SchoolCode)
	fmt.Fprintf(out, "owner_login: %s / %s\n", summary.OwnerEmail, summary.OwnerPassword)
	fmt.Fprintf(out, "owner_phone: %s\n", summary.OwnerPhone)
	fmt.Fprintf(out, "subscription_plan: %s\n", summary.PlanCode)
	fmt.Fprintf(out, "students_import_rows: %d\n", summary.StudentRows)
	fmt.Fprintf(out, "fee_schedules: %d\n", summary.FeeScheduleCount)
	fmt.Fprintf(out, "invoices: %d\n", summary.InvoiceCount)
	fmt.Fprintf(out, "payment_transactions: %d\n", summary.PaymentTransactionCount)
	fmt.Fprintf(out, "notification_campaigns: %d\n", summary.NotificationCount)
}

func seedFinanceHubDemo(ctx context.Context, db *sql.DB) (demoSeedSummary, error) {
	if err := seedDemoSubscriptionPlans(ctx, db); err != nil {
		return demoSeedSummary{}, err
	}
	tenantID, ownerID, err := seedDemoTenantOwner(ctx, db)
	if err != nil {
		return demoSeedSummary{}, err
	}
	if err := resetDemoTenantOperationalData(ctx, db, tenantID); err != nil {
		return demoSeedSummary{}, err
	}
	if err := seedDemoPaymentProviderConfig(ctx, db, tenantID, ownerID); err != nil {
		return demoSeedSummary{}, err
	}

	rows, err := loadFinanceHubDemoMasterDataRows()
	if err != nil {
		return demoSeedSummary{}, err
	}
	importResponse, err := importMasterDataRows(ctx, db, rows, true, tenantID)
	if err != nil {
		return demoSeedSummary{}, err
	}
	if len(importResponse.Issues) > 0 {
		return demoSeedSummary{}, fmt.Errorf("demo master data has %d import issues; first issue: %s", len(importResponse.Issues), importResponse.Issues[0].Message)
	}

	feeScheduleIDs, err := seedDemoFeeSchedules(ctx, db, tenantID, ownerID)
	if err != nil {
		return demoSeedSummary{}, err
	}
	invoiceCount, err := seedDemoInvoices(ctx, db, tenantID, feeScheduleIDs)
	if err != nil {
		return demoSeedSummary{}, err
	}
	paymentCount, err := seedDemoPayments(ctx, db, tenantID, ownerID)
	if err != nil {
		return demoSeedSummary{}, err
	}
	notificationCount, err := seedDemoNotifications(ctx, db, tenantID)
	if err != nil {
		return demoSeedSummary{}, err
	}
	if err := seedDemoOperationLog(ctx, db, tenantID); err != nil {
		return demoSeedSummary{}, err
	}
	if err := rebuildTenantUsageCounter(ctx, db, tenantID, subscriptionMetricSchools, time.Now()); err != nil {
		return demoSeedSummary{}, err
	}
	if err := rebuildTenantUsageCounter(ctx, db, tenantID, subscriptionMetricOperators, time.Now()); err != nil {
		return demoSeedSummary{}, err
	}
	if err := rebuildTenantUsageCounter(ctx, db, tenantID, subscriptionMetricStudents, time.Now()); err != nil {
		return demoSeedSummary{}, err
	}

	return demoSeedSummary{
		TenantID:                tenantID,
		TenantCode:              demoTenantCode,
		TenantName:              demoTenantName,
		SchoolCode:              demoSchoolCode,
		OwnerID:                 ownerID,
		OwnerEmail:              demoOwnerEmail,
		OwnerPhone:              demoOwnerPhone,
		OwnerPassword:           demoOwnerPassword,
		PlanCode:                "pro",
		StudentRows:             len(rows),
		FeeScheduleCount:        len(feeScheduleIDs),
		InvoiceCount:            invoiceCount,
		PaymentTransactionCount: paymentCount,
		NotificationCount:       notificationCount,
	}, nil
}

func demoSubscriptionPlans() []demoSubscriptionPlanSeed {
	return []demoSubscriptionPlanSeed{
		{
			Code:        "free",
			Name:        "Free",
			Description: "Demo starter plan for small trial workspaces",
			Limits: map[string]int{
				subscriptionMetricSchools:              1,
				subscriptionMetricOperators:            2,
				subscriptionMetricStudents:             100,
				subscriptionMetricMonthlyNotifications: 100,
			},
		},
		{
			Code:        "go",
			Name:        "Go",
			Description: "Demo entry plan for one active school operation",
			Limits: map[string]int{
				subscriptionMetricSchools:              1,
				subscriptionMetricOperators:            5,
				subscriptionMetricStudents:             300,
				subscriptionMetricMonthlyNotifications: 1000,
			},
		},
		{
			Code:        "plus",
			Name:        "Plus",
			Description: "Demo growth plan for multi-team finance operations",
			Limits: map[string]int{
				subscriptionMetricSchools:              3,
				subscriptionMetricOperators:            20,
				subscriptionMetricStudents:             1500,
				subscriptionMetricMonthlyNotifications: 5000,
			},
		},
		{
			Code:        "pro",
			Name:        "Pro",
			Description: "Demo full management plan for customer presentations",
			Limits: map[string]int{
				subscriptionMetricSchools:              10,
				subscriptionMetricOperators:            100,
				subscriptionMetricStudents:             5000,
				subscriptionMetricMonthlyNotifications: 10000,
			},
		},
	}
}

func financeHubDemoMasterDataRows() []masterDataImportRow {
	return []masterDataImportRow{
		{RowNumber: 2, StudentCode: "FH-S001", StudentName: "An Nguyen Minh", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "1", ClassName: "1.01", ParentName: "Linh Pham", ParentEmail: "parent-fh01@example.com", ParentPhone: "0902000001", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 3, StudentCode: "FH-S002", StudentName: "Bao Tran Gia", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "1", ClassName: "1.01", ParentName: "Nam Tran", ParentEmail: "parent-fh02@example.com", ParentPhone: "0902000002", Relationship: "father", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 4, StudentCode: "FH-S003", StudentName: "Chi Le Bao", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "1", ClassName: "1.01", ParentName: "Hoa Le", ParentEmail: "parent-fh03a@example.com", ParentPhone: "0902000003", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 5, StudentCode: "FH-S003", StudentName: "Chi Le Bao", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "1", ClassName: "1.01", ParentName: "Phuc Le", ParentEmail: "parent-fh03b@example.com", ParentPhone: "0902000004", Relationship: "father", IsPrimary: false, ParentActive: true, ReceivesBillingEmail: false},
		{RowNumber: 6, StudentCode: "FH-S004", StudentName: "Khue Do Gia", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "1", ClassName: "1.01", ParentName: "Thao Do", ParentEmail: "", ParentPhone: "0902000005", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: false},
		{RowNumber: 7, StudentCode: "FH-S005", StudentName: "Duc Pham Anh", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "2", ClassName: "2.01", ParentName: "Huong Pham", ParentEmail: "parent-fh05@example.com", ParentPhone: "0902000006", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 8, StudentCode: "FH-S006", StudentName: "Em Pham Minh", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "2", ClassName: "2.01", ParentName: "Huong Pham", ParentEmail: "parent-fh05@example.com", ParentPhone: "0902000006", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 9, StudentCode: "FH-S007", StudentName: "Gia Vu Long", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "2", ClassName: "2.01", ParentName: "Quang Vu", ParentEmail: "parent-fh07@example.com", ParentPhone: "0902000007", Relationship: "father", IsPrimary: true, ParentActive: false, ReceivesBillingEmail: false},
		{RowNumber: 10, StudentCode: "FH-S008", StudentName: "Ha Hoang Yen", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.02", ParentName: "Mai Hoang", ParentEmail: "parent-fh08@example.com", ParentPhone: "0902000008", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 11, StudentCode: "FH-S009", StudentName: "Khang Nguyen Bao", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.02", ParentName: "Son Nguyen", ParentEmail: "parent-fh09@example.com", ParentPhone: "0902000009", Relationship: "father", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 12, StudentCode: "FH-S010", StudentName: "Lan Tran Ngoc", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.02", ParentName: "Hanh Tran", ParentEmail: "parent-fh10@example.com", ParentPhone: "0902000010", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 13, StudentCode: "FH-S011", StudentName: "An Nguyen Minh", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.03", ParentName: "Van Nguyen", ParentEmail: "parent-fh11@example.com", ParentPhone: "0902000011", Relationship: "father", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 14, StudentCode: "FH-S012", StudentName: "Nhi Dao Mai", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.03", ParentName: "Thu Dao", ParentEmail: "parent-fh12@example.com", ParentPhone: "0902000012", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 15, StudentCode: "FH-S013", StudentName: "Oanh Bui Khanh", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.03", ParentName: "Tuan Bui", ParentEmail: "parent-fh13@example.com", ParentPhone: "0902000013", Relationship: "father", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 16, StudentCode: "FH-S014", StudentName: "Phuc Dang Khoa", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "3", ClassName: "3.03", ParentName: "Lam Dang", ParentEmail: "", ParentPhone: "0902000014", Relationship: "guardian", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: false},
		{RowNumber: 17, StudentCode: "FH-S015", StudentName: "Quynh Vo Anh", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "4", ClassName: "4.01", ParentName: "Ngoc Vo", ParentEmail: "parent-fh15@example.com", ParentPhone: "0902000015", Relationship: "mother", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
		{RowNumber: 18, StudentCode: "FH-S016", StudentName: "Son Ly Thai", SchoolCode: demoSchoolCode, SchoolYearCode: demoSchoolYearCode, Grade: "4", ClassName: "4.01", ParentName: "Binh Ly", ParentEmail: "parent-fh16@example.com", ParentPhone: "0902000016", Relationship: "father", IsPrimary: true, ParentActive: true, ReceivesBillingEmail: true},
	}
}

func demoFeeProfiles() []demoFeeProfileRow {
	return []demoFeeProfileRow{
		{Profile: "grade_1_apr_may", Grade: "1", ClassName: "1.01", FeeTypeCode: "tuition", LabelVI: "Hoc phi T04-T05", LabelEN: "Tuition Apr-May", Amount: 8400000, DisplayOrder: 10},
		{Profile: "grade_1_apr_may", Grade: "1", ClassName: "1.01", FeeTypeCode: "shuttle", LabelVI: "Phi xe T04", LabelEN: "Shuttle Apr", Amount: 1200000, DisplayOrder: 20},
		{Profile: "grade_1_apr_may", Grade: "1", ClassName: "1.01", FeeTypeCode: "insurance", LabelVI: "Bao hiem y te", LabelEN: "Health insurance", Amount: 563850, DisplayOrder: 30},
		{Profile: "grade_1_apr_may", Grade: "1", ClassName: "1.01", FeeTypeCode: "materials", LabelVI: "Sach CTQT", LabelEN: "International material", Amount: 800000, DisplayOrder: 40},
		{Profile: "grade_2_apr_may", Grade: "2", ClassName: "2.01", FeeTypeCode: "tuition", LabelVI: "Hoc phi T04-T05", LabelEN: "Tuition Apr-May", Amount: 9000000, DisplayOrder: 10},
		{Profile: "grade_2_apr_may", Grade: "2", ClassName: "2.01", FeeTypeCode: "shuttle", LabelVI: "Phi xe T04", LabelEN: "Shuttle Apr", Amount: 1500000, DisplayOrder: 20},
		{Profile: "grade_2_apr_may", Grade: "2", ClassName: "2.01", FeeTypeCode: "insurance", LabelVI: "Bao hiem y te", LabelEN: "Health insurance", Amount: 563850, DisplayOrder: 30},
		{Profile: "grade_2_apr_may", Grade: "2", ClassName: "2.01", FeeTypeCode: "materials", LabelVI: "Sach CTQT", LabelEN: "International material", Amount: 950000, DisplayOrder: 40},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.02", FeeTypeCode: "tuition", LabelVI: "Hoc phi T04-T05", LabelEN: "Tuition Apr-May", Amount: 9600000, DisplayOrder: 10},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.02", FeeTypeCode: "shuttle", LabelVI: "Phi xe T04", LabelEN: "Shuttle Apr", Amount: 1700000, DisplayOrder: 20},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.02", FeeTypeCode: "insurance", LabelVI: "Bao hiem y te", LabelEN: "Health insurance", Amount: 563850, DisplayOrder: 30},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.02", FeeTypeCode: "materials", LabelVI: "Sach CTQT", LabelEN: "International material", Amount: 1200000, DisplayOrder: 40},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.03", FeeTypeCode: "tuition", LabelVI: "Hoc phi T04-T05", LabelEN: "Tuition Apr-May", Amount: 9600000, DisplayOrder: 10},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.03", FeeTypeCode: "shuttle", LabelVI: "Phi xe T04", LabelEN: "Shuttle Apr", Amount: 1700000, DisplayOrder: 20},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.03", FeeTypeCode: "insurance", LabelVI: "Bao hiem y te", LabelEN: "Health insurance", Amount: 563850, DisplayOrder: 30},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.03", FeeTypeCode: "uniform", LabelVI: "Dong phuc", LabelEN: "Uniform", Amount: 850000, DisplayOrder: 40},
		{Profile: "grade_3_apr_may", Grade: "3", ClassName: "3.03", FeeTypeCode: "materials", LabelVI: "Sach CTQT", LabelEN: "International material", Amount: 1200000, DisplayOrder: 50},
		{Profile: "grade_4_apr_may", Grade: "4", ClassName: "4.01", FeeTypeCode: "tuition", LabelVI: "Hoc phi T04-T05", LabelEN: "Tuition Apr-May", Amount: 10200000, DisplayOrder: 10},
		{Profile: "grade_4_apr_may", Grade: "4", ClassName: "4.01", FeeTypeCode: "lunch", LabelVI: "Phi an ban tru", LabelEN: "Lunch", Amount: 1800000, DisplayOrder: 20},
		{Profile: "grade_4_apr_may", Grade: "4", ClassName: "4.01", FeeTypeCode: "materials", LabelVI: "Sach CTQT", LabelEN: "International material", Amount: 1300000, DisplayOrder: 30},
	}
}

func demoFeeAdjustments() []studentFeeAdjustmentInput {
	return []studentFeeAdjustmentInput{
		{StudentCode: "FH-S003", AdjustmentType: adjustmentTypeDiscount, FeeTypeCode: "tuition", Amount: 500000, Reason: "Sibling discount"},
		{StudentCode: "FH-S004", AdjustmentType: adjustmentTypeWaiver, FeeTypeCode: "shuttle", Amount: 0, Reason: "No shuttle service"},
		{StudentCode: "FH-S005", AdjustmentType: adjustmentTypeDiscount, FeeTypeCode: "tuition", Amount: 400000, Reason: "Sibling discount"},
		{StudentCode: "FH-S006", AdjustmentType: adjustmentTypeDiscount, FeeTypeCode: "tuition", Amount: 400000, Reason: "Sibling discount"},
		{StudentCode: "FH-S007", AdjustmentType: adjustmentTypeCarryOver, FeeTypeCode: "previous_fees", Amount: 600000, Reason: "Prior balance"},
		{StudentCode: "FH-S010", AdjustmentType: adjustmentTypeSurcharge, FeeTypeCode: "materials", Amount: 250000, Reason: "Extra workbook set"},
		{StudentCode: "FH-S011", AdjustmentType: adjustmentTypeWaiver, FeeTypeCode: "uniform", Amount: 0, Reason: "Uniform purchased earlier"},
		{StudentCode: "FH-S012", AdjustmentType: adjustmentTypeDiscount, FeeTypeCode: "tuition", Amount: 300000, Reason: "Early payment policy"},
		{StudentCode: "FH-S014", AdjustmentType: adjustmentTypeSurcharge, FeeTypeCode: "custom", Amount: 200000, Reason: "Late bus registration"},
		{StudentCode: "FH-S015", AdjustmentType: adjustmentTypeDiscount, FeeTypeCode: "lunch", Amount: 250000, Reason: "Lunch scholarship"},
		{StudentCode: "FH-S016", AdjustmentType: adjustmentTypeCarryOver, FeeTypeCode: "previous_fees", Amount: 450000, Reason: "Prior balance"},
	}
}

func seedDemoSubscriptionPlans(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, plan := range demoSubscriptionPlans() {
		limits, err := json.Marshal(plan.Limits)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO subscription_plans (code, name, status, description, limits)
VALUES ($1, $2, 'active', $3, $4::jsonb)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	status = 'active',
	description = EXCLUDED.description,
	limits = EXCLUDED.limits,
	updated_at = now()`,
			plan.Code,
			plan.Name,
			plan.Description,
			string(limits),
		); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `
UPDATE subscription_plans
SET status = 'archived',
	updated_at = now()
WHERE code IN ('free_trial', 'standard')`); err != nil {
		return err
	}
	return tx.Commit()
}

func seedDemoTenantOwner(ctx context.Context, db *sql.DB) (string, string, error) {
	ownerInput := adminUserSaveInput{
		Email:       demoOwnerEmail,
		Phone:       demoOwnerPhone,
		DisplayName: demoOwnerName,
		Status:      "active",
		Password:    demoOwnerPassword,
	}
	if err := validateAdminUserSaveInput(&ownerInput); err != nil {
		return "", "", err
	}
	passwordHash, err := hashPassword(ownerInput.Password)
	if err != nil {
		return "", "", err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	var ownerID string
	if err := tx.QueryRowContext(ctx, `
WITH updated AS (
	UPDATE app_users
	SET email = $1,
		phone = $2,
		display_name = $3,
		status = 'active',
		password_hash = $4,
		password_updated_at = now(),
		updated_at = now()
	WHERE lower(COALESCE(email, '')) = lower($1)
		OR phone = $2
	RETURNING id::text
), inserted AS (
	INSERT INTO app_users (email, phone, display_name, status, password_hash, password_updated_at)
	SELECT $1, $2, $3, 'active', $4, now()
	WHERE NOT EXISTS (SELECT 1 FROM updated)
	RETURNING id::text
)
SELECT id FROM updated
UNION ALL
SELECT id FROM inserted
LIMIT 1`,
		ownerInput.Email,
		ownerInput.Phone,
		ownerInput.DisplayName,
		passwordHash,
	).Scan(&ownerID); err != nil {
		return "", "", err
	}

	var tenantID string
	if err := tx.QueryRowContext(ctx, `
INSERT INTO tenants (code, name, status, created_by_user_id, updated_by_user_id)
VALUES ($1, $2, 'active', $3::uuid, $3::uuid)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	status = 'active',
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()
RETURNING id::text`, demoTenantCode, demoTenantName, ownerID).Scan(&tenantID); err != nil {
		return "", "", err
	}
	if err := ensureTenantMembership(ctx, tx, tenantID, ownerID, true); err != nil {
		return "", "", err
	}
	if err := seedDemoTenantOwnerFullRolePermissions(ctx, tx); err != nil {
		return "", "", err
	}
	if err := ensureTenantUserRole(ctx, tx, tenantID, ownerID, "tenant_owner", ownerID); err != nil {
		return "", "", err
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO schools (tenant_id, code, name, status, created_by_user_id, updated_by_user_id)
VALUES ($1::uuid, $2, $3, 'active', $4::uuid, $4::uuid)
ON CONFLICT (tenant_id, code) DO UPDATE
SET name = EXCLUDED.name,
	status = 'active',
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()`, tenantID, demoSchoolCode, demoSchoolName, ownerID); err != nil {
		return "", "", err
	}
	if err := seedDemoTenantSubscription(ctx, tx, tenantID, ownerID); err != nil {
		return "", "", err
	}
	if err := ensureDefaultPaymentProvidersForTenant(ctx, tx, tenantID, ownerID); err != nil {
		return "", "", err
	}
	if err := insertAuditLog(ctx, tx, auditLogInput{
		Context:    demoAuditContext(tenantID, ownerID),
		Action:     "demo.seed.tenant",
		EntityType: "tenant",
		EntityID:   tenantID,
		Metadata: map[string]any{
			"tenantCode": demoTenantCode,
			"schoolCode": demoSchoolCode,
			"planCode":   "pro",
		},
	}); err != nil {
		return "", "", err
	}
	return tenantID, ownerID, tx.Commit()
}

func seedDemoTenantOwnerFullRolePermissions(ctx context.Context, exec masterDataExecutor) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO app_role_permissions (role_id, permission_id)
SELECT role.id, permission.id
FROM app_roles role
JOIN app_permissions permission ON permission.code IN (
	'dashboard.view',
	'user.view',
	'user.create',
	'user.update',
	'user.assign_role',
	'role.view',
	'student.view',
	'student.create',
	'student.update',
	'school_tree.view',
	'school_tree.update',
	'fee.view',
	'fee.create',
	'fee.update',
	'invoice.view',
	'invoice.create',
	'invoice.update',
	'payment.view',
	'payment.create',
	'payment.reconcile',
	'notification.view',
	'notification.create',
	'notification.send',
	'email_config.view',
	'email_config.update',
	'email_cron.view',
	'email_cron.update',
	'report.view',
	'report.export',
	'operation_log.view',
	'audit_log.view',
	'tenant.switch',
	'subscription.view',
	'subscription.update'
)
WHERE role.code = 'tenant_owner'
ON CONFLICT (role_id, permission_id) DO NOTHING`)
	return err
}

func seedDemoTenantSubscription(ctx context.Context, exec masterDataExecutor, tenantID string, ownerID string) error {
	metadata := map[string]any{
		"amount":                3990000,
		"interval_months":       1,
		"due_days":              10,
		"auto_renew":            true,
		"renewal_mode":          "auto_generate",
		"finance_note":          "Demo Pro monthly subscription",
		"automation_enabled":    true,
		"renewal_lead_days":     7,
		"dunning_enabled":       true,
		"dunning_interval_days": 3,
		"suspend_enabled":       true,
		"suspend_after_days":    14,
	}
	metadataJSON, err := jsonObjectString(metadata)
	if err != nil {
		return err
	}
	_, err = exec.ExecContext(ctx, `
INSERT INTO tenant_subscriptions (
	tenant_id, plan_id, status, current_period_starts_at, current_period_ends_at,
	billing_metadata, created_by_user_id, updated_by_user_id
)
SELECT $1::uuid, plan.id, 'active', '2026-01-01'::timestamptz, '2026-12-31'::timestamptz,
	$3::jsonb, $2::uuid, $2::uuid
FROM subscription_plans plan
WHERE plan.code = 'pro'
ON CONFLICT (tenant_id) DO UPDATE
SET plan_id = EXCLUDED.plan_id,
	status = 'active',
	trial_ends_at = NULL,
	current_period_starts_at = EXCLUDED.current_period_starts_at,
	current_period_ends_at = EXCLUDED.current_period_ends_at,
	billing_metadata = EXCLUDED.billing_metadata,
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()`, tenantID, ownerID, metadataJSON)
	return err
}

func seedDemoPaymentProviderConfig(ctx context.Context, db *sql.DB, tenantID string, ownerID string) error {
	if _, err := db.ExecContext(ctx, `
UPDATE payment_providers
SET config = config || $3::jsonb,
	status = CASE WHEN code = 'manual_vietqr' THEN 'active' ELSE status END,
	updated_by_user_id = $2::uuid,
	updated_at = now()
WHERE tenant_id = $1::uuid
	AND code IN ('manual_vietqr', 'sepay')`,
		tenantID,
		ownerID,
		`{"bankBin":"970415","bank_bin":"970415","accountNumber":"FHCOLLECT001","account_number":"FHCOLLECT001","accountName":"SUNRISE DEMO","account_name":"SUNRISE DEMO","collection":{"bankBin":"970415","accountNumber":"FHCOLLECT001","accountName":"SUNRISE DEMO"}}`,
	); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `
INSERT INTO tenant_payment_settings (tenant_id, default_provider_code, updated_by_user_id)
VALUES ($1::uuid, 'manual_vietqr', $2::uuid)
ON CONFLICT (tenant_id) DO UPDATE
SET default_provider_code = 'manual_vietqr',
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()`, tenantID, ownerID); err != nil {
		return err
	}
	return nil
}

func resetDemoTenantOperationalData(ctx context.Context, db *sql.DB, tenantID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []string{
		`DELETE FROM notification_campaigns WHERE tenant_id = $1::uuid`,
		`DELETE FROM manual_cash_receipts m
			USING invoices i, school_years sy, schools sc
			WHERE m.invoice_id = i.id AND i.school_year_id = sy.id AND sy.school_id = sc.id AND sc.tenant_id = $1::uuid`,
		`DELETE FROM reconciliation_matches rm
			USING invoices i, school_years sy, schools sc
			WHERE rm.invoice_id = i.id AND i.school_year_id = sy.id AND sy.school_id = sc.id AND sc.tenant_id = $1::uuid`,
		`DELETE FROM payment_transactions pt
			USING payment_providers pp
			WHERE pt.provider_id = pp.id AND pp.tenant_id = $1::uuid`,
		`DELETE FROM payment_intents pi
			USING invoices i, school_years sy, schools sc
			WHERE pi.invoice_id = i.id AND i.school_year_id = sy.id AND sy.school_id = sc.id AND sc.tenant_id = $1::uuid`,
		`DELETE FROM invoices i
			USING school_years sy, schools sc
			WHERE i.school_year_id = sy.id AND sy.school_id = sc.id AND sc.tenant_id = $1::uuid`,
		`DELETE FROM operation_logs
			WHERE tenant_id = $1::uuid AND operation LIKE 'demo.seed%'`,
	}
	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt, tenantID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func loadFinanceHubDemoMasterDataRows() ([]masterDataImportRow, error) {
	rows := financeHubDemoMasterDataRows()
	return rows, nil
}

func seedDemoFeeSchedules(ctx context.Context, db *sql.DB, tenantID string, ownerID string) ([]string, error) {
	profiles := demoFeeProfiles()
	adjustments := demoFeeAdjustments()
	adjustmentsByClass, err := demoFeeAdjustmentsByClass(ctx, db, tenantID, adjustments)
	if err != nil {
		return nil, err
	}

	type classKey struct {
		Grade     string
		ClassName string
	}
	grouped := map[classKey][]demoFeeProfileRow{}
	for _, profile := range profiles {
		key := classKey{Grade: profile.Grade, ClassName: profile.ClassName}
		grouped[key] = append(grouped[key], profile)
	}
	keys := make([]classKey, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Grade != keys[j].Grade {
			return keys[i].Grade < keys[j].Grade
		}
		return keys[i].ClassName < keys[j].ClassName
	})

	scheduleIDs := []string{}
	for _, key := range keys {
		schoolYearID, classID, err := loadDemoClassScope(ctx, db, tenantID, key.Grade, key.ClassName)
		if err != nil {
			return nil, err
		}
		input := feeScheduleInput{
			TenantID:     tenantID,
			SchoolYearID: schoolYearID,
			ClassID:      classID,
			PeriodCode:   demoPeriodCode,
			Month:        4,
			Name:         "Demo " + key.ClassName + " - " + demoPeriodCode,
			Notes:        "Finance Hub demo fee schedule",
			Status:       feeScheduleStatusActive,
			OperatorName: demoOwnerName,
			Items:        demoFeeScheduleItems(grouped[key]),
			Adjustments:  adjustmentsByClass[key.ClassName],
		}
		if existingID, err := findExistingDemoFeeScheduleID(ctx, db, tenantID, schoolYearID, classID); err != nil {
			return nil, err
		} else {
			input.ID = existingID
		}
		id, err := saveFeeSchedule(ctx, db, input, demoAuditContext(tenantID, ownerID))
		if err != nil {
			return nil, err
		}
		scheduleIDs = append(scheduleIDs, id)
	}
	return scheduleIDs, nil
}

func demoFeeAdjustmentsByClass(ctx context.Context, db *sql.DB, tenantID string, adjustments []studentFeeAdjustmentInput) (map[string][]studentFeeAdjustmentInput, error) {
	out := map[string][]studentFeeAdjustmentInput{}
	for _, adjustment := range adjustments {
		className, err := loadDemoStudentClassName(ctx, db, tenantID, adjustment.StudentCode)
		if err != nil {
			return nil, err
		}
		out[className] = append(out[className], adjustment)
	}
	return out, nil
}

func loadDemoStudentClassName(ctx context.Context, db *sql.DB, tenantID string, studentCode string) (string, error) {
	var className string
	err := db.QueryRowContext(ctx, `
SELECT c.name
FROM students s
JOIN classes c ON c.id = s.class_id
WHERE s.tenant_id = $1::uuid
	AND s.student_code = $2`, tenantID, studentCode).Scan(&className)
	return className, err
}

func loadDemoClassScope(ctx context.Context, db *sql.DB, tenantID string, grade string, className string) (string, string, error) {
	var schoolYearID string
	var classID string
	err := db.QueryRowContext(ctx, `
SELECT sy.id::text, c.id::text
FROM classes c
JOIN school_years sy ON sy.id = c.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
	AND sc.code = $2
	AND sy.code = $3
	AND c.grade = $4
	AND c.name = $5`,
		tenantID,
		demoSchoolCode,
		demoSchoolYearCode,
		grade,
		className,
	).Scan(&schoolYearID, &classID)
	return schoolYearID, classID, err
}

func findExistingDemoFeeScheduleID(ctx context.Context, db *sql.DB, tenantID string, schoolYearID string, classID string) (string, error) {
	var id string
	err := db.QueryRowContext(ctx, `
SELECT fs.id::text
FROM fee_schedules fs
JOIN school_years sy ON sy.id = fs.school_year_id
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
	AND fs.school_year_id = $2::uuid
	AND fs.scope_type = 'class'
	AND fs.class_id = $3::uuid
	AND fs.period_code = $4
	AND COALESCE(fs.month, 0) = 4
	AND fs.status <> 'archived'
LIMIT 1`, tenantID, schoolYearID, classID, demoPeriodCode).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

func demoFeeScheduleItems(rows []demoFeeProfileRow) []feeScheduleItemInput {
	items := make([]feeScheduleItemInput, 0, len(rows))
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].DisplayOrder < rows[j].DisplayOrder
	})
	for _, row := range rows {
		items = append(items, feeScheduleItemInput{
			FeeTypeCode:  row.FeeTypeCode,
			LabelVI:      row.LabelVI,
			LabelEN:      row.LabelEN,
			Amount:       row.Amount,
			DisplayOrder: row.DisplayOrder,
		})
	}
	return items
}

func seedDemoInvoices(ctx context.Context, db *sql.DB, tenantID string, feeScheduleIDs []string) (int, error) {
	count := 0
	for _, scheduleID := range feeScheduleIDs {
		input := normalizeInvoiceGenerateInput(invoiceGenerateInput{
			TenantID:      tenantID,
			FeeScheduleID: scheduleID,
			BankBIN:       demoCollectionBIN,
			BankAccount:   demoCollectionAcct,
			IssueDate:     "2026-04-05",
			DueDate:       "2026-04-20",
		})
		preview, err := buildInvoicePreviewFromDB(ctx, db, input)
		if err != nil {
			return 0, err
		}
		saved, err := saveGeneratedInvoices(ctx, db, input, preview)
		if err != nil {
			return 0, err
		}
		count += saved.Summary.GeneratedCount
	}
	invoices, err := listInvoiceSummaries(ctx, db, invoiceListFilters{TenantID: tenantID, PeriodCode: demoPeriodCode})
	if err != nil {
		return 0, err
	}
	provider, err := loadPaymentProviderByCode(ctx, db, tenantID, paymentProviderManualVietQR)
	if err != nil {
		return 0, err
	}
	for _, invoice := range invoices {
		document, err := loadInvoiceDocument(ctx, db, invoice.ID, tenantID)
		if err != nil {
			return 0, err
		}
		if _, err := createPaymentIntentForInvoice(ctx, db, provider, document); err != nil {
			return 0, err
		}
	}
	return len(invoices), nil
}

func seedDemoPayments(ctx context.Context, db *sql.DB, tenantID string, ownerID string) (int, error) {
	invoices, err := listInvoiceSummaries(ctx, db, invoiceListFilters{TenantID: tenantID, PeriodCode: demoPeriodCode})
	if err != nil {
		return 0, err
	}
	byStudent := map[string]invoiceSummary{}
	for _, invoice := range invoices {
		byStudent[invoice.StudentCode] = invoice
	}
	cases := []struct {
		StudentCode string
		Provider    string
		TxnID       string
		ReceiptRef  string
		Amount      func(invoiceSummary) int
		Reason      string
		PaidAt      string
	}{
		{StudentCode: "FH-S001", Provider: paymentProviderSePay, TxnID: "FH-DEMO-SEPAY-001", Amount: func(i invoiceSummary) int { return i.TotalAmount }, Reason: "SePay demo exact payment", PaidAt: "2026-04-08T09:12:00+07:00"},
		{StudentCode: "FH-S003", Provider: paymentProviderManualVietQR, TxnID: "cash:FH-CASH-003", ReceiptRef: "FH-CASH-003", Amount: func(invoiceSummary) int { return 2500000 }, Reason: "Phu huynh nop mot phan", PaidAt: "2026-04-09T10:30:00+07:00"},
		{StudentCode: "FH-S005", Provider: paymentProviderManualVietQR, TxnID: "cash:FH-CASH-005", ReceiptRef: "FH-CASH-005", Amount: func(i invoiceSummary) int { return i.TotalAmount }, Reason: "Thu tien mat tai van phong", PaidAt: "2026-04-10T14:15:00+07:00"},
		{StudentCode: "FH-S009", Provider: paymentProviderSePay, TxnID: "FH-DEMO-SEPAY-009", Amount: func(i invoiceSummary) int { return i.TotalAmount + 500000 }, Reason: "SePay demo overpayment", PaidAt: "2026-04-11T08:20:00+07:00"},
		{StudentCode: "FH-S012", Provider: paymentProviderManualVietQR, TxnID: "cash:FH-CASH-012", ReceiptRef: "FH-CASH-012", Amount: func(i invoiceSummary) int { return i.TotalAmount + 600000 }, Reason: "Phu huynh nop du tien va du phu thu", PaidAt: "2026-04-12T09:45:00+07:00"},
	}
	count := 0
	for _, item := range cases {
		invoice, ok := byStudent[item.StudentCode]
		if !ok {
			return 0, fmt.Errorf("demo invoice for %s not found", item.StudentCode)
		}
		if err := seedDemoMatchedPayment(ctx, db, tenantID, ownerID, invoice, item.Provider, item.TxnID, item.ReceiptRef, item.Amount(invoice), item.Reason, item.PaidAt); err != nil {
			return 0, err
		}
		count++
	}
	if err := seedDemoUnmatchedPayment(ctx, db, tenantID); err != nil {
		return 0, err
	}
	return count + 1, nil
}

func seedDemoMatchedPayment(ctx context.Context, db *sql.DB, tenantID string, ownerID string, invoice invoiceSummary, providerCode string, transactionID string, receiptRef string, amount int, reason string, paidAt string) error {
	provider, err := loadPaymentProviderByCode(ctx, db, tenantID, providerCode)
	if err != nil {
		return err
	}
	paidTime := parseProviderTime(paidAt, time.Now())
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rawPayload, _ := jsonObjectString(map[string]any{
		"source":                "demo_seed",
		"provider":              providerCode,
		"providerTransactionId": transactionID,
		"receiptReference":      receiptRef,
		"invoiceCode":           invoice.InvoiceCode,
	})
	var transactionRowID string
	if err := tx.QueryRowContext(ctx, `
INSERT INTO payment_transactions (
	provider_id, invoice_id, provider_transaction_id, direction, amount, currency,
	transaction_time, account_number, bank_bin, bank_name, description, reference_code, status, raw_payload,
	created_by_user_id, updated_by_user_id
)
VALUES ($1::uuid, $2::uuid, $3, 'in', $4, 'VND', $5, $6, $7, $8, $9, $10, 'matched', $11::jsonb, $12::uuid, $12::uuid)
ON CONFLICT (provider_id, provider_transaction_id) WHERE provider_transaction_id <> '' DO UPDATE
SET invoice_id = EXCLUDED.invoice_id,
	amount = EXCLUDED.amount,
	transaction_time = EXCLUDED.transaction_time,
	description = EXCLUDED.description,
	reference_code = EXCLUDED.reference_code,
	status = 'matched',
	raw_payload = EXCLUDED.raw_payload,
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()
RETURNING id::text`,
		provider.ID,
		invoice.ID,
		transactionID,
		amount,
		paidTime,
		invoice.CollectionBankAccount,
		invoice.CollectionBankBIN,
		firstNonEmpty(provider.DisplayName, providerCode),
		"Demo payment "+invoice.InvoiceCode,
		firstNonEmpty(receiptRef, transactionID),
		rawPayload,
		ownerID,
	).Scan(&transactionRowID); err != nil {
		return err
	}
	metadata, _ := jsonObjectString(map[string]any{
		"source":           "demo_seed",
		"provider":         providerCode,
		"receiptReference": receiptRef,
	})
	if _, err := tx.ExecContext(ctx, `
INSERT INTO reconciliation_matches (invoice_id, transaction_id, match_type, status, score, amount_applied, reason, metadata, created_by_user_id)
VALUES ($1::uuid, $2::uuid, $3, 'matched', 100, $4, $5, $6::jsonb, $7::uuid)
ON CONFLICT (transaction_id, invoice_id) WHERE status <> 'reversed' DO UPDATE
SET match_type = EXCLUDED.match_type,
	status = 'matched',
	score = 100,
	amount_applied = EXCLUDED.amount_applied,
	reason = EXCLUDED.reason,
	metadata = EXCLUDED.metadata`,
		invoice.ID,
		transactionRowID,
		demoPaymentMatchType(providerCode, receiptRef),
		amount,
		reason,
		metadata,
		ownerID,
	); err != nil {
		return err
	}
	if receiptRef != "" {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO manual_cash_receipts (
	invoice_id, payment_transaction_id, collector_user_id, collector_name,
	amount, currency, paid_at, receipt_reference, reason, note, created_by_user_id
)
VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, 'VND', $6, $7, $8, $9, $3::uuid)
ON CONFLICT (receipt_reference) DO UPDATE
SET invoice_id = EXCLUDED.invoice_id,
	payment_transaction_id = EXCLUDED.payment_transaction_id,
	collector_user_id = EXCLUDED.collector_user_id,
	collector_name = EXCLUDED.collector_name,
	amount = EXCLUDED.amount,
	paid_at = EXCLUDED.paid_at,
	reason = EXCLUDED.reason,
	note = EXCLUDED.note`,
			invoice.ID,
			transactionRowID,
			ownerID,
			"Thu Ngan Demo",
			amount,
			paidTime,
			receiptRef,
			reason,
			"Demo cash collection",
		); err != nil {
			return err
		}
	}
	if _, err := refreshInvoicePaymentStatus(ctx, tx, invoice.ID, reason); err != nil {
		return err
	}
	if err := insertAuditLog(ctx, tx, auditLogInput{
		Context:    demoAuditContext(tenantID, ownerID),
		Action:     "demo.seed.payment",
		EntityType: "payment_transaction",
		EntityID:   transactionRowID,
		Reason:     reason,
		Metadata: map[string]any{
			"invoiceCode": invoice.InvoiceCode,
			"amount":      amount,
			"provider":    providerCode,
		},
	}); err != nil {
		return err
	}
	return tx.Commit()
}

func demoPaymentMatchType(providerCode string, receiptRef string) string {
	if receiptRef != "" {
		return "cash"
	}
	if providerCode == paymentProviderSePay {
		return "provider_reference"
	}
	return "manual"
}

func seedDemoUnmatchedPayment(ctx context.Context, db *sql.DB, tenantID string) error {
	provider, err := loadPaymentProviderByCode(ctx, db, tenantID, paymentProviderSePay)
	if err != nil {
		return err
	}
	rawPayload, _ := jsonObjectString(map[string]any{
		"source":                "demo_seed",
		"providerTransactionId": "FH-DEMO-SEPAY-UNMATCHED",
		"description":           "Chuyen khoan hoc phi khong co ma hoa don",
	})
	_, err = db.ExecContext(ctx, `
INSERT INTO payment_transactions (
	provider_id, provider_transaction_id, direction, amount, currency,
	transaction_time, account_number, bank_name, description, reference_code, status, raw_payload
)
VALUES ($1::uuid, 'FH-DEMO-SEPAY-UNMATCHED', 'in', 5000000, 'VND',
	'2026-04-13T08:20:00+07:00'::timestamptz, $2, 'DEMO_BANK',
	'Chuyen khoan hoc phi khong co ma hoa don', 'FH-DEMO-SEPAY-UNMATCHED', 'manual_review', $3::jsonb)
ON CONFLICT (provider_id, provider_transaction_id) WHERE provider_transaction_id <> '' DO UPDATE
SET amount = EXCLUDED.amount,
	transaction_time = EXCLUDED.transaction_time,
	account_number = EXCLUDED.account_number,
	description = EXCLUDED.description,
	reference_code = EXCLUDED.reference_code,
	status = 'manual_review',
	raw_payload = EXCLUDED.raw_payload,
	updated_at = now()`,
		provider.ID,
		demoCollectionAcct,
		rawPayload,
	)
	return err
}

func seedDemoNotifications(ctx context.Context, db *sql.DB, tenantID string) (int, error) {
	schoolYearID, err := loadDemoSchoolYearID(ctx, db, tenantID)
	if err != nil {
		return 0, err
	}
	inputs := []notificationCampaignInput{
		{
			TenantID:      tenantID,
			Name:          "Sunrise Demo first notice 2025-04",
			CampaignType:  notificationCampaignFirstNotice,
			SchoolYearID:  schoolYearID,
			PeriodCode:    demoPeriodCode,
			InvoiceStatus: invoiceStatusUnpaid,
			DryRun:        true,
		},
		{
			TenantID:      tenantID,
			Name:          "Sunrise Demo reminder partial 2025-04",
			CampaignType:  notificationCampaignReminder,
			SchoolYearID:  schoolYearID,
			PeriodCode:    demoPeriodCode,
			InvoiceStatus: invoiceStatusPartial,
			DueOnOrBefore: "2026-04-20",
			DryRun:        true,
			ForceResend:   false,
		},
	}
	count := 0
	for _, input := range inputs {
		input = normalizeNotificationCampaignInput(input)
		preview, err := buildNotificationPreviewFromDB(ctx, db, input)
		if err != nil {
			return 0, err
		}
		if len(preview.Issues) > 0 {
			return 0, fmt.Errorf("demo notification %q has issue: %s", input.Name, preview.Issues[0].Message)
		}
		campaign, err := saveNotificationCampaign(ctx, db, input, preview.Template, preview.Recipients)
		if err != nil {
			return 0, err
		}
		input.CampaignID = campaign.ID
		recipients, err := loadNotificationRecipients(ctx, db, campaign.ID, tenantID)
		if err != nil {
			return 0, err
		}
		cfg := defaultEmailConfig()
		cfg.SchoolName = demoSchoolName
		cfg.SchoolNameEN = "Sunrise International School"
		if _, err := sendNotificationCampaign(ctx, db, cfg, preview.Template, input, recipients, "http://localhost:18080", 0); err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

func loadDemoSchoolYearID(ctx context.Context, db *sql.DB, tenantID string) (string, error) {
	var id string
	err := db.QueryRowContext(ctx, `
SELECT sy.id::text
FROM school_years sy
JOIN schools sc ON sc.id = sy.school_id
WHERE sc.tenant_id = $1::uuid
	AND sc.code = $2
	AND sy.code = $3`, tenantID, demoSchoolCode, demoSchoolYearCode).Scan(&id)
	return id, err
}

func seedDemoOperationLog(ctx context.Context, db *sql.DB, tenantID string) error {
	return recordOperationLog(ctx, db, operationLogInput{
		TenantID:   tenantID,
		Source:     "demo",
		Level:      "warn",
		Operation:  "demo.seed.payment_reconciliation",
		Status:     "manual_review",
		Message:    "Demo unmatched transfer ready for reconciliation review",
		EntityType: "payment_transaction",
		Metadata: map[string]any{
			"providerTransactionId": "FH-DEMO-SEPAY-UNMATCHED",
			"tenantCode":            demoTenantCode,
		},
	})
}

func demoAuditContext(tenantID string, ownerID string) requestAuditContext {
	return requestAuditContext{
		TenantID:    tenantID,
		ActorUserID: ownerID,
		ActorName:   demoOwnerName,
		RequestID:   "demo-seed-finance-hub",
		UserAgent:   "dekisugi-demo-seed",
	}
}
