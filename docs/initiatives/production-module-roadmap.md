# ABC SUN Production Module Roadmap

Last reviewed: 2026-05-22

This document records the production-oriented module analysis for the ABC SUN school payment system. It is intended to be used as a launch point for separate implementation initiatives on different days.

## Production Direction

The production system should not be built around temporary payment rows. The durable business model should be:

```text
Student -> Invoice -> PaymentTransaction -> Reconciliation -> NotificationLog
```

The current `paymentRow` flow remains useful as an import, preview, and QR/email engine, but it should not be the long-term source of truth.

Core production rules:

- Use PostgreSQL or another transactional database for business data.
- Keep QR generation and email rendering as engines behind durable invoice/payment records.
- Make every money movement auditable.
- Do not delete payment transactions; use reversal or adjustment records.
- Enforce permissions at API level, not only in the UI.
- Keep staging and production separated.
- Never send real email, run real cron batches, or store real secrets in tracked files unless explicitly operating production with approved procedures.

## Payment Gateway Direction

For school tuition, prefer VietQR plus bank/webhook reconciliation over percentage-fee hosted gateways. Tuition transactions are usually high value, so percentage fees can become expensive quickly.

Recommended production posture:

- Primary long-term path: negotiate direct bank collection, virtual accounts, or tuition collection service when school scale justifies it.
- Practical production path: integrate a payment provider adapter with payOS or SePay first, then add direct bank collection later without changing the invoice ledger.
- Do not lock the domain model to one provider. Use a `payment_providers` abstraction and store raw provider events.
- Use VNPAY, MoMo, or ZaloPay as optional secondary channels only when the school needs wallets, cards, or hosted checkout.

Provider notes checked on 2026-05-22:

- payOS: official site states free setup, maintenance, and transaction fees from 2026-01-23, with money going directly to the bank account. Sources: https://payos.vn/ and https://payos.vn/thu-ho/
- SePay: supports VietQR webhook/API reconciliation; pricing page includes a free tier and paid plans by transaction volume. Sources: https://sepay.vn/bang-gia.html and https://sepay.vn/xac-thuc-thanh-toan.html
- Casso: useful for cash-flow management and bank API/webhook; pricing depends on plan and API needs. Source: https://api.casso.vn/pricing-table
- VNPAY: official fee page states free or up to 2.5% transaction value. Source: https://vnpay.vn/bieu-phi-dich-vu-0f61zmr3hwgs

## Initiative 1: Foundation And Persistence

Goal: establish production-grade persistence and cross-cutting controls before building modules.

Scope:

- Add database migration system.
- Add DB connection/configuration with separate local/staging/production settings.
- Add base tables for audit logging, users, roles, and permissions.
- Define ID conventions and timestamp/user audit columns.
- Add backup and restore runbook.

Acceptance:

- Database schema can be created from scratch by migrations.
- Migrations can run repeatedly without data loss.
- Business tables have `created_at`, `updated_at`, and auditable actor fields where needed.
- Sensitive config is environment-based or secret-managed, not committed.

Launch prompt:

```text
Start Initiative 1 from docs/initiatives/production-module-roadmap.md. Design and implement the production persistence foundation for ABC SUN, preserving current VietQR/email behavior.
```

## Initiative 2: Student, Parent, Class Master Data

Goal: replace embedded student/parent fields in payment rows with durable master data.

Production entities:

- `school_years`
- `classes`
- `students`
- `parents`
- `student_parents`

Required behavior:

- Student code is mandatory and unique per school scope.
- A student can have multiple parents/guardians.
- A parent can be linked to multiple students.
- Parent email has delivery flags: primary, active, receives billing email.
- Students can be listed by class, grade, and school year.
- CSV import maps legacy payment-row fields into master data with conflict reporting.

Acceptance:

- No production workflow relies on student name as an identifier.
- Duplicate names are supported safely.
- Import does not silently overwrite conflicting parent/student data.
- UI supports filtering students by class and grade.

Launch prompt:

```text
Start Initiative 2 from docs/initiatives/production-module-roadmap.md. Build production student, parent, and class master data with imports and class/grade listing.
```

## Initiative 3: Fee Types And Fee Schedules

Goal: define reusable tuition and fee rules before invoices are generated.

Production entities:

- `fee_types`
- `fee_schedules`
- `fee_schedule_items`
- optional `student_fee_adjustments`

Required behavior:

- Fee types support Vietnamese and English labels.
- Fee schedules are scoped by class, grade, school year, period, and month where applicable.
- Fee schedule items define default amounts for tuition, lunch, shuttle, uniform, insurance, previous fees, and custom fees.
- Per-student adjustments support discount, surcharge, waiver, and carry-over.
- Adjustments require reason and audit trail.

Acceptance:

- A schedule can be previewed before invoice generation.
- Fee totals are deterministic and test-covered.
- Existing `PaymentItems` total-overrides-amount invariant remains preserved when invoice items are converted to QR/payment data.

Launch prompt:

```text
Start Initiative 3 from docs/initiatives/production-module-roadmap.md. Build production fee types and fee schedules for class/period/month tuition setup.
```

## Initiative 4: Invoice And PDF Receipt

Goal: make invoices the source of truth for payment requests and receipts.

Production entities:

- `invoices`
- `invoice_items`
- `invoice_adjustments`
- `invoice_status_history`
- `receipt_documents`

Required behavior:

- Generate invoices from fee schedules for selected classes or students.
- Invoice codes are stable and unique.
- Invoice items preserve labels and amounts at generation time.
- QR payloads use invoice code as the payment reference.
- PDF receipts are generated from invoice data, not from ad hoc UI rows.
- Invoice status is derived from payment totals unless an authorized adjustment is recorded.

Acceptance:

- Invoice generation is idempotent for the same class/period unless explicitly regenerated.
- Generated PDF contains school name, student, class, period, line items, total, QR, status, and issue timestamp.
- Invoice code maps cleanly to VietQR `BillNumber`.

Launch prompt:

```text
Start Initiative 4 from docs/initiatives/production-module-roadmap.md. Build invoice generation and PDF receipt output as the production source of payment requests.
```

## Initiative 5: Payment And Reconciliation

Goal: record real money movements and reconcile them to invoices.

Production entities:

- `payment_providers`
- `payment_intents`
- `payment_transactions`
- `provider_events`
- `reconciliation_matches`
- `manual_cash_receipts`

Required behavior:

- Create a payment intent per invoice and provider.
- Generate QR or payment link through the selected provider adapter.
- Accept provider webhooks with idempotency.
- Store raw webhook payloads before parsing.
- Match incoming bank transactions by invoice code, amount, account, and provider reference.
- Support cash payment entry by authorized staff.
- Status outcomes: `unpaid`, `partial`, `paid`, `overpaid`, `manual_review`.

Provider adapter baseline:

- `manual_vietqr`: current QR generation, manual reconciliation.
- `sepay`: webhook/API based reconciliation.
- `payos`: payment link/QR and webhook based reconciliation.
- Future adapters: direct bank collection, VNPAY, MoMo, ZaloPay.

Acceptance:

- Duplicate webhook does not duplicate a payment.
- Underpayment marks invoice `partial`.
- Exact payment marks invoice `paid`.
- Overpayment marks invoice `overpaid` or `manual_review` based on policy.
- Cash payments are auditable with collector, timestamp, amount, and receipt reference.

Launch prompt:

```text
Start Initiative 5 from docs/initiatives/production-module-roadmap.md. Build production payment provider adapters, webhook ingestion, transaction ledger, and invoice reconciliation.
```

## Initiative 6: Notification Campaigns

Goal: move email sending from row-based batches to invoice-based campaigns.

Production entities:

- `notification_templates`
- `notification_campaigns`
- `notification_recipients`
- `notification_logs`

Required behavior:

- Templates are configurable and versioned.
- Campaigns target invoices by school year, period, class, status, and due date.
- First notice targets newly issued invoices.
- Reminder notice targets `unpaid` and `partial` invoices.
- Logs record recipient, template version, provider, status, provider message ID, and error.
- Existing preview and dry-run behavior remains default validation path.

Acceptance:

- Campaign dry-run shows exact recipients and invoice totals before sending.
- Reminder campaign cannot accidentally include already paid invoices.
- Sending is idempotent per campaign/template/invoice/recipient unless explicitly re-sent.
- Real sends remain protected by confirmation and environment safeguards.

Launch prompt:

```text
Start Initiative 6 from docs/initiatives/production-module-roadmap.md. Build invoice-based notification templates, campaigns, dry-run, send logs, and reminders.
```

## Initiative 7: Web Admin

Goal: replace the single batch tool with production admin workflows.

Required areas:

- Students and parents
- Classes and school years
- Fee schedules
- Invoices and PDF receipts
- Payment reconciliation
- Notification campaigns
- Dashboard and reports
- User and role administration

Dashboard metrics:

- Total receivable
- Total collected
- Outstanding amount
- Collection rate
- Unpaid student count
- Partial payment count
- Overpaid/manual-review count
- Top classes by outstanding amount

Acceptance:

- Admin workflows are dense and operational, not marketing-style pages.
- All writes go through API permissions.
- Tables support filtering by class, grade, period, month, and status.
- Payment reconciliation screen clearly separates matched, unmatched, partial, and overpaid transactions.

Launch prompt:

```text
Start Initiative 7 from docs/initiatives/production-module-roadmap.md. Build production Web Admin workflows for students, fees, invoices, payments, notifications, dashboard, reports, and roles.
```

## Initiative 8: Reports, Audit, And Operations

Goal: make the system usable for accounting, school management, and production operations.

Required behavior:

- Export reports by class, period, month, and payment status.
- Export invoice/payment data to Excel or CSV.
- Maintain immutable audit logs for money and fee changes.
- Add operational logs for webhook, email, and background job failures.
- Add backup, restore, deployment, and incident runbooks.
- Add staging smoke tests and production readiness checklist.

Acceptance:

- Accounting can reconcile total invoices against total transactions.
- Every manual payment and fee adjustment has an actor and reason.
- Production deployment has documented rollback and backup verification.

Launch prompt:

```text
Start Initiative 8 from docs/initiatives/production-module-roadmap.md. Build reporting, audit, exports, and operational readiness for production.
```

## Recommended Build Order

1. Foundation and persistence.
2. Student, parent, class master data.
3. Fee types and fee schedules.
4. Invoice and PDF receipt.
5. Payment and reconciliation.
6. Notification campaigns.
7. Web Admin.
8. Reports, audit, and operations.

This order keeps the domain model stable before provider-specific payment work and avoids building dashboard or email flows on temporary data.

## Current Code To Preserve As Engines

- VietQR payload and PNG generation in `main.go` and `vietqr_standard.go`.
- `PaymentItems` total overriding raw `Amount`.
- `BillNumber` mapping to VietQR Additional Data `62-01`.
- `Note` mapping to VietQR Additional Data `62-08`.
- Email preview and dry-run behavior.
- Gmail/Resend provider separation behind `sendRenderedEmail`.
- Cron rolling 24-hour quota rules until replaced by production notification jobs.

