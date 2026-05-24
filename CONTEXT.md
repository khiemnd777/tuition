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

## Flagged Ambiguities

- "QR" can mean payload, PNG image, or email link. Use the precise term above.
- "Amount" can mean raw CSV/JSON amount or payment-item total. In this app, payment-item total wins when payment items exist.
