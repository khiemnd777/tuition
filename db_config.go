package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	appEnvLocal      = "local"
	appEnvStaging    = "staging"
	appEnvProduction = "production"

	defaultMigrationTable   = "schema_migrations"
	defaultMaxOpenConns     = 5
	defaultMaxIdleConns     = 2
	defaultConnMaxLifetime  = 30 * time.Minute
	defaultDatabasePingTime = 5 * time.Second
)

var sqlIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type databaseConfig struct {
	Environment     string
	URL             string
	URLSource       string
	MigrationsTable string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func loadDatabaseConfig() (databaseConfig, error) {
	return loadDatabaseConfigFromLookup(os.LookupEnv)
}

func loadDatabaseConfigFromLookup(lookup func(string) (string, bool)) (databaseConfig, error) {
	env, err := loadAppEnvironment(lookup)
	if err != nil {
		return databaseConfig{}, err
	}

	url, source := firstEnvValue(lookup, databaseURLCandidates(env)...)
	table := envString(lookup, "DEKISUGI_DB_MIGRATIONS_TABLE", defaultMigrationTable)
	if !sqlIdentifierPattern.MatchString(table) {
		return databaseConfig{}, fmt.Errorf("DEKISUGI_DB_MIGRATIONS_TABLE must be a simple SQL identifier")
	}

	return databaseConfig{
		Environment:     env,
		URL:             url,
		URLSource:       source,
		MigrationsTable: table,
		MaxOpenConns:    envPositiveInt(lookup, "DEKISUGI_DB_MAX_OPEN_CONNS", defaultMaxOpenConns),
		MaxIdleConns:    envPositiveInt(lookup, "DEKISUGI_DB_MAX_IDLE_CONNS", defaultMaxIdleConns),
		ConnMaxLifetime: envDuration(lookup, "DEKISUGI_DB_CONN_MAX_LIFETIME", defaultConnMaxLifetime),
	}, nil
}

func loadAppEnvironment(lookup func(string) (string, bool)) (string, error) {
	raw, _ := firstEnvValue(lookup, "DEKISUGI_ENV", "APP_ENV")
	if raw == "" {
		return appEnvLocal, nil
	}
	return normalizeAppEnvironment(raw)
}

func normalizeAppEnvironment(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "local", "dev", "development":
		return appEnvLocal, nil
	case "stage", "staging":
		return appEnvStaging, nil
	case "prod", "production":
		return appEnvProduction, nil
	default:
		return "", fmt.Errorf("unsupported environment %q; use local, staging, or production", raw)
	}
}

func databaseURLCandidates(env string) []string {
	return []string{
		"DEKISUGI_DATABASE_URL_" + strings.ToUpper(env),
		"DEKISUGI_DATABASE_URL",
		"DATABASE_URL",
	}
}

func (cfg databaseConfig) requireURL() error {
	if strings.TrimSpace(cfg.URL) != "" {
		return nil
	}
	return fmt.Errorf("database URL is required; set %s, DEKISUGI_DATABASE_URL, or DATABASE_URL", databaseURLCandidates(cfg.Environment)[0])
}

func openConfiguredDatabase(ctx context.Context, cfg databaseConfig) (*sql.DB, error) {
	if err := cfg.requireURL(); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, defaultDatabasePingTime)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func firstEnvValue(lookup func(string) (string, bool), keys ...string) (string, string) {
	for _, key := range keys {
		if value, ok := lookup(key); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value), key
		}
	}
	return "", ""
}

func envString(lookup func(string) (string, bool), key string, fallback string) string {
	if value, ok := lookup(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func envPositiveInt(lookup func(string) (string, bool), key string, fallback int) int {
	value, ok := lookup(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func envDuration(lookup func(string) (string, bool), key string, fallback time.Duration) time.Duration {
	value, ok := lookup(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
