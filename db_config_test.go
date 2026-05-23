package main

import (
	"strings"
	"testing"
	"time"
)

func TestLoadDatabaseConfigUsesEnvironmentSpecificURL(t *testing.T) {
	cfg, err := loadDatabaseConfigFromLookup(mapLookup(map[string]string{
		"ABC_ENV":                  "staging",
		"ABC_DATABASE_URL":         "postgres://fallback",
		"ABC_DATABASE_URL_STAGING": "postgres://staging",
		"ABC_DB_MAX_OPEN_CONNS":    "12",
		"ABC_DB_MAX_IDLE_CONNS":    "4",
		"ABC_DB_CONN_MAX_LIFETIME": "45m",
	}))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Environment != appEnvStaging {
		t.Fatalf("expected staging environment, got %q", cfg.Environment)
	}
	if cfg.URL != "postgres://staging" || cfg.URLSource != "ABC_DATABASE_URL_STAGING" {
		t.Fatalf("expected staging URL source, got %q from %q", cfg.URL, cfg.URLSource)
	}
	if cfg.MaxOpenConns != 12 || cfg.MaxIdleConns != 4 || cfg.ConnMaxLifetime != 45*time.Minute {
		t.Fatalf("unexpected pool config: %+v", cfg)
	}
}

func TestDatabaseConfigRequiresURLForDatabaseCommands(t *testing.T) {
	cfg, err := loadDatabaseConfigFromLookup(mapLookup(map[string]string{
		"ABC_ENV": "production",
	}))
	if err != nil {
		t.Fatal(err)
	}

	err = cfg.requireURL()
	if err == nil {
		t.Fatal("expected missing database URL error")
	}
	if !strings.Contains(err.Error(), "ABC_DATABASE_URL_PRODUCTION") {
		t.Fatalf("expected production URL variable in error, got %v", err)
	}
}

func TestLoadDatabaseConfigRejectsUnsafeMigrationTable(t *testing.T) {
	_, err := loadDatabaseConfigFromLookup(mapLookup(map[string]string{
		"ABC_DB_MIGRATIONS_TABLE": "schema_migrations; drop table app_users",
	}))
	if err == nil {
		t.Fatal("expected invalid migration table error")
	}
}

func mapLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
