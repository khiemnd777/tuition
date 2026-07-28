# DEKISUGI QR Tool

## Purpose

`qr-tool/` is the official DEKISUGI product. It is a standalone static web app for turning user-owned spreadsheet rows into VietQR PNG files and email drafts. It deliberately does not become a source of truth: the browser session owns all imported rows, mappings, QR images, templates, and generated exports.

The former Go Finance Hub—including its API, Web Admin, PostgreSQL migrations, tenant/subscription modules, email sender, and cron scheduler—is obsolete and planned for complete removal. New product work defaults to `qr-tool/`; the legacy stack is not required to run the official app and must not be extended unless explicitly requested.

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
2. Review automatic header aliases and map each source column to a grouped field under `Student & school`, `Parent`, `Payment`, or `Fees / advanced`, or choose `Ignore`.
3. Optionally apply a default bank BIN and account number when those columns are not present per row.
4. Map multiple arbitrary columns to `Fee item`; each source header becomes the fee label.
5. Generate and review rows. Errors remain visible and are never silently dropped.
6. Export valid QR images with a manifest and a separate error report.
7. Select a row to preview a custom email template, copy content/QR, or download an unsent EML draft.
8. Prefer exporting one versioned Gmail JSON data file and importing it into a copied, user-owned Google Sheet template; keep the manual Gmail Free ZIP and portable provider/email bundle as fallbacks.
9. Optionally open the coffee-support dialog and generate a local VietQR for VPBank without sending app or payment data to a DEKISUGI backend.

## Supported Payment Fields

- `student_name`
- `student_code`
- `school_name`
- `cohort`
- `year`
- `parent_name`
- `class_name`
- `bank_bin`
- `bank_account`
- `email`
- `amount`
- `payment_items`
- `bill_number`
- `note`
- repeatable `fee_item`, whose label is derived from the source spreadsheet header

Legacy named fee headers remain compatible: automatic mapping routes them to repeatable `fee_item` instead of presenting each old fee preset as a separate system field.

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

Templates contain a plain subject plus sanitized email-safe HTML. Supported merge fields include recipient, student code/name, school, cohort, year, class, parent, amount/payment items, bank/account, bill number, payment note, and the per-row QR image.

Generated output never sends email:

- Rich clipboard payload contains `text/html` and `text/plain`; QR can also be copied independently.
- EML output uses `X-Unsent: 1`, multipart alternative text/HTML, and a matching inline QR Content-ID.
- Provider bundle includes EML, HTML, text, QR, CSV, JSONL, manifest, template, and human-readable instructions.
- The recommended Gmail path exports `DEKISUGI_GMAIL_DATA_YYYY-MM-DD.json` with schema kind `dekisugi.gmail-data`, version `1`, pre-rendered email bodies, per-row QR Base64, status, and validation errors. The static utility does not upload this file.
- After one-time setup, the browser may persist only the normalized personal Google Sheet `/edit` URL. No payment row, recipient, QR, email body, template, or Google credential is persisted. Disconnecting removes this URL without deleting the user's Sheet.
- A separately published Google Sheet template owns `Code.gs`, `Sidebar.html`, and an explicit Apps Script manifest. The user clicks the configured `/copy` URL, becomes owner of the copy, and selects the JSON file in the sidebar.
- The sidebar uses a custom confirmation surface, requires a successful self-addressed test before real sending or scheduling, caps each run at 90 recipients, preserves ten recipients of reported quota, locks concurrent batches, and never automatically retries a row left at `SENDING`.
- Daily scheduling is opt-in. Its time trigger runs in the user's copied Sheet and Google account, processes only `READY` rows, skips `SENT`, and deletes itself after no pending rows remain.
- Importing a new dataset removes an existing schedule, clears test approval, validates the schema/row/cell limits, and neutralizes leading spreadsheet-formula characters before writing cells.
- Gmail Free bundle contains exactly three files: `01_DANH_SACH_GUI.xlsx`, `02_CODE_GUI_EMAIL.gs`, and `BAT_DAU_TU_DAY.html`. Its script exposes `0. Nhập dữ liệu mới`, validates the versioned JSON, neutralizes formula-leading values, replaces the current `EMAILS` dataset, and retains the installed code and Google authorization for future collection periods.
- Its workbook contains a visible guide sheet plus one `EMAILS` row per imported record, with pre-rendered subject, HTML/text body, per-recipient QR Base64, send flag, status, sent time, and error details.
- The Apps Script runs only in the user's own Google Sheet and Gmail account. It uses `MailApp`, embeds each QR through a matching inline Content-ID, requires a manual confirmation before real sending, and never asks for a Gmail password or App Password.
- The popup, workbook guide sheet, and offline HTML all use a non-technical five-stage walkthrough. They name the exact Apps Script function selector, `Run`, `Review permissions`, optional unverified-app path, consent buttons, expected `Execution completed` result, Sheet reload/menu, self-addressed test send, and final-send confirmation.
- The offline instructions explicitly convert the exported `.xlsx` to native Google Sheets before opening Apps Script, explain popup/admin-policy failure cases, and keep technical body/QR columns hidden after setup. `setup()` reports through the execution log and a non-blocking Sheet toast instead of leaving a blocking dialog in another tab.
- The script sends at most 90 recipients per manual batch/day, preserves ten recipients of reported quota, marks rows `READY`, `SENDING`, `SENT`, `ERROR`, or `SKIP`, and does not automatically resend `SENT` rows.
- The legacy ZIP script creates no time trigger. The copied-template script only creates a time trigger after an explicit sidebar confirmation. Neither script makes an external `UrlFetchApp` request. The static utility still has no backend, credentials, Google authorization, or real email-send path of its own.

`VITE_GMAIL_SHEET_TEMPLATE_URL` is compiled into the static build and must be an HTTPS `docs.google.com/spreadsheets/d/.../copy` URL. When it is absent or invalid, the first-use action downloads the reusable manual Gmail ZIP. When it is present, the first-use action opens the published template. After either setup path, the user connects their personal Sheet URL and the monthly flow always opens that saved Sheet. The one-time publisher workflow is documented in `qr-tool/google-sheet-template/SETUP.md`; the repository does not create or mutate an external Google Drive file.

Native Gmail Mail Merge is intentionally not exposed because it is unavailable to most free Gmail accounts and cannot attach a different local QR to every recipient. The provider export remains available for advanced users who need EML/CSV/JSONL integration instead.

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
