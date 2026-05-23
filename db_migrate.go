package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"time"
)

//go:embed migrations/*.sql
var embeddedMigrationFiles embed.FS

var migrationFilePattern = regexp.MustCompile(`^([0-9]{4,})_([a-z0-9_]+)\.sql$`)

type migration struct {
	Version  string
	Name     string
	FileName string
	SQL      string
	Checksum string
}

type appliedMigration struct {
	Version   string
	Name      string
	Checksum  string
	AppliedAt time.Time
}

type migrationStatus struct {
	Migration migration
	Applied   bool
	Drifted   bool
	AppliedAt time.Time
}

func loadEmbeddedMigrations() ([]migration, error) {
	entries, err := fs.ReadDir(embeddedMigrationFiles, "migrations")
	if err != nil {
		return nil, err
	}

	migrations := make([]migration, 0, len(entries))
	seenVersions := map[string]string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		matches := migrationFilePattern.FindStringSubmatch(fileName)
		if matches == nil {
			return nil, fmt.Errorf("invalid migration filename %q", fileName)
		}

		content, err := fs.ReadFile(embeddedMigrationFiles, path.Join("migrations", fileName))
		if err != nil {
			return nil, err
		}
		if len(content) == 0 {
			return nil, fmt.Errorf("migration %q is empty", fileName)
		}

		version := matches[1]
		if previous := seenVersions[version]; previous != "" {
			return nil, fmt.Errorf("duplicate migration version %s in %s and %s", version, previous, fileName)
		}
		seenVersions[version] = fileName

		migrations = append(migrations, migration{
			Version:  version,
			Name:     matches[2],
			FileName: fileName,
			SQL:      string(content),
			Checksum: migrationChecksum(content),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	if len(migrations) == 0 {
		return nil, fmt.Errorf("no database migrations found")
	}
	return migrations, nil
}

func runMigrations(ctx context.Context, db *sql.DB, table string, migrations []migration, out io.Writer) error {
	if err := validateMigrationTableName(table); err != nil {
		return err
	}
	if len(migrations) == 0 {
		return fmt.Errorf("no database migrations to apply")
	}
	if err := ensureMigrationTable(ctx, db, table); err != nil {
		return err
	}

	applied, err := loadAppliedMigrations(ctx, db, table)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if existing, ok := applied[migration.Version]; ok {
			if existing.Checksum != migration.Checksum {
				return fmt.Errorf("migration %s was already applied with a different checksum", migration.Version)
			}
			if out != nil {
				fmt.Fprintf(out, "skip %s %s\n", migration.Version, migration.Name)
			}
			continue
		}

		if err := applyMigration(ctx, db, table, migration); err != nil {
			return err
		}
		if out != nil {
			fmt.Fprintf(out, "applied %s %s\n", migration.Version, migration.Name)
		}
	}
	return nil
}

func ensureMigrationTable(ctx context.Context, db *sql.DB, table string) error {
	if err := validateMigrationTableName(table); err != nil {
		return err
	}
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	version text PRIMARY KEY,
	name text NOT NULL,
	checksum text NOT NULL,
	applied_at timestamptz NOT NULL DEFAULT now()
)`, table)
	_, err := db.ExecContext(ctx, query)
	return err
}

func loadAppliedMigrations(ctx context.Context, db *sql.DB, table string) (map[string]appliedMigration, error) {
	if err := validateMigrationTableName(table); err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT version, name, checksum, applied_at FROM %s ORDER BY version`, table)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := map[string]appliedMigration{}
	for rows.Next() {
		var item appliedMigration
		if err := rows.Scan(&item.Version, &item.Name, &item.Checksum, &item.AppliedAt); err != nil {
			return nil, err
		}
		applied[item.Version] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return applied, nil
}

func applyMigration(ctx context.Context, db *sql.DB, table string, migration migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("apply migration %s: %w", migration.FileName, err)
	}

	query := fmt.Sprintf(`INSERT INTO %s (version, name, checksum) VALUES ($1, $2, $3)`, table)
	if _, err := tx.ExecContext(ctx, query, migration.Version, migration.Name, migration.Checksum); err != nil {
		return fmt.Errorf("record migration %s: %w", migration.FileName, err)
	}

	return tx.Commit()
}

func buildMigrationStatuses(migrations []migration, applied map[string]appliedMigration) []migrationStatus {
	statuses := make([]migrationStatus, 0, len(migrations))
	for _, migration := range migrations {
		existing, ok := applied[migration.Version]
		statuses = append(statuses, migrationStatus{
			Migration: migration,
			Applied:   ok,
			Drifted:   ok && existing.Checksum != migration.Checksum,
			AppliedAt: existing.AppliedAt,
		})
	}
	return statuses
}

func validateMigrationTableName(table string) error {
	if !sqlIdentifierPattern.MatchString(table) {
		return fmt.Errorf("migration table %q must be a simple SQL identifier", table)
	}
	return nil
}

func migrationChecksum(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}
