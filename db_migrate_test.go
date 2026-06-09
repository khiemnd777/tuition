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

func TestLoadEmbeddedMigrationsIncludesReportsAuditOperations(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var operations migration
	for _, item := range migrations {
		if item.Version == "0008" {
			operations = item
			break
		}
	}
	if operations.Name != "reports_audit_operations" {
		t.Fatalf("expected operations migration 0008, got %+v", operations)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS operation_logs",
		"manual_cash_receipts_reason_not_blank",
		"admin.reports.export",
		"operations.read",
		"operation_logs_metadata_object",
	} {
		if !strings.Contains(operations.SQL, want) {
			t.Fatalf("expected reports/audit/operations migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesAuthSessions(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var auth migration
	for _, item := range migrations {
		if item.Version == "0009" {
			auth = item
			break
		}
	}
	if auth.Name != "auth_sessions" {
		t.Fatalf("expected auth migration 0009, got %+v", auth)
	}

	for _, want := range []string{
		"ADD COLUMN IF NOT EXISTS password_hash",
		"CREATE TABLE IF NOT EXISTS app_auth_sessions",
		"CREATE TABLE IF NOT EXISTS app_auth_access_tokens",
		"CREATE TABLE IF NOT EXISTS app_auth_refresh_tokens",
		"app_auth_access_tokens_hash_key",
		"app_auth_refresh_tokens_hash_key",
		"used_at timestamptz",
		"revoked_at timestamptz",
	} {
		if !strings.Contains(auth.SQL, want) {
			t.Fatalf("expected auth migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesRBACPermissions(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var rbac migration
	for _, item := range migrations {
		if item.Version == "0010" {
			rbac = item
			break
		}
	}
	if rbac.Name != "rbac_permissions" {
		t.Fatalf("expected RBAC migration 0010, got %+v", rbac)
	}

	for _, want := range []string{
		"email.config.read",
		"email.config.write",
		"email.send",
		"email.cron.manage",
		"super_admin",
		"billing_admin",
	} {
		if !strings.Contains(rbac.SQL, want) {
			t.Fatalf("expected RBAC migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesSchoolTreeManagement(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var schoolTree migration
	for _, item := range migrations {
		if item.Version == "0011" {
			schoolTree = item
			break
		}
	}
	if schoolTree.Name != "school_tree_management" {
		t.Fatalf("expected school tree migration 0011, got %+v", schoolTree)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS schools",
		"ALTER TABLE school_years ADD COLUMN IF NOT EXISTS school_id",
		"school_years_school_id_fkey",
		"school_years_school_code_key",
		"schools_set_updated_at",
		"DEKISUGI",
	} {
		if !strings.Contains(schoolTree.SQL, want) {
			t.Fatalf("expected school tree migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesUserContactsAndRoles(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var users migration
	for _, item := range migrations {
		if item.Version == "0012" {
			users = item
			break
		}
	}
	if users.Name != "user_contacts_and_roles" {
		t.Fatalf("expected users migration 0012, got %+v", users)
	}

	for _, want := range []string{
		"ADD COLUMN IF NOT EXISTS phone",
		"app_users_contact_required",
		"app_users_phone_key",
		"Admin / Quản trị viên",
		"Staff / Nhân sự",
		"Accountant / Kế toán",
		"user.view",
		"payment.reconcile",
		"email_cron.update",
	} {
		if !strings.Contains(users.SQL, want) {
			t.Fatalf("expected user contacts/roles migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesPaidConfirmationNotifications(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var paidConfirmation migration
	for _, item := range migrations {
		if item.Version == "0013" {
			paidConfirmation = item
			break
		}
	}
	if paidConfirmation.Name != "paid_confirmation_notifications" {
		t.Fatalf("expected paid confirmation migration 0013, got %+v", paidConfirmation)
	}

	for _, want := range []string{
		"payment_confirmation",
		"payment_paid",
		"auto_paid_confirmation",
		"notification_campaigns_type_check",
	} {
		if !strings.Contains(paidConfirmation.SQL, want) {
			t.Fatalf("expected paid confirmation migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantFoundation(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var tenant migration
	for _, item := range migrations {
		if item.Version == "0014" {
			tenant = item
			break
		}
	}
	if tenant.Name != "tenant_foundation" {
		t.Fatalf("expected tenant migration 0014, got %+v", tenant)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS tenants",
		"CREATE TABLE IF NOT EXISTS tenant_memberships",
		"ALTER TABLE schools ADD COLUMN IF NOT EXISTS tenant_id",
		"schools_tenant_code_key",
		"tenant_memberships_tenant_id_fkey",
		"tenant_memberships_user_id_fkey",
		"DEKISUGI",
	} {
		if !strings.Contains(tenant.SQL, want) {
			t.Fatalf("expected tenant foundation migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantAwareAuthRBAC(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var tenantAuth migration
	for _, item := range migrations {
		if item.Version == "0015" {
			tenantAuth = item
			break
		}
	}
	if tenantAuth.Name != "tenant_aware_auth_rbac" {
		t.Fatalf("expected tenant-aware auth/RBAC migration 0015, got %+v", tenantAuth)
	}

	for _, want := range []string{
		"ALTER TABLE app_auth_sessions ADD COLUMN IF NOT EXISTS tenant_id",
		"CREATE TABLE IF NOT EXISTS tenant_user_roles",
		"app_auth_sessions_tenant_id_fkey",
		"tenant_user_roles_tenant_id_fkey",
		"tenant_user_roles_user_id_fkey",
		"tenant_user_roles_role_id_fkey",
		"app_user_roles",
		"DEKISUGI",
	} {
		if !strings.Contains(tenantAuth.SQL, want) {
			t.Fatalf("expected tenant-aware auth/RBAC migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesActorModelRoleSplit(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var actorModel migration
	for _, item := range migrations {
		if item.Version == "0024" {
			actorModel = item
			break
		}
	}
	if actorModel.Name != "actor_model_role_split" {
		t.Fatalf("expected actor model role split migration 0024, got %+v", actorModel)
	}

	for _, want := range []string{
		"platform_admin",
		"tenant_owner",
		"tenant_admin",
		"tenant_staff",
		"tenant_accountant",
		"INSERT INTO app_user_roles",
		"INSERT INTO tenant_user_roles",
		"DELETE FROM tenant_user_roles",
	} {
		if !strings.Contains(actorModel.SQL, want) {
			t.Fatalf("expected actor model role split migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesSubscriptionPaymentConfirmationLink(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var item migration
	for _, current := range migrations {
		if current.Version == "0025" {
			item = current
			break
		}
	}
	if item.Name != "subscription_payment_confirmation_link" {
		t.Fatalf("expected subscription payment confirmation link migration 0025, got %+v", item)
	}
	for _, want := range []string{
		"ALTER TABLE payment_transactions",
		"subscription_invoice_id",
		"payment_transactions_subscription_invoice_idx",
		"payment_transactions_subscription_invoice_id_fkey",
	} {
		if !strings.Contains(item.SQL, want) {
			t.Fatalf("expected subscription payment confirmation link migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesPlatformAuthSessionsNullableTenant(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var item migration
	for _, current := range migrations {
		if current.Version == "0026" {
			item = current
			break
		}
	}
	if item.Name != "platform_auth_sessions_nullable_tenant" {
		t.Fatalf("expected platform auth sessions nullable tenant migration 0026, got %+v", item)
	}
	for _, want := range []string{
		"ALTER TABLE app_auth_sessions",
		"ALTER COLUMN tenant_id DROP NOT NULL",
	} {
		if !strings.Contains(item.SQL, want) {
			t.Fatalf("expected platform auth sessions nullable tenant migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantDataIsolation(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var tenantData migration
	for _, item := range migrations {
		if item.Version == "0016" {
			tenantData = item
			break
		}
	}
	if tenantData.Name != "tenant_data_isolation" {
		t.Fatalf("expected tenant data isolation migration 0016, got %+v", tenantData)
	}

	for _, want := range []string{
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS tenant_id",
		"ALTER TABLE parents ADD COLUMN IF NOT EXISTS tenant_id",
		"ALTER TABLE notification_campaigns ADD COLUMN IF NOT EXISTS tenant_id",
		"ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS tenant_id",
		"ALTER TABLE operation_logs ADD COLUMN IF NOT EXISTS tenant_id",
		"ALTER TABLE audit_logs DISABLE TRIGGER audit_logs_prevent_update",
		"ALTER TABLE audit_logs ENABLE TRIGGER audit_logs_prevent_update",
		"students_tenant_code_key",
		"parents_tenant_email_key",
		"notification_campaigns_tenant_code_key",
		"operation_logs_tenant_id_fkey",
		"DEKISUGI",
	} {
		if !strings.Contains(tenantData.SQL, want) {
			t.Fatalf("expected tenant data isolation migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantOnboardingAndSwitching(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var tenantOnboarding migration
	for _, item := range migrations {
		if item.Version == "0017" {
			tenantOnboarding = item
			break
		}
	}
	if tenantOnboarding.Name != "tenant_onboarding_and_switching" {
		t.Fatalf("expected tenant onboarding migration 0017, got %+v", tenantOnboarding)
	}

	for _, want := range []string{
		"tenant.view",
		"tenant.create",
		"tenant.update",
		"tenant.switch",
		"admin",
		"staff",
		"accountant",
		"tenant_memberships_user_status_idx",
	} {
		if !strings.Contains(tenantOnboarding.SQL, want) {
			t.Fatalf("expected tenant onboarding migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantScopedPaymentProvidersAndWebhooks(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "0018" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "tenant_scoped_payment_providers_and_webhooks" {
		t.Fatalf("expected tenant-scoped payment providers migration 0018, got %+v", migrationItem)
	}

	for _, want := range []string{
		"payment_providers ADD COLUMN IF NOT EXISTS tenant_id",
		"CREATE UNIQUE INDEX IF NOT EXISTS payment_providers_tenant_code_key",
		"CREATE INDEX IF NOT EXISTS payment_providers_tenant_id_idx",
		"payment_providers_tenant_id_fkey",
		"INSERT INTO payment_providers (tenant_id, code, display_name, provider_type, status, config)",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected tenant-scoped payment providers migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesSubscriptionHardeningCrossTenantOperations(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "0019" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "subscription_hardening_cross_tenant_operations" {
		t.Fatalf("expected subscription hardening migration 0019, got %+v", migrationItem)
	}

	for _, want := range []string{
		"operation_log.cross_tenant_view",
		"audit_log.cross_tenant_view",
		"operations.cross_tenant.read",
		"audit.cross_tenant.read",
		"admin",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected subscription hardening migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantSubscriptionBilling(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "0020" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "tenant_subscription_billing" {
		t.Fatalf("expected tenant subscription billing migration 0020, got %+v", migrationItem)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS subscription_plans",
		"CREATE TABLE IF NOT EXISTS tenant_subscriptions",
		"free_trial",
		"standard",
		"subscription.view",
		"subscription.update",
		"tenant_subscriptions_tenant_id_key",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected tenant subscription billing migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesSetUpdatedAtAliasCompatibility(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "00195" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "set_updated_at_alias" {
		t.Fatalf("expected compatibility migration 00195, got %+v", migrationItem)
	}

	for _, want := range []string{
		"CREATE OR REPLACE FUNCTION set_updated_at()",
		"RETURN abc_set_updated_at();",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected compatibility migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesSetUpdatedAtAliasFix(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "00196" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "fix_set_updated_at_alias" {
		t.Fatalf("expected compatibility fix migration 00196, got %+v", migrationItem)
	}

	for _, want := range []string{
		"CREATE OR REPLACE FUNCTION set_updated_at()",
		"NEW.updated_at = now();",
		"RETURN NEW;",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected compatibility fix migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesTenantUsageEntitlements(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "0021" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "tenant_usage_entitlements" {
		t.Fatalf("expected tenant usage entitlements migration 0021, got %+v", migrationItem)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS tenant_usage_counters",
		"'schools'",
		"'operators'",
		"'students'",
		"'monthly_notifications'",
		`"students":200`,
		`"monthly_notifications":10000`,
		"tenant_usage_counters_metric_idx",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected tenant usage entitlements migration to contain %q", want)
		}
	}
}

func TestLoadEmbeddedMigrationsIncludesSubscriptionInvoicingAndDunning(t *testing.T) {
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}

	var migrationItem migration
	for _, item := range migrations {
		if item.Version == "0022" {
			migrationItem = item
			break
		}
	}
	if migrationItem.Name != "subscription_invoicing_and_dunning" {
		t.Fatalf("expected subscription invoicing migration 0022, got %+v", migrationItem)
	}

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS subscription_invoices",
		"CREATE TABLE IF NOT EXISTS subscription_invoice_status_history",
		"CREATE TABLE IF NOT EXISTS subscription_dunning_runs",
		"subscription_invoices_tenant_invoice_code_key",
		"subscription_invoices_subscription_period_key",
		"'past_due'",
		"'paid'",
	} {
		if !strings.Contains(migrationItem.SQL, want) {
			t.Fatalf("expected subscription invoicing migration to contain %q", want)
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

func TestDatabaseMigrationsReadyRejectsPendingAndDriftedMigrations(t *testing.T) {
	migrations := []migration{
		{Version: "0001", Name: "foundation", FileName: "0001_foundation.sql", Checksum: "checksum-1"},
		{Version: "0002", Name: "students", FileName: "0002_students.sql", Checksum: "checksum-2"},
	}

	t.Run("ready", func(t *testing.T) {
		db, _ := openFakeMigrationDB(t, "ready", map[string]appliedMigration{
			"0001": {Version: "0001", Name: "foundation", Checksum: "checksum-1", AppliedAt: time.Now()},
			"0002": {Version: "0002", Name: "students", Checksum: "checksum-2", AppliedAt: time.Now()},
		})
		if err := databaseMigrationsReady(context.Background(), db, defaultMigrationTable, migrations); err != nil {
			t.Fatalf("expected database readiness, got %v", err)
		}
	})

	t.Run("pending", func(t *testing.T) {
		db, _ := openFakeMigrationDB(t, "pending", map[string]appliedMigration{
			"0001": {Version: "0001", Name: "foundation", Checksum: "checksum-1", AppliedAt: time.Now()},
		})
		err := databaseMigrationsReady(context.Background(), db, defaultMigrationTable, migrations)
		if err == nil || !strings.Contains(err.Error(), "pending: 0002") {
			t.Fatalf("expected pending migration error, got %v", err)
		}
	})

	t.Run("drifted", func(t *testing.T) {
		db, _ := openFakeMigrationDB(t, "drifted", map[string]appliedMigration{
			"0001": {Version: "0001", Name: "foundation", Checksum: "old-checksum", AppliedAt: time.Now()},
			"0002": {Version: "0002", Name: "students", Checksum: "checksum-2", AppliedAt: time.Now()},
		})
		err := databaseMigrationsReady(context.Background(), db, defaultMigrationTable, migrations)
		if err == nil || !strings.Contains(err.Error(), "drifted: 0001") {
			t.Fatalf("expected drifted migration error, got %v", err)
		}
	})
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
