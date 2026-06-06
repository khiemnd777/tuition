# ABC SUN Initiative State

Last updated: 2026-06-06

## Current Status

Production roadmap implementation has student, parent, class master data, fee schedule setup, invoice/PDF receipt output, invoice issuance workbench, payment/reconciliation ledger, notification campaigns, communication campaign workbench, production Web Admin screens, reports/export, audit review, operational readiness, responsive/accessibility hardening, and operator guardrails complete.

Advanced Production work in the current roadmap is complete. Subscription conversion has tenant foundation, tenant-aware auth/RBAC, backend data isolation, tenant onboarding/switching, subscription hardening, tenant billing lifecycle enforcement, tenant entitlement/metering, subscription billing operations, subscription finance controls, cross-tenant finance operations, and subscription background automation complete.

Current phase: `platform_admin_platform_only_ux_hardening_complete`

Current initiative: Platform-only UX hardening is complete, clarifying control-plane session state versus tenant workspace state in Web Admin.

Next recommended initiative: No additional subscription/platform split phase is currently queued; the remaining follow-up is optional browser verification of the new platform-only UX.

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
- Subscription Phase 1: Tenant Foundation is complete:
  - Added tenant foundation migration `0014_tenant_foundation`.
  - Created `tenants` and `tenant_memberships`.
  - Seeded the default `ABC_SUN` tenant.
  - Backfilled all existing schools to the default tenant.
  - Backfilled existing app users as active default-tenant members.
  - Changed school code uniqueness from global `code` to `(tenant_id, code)`.
  - Updated school creation paths to write schools under the default tenant until tenant-aware auth/UI is introduced.
- Subscription Phase 2: Tenant-Aware Auth And RBAC is complete:
  - Added migration `0015_tenant_aware_auth_rbac`.
  - Added `tenant_id` to `app_auth_sessions` and backfilled existing sessions to default tenant.
  - Created `tenant_user_roles` and backfilled role assignments from `app_user_roles`.
  - Added active tenant data to auth session responses.
  - Made route-level RBAC require an active tenant context.
  - Loaded roles and permissions from tenant-scoped assignments.
  - Updated admin user/role management to read and write role assignments inside the active tenant.
- Subscription Phase 3: Data Isolation is complete:
  - Added migration `0016_tenant_data_isolation`.
  - Added `tenant_id` to `students`, `parents`, `notification_campaigns`, `audit_logs`, and `operation_logs`.
  - Backfilled tenant-owned records to the default `ABC_SUN` tenant or their owning school tenant.
  - Changed student code, parent email, and notification campaign code uniqueness to tenant scope.
  - Added active-tenant fail-fast checks to tenant-scoped API handlers.
  - Scoped master data, school tree, fee schedule, invoice, notification, payment dashboard, admin report/export, audit, operation, and readiness queries by active tenant.
  - Added tenant guards on write paths for school/year/class/student/parent/fee schedule/invoice/notification flows.
  - Kept public payment webhooks and provider credentials as a follow-up because provider tenant routing needs separate configuration ownership.
- Subscription Phase 4: Tenant Onboarding And Switching is complete:
  - Added migration `0017_tenant_onboarding_and_switching`.
  - Seeded tenant view/create/update/switch permissions and membership switcher index.
  - Added tenant list/save APIs for tenant onboarding and active tenant updates.
  - Added active tenant switch API that validates membership, issues a fresh browser session, and revokes the previous session.
  - Added Web Admin tenant switcher in the top bar.
  - Added tenant onboarding panel in `Người dùng & quyền` with create/edit tenant dialog and initial school creation.
  - Made frontend session recovery tenant-scoped and reset tenant-owned caches after switching.
  - Tenant-scoped payment provider credentials and public webhook routing were moved to Phase 5.
- Subscription Phase 5: Tenant-Scoped Payment Provider And Webhook Ownership is complete:
  - Added migration `0018_tenant_scoped_payment_providers_and_webhooks`.
  - Added `tenant_id` to `payment_providers`, seeded defaults for all tenants, and added tenant-level unique constraints and indexing.
  - Updated payment provider loading, provider list, webhook parsing, webhook event flow, and reconciliation queries to load/store by active tenant.
  - Updated payment reconciliation query paths so tenant isolation applies to transaction/list/match flows.
- Subscription Phase 6: Tenant Subscription Hardening is complete:
  - Reused `payment_providers.config` as the tenant-owned credential source for `payOS`, with env fallback kept only for missing fields and default-tenant compatibility.
  - Updated payment intent creation and webhook signature verification to resolve `payOS` secrets per tenant provider instead of from one global config only.
  - Added cross-tenant operation/audit read permissions for admin operators through migration `0019_subscription_hardening_cross_tenant_operations`.
  - Added tenant-scoped or cross-tenant filtering for operation and audit log APIs, including tenant identity in returned log rows.
  - Added tenant filtering in the `Vận hành` screen for operators with cross-tenant monitoring permission.
- Subscription Phase 7: Tenant Billing, Plan Lifecycle, And Subscription Enforcement is complete:
  - Added migration `0020_tenant_subscription_billing`.
  - Added `subscription_plans` and `tenant_subscriptions`, seeded `free_trial` and `standard`, and backfilled all existing tenants with a current subscription row.
  - Added subscription permissions and APIs for viewing plans and updating the active tenant subscription.
  - Extended auth session and tenant-admin responses with subscription status, plan, and lifecycle dates.
  - Added tenant subscription editing in Web Admin under the existing tenant management area.
  - Added centralized enforcement so write-heavy production workflows require the active tenant subscription to be `active` or `trial`, while read workflows and subscription repair remain available.
- Subscription Phase 8: Tenant feature entitlements, usage metering, and billing operations workflow is complete:
  - Added migration `0021_tenant_usage_entitlements`.
  - Added `tenant_usage_counters`, backfilled school/operator/student usage, and seeded monthly notification usage counters from sent notification logs.
  - Expanded default `free_trial` and `standard` plan limits to include `students` and `monthly_notifications` alongside `schools` and `operators`.
  - Added backend quota enforcement for school creation, student creation/import, tenant operator role assignment, and real notification sends.
  - Counted monthly notification usage only for real sent recipients, not dry runs or skipped duplicate recipients.
  - Extended tenant admin responses with per-tenant usage summaries so operators can see current plan consumption in Web Admin.
- Subscription Phase 9: Subscription invoicing, collections, and dunning automation is complete:
  - Added migration `0022_subscription_invoicing_and_dunning`.
  - Added `subscription_invoices`, `subscription_invoice_status_history`, and `subscription_dunning_runs` as a dedicated billing lane separate from school tuition invoices.
  - Added subscription billing APIs for listing invoices, generating the current period invoice, marking invoice payment manually, and running dry-run or real dunning sends for overdue subscription invoices.
  - Added backend sync so overdue open subscription invoices move to `past_due`, while paid subscription invoices move the tenant subscription back to `active` and update the current paid period window.
  - Added tenant-admin billing UI with invoice list, suggested next billing period, generate invoice action, manual mark-paid action, and dunning run dialog.
  - Reused the existing Gmail/Resend send path for dunning emails to tenant owner/admin/accountant recipients, while keeping manual collection as the payment confirmation path for this phase.
- Subscription Phase 10: Subscription self-service billing configuration, renewal controls, and finance exports is complete:
  - Added billing configuration save API over `tenant_subscriptions.billing_metadata` for amount, interval months, due days, auto renew, renewal mode, and finance note.
  - Exposed normalized billing config and renewal preview in the subscription billing response so tenant admins can manage billing without editing metadata manually.
  - Added finance CSV exports for subscription invoices, overdue invoices, paid invoices, and dunning history under the active tenant.
  - Added tenant-admin UI controls for billing config, renewal-oriented invoice generation from the suggested next period, and CSV export filters/actions.
  - Kept renewal execution manual-triggered while making the config explicit for a future scheduled automation phase.
- Subscription Phase 11: Cross-tenant finance console, scheduled renewals, and collection automation is complete:
  - Added cross-tenant subscription finance console APIs with filters for tenant search, subscription status, invoice status, auto renew, renewal mode, overdue-only, and missing-config-only scope.
  - Added scope-aware enforcement so `all` tenant finance actions require existing cross-tenant audit or operation permissions.
  - Added batch preview/run APIs for renewal invoice generation and dunning execution across multiple tenants using the same billing and email flows already implemented for single-tenant operations.
  - Added cross-tenant finance CSV exports for overview, overdue tenants, renewal candidates, and dunning-target tenants.
  - Added Web Admin cross-tenant finance console panel with summary cards, filters, renewal batch actions, dunning batch actions, and export trigger.
- Subscription Phase 12: Background scheduler for renewals/dunning and subscription suspension policy automation is complete:
  - Added migration `0023_subscription_automation_scheduler` to persist automation run logs across active-tenant and cross-tenant scopes.
  - Extended tenant subscription billing config with automation policy fields for enablement, renewal lead window, dunning cooldown, and overdue suspension threshold.
  - Added subscription automation status and run APIs so finance operators can preview or execute a combined automation cycle without switching tenant workflows manually.
  - Added background scheduler startup controlled by `ABC_SUBSCRIPTION_AUTOMATION_ENABLED` and `ABC_SUBSCRIPTION_AUTOMATION_INTERVAL`.
  - Reused the existing renewal invoice and dunning send flows, while adding cooldown-aware candidate selection for automated dunning and automatic suspension after overdue grace periods.
  - Added Web Admin automation panel in the subscription finance console to show scheduler state, latest run summary, and manual preview/run actions.
- Platform/Admin auth-session decoupling is complete:
  - Added migration `0026_platform_auth_sessions_nullable_tenant` so `app_auth_sessions.tenant_id` can be null for platform-only sessions.
  - Allowed `platform_admin` users to login and refresh without forcing an active tenant binding, while keeping tenant-required validation for non-platform users.
  - Kept finance and operations control-plane defaults on `all` scope when the session has no active tenant.
- Platform-only Web Admin UX hardening is complete:
  - Updated the auth badge to identify platform-only sessions explicitly as `Platform admin`.
  - Hid the tenant switcher cleanly when a platform-only session has no tenant memberships to switch into.
  - Replaced tenant-scoped empty states in tenant admin and subscription panels with control-plane copy that explains a tenant must be selected for tenant billing actions.
  - Disabled tenant-scoped subscription actions in platform-only mode until an active tenant is selected.
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
- Advanced 9: Tuition Setup Guided Workflow is complete:
  - Added school-aware fee schedule setup, filters, and saved schedule summaries with scope, period, item count, student count, preview total, updated actor, and updated timestamp.
  - Added guided setup steps, structured per-student adjustment rows, CSV paste fallback, preview issue rendering, and billing-recipient readiness in fee preview rows.
  - Added preview issues for empty/missing student scope and missing active billing recipients while preserving existing validation for invalid amounts and adjustment reasons.
  - Added quick handoff from a saved fee schedule into invoice preview/generation.
  - Added tests that keep fee preview totals aligned with invoice/payment item totals.
- Advanced 10: Invoice Issuance Workbench is complete:
  - Added step-based invoice workbench states for scope, preview, issues, and generation.
  - Added invoice preview idempotency metadata for ready-to-generate, already-generated, ready-to-regenerate, and blocked paid regeneration states.
  - Added blocking/warning issue rendering for paid regeneration and missing billing recipients.
  - Expanded issued invoice list metadata with base/adjusted totals, outstanding amount, item/adjustment counts, payment intent/match counts, sent counts, QR/PDF readiness, and issue state.
  - Added invoice detail API and UI snapshot for immutable line items, adjustments, operational counts, PDF/QR actions, notification handoff, payment intent handoff, and status history.
  - Added invoice CSV export handoff through the existing report export API and regression tests for invoice preview idempotency.
- Advanced 11: Communication Campaign Workbench is complete:
  - Added notification recipient metadata for outstanding amount, QR readiness, send count, last sent time, last log status, and retry eligibility.
  - Added selected-recipient email preview API and UI rendering for subject/HTML without sending real email.
  - Added step-based notification workbench states for target filters, recipients, email preview, send/logs, and cron.
  - Expanded recipient preview with parent email/name, student, invoice, outstanding amount, QR state, sent/retry state, status, and per-recipient selection.
  - Added explicit retry selected flow using saved campaign recipients, confirm dialog, `recipientIds`, `confirmSend=true`, and `forceResend=true`.
  - Added read-only cron snapshot in the notification workflow for enabled state, send time, daily limit, queued/sent/errors, rolling 24h sent count, next/last run, and recent results.
  - Kept reminder paid-invoice safeguards and real-send confirmation intact, with focused tests for summary retry/QR state, recipient filtering, and email-preview recipient selection.
- Advanced 12: Collection And Reconciliation Workbench is complete:
  - Added school, school year, grade, class, period, provider, invoice status, and transaction status filters to the reconciliation workflow.
  - Expanded reconciliation summary with receivable, collected, outstanding, collection rate, unpaid, partial, paid, overpaid, manual-review, and unmatched counts.
  - Added reconciliation match metadata to transaction rows and invoice detail: match type, status, score, applied amount, reason, and matched provider reference.
  - Added step-based reconciliation workbench states for scope, invoice ledger, transactions, match detail, and manual review.
  - Expanded invoice ledger actions for detail, QR/payment intent, payOS when configured, cash receipt, PDF, and notification handoff.
  - Added a manual-review queue for partial, overpaid, manual-review invoices, unmatched transactions, and transaction manual review.
  - Preserved idempotent payment ledger behavior and added focused tests for collection summary counts and stable invoice match scope.
- Advanced 13: Reports, Audit, And Operations Command Center is complete:
  - Added provider-aware report filtering and provider options to the reports API.
  - Expanded invoice report rows with outstanding amount, QR/PDF readiness, payment intent count, match count, notification sent count, and issue state.
  - Added detailed payment transaction rows to the reports response and CSV export, including match type, match status, score, applied amount, and match reason.
  - Expanded operation and audit log filters for operation/status/action/entity type/entity id and added command-center summary counts.
  - Sanitized operation/audit metadata before returning it to the UI so secret-like keys and raw payloads are redacted.
  - Added Reports UI transaction table and Operations UI summary/detail panel with sanitized metadata and workflow drilldown actions.
  - Added focused tests for transaction export match explanation, operation grouping, and metadata redaction.
- Advanced 14: Responsive, Accessibility, And UI Polish is complete:
  - Hardened app dialog accessibility with focus restoration, `role=alert` error output, busy state locking, loading affordance, and viewport-safe sizing.
  - Added keyboard activation and `aria-selected` synchronization for selectable workflow table rows.
  - Tightened responsive CSS for panel headers, action bars, dialog actions, status pills, metrics, table wrappers, and small-screen one-column layouts.
  - Prevented wide tables from expanding the page viewport while preserving horizontal scroll inside table wrappers.
  - Converted compact fee remove controls to icon-only buttons with screen-reader text and labels to avoid button text overflow.
  - Kept native browser `alert`, `confirm`, and `prompt` out of the UI and preserved app dialog confirmations.
  - Verified desktop workflow navigation across production tabs with no page-level horizontal overflow and no browser console warnings/errors.
- Advanced 15: Onboarding And Operator Guardrails is complete:
  - Added a Dashboard operator setup checklist for first-run path: school, school year/classes, student/parent import, billing recipients, email provider, fee schedule, and invoice preview.
  - Added role clarity in the User & Roles screen by summarizing Admin, Staff, Accountant, and custom-role permissions by view/create/update/send/reconcile/export/administer capability.
  - Added session recovery for active workflow and global school/year/period/month context after session expiry and login refresh, while clearing it on explicit logout.
  - Expanded risky action confirmations with actor, scope, provider, recipient, queue, invoice, import, and audit context for real sends, cron run/disable, invoice generation, imports, and fee schedule saves.
  - Added audit-bound cash receipt guardrails showing current actor and requiring reason before creating the manual cash ledger entry.
  - Kept provider secrets masked/blank in UI and continued using app dialogs instead of native browser prompts.

## Advanced Production Roadmap

| Initiative | Name | Goal | Status |
| --- | --- | --- | --- |
| Advanced 1 | Restructure Layout/UI | Rework the Web Admin into straightforward main menu, submenu, list screens, and detail pages so users do not have too many workflows on one screen. | Complete |
| Advanced 2 | Access/Refresh Token Authentication | Add production-grade login with short-lived access tokens, refresh token rotation/revocation, logout, and session expiry. | Complete |
| Advanced 3 | RBAC Enforcement | Enforce roles and permissions at the API level, map permissions to routes/actions, hide unauthorized UI actions, and preserve authenticated audit actors. | Complete |
| Advanced 4 | School Tree Management | Manage the hierarchy `school > school year/cohort > class > tuition/fee schedule/surcharges > students/parents`. | Complete |
| Advanced 5 | Task-Based Workflow And Navigation | Reframe the admin UI around daily school-accounting tasks, production navigation, context selectors, breadcrumbs, and permission-aware quick actions. | Complete |
| Advanced 6 | School Tree Relationship Workspace | Make `school > school year/cohort > grade > class` visible, navigable, editable, and connected to students, fee schedules, invoices, and readiness counts. | Complete |
| Advanced 7 | Student And Parent Relationship Workspace | Make student, parent, guardian, sibling, and billing-recipient relationships explicit with a clear list/detail UI and app-dialog editing. | Complete |
| Advanced 8 | Data Quality And Readiness Center | Surface blocking and warning issues before fee setup, invoice generation, notification sending, and reconciliation. | Complete |
| Advanced 9 | Tuition Setup Guided Workflow | Guide operators from school-tree scope selection through fee items, adjustments, preview, save, audit reason, and invoice handoff. | Complete |
| Advanced 10 | Invoice Issuance Workbench | Turn invoice generation into a step-based workbench with preview, blocking issues, idempotency visibility, QR/PDF detail, and bulk actions. | Complete |
| Advanced 11 | Communication Campaign Workbench | Make campaign targeting, billing-recipient resolution, dry-run preview, real-send confirmation, send logs, retries, and cron queues explicit. | Complete |
| Advanced 12 | Collection And Reconciliation Workbench | Center collection and reconciliation around invoices, payment intents, cash receipts, transaction matching, and manual-review queues. | Complete |
| Advanced 13 | Reports, Audit, And Operations Command Center | Consolidate filtered reports, exports, audit trails, provider/email/cron failures, and operational drilldowns. | Complete |
| Advanced 14 | Responsive, Accessibility, And UI Polish | Verify and harden all production workflows on desktop and mobile with stable layouts, dialogs, states, and accessible controls. | Complete |
| Advanced 15 | Onboarding And Operator Guardrails | Add first-run/operator guardrails, role clarity, session recovery, risky-action confirmations, and audit-bound action context. | Complete |

Recommended order:

1. Advanced 5: establish task-based navigation before deeper screen work.
2. Advanced 6: make the school/year/grade/class tree the primary setup context.
3. Advanced 7: make student-parent-billing relationships explicit and actionable.
4. Advanced 8: add readiness checks so missing data is visible before operators run production workflows.
5. Advanced 9 through Advanced 12: refine the fee, invoice, notification, and reconciliation workbenches in business-flow order.
6. Advanced 13 through Advanced 15: harden reports, operations, responsive behavior, onboarding, and guardrails.

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

Advanced 5 launch prompt:

```text
Start Advanced 5 from docs/initiatives/current-state.md. Reframe the ABC SUN admin UI around task-based school-accounting workflows, production navigation groups, context selectors, breadcrumbs, quick actions, legacy QR-tool demotion, and permission-aware menu/action visibility.
```

Advanced 5 planning:

- Goal: make the app read like a daily operations workflow instead of a collection of independent modules.
- Navigation groups:
  - `Tổng quan`: dashboard and work queue.
  - `Thiết lập`: school tree, students, parents, users, roles, permissions.
  - `Học phí`: fee schedules, invoice generation, invoice detail, QR/PDF.
  - `Thu tiền`: reconciliation, incoming transactions, cash receipt entry, manual review.
  - `Liên lạc`: notification campaigns, email preview, email config, cron queue.
  - `Báo cáo & vận hành`: reports, exports, audit logs, operation logs.
- Demote the legacy `Thanh toán` workflow into a secondary QR/import tool so it does not compete with production invoice workflows.
- Add top-bar context for current school, school year, period/month, and authenticated operator.
- Add breadcrumbs for deeper screens, such as `Thiết lập / Học sinh & phụ huynh / Lớp 1A`.
- Add dashboard quick actions for setup, fee schedule creation, invoice generation, notification preview/send, and reconciliation.
- Keep existing RBAC behavior by hiding inaccessible menus and actions and avoiding unauthorized data loads.
- Preserve current tab IDs and data contracts where possible to keep the change low risk.
- Use app dialog and confirm components for all production create/update/real-action workflows.
- Do not use native browser `alert`, `confirm`, or `prompt`.
- Mock up as text:

```text
Tổng quan / Việc cần xử lý

┌───────────────┬────────────────────────────────────────────────────┐
│ Tổng quan     │ ABC SUN · 2025-2026 · Kỳ 05/2026 · Operator        │
│  Dashboard    │                                                    │
│  Việc cần xử lý│ [12 thiếu recipient] [8 giao dịch chưa khớp]       │
│               │ [5 lớp chưa có bảng phí] [3 email lỗi]             │
│ Thiết lập     │                                                    │
│  Trường/lớp   │ Việc cần xử lý                                     │
│  Học sinh     │ Lớp 1A chưa có bảng phí kỳ hiện tại      [Mở]      │
│ Học phí       │ HS001 thiếu phụ huynh nhận phí           [Sửa]     │
│  Bảng phí     │ Giao dịch GD-889 chưa khớp invoice       [Đối soát]│
│  Hóa đơn      │                                                    │
│ Thu tiền      │ Bước nhanh                                         │
│  Đối soát     │ [Lập bảng phí] [Sinh hóa đơn] [Gửi thông báo]      │
└───────────────┴────────────────────────────────────────────────────┘
```

Advanced 5 progress:

- Advanced 5 is complete.
  - Reframed the sidebar into task-based production groups: `Tổng quan`, `Thiết lập`, `Học phí`, `Thu tiền`, `Liên lạc`, and `Báo cáo & vận hành`.
  - Renamed the dashboard entry to `Việc cần xử lý` and added dashboard work queue rows derived from existing receivable/reconciliation summary data.
  - Added permission-aware quick actions for student setup, fee schedule setup, invoice generation, notification campaign work, reconciliation, and the legacy QR/import tool.
  - Demoted the legacy `Thanh toán` tab into `Công cụ QR/import` under `Thu tiền`; legacy import/generate top-bar actions now only show on that tool.
  - Added top-bar breadcrumbs and compact school/year/period/month context controls that sync into the active workflow's existing filters where supported.
  - Preserved existing tab IDs, API contracts, app dialogs, confirmation components, and backend RBAC enforcement.
  - Updated README operator notes for the task-based Web Admin navigation.

Advanced 5 completion prompt:

```text
Advanced 5 is complete. Continue with Advanced 6 from docs/initiatives/current-state.md when ready.
```

Advanced 6 launch prompt:

```text
Start Advanced 6 from docs/initiatives/current-state.md. Build a clear School Tree Relationship Workspace for school, school year/cohort, grade, and class, with node details, counts, app-dialog editing, student handoff, fee-schedule handoff, and readiness indicators.
```

Advanced 6 planning:

- Goal: make `school > school year/cohort > grade > class` visible as the primary setup model.
- Left pane: tree with school, school year/cohort, grade, and class nodes.
- Tree node labels should show operational counts where available, such as class count, student count, missing-student count, fee-schedule count, invoice count, and issue count.
- Selecting a school shows school details, active years, class totals, student totals, and setup warnings.
- Selecting a school year/cohort shows grades, classes, active period context, student totals, fee schedule readiness, and invoice readiness.
- Selecting a grade shows all classes in that grade and aggregate totals.
- Selecting a class shows class detail, student roster, billing readiness, related fee schedules, related invoices, and quick actions.
- Add/edit school, school year/cohort, and class through app dialogs only.
- Keep inline tree selection fast and compact, but do not use inline upsert panels for production forms.
- Add empty states for missing school year, missing class, missing students, and missing fee schedule.
- Add quick action from class detail to filtered student list.
- Add quick action from class detail to fee schedule setup with scope preselected.
- Add quick action from class detail to invoice list/generation for the selected period.
- Preserve CSV compatibility: missing school in legacy imports defaults to the configured/default school.
- Preserve RBAC: read-only users can inspect, write users can open save dialogs.
- Mock up as text:

```text
Thiết lập / Trường & lớp

┌────────────────────┬───────────────────────────────────────────────┐
│ Cây trường         │ Chi tiết node đang chọn                       │
│ ABC SUN            │ Lớp 1A                                        │
│ └─ 2025-2026       │ Năm học: 2025-2026 · Khối 1 · 32 học sinh     │
│    ├─ Khối 1       │                                               │
│    │  ├─ Lớp 1A    │ Readiness                                     │
│    │  └─ Lớp 1B    │ Billing recipients: 30/32 đủ                 │
│    └─ Khối 2       │ Bảng phí kỳ 05/2026: chưa có                  │
│                    │ Hóa đơn kỳ 05/2026: chưa sinh                 │
│ [+ Trường]         │                                               │
│ [+ Năm học]        │ [Xem học sinh] [Lập bảng phí] [Sinh hóa đơn]  │
│ [+ Lớp]            │                                               │
└────────────────────┴───────────────────────────────────────────────┘
```

Advanced 6 progress:

- Advanced 6 is complete.
  - Extended the school tree API with period/month-aware readiness counts for billing recipients, active fee schedules, invoices, open invoice attention, and issue totals.
  - Aggregated readiness from class to grade, school year, and school levels while preserving existing school/year/class/student/fee aggregate behavior.
  - Upgraded the `Cây trường` panel into a relationship workspace with scan-friendly node badges, detail metrics, readiness rows, roster preview, and quick actions.
  - Connected class/year/grade detail actions to filtered students, scoped fee schedule setup, and invoice generation using the existing app dialog flows.
  - Kept school, school year, and class create/update forms in app dialogs instead of inline production upsert panels.
  - Updated README operator notes and added focused tests for readiness aggregation and period/month request validation.

Advanced 6 completion prompt:

```text
Advanced 6 is complete. Continue with Advanced 7 from docs/initiatives/current-state.md when ready.
```

Advanced 7 launch prompt:

```text
Start Advanced 7 from docs/initiatives/current-state.md. Build a Student And Parent Relationship Workspace that makes class membership, guardians, sibling links, billing-recipient rules, contact health, import conflicts, and dialog-based editing explicit.
```

Advanced 7 planning:

- Goal: make student, parent, guardian, sibling, and billing-recipient relationships understandable at a glance.
- Layout: school-tree/context pane, student list pane, student detail pane.
- Student filters: school, school year, grade, class, text search, billing readiness, missing contact, import conflict state.
- Search should cover student code, student name, parent name, parent email, and parent phone.
- Student list should show durable student code, student name, class, billing-recipient count, contact warning, invoice attention count where available.
- Student detail should show student code, name, school, school year, grade, class, and data quality state.
- Parent relationship table should show parent name, relationship label, email, phone, primary flag, active flag, email active flag, billing-recipient flag, and warning state.
- Billing-recipient rule must be visible: a recipient is valid only when parent relationship is active, parent email is active, receives-billing flag is true, and email exists.
- Support one student with many parents/guardians.
- Support one parent linked to many students, including siblings.
- Show sibling links when a parent is shared across multiple students.
- Dialog editing:
  - Student identity fields: student code, name, class.
  - Parent rows: name, relationship, email, phone, primary, active, email active, receives billing.
  - Add/remove parent rows without layout shifts.
  - Validate at least one contact method where required by backend rules.
  - Warn when no valid billing recipient remains.
- Import workflow:
  - Keep mapping preview before apply.
  - Show conflicts by student code, parent identity, class mismatch, email mismatch, and duplicate relationship.
  - Never silently overwrite mismatched student/parent data.
- Quick actions: edit student, open invoices, open notifications, open class detail.
- Mock up as text:

```text
Thiết lập / Học sinh & phụ huynh

┌────────────────────┬──────────────────────────┬──────────────────────┐
│ Cây trường         │ Danh sách học sinh        │ Chi tiết học sinh     │
│ ABC SUN            │ [Search] [Khối] [Lớp]     │ HS001 - Nguyễn An     │
│ └─ 2025-2026       │ [Thiếu contact ▾]         │ Lớp 1A · 2025-2026   │
│    └─ Khối 1       │ HS001 Nguyễn An      1A   │                      │
│       └─ Lớp 1A    │ HS002 Trần Bình      1A   │ Phụ huynh             │
│                    │ HS003 Lê Chi         1A   │ Mẹ · chính · nhận phí │
│                    │                          │ Bố · phụ · không gửi  │
│ [+ Học sinh]       │ [+ Học sinh] [Import]     │ [Sửa] [Xem hóa đơn]   │
└────────────────────┴──────────────────────────┴──────────────────────┘
```

Advanced 7 progress:

- Advanced 7 is complete.
  - Extended the master-data student API with parent count, billing-recipient count, missing billing state, contact warning, invoice attention count, relationship labels, parent phone, parent billing-ready flags, and sibling links through shared active parents.
  - Extended master-data import mapping for parent phone and relationship while preserving existing CSV compatibility.
  - Upgraded the `Học sinh & phụ huynh` table into a relationship scan view with billing/contact/invoice attention indicators and a billing readiness filter.
  - Upgraded student detail with relationship metrics, visible billing-recipient rule, parent relationship table, sibling links, and quick actions to invoices, notifications, and class tree scope.
  - Updated the app-dialog student editor to capture parent phone and relationship labels, and to warn when no valid billing recipient remains.
  - Added focused tests for relationship-state derivation, parent contact normalization, and UUID placeholder generation.

Advanced 7 completion prompt:

```text
Advanced 7 is complete. Continue with Advanced 8 from docs/initiatives/current-state.md when ready.
```

Advanced 8 launch prompt:

```text
Start Advanced 8 from docs/initiatives/current-state.md. Add a Data Quality And Readiness Center that surfaces blocking and warning issues before fee setup, invoice generation, notification sending, and reconciliation.
```

Advanced 8 planning:

- Goal: prevent operators from discovering missing data only after running a workflow.
- Add a dashboard/work-queue view with severity groups: blocking, warning, informational.
- Readiness checks:
  - Student has no class.
  - Class has no students.
  - Student has no parent/guardian.
  - Student has no valid billing recipient.
  - Billing recipient has no email.
  - Billing recipient email is marked inactive.
  - Parent relationship is inactive but selected for billing.
  - Duplicate student code in import preview.
  - Parent/student import conflict.
  - Class has no fee schedule for selected period.
  - Fee schedule has empty or zero-value required fee items.
  - Student adjustment lacks reason.
  - Invoice preview contains blocked students.
  - Invoice exists but has no QR/payment data.
  - Invoice is unpaid, partial, overpaid, or manual review.
  - Incoming transaction is unmatched.
  - Notification campaign has failed recipients.
  - Email provider is not configured.
  - Cron queue has errors or over-limit sends.
- Each issue should have entity type, entity id, scope, severity, message, and action target.
- Filters: school, school year, grade, class, period/month, issue type, severity.
- Each issue row should deep-link to the relevant edit/detail workflow with the right filters preselected.
- Acceptance: operators can resolve readiness blockers before generating invoices or sending email.

Advanced 8 progress:

- Advanced 8 is complete.
  - Extended `/api/v1/admin/dashboard` with a `readiness` payload that groups issues by blocking, warning, and info severity.
  - Added school-aware dashboard/report filters and preserved existing school year, grade, class, period, month, and invoice status filters.
  - Added readiness checks for missing student relationships, invalid billing recipients, inactive billing contacts, empty classes, missing active fee schedules for the selected period, empty/zero-value fee schedules, invoice payment/status issues, unmatched/manual-review transactions, failed notification recipients, email provider readiness, cron queue/quota problems, and recent cron operation errors.
  - Added dashboard UI for the Data Quality & Readiness Center with summary cards, severity/type filters, grouped issue rows, and action handoff to students, fees, invoices, notification, reconciliation, email/cron, and operation logs.
  - Added focused tests for readiness summary ordering and invoice issue classification.

Advanced 8 completion prompt:

```text
Advanced 8 is complete. Continue with Advanced 9 from docs/initiatives/current-state.md when ready.
```

Advanced 9 launch prompt:

```text
Start Advanced 9 from docs/initiatives/current-state.md. Build a guided Tuition Setup workflow from school-tree scope selection through fee items, adjustments, preview, save, audit reason, and invoice handoff.
```

Advanced 9 planning:

- Goal: make fee schedule setup deterministic, scoped, previewable, and auditable.
- Scope selection should come from school, school year, grade, class, period, and month.
- Fee items should include tuition, lunch, shuttle, uniform, insurance, materials, previous fees, and custom fees.
- Fee item labels should keep Vietnamese and English labels where supported.
- Allow per-student adjustments: discount, surcharge, waiver, carry-over.
- Require reason for adjustments.
- Preview totals by student before save.
- Show issues for missing students, empty class, invalid amount, missing adjustment reason, and no valid billing recipient.
- Saving should preserve `PaymentItems` total-overrides-amount behavior downstream.
- Saved schedules list should show scope, period, item count, student count, total preview, updated actor, and updated timestamp.
- Quick action from saved schedule to invoice preview/generation.
- Acceptance: previewed totals match invoice item totals after generation.

Advanced 10 launch prompt:

```text
Start Advanced 10 from docs/initiatives/current-state.md. Build an Invoice Issuance Workbench with step-based preview, blocking issues, idempotency visibility, generated invoice list, QR/PDF detail, and bulk actions.
```

Advanced 10 planning:

- Goal: make invoice generation safe, transparent, and invoice-centered.
- Step 1: select scope and saved fee schedule.
- Step 2: preview generated invoice rows.
- Step 3: show blocking issues and warnings.
- Step 4: confirm invoice generation.
- Show idempotency state: not generated, already generated, regenerated only if supported by explicit operator action.
- Invoice preview rows should show student, class, line-item count, total, billing-recipient state, and issue state.
- Generated invoice list should show invoice code, student, class, period, total, paid amount, outstanding, status, sent state, QR/payment state, PDF state.
- Invoice detail should show immutable line-item snapshot, adjustments, QR, bill number, PDF link, payment history, notification history, and status history.
- Bulk actions: export CSV, open PDF, generate missing QR/payment intent, open notification flow for selected unpaid invoices.
- Acceptance: no duplicate invoices for the same class/period/default generation path.

Advanced 10 completion prompt:

```text
Advanced 10 is complete. Continue with Advanced 11 from docs/initiatives/current-state.md when ready.
```

Advanced 11 launch prompt:

```text
Start Advanced 11 from docs/initiatives/current-state.md. Build a Communication Campaign Workbench with recipient resolution, dry-run preview, email preview, send confirmation, logs, retries, cron queues, and paid-invoice safeguards.
```

Advanced 11 planning:

- Goal: make invoice-based communication safe and inspectable before any real email is sent.
- Campaign types: first payment notice, reminder, paid receipt/confirmation where supported.
- Target filters: school, school year, grade, class, period, due date, invoice status, campaign type.
- Recipient resolver should use active parent relationships with active billing email flags.
- Dry-run is the default validation path.
- Preview should show recipient email, parent name, student name, invoice code, amount, status, and QR availability.
- Email preview should show subject and rendered HTML for selected template and selected recipient/invoice.
- Real send must use confirm dialog and `confirmSend`.
- Paid invoices must be excluded from reminder campaigns.
- Logs should show campaign, template version, invoice, recipient, provider, status, provider message id, error, and timestamp.
- Retry should be explicit and idempotent unless operator deliberately re-sends.
- Cron controls should show enabled state, send time, daily limit, queued count, sent count, error count, and recent results.
- Acceptance: operators know exactly which parent email receives which invoice before sending.

Advanced 11 completion prompt:

```text
Advanced 11 is complete. Continue with Advanced 12 from docs/initiatives/current-state.md when ready.
```

Advanced 12 launch prompt:

```text
Start Advanced 12 from docs/initiatives/current-state.md. Build a Collection And Reconciliation Workbench centered on invoices, payment intents, cash receipts, transaction matching, underpayment, overpayment, and manual-review queues.
```

Advanced 12 planning:

- Goal: make every money movement visible from the invoice and every unmatched transaction actionable.
- Summary should show receivable, collected, outstanding, collection rate, unpaid, partial, overpaid, manual review, and unmatched transaction counts.
- Invoice ledger filters: school, school year, grade, class, period, provider, invoice status, transaction status.
- Invoice row actions: view detail, create payment intent/QR, enter cash receipt, open PDF, open notification.
- Invoice detail should show total, paid amount, outstanding, payment intents, transactions, reconciliation matches, status history, and collection account/reference.
- Payment intent detail should show provider, provider reference, payment URL when available, QR, expiry/status where available.
- Cash receipt dialog must require amount, collector/operator, reason, and optional receipt reference.
- Transaction table should show provider, reference, account, amount, description, received time, matched invoice, match confidence/reason, and status.
- Matching explanation should show invoice code, amount, provider reference, collection account, and mismatch reasons.
- Underpayment should mark invoice partial.
- Exact payment should mark invoice paid.
- Overpayment should mark invoice overpaid or manual review according to existing policy.
- Duplicate webhook/event should not duplicate ledger entries.
- Acceptance: unmatched and manual-review items are easy to find and do not hide among paid invoices.

Advanced 12 completion prompt:

```text
Advanced 12 is complete. Continue with Advanced 13 from docs/initiatives/current-state.md when ready.
```

Advanced 13 launch prompt:

```text
Start Advanced 13 from docs/initiatives/current-state.md. Build a Reports, Audit, And Operations Command Center with scoped report filters, exports, audit drilldowns, operation logs, and provider/email/cron failure review.
```

Advanced 13 planning:

- Goal: make end-of-day review and troubleshooting possible from one operations area.
- Report filters should reuse school, school year, grade, class, period, month, status, and provider context.
- Class summary report: receivable, collected, outstanding, collection rate, unpaid count, partial count, paid count.
- Invoice detail report: student, class, invoice code, period, total, paid amount, outstanding, status, sent state, provider state.
- Payment transaction export: provider, reference, amount, account, matched invoice, received timestamp, status.
- Export CSV should preserve current filters.
- Audit log should show actor, action, entity type, entity id, reason, timestamp, and immutable details.
- Operation logs should group webhook, email, cron, and background job failures.
- Log detail should avoid secrets and raw credentials.
- Drilldowns should open the related invoice, student, class, campaign, or transaction when ids are available.
- Acceptance: operator can explain who changed data, what failed, and what still needs action.

Advanced 13 completion prompt:

```text
Advanced 14 is complete. Continue with Advanced 15 from docs/initiatives/current-state.md when ready.
```

Advanced 14 launch prompt:

```text
Start Advanced 14 from docs/initiatives/current-state.md. Verify and harden responsive behavior, accessibility, layout stability, dialog behavior, empty/loading/error states, and visual polish across all production workflows.
```

Advanced 14 planning:

- Goal: make the production UI reliable across desktop and mobile without layout overlap or confusing state changes.
- Check desktop and mobile layouts for every production workflow.
- Tables may scroll horizontally but must not expand the page viewport.
- Detail panels collapse cleanly on mobile.
- Fixed-format UI such as toolbars, status pills, icon buttons, counters, and table cells should not resize unexpectedly on hover or status update.
- Text must not overflow buttons, cards, pills, table cells, or dialog actions.
- Dialogs should manage focus, error display, close behavior, and action loading states consistently.
- Buttons should use icons where familiar and labels where actions need text.
- Empty states, loading states, success states, and error states should be consistent.
- Selected row styling should be consistent across student, invoice, and reconciliation tables.
- No nested cards for page sections.
- No native browser `alert`, `confirm`, or `prompt`.
- Browser verification should cover the changed workflows.
- Acceptance: no overlapping UI, no major layout shift, and no unusable workflow on mobile.

Advanced 14 completion prompt:

```text
Advanced 14 is complete. Continue with Advanced 15 from docs/initiatives/current-state.md when ready.
```

Advanced 15 launch prompt:

```text
Start Advanced 15 from docs/initiatives/current-state.md. Add onboarding and operator guardrails for first-run setup, role clarity, session recovery, risky-action confirmations, audit-bound actions, and secret-safe operational UI.
```

Advanced 15 planning:

- Goal: reduce operator mistakes during setup and production use.
- First-run checklist after bootstrap admin:
  - Create or verify school.
  - Create school year/cohort.
  - Create classes.
  - Import students and parents.
  - Resolve missing billing recipients.
  - Configure email provider.
  - Create first fee schedule.
  - Preview first invoice batch.
- Role templates should be understandable for Admin, Staff, and Accountant.
- Permission summaries should show what a user can view, create, update, send, reconcile, export, and administer.
- Session expiry should return to login cleanly and recover context where practical.
- Risky actions should have clear confirm dialogs: send real email, run cron now, generate invoices, apply import, save fee adjustments, enter cash receipt, disable cron.
- Audit-bound actions should show current actor and reason requirements when applicable.
- Do not show secrets or local config values in UI/logs.
- Acceptance: first-time operators know the minimum setup path and risky actions remain intentional.

Advanced 15 completion prompt:

```text
Advanced 15 is complete. Review, commit, and push the Advanced Production hardening changes when ready, or define the next roadmap initiative.
```

## Completed Follow-up: User Contact Bootstrap And Canonical RBAC

User request:

```text
Add user management with Name, Phone, Email where Phone or Email is required, role dropdown multi-select for Admin/Staff/Accountant, canonical permissions, and first-admin bootstrap UI when app_users is empty.
```

Completed:

- Phase B tenant signup and owner onboarding is complete:
  - Added public tenant signup API for creating a tenant, first school, and first `tenant_owner` account without platform-admin intervention.
  - Preserved first-user platform bootstrap separately for `platform_admin`, so system bootstrap and tenant self-serve signup are no longer conflated.
  - Extended auth session responses with `onboarding` metadata for newly created tenant-owner sessions.
  - Added public login-screen flow for `Đăng ký trường`, auto-login after signup, and tenant onboarding banner after first session.
- Phase C self-serve subscription purchase is complete:
  - Added tenant-facing `Gói & Thanh toán` workspace for `subscription.view|subscription.update` actors, independent from platform control-plane tabs.
  - Added self-serve purchase API for owner checkout discovery and checkout execution over the existing subscription invoice lane.
  - Reused tenant payment providers for checkout generation: `manual_vietqr` / `sepay` return VietQR payloads, while `payOS` returns checkout link data when configured.
  - Reused existing subscription invoice generation so checkout stays aligned with subscription billing, lifecycle status, and finance visibility already built in earlier phases.
- Phase D subscription payment confirmation automation is complete:
  - Added migration `0025_subscription_payment_confirmation_link` so `payment_transactions` can link directly to `subscription_invoices`.
  - Extended webhook/payment reconciliation to try tuition invoice matching first, then subscription invoice matching by invoice code / stored provider reference and exact amount.
  - Added automatic subscription invoice confirmation so matched checkout payments move the subscription invoice to `paid` and reactivate/update the tenant subscription period without requiring manual mark-paid fallback.
  - Updated tenant-facing subscription refresh flow so owners can reload checkout/billing state and see the paid transition after provider events are processed.
- Phase A actor-model correction is complete:
  - Added migration `0024_actor_model_role_split` to seed `platform_admin`, `tenant_owner`, `tenant_admin`, `tenant_staff`, `tenant_accountant`, backfill legacy role assignments, and remove legacy tenant-scoped `admin|staff|accountant` rows.
  - Updated auth session enrichment to expose `platformRoles`, `activeTenantRoles`, `platformRoleCodes`, `activeTenantRoleCodes`, `isPlatformAdmin`, and `isTenantOwner`, while merging platform + tenant permission sets at runtime.
  - Moved tenant user/role management out of `Platform Admin` UI into a dedicated `Tenant Access` workspace so platform control plane is no longer sharing the tenant admin surface.
- Platform/Admin split cleanup continuation is complete:
  - Added dedicated platform user management APIs for listing, saving, and assigning `platform_admin` roles independently from tenant role assignment.
  - Added a separate `Platform users` panel under `Platform Admin`, while keeping `Tenant Access` scoped to tenant operator roles only.
  - Extended RBAC route metadata so platform-only control-plane endpoints can be authorized without reusing the tenant-scoped active-tenant gate.
  - Kept tenant onboarding, tenant subscription management, and tenant operator management logically separated in the UI and API surface.
- Platform/Admin split hardening follow-up is complete:
  - Removed the early active-tenant requirement from cross-tenant finance console, finance batch, automation status/run, and audit/operation scope resolution.
  - Defaulted platform control-plane finance and automation queries to `scope=all` when the current session has no active tenant context.
  - Updated operations filtering so `platform_admin` can inspect all-tenant logs without being forced through the active-tenant fallback first.
  - Kept tenant-scoped billing, checkout, invoice, and operator workflows unchanged so tenant-owned write paths still require a concrete active tenant.
- Platform/Admin auth-session decoupling is complete:
  - Added migration `0026_platform_auth_sessions_nullable_tenant` so `app_auth_sessions.tenant_id` is no longer mandatory for platform-only sessions.
  - Updated login and session issuance to prefer the active tenant when present, but allow `platform_admin` to create a valid session with no tenant binding.
  - Updated access-token and refresh-token session loading/validation so platform-only sessions remain refreshable and authenticated without a synthetic tenant membership.
  - Kept tenant-owner and tenant-admin sessions backward compatible by preserving tenant-bound session creation when an active tenant exists.
- Added migration `0012_user_contacts_and_roles` to add `app_users.phone`, relax email-only identity, enforce Email-or-SĐT contact, seed roles `admin`, `staff`, `accountant`, and seed canonical `{module}.{action}` permissions.
- Added public auth bootstrap status/create API at `/api/v1/auth/bootstrap`; when no users exist, the UI shows the first Admin creation form before login.
- Updated login to accept Email or SĐT, while preserving HttpOnly access/refresh token behavior.
- Updated user admin API and UI to manage Tên, Email, SĐT, Password, Status, and multi-select role dropdown.
- Mapped protected routes to canonical permissions such as `user.view`, `student.update`, `invoice.create`, `payment.reconcile`, `report.export`, and kept legacy permission aliases for migrated roles.
- Updated tests and docs for user contact validation, bootstrap/login identifier behavior, migration coverage, and RBAC route mapping.

## Not Started

- No Advanced Production roadmap initiatives remain in the current plan.

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
Advanced 15 is complete. Review, commit, and push the Advanced Production hardening changes when ready, or define the next roadmap initiative.
```

## Known Safety Constraints

- Do not send real email unless explicitly requested.
- Do not run real cron batches unless explicitly requested.
- Do not commit or print real secrets from local config files.
- Preserve `PaymentItems` total overriding raw `Amount`.
- Preserve VietQR TLV/CRC behavior with exact tests when touching QR generation.
