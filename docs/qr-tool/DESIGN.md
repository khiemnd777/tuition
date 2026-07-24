# DEKISUGI QR Export Utility

## Purpose

`qr-tool/` is a standalone static web utility for turning user-owned spreadsheet rows into VietQR PNG files and email drafts. It deliberately does not become a source of truth: the browser session owns all imported rows, mappings, QR images, templates, and generated exports.

The existing Go API, Web Admin, PostgreSQL migrations, tenant/subscription modules, email sender, and cron scheduler remain unchanged and are not required to run this utility.

## Runtime And Privacy Contract

- All spreadsheet parsing, mapping, VietQR generation, PNG rendering, template rendering, MIME assembly, and ZIP generation happen in the browser.
- Runtime code makes no DEKISUGI API, Gmail, provider, analytics, font, or CDN request.
- Dependencies are bundled into the production static assets.
- Payment data and templates are not written to cookies, `localStorage`, `sessionStorage`, `IndexedDB`, a service worker, or a server.
- Reloading or closing the page discards the in-memory session.
- A reusable template must be exported explicitly as JSON and imported again by the user.
- The Content Security Policy blocks third-party scripts, frames, objects, and form submissions.

## User Flow

1. Choose or drag an XLSX, XLS, or CSV file containing at most 500 payment rows.
2. Review automatic header aliases and map each source column to a supported field or `Ignore`.
3. Optionally apply a default bank BIN and account number when those columns are not present per row.
4. Map multiple arbitrary columns to `Fee item`; each source header becomes the fee label.
5. Generate and review rows. Errors remain visible and are never silently dropped.
6. Export valid QR images with a manifest and a separate error report.
7. Select a row to preview a custom email template, copy content/QR, or download an unsent EML draft.
8. Export a Gmail Mail Merge workbook or a portable provider/email bundle.

## Supported Payment Fields

- `student_name`
- `parent_name`
- `class_name`
- `bank_bin`
- `bank_account`
- `email`
- `amount`
- `payment_items`
- `bill_number`
- `note`
- the legacy named fee columns documented in the root README
- repeatable `fee_item`, whose label is derived from the source spreadsheet header

`payment_items` override raw `amount` whenever at least one payment item exists.

## VietQR Contract

The JavaScript engine mirrors `vietqr_standard.go`:

- NAPAS AID `A000000727` and account-transfer service `QRIBFTTA`.
- Static method `11`, or dynamic method `12` when an amount is present.
- Currency `704`, country `VN`.
- Bill Number at Additional Data `62-01`.
- Payment Purpose at Additional Data `62-08`.
- ASCII-normalized ANS values capped at 25 characters.
- Bank BIN must be six digits and exist in the bundled bank list.
- Account number must be one to nineteen alphanumeric characters.
- CRC-16/CCITT-FALSE, uppercase and zero-padded to four characters.

Regression tests share the exact static, dynamic, and CRC fixtures used by the Go test suite.

## Email Template And Export Contract

Templates contain a plain subject plus sanitized email-safe HTML. Supported merge fields include recipient/student/parent/class, amount/payment items, bank/account, bill number, payment note, and the per-row QR image.

Generated output never sends email:

- Rich clipboard payload contains `text/html` and `text/plain`; QR can also be copied independently.
- EML output uses `X-Unsent: 1`, multipart alternative text/HTML, and a matching inline QR Content-ID.
- Provider bundle includes EML, HTML, text, QR, CSV, JSONL, manifest, template, and human-readable instructions.
- Gmail bundle places text-only recipient fields in the first worksheet and converts supported body placeholders into Gmail merge tags.

Gmail Mail Merge cannot personalize the subject or attach a different local QR to every recipient. When a template contains `qr_image`, the Gmail-specific body replaces it with an explicit compatibility warning. EML/provider exports preserve the per-recipient QR.

## Security And Failure Behavior

- Spreadsheet and template values are HTML-escaped before merge rendering.
- Imported/editor HTML is sanitized; scripts, embedded objects, event handlers, remote style URLs, and unsafe link schemes are removed.
- Invalid JSON payment items become row errors instead of being ignored.
- Duplicate singleton field mappings are blocked; repeatable fee mappings are allowed.
- Invalid QR rows are included in `errors.csv` and excluded from PNG/EML generation.
- Rows with missing or invalid email remain in the QR export and are reported in `email-errors.csv` for email exports.
- Filenames are ASCII-normalized, length-capped, collision-safe, and do not contain path separators.

## Development Verification

```sh
cd qr-tool
npm ci
npm test
npm run build
```

The repository-level Go regression remains:

```sh
go test ./...
```
