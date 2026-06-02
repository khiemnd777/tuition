# Finance Hub Demo Sample Pack

This folder contains fictional data for an end-to-end Finance Hub customer demo. All emails use `example.com`, all collection accounts are placeholders, and all webhook/payment references are demo-only.

## Files

- `master_data.csv`: importable student and parent master data across classes `1.01`, `2.01`, `3.02`, `3.03`, and `4.01`.
- `fee_schedule_profiles.csv`: fee item profiles to copy into the `Bang phi` fee item table by grade/class.
- `fee_adjustments.csv`: paste-ready student fee adjustments for discounts, waivers, surcharges, and carry-over balances.
- `legacy_qr_payments.csv`: importable legacy QR/payment rows for the secondary QR/import tool.
- `invoice_generation_request.json`: API request template for invoice preview/generation after a fee schedule is saved.
- `reconciliation_webhooks.jsonl`: SePay webhook payload templates for paid, partial, overpaid, and unmatched reconciliation cases.
- `manual_cash_receipts.csv`: cash receipt cases for manual payment collection demos.
- `notification_campaigns.json`: notification campaign preview, reminder, and email-preview request templates.

## Demo Flow

1. Import `master_data.csv` in `Hoc sinh & phu huynh`.
2. Review relationship cases: duplicate student display name (`FH-S001` and `FH-S011`), siblings sharing one parent (`FH-S005` and `FH-S006`), and missing billing recipients (`FH-S004`, `FH-S007`, `FH-S014`).
3. In `Bang phi`, choose school year `2025-2026`, period `2025-04`, month `4`, and status `active`.
4. Use `fee_schedule_profiles.csv` as the amount guide for the selected grade/class. Paste matching rows from `fee_adjustments.csv` into the adjustment CSV box.
5. Preview and save the fee schedule, then open the invoice handoff.
6. In `Hoa don`, preview and generate invoices with bank BIN `970415` and collection account `FHCOLLECT001`.
7. In `Doi soat`, create QR/payment intents from invoice rows. For webhook demos, replace each `<INVOICE_CODE_...>` placeholder in `reconciliation_webhooks.jsonl` with an invoice code shown by the app, then post to the listed path.
8. Use `manual_cash_receipts.csv` to demo cash receipt entry and manual review states.
9. In `Thong bao`, use `notification_campaigns.json` as the target guide. Always preview or dry-run before any real send.
10. In `Cong cu QR/import`, import `legacy_qr_payments.csv` to show the legacy QR batch and email preview path.

## Demo Cases Included

- Student and parent relationship readiness.
- Class/grade fee profiles and per-student adjustments.
- Invoice preview/generation and idempotency.
- VietQR-ready invoice payment references.
- Reconciliation states: paid, partial, overpaid, unmatched, and cash.
- Notification targeting for first notices and reminders.

## Safety

- Do not send real email or run real cron for this sample pack.
- Do not use the placeholder bank accounts as real payment instructions.
- Keep webhook payloads in preview/demo environments only.
