# ABC SUN VietQR Context

## Language

**Payment row**:
A single student payment instruction. It contains student, parent, class, bank BIN, account, email, amount, payment items, bill number, and note.
Avoid: record, line item when referring to the full row.

**Payment item**:
One fee component inside a payment row, such as tuition, shuttle, insurance, uniform, material, or previous fees. Payment items sum to the amount used in the QR.
Avoid: item when the context could mean a QR result.

**QR item**:
The generated output for a payment row. It contains the cleaned payment row, bank display name, VietQR payload, PNG data URL, and validation errors.

**VietQR payload**:
The EMV/NAPAS TLV string that banking apps scan. It is not the PNG image.
Avoid: QR image when discussing the payload string.

**Bill number**:
The payment reference stored in VietQR Additional Data `62-01`. If the user omits it, the app derives a stable `SUN...` value from the payment row.

**Payment purpose**:
The note stored in VietQR Additional Data `62-08`. The app defaults it from the student name when absent.

**Email config**:
Local provider settings used for preview and send. Secrets are stored only in local ignored JSON files and exposed through masked public responses.

**Email provider**:
The delivery backend. Current providers are Gmail SMTP and Resend.

**Email cron**:
The local queue and scheduler for sending payment emails gradually. It respects a rolling 24-hour quota.

**Public base URL**:
The externally reachable app URL embedded in email QR links. If missing, request host or localhost fallback is used.

**Student code**:
The durable production identifier for a student. Student names are display data only and are not unique.

**School year**:
The academic year used to group classes and student listings.

**Class master data**:
The production class record for a school year and grade. Payment rows may still carry a display class name during import or preview.

**Parent contact**:
A parent or guardian linked to one or more students. Billing-email delivery uses explicit active, primary, and receives-billing flags.

**Fee type**:
A reusable production fee category such as tuition, lunch, shuttle, uniform, insurance, materials, previous fees, or custom fees. Fee types carry Vietnamese and English labels.

**Bảng phí theo kỳ**:
The production fee setup for a school year, grade, or class in a specific collection period or month. It is previewed before invoice generation and is not the same as the temporary payment-row fee template.

**Fee schedule item**:
One default amount inside a bảng phí theo kỳ, tied to a fee type and carrying Vietnamese and English labels.

**Student fee adjustment**:
A per-student discount, surcharge, waiver, or carry-over applied on top of a bảng phí theo kỳ. Every adjustment needs a reason for auditability.

**Invoice**:
The production payment request generated from a fee schedule for one student. It snapshots student, class, period, fee items, adjustments, total, bank account, QR bill number, and status.

**Invoice code**:
The stable production invoice reference. It maps directly to VietQR Bill Number `62-01` so reconciliation can match payments by invoice.

**PDF receipt**:
A generated PDF document rendered from invoice data. It includes school, student, class, period, invoice items, total, payment status, issue timestamp, and VietQR payment QR.

**Payment provider**:
The production adapter that creates payment intents or receives transaction events. Current baseline providers are `manual_vietqr`, `sepay`, and `payos`.

**Payment intent**:
The provider-specific payment request for one invoice. For manual VietQR it stores the generated QR payload; for payOS it stores the provider reference and payment URL.

**Provider event**:
The raw webhook payload received from a payment provider. It is stored before normalized transaction parsing so retries and provider data remain auditable.

**Payment transaction**:
A ledger entry for money movement from a provider webhook or manual cash receipt. It is not deleted; matching or reversal records describe how it affects an invoice.

**Reconciliation match**:
The auditable link between a payment transaction and an invoice. It records match type, score, applied amount, and reason.

**Manual cash receipt**:
A staff-entered cash collection record for an invoice. It creates a payment transaction and a reconciliation match with collector, timestamp, amount, receipt reference, and audit reason.

**Notification template**:
A versioned invoice-email template used by notification campaigns. Baseline templates are first payment notice and payment reminder.

**Notification campaign**:
An invoice-based email batch targeting invoices by school year, period, class, invoice status, and due date. It replaces row-based batches for production billing notices.

**Notification recipient**:
One parent billing contact selected for one invoice inside a notification campaign. Recipients come from active parent links that receive billing email.

**Notification log**:
The auditable delivery record for a campaign/template/invoice/recipient. It stores provider, status, provider message ID, error, dry-run flag, and sent timestamp.

**Admin dashboard**:
The production overview built from invoices and payment transactions. It reports receivable, collected, outstanding, collection rate, unpaid/partial/review counts, and top classes by outstanding amount.

**Admin report**:
A filterable production report over durable invoice and payment data, grouped by class or listed by invoice for accounting review and CSV export preparation.

**Report export**:
A CSV output generated from the current admin report filter. Current export datasets are class summaries, invoice detail, and payment transactions.

**Audit log**:
An immutable append-only record in `audit_logs` for money and fee changes. It stores actor context when available, action, entity, reason, and metadata.

**Operation log**:
A production failure log in `operation_logs` for webhook, email, and background-job issues. It is used for incident review and does not replace the immutable money/fee audit log.

**App user**:
An administrative user stored in `app_users`. This is separate from students and parents.

**App role**:
A named set of app permissions assigned to app users through `app_user_roles`.

**App permission**:
A code seeded in `app_permissions` and attached to roles. Admin write APIs declare required permission codes through their request contract.

## Relationships

- A payment row can have many payment items.
- Payment items, when present, determine the QR amount.
- A payment row generates one QR item.
- A QR item can render into an email preview or email send payload.
- Manual email sends and cron sends share the same rolling quota.
- A student can have many parent contacts.
- A parent contact can be linked to many students.
- A class belongs to one school year and can have many students.
- A bảng phí theo kỳ has many fee schedule items.
- A student can have many fee adjustments for a bảng phí theo kỳ.
- A bảng phí theo kỳ can generate one active invoice per student.
- An invoice has many invoice items and optional invoice adjustments.
- An invoice generates VietQR payment data using invoice code as the bill number.
- An invoice can have payment intents through multiple providers.
- A provider event can create at most one idempotent payment transaction per provider transaction reference.
- A payment transaction can match one invoice through reconciliation.
- Manual cash receipts create payment transactions, reconciliation matches, and audit logs.
- A notification campaign targets many invoices and many notification recipients.
- A notification recipient belongs to one invoice and one parent billing contact.
- A notification log records one dry-run, send, skip, or error outcome for one campaign/template/invoice/recipient.
- An admin dashboard summarizes many invoices and payment transactions.
- An admin report can export class, invoice, or payment transaction CSV.
- An operation log can reference a provider event, notification recipient, or background job result.
- An app user can have many app roles.
- An app role can have many app permissions.

## Flagged Ambiguities

- "QR" can mean payload, PNG image, or email link. Use the precise term above.
- "Amount" can mean raw CSV/JSON amount or payment-item total. In this app, payment-item total wins when payment items exist.
- "Payment" can mean a payment intent, a provider transaction, a reconciliation match, or an invoice status. Use the precise term above.
- "Notification" can mean a template, campaign, recipient, or delivery log. Use the precise term above.
- "Log" can mean audit log, operation log, notification log, provider event, or local cron state. Use the precise term above.
- "User" can mean an app user, student, parent contact, or operator. Use the precise term above.
