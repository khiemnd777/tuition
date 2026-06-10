package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDemoSubscriptionPlansUseCustomerFacingPlanSet(t *testing.T) {
	plans := demoSubscriptionPlans()
	want := []string{"free", "go", "plus", "pro"}
	if len(plans) != len(want) {
		t.Fatalf("expected %d demo plans, got %d", len(want), len(plans))
	}
	for idx, code := range want {
		if plans[idx].Code != code {
			t.Fatalf("expected plan %d to be %q, got %q", idx, code, plans[idx].Code)
		}
		if plans[idx].Name == "" || plans[idx].Description == "" {
			t.Fatalf("expected plan %q to include display copy", plans[idx].Code)
		}
		for _, metric := range subscriptionUsageMetricCodes() {
			if plans[idx].Limits[metric] <= 0 {
				t.Fatalf("expected plan %q to define positive %s limit", plans[idx].Code, metric)
			}
		}
	}
	if plans[0].Limits[subscriptionMetricStudents] >= plans[len(plans)-1].Limits[subscriptionMetricStudents] {
		t.Fatalf("expected Pro student limit to exceed Free")
	}
}

func TestDemoSeedOwnerCredentialIsFictionalAndValid(t *testing.T) {
	input := adminUserSaveInput{
		Email:       demoOwnerEmail,
		Phone:       demoOwnerPhone,
		DisplayName: demoOwnerName,
		Status:      "active",
		Password:    demoOwnerPassword,
	}
	if err := validateAdminUserSaveInput(&input); err != nil {
		t.Fatalf("expected demo owner credential to validate: %v", err)
	}
	if input.Email != "owner.demo@example.com" {
		t.Fatalf("expected example.com demo email, got %q", input.Email)
	}
}

func TestFinanceHubDemoSeedRowsTargetSunriseSchool(t *testing.T) {
	rows, err := loadFinanceHubDemoMasterDataRows()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 17 {
		t.Fatalf("expected 17 finance hub demo rows, got %d", len(rows))
	}
	for _, row := range rows {
		if row.SchoolCode != demoSchoolCode {
			t.Fatalf("expected row %s to target %s, got %s", row.StudentCode, demoSchoolCode, row.SchoolCode)
		}
		if row.ParentEmail != "" && !hasExampleDomain(row.ParentEmail) {
			t.Fatalf("expected fictional parent email for %s, got %q", row.StudentCode, row.ParentEmail)
		}
	}
}

func TestParseDemoSeedCommandOptions(t *testing.T) {
	options, err := parseDemoSeedCommandOptions([]string{"seed", "finance-hub"})
	if err != nil {
		t.Fatal(err)
	}
	if options.Refresh {
		t.Fatal("expected default demo seed command to apply once without refresh")
	}

	options, err = parseDemoSeedCommandOptions([]string{"seed", "finance-hub", "--refresh"})
	if err != nil {
		t.Fatal(err)
	}
	if !options.Refresh {
		t.Fatal("expected --refresh to enable demo seed refresh")
	}

	if _, err := parseDemoSeedCommandOptions([]string{"seed", "finance-hub", "--force"}); err == nil {
		t.Fatal("expected unsupported demo seed flag to fail")
	}
}

func TestRunDemoDataMigrationAppliesSkipsAndRefreshes(t *testing.T) {
	db, state := openFakeDemoDataMigrationDB(t, "apply-skip-refresh", nil)
	applyCount := 0
	migration := demoDataMigration{
		Version:  "0001",
		Name:     "finance_hub_demo",
		Checksum: "checksum-1",
		Apply: func(ctx context.Context, db *sql.DB) (demoSeedSummary, error) {
			applyCount++
			return demoSeedSummary{TenantCode: demoTenantCode}, nil
		},
	}

	result, err := runDemoDataMigration(context.Background(), db, migration, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "applied" || result.Summary == nil {
		t.Fatalf("expected applied result with summary, got %+v", result)
	}
	if applyCount != 1 {
		t.Fatalf("expected first run to apply once, got %d", applyCount)
	}
	if _, ok := state.applied["0001"]; !ok {
		t.Fatal("expected migration registry row after first run")
	}

	result, err = runDemoDataMigration(context.Background(), db, migration, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "skipped" || result.Summary != nil {
		t.Fatalf("expected skipped result without summary, got %+v", result)
	}
	if applyCount != 1 {
		t.Fatalf("expected second run without refresh to skip apply, got %d", applyCount)
	}

	result, err = runDemoDataMigration(context.Background(), db, migration, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "refreshed" || result.Summary == nil {
		t.Fatalf("expected refreshed result with summary, got %+v", result)
	}
	if applyCount != 2 {
		t.Fatalf("expected refresh to reapply once, got %d", applyCount)
	}
	if state.refreshCount != 1 {
		t.Fatalf("expected registry refresh count 1, got %d", state.refreshCount)
	}
}

func TestRunDemoDataMigrationRejectsChecksumDrift(t *testing.T) {
	db, _ := openFakeDemoDataMigrationDB(t, "checksum-drift", map[string]appliedDemoDataMigration{
		"0001": {
			Version:     "0001",
			Name:        "finance_hub_demo",
			Checksum:    "old-checksum",
			AppliedAt:   time.Now(),
			RefreshedAt: time.Now(),
		},
	})
	applyCount := 0
	migration := demoDataMigration{
		Version:  "0001",
		Name:     "finance_hub_demo",
		Checksum: "new-checksum",
		Apply: func(ctx context.Context, db *sql.DB) (demoSeedSummary, error) {
			applyCount++
			return demoSeedSummary{}, nil
		},
	}

	err := func() error {
		_, err := runDemoDataMigration(context.Background(), db, migration, true)
		return err
	}()
	if err == nil || !strings.Contains(err.Error(), "different checksum") {
		t.Fatalf("expected checksum drift error, got %v", err)
	}
	if applyCount != 0 {
		t.Fatalf("expected checksum drift to block apply, got %d applies", applyCount)
	}
}

func hasExampleDomain(email string) bool {
	return len(email) >= len("@example.com") && email[len(email)-len("@example.com"):] == "@example.com"
}

var registerFakeDemoDataMigrationDriver sync.Once
var fakeDemoDataMigrationStates sync.Map

type fakeDemoDataMigrationState struct {
	mu           sync.Mutex
	applied      map[string]appliedDemoDataMigration
	refreshCount int
}

func openFakeDemoDataMigrationDB(t *testing.T, name string, applied map[string]appliedDemoDataMigration) (*sql.DB, *fakeDemoDataMigrationState) {
	t.Helper()

	registerFakeDemoDataMigrationDriver.Do(func() {
		sql.Register("fake_demo_data_migrations", fakeDemoDataMigrationDriver{})
	})

	copiedApplied := map[string]appliedDemoDataMigration{}
	for version, item := range applied {
		copiedApplied[version] = item
	}
	state := &fakeDemoDataMigrationState{applied: copiedApplied}
	fakeDemoDataMigrationStates.Store(name, state)

	db, err := sql.Open("fake_demo_data_migrations", name)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		fakeDemoDataMigrationStates.Delete(name)
	})
	return db, state
}

type fakeDemoDataMigrationDriver struct{}

func (fakeDemoDataMigrationDriver) Open(name string) (driver.Conn, error) {
	value, ok := fakeDemoDataMigrationStates.Load(name)
	if !ok {
		return nil, errors.New("unknown fake demo data migration database")
	}
	return &fakeDemoDataMigrationConn{state: value.(*fakeDemoDataMigrationState)}, nil
}

type fakeDemoDataMigrationConn struct {
	state *fakeDemoDataMigrationState
}

func (conn *fakeDemoDataMigrationConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not implemented")
}

func (conn *fakeDemoDataMigrationConn) Close() error {
	return nil
}

func (conn *fakeDemoDataMigrationConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not implemented")
}

func (conn *fakeDemoDataMigrationConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	conn.state.mu.Lock()
	defer conn.state.mu.Unlock()

	normalized := strings.TrimSpace(query)
	switch {
	case strings.HasPrefix(normalized, "CREATE TABLE IF NOT EXISTS demo_data_migrations"):
		return driver.RowsAffected(0), nil
	case strings.HasPrefix(normalized, "INSERT INTO demo_data_migrations"):
		version := args[0].Value.(string)
		conn.state.applied[version] = appliedDemoDataMigration{
			Version:     version,
			Name:        args[1].Value.(string),
			Checksum:    args[2].Value.(string),
			AppliedAt:   time.Now(),
			RefreshedAt: time.Now(),
		}
		return driver.RowsAffected(1), nil
	case strings.HasPrefix(normalized, "UPDATE demo_data_migrations"):
		version := args[0].Value.(string)
		item := conn.state.applied[version]
		item.RefreshedAt = time.Now()
		conn.state.applied[version] = item
		conn.state.refreshCount++
		return driver.RowsAffected(1), nil
	default:
		return nil, errors.New("unexpected exec: " + query)
	}
}

func (conn *fakeDemoDataMigrationConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if !strings.HasPrefix(strings.TrimSpace(query), "SELECT version, name, checksum, applied_at, refreshed_at") {
		return nil, errors.New("unexpected query: " + query)
	}

	conn.state.mu.Lock()
	defer conn.state.mu.Unlock()

	values := make([][]driver.Value, 0, len(conn.state.applied))
	for _, item := range conn.state.applied {
		values = append(values, []driver.Value{item.Version, item.Name, item.Checksum, item.AppliedAt, item.RefreshedAt})
	}
	return &fakeDemoDataMigrationRows{
		columns: []string{"version", "name", "checksum", "applied_at", "refreshed_at"},
		values:  values,
	}, nil
}

type fakeDemoDataMigrationRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (rows *fakeDemoDataMigrationRows) Columns() []string {
	return rows.columns
}

func (rows *fakeDemoDataMigrationRows) Close() error {
	return nil
}

func (rows *fakeDemoDataMigrationRows) Next(dest []driver.Value) error {
	if rows.index >= len(rows.values) {
		return io.EOF
	}
	copy(dest, rows.values[rows.index])
	rows.index++
	return nil
}
