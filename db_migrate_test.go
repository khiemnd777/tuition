package main

import (
	"bytes"
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

func TestLoadEmbeddedMigrationsIncludesFoundationSchema(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if migrations[0].Version != "0001" || migrations[0].Name != "foundation" {
		t.Fatalf("expected first foundation migration, got %+v", migrations[0])
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS app_users",
		"CREATE TABLE IF NOT EXISTS app_roles",
		"CREATE TABLE IF NOT EXISTS app_permissions",
		"CREATE TABLE IF NOT EXISTS audit_logs",
		"created_by_user_id",
		"updated_by_user_id",
		"abc_prevent_audit_log_mutation",
	} {
		if !strings.Contains(migrations[0].SQL, want) {
			t.Fatalf("expected foundation migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesMasterDataSchema(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var masterData migration
	for _, item := range migrations {
		if item.Version == "0002" {
			masterData = item
			break
		}
	}
	if masterData.Name != "student_parent_class_master_data" {
		t.Fatalf("expected master data migration 0002, got %+v", masterData)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS school_years",
		"CREATE TABLE IF NOT EXISTS classes",
		"CREATE TABLE IF NOT EXISTS students",
		"CREATE TABLE IF NOT EXISTS parents",
		"CREATE TABLE IF NOT EXISTS student_parents",
		"CONSTRAINT students_student_code_key UNIQUE (student_code)",
		"CREATE UNIQUE INDEX IF NOT EXISTS parents_email_key",
		"receives_billing_email boolean NOT NULL DEFAULT true",
		"student_parents_one_primary_per_student_idx",
	} {
		if !strings.Contains(masterData.SQL, want) {
			t.Fatalf("expected master data migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesFeeScheduleSchema(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var feeSchedule migration
	for _, item := range migrations {
		if item.Version == "0003" {
			feeSchedule = item
			break
		}
	}
	if feeSchedule.Name != "fee_types_and_fee_schedules" {
		t.Fatalf("expected fee schedule migration 0003, got %+v", feeSchedule)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS fee_types",
		"CREATE TABLE IF NOT EXISTS fee_schedules",
		"CREATE TABLE IF NOT EXISTS fee_schedule_items",
		"CREATE TABLE IF NOT EXISTS student_fee_adjustments",
		"label_vi text NOT NULL",
		"label_en text NOT NULL",
		"CONSTRAINT fee_schedules_scope_check",
		"CONSTRAINT student_fee_adjustments_type_check",
		"CONSTRAINT student_fee_adjustments_reason_not_blank",
		"fee_schedules.read",
		"fee_schedules.write",
	} {
		if !strings.Contains(feeSchedule.SQL, want) {
			t.Fatalf("expected fee schedule migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesInvoiceSchema(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var invoice migration
	for _, item := range migrations {
		if item.Version == "0004" {
			invoice = item
			break
		}
	}
	if invoice.Name != "invoices_and_receipts" {
		t.Fatalf("expected invoice migration 0004, got %+v", invoice)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS invoices",
		"CREATE TABLE IF NOT EXISTS invoice_items",
		"CREATE TABLE IF NOT EXISTS invoice_adjustments",
		"CREATE TABLE IF NOT EXISTS invoice_status_history",
		"CREATE TABLE IF NOT EXISTS receipt_documents",
		"qr_bill_number text NOT NULL",
		"CONSTRAINT invoices_status_check",
		"invoices_schedule_student_active_key",
		"invoices.read",
		"invoices.write",
	} {
		if !strings.Contains(invoice.SQL, want) {
			t.Fatalf("expected invoice migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesPaymentSchema(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var payment migration
	for _, item := range migrations {
		if item.Version == "0005" {
			payment = item
			break
		}
	}
	if payment.Name != "payments_and_reconciliation" {
		t.Fatalf("expected payment migration 0005, got %+v", payment)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS payment_providers",
		"CREATE TABLE IF NOT EXISTS payment_intents",
		"CREATE TABLE IF NOT EXISTS provider_events",
		"CREATE TABLE IF NOT EXISTS payment_transactions",
		"CREATE TABLE IF NOT EXISTS reconciliation_matches",
		"CREATE TABLE IF NOT EXISTS manual_cash_receipts",
		"provider_events_provider_event_id_key",
		"payment_transactions_provider_txn_key",
		"payments.reconcile",
		"'manual_vietqr'",
		"'sepay'",
		"'payos'",
	} {
		if !strings.Contains(payment.SQL, want) {
			t.Fatalf("expected payment migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesNotificationSchema(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var notification migration
	for _, item := range migrations {
		if item.Version == "0006" {
			notification = item
			break
		}
	}
	if notification.Name != "notification_campaigns" {
		t.Fatalf("expected notification migration 0006, got %+v", notification)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS notification_templates",
		"CREATE TABLE IF NOT EXISTS notification_campaigns",
		"CREATE TABLE IF NOT EXISTS notification_recipients",
		"CREATE TABLE IF NOT EXISTS notification_logs",
		"notification_logs_send_idempotency_key",
		"'first_notice'",
		"'reminder'",
		"notifications.read",
		"notifications.write",
		"notifications.send",
	} {
		if !strings.Contains(notification.SQL, want) {
			t.Fatalf("expected notification migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesWebAdminPermissions(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var admin migration
	for _, item := range migrations {
		if item.Version == "0007" {
			admin = item
			break
		}
	}
	if admin.Name != "web_admin" {
		t.Fatalf("expected web admin migration 0007, got %+v", admin)
	}

	for _, want := range []string{
		"admin.dashboard.read",
		"admin.reports.read",
		"system.users.assign_roles",
		"super_admin",
		"billing_admin",
		"viewer",
	} {
		if !strings.Contains(admin.SQL, want) {
			t.Fatalf("expected web admin migration to contain %q", want)
		}
	}
}

func TestRunMigrationsAppliesOnlyPendingMigrations(t *testing.T) {
	sqlText := "CREATE TABLE app_users (id uuid PRIMARY KEY);"
	migrations := []migration{{
		Version:  "0001",
		Name:     "foundation",
		FileName: "0001_foundation.sql",
		SQL:      sqlText,
		Checksum: migrationChecksum([]byte(sqlText)),
	}}
	db, state := openFakeMigrationDB(t, "apply-pending", nil)

	var out bytes.Buffer
	if err := runMigrations(context.Background(), db, defaultMigrationTable, migrations, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "applied 0001 foundation") {
		t.Fatalf("expected applied output, got %q", out.String())
	}
	if _, ok := state.applied["0001"]; !ok {
		t.Fatal("expected migration to be recorded")
	}
	if got := state.countExec(sqlText); got != 1 {
		t.Fatalf("expected migration SQL to run once, got %d", got)
	}

	out.Reset()
	if err := runMigrations(context.Background(), db, defaultMigrationTable, migrations, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "skip 0001 foundation") {
		t.Fatalf("expected skip output on repeated run, got %q", out.String())
	}
	if got := state.countExec(sqlText); got != 1 {
		t.Fatalf("expected repeated migration run to avoid data changes, got %d executions", got)
	}
}

func TestRunMigrationsDetectsChecksumDrift(t *testing.T) {
	sqlText := "CREATE TABLE app_users (id uuid PRIMARY KEY);"
	migrations := []migration{{
		Version:  "0001",
		Name:     "foundation",
		FileName: "0001_foundation.sql",
		SQL:      sqlText,
		Checksum: migrationChecksum([]byte(sqlText)),
	}}
	db, _ := openFakeMigrationDB(t, "checksum-drift", map[string]appliedMigration{
		"0001": {Version: "0001", Name: "foundation", Checksum: "old-checksum", AppliedAt: time.Now()},
	})

	err := runMigrations(context.Background(), db, defaultMigrationTable, migrations, io.Discard)
	if err == nil {
		t.Fatal("expected checksum drift error")
	}
	if !strings.Contains(err.Error(), "different checksum") {
		t.Fatalf("expected checksum drift error, got %v", err)
	}
}

var registerFakeMigrationDriver sync.Once
var fakeMigrationStates sync.Map

type fakeMigrationState struct {
	mu      sync.Mutex
	applied map[string]appliedMigration
	execs   []string
}

func (state *fakeMigrationState) countExec(query string) int {
	state.mu.Lock()
	defer state.mu.Unlock()

	count := 0
	for _, exec := range state.execs {
		if exec == query {
			count++
		}
	}
	return count
}

func openFakeMigrationDB(t *testing.T, name string, applied map[string]appliedMigration) (*sql.DB, *fakeMigrationState) {
	t.Helper()

	registerFakeMigrationDriver.Do(func() {
		sql.Register("fake_migrations", fakeMigrationDriver{})
	})

	copiedApplied := map[string]appliedMigration{}
	for version, item := range applied {
		copiedApplied[version] = item
	}
	state := &fakeMigrationState{applied: copiedApplied}
	fakeMigrationStates.Store(name, state)

	db, err := sql.Open("fake_migrations", name)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		fakeMigrationStates.Delete(name)
	})
	return db, state
}

type fakeMigrationDriver struct{}

func (fakeMigrationDriver) Open(name string) (driver.Conn, error) {
	value, ok := fakeMigrationStates.Load(name)
	if !ok {
		return nil, errors.New("unknown fake migration database")
	}
	return &fakeMigrationConn{state: value.(*fakeMigrationState)}, nil
}

type fakeMigrationConn struct {
	state *fakeMigrationState
}

func (conn *fakeMigrationConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not implemented")
}

func (conn *fakeMigrationConn) Close() error {
	return nil
}

func (conn *fakeMigrationConn) Begin() (driver.Tx, error) {
	return &fakeMigrationTx{}, nil
}

func (conn *fakeMigrationConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return &fakeMigrationTx{}, nil
}

func (conn *fakeMigrationConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	conn.state.mu.Lock()
	defer conn.state.mu.Unlock()

	conn.state.execs = append(conn.state.execs, query)
	if strings.HasPrefix(strings.TrimSpace(query), "INSERT INTO "+defaultMigrationTable) {
		conn.state.applied[args[0].Value.(string)] = appliedMigration{
			Version:   args[0].Value.(string),
			Name:      args[1].Value.(string),
			Checksum:  args[2].Value.(string),
			AppliedAt: time.Now(),
		}
	}
	return driver.RowsAffected(1), nil
}

func (conn *fakeMigrationConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if !strings.HasPrefix(strings.TrimSpace(query), "SELECT version, name, checksum, applied_at FROM "+defaultMigrationTable) {
		return nil, errors.New("unexpected query: " + query)
	}

	conn.state.mu.Lock()
	defer conn.state.mu.Unlock()

	values := make([][]driver.Value, 0, len(conn.state.applied))
	for _, item := range conn.state.applied {
		values = append(values, []driver.Value{item.Version, item.Name, item.Checksum, item.AppliedAt})
	}
	return &fakeMigrationRows{
		columns: []string{"version", "name", "checksum", "applied_at"},
		values:  values,
	}, nil
}

type fakeMigrationTx struct{}

func (fakeMigrationTx) Commit() error {
	return nil
}

func (fakeMigrationTx) Rollback() error {
	return nil
}

type fakeMigrationRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (rows *fakeMigrationRows) Columns() []string {
	return rows.columns
}

func (rows *fakeMigrationRows) Close() error {
	return nil
}

func (rows *fakeMigrationRows) Next(dest []driver.Value) error {
	if rows.index >= len(rows.values) {
		return io.EOF
	}
	copy(dest, rows.values[rows.index])
	rows.index++
	return nil
}
