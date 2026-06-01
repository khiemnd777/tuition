# Demo Sample Data

These files are fictional data for customer demos. They are safe for preview, QR generation, invoice dry-runs, and email dry-runs. Do not use them as real payment instructions: emails use `example.com`, and bank accounts use `DEMOACC...` placeholders.

## Files

- `demo_master_data.csv`: production master data for 14 demo students across classes `1.01`, `2.01`, `3.02`, and `3.03`.
- `demo_payments.csv`: legacy QR/import payment rows with fee columns that override the raw `amount` field.
- `demo_fee_adjustments.csv`: paste-ready student fee adjustments for the tuition setup workflow.

The smaller `students.csv` and `master_data.csv` files are still kept as compact technical samples.

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
