# PostgreSQL Backup And Restore Runbook

This runbook covers the ABC SUN production persistence foundation. Keep database URLs and backup files outside git.

## Environment

Set the active environment before running DB commands:

```sh
export ABC_ENV=production
export ABC_DATABASE_URL_PRODUCTION='postgres://user:password@host:5432/abcsun?sslmode=require'
```

Use `ABC_ENV=local` with `ABC_DATABASE_URL_LOCAL` for local development and `ABC_ENV=staging` with `ABC_DATABASE_URL_STAGING` for staging.

Check the app sees the intended environment without printing secrets:

```sh
go run . db config
go run . db ping
```

## Migrations

Apply migrations before starting a new environment and after every deploy that includes schema changes:

```sh
go run . migrate status
go run . migrate up
go run . migrate status
```

The migration runner records applied versions and checksums in `schema_migrations`. Re-running `migrate up` skips already-applied migrations.

## Backup

Create a custom-format backup with PostgreSQL `pg_dump`:

```sh
mkdir -p backups
pg_dump "$ABC_DATABASE_URL_PRODUCTION" \
  --format=custom \
  --no-owner \
  --no-acl \
  --file="backups/abcsun-$(date +%Y%m%d-%H%M%S).dump"
```

Before risky production changes, take a fresh backup and verify that the file exists and is non-empty:

```sh
ls -lh backups/*.dump
pg_restore --list backups/abcsun-YYYYMMDD-HHMMSS.dump >/tmp/abcsun-backup-list.txt
test -s /tmp/abcsun-backup-list.txt
```

## Restore Drill

Restore into a new empty database, not over the live database:

```sh
createdb abcsun_restore_check
pg_restore \
  --dbname='postgres://user:password@host:5432/abcsun_restore_check?sslmode=require' \
  --clean \
  --if-exists \
  --no-owner \
  --no-acl \
  backups/abcsun-YYYYMMDD-HHMMSS.dump
```

Point a local shell at the restored database and verify migrations:

```sh
export ABC_ENV=local
export ABC_DATABASE_URL_LOCAL='postgres://user:password@host:5432/abcsun_restore_check?sslmode=require'
go run . migrate status
```

Also verify that core accounting tables are readable in the restored database:

```sh
psql "$ABC_DATABASE_URL_LOCAL" -c "select count(*) from invoices;"
psql "$ABC_DATABASE_URL_LOCAL" -c "select count(*) from payment_transactions;"
psql "$ABC_DATABASE_URL_LOCAL" -c "select count(*) from audit_logs;"
psql "$ABC_DATABASE_URL_LOCAL" -c "select count(*) from operation_logs;"
```

## Restore Incident Procedure

1. Stop app processes that write to the damaged database.
2. Preserve the damaged database for investigation by renaming it or taking a final dump if possible.
3. Create a new empty replacement database.
4. Restore the selected backup with `pg_restore`.
5. Run `go run . migrate up` against the replacement database.
6. Run `go run . db ping` and the staging smoke checks before pointing traffic back.
7. Record the backup file, restore timestamp, operator, and incident reason in the operations log.

Never restore production from an unverified backup file.
