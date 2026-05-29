# ABC SUN Initiative State

Last updated: 2026-05-29

## Current Status

Production roadmap implementation has student, parent, class master data, fee schedule setup, invoice/PDF receipt output, payment/reconciliation ledger, and notification campaigns complete.

Current phase: `initiative_6_complete`

Current initiative: none

Next recommended initiative: `Initiative 7: Web Admin`

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
- Initiative 3: Fee Types And Fee Schedules is complete:
  - Added production schema for `fee_types`, `fee_schedules`, `fee_schedule_items`, and `student_fee_adjustments`.
  - Seeded default fee types for tuition, lunch, shuttle, uniform, insurance, materials, previous fees, and custom fees with Vietnamese and English labels.
  - Added API endpoints for fee schedule options, list, preview, and save.
  - Added deterministic preview totals for class/grade/school-year scopes and student adjustments.
  - Added adjustment support for discount, surcharge, waiver, and carry-over with required reasons.
  - Added the `Bảng phí` UI workflow for period fee setup, preview before invoice generation, and saved schedule list.
  - Preserved the legacy payment-row fee template and the `PaymentItems` total-overrides-amount invariant.
  - Added tests, README docs, and glossary terms.
- Initiative 4: Invoice And PDF Receipt is complete:
  - Added production schema for `invoices`, `invoice_items`, `invoice_adjustments`, `invoice_status_history`, and `receipt_documents`.
  - Added invoice APIs for options, list, preview, idempotent generation, QR payment data, and PDF receipt output.
  - Generated stable invoice codes that map directly to VietQR `BillNumber`.
  - Snapshot invoice line items and adjustments from saved fee schedules so issued invoices do not change when schedules are edited later.
  - Added PDF receipt rendering from invoice data with school, student, class, period, line items, total, status, issue timestamp, and VietQR QR.
  - Added the `Hóa đơn` UI workflow with preview, generation, invoice list, QR preview, and PDF links.
  - Kept the legacy `Thanh toán` tab before `Email & Cron`.
  - Added tests, README docs, and glossary terms.
- Initiative 5: Payment And Reconciliation is complete:
  - Added production schema for `payment_providers`, `payment_intents`, `provider_events`, `payment_transactions`, `reconciliation_matches`, and `manual_cash_receipts`.
  - Seeded baseline providers for `manual_vietqr`, `sepay`, and `payos`.
  - Added payment APIs for providers, payment intent creation, transaction listing, reconciliation dashboard data, provider webhooks, and manual cash receipts.
  - Implemented idempotent webhook/event and transaction storage so retries do not duplicate payment ledger entries.
  - Added SePay webhook normalization, payOS payment link request/signature support, payOS webhook signature verification when checksum key is configured, and manual VietQR intents.
  - Reconciled incoming transactions to invoices by invoice code/provider reference, collection account, and amount; invoice `paid_amount` and status now derive from active reconciliation matches.
  - Added the `Đối soát` UI workflow with provider filters, receivable summary, invoice ledger actions, transaction table, QR/payment intent detail, and manual cash receipt entry.
  - Added tests, README docs, env examples, and glossary terms.
- Initiative 6: Notification Campaigns is complete:
  - Added production schema for `notification_templates`, `notification_campaigns`, `notification_recipients`, and `notification_logs`.
  - Seeded versioned baseline templates for first notice and reminder campaigns.
  - Added notification APIs for options/templates, campaign listing, dry-run preview, save, send, and delivery logs.
  - Targeted invoice recipients through active parent billing contacts instead of temporary payment rows.
  - Enforced reminder targeting so paid invoices are not accidentally included.
  - Reused the existing email renderer/provider path so notification sends keep inline VietQR CID images and Gmail/Resend behavior.
  - Added idempotent send behavior per campaign/template/invoice/recipient unless explicitly re-sent.
  - Added the `Thông báo` UI workflow with filters, recipient dry-run, saved campaign list, real-send confirmation, and delivery logs.
  - Added tests, README docs, and glossary terms.

## Not Started

- Dashboard, reporting, role administration, and production operations screens have not been implemented.
- Reports and operations modules have not started.

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
Start Initiative 7 from docs/initiatives/production-module-roadmap.md. Build production Web Admin workflows for dashboard, reports, roles, and any remaining admin screens. Preserve current VietQR/payment reconciliation and notification campaign behavior, and update docs/initiatives/current-state.md when finished.
```

## Known Safety Constraints

- Do not send real email unless explicitly requested.
- Do not run real cron batches unless explicitly requested.
- Do not commit or print real secrets from local config files.
- Preserve `PaymentItems` total overriding raw `Amount`.
- Preserve VietQR TLV/CRC behavior with exact tests when touching QR generation.
