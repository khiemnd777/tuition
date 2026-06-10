package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

const databaseCommandTimeout = 2 * time.Minute

func runDBCommand(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	switch args[0] {
	case "db":
		return true, runDatabaseUtilityCommand(ctx, args[1:], stdout)
	case "migrate":
		return true, runMigrationCommand(ctx, args[1:], stdout)
	case "demo":
		return true, runDemoCommand(ctx, args[1:], stdout)
	default:
		return false, nil
	}
}

func runDatabaseUtilityCommand(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: go run . db [config|ping]")
	}

	cfg, err := loadDatabaseConfig()
	if err != nil {
		return err
	}

	switch args[0] {
	case "config":
		printDatabaseConfig(stdout, cfg)
		return nil
	case "ping":
		db, err := openConfiguredDatabase(ctx, cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		fmt.Fprintf(stdout, "database ping ok for %s\n", cfg.Environment)
		return nil
	default:
		return fmt.Errorf("usage: go run . db [config|ping]")
	}
}

func runMigrationCommand(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: go run . migrate [up|status]")
	}

	cfg, err := loadDatabaseConfig()
	if err != nil {
		return err
	}
	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		return err
	}

	commandCtx, cancel := context.WithTimeout(ctx, databaseCommandTimeout)
	defer cancel()

	db, err := openConfiguredDatabase(commandCtx, cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	switch args[0] {
	case "up":
		if err := runMigrations(commandCtx, db, cfg.MigrationsTable, migrations, stdout); err != nil {
			return err
		}
		return runDefaultDemoDataMigrations(commandCtx, db, stdout)
	case "status":
		if err := ensureMigrationTable(commandCtx, db, cfg.MigrationsTable); err != nil {
			return err
		}
		applied, err := loadAppliedMigrations(commandCtx, db, cfg.MigrationsTable)
		if err != nil {
			return err
		}
		printMigrationStatuses(stdout, buildMigrationStatuses(migrations, applied))
		return nil
	default:
		return fmt.Errorf("usage: go run . migrate [up|status]")
	}
}

func printDatabaseConfig(out io.Writer, cfg databaseConfig) {
	urlState := "not configured"
	if cfg.URLSource != "" {
		urlState = "configured via " + cfg.URLSource
	}

	fmt.Fprintf(out, "environment: %s\n", cfg.Environment)
	fmt.Fprintf(out, "database_url: %s\n", urlState)
	fmt.Fprintf(out, "migrations_table: %s\n", cfg.MigrationsTable)
	fmt.Fprintf(out, "max_open_conns: %d\n", cfg.MaxOpenConns)
	fmt.Fprintf(out, "max_idle_conns: %d\n", cfg.MaxIdleConns)
	fmt.Fprintf(out, "conn_max_lifetime: %s\n", cfg.ConnMaxLifetime)
}

func printMigrationStatuses(out io.Writer, statuses []migrationStatus) {
	fmt.Fprintln(out, "version  name                         status   applied_at")
	for _, status := range statuses {
		state := "pending"
		appliedAt := "-"
		if status.Applied {
			state = "applied"
			appliedAt = status.AppliedAt.Format(time.RFC3339)
		}
		if status.Drifted {
			state = "drifted"
		}
		fmt.Fprintf(out, "%-8s %-28s %-8s %s\n", status.Migration.Version, strings.ReplaceAll(status.Migration.Name, "_", "-"), state, appliedAt)
	}
}
