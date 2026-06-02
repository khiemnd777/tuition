# Demo Sample Data

These files are fictional data for customer demos. They are safe for preview, QR generation, invoice dry-runs, and email dry-runs. Do not use them as real payment instructions: emails use `example.com`, and bank accounts use `DEMOACC...` placeholders.

## Files

- `finance_hub_demo/`: end-to-end Finance Hub demo pack covering master data, fee schedules, invoices, reconciliation, notifications, and the legacy QR/import tool.
- `demo_master_data.csv`: production master data for 14 demo students across classes `1.01`, `2.01`, `3.02`, and `3.03`.
- `demo_payments.csv`: legacy QR/import payment rows with fee columns that override the raw `amount` field.
- `demo_fee_adjustments.csv`: paste-ready student fee adjustments for the tuition setup workflow.

The smaller `students.csv` and `master_data.csv` files are still kept as compact technical samples.

## Finance Hub Demo Pack

Use `finance_hub_demo/` when demoing the full Finance Hub story to customers:

1. `finance_hub_demo/master_data.csv` for `Hoc sinh & phu huynh`.
2. `finance_hub_demo/fee_schedule_profiles.csv` and `finance_hub_demo/fee_adjustments.csv` for `Bang phi`.
3. `finance_hub_demo/invoice_generation_request.json` as the invoice preview/generation API template after a fee schedule is saved.
4. `finance_hub_demo/reconciliation_webhooks.jsonl` and `finance_hub_demo/manual_cash_receipts.csv` for `Doi soat`.
5. `finance_hub_demo/notification_campaigns.json` for notification preview, reminder, and email-preview demos.
6. `finance_hub_demo/legacy_qr_payments.csv` for the secondary QR/import workflow.

See `finance_hub_demo/README.md` for the step-by-step runbook and the placeholder replacement points.

## Suggested Demo Flow

1. Start the app and open the Web Admin UI.
2. In `Hoc sinh & phu huynh`, import `samples/demo_master_data.csv`.
3. Use import preview first. The demo data is intended to be applyable, while still containing students that later show readiness warnings because they do not have an active billing recipient.
4. In `Bang phi`, create a fee schedule for school year `2025-2026`. For a quick demo, use one class such as `3.02`, period `2025-04`, month `4`, and status `active`.
5. Paste rows from `samples/demo_fee_adjustments.csv` into the fee adjustments CSV box. Keep only rows for the selected class when demoing a single-class schedule.
6. Preview the fee schedule, then open invoice preview/generation if the environment has a configured local database.
7. In `Cong cu QR/import`, import `samples/demo_payments.csv` to show the legacy QR batch table and email preview. Use preview or dry-run only unless real sending is explicitly intended.

## Demo Cases Included

- Duplicate student display name with different student codes: `DEMO-S001` and `DEMO-S011`.
- Siblings sharing one parent contact: `DEMO-S005` and `DEMO-S006`.
- Students without billing-ready recipients for readiness demo: `DEMO-S004`, `DEMO-S007`, and `DEMO-S014`.
- Payment rows with mixed fee components: tuition, shuttle, health insurance, uniform, materials, and previous fees.
- Fee adjustments covering discount, waiver, surcharge, and carry-over cases.

## Notes

- `PaymentItems` totals intentionally determine QR amounts. The `amount` column in `demo_payments.csv` is `0` so the fee columns are visibly the source of truth.
- `bill_number` values are under the VietQR bill number limit and use `DEMO2504...` references for easy reconciliation demos.
- `note` values are ASCII and short enough for VietQR purpose normalization.
