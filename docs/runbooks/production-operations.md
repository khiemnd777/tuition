# Production Operations Runbook

This runbook covers deployment readiness, staging smoke tests, rollback, and incident response for the ABC SUN production modules.

## Pre-Deploy Checklist

Run these checks before deploying a build that changes production modules:

```sh
go test ./...
go run . migrate status
```

Verify environment configuration without printing secrets:

```sh
ABC_ENV=staging go run . db config
ABC_ENV=staging go run . db ping
```

Before production schema changes, take and verify a fresh backup using `docs/runbooks/backup-restore.md`.

## Deployment Sequence

1. Confirm the target environment: `ABC_ENV=staging` or `ABC_ENV=production`.
2. Confirm database URL comes from the environment or secret manager, not a tracked file.
3. Run `go run . migrate status`.
4. Run `go run . migrate up`.
5. Start the API and admin UI.
6. Run the staging smoke tests below before routing real production traffic.
7. Record deploy timestamp, operator, git revision, migration status, and backup file in the deployment notes.

## Staging Smoke Tests

Use staging data only. Do not send real parent email during smoke tests.

1. Open the admin UI and load `Dashboard`.
2. Open `Báo cáo`, apply a year/class/period filter, and export class, invoice, and transaction CSV.
3. Open `Vận hành` and confirm operation/audit log endpoints load.
4. Preview a fee schedule with at least one adjustment and save it with an operator name.
5. Generate or open an unpaid invoice, then create a manual cash receipt in staging with a clear reason.
6. Confirm invoice status and `Đối soát` ledger update.
7. Run notification campaign preview or dry-run only.
8. Check `go run . migrate status` remains clean after the smoke test.

## Rollback

The migration runner is forward-only. Rollback means reverting the app binary/container and restoring data only when required.

1. Stop the newly deployed app.
2. Start the previous known-good app version.
3. If the new migration is backward-compatible, keep the migrated database and monitor.
4. If data corruption occurred, stop writers and follow the restore incident procedure in `backup-restore.md`.
5. Preserve failed webhook payloads, notification logs, operation logs, and audit logs for review.

## Incident Response

Use `Vận hành` first:

- `operation_logs`: webhook parsing/signature/reconciliation errors, notification send failures, cron/background job failures.
- `audit_logs`: manual cash receipts and fee adjustment changes with actor/reason metadata.

Incident steps:

1. Classify impact: invoices, payments, notifications, exports, or admin access.
2. Pause risky writes if money movement may be affected.
3. Export current reports from `Báo cáo` for the affected filter.
4. Review operation logs by source and level.
5. Review audit logs for the affected invoice, manual receipt, fee schedule, or adjustment.
6. Reconcile invoice totals against payment transactions before marking the incident resolved.
7. Record root cause, fix, validation, and follow-up owner.

## Production Readiness Checklist

- Database migrations apply cleanly from an empty database.
- Backup and restore drill has been performed for the target environment.
- Report exports open in spreadsheet software and preserve amounts as integers.
- Manual cash receipts require collector, receipt reference, amount, and reason.
- Fee adjustments require reason and operator name when saved through the UI.
- Notification sends are dry-run validated before any real send.
- Webhook secrets and provider credentials are stored outside git.
- `Vận hành` can load audit and operation logs.
- Staging smoke tests passed for the same build before production traffic.
