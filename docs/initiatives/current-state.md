# ABC SUN Initiative State

Last updated: 2026-05-23

## Current Status

Production roadmap implementation has student, parent, and class master data complete.

Current phase: `initiative_2_complete`

Current initiative: none

Next recommended initiative: `Initiative 3: Fee Types And Fee Schedules`

Roadmap source: `docs/initiatives/production-module-roadmap.md`

## Completed

- Production module analysis was captured in `docs/initiatives/production-module-roadmap.md`.
- The roadmap is split into independent initiatives:
  - Initiative 1: Foundation And Persistence
  - Initiative 2: Student, Parent, Class Master Data
  - Initiative 3: Fee Types And Fee Schedules
  - Initiative 4: Invoice And PDF Receipt
  - Initiative 5: Payment And Reconciliation
  - Initiative 6: Notification Campaigns
  - Initiative 7: Web Admin
  - Initiative 8: Reports, Audit, And Operations
- Payment gateway direction was recorded:
  - Prefer VietQR plus bank/webhook reconciliation for tuition.
  - Keep provider adapters flexible.
  - Consider payOS or SePay as practical production providers.
  - Avoid percentage-fee hosted gateways as the primary low-cost tuition channel.
- README now links to the production roadmap.
- Initiative 1: Foundation And Persistence is complete:
  - Added PostgreSQL configuration through environment variables for local, staging, and production.
  - Added `go run . db config`, `go run . db ping`, `go run . migrate status`, and `go run . migrate up`.
  - Added embedded SQL migrations under `migrations/`.
  - Added foundation schema for `app_users`, `app_roles`, `app_permissions`, role assignment tables, and immutable `audit_logs`.
  - Defined UUID primary key and timestamp/actor audit conventions in `docs/adr/0002-production-persistence-foundation.md`.
  - Added PostgreSQL backup/restore runbook in `docs/runbooks/backup-restore.md`.
- Added local Docker stack for `api`, `admin`, `postgres`, and `redis`:
  - API builds the Go app image and keeps local runtime JSON state in `api_data`.
  - Admin runs Vite dev server in local Docker with source bind mount, HMR, API proxying, and `admin_node_modules` volume.
  - Admin production image still builds static assets with `npm ci`/`npm run build` and serves them through nginx.
  - API default HTTP port is `18080`; port `8080` is intentionally avoided.
  - PostgreSQL and Redis use named local volumes and health checks.
  - Migrations can run through `docker compose run --rm api migrate up`.
- Initiative 2: Student, Parent, Class Master Data is complete:
  - Added production schema for `school_years`, `classes`, `students`, `parents`, and `student_parents`.
  - Made `student_code` the mandatory durable student identifier; duplicate student names are supported.
  - Added parent contact delivery flags for primary, active, and billing email behavior.
  - Added CSV master data import preview/apply with conflict reporting and no silent overwrite of mismatched records.
  - Added API endpoints for master-data options, student listing by school year/class/grade/search, and CSV import.
  - Added the `Học sinh` UI tab with filters, import preview/apply controls, student table, and conflict report.
  - Added `samples/master_data.csv`, tests, README docs, and glossary terms.

## Not Started

- No payment provider integration has been implemented.
- Web Admin production screens beyond the student master-data tab have not been implemented.
- Fee schedules, invoices, reconciliation, notification campaigns, reports, and operations modules have not started.

## Agent Protocol

When the user says:

- `Đến đâu rồi?`
- `bắt đầu tiếp công việc`
- `tiếp tục`
- `continue`
- `status`

The agent must:

1. Read this file first.
2. Read `docs/initiatives/production-module-roadmap.md`.
3. Briefly tell the user:
   - current phase
   - completed work
   - next recommended initiative
   - any known blockers
4. If the user asked to continue, start the next recommended initiative unless the user specifies a different initiative.
5. Use repo skills before implementation:
   - `$abcsun-change-workflow` by default
   - `$abcsun-vietqr-payments` for QR/payment/invoice amount behavior
   - `$abcsun-email-delivery` for email/cron/notification behavior
   - `$abcsun-frontend-ui` for Web Admin UI work
6. Update this file at the end of every initiative or meaningful partial implementation.

## Next Launch Prompt

Use this when the user says to continue without specifying a module:

```text
Start Initiative 3 from docs/initiatives/production-module-roadmap.md. Build production fee types and fee schedules for class/period/month tuition setup. Preserve current VietQR/email behavior and update docs/initiatives/current-state.md when finished.
```

## Known Safety Constraints

- Do not send real email unless explicitly requested.
- Do not run real cron batches unless explicitly requested.
- Do not commit or print real secrets from local config files.
- Preserve `PaymentItems` total overriding raw `Amount`.
- Preserve VietQR TLV/CRC behavior with exact tests when touching QR generation.
