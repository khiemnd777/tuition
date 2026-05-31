# ABC SUN Initiative State

Last updated: 2026-05-31

## Current Status

Production roadmap implementation has student, parent, class master data, fee schedule setup, invoice/PDF receipt output, payment/reconciliation ledger, notification campaigns, production Web Admin screens, reports/export, audit review, and operational readiness complete.

Advanced Production planning has started so the project can move into production hardening and usability work one initiative at a time.

Current phase: `advanced_4_complete`

Current initiative: Advanced 4 - School Tree Management is complete.

Next recommended initiative: define the next Advanced Production initiative before implementation.

Roadmap source: `docs/initiatives/production-module-roadmap.md` for completed production modules; Advanced Production roadmap is currently recorded in this file.

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
- Initiative 7: Web Admin is complete:
  - Added admin dashboard APIs for receivable, collected, outstanding, collection rate, unpaid/partial/review counts, unmatched transactions, top classes by outstanding amount, and attention invoices.
  - Added report APIs for class-level summaries and invoice-level detail with filters for school year, grade, class, period, month, and status.
  - Added user, role, and permission APIs over `app_users`, `app_roles`, `app_permissions`, and `app_user_roles`.
  - Added a `0007_web_admin` migration to seed dashboard/report and role-assignment permissions without changing prior migration checksums.
  - Added the `Dashboard`, `Báo cáo`, and `Người dùng` UI workflows while preserving existing `Học sinh`, `Bảng phí`, `Hóa đơn`, `Đối soát`, and `Thông báo` workflows.
  - Added focused tests for dashboard/report aggregation, attention invoice ordering, user input normalization, and migration permission seeds.
  - Updated README docs and glossary terms.
- Initiative 8: Reports, Audit, And Operations is complete:
  - Added report CSV exports for class summaries, invoice detail, and payment transactions using the existing admin report filters.
  - Added immutable audit capture for fee schedule saves, student fee adjustments, and manual cash receipts with actor/reason metadata.
  - Added `operation_logs` and the `Vận hành` UI for webhook, notification email, and email cron/background job failure review.
  - Required manual cash receipt reasons and UI operator names for saved fee adjustments.
  - Added operation/audit APIs, report export permissions, tests, README/API docs, glossary terms, backup verification notes, and production operations runbook.

## Advanced Production Roadmap

| Initiative | Name | Goal | Status |
| --- | --- | --- | --- |
| Advanced 1 | Restructure Layout/UI | Rework the Web Admin into straightforward main menu, submenu, list screens, and detail pages so users do not have too many workflows on one screen. | Complete |
| Advanced 2 | Access/Refresh Token Authentication | Add production-grade login with short-lived access tokens, refresh token rotation/revocation, logout, and session expiry. | Complete |
| Advanced 3 | RBAC Enforcement | Enforce roles and permissions at the API level, map permissions to routes/actions, hide unauthorized UI actions, and preserve authenticated audit actors. | Complete |
| Advanced 4 | School Tree Management | Manage the hierarchy `school > school year/cohort > class > tuition/fee schedule/surcharges > students/parents`. | Complete |

Recommended order:

No remaining Advanced Production initiatives are currently defined.

Advanced 1 progress:

- Advanced 1A: App Shell is complete.
  - Replaced the horizontal tab strip with a sidebar/top-bar app shell.
  - Preserved existing tab IDs, buttons, and API behavior.
- Advanced 1B: Menu/Submenu Grouping is complete.
  - Grouped workflows into Tổng quan, Trường & học sinh, Học phí, Thanh toán, Liên lạc, and Quản trị.
  - Added top-bar screen title metadata that updates when switching workflows.
- Advanced 1C: List/detail patterns are complete.
  - Added a student detail panel for Học sinh.
  - Added an invoice detail panel with QR/PDF actions for Hóa đơn.
  - Added selected-row states for Đối soát invoice and transaction rows.
- Advanced 1D: Responsive cleanup and visual polish are complete.
  - Collapsed list/detail layouts to one column on mobile.
  - Prevented wide tables from expanding the page viewport.
  - Standardized selected rows, detail placeholders, and detail metadata spacing.

Advanced 1 launch prompt:

```text
Advanced 1 is complete. Continue with Advanced 2 from docs/initiatives/current-state.md when ready.
```

Advanced 2 launch prompt:

```text
Start Advanced 2 from docs/initiatives/current-state.md. Build production-grade login for ABC SUN with access tokens, refresh token rotation and revocation, logout, session expiry, and secure browser session handling.
```

Advanced 2 progress:

- Advanced 2 is complete.
  - Added auth schema with `password_hash` on `app_users`, auth sessions, access tokens, and rotated refresh tokens.
  - Added login, session, refresh, and logout APIs with HttpOnly SameSite browser cookies.
  - Stored only token hashes in PostgreSQL; refresh tokens are single-use and session revocation is recorded.
  - Added bootstrap admin support through `ABC_AUTH_BOOTSTRAP_EMAIL`, `ABC_AUTH_BOOTSTRAP_PASSWORD`, and `ABC_AUTH_BOOTSTRAP_DISPLAY_NAME`.
  - Protected production API routes with access-token middleware while keeping auth endpoints, provider webhooks, and QR PNG public.
  - Preserved authenticated audit actor context and added optional password setting in the user admin screen.
  - Added login/logout UI, session refresh handling, expiry recovery, tests, and docs.

Advanced 2 completion prompt:

```text
Advanced 2 is complete. Continue with Advanced 3 from docs/initiatives/current-state.md when ready.
```

Advanced 3 launch prompt:

```text
Start Advanced 3 from docs/initiatives/current-state.md. Enforce RBAC at the API level for ABC SUN, map permissions to routes and UI actions, and preserve authenticated audit actors.
```

Advanced 3 progress:

- Advanced 3 is complete.
  - Added a central server-side RBAC route map for protected `/api/v1` endpoints.
  - Enforced permissions from authenticated user roles instead of trusting `X-ABC-Admin-Permission`.
  - Kept auth endpoints, provider webhooks, and QR PNG public; kept bank options authenticated-only.
  - Added dynamic permission resolution for shared import-field and email-config endpoints.
  - Added email config/send/cron permission seeds in migration `0010_rbac_permissions`.
  - Preserved authenticated audit actor context for audited writes.
  - Updated the UI to hide unauthorized menus/actions, remove spoofable permission headers, and avoid loading workflows the user cannot access.
  - Added RBAC tests for route coverage, sensitive route mappings, dynamic resolvers, and header spoof rejection.

Advanced 3 completion prompt:

```text
Advanced 3 is complete. Continue with Advanced 4 from docs/initiatives/current-state.md when ready.
```

Advanced 4 launch prompt:

```text
Start Advanced 4 from docs/initiatives/current-state.md. Build tree-based management for school, school year/cohort, class, tuition and fee structures, students, and parents.
```

Advanced 4 progress:

- Advanced 4 is complete.
  - Added `schools` and `school_id` on `school_years` through migration `0011_school_tree_management`.
  - Backfilled existing school-year data into default school `ABC_SUN`.
  - Added server-side school tree APIs for reading the tree and saving schools, school years/cohorts, and classes.
  - Kept legacy master-data CSV compatible by making `school` optional and defaulting to `ABC_SUN`.
  - Extended master-data options/students with school metadata and `schoolId` filtering.
  - Added route-level RBAC for school tree read/write APIs using `master_data.read` and `master_data.write`.
  - Added a `Học sinh` school-tree panel with node-driven filters, compact edit forms, and a quick path into related fee schedules.
  - Added tests for migration coverage, school tree aggregation, input normalization, and RBAC route mapping.

Advanced 4 completion prompt:

```text
Advanced 4 is complete. Define the next Advanced Production initiative in docs/initiatives/current-state.md before implementation.
```

## Completed Follow-up: User Contact Bootstrap And Canonical RBAC

User request:

```text
Add user management with Name, Phone, Email where Phone or Email is required, role dropdown multi-select for Admin/Staff/Accountant, canonical permissions, and first-admin bootstrap UI when app_users is empty.
```

Completed:

- Added migration `0012_user_contacts_and_roles` to add `app_users.phone`, relax email-only identity, enforce Email-or-SĐT contact, seed roles `admin`, `staff`, `accountant`, and seed canonical `{module}.{action}` permissions.
- Added public auth bootstrap status/create API at `/api/v1/auth/bootstrap`; when no users exist, the UI shows the first Admin creation form before login.
- Updated login to accept Email or SĐT, while preserving HttpOnly access/refresh token behavior.
- Updated user admin API and UI to manage Tên, Email, SĐT, Password, Status, and multi-select role dropdown.
- Mapped protected routes to canonical permissions such as `user.view`, `student.update`, `invoice.create`, `payment.reconcile`, `report.export`, and kept legacy permission aliases for migrated roles.
- Updated tests and docs for user contact validation, bootstrap/login identifier behavior, migration coverage, and RBAC route mapping.

## Not Started

- None.

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
4. If the user asked to continue, start the next recommended Advanced Production initiative unless they specify a different initiative.
5. Use repo skills before implementation:
   - `$abcsun-change-workflow` by default
   - `$abcsun-vietqr-payments` for QR/payment/invoice amount behavior
   - `$abcsun-email-delivery` for email/cron/notification behavior
   - `$abcsun-frontend-ui` for Web Admin UI work
6. Update this file at the end of every initiative or meaningful partial implementation.

## Next Launch Prompt

Use this when the user says to continue without specifying a module:

```text
The production module roadmap is complete through Initiative 8 and Advanced 4 is complete. Review docs/initiatives/current-state.md and define the next Advanced Production initiative before implementation.
```

## Known Safety Constraints

- Do not send real email unless explicitly requested.
- Do not run real cron batches unless explicitly requested.
- Do not commit or print real secrets from local config files.
- Preserve `PaymentItems` total overriding raw `Amount`.
- Preserve VietQR TLV/CRC behavior with exact tests when touching QR generation.
